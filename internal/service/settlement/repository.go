package settlement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/openpayment/gateway/internal/database"
)

type Repository struct {
	*database.BaseRepository
}

func NewRepository(db *database.PostgresDB) *Repository {
	return &Repository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *Repository) CreateSettlement(ctx context.Context, s *Settlement) error {
	s.ID = uuid.New().String()
	s.CreatedAt = time.Now()

	query := `INSERT INTO settlements (id, merchant_id, amount, currency, status, fee, net_amount, period_start, period_end, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.Pool.Exec(ctx, query,
		s.ID, s.MerchantID, s.Amount, s.Currency, s.Status, s.Fee, s.NetAmount,
		s.PeriodStart, s.PeriodEnd, s.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert settlement: %w", err)
	}
	return nil
}

func (r *Repository) GetSettlement(ctx context.Context, id, merchantID string) (*Settlement, error) {
	query := `SELECT id, merchant_id, amount, currency, status, fee, net_amount, payout_ref, period_start, period_end, completed_at, created_at
		FROM settlements WHERE id = $1 AND merchant_id = $2`

	var s Settlement
	err := r.Pool.QueryRow(ctx, query, id, merchantID).Scan(
		&s.ID, &s.MerchantID, &s.Amount, &s.Currency, &s.Status, &s.Fee, &s.NetAmount,
		&s.PayoutRef, &s.PeriodStart, &s.PeriodEnd, &s.CompletedAt, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get settlement: %w", err)
	}
	return &s, nil
}

func (r *Repository) ListSettlements(ctx context.Context, merchantID string, limit, offset int) ([]Settlement, error) {
	query := `SELECT id, merchant_id, amount, currency, status, fee, net_amount, payout_ref, period_start, period_end, completed_at, created_at
		FROM settlements WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, merchantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list settlements: %w", err)
	}
	defer rows.Close()

	var settlements []Settlement
	for rows.Next() {
		var s Settlement
		if err := rows.Scan(
			&s.ID, &s.MerchantID, &s.Amount, &s.Currency, &s.Status, &s.Fee, &s.NetAmount,
			&s.PayoutRef, &s.PeriodStart, &s.PeriodEnd, &s.CompletedAt, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan settlement: %w", err)
		}
		settlements = append(settlements, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if settlements == nil {
		settlements = []Settlement{}
	}
	return settlements, nil
}

func (r *Repository) GetUnsettledVolume(ctx context.Context, merchantID, currency string) (int64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) - COALESCE((
		SELECT COALESCE(SUM(amount), 0) FROM settlements
		WHERE merchant_id = $1 AND currency = $2 AND status = 'completed'
	), 0) FROM payment_intents
	WHERE merchant_id = $1 AND currency = $2 AND status IN ('succeeded', 'captured')`

	var volume int64
	err := r.Pool.QueryRow(ctx, query, merchantID, currency).Scan(&volume)
	if err != nil {
		return 0, fmt.Errorf("get unsettled volume: %w", err)
	}
	return volume, nil
}

func (r *Repository) UpdateSettlementStatus(ctx context.Context, id, status string) error {
	now := time.Now()
	query := `UPDATE settlements SET status = $1, completed_at = $2 WHERE id = $3`
	_, err := r.Pool.Exec(ctx, query, status, now, id)
	if err != nil {
		return fmt.Errorf("update settlement status: %w", err)
	}
	return nil
}
