package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"

	"github.com/openpayment/gateway/internal/api/middleware"
	"github.com/openpayment/gateway/internal/config"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*", "https://openpaymentweb-production.up.railway.app"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Idempotency-Key", "X-Signature", "X-Timestamp", "X-Nonce"},
		ExposedHeaders:   []string{"X-Request-Id", "X-RateLimit-Remaining"},
		MaxAge:           3600,
	}))
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", healthHandler)

	log.Info().Str("env", cfg.Environment).Msg("router initialized")
	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "open-payment-gateway",
	})
}
