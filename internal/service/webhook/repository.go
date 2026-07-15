package webhook

import (
	"context"
	"time"

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

func (r *Repository) CreateEndpoint(ctx context.Context, e *Endpoint) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO webhooks (merchant_id, event, url, secret, status)
		 VALUES ($1, $2, $3, $4, 'active')`,
		e.MerchantID, e.Event, e.URL, e.Secret,
	)
	return err
}

func (r *Repository) GetEndpoint(ctx context.Context, id, merchantID string) (*Endpoint, error) {
	var e Endpoint
	err := r.Pool.QueryRow(ctx,
		`SELECT id, merchant_id, event, url, secret, previous_secret, previous_secret_expires_at, status, created_at, updated_at
		 FROM webhooks WHERE id = $1 AND merchant_id = $2`,
		id, merchantID,
	).Scan(&e.ID, &e.MerchantID, &e.Event, &e.URL, &e.Secret, &e.PreviousSecret, &e.PreviousSecretExpires, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repository) ListEndpoints(ctx context.Context, merchantID string) ([]Endpoint, error) {
	rows, err := r.Pool.Query(ctx,
		`SELECT id, merchant_id, event, url, secret, previous_secret, previous_secret_expires_at, status, created_at, updated_at
		 FROM webhooks WHERE merchant_id = $1 AND status != 'deleted'
		 ORDER BY created_at DESC`,
		merchantID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []Endpoint
	for rows.Next() {
		var e Endpoint
		if err := rows.Scan(&e.ID, &e.MerchantID, &e.Event, &e.URL, &e.Secret, &e.PreviousSecret, &e.PreviousSecretExpires, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, e)
	}
	if endpoints == nil {
		endpoints = []Endpoint{}
	}
	return endpoints, rows.Err()
}

func (r *Repository) DeleteEndpoint(ctx context.Context, id, merchantID string) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE webhooks SET status = 'deleted', updated_at = NOW()
		 WHERE id = $1 AND merchant_id = $2`,
		id, merchantID,
	)
	return err
}

func (r *Repository) GetEndpointsForEvent(ctx context.Context, merchantID, event string) ([]Endpoint, error) {
	rows, err := r.Pool.Query(ctx,
		`SELECT id, merchant_id, event, url, secret, previous_secret, previous_secret_expires_at, status, created_at, updated_at
		 FROM webhooks WHERE merchant_id = $1 AND event = $2 AND status = 'active'`,
		merchantID, event,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []Endpoint
	for rows.Next() {
		var e Endpoint
		if err := rows.Scan(&e.ID, &e.MerchantID, &e.Event, &e.URL, &e.Secret, &e.PreviousSecret, &e.PreviousSecretExpires, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, e)
	}
	return endpoints, rows.Err()
}

func (r *Repository) CreateDelivery(ctx context.Context, d *Delivery) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO webhook_deliveries (webhook_id, event, payload, status, attempt, max_attempts, response_code, response_body, next_attempt_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		d.WebhookID, d.Event, d.Payload, d.Status, d.Attempt, d.MaxAttempts, d.ResponseCode, nil, d.NextAttemptAt,
	)
	return err
}

func (r *Repository) GetPendingDeliveries(ctx context.Context, limit int) ([]Delivery, error) {
	rows, err := r.Pool.Query(ctx,
		`SELECT id, webhook_id, event, payload, status, attempt, max_attempts, response_code, next_attempt_at, created_at
		 FROM webhook_deliveries
		 WHERE status = 'pending' AND (next_attempt_at IS NULL OR next_attempt_at <= NOW())
		 ORDER BY created_at ASC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []Delivery
	for rows.Next() {
		var d Delivery
		if err := rows.Scan(&d.ID, &d.WebhookID, &d.Event, &d.Payload, &d.Status, &d.Attempt, &d.MaxAttempts, &d.ResponseCode, &d.NextAttemptAt, &d.CreatedAt); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}
	return deliveries, rows.Err()
}

func (r *Repository) UpdateDelivery(ctx context.Context, d *Delivery) error {
	var responseBody *string
	_, err := r.Pool.Exec(ctx,
		`UPDATE webhook_deliveries
		 SET status = $1, attempt = $2, response_code = $3, response_body = $4, next_attempt_at = $5
		 WHERE id = $6`,
		d.Status, d.Attempt, d.ResponseCode, responseBody, d.NextAttemptAt, d.ID,
	)
	return err
}

func (r *Repository) GetDeliveryWebhook(ctx context.Context, webhookID string) (*Endpoint, error) {
	var e Endpoint
	err := r.Pool.QueryRow(ctx,
		`SELECT id, merchant_id, event, url, secret, previous_secret, previous_secret_expires_at, status, created_at, updated_at
		 FROM webhooks WHERE id = $1`,
		webhookID,
	).Scan(&e.ID, &e.MerchantID, &e.Event, &e.URL, &e.Secret, &e.PreviousSecret, &e.PreviousSecretExpires, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *Repository) RotateSecret(ctx context.Context, id, currentSecret, newSecret string, expiresAt time.Time) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE webhooks
		 SET previous_secret = $1, previous_secret_expires_at = $2, secret = $3, updated_at = NOW()
		 WHERE id = $4`,
		currentSecret, expiresAt, newSecret, id,
	)
	return err
}
