package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"

	"github.com/openpayment/gateway/internal/api/middleware"
	"github.com/openpayment/gateway/internal/config"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	allowedOrigins := getAllowedOrigins(cfg.Environment)

	rateLimiter := middleware.NewRateLimiter(100, 200, time.Second)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Idempotency-Key", "X-Signature", "X-Timestamp", "X-Nonce"},
		ExposedHeaders:   []string{"X-Request-Id", "X-RateLimit-Remaining"},
		MaxAge:           3600,
	}))
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RequestSizeLimiter(1 << 20))

	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(rateLimiter))
		r.Get("/health", healthHandler)
	})

	log.Info().Str("env", cfg.Environment).Strs("cors_origins", allowedOrigins).Msg("router initialized")
	return r
}

func getAllowedOrigins(env string) []string {
	if env == "production" {
		return []string{
			"https://openpaymentweb-production.up.railway.app",
		}
	}
	return []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"https://openpaymentweb-production.up.railway.app",
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "open-payment-gateway",
	})
}
