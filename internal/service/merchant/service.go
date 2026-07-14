package merchant

import (
	"context"
	"errors"

	"github.com/openpayment/gateway/internal/service/auth"
)

var (
	ErrNameRequired    = errors.New("name is required")
)

type Service struct {
	repo MerchantRepository
}

func NewService(repo MerchantRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, merchantID string) (*Merchant, error) {
	return s.repo.GetByID(ctx, merchantID)
}

func (s *Service) UpdateProfile(ctx context.Context, merchantID string, req UpdateMerchantRequest) (*Merchant, error) {
	if req.Name != nil && *req.Name == "" {
		return nil, ErrNameRequired
	}

	return s.repo.Update(ctx, merchantID, req)
}

func (s *Service) ListAPIKeys(ctx context.Context, merchantID string) ([]APIKey, error) {
	return s.repo.ListAPIKeys(ctx, merchantID)
}

func (s *Service) CreateAPIKey(ctx context.Context, merchantID string, req CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	if req.Name == "" {
		return nil, ErrNameRequired
	}

	fullKey, keyHash := auth.GenerateAPIKey("sk_", false)

	permissions := []string{"read", "write"}

	key, err := s.repo.CreateAPIKey(ctx, merchantID, req.Name, "sk_test_", keyHash, permissions)
	if err != nil {
		return nil, err
	}

	return &CreateAPIKeyResponse{
		APIKey:  *key,
		FullKey: fullKey,
	}, nil
}

func (s *Service) RevokeAPIKey(ctx context.Context, keyID, merchantID string) error {
	return s.repo.RevokeAPIKey(ctx, keyID, merchantID)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Merchant, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*Merchant, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) UpdateLastUsed(ctx context.Context, keyHash string) error {
	return s.repo.UpdateLastUsed(ctx, keyHash)
}

func (s *Service) GetAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	return s.repo.GetAPIKeyByHash(ctx, keyHash)
}
