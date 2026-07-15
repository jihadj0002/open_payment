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

	if cfg.EncryptionKey != "" {
		os.Setenv("ENCRYPTION_KEY", cfg.EncryptionKey)
		if err := encrypt.Init(); err != nil {
			log.Warn().Err(err).Msg("encryption initialization failed - sensitive data will not be encrypted")
		} else {
			log.Info().Msg("encryption initialized")
		}
	} else {
		log.Warn().Msg("ENCRYPTION_KEY not set - sensitive data will not be encrypted")
	}

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

	router := api.NewRouter(cfg, &api.HealthChecker{
		DB:      db.Pool,
		Uptime:  time.Now(),
		Version: "1.0.0",
	})
	auth.RegisterAuthRoutes(router, authSvc)
	merchant.RegisterMerchantRoutes(router, merchantSvc, auth.AuthMiddleware(authSvc))
	payment.RegisterPaymentRoutes(router, paymentSvc, auth.AuthMiddleware(authSvc))

	customerRepo := customer.NewRepository(db)
	customerSvc := customer.NewService(customerRepo)
	customer.RegisterCustomerRoutes(router, customerSvc, auth.AuthMiddleware(authSvc))

	webhook.RegisterWebhookRoutes(router, webhookSvc, auth.AuthMiddleware(authSvc))
	ledger.RegisterLedgerRoutes(router, ledgerSvc, auth.AuthMiddleware(authSvc))

	fraudRepo := fraud.NewRepository(db)
	fraudSvc := fraud.NewService(fraudRepo)
	fraud.RegisterFraudRoutes(router, fraudSvc, auth.AuthMiddleware(authSvc))

	adminRepo := admin.NewRepository(db)
	adminSvc := admin.NewService(adminRepo)
	admin.RegisterAdminRoutes(router, adminSvc, auth.AuthMiddleware(authSvc))

	settlementRepo := settlement.NewRepository(db)
	settlementSvc := settlement.NewService(settlementRepo, ledgerSvc)
	settlement.RegisterSettlementRoutes(router, settlementSvc, auth.AuthMiddleware(authSvc))

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
