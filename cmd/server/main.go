package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/api/middleware"
	"github.com/openpayment/gateway/internal/api/monitoring"
	"github.com/openpayment/gateway/internal/config"
	"github.com/openpayment/gateway/internal/database"
	"github.com/openpayment/gateway/internal/pkg/encrypt"
	"github.com/openpayment/gateway/internal/pkg/logger"
	"github.com/openpayment/gateway/internal/service/admin"
	"github.com/openpayment/gateway/internal/service/auth"
	"github.com/openpayment/gateway/internal/service/customer"
	"github.com/openpayment/gateway/internal/service/fraud"
	"github.com/openpayment/gateway/internal/service/ledger"
	"github.com/openpayment/gateway/internal/service/merchant"
	"github.com/openpayment/gateway/internal/service/payment"
	"github.com/openpayment/gateway/internal/service/settlement"
	"github.com/openpayment/gateway/internal/service/webhook"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogLevel)

	if cfg.JWTSecret == "" {
		log.Fatal().Msg("JWT_SECRET is required - set it to a random 64-hex-char string")
	}

	os.Setenv("ENCRYPTION_KEY", cfg.EncryptionKey)
	if err := encrypt.Init(); err != nil {
		log.Fatal().Err(err).Msg("encryption initialization failed - cannot run without encryption")
	}
	log.Info().Msg("encryption initialized")

	log.Info().Str("port", cfg.Port).Msg("Starting Open Payment Gateway...")

	var db *database.PostgresDB
	var err error
	for i := 0; i < 30; i++ {
		db, err = database.NewPostgres(cfg)
		if err == nil {
			break
		}
		log.Warn().Err(err).Int("attempt", i+1).Msg("waiting for database...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database after 30 attempts")
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	authSvc := auth.NewAuthService(cfg, db.Pool)
	merchantRepo := merchant.NewRepository(db)
	merchantSvc := merchant.NewService(merchantRepo)

	webhookRepo := webhook.NewRepository(db)
	webhookSvc := webhook.NewService(webhookRepo)

	ledgerRepo := ledger.NewRepository(db)
	ledgerSvc := ledger.NewService(ledgerRepo)

	paymentRepo := payment.NewRepository(db)
	paymentSvc := payment.NewService(paymentRepo, payment.NewProcessorClient("http://mock-processor:9000")).
		WithWebhook(webhookSvc).
		WithLedger(ledgerSvc)

	hc := &api.HealthChecker{
		DB:      db.Pool,
		Uptime:  time.Now(),
		Version: "1.0.0",
	}

	router := api.NewRouter(cfg, hc)

	authRateLimiter := middleware.NewRateLimiter(10, 20, time.Second)

	v1 := hc.V1
	auth.RegisterAuthRoutes(v1, authSvc, authRateLimiter)
	merchant.RegisterMerchantRoutes(v1, merchantSvc, auth.AuthMiddleware(authSvc))
	payment.RegisterPaymentRoutes(v1, paymentSvc, auth.AuthMiddleware(authSvc))

	customerRepo := customer.NewRepository(db)
	customerSvc := customer.NewService(customerRepo)
	customer.RegisterCustomerRoutes(v1, customerSvc, auth.AuthMiddleware(authSvc))

	webhook.RegisterWebhookRoutes(v1, webhookSvc, auth.AuthMiddleware(authSvc))
	ledger.RegisterLedgerRoutes(v1, ledgerSvc, auth.AuthMiddleware(authSvc))

	fraudRepo := fraud.NewRepository(db)
	fraudSvc := fraud.NewService(fraudRepo)
	fraud.RegisterFraudRoutes(v1, fraudSvc, auth.AuthMiddleware(authSvc))

	adminRepo := admin.NewRepository(db)
	adminSvc := admin.NewService(adminRepo)

	statsRepo := admin.NewStatsRepository(db)
	statsSvc := admin.NewStatsService(statsRepo)

	slaSvc := monitoring.NewSLAService(db.Pool)

	admin.RegisterAdminRoutes(v1, adminSvc, statsSvc, merchantSvc, auth.AuthMiddleware(authSvc))
	monitoring.RegisterSLARoutes(v1, slaSvc, auth.AuthMiddleware(authSvc))

	settlementRepo := settlement.NewRepository(db)
	settlementSvc := settlement.NewService(settlementRepo, ledgerSvc)
	settlement.RegisterSettlementRoutes(v1, settlementSvc, auth.AuthMiddleware(authSvc))

	webhookWorker := webhook.NewRetryWorker(webhookSvc, 60*time.Second)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	go webhookWorker.Start(workerCtx)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	workerCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}

	db.Close()
	log.Info().Msg("server stopped gracefully")
}
