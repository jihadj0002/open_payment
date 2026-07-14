package customer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

func (r *Repository) Create(ctx context.Context, c *Customer) error {
	metaBytes, err := json.Marshal(c.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `INSERT INTO customers (id, merchant_id, email, phone, name, metadata) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at, updated_at`

	return r.Pool.QueryRow(ctx, query, c.ID, c.MerchantID, c.Email, c.Phone, c.Name, metaBytes).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id, merchantID string) (*Customer, error) {
	query := `SELECT id, merchant_id, email, phone, name, metadata, created_at, updated_at FROM customers WHERE id = $1 AND merchant_id = $2`

	var c Customer
	var metaBytes []byte
	err := r.Pool.QueryRow(ctx, query, id, merchantID).Scan(
		&c.ID, &c.MerchantID, &c.Email, &c.Phone, &c.Name, &metaBytes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get customer by id: %w", err)
	}

	if err := json.Unmarshal(metaBytes, &c.Metadata); err != nil {
		c.Metadata = map[string]string{}
	}

	return &c, nil
}

func (r *Repository) List(ctx context.Context, merchantID string, limit, offset int) ([]Customer, error) {
	query := `SELECT id, merchant_id, email, phone, name, metadata, created_at, updated_at FROM customers WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, merchantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		var metaBytes []byte
		if err := rows.Scan(&c.ID, &c.MerchantID, &c.Email, &c.Phone, &c.Name, &metaBytes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		if err := json.Unmarshal(metaBytes, &c.Metadata); err != nil {
			c.Metadata = map[string]string{}
		}
		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if customers == nil {
		customers = []Customer{}
	}

	return customers, nil
}

func (r *Repository) Update(ctx context.Context, id, merchantID string, req UpdateCustomerRequest) (*Customer, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, *req.Email)
		argIdx++
	}
	if req.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argIdx))
		args = append(args, *req.Phone)
		argIdx++
	}
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Metadata != nil {
		metaBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("marshal metadata: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("metadata = $%d", argIdx))
		args = append(args, metaBytes)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id, merchantID)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, id, merchantID)

	query := fmt.Sprintf(
		`UPDATE customers SET %s WHERE id = $%d AND merchant_id = $%d RETURNING id, merchant_id, email, phone, name, metadata, created_at, updated_at`,
		strings.Join(setClauses, ", "),
		argIdx, argIdx+1,
	)

	var c Customer
	var metaBytes []byte
	err := r.Pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.MerchantID, &c.Email, &c.Phone, &c.Name, &metaBytes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("update customer: %w", err)
	}

	if err := json.Unmarshal(metaBytes, &c.Metadata); err != nil {
		c.Metadata = map[string]string{}
	}

	return &c, nil
}

func (r *Repository) CreatePaymentMethod(ctx context.Context, merchantID string, pm *PaymentMethod) error {
	query := `INSERT INTO payment_methods (id, merchant_id, customer_id, type, token, last4, brand, exp_month, exp_year, cardholder_name, fingerprint, wallet_type, wallet_phone, is_default, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING created_at, updated_at`

	return r.Pool.QueryRow(ctx, query,
		pm.ID, merchantID, pm.CustomerID, pm.Type, pm.Token, pm.Last4, pm.Brand,
		pm.ExpMonth, pm.ExpYear, pm.CardholderName, pm.Fingerprint, pm.WalletType, pm.WalletPhone,
		pm.IsDefault, pm.IsActive,
	).Scan(&pm.CreatedAt, &pm.UpdatedAt)
}

func (r *Repository) ListPaymentMethods(ctx context.Context, customerID, merchantID string) ([]PaymentMethod, error) {
	query := `SELECT id, customer_id, type, token, last4, brand, exp_month, exp_year, cardholder_name, wallet_type, wallet_phone, fingerprint, is_default, is_active, created_at, updated_at FROM payment_methods WHERE customer_id = $1 AND merchant_id = $2 AND is_active = true ORDER BY is_default DESC, created_at DESC`

	rows, err := r.Pool.Query(ctx, query, customerID, merchantID)
	if err != nil {
		return nil, fmt.Errorf("list payment methods: %w", err)
	}
	defer rows.Close()

	var methods []PaymentMethod
	for rows.Next() {
		var pm PaymentMethod
		if err := rows.Scan(
			&pm.ID, &pm.CustomerID, &pm.Type, &pm.Token, &pm.Last4, &pm.Brand,
			&pm.ExpMonth, &pm.ExpYear, &pm.CardholderName, &pm.WalletType, &pm.WalletPhone,
			&pm.Fingerprint, &pm.IsDefault, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan payment method: %w", err)
		}
		methods = append(methods, pm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if methods == nil {
		methods = []PaymentMethod{}
	}

	return methods, nil
}

func (r *Repository) GetPaymentMethod(ctx context.Context, id, customerID, merchantID string) (*PaymentMethod, error) {
	query := `SELECT id, customer_id, type, token, last4, brand, exp_month, exp_year, cardholder_name, wallet_type, wallet_phone, fingerprint, is_default, is_active, created_at, updated_at FROM payment_methods WHERE id = $1 AND customer_id = $2 AND merchant_id = $3`

	var pm PaymentMethod
	err := r.Pool.QueryRow(ctx, query, id, customerID, merchantID).Scan(
		&pm.ID, &pm.CustomerID, &pm.Type, &pm.Token, &pm.Last4, &pm.Brand,
		&pm.ExpMonth, &pm.ExpYear, &pm.CardholderName, &pm.WalletType, &pm.WalletPhone,
		&pm.Fingerprint, &pm.IsDefault, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get payment method: %w", err)
	}

	return &pm, nil
}

func (r *Repository) DeactivatePaymentMethod(ctx context.Context, id, customerID, merchantID string) error {
	query := `UPDATE payment_methods SET is_active = false, updated_at = NOW() WHERE id = $1 AND customer_id = $2 AND merchant_id = $3`
	result, err := r.Pool.Exec(ctx, query, id, customerID, merchantID)
	if err != nil {
		return fmt.Errorf("deactivate payment method: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}
	return nil
}

func (r *Repository) UnsetDefaultPaymentMethods(ctx context.Context, customerID, merchantID string) error {
	_, err := r.Pool.Exec(ctx, `UPDATE payment_methods SET is_default = false, updated_at = NOW() WHERE customer_id = $1 AND merchant_id = $2`, customerID, merchantID)
	if err != nil {
		return fmt.Errorf("unset default payment methods: %w", err)
	}
	return nil
}
