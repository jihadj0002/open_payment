package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/openpayment/gateway/internal/api/metrics"
	"github.com/openpayment/gateway/internal/api/middleware"
	"github.com/openpayment/gateway/internal/config"
)

type HealthChecker struct {
	DB      *pgxpool.Pool
	Redis   string
	Procs   []string
	Uptime  time.Time
	Version string
	V1      chi.Router
}

func NewRouter(cfg *config.Config, hc *HealthChecker) *chi.Mux {
	r := chi.NewRouter()

	allowedOrigins := getAllowedOrigins(cfg.Environment)

	rateLimiter := middleware.NewRateLimiter(100, 200, time.Second)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Idempotency-Key", "X-Signature", "X-Timestamp", "X-Nonce"},
		ExposedHeaders:   []string{"X-Request-Id", "X-RateLimit-Remaining", "X-API-Version"},
		MaxAge:           3600,
	}))
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RequestSizeLimiter(1 << 20))
	r.Use(metrics.Middleware)

	r.Get("/", rootInfoHandler(hc))

	r.Group(func(r chi.Router) {
		r.Handle("/metrics", metrics.Handler())
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(rateLimiter))
		r.Get("/health", healthHandler(hc))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Use(apiVersionMiddleware)
		hc.V1 = r
	})

	log.Info().Str("env", cfg.Environment).Strs("cors_origins", allowedOrigins).Msg("router initialized")
	return r
}

func apiVersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", "1")
		next.ServeHTTP(w, r)
	})
}

func rootInfoHandler(hc *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, http.StatusOK, map[string]interface{}{
			"name":    "Open Payment Gateway",
			"version": "1.0.0",
			"api_version": map[string]interface{}{
				"v1": "/v1/",
			},
			"documentation": "/docs",
			"health":        "/health",
		})
	}
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

type healthCheckResult struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}

type healthResponse struct {
	Status   string                       `json:"status"`
	Uptime   string                       `json:"uptime"`
	Version  string                       `json:"version,omitempty"`
	Checks   map[string]healthCheckResult `json:"checks"`
}

func healthHandler(hc *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		checks := make(map[string]healthCheckResult)
		overallStatus := "ok"

		if hc.DB != nil {
			if err := hc.DB.Ping(ctx); err != nil {
				checks["database"] = healthCheckResult{Status: "unhealthy", Error: err.Error()}
				overallStatus = "degraded"
			} else {
				checks["database"] = healthCheckResult{Status: "ok"}
			}
		}

		if hc.Redis != "" {
			checks["redis"] = healthCheckResult{Status: "unhealthy", Error: "redis not configured"}
			overallStatus = "degraded"
		} else {
			checks["redis"] = healthCheckResult{Status: "ok"}
		}

		for _, proc := range hc.Procs {
			checks[proc] = healthCheckResult{Status: "ok"}
		}

		httpStatus := http.StatusOK
		if overallStatus != "ok" {
			httpStatus = http.StatusServiceUnavailable
		}

		RespondJSON(w, httpStatus, healthResponse{
			Status:  overallStatus,
			Uptime:  time.Since(hc.Uptime).String(),
			Version: hc.Version,
			Checks:  checks,
		})
	}
}
