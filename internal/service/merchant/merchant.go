package merchant

import "context"

type MerchantRepository interface {
	GetByID(ctx context.Context, id string) (*Merchant, error)
	GetByEmail(ctx context.Context, email string) (*Merchant, error)
	Update(ctx context.Context, id string, req UpdateMerchantRequest) (*Merchant, error)
	ListAPIKeys(ctx context.Context, merchantID string) ([]APIKey, error)
	CreateAPIKey(ctx context.Context, merchantID, name, keyPrefix, keyHash string, permissions []string) (*APIKey, error)
	RevokeAPIKey(ctx context.Context, keyID, merchantID string) error
	UpdateLastUsed(ctx context.Context, keyHash string) error
	GetAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error)
}
