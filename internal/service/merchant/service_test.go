//go:build unit

package merchant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
)

type mockMerchantRepo struct {
	mock.Mock
}

func (m *mockMerchantRepo) GetByID(ctx context.Context, id string) (*Merchant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Merchant), args.Error(1)
}

func (m *mockMerchantRepo) GetByEmail(ctx context.Context, email string) (*Merchant, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Merchant), args.Error(1)
}

func (m *mockMerchantRepo) Update(ctx context.Context, id string, req UpdateMerchantRequest) (*Merchant, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Merchant), args.Error(1)
}

func (m *mockMerchantRepo) ListAPIKeys(ctx context.Context, merchantID string) ([]APIKey, error) {
	args := m.Called(ctx, merchantID)
	return args.Get(0).([]APIKey), args.Error(1)
}

func (m *mockMerchantRepo) CreateAPIKey(ctx context.Context, merchantID, name, keyPrefix, keyHash string, permissions []string) (*APIKey, error) {
	args := m.Called(ctx, merchantID, name, keyPrefix, keyHash, permissions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*APIKey), args.Error(1)
}

func (m *mockMerchantRepo) RevokeAPIKey(ctx context.Context, keyID, merchantID string) error {
	args := m.Called(ctx, keyID, merchantID)
	return args.Error(0)
}

func (m *mockMerchantRepo) UpdateLastUsed(ctx context.Context, keyHash string) error {
	args := m.Called(ctx, keyHash)
	return args.Error(0)
}

func (m *mockMerchantRepo) GetAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	args := m.Called(ctx, keyHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*APIKey), args.Error(1)
}

func TestGetProfile_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	expected := &Merchant{
		ID:     "merch_1",
		Name:   "Test Merchant",
		Email:  "test@example.com",
		Status: "active",
	}

	repo.On("GetByID", mock.Anything, "merch_1").Return(expected, nil).Once()

	result, err := svc.GetProfile(context.Background(), "merch_1")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, pkgErr.ErrNotFound).Once()

	result, err := svc.GetProfile(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.ErrorIs(t, err, pkgErr.ErrNotFound)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateProfile_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	name := "Updated Merchant"
	req := UpdateMerchantRequest{
		Name: &name,
	}

	expected := &Merchant{
		ID:     "merch_1",
		Name:   "Updated Merchant",
		Email:  "test@example.com",
		Status: "active",
	}

	repo.On("Update", mock.Anything, "merch_1", req).Return(expected, nil).Once()

	result, err := svc.UpdateProfile(context.Background(), "merch_1", req)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Merchant", result.Name)
	repo.AssertExpectations(t)
}

func TestUpdateProfile_EmptyName(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	empty := ""
	req := UpdateMerchantRequest{
		Name: &empty,
	}

	result, err := svc.UpdateProfile(context.Background(), "merch_1", req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNameRequired)
	assert.Nil(t, result)
}

func TestCreateAPIKey_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	req := CreateAPIKeyRequest{Name: "My API Key"}

	repo.On("CreateAPIKey", mock.Anything, "merch_1", "My API Key", "sk_test_", mock.AnythingOfType("string"), []string{"read", "write"}).Return(&APIKey{
		ID:     "key_1",
		Name:   "My API Key",
		KeyPrefix: "sk_test_",
		Status: "active",
	}, nil).Once()

	result, err := svc.CreateAPIKey(context.Background(), "merch_1", req)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.FullKey)
	assert.Equal(t, "My API Key", result.Name)
	assert.Len(t, result.FullKey, len("sk_test_")+64)
	repo.AssertExpectations(t)
}

func TestCreateAPIKey_EmptyName(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	req := CreateAPIKeyRequest{Name: ""}

	result, err := svc.CreateAPIKey(context.Background(), "merch_1", req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNameRequired)
	assert.Nil(t, result)
}

func TestRevokeAPIKey_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	repo.On("RevokeAPIKey", mock.Anything, "key_1", "merch_1").Return(nil).Once()

	err := svc.RevokeAPIKey(context.Background(), "key_1", "merch_1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRevokeAPIKey_NotFound(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	repo.On("RevokeAPIKey", mock.Anything, "key_nonexistent", "merch_1").Return(pkgErr.ErrNotFound).Once()

	err := svc.RevokeAPIKey(context.Background(), "key_nonexistent", "merch_1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, pkgErr.ErrNotFound)
	repo.AssertExpectations(t)
}

func TestListAPIKeys_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	keys := []APIKey{
		{ID: "key_1", Name: "Key 1", Status: "active"},
		{ID: "key_2", Name: "Key 2", Status: "revoked"},
	}

	repo.On("ListAPIKeys", mock.Anything, "merch_1").Return(keys, nil).Once()

	result, err := svc.ListAPIKeys(context.Background(), "merch_1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	expected := &Merchant{ID: "merch_1", Name: "Test"}
	repo.On("GetByID", mock.Anything, "merch_1").Return(expected, nil).Once()

	result, err := svc.GetByID(context.Background(), "merch_1")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetByEmail_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockMerchantRepo)
	svc := NewService(repo)

	expected := &Merchant{ID: "merch_1", Email: "test@example.com"}
	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(expected, nil).Once()

	result, err := svc.GetByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}
