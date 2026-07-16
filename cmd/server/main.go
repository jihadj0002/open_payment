package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/api/middleware"
	"github.com/openpayment/gateway/internal/api/monitoring"
	"github.com/openpayment/gateway/internal/config"
	"github.com/openpayment/gateway/internal/database"
	"github.com/openpayment/gateway/internal/pkg/encrypt"
	"github.com/openpayment/gateway/internal/pkg/logger"
	"github.com/openpayment/gateway/internal/processor/bkash"
	"github.com/openpayment/gateway/internal/processor/nagad"
	"github.com/openpayment/gateway/internal/processor/wallet"
	"github.com/openpayment/gateway/internal/service/admin"
	"github.com/openpayment/gateway/internal/service/auth"
	"github.com/openpayment/gateway/internal/service/customer"
	"github.com/openpayment/gateway/internal/service/fraud"
	"github.com/openpayment/gateway/internal/service/ledger"
	"github.com/openpayment/gateway/internal/service/merchant"
	"github.com/openpayment/gateway/internal/service/notification"
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

	adminRepo := admin.NewRepository(db)
	adminSvc := admin.NewService(adminRepo)

	emailSender := notification.NewEmailSender(cfg.SMTP, cfg.FrontendURL)
	authSvc := auth.NewAuthService(cfg, db.Pool).WithAuditor(adminSvc).WithEmailSender(emailSender)
	merchantRepo := merchant.NewRepository(db)
	merchantSvc := merchant.NewService(merchantRepo)

	webhookRepo := webhook.NewRepository(db)
	webhookSvc := webhook.NewService(webhookRepo)

	ledgerRepo := ledger.NewRepository(db)
	ledgerSvc := ledger.NewService(ledgerRepo)

	paymentRepo := payment.NewRepository(db)
	paymentSvc := payment.NewService(paymentRepo, payment.NewProcessorClient("http://mock-processor:9000")).
		WithWebhook(webhookSvc).
		WithLedger(ledgerSvc).
		WithAuditor(adminSvc)

	bkashAppKey := os.Getenv("BKASH_APP_KEY")
	bkashAppSecret := os.Getenv("BKASH_APP_SECRET")
	bkashUsername := os.Getenv("BKASH_USERNAME")
	bkashPassword := os.Getenv("BKASH_PASSWORD")
	bkashBaseURL := os.Getenv("BKASH_BASE_URL")
	if bkashBaseURL == "" {
		bkashBaseURL = "https://checkout.sandbox.bka.sh/v1.2.0-beta"
	}

	nagadMerchantID := os.Getenv("NAGAD_MERCHANT_ID")
	nagadMerchantPrivateKey := os.Getenv("NAGAD_MERCHANT_PRIVATE_KEY")
	nagadPGPublicKey := os.Getenv("NAGAD_PG_PUBLIC_KEY")
	nagadBaseURL := os.Getenv("NAGAD_BASE_URL")
	if nagadBaseURL == "" {
		nagadBaseURL = "https://sandbox.nagad.com"
	}

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:8080"
	}

	if bkashAppKey != "" && bkashAppSecret != "" {
		tokenManager := bkash.NewTokenManager(bkashBaseURL, bkashAppKey, bkashAppSecret, bkashUsername, bkashPassword)
		if err := tokenManager.Start(); err != nil {
			log.Warn().Err(err).Msg("failed to start bKash token manager (running without bKash)")
		} else {
			bkashAdapter := bkash.NewAdapter(tokenManager, bkashAppKey, bkashBaseURL, publicURL)
			bkashWalletProvider := payment.NewBkashWalletProvider(bkashAdapter)
			paymentSvc.WithWalletProvider(bkashWalletProvider)
			log.Info().Msg("bKash wallet provider initialized")
		}
	} else {
		log.Warn().Msg("BKASH_APP_KEY and BKASH_APP_SECRET not set - bKash integration disabled")
	}

	if nagadMerchantID != "" && nagadMerchantPrivateKey != "" {
		nagadCreds := nagad.Credentials{
			MerchantID:         nagadMerchantID,
			MerchantPrivateKey: nagadMerchantPrivateKey,
			PGPublicKey:        nagadPGPublicKey,
			BaseURL:            nagadBaseURL,
		}
		nagadAdapter, err := nagad.NewAdapter(nagadCreds, publicURL)
		if err != nil {
			log.Warn().Err(err).Msg("failed to create Nagad adapter (running without Nagad)")
		} else {
			nagadWalletProvider := payment.NewNagadWalletProvider(nagadAdapter)
			paymentSvc.WithWalletProvider(nagadWalletProvider)
			log.Info().Msg("Nagad wallet provider initialized")
		}
	} else {
		log.Warn().Msg("NAGAD_MERCHANT_ID and NAGAD_MERCHANT_PRIVATE_KEY not set - Nagad integration disabled")
	}

	if os.Getenv("MOCK_WALLET_ENABLED") == "true" {
		mockAdapter := wallet.NewMockAdapter(publicURL)
		mockWalletProvider := payment.NewMockWalletProvider(mockAdapter)
		paymentSvc.WithWalletProvider(mockWalletProvider)
		log.Info().Msg("mock wallet provider initialized")
	}

	statsRepo := admin.NewStatsRepository(db)
	statsSvc := admin.NewStatsService(statsRepo)

	apiUsageMW := middleware.NewAPIUsageMiddleware(statsSvc, func(ctx context.Context) string {
		claims := auth.GetClaims(ctx)
		if claims != nil {
			return claims.MerchantID
		}
		return ""
	})

	var redisClient *redis.Client
	if cfg.RedisURL != "" {
		opts, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			log.Warn().Err(err).Msg("invalid REDIS_URL, falling back to in-memory rate limiter")
		} else {
			redisClient = redis.NewClient(opts)
		}
	}

	hc := &api.HealthChecker{
		DB:          db.Pool,
		RedisClient: redisClient,
		Uptime:      time.Now(),
		Version:     "1.0.0",
		APIUsageMW:  apiUsageMW,
	}

	router := api.NewRouter(cfg, hc)

	var authRateLimiter middleware.Limiter
	if redisClient != nil {
		authRateLimiter = middleware.NewRedisRateLimiter(redisClient, 20, time.Second)
	} else {
		authRateLimiter = middleware.NewMemoryRateLimiter(10, 20, time.Second)
	}

	v1 := hc.V1
	auth.RegisterAuthRoutes(v1, authSvc, authRateLimiter)
	merchant.RegisterMerchantRoutes(v1, merchantSvc, auth.AuthMiddleware(authSvc))
	payment.RegisterPaymentRoutes(v1, paymentSvc, auth.AuthMiddleware(authSvc))

	payment.RegisterCheckoutRoutes(v1, paymentSvc)
	payment.RegisterBkashRoutes(v1, paymentSvc)
	payment.RegisterNagadRoutes(v1, paymentSvc)

	customerRepo := customer.NewRepository(db)
	customerSvc := customer.NewService(customerRepo)
	customer.RegisterCustomerRoutes(v1, customerSvc, auth.AuthMiddleware(authSvc))

	webhook.RegisterWebhookRoutes(v1, webhookSvc, auth.AuthMiddleware(authSvc))
	ledger.RegisterLedgerRoutes(v1, ledgerSvc, auth.AuthMiddleware(authSvc))

	fraudRepo := fraud.NewRepository(db)
	fraudSvc := fraud.NewService(fraudRepo)
	fraud.RegisterFraudRoutes(v1, fraudSvc, auth.AuthMiddleware(authSvc))

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
		if cfg.TLSCert != "" && cfg.TLSKey != "" {
			log.Info().Str("addr", srv.Addr).Msg("server listening with TLS")
			if err := srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("server failed to start with TLS")
			}
		} else {
			log.Info().Str("addr", srv.Addr).Msg("server listening (without TLS, expected behind TLS-terminating proxy)")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("server failed to start")
			}
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

	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close Redis client")
		} else {
			log.Info().Msg("Redis client closed")
		}
	}

	db.Close()
	log.Info().Msg("server stopped gracefully")
}
