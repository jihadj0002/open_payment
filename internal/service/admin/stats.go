package admin

import (
	"context"
	"time"

	"github.com/openpayment/gateway/internal/database"
)

type UsageStats struct {
	MerchantID    string  `json:"merchant_id"`
	Period        string  `json:"period"`
	RequestCount  int     `json:"request_count"`
	ErrorCount    int     `json:"error_count"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
	PaymentVolume int64   `json:"payment_volume"`
	PaymentCount  int     `json:"payment_count"`
}

type StatsRepository struct {
	*database.BaseRepository
}

func NewStatsRepository(db *database.PostgresDB) *StatsRepository {
	return &StatsRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *StatsRepository) GetUsageStats(ctx context.Context, merchantID, period string, from, to time.Time) (*UsageStats, error) {
	stats := &UsageStats{
		MerchantID: merchantID,
		Period:     period,
	}

	err := r.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) AS request_count,
			COUNT(*) FILTER (WHERE status_code >= 400) AS error_count,
			COALESCE(AVG(duration_ms), 0) AS avg_latency_ms
		FROM api_usage_logs
		WHERE merchant_id = $1
			AND created_at >= $2
			AND created_at <= $3`,
		merchantID, from, to,
	).Scan(&stats.RequestCount, &stats.ErrorCount, &stats.AvgLatencyMs)
	if err != nil {
		return stats, nil
	}

	err = r.Pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(amount), 0) AS payment_volume,
			COUNT(*) AS payment_count
		FROM payment_intents
		WHERE merchant_id = $1
			AND status = 'succeeded'
			AND created_at >= $2
			AND created_at <= $3`,
		merchantID, from, to,
	).Scan(&stats.PaymentVolume, &stats.PaymentCount)
	if err != nil {
		return stats, nil
	}

	return stats, nil
}

func (r *StatsRepository) LogAPIUsage(ctx context.Context, merchantID, endpoint, method string, statusCode int, durationMs int, ipAddress, userAgent string) error {
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO api_usage_logs (merchant_id, endpoint, method, status_code, duration_ms, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		merchantID, endpoint, method, statusCode, durationMs, ipAddress, userAgent,
	)
	return err
}

type StatsService struct {
	repo *StatsRepository
}

func NewStatsService(repo *StatsRepository) *StatsService {
	return &StatsService{repo: repo}
}

func (s *StatsService) GetUsageStats(ctx context.Context, merchantID, period string, from, to time.Time) (*UsageStats, error) {
	return s.repo.GetUsageStats(ctx, merchantID, period, from, to)
}

func (s *StatsService) LogAPIUsage(ctx context.Context, merchantID, endpoint, method string, statusCode int, durationMs int, ipAddress, userAgent string) error {
	return s.repo.LogAPIUsage(ctx, merchantID, endpoint, method, statusCode, durationMs, ipAddress, userAgent)
}
