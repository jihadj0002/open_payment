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
	"github.com/openpayment/gateway/internal/pkg/logger"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogLevel)

	log.Info().Str("port", cfg.Port).Msg("Starting Open Payment Gateway...")

	router := api.NewRouter(cfg)

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
