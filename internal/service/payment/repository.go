package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/openpayment/gateway/internal/database"
	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
)

type Repository struct {
	*database.BaseRepository
}

func NewRepository(db *database.PostgresDB) *Repository {
	return &Repository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *Repository) CreatePaymentIntent(ctx context.Context, pi *PaymentIntent) error {
	pi.ID = uuid.New().String()
	pi.ClientSecret = "pi_" + uuid.New().String()
	pi.CreatedAt = time.Now()
	pi.UpdatedAt = time.Now()

	metaBytes, err := json.Marshal(pi.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `INSERT INTO payment_intents (id, merchant_id, customer_id, amount, currency, status, idempotency_key, description, metadata, processor, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = r.Pool.Exec(ctx, query,
		pi.ID, pi.MerchantID, pi.CustomerID, pi.Amount, pi.Currency, pi.Status,
		pi.IdempotencyKey, pi.Description, metaBytes, pi.PaymentMethod,
		pi.CreatedAt, pi.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert payment intent: %w", err)
	}

	return nil
}

func (r *Repository) GetPaymentIntent(ctx context.Context, id, merchantID string) (*PaymentIntent, error) {
	query := `SELECT id, merchant_id, customer_id, amount, currency, status, idempotency_key, description, metadata, failure_reason, processor, created_at, updated_at
		FROM payment_intents WHERE id = $1 AND merchant_id = $2`

	var pi PaymentIntent
	var metaBytes []byte
	var failureReason *string

	err := r.Pool.QueryRow(ctx, query, id, merchantID).Scan(
		&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.Currency, &pi.Status,
		&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
		&pi.PaymentMethod, &pi.CreatedAt, &pi.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get payment intent: %w", err)
	}

	if metaBytes != nil {
		if err := json.Unmarshal(metaBytes, &pi.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	if failureReason != nil {
		pi.ErrorMessage = failureReason
	}

	pi.AmountCapturable = pi.Amount
	pi.CaptureMethod = "automatic"
	pi.ClientSecret = "pi_" + pi.ID

	return &pi, nil
}

func (r *Repository) UpdatePaymentIntentStatus(ctx context.Context, id, status string) error {
	return r.updatePaymentIntentStatusWithReason(ctx, id, status, "system", nil)
}

func (r *Repository) updatePaymentIntentStatusWithReason(ctx context.Context, id, status, changedBy string, reason *string) error {
	oldStatus, err := r.getPaymentIntentStatus(ctx, id)
	if err != nil {
		return err
	}

	query := `UPDATE payment_intents SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.Pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update payment intent status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}

	if err := r.recordStatusChange(ctx, id, oldStatus, status, changedBy, reason); err != nil {
		return fmt.Errorf("record status change: %w", err)
	}

	return nil
}

func (r *Repository) getPaymentIntentStatus(ctx context.Context, id string) (string, error) {
	var status string
	err := r.Pool.QueryRow(ctx, `SELECT status FROM payment_intents WHERE id = $1`, id).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", pkgErr.ErrNotFound
		}
		return "", fmt.Errorf("get current status: %w", err)
	}
	return status, nil
}

func (r *Repository) recordStatusChange(ctx context.Context, paymentIntentID, oldStatus, newStatus, changedBy string, reason *string) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO status_history (payment_intent_id, old_status, new_status, changed_by, reason) VALUES ($1, $2, $3, $4, $5)`,
		paymentIntentID, oldStatus, newStatus, changedBy, reason,
	)
	return err
}

func (r *Repository) ListStatusHistory(ctx context.Context, paymentIntentID string) ([]StatusHistoryEntry, error) {
	rows, err := r.Pool.Query(ctx,
		`SELECT id, payment_intent_id, old_status, new_status, changed_by, reason, created_at
		 FROM status_history WHERE payment_intent_id = $1 ORDER BY created_at ASC`,
		paymentIntentID,
	)
	if err != nil {
		return nil, fmt.Errorf("list status history: %w", err)
	}
	defer rows.Close()

	var history []StatusHistoryEntry
	for rows.Next() {
		var h StatusHistoryEntry
		if err := rows.Scan(&h.ID, &h.PaymentIntentID, &h.OldStatus, &h.NewStatus, &h.ChangedBy, &h.Reason, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan status history: %w", err)
		}
		history = append(history, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if history == nil {
		history = []StatusHistoryEntry{}
	}
	return history, nil
}

func (r *Repository) UpdatePaymentIntentCapture(ctx context.Context, id string, amountCapturable, amountReceived int64, status string) error {
	oldStatus, err := r.getPaymentIntentStatus(ctx, id)
	if err != nil {
		return err
	}

	query := `UPDATE payment_intents SET amount_capturable = $1, amount_received = $2, status = $3, updated_at = NOW() WHERE id = $4`
	result, err := r.Pool.Exec(ctx, query, amountCapturable, amountReceived, status, id)
	if err != nil {
		return fmt.Errorf("update payment intent capture: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}

	if err := r.recordStatusChange(ctx, id, oldStatus, status, "system", nil); err != nil {
		return fmt.Errorf("record status change: %w", err)
	}

	return nil
}

func (r *Repository) ListPaymentIntents(ctx context.Context, merchantID string, limit, offset int) ([]PaymentIntent, error) {
	query := `SELECT id, merchant_id, customer_id, amount, currency, status, idempotency_key, description, metadata, failure_reason, processor, created_at, updated_at
		FROM payment_intents WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, merchantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list payment intents: %w", err)
	}
	defer rows.Close()

	var intents []PaymentIntent
	for rows.Next() {
		var pi PaymentIntent
		var metaBytes []byte
		var failureReason *string

		err := rows.Scan(
			&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.Currency, &pi.Status,
			&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
			&pi.PaymentMethod, &pi.CreatedAt, &pi.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan payment intent: %w", err)
		}

		if metaBytes != nil {
			if err := json.Unmarshal(metaBytes, &pi.Metadata); err != nil {
				return nil, fmt.Errorf("unmarshal metadata: %w", err)
			}
		}

		if failureReason != nil {
			pi.ErrorMessage = failureReason
		}

		pi.AmountCapturable = pi.Amount
		pi.CaptureMethod = "automatic"
		pi.ClientSecret = "pi_" + pi.ID

		intents = append(intents, pi)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if intents == nil {
		intents = []PaymentIntent{}
	}

	return intents, nil
}

func (r *Repository) GetByIdempotencyKey(ctx context.Context, key, merchantID string) (*PaymentIntent, error) {
	query := `SELECT id, merchant_id, customer_id, amount, currency, status, idempotency_key, description, metadata, failure_reason, processor, created_at, updated_at
		FROM payment_intents WHERE idempotency_key = $1 AND merchant_id = $2`

	var pi PaymentIntent
	var metaBytes []byte
	var failureReason *string

	err := r.Pool.QueryRow(ctx, query, key, merchantID).Scan(
		&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.Currency, &pi.Status,
		&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
		&pi.PaymentMethod, &pi.CreatedAt, &pi.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get by idempotency key: %w", err)
	}

	if metaBytes != nil {
		if err := json.Unmarshal(metaBytes, &pi.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	if failureReason != nil {
		pi.ErrorMessage = failureReason
	}

	pi.AmountCapturable = pi.Amount
	pi.CaptureMethod = "automatic"
	pi.ClientSecret = "pi_" + pi.ID

	return &pi, nil
}

func (r *Repository) CreateTransaction(ctx context.Context, tx *Transaction) error {
	tx.ID = uuid.New().String()
	tx.CreatedAt = time.Now()

	respBytes, err := json.Marshal(tx.ProcessorResponse)
	if err != nil {
		return fmt.Errorf("marshal processor response: %w", err)
	}

	query := `INSERT INTO transactions (id, payment_intent_id, merchant_id, type, amount, currency, status, processor_ref, processor_response, fee, net_amount, idempotency_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = r.Pool.Exec(ctx, query,
		tx.ID, tx.PaymentIntentID, tx.MerchantID, tx.Type, tx.Amount, tx.Currency, tx.Status,
		tx.ProcessorRef, respBytes, tx.Fee, tx.NetAmount, tx.IdempotencyKey, tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetTransactionByIdempotencyKey(ctx context.Context, key, merchantID string) (*Transaction, error) {
	query := `SELECT id, payment_intent_id, merchant_id, type, amount, currency, status, processor_ref, processor_response, fee, net_amount, idempotency_key, created_at
		FROM transactions WHERE idempotency_key = $1 AND merchant_id = $2`

	var t Transaction
	var respBytes []byte

	err := r.Pool.QueryRow(ctx, query, key, merchantID).Scan(
		&t.ID, &t.PaymentIntentID, &t.MerchantID, &t.Type, &t.Amount, &t.Currency, &t.Status,
		&t.ProcessorRef, &respBytes, &t.Fee, &t.NetAmount, &t.IdempotencyKey, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get transaction by idempotency key: %w", err)
	}

	if respBytes != nil {
		t.ProcessorResponse = json.RawMessage(respBytes)
	}

	return &t, nil
}

func (r *Repository) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	query := `SELECT id, payment_intent_id, merchant_id, type, amount, currency, status, processor_ref, processor_response, fee, net_amount, created_at
		FROM transactions WHERE id = $1`

	var t Transaction
	var respBytes []byte

	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.PaymentIntentID, &t.MerchantID, &t.Type, &t.Amount, &t.Currency, &t.Status,
		&t.ProcessorRef, &respBytes, &t.Fee, &t.NetAmount, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get transaction: %w", err)
	}

	if respBytes != nil {
		t.ProcessorResponse = json.RawMessage(respBytes)
	}

	return &t, nil
}
