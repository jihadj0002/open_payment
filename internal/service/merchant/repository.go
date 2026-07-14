package merchant

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

func (r *Repository) GetByID(ctx context.Context, id string) (*Merchant, error) {
	query := `SELECT id, name, email, webhook_url, status, created_at, updated_at FROM merchants WHERE id = $1`

	var m Merchant
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.Email, &m.WebhookURL, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get merchant by id: %w", err)
	}

	return &m, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*Merchant, error) {
	query := `SELECT id, name, email, webhook_url, status, created_at, updated_at FROM merchants WHERE email = $1`

	var m Merchant
	err := r.Pool.QueryRow(ctx, query, email).Scan(
		&m.ID, &m.Name, &m.Email, &m.WebhookURL, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get merchant by email: %w", err)
	}

	return &m, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateMerchantRequest) (*Merchant, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.WebhookURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("webhook_url = $%d", argIdx))
		args = append(args, *req.WebhookURL)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(
		`UPDATE merchants SET %s WHERE id = $%d RETURNING id, name, email, webhook_url, status, created_at, updated_at`,
		strings.Join(setClauses, ", "),
		argIdx,
	)

	var m Merchant
	err := r.Pool.QueryRow(ctx, query, args...).Scan(
		&m.ID, &m.Name, &m.Email, &m.WebhookURL, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("update merchant: %w", err)
	}

	return &m, nil
}

func (r *Repository) ListAPIKeys(ctx context.Context, merchantID string) ([]APIKey, error) {
	query := `SELECT id, merchant_id, key_prefix, name, permissions, last_used_at, status, created_at FROM api_keys WHERE merchant_id = $1 ORDER BY created_at DESC`

	rows, err := r.Pool.Query(ctx, query, merchantID)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var k APIKey
		var permBytes []byte
		err := rows.Scan(&k.ID, &k.MerchantID, &k.KeyPrefix, &k.Name, &permBytes, &k.LastUsedAt, &k.Status, &k.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan api key: %w", err)
		}
		if err := json.Unmarshal(permBytes, &k.Permissions); err != nil {
			return nil, fmt.Errorf("unmarshal permissions: %w", err)
		}
		keys = append(keys, k)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if keys == nil {
		keys = []APIKey{}
	}

	return keys, nil
}

func (r *Repository) CreateAPIKey(ctx context.Context, merchantID, name, keyPrefix, keyHash string, permissions []string) (*APIKey, error) {
	permBytes, err := json.Marshal(permissions)
	if err != nil {
		return nil, fmt.Errorf("marshal permissions: %w", err)
	}

	query := `INSERT INTO api_keys (merchant_id, key_prefix, key_hash, name, permissions, status) VALUES ($1, $2, $3, $4, $5, 'active') RETURNING id, created_at`

	var k APIKey
	k.MerchantID = merchantID
	k.KeyPrefix = keyPrefix
	k.Name = name
	k.Permissions = permissions
	k.Status = "active"

	err = r.Pool.QueryRow(ctx, query, merchantID, keyPrefix, keyHash, name, permBytes).Scan(&k.ID, &k.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	return &k, nil
}

func (r *Repository) RevokeAPIKey(ctx context.Context, keyID, merchantID string) error {
	query := `UPDATE api_keys SET status = 'revoked' WHERE id = $1 AND merchant_id = $2`
	result, err := r.Pool.Exec(ctx, query, keyID, merchantID)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pkgErr.ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateLastUsed(ctx context.Context, keyHash string) error {
	_, err := r.Pool.Exec(ctx, `UPDATE api_keys SET last_used_at = NOW() WHERE key_hash = $1`, keyHash)
	if err != nil {
		return fmt.Errorf("update last used: %w", err)
	}
	return nil
}

func (r *Repository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	query := `SELECT id, merchant_id, key_prefix, name, permissions, last_used_at, status, created_at FROM api_keys WHERE key_hash = $1`

	var k APIKey
	var permBytes []byte
	err := r.Pool.QueryRow(ctx, query, keyHash).Scan(
		&k.ID, &k.MerchantID, &k.KeyPrefix, &k.Name, &permBytes, &k.LastUsedAt, &k.Status, &k.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErr.ErrNotFound
		}
		return nil, fmt.Errorf("get api key by hash: %w", err)
	}

	if err := json.Unmarshal(permBytes, &k.Permissions); err != nil {
		return nil, fmt.Errorf("unmarshal permissions: %w", err)
	}

	return &k, nil
}

