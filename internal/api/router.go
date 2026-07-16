package api

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/openpayment/gateway/internal/api/metrics"
	"github.com/openpayment/gateway/internal/api/middleware"
	"github.com/openpayment/gateway/internal/config"
)

type HealthChecker struct {
	DB           *pgxpool.Pool
	RedisClient  *redis.Client
	Procs        []string
	Uptime       time.Time
	Version      string
	V1           chi.Router
	APIUsageMW   func(http.Handler) http.Handler
}

func NewRouter(cfg *config.Config, hc *HealthChecker) *chi.Mux {
	r := chi.NewRouter()

	allowedOrigins := getAllowedOrigins()

	var rateLimiter middleware.Limiter
	if cfg.RedisURL != "" {
		opts, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			log.Warn().Err(err).Msg("invalid REDIS_URL, falling back to in-memory rate limiter")
			rateLimiter = middleware.NewMemoryRateLimiter(100, 200, time.Second)
		} else {
			rdb := redis.NewClient(opts)
			rateLimiter = middleware.NewRedisRateLimiter(rdb, 200, time.Second)
			log.Info().Msg("using Redis-backed rate limiter")
		}
	} else {
		rateLimiter = middleware.NewMemoryRateLimiter(100, 200, time.Second)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Idempotency-Key", "X-Signature", "X-Timestamp", "X-Nonce"},
		ExposedHeaders:   []string{"X-Request-Id", "X-RateLimit-Remaining", "X-API-Version"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RequestSizeLimiter(1 << 20))
	r.Use(metrics.Middleware)
	if hc.APIUsageMW != nil {
		r.Use(hc.APIUsageMW)
	}

	r.Get("/", rootInfoHandler(hc))

	r.Group(func(r chi.Router) {
		r.Handle("/metrics", metrics.Handler())
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(rateLimiter))
		r.Get("/health", healthHandler(hc))
		r.Get("/ready", readyHandler(hc))
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

func readyHandler(hc *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		checks := make(map[string]healthCheckResult)
		overallStatus := "ready"

		if hc.DB != nil {
			if err := hc.DB.Ping(ctx); err != nil {
				checks["database"] = healthCheckResult{Status: "not_ready", Error: err.Error()}
				overallStatus = "not_ready"
			} else {
				checks["database"] = healthCheckResult{Status: "ready"}
			}
		}

		if hc.RedisClient != nil {
			if err := hc.RedisClient.Ping(ctx).Err(); err != nil {
				checks["redis"] = healthCheckResult{Status: "not_ready", Error: err.Error()}
				overallStatus = "not_ready"
			} else {
				checks["redis"] = healthCheckResult{Status: "ready"}
			}
		}

		httpStatus := http.StatusOK
		if overallStatus != "ready" {
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

func getAllowedOrigins() []string {
	envOrigins := os.Getenv("CORS_ORIGINS")
	if envOrigins != "" {
		origins := strings.Split(envOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		return origins
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

		if hc.RedisClient != nil {
			if err := hc.RedisClient.Ping(ctx).Err(); err != nil {
				checks["redis"] = healthCheckResult{Status: "unhealthy", Error: err.Error()}
				overallStatus = "degraded"
			} else {
				checks["redis"] = healthCheckResult{Status: "ok"}
			}
		} else {
			checks["redis"] = healthCheckResult{Status: "not_configured"}
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
