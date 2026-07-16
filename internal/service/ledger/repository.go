package ledger

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/openpayment/gateway/internal/database"
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

func (r *Repository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.Pool.Begin(ctx)
}

func (r *Repository) CreateEntry(ctx context.Context, e *Entry) error {
	return r.createEntryTx(ctx, r.Pool, e)
}

func (r *Repository) CreateEntryTx(ctx context.Context, tx pgx.Tx, e *Entry) error {
	return r.createEntryTx(ctx, tx, e)
}

func (r *Repository) createEntryTx(ctx context.Context, q execQuerier, e *Entry) error {
	e.ID = uuid.New().String()
	e.CreatedAt = time.Now()

	query := `INSERT INTO ledger_entries (id, transaction_id, merchant_id, entry_type, amount, currency, balance_before, balance_after, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.Pool.Exec(ctx, query,
		e.ID, e.TransactionID, e.MerchantID, e.EntryType, e.Amount, e.Currency,
		e.BalanceBefore, e.BalanceAfter, e.Description, e.CreatedAt,
	)
	return err
}

func (r *Repository) GetCurrentBalance(ctx context.Context, merchantID, currency string) (int64, error) {
	return r.getCurrentBalanceTx(ctx, r.Pool, merchantID, currency)
}

func (r *Repository) GetCurrentBalanceTx(ctx context.Context, tx pgx.Tx, merchantID, currency string) (int64, error) {
	return r.getCurrentBalanceTx(ctx, tx, merchantID, currency)
}

func (r *Repository) getCurrentBalanceTx(ctx context.Context, q execQuerier, merchantID, currency string) (int64, error) {
	query := `SELECT balance_after FROM ledger_entries WHERE merchant_id=$1 AND currency=$2 ORDER BY created_at DESC LIMIT 1`

	var balance int64
	err := q.QueryRow(ctx, query, merchantID, currency).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return balance, nil
}

func (r *Repository) GetPendingBalance(ctx context.Context, merchantID, currency string) (int64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM payment_intents WHERE merchant_id=$1 AND currency=$2 AND status IN ('authorized','processing')`

	var pending int64
	err := r.Pool.QueryRow(ctx, query, merchantID, currency).Scan(&pending)
	return pending, err
}

func (r *Repository) GetEntries(ctx context.Context, merchantID, currency string, limit, offset int) ([]Entry, error) {
	query := `SELECT id, transaction_id, merchant_id, entry_type, amount, currency, balance_before, balance_after, COALESCE(description, ''), created_at
		FROM ledger_entries WHERE merchant_id=$1 AND currency=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`

	rows, err := r.Pool.Query(ctx, query, merchantID, currency, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.TransactionID, &e.MerchantID, &e.EntryType, &e.Amount, &e.Currency,
			&e.BalanceBefore, &e.BalanceAfter, &e.Description, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []Entry{}
	}
	return entries, rows.Err()
}

func (r *Repository) GetEntriesByTransaction(ctx context.Context, transactionID string) ([]Entry, error) {
	query := `SELECT id, transaction_id, merchant_id, entry_type, amount, currency, balance_before, balance_after, COALESCE(description, ''), created_at
		FROM ledger_entries WHERE transaction_id=$1 ORDER BY created_at ASC`

	rows, err := r.Pool.Query(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.TransactionID, &e.MerchantID, &e.EntryType, &e.Amount, &e.Currency,
			&e.BalanceBefore, &e.BalanceAfter, &e.Description, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []Entry{}
	}
	return entries, rows.Err()
}
