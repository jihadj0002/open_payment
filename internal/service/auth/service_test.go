//go:build unit

package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openpayment/gateway/internal/config"
)

func TestGenerateTokenPair_ValidClaims(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key"}
	svc := NewAuthService(cfg, nil)

	claims := Claims{
		MerchantID:  "merch_1",
		UserID:      "user_1",
		Role:        "merchant",
		Permissions: []string{"read", "write"},
	}

	pair, err := svc.GenerateTokenPair(claims)
	assert.NoError(t, err)
	assert.NotNil(t, pair)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, "Bearer", pair.TokenType)
	assert.Equal(t, int(AccessTokenExpiry.Seconds()), pair.ExpiresIn)

	token, err := jwt.ParseWithClaims(pair.AccessToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	require.NoError(t, err)

	customClaims, ok := token.Claims.(*jwtCustomClaims)
	require.True(t, ok)
	assert.True(t, token.Valid)
	assert.Equal(t, "merch_1", customClaims.MerchantID)
	assert.Equal(t, "user_1", customClaims.UserID)
	assert.Equal(t, "merchant", customClaims.Role)
	assert.Equal(t, []string{"read", "write"}, customClaims.Permissions)
}

func TestGenerateTokenPair_DifferentSecrets(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key"}
	svc := NewAuthService(cfg, nil)

	pair1, _ := svc.GenerateTokenPair(Claims{MerchantID: "m1", Role: "merchant"})
	pair2, _ := svc.GenerateTokenPair(Claims{MerchantID: "m2", Role: "merchant"})

	assert.NotEqual(t, pair1.AccessToken, pair2.AccessToken)
}

func TestValidateToken_ValidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key"}
	svc := NewAuthService(cfg, nil)

	claims := Claims{
		MerchantID:  "merch_1",
		Role:        "merchant",
		Permissions: []string{"read"},
	}

	pair, err := svc.GenerateTokenPair(claims)
	require.NoError(t, err)

	parsed, err := svc.ValidateToken(pair.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, "merch_1", parsed.MerchantID)
	assert.Equal(t, "merchant", parsed.Role)
	assert.Equal(t, []string{"read"}, parsed.Permissions)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key"}
	svc := NewAuthService(cfg, nil)

	now := time.Now()
	jwtClaims := jwtCustomClaims{
		MerchantID: "merch_1",
		Role:       "merchant",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
			Issuer:    "open-payment-gateway",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)

	parsed, err := svc.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsed)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key"}
	svc := NewAuthService(cfg, nil)

	claims := jwtCustomClaims{
		MerchantID: "merch_1",
		Role:       "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "open-payment-gateway",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("different-secret"))

	parsed, err := svc.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsed)
}

func TestValidateToken_MalformedToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key"}
	svc := NewAuthService(cfg, nil)

	parsed, err := svc.ValidateToken("not-a-valid-token")
	assert.Error(t, err)
	assert.Nil(t, parsed)
}

func TestHashAPIKey_Consistent(t *testing.T) {
	key := "sk_test_abc123def456"
	hash1 := HashAPIKey(key)
	hash2 := HashAPIKey(key)
	assert.Equal(t, hash1, hash2)
	assert.NotEmpty(t, hash1)
}

func TestHashAPIKey_DifferentKeys(t *testing.T) {
	hash1 := HashAPIKey("sk_test_abc")
	hash2 := HashAPIKey("sk_test_def")
	assert.NotEqual(t, hash1, hash2)
}

func TestGenerateAPIKey_Format(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		live   bool
		expect string
	}{
		{"secret test", "sk_", false, "sk_test_"},
		{"secret live", "sk_", true, "sk_live_"},
		{"publishable test", "pk_", false, "pk_test_"},
		{"publishable live", "pk_", true, "pk_live_"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullKey, keyHash := GenerateAPIKey(tt.prefix, tt.live)
			assert.True(t, strings.HasPrefix(fullKey, tt.expect))
			assert.Len(t, fullKey, len(tt.expect)+64)
			assert.NotEmpty(t, keyHash)
			assert.Equal(t, HashAPIKey(fullKey), keyHash)
		})
	}
}

func TestGenerateSecretAndPublishableKeys(t *testing.T) {
	secretKey, pubKey, secretHash, pubHash := GenerateSecretAndPublishableKeys()

	assert.True(t, strings.HasPrefix(secretKey, "sk_test_"))
	assert.True(t, strings.HasPrefix(pubKey, "pk_test_"))
	assert.Equal(t, HashAPIKey(secretKey), secretHash)
	assert.Equal(t, HashAPIKey(pubKey), pubHash)
	assert.NotEqual(t, secretKey, pubKey)
}


