package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/config"
	"github.com/openpayment/gateway/internal/database"
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

	log.Info().Str("port", cfg.Port).Msg("Starting Open Payment Gateway...")

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := database.RunMigrations(db, "internal/database/migrations"); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	authSvc := auth.NewAuthService(cfg, db.Pool)
	merchantRepo := merchant.NewRepository(db)
	merchantSvc := merchant.NewService(merchantRepo)
	paymentRepo := payment.NewRepository(db)
	paymentSvc := payment.NewService(paymentRepo, payment.NewProcessorClient("http://mock-processor:9000"))

	router := api.NewRouter(cfg)
	auth.RegisterAuthRoutes(router, authSvc)
	merchant.RegisterMerchantRoutes(router, merchantSvc, auth.AuthMiddleware(authSvc))
	payment.RegisterPaymentRoutes(router, paymentSvc, auth.AuthMiddleware(authSvc))

	customerRepo := customer.NewRepository(db)
	customerSvc := customer.NewService(customerRepo)
	customer.RegisterCustomerRoutes(router, customerSvc, auth.AuthMiddleware(authSvc))

	webhookRepo := webhook.NewRepository(db)
	webhookSvc := webhook.NewService(webhookRepo)
	webhook.RegisterWebhookRoutes(router, webhookSvc, auth.AuthMiddleware(authSvc))

	ledgerRepo := ledger.NewRepository(db)
	ledgerSvc := ledger.NewService(ledgerRepo)
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
}
