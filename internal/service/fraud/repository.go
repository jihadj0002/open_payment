package fraud

import (
	"context"
	"errors"
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

func (r *Repository) SaveFraudCheck(ctx context.Context, fc *FraudCheck) error {
	fc.ID = uuid.New().String()
	fc.CheckedAt = time.Now()

	query := `INSERT INTO fraud_checks (id, payment_intent_id, merchant_id, score, threshold, verdict, flags, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.Pool.Exec(ctx, query,
		fc.ID, fc.PaymentIntentID, fc.MerchantID, fc.Score, fc.Threshold, fc.Verdict, fc.Flags, fc.CheckedAt,
	)
	return err
}

func (r *Repository) GetRecentTransactionsByCustomer(ctx context.Context, customerID, merchantID string, withinMinutes int) (int, error) {
	query := `SELECT COUNT(*) FROM payment_intents
		WHERE customer_id = $1 AND merchant_id = $2 AND created_at >= NOW() - ($3 || ' minutes')::INTERVAL`

	var count int
	err := r.Pool.QueryRow(ctx, query, customerID, merchantID, withinMinutes).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) GetFraudConfig(ctx context.Context, merchantID string) (*FraudConfig, error) {
	query := `SELECT merchant_id, max_amount, block_vpn, max_ip_count, enabled
		FROM fraud_configs WHERE merchant_id = $1`

	var cfg FraudConfig
	err := r.Pool.QueryRow(ctx, query, merchantID).Scan(
		&cfg.MerchantID, &cfg.MaxAmount, &cfg.BlockVPN, &cfg.MaxIpCount, &cfg.Enabled,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &FraudConfig{
				MerchantID: merchantID,
				MaxAmount:  1000000,
				BlockVPN:   false,
				MaxIpCount: 10,
				Enabled:    true,
			}, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *Repository) UpsertFraudConfig(ctx context.Context, cfg *FraudConfig) error {
	query := `INSERT INTO fraud_configs (merchant_id, max_amount, block_vpn, max_ip_count, enabled, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (merchant_id) DO UPDATE SET
			max_amount = EXCLUDED.max_amount,
			block_vpn = EXCLUDED.block_vpn,
			max_ip_count = EXCLUDED.max_ip_count,
			enabled = EXCLUDED.enabled,
			updated_at = NOW()`

	_, err := r.Pool.Exec(ctx, query,
		cfg.MerchantID, cfg.MaxAmount, cfg.BlockVPN, cfg.MaxIpCount, cfg.Enabled,
	)
	return err
}
