package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/openpayment/gateway/internal/config"
)

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidAPIKey       = errors.New("invalid API key")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrPublishableKeyNotAllowed = errors.New("publishable keys are not allowed for this operation")
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 30 * 24 * time.Hour
)

type jwtCustomClaims struct {
	MerchantID  string   `json:"merchant_id"`
	UserID      string   `json:"user_id,omitempty"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type AuthService struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewAuthService(cfg *config.Config, db *pgxpool.Pool) *AuthService {
	return &AuthService{cfg: cfg, db: db}
}

func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func GenerateAPIKey(prefix string, isLive bool) (string, string) {
	mode := "test"
	if isLive {
		mode = "live"
	}

	fullPrefix := prefix + mode + "_"

	b := make([]byte, 32)
	rand.Read(b)
	key := fullPrefix + hex.EncodeToString(b)

	return key, HashAPIKey(key)
}

func GenerateSecretAndPublishableKeys() (secretKey, pubKey, secretHash, pubHash string) {
	secretKey, secretHash = GenerateAPIKey("sk_", false)
	pubKey, pubHash = GenerateAPIKey("pk_", false)
	return
}

func (s *AuthService) GenerateTokenPair(claims Claims) (*TokenPair, error) {
	now := time.Now()

	jwtClaims := jwtCustomClaims{
		MerchantID:  claims.MerchantID,
		UserID:      claims.UserID,
		Role:        claims.Role,
		Permissions: claims.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpiry)),
			Issuer:    "open-payment-gateway",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}
	refreshToken := hex.EncodeToString(refreshBytes)

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return &Claims{
		MerchantID:  claims.MerchantID,
		UserID:      claims.UserID,
		Role:        claims.Role,
		Permissions: claims.Permissions,
	}, nil
}

func (s *AuthService) ValidateAPIKey(key string) (*Claims, error) {
	var (
		isSecret bool
	)

	if strings.HasPrefix(key, "sk_live_") || strings.HasPrefix(key, "sk_test_") {
		isSecret = true
	} else if strings.HasPrefix(key, "pk_live_") || strings.HasPrefix(key, "pk_test_") {
		isSecret = false
	} else {
		return nil, ErrInvalidAPIKey
	}

	keyHash := HashAPIKey(key)

	var (
		merchantID  string
		permissions []string
		status      string
		permBytes   []byte
	)

	err := s.db.QueryRow(
		context.Background(),
		`SELECT merchant_id, permissions, status FROM api_keys WHERE key_hash = $1`,
		keyHash,
	).Scan(&merchantID, &permBytes, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidAPIKey
		}
		return nil, fmt.Errorf("querying API key: %w", err)
	}

	if status != "active" {
		return nil, ErrInvalidAPIKey
	}

	if err := json.Unmarshal(permBytes, &permissions); err != nil {
		return nil, fmt.Errorf("parsing permissions: %w", err)
	}

	role := "api_secret"
	if !isSecret {
		role = "api_publishable"
	}

	if !isSecret {
		permissions = []string{"token:create"}
	}

	_, _ = s.db.Exec(
		context.Background(),
		`UPDATE api_keys SET last_used_at = NOW() WHERE key_hash = $1`,
		keyHash,
	)

	return &Claims{
		MerchantID:  merchantID,
		Role:        role,
		Permissions: permissions,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	var (
		id           string
		name         string
		passwordHash string
		status       string
	)

	err := s.db.QueryRow(
		ctx,
		`SELECT id, name, password_hash, status FROM merchants WHERE email = $1`,
		email,
	).Scan(&id, &name, &passwordHash, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("querying merchant: %w", err)
	}

	if status != "active" {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.GenerateTokenPair(Claims{
		MerchantID:  id,
		Role:        "merchant",
		Permissions: []string{"read", "write"},
	})
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*TokenPair, error) {
	var existingID string
	err := s.db.QueryRow(ctx, `SELECT id FROM merchants WHERE email = $1`, req.Email).Scan(&existingID)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("checking existing email: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var merchantID string
	err = tx.QueryRow(
		ctx,
		`INSERT INTO merchants (name, email, password_hash, status) VALUES ($1, $2, $3, 'active') RETURNING id`,
		req.Name, req.Email, string(hashedPassword),
	).Scan(&merchantID)
	if err != nil {
		return nil, fmt.Errorf("creating merchant: %w", err)
	}

	secretKey, pubKey, secretHash, pubHash := GenerateSecretAndPublishableKeys()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO api_keys (merchant_id, key_prefix, key_hash, name, permissions, status) VALUES ($1, $2, $3, $4, $5, 'active')`,
		merchantID, "sk_test_", secretHash, "Default Secret Key", `["read","write"]`,
	)
	if err != nil {
		return nil, fmt.Errorf("creating secret API key: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO api_keys (merchant_id, key_prefix, key_hash, name, permissions, status) VALUES ($1, $2, $3, $4, $5, 'active')`,
		merchantID, "pk_test_", pubHash, "Default Publishable Key", `["read"]`,
	)
	if err != nil {
		return nil, fmt.Errorf("creating publishable API key: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	_ = secretKey
	_ = pubKey

	return s.GenerateTokenPair(Claims{
		MerchantID:  merchantID,
		Role:        "merchant",
		Permissions: []string{"read", "write"},
	})
}
