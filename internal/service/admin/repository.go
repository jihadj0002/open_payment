package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/openpayment/gateway/internal/database"
	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
	"github.com/openpayment/gateway/internal/service/merchant"
)

type Repository struct {
	*database.BaseRepository
}

func NewRepository(db *database.PostgresDB) *Repository {
	return &Repository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *Repository) ListMerchants(ctx context.Context, status, query string, limit, offset int) ([]MerchantListItem, error) {
	args := []interface{}{}
	argIdx := 1

	sql := `SELECT m.id, m.name, m.email, m.status,
		COALESCE(v.volume, 0), COALESCE(v.count, 0), NULL, m.created_at
		FROM merchants m
		LEFT JOIN (
			SELECT merchant_id,
				SUM(amount) AS volume,
				COUNT(*) AS count
			FROM payment_intents
			WHERE status = 'succeeded'
				AND created_at >= $` + fmt.Sprintf("%d", argIdx) + `
			GROUP BY merchant_id
		) v ON v.merchant_id = m.id
		WHERE 1=1`

	argIdx++
	monthStart := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	args = append(args, monthStart)

	if status != "" {
		sql += fmt.Sprintf(" AND m.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if query != "" {
		sql += fmt.Sprintf(" AND (m.name ILIKE $%d OR m.email ILIKE $%d)", argIdx, argIdx+1)
		like := "%" + query + "%"
		args = append(args, like, like)
		argIdx += 2
	}

	sql += " ORDER BY m.created_at DESC"
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list merchants: %w", err)
	}
	defer rows.Close()

	var items []MerchantListItem
	for rows.Next() {
		var item MerchantListItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.Status,
			&item.VolumeThisMonth, &item.TransactionCount, &item.FeeConfigID, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan merchant: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if items == nil {
		items = []MerchantListItem{}
	}
	return items, nil
}

func (r *Repository) GetMerchantDetail(ctx context.Context, id string) (*merchant.Merchant, error) {
	query := `SELECT id, name, email, webhook_url, status, created_at, updated_at FROM merchants WHERE id = $1`

	var m merchant.Merchant
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.Email, &m.WebhookURL, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get merchant detail: %w", err)
	}
	return &m, nil
}

func (r *Repository) ApproveMerchant(ctx context.Context, id string) error {
	query := `UPDATE merchants SET status = 'active', updated_at = NOW() WHERE id = $1`
	result, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("approve merchant: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}
	return nil
}

func (r *Repository) SuspendMerchant(ctx context.Context, id, reason string) error {
	query := `UPDATE merchants SET status = 'suspended', updated_at = NOW() WHERE id = $1`
	result, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("suspend merchant: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}
	return nil
}

func (r *Repository) TerminateMerchant(ctx context.Context, id string) error {
	query := `UPDATE merchants SET status = 'terminated', updated_at = NOW() WHERE id = $1`
	result, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("terminate merchant: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}
	return nil
}

func (r *Repository) ListTransactions(ctx context.Context, merchantID, status, paymentMethod string, limit, offset int) ([]TransactionListItem, error) {
	args := []interface{}{}
	argIdx := 1

	sql := `SELECT t.id, t.payment_intent_id, t.merchant_id, m.name, t.type, t.status,
		t.amount, t.currency, COALESCE(t.processor_ref, ''), EXTRACT(EPOCH FROM t.created_at)::BIGINT
		FROM transactions t
		JOIN merchants m ON m.id = t.merchant_id
		WHERE 1=1`

	if merchantID != "" {
		sql += fmt.Sprintf(" AND t.merchant_id = $%d", argIdx)
		args = append(args, merchantID)
		argIdx++
	}
	if status != "" {
		sql += fmt.Sprintf(" AND t.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	sql += " ORDER BY t.created_at DESC"
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	var items []TransactionListItem
	for rows.Next() {
		var item TransactionListItem
		if err := rows.Scan(&item.ID, &item.PaymentID, &item.MerchantID, &item.MerchantName,
			&item.Type, &item.Status, &item.Amount, &item.Currency, &item.ProcessorRef, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if items == nil {
		items = []TransactionListItem{}
	}
	return items, nil
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
		t.ProcessorResponse = respBytes
	}

	return &t, nil
}

type Transaction struct {
	ID                string `json:"id"`
	PaymentIntentID   string `json:"payment_intent_id"`
	MerchantID        string `json:"merchant_id"`
	Type              string `json:"type"`
	Amount            int64  `json:"amount"`
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	ProcessorRef      *string `json:"processor_ref,omitempty"`
	ProcessorResponse []byte `json:"processor_response,omitempty"`
	Fee               int64  `json:"fee"`
	NetAmount         int64  `json:"net_amount"`
	CreatedAt         time.Time `json:"created_at"`
}

func (r *Repository) ListDisputes(ctx context.Context, status string, limit, offset int) ([]Dispute, error) {
	args := []interface{}{}
	argIdx := 1

	sql := `SELECT id, payment_id, merchant_id, amount, currency, status, reason, respond_by, resolved_at, created_at
		FROM disputes WHERE 1=1`

	if status != "" {
		sql += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	sql += " ORDER BY created_at DESC"
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list disputes: %w", err)
	}
	defer rows.Close()

	var items []Dispute
	for rows.Next() {
		var item Dispute
		if err := rows.Scan(&item.ID, &item.PaymentID, &item.MerchantID, &item.Amount,
			&item.Currency, &item.Status, &item.Reason, &item.RespondBy, &item.ResolvedAt, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan dispute: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if items == nil {
		items = []Dispute{}
	}
	return items, nil
}

func (r *Repository) GetDispute(ctx context.Context, id string) (*Dispute, error) {
	query := `SELECT id, payment_id, merchant_id, amount, currency, status, reason, respond_by, resolved_at, created_at
		FROM disputes WHERE id = $1`

	var d Dispute
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.PaymentID, &d.MerchantID, &d.Amount, &d.Currency, &d.Status, &d.Reason,
		&d.RespondBy, &d.ResolvedAt, &d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get dispute: %w", err)
	}
	return &d, nil
}

func (r *Repository) ResolveDispute(ctx context.Context, id, resolution, notes string) error {
	query := `UPDATE disputes SET status = 'resolved', resolved_at = NOW() WHERE id = $1`
	result, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("resolve dispute: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}
	return nil
}

func (r *Repository) ListFeeConfigs(ctx context.Context) ([]FeeConfig, error) {
	query := `SELECT id, name, rate, fixed_fee, cross_border_rate, monthly_fee, min_monthly_volume, status, created_at, updated_at
		FROM fee_configs ORDER BY created_at DESC`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list fee configs: %w", err)
	}
	defer rows.Close()

	var items []FeeConfig
	for rows.Next() {
		var item FeeConfig
		if err := rows.Scan(&item.ID, &item.Name, &item.Rate, &item.FixedFee,
			&item.CrossBorderRate, &item.MonthlyFee, &item.MinMonthlyVolume,
			&item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan fee config: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if items == nil {
		items = []FeeConfig{}
	}
	return items, nil
}

func (r *Repository) GetFeeConfig(ctx context.Context, id string) (*FeeConfig, error) {
	query := `SELECT id, name, rate, fixed_fee, cross_border_rate, monthly_fee, min_monthly_volume, status, created_at, updated_at
		FROM fee_configs WHERE id = $1`

	var f FeeConfig
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.Name, &f.Rate, &f.FixedFee,
		&f.CrossBorderRate, &f.MonthlyFee, &f.MinMonthlyVolume,
		&f.Status, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get fee config: %w", err)
	}
	return &f, nil
}

func (r *Repository) CreateFeeConfig(ctx context.Context, f *FeeConfig) error {
	query := `INSERT INTO fee_configs (name, rate, fixed_fee, cross_border_rate, monthly_fee, min_monthly_volume, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.Pool.QueryRow(ctx, query,
		f.Name, f.Rate, f.FixedFee, f.CrossBorderRate, f.MonthlyFee, f.MinMonthlyVolume, f.Status,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create fee config: %w", err)
	}
	return nil
}

func (r *Repository) UpdateFeeConfig(ctx context.Context, id string, f *FeeConfig) error {
	query := `UPDATE fee_configs SET
		name = COALESCE($2, name),
		rate = COALESCE($3, rate),
		fixed_fee = COALESCE($4, fixed_fee),
		cross_border_rate = COALESCE($5, cross_border_rate),
		monthly_fee = COALESCE($6, monthly_fee),
		min_monthly_volume = COALESCE($7, min_monthly_volume),
		status = COALESCE($8, status),
		updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, rate, fixed_fee, cross_border_rate, monthly_fee, min_monthly_volume, status, created_at, updated_at`

	err := r.Pool.QueryRow(ctx, query,
		id, f.Name, f.Rate, f.FixedFee, f.CrossBorderRate, f.MonthlyFee, f.MinMonthlyVolume, f.Status,
	).Scan(&f.ID, &f.Name, &f.Rate, &f.FixedFee,
		&f.CrossBorderRate, &f.MonthlyFee, &f.MinMonthlyVolume,
		&f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkgErr.ErrNotFound
		}
		return fmt.Errorf("update fee config: %w", err)
	}
	return nil
}

func (r *Repository) GetSystemConfig(ctx context.Context) (*SystemConfig, error) {
	query := `SELECT id, supported_currencies, supported_countries, max_transaction_amount, maintenance_mode, updated_at
		FROM system_config WHERE id = 'default'`

	var c SystemConfig
	var currenciesBytes, countriesBytes []byte

	err := r.Pool.QueryRow(ctx, query).Scan(
		&c.ID, &currenciesBytes, &countriesBytes,
		&c.MaxTransactionAmount, &c.MaintenanceMode, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get system config: %w", err)
	}

	if len(currenciesBytes) > 0 {
		if err := json.Unmarshal(currenciesBytes, &c.SupportedCurrencies); err != nil {
			return nil, fmt.Errorf("unmarshal currencies: %w", err)
		}
	}
	if len(countriesBytes) > 0 {
		if err := json.Unmarshal(countriesBytes, &c.SupportedCountries); err != nil {
			return nil, fmt.Errorf("unmarshal countries: %w", err)
		}
	}

	return &c, nil
}

func (r *Repository) UpdateSystemConfig(ctx context.Context, cfg *SystemConfig) error {
	currenciesBytes, err := json.Marshal(cfg.SupportedCurrencies)
	if err != nil {
		return fmt.Errorf("marshal currencies: %w", err)
	}
	countriesBytes, err := json.Marshal(cfg.SupportedCountries)
	if err != nil {
		return fmt.Errorf("marshal countries: %w", err)
	}

	query := `UPDATE system_config SET
		supported_currencies = $1,
		supported_countries = $2,
		max_transaction_amount = $3,
		maintenance_mode = $4,
		updated_at = NOW()
		WHERE id = 'default'
		RETURNING id, supported_currencies, supported_countries, max_transaction_amount, maintenance_mode, updated_at`

	var currenciesRes, countriesRes []byte
	err = r.Pool.QueryRow(ctx, query,
		currenciesBytes, countriesBytes, cfg.MaxTransactionAmount, cfg.MaintenanceMode,
	).Scan(&cfg.ID, &currenciesRes, &countriesRes, &cfg.MaxTransactionAmount, &cfg.MaintenanceMode, &cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update system config: %w", err)
	}

	if len(currenciesRes) > 0 {
		if err := json.Unmarshal(currenciesRes, &cfg.SupportedCurrencies); err != nil {
			return fmt.Errorf("unmarshal currencies: %w", err)
		}
	}
	if len(countriesRes) > 0 {
		if err := json.Unmarshal(countriesRes, &cfg.SupportedCountries); err != nil {
			return fmt.Errorf("unmarshal countries: %w", err)
		}
	}

	return nil
}

func (r *Repository) CreateAuditLog(ctx context.Context, e *AuditLogEntry) error {
	query := `INSERT INTO audit_logs (actor_id, action, resource_type, resource_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err := r.Pool.QueryRow(ctx, query,
		e.ActorID, e.Action, e.ResourceType, e.ResourceID, e.Details, e.IPAddress,
	).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *Repository) ListAuditLogs(ctx context.Context, actorID, action, resourceType string, limit, offset int) ([]AuditLogEntry, error) {
	args := []interface{}{}
	argIdx := 1

	sql := `SELECT id, actor_id, action, resource_type, COALESCE(resource_id, ''), COALESCE(details, ''), COALESCE(ip_address, ''), created_at
		FROM audit_logs WHERE 1=1`

	if actorID != "" {
		sql += fmt.Sprintf(" AND actor_id = $%d", argIdx)
		args = append(args, actorID)
		argIdx++
	}
	if action != "" {
		sql += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, action)
		argIdx++
	}
	if resourceType != "" {
		sql += fmt.Sprintf(" AND resource_type = $%d", argIdx)
		args = append(args, resourceType)
		argIdx++
	}

	sql += " ORDER BY created_at DESC"
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	defer rows.Close()

	var items []AuditLogEntry
	for rows.Next() {
		var item AuditLogEntry
		if err := rows.Scan(&item.ID, &item.ActorID, &item.Action, &item.ResourceType,
			&item.ResourceID, &item.Details, &item.IPAddress, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if items == nil {
		items = []AuditLogEntry{}
	}
	return items, nil
}
