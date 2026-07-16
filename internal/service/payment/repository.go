package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/openpayment/gateway/internal/database"
	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
)

type execQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

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
	pi.ClientSecret = "pi_" + uuid.New().String() + "_secret_" + uuid.New().String()
	pi.CreatedAt = time.Now()
	pi.UpdatedAt = time.Now()

	metaBytes, err := json.Marshal(pi.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `INSERT INTO payment_intents 
		(id, merchant_id, customer_id, amount, amount_capturable, amount_received, capture_method, currency, status, 
		 idempotency_key, description, metadata, processor, wallet_provider,
		 return_url, cancel_url, client_secret, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	_, err = r.Pool.Exec(ctx, query,
		pi.ID, pi.MerchantID, pi.CustomerID, pi.Amount, pi.AmountCapturable, pi.AmountReceived, pi.CaptureMethod,
		pi.Currency, pi.Status,
		pi.IdempotencyKey, pi.Description, metaBytes, pi.PaymentMethod, pi.WalletProvider,
		pi.ReturnURL, pi.CancelURL, pi.ClientSecret,
		pi.CreatedAt, pi.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert payment intent: %w", err)
	}

	return nil
}

func (r *Repository) GetPaymentIntent(ctx context.Context, id, merchantID string) (*PaymentIntent, error) {
	var whereClause string
	var args []interface{}

	if merchantID != "" {
		whereClause = "WHERE pi.id = $1 AND pi.merchant_id = $2"
		args = []interface{}{id, merchantID}
	} else {
		whereClause = "WHERE pi.id = $1"
		args = []interface{}{id}
	}

	query := `SELECT pi.id, pi.merchant_id, pi.customer_id, pi.amount, pi.amount_capturable, pi.amount_received,
		pi.capture_method, pi.currency, pi.status, 
		pi.idempotency_key, pi.description, pi.metadata, pi.failure_reason, pi.processor,
		COALESCE(pi.wallet_provider, ''), pi.return_url, pi.cancel_url, pi.client_secret,
		COALESCE(pi.redirect_url, ''), COALESCE(pi.provider_ref, ''),
		pi.created_at, pi.updated_at
		FROM payment_intents pi ` + whereClause

	var pi PaymentIntent
	var metaBytes []byte
	var failureReason *string
	var redirectURL string
	var providerRef string
	var walletProvider string

	err := r.Pool.QueryRow(ctx, query, args...).Scan(
		&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.AmountCapturable, &pi.AmountReceived,
		&pi.CaptureMethod, &pi.Currency, &pi.Status,
		&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
		&pi.PaymentMethod, &walletProvider,
		&pi.ReturnURL, &pi.CancelURL, &pi.ClientSecret,
		&redirectURL, &providerRef,
		&pi.CreatedAt, &pi.UpdatedAt,
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

	if redirectURL != "" {
		pi.RedirectURL = &redirectURL
	}
	if providerRef != "" {
		pi.ProviderRef = &providerRef
	}
	pi.WalletProvider = walletProvider

	return &pi, nil
}

func (r *Repository) UpdatePaymentIntentProvider(ctx context.Context, id, providerRef, redirectURL string) error {
	return r.execUpdatePaymentIntentProvider(ctx, r.Pool, id, providerRef, redirectURL)
}

func (r *Repository) UpdatePaymentIntentProviderTx(ctx context.Context, tx pgx.Tx, id, providerRef, redirectURL string) error {
	return r.execUpdatePaymentIntentProvider(ctx, tx, id, providerRef, redirectURL)
}

func (r *Repository) execUpdatePaymentIntentProvider(ctx context.Context, q execQuerier, id, providerRef, redirectURL string) error {
	query := `UPDATE payment_intents SET provider_ref = $1, redirect_url = $2, updated_at = NOW() WHERE id = $3`
	_, err := q.Exec(ctx, query, providerRef, redirectURL, id)
	if err != nil {
		return fmt.Errorf("update payment intent provider: %w", err)
	}
	return nil
}

func (r *Repository) UpdatePaymentIntentRedirect(ctx context.Context, id, redirectURL string) error {
	return r.execUpdatePaymentIntentRedirect(ctx, r.Pool, id, redirectURL)
}

func (r *Repository) UpdatePaymentIntentRedirectTx(ctx context.Context, tx pgx.Tx, id, redirectURL string) error {
	return r.execUpdatePaymentIntentRedirect(ctx, tx, id, redirectURL)
}

func (r *Repository) execUpdatePaymentIntentRedirect(ctx context.Context, q execQuerier, id, redirectURL string) error {
	query := `UPDATE payment_intents SET redirect_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := q.Exec(ctx, query, redirectURL, id)
	if err != nil {
		return fmt.Errorf("update payment intent redirect: %w", err)
	}
	return nil
}

func (r *Repository) UpdatePaymentIntentStatus(ctx context.Context, id, status string) error {
	return r.updatePaymentIntentStatusWithReason(ctx, r.Pool, id, status, "system", nil)
}

func (r *Repository) UpdatePaymentIntentStatusTx(ctx context.Context, tx pgx.Tx, id, status string) error {
	return r.updatePaymentIntentStatusWithReason(ctx, tx, id, status, "system", nil)
}

func (r *Repository) updatePaymentIntentStatusWithReason(ctx context.Context, q execQuerier, id, status, changedBy string, reason *string) error {
	oldStatus, err := r.getPaymentIntentStatusTx(ctx, q, id)
	if err != nil {
		return err
	}

	query := `UPDATE payment_intents SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := q.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update payment intent status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}

	if err := r.recordStatusChangeTx(ctx, q, id, oldStatus, status, changedBy, reason); err != nil {
		return fmt.Errorf("record status change: %w", err)
	}

	return nil
}

func (r *Repository) getPaymentIntentStatusTx(ctx context.Context, q execQuerier, id string) (string, error) {
	var status string
	err := q.QueryRow(ctx, `SELECT status FROM payment_intents WHERE id = $1`, id).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", pkgErr.ErrNotFound
		}
		return "", fmt.Errorf("get current status: %w", err)
	}
	return status, nil
}

func (r *Repository) recordStatusChangeTx(ctx context.Context, q execQuerier, paymentIntentID, oldStatus, newStatus, changedBy string, reason *string) error {
	_, err := q.Exec(ctx,
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
	return r.updatePaymentIntentCaptureTx(ctx, r.Pool, id, amountCapturable, amountReceived, status)
}

func (r *Repository) UpdatePaymentIntentCaptureTx(ctx context.Context, tx pgx.Tx, id string, amountCapturable, amountReceived int64, status string) error {
	return r.updatePaymentIntentCaptureTx(ctx, tx, id, amountCapturable, amountReceived, status)
}

func (r *Repository) updatePaymentIntentCaptureTx(ctx context.Context, q execQuerier, id string, amountCapturable, amountReceived int64, status string) error {
	oldStatus, err := r.getPaymentIntentStatusTx(ctx, q, id)
	if err != nil {
		return err
	}

	query := `UPDATE payment_intents SET amount_capturable = $1, amount_received = $2, status = $3, updated_at = NOW() WHERE id = $4`
	result, err := q.Exec(ctx, query, amountCapturable, amountReceived, status, id)
	if err != nil {
		return fmt.Errorf("update payment intent capture: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}

	if err := r.recordStatusChangeTx(ctx, q, id, oldStatus, status, "system", nil); err != nil {
		return fmt.Errorf("record status change: %w", err)
	}

	return nil
}

func (r *Repository) ListPaymentIntents(ctx context.Context, merchantID string, limit, offset int) ([]PaymentIntent, error) {
	query := `SELECT pi.id, pi.merchant_id, pi.customer_id, pi.amount, pi.amount_capturable, pi.amount_received,
		pi.capture_method, pi.currency, pi.status,
		pi.idempotency_key, pi.description, pi.metadata, pi.failure_reason, pi.processor,
		COALESCE(pi.wallet_provider, ''), pi.return_url, pi.cancel_url, pi.client_secret,
		COALESCE(pi.redirect_url, ''), COALESCE(pi.provider_ref, ''),
		pi.created_at, pi.updated_at
		FROM payment_intents pi WHERE pi.merchant_id = $1 ORDER BY pi.created_at DESC LIMIT $2 OFFSET $3`

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
		var redirectURL string
		var providerRef string
		var walletProvider string

		err := rows.Scan(
			&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.AmountCapturable, &pi.AmountReceived,
			&pi.CaptureMethod, &pi.Currency, &pi.Status,
			&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
			&pi.PaymentMethod, &walletProvider,
			&pi.ReturnURL, &pi.CancelURL, &pi.ClientSecret,
			&redirectURL, &providerRef,
			&pi.CreatedAt, &pi.UpdatedAt,
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

		if redirectURL != "" {
			pi.RedirectURL = &redirectURL
		}
		if providerRef != "" {
			pi.ProviderRef = &providerRef
		}
		pi.WalletProvider = walletProvider

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
	query := `SELECT pi.id, pi.merchant_id, pi.customer_id, pi.amount, pi.amount_capturable, pi.amount_received,
		pi.capture_method, pi.currency, pi.status,
		pi.idempotency_key, pi.description, pi.metadata, pi.failure_reason, pi.processor,
		COALESCE(pi.wallet_provider, ''), pi.return_url, pi.cancel_url, pi.client_secret,
		COALESCE(pi.redirect_url, ''), COALESCE(pi.provider_ref, ''),
		pi.created_at, pi.updated_at
		FROM payment_intents pi WHERE pi.idempotency_key = $1 AND pi.merchant_id = $2`

	var pi PaymentIntent
	var metaBytes []byte
	var failureReason *string
	var redirectURL string
	var providerRef string
	var walletProvider string

	err := r.Pool.QueryRow(ctx, query, key, merchantID).Scan(
		&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.AmountCapturable, &pi.AmountReceived,
		&pi.CaptureMethod, &pi.Currency, &pi.Status,
		&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
		&pi.PaymentMethod, &walletProvider,
		&pi.ReturnURL, &pi.CancelURL, &pi.ClientSecret,
		&redirectURL, &providerRef,
		&pi.CreatedAt, &pi.UpdatedAt,
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

	if redirectURL != "" {
		pi.RedirectURL = &redirectURL
	}
	if providerRef != "" {
		pi.ProviderRef = &providerRef
	}
	pi.WalletProvider = walletProvider
	pi.ClientSecret = "pi_" + pi.ID

	return &pi, nil
}

func (r *Repository) CreateTransaction(ctx context.Context, txModel *Transaction) error {
	return r.createTransactionTx(ctx, r.Pool, txModel)
}

func (r *Repository) CreateTransactionTx(ctx context.Context, tx pgx.Tx, txModel *Transaction) error {
	return r.createTransactionTx(ctx, tx, txModel)
}

func (r *Repository) createTransactionTx(ctx context.Context, q execQuerier, txModel *Transaction) error {
	txModel.ID = uuid.New().String()
	txModel.CreatedAt = time.Now()

	respBytes, err := json.Marshal(txModel.ProcessorResponse)
	if err != nil {
		return fmt.Errorf("marshal processor response: %w", err)
	}

	query := `INSERT INTO transactions (id, payment_intent_id, merchant_id, type, amount, currency, status, processor_ref, processor_response, fee, net_amount, idempotency_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = q.Exec(ctx, query,
		txModel.ID, txModel.PaymentIntentID, txModel.MerchantID, txModel.Type, txModel.Amount, txModel.Currency, txModel.Status,
		txModel.ProcessorRef, respBytes, txModel.Fee, txModel.NetAmount, txModel.IdempotencyKey, txModel.CreatedAt,
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

func (r *Repository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.Pool.Begin(ctx)
}

func (r *Repository) CountPaymentIntents(ctx context.Context, merchantID string) (int, error) {
	var count int
	err := r.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM payment_intents WHERE merchant_id = $1`, merchantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count payment intents: %w", err)
	}
	return count, nil
}

func (r *Repository) GetPaymentIntentForUpdate(ctx context.Context, id, merchantID string) (*PaymentIntent, error) {
	pi, err := r.getPaymentIntentTx(ctx, r.Pool, id, merchantID, "FOR UPDATE")
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (r *Repository) GetPaymentIntentForUpdateTx(ctx context.Context, tx pgx.Tx, id, merchantID string) (*PaymentIntent, error) {
	pi, err := r.getPaymentIntentTx(ctx, tx, id, merchantID, "FOR UPDATE")
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (r *Repository) getPaymentIntentTx(ctx context.Context, q execQuerier, id, merchantID, lockClause string) (*PaymentIntent, error) {
	var whereClause string
	var args []interface{}

	if merchantID != "" {
		whereClause = fmt.Sprintf("WHERE pi.id = $1 AND pi.merchant_id = $2 %s", lockClause)
		args = []interface{}{id, merchantID}
	} else {
		whereClause = fmt.Sprintf("WHERE pi.id = $1 %s", lockClause)
		args = []interface{}{id}
	}

	query := `SELECT pi.id, pi.merchant_id, pi.customer_id, pi.amount, pi.amount_capturable, pi.amount_received,
		pi.capture_method, pi.currency, pi.status, 
		pi.idempotency_key, pi.description, pi.metadata, pi.failure_reason, pi.processor,
		COALESCE(pi.wallet_provider, ''), pi.return_url, pi.cancel_url, pi.client_secret,
		COALESCE(pi.redirect_url, ''), COALESCE(pi.provider_ref, ''),
		pi.created_at, pi.updated_at
		FROM payment_intents pi ` + whereClause

	var pi PaymentIntent
	var metaBytes []byte
	var failureReason *string
	var redirectURL string
	var providerRef string
	var walletProvider string

	err := q.QueryRow(ctx, query, args...).Scan(
		&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.AmountCapturable, &pi.AmountReceived,
		&pi.CaptureMethod, &pi.Currency, &pi.Status,
		&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
		&pi.PaymentMethod, &walletProvider,
		&pi.ReturnURL, &pi.CancelURL, &pi.ClientSecret,
		&redirectURL, &providerRef,
		&pi.CreatedAt, &pi.UpdatedAt,
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

	if redirectURL != "" {
		pi.RedirectURL = &redirectURL
	}
	if providerRef != "" {
		pi.ProviderRef = &providerRef
	}
	pi.WalletProvider = walletProvider

	return &pi, nil
}

func (r *Repository) GetPaymentIntentByProviderRef(ctx context.Context, providerRef string) (*PaymentIntent, error) {
	query := `SELECT pi.id, pi.merchant_id, pi.customer_id, pi.amount, pi.amount_capturable, pi.amount_received,
		pi.capture_method, pi.currency, pi.status,
		pi.idempotency_key, pi.description, pi.metadata, pi.failure_reason, pi.processor,
		COALESCE(pi.wallet_provider, ''), pi.return_url, pi.cancel_url, pi.client_secret,
		COALESCE(pi.redirect_url, ''), COALESCE(pi.provider_ref, ''),
		pi.created_at, pi.updated_at
		FROM payment_intents pi WHERE pi.provider_ref = $1`

	var pi PaymentIntent
	var metaBytes []byte
	var failureReason *string
	var redirectURL string
	var providerRefScanned string
	var walletProvider string

	err := r.Pool.QueryRow(ctx, query, providerRef).Scan(
		&pi.ID, &pi.MerchantID, &pi.CustomerID, &pi.Amount, &pi.AmountCapturable, &pi.AmountReceived,
		&pi.CaptureMethod, &pi.Currency, &pi.Status,
		&pi.IdempotencyKey, &pi.Description, &metaBytes, &failureReason,
		&pi.PaymentMethod, &walletProvider,
		&pi.ReturnURL, &pi.CancelURL, &pi.ClientSecret,
		&redirectURL, &providerRefScanned,
		&pi.CreatedAt, &pi.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get payment intent by provider ref: %w", err)
	}

	if metaBytes != nil {
		if err := json.Unmarshal(metaBytes, &pi.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	if failureReason != nil {
		pi.ErrorMessage = failureReason
	}

	if redirectURL != "" {
		pi.RedirectURL = &redirectURL
	}
	if providerRefScanned != "" {
		pi.ProviderRef = &providerRefScanned
	}
	pi.WalletProvider = walletProvider

	return &pi, nil
}
