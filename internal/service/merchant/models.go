package merchant

import "time"

type Merchant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	WebhookURL *string  `json:"webhook_url,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateMerchantRequest struct {
	Name       *string `json:"name,omitempty"`
	WebhookURL *string `json:"webhook_url,omitempty"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type APIKey struct {
	ID          string     `json:"id"`
	MerchantID  string     `json:"merchant_id"`
	KeyPrefix   string     `json:"key_prefix"`
	Name        string     `json:"name"`
	Permissions []string   `json:"permissions"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateAPIKeyResponse struct {
	APIKey
	FullKey string `json:"full_key"`
}
