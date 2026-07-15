package monitoring

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openpayment/gateway/internal/api"
)

type SLAStats struct {
	Availability    float64          `json:"availability"`
	ErrorRate       float64          `json:"error_rate"`
	LatencyP50      float64          `json:"latency_p50_ms"`
	LatencyP95      float64          `json:"latency_p95_ms"`
	LatencyP99      float64          `json:"latency_p99_ms"`
	TotalRequests   int64            `json:"total_requests"`
	TotalErrors     int64            `json:"total_errors"`
	Breakdown2xx    int64            `json:"breakdown_2xx"`
	Breakdown4xx    int64            `json:"breakdown_4xx"`
	Breakdown5xx    int64            `json:"breakdown_5xx"`
	Period          string           `json:"period"`
	Healthy         bool             `json:"healthy"`
}

type SLAService struct {
	db *pgxpool.Pool
}

func NewSLAService(db *pgxpool.Pool) *SLAService {
	return &SLAService{db: db}
}

func (s *SLAService) GetSLAStats(ctx context.Context, period string) (*SLAStats, error) {
	hours := int64(1)
	switch period {
	case "24h":
		hours = 24
	case "7d":
		hours = 168
	case "30d":
		hours = 720
	}

	stats := &SLAStats{
		Period: period,
	}

	err := s.db.QueryRow(ctx, fmt.Sprintf(`
		SELECT
			COUNT(*)::BIGINT,
			COUNT(*) FILTER (WHERE status_code >= 500)::BIGINT,
			COUNT(*) FILTER (WHERE status_code >= 400 AND status_code < 500)::BIGINT,
			COUNT(*) FILTER (WHERE status_code >= 200 AND status_code < 300)::BIGINT,
			COALESCE(PERCENTILE_CONT(0.50) WITHIN GROUP (ORDER BY duration_ms), 0),
			COALESCE(PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY duration_ms), 0),
			COALESCE(PERCENTILE_CONT(0.99) WITHIN GROUP (ORDER BY duration_ms), 0)
		FROM api_usage_logs
		WHERE created_at >= NOW() - $1::INTERVAL
	`), fmt.Sprintf("%d hours", hours),
	).Scan(
		&stats.TotalRequests, &stats.Breakdown5xx, &stats.Breakdown4xx,
		&stats.Breakdown2xx, &stats.LatencyP50, &stats.LatencyP95, &stats.LatencyP99,
	)
	if err != nil {
		return nil, fmt.Errorf("query sla stats: %w", err)
	}

	stats.TotalErrors = stats.Breakdown4xx + stats.Breakdown5xx

	if stats.TotalRequests > 0 {
		stats.ErrorRate = float64(stats.TotalErrors) / float64(stats.TotalRequests) * 100
		stats.Availability = float64(stats.TotalRequests-stats.TotalErrors) / float64(stats.TotalRequests) * 100
	} else {
		stats.Availability = 100
	}

	stats.Healthy = stats.ErrorRate < 1.0 && stats.LatencyP99 < 2000 && stats.Availability >= 99.9

	return stats, nil
}

func RegisterSLARoutes(r chi.Router, svc *SLAService, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Get("/admin/sla", HandleGetSLA(svc))
}

func HandleGetSLA(svc *SLAService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if period == "" {
			period = "1h"
		}

		stats, err := svc.GetSLAStats(r.Context(), period)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "failed to get SLA stats")
			return
		}

		api.RespondJSON(w, http.StatusOK, stats)
	}
}
