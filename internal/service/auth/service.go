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
	"github.com/openpayment/gateway/internal/pkg/encrypt"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidToken           = errors.New("invalid token")
	ErrInvalidAPIKey          = errors.New("invalid API key")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrPublishableKeyNotAllowed = errors.New("publishable keys are not allowed for this operation")
	ErrInvalidResetToken      = errors.New("invalid or expired reset token")
	ErrResetTokenUsed         = errors.New("reset token has already been used")
)

type Auditor interface {
	Log(ctx context.Context, actorID, action, resourceType, resourceID, details string) error
}

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
	cfg          *config.Config
	db           *pgxpool.Pool
	auditor      Auditor
	emailSender  ResetTokenSender
}

type ResetTokenSender interface {
	SendResetToken(to, token string) error
}

func NewAuthService(cfg *config.Config, db *pgxpool.Pool) *AuthService {
	return &AuthService{cfg: cfg, db: db}
}

func (s *AuthService) WithAuditor(auditor Auditor) *AuthService {
	s.auditor = auditor
	return s
}

func (s *AuthService) WithEmailSender(es ResetTokenSender) *AuthService {
	s.emailSender = es
	return s
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

	// 28 bytes = 56 hex chars, + 8 prefix = 64 chars (fits VARCHAR(64))
	b := make([]byte, 28)
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
	refreshHash := HashAPIKey(refreshToken)

	if claims.MerchantID != "" && s.db != nil {
		_, _ = s.db.Exec(
			context.Background(),
			`INSERT INTO refresh_tokens (merchant_id, token_hash, expires_at) VALUES ($1, $2, NOW() + $3::INTERVAL)`,
			claims.MerchantID, refreshHash, "30 days",
		)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	refreshHash := HashAPIKey(refreshToken)

	var merchantID string
	var revokedAt *time.Time
	var expiresAt time.Time
	err := s.db.QueryRow(
		ctx,
		`SELECT merchant_id, revoked_at, expires_at FROM refresh_tokens WHERE token_hash = $1`,
		refreshHash,
	).Scan(&merchantID, &revokedAt, &expiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("querying refresh token: %w", err)
	}

	if revokedAt != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().After(expiresAt) {
		return nil, ErrInvalidToken
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1`, refreshHash)
	if err != nil {
		return nil, fmt.Errorf("revoking old token: %w", err)
	}

	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return nil, fmt.Errorf("generating new refresh token: %w", err)
	}
	newRefreshToken := hex.EncodeToString(refreshBytes)
	newRefreshHash := HashAPIKey(newRefreshToken)

	_, err = tx.Exec(
		ctx,
		`INSERT INTO refresh_tokens (merchant_id, token_hash, expires_at) VALUES ($1, $2, NOW() + $3::INTERVAL)`,
		merchantID, newRefreshHash, "30 days",
	)
	if err != nil {
		return nil, fmt.Errorf("storing new refresh token: %w", err)
	}

	jwtClaims := jwtCustomClaims{
		MerchantID: merchantID,
		Role:       "merchant",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			Issuer:    "open-payment-gateway",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) RevokeRefreshTokens(ctx context.Context, merchantID string) error {
	_, err := s.db.Exec(
		ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE merchant_id = $1 AND revoked_at IS NULL`,
		merchantID,
	)
	if err != nil {
		return fmt.Errorf("revoking refresh tokens: %w", err)
	}
	if s.auditor != nil {
		s.auditor.Log(ctx, merchantID, "refresh_tokens.revoked", "merchant", merchantID, "all refresh tokens revoked")
	}
	return nil
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

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	var merchantID string
	err := s.db.QueryRow(ctx, `SELECT id FROM merchants WHERE email = $1`, email).Scan(&merchantID)

	exists := err == nil
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("querying merchant: %w", err)
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("generating reset token: %w", err)
	}
	rawToken := hex.EncodeToString(tokenBytes)
	tokenHash := HashAPIKey(rawToken)

	if exists {
		_, err = s.db.Exec(
			ctx,
			`INSERT INTO password_reset_tokens (merchant_id, token_hash, expires_at) VALUES ($1, $2, NOW() + INTERVAL '1 hour')`,
			merchantID, tokenHash,
		)
		if err != nil {
			return fmt.Errorf("storing reset token: %w", err)
		}
	}

	if s.emailSender != nil {
		if err := s.emailSender.SendResetToken(email, rawToken); err != nil {
			if s.cfg.Environment == "production" {
				return fmt.Errorf("sending reset email: %w", err)
			}
			log.Warn().Err(err).Str("email", email).Msg("failed to send reset email (non-fatal in dev)")
		}
	} else {
		log.Info().Str("email", email).Str("reset_token", rawToken).Msg("password reset token generated (no email sender configured)")
	}

	log.Info().Str("email", email).Msg("password reset token generation attempted")
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, password string) error {
	tokenHash := HashAPIKey(token)

	var merchantID string
	var usedAt *time.Time
	err := s.db.QueryRow(
		ctx,
		`SELECT merchant_id, used_at FROM password_reset_tokens WHERE token_hash = $1 AND expires_at > NOW()`,
		tokenHash,
	).Scan(&merchantID, &usedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidResetToken
		}
		return fmt.Errorf("querying reset token: %w", err)
	}

	if usedAt != nil {
		return ErrResetTokenUsed
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), encrypt.BcryptCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `UPDATE merchants SET password_hash = $1 WHERE id = $2`, string(hashedPassword), merchantID)
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE password_reset_tokens SET used_at = NOW() WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return fmt.Errorf("invalidating token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	committed = true

	if s.auditor != nil {
		s.auditor.Log(ctx, merchantID, "password.reset", "merchant", merchantID, "password reset completed")
	}

	log.Info().Str("merchant_id", merchantID).Msg("password reset successful")
	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
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
			if s.auditor != nil {
				s.auditor.Log(ctx, "", "login.failed", "merchant", "", fmt.Sprintf("failed login for email: %s (not found)", email))
			}
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("querying merchant: %w", err)
	}

	if status != "active" {
		if s.auditor != nil {
			s.auditor.Log(ctx, id, "login.failed", "merchant", id, fmt.Sprintf("failed login for inactive account: %s", email))
		}
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		if s.auditor != nil {
			s.auditor.Log(ctx, id, "login.failed", "merchant", id, fmt.Sprintf("failed login for email: %s (wrong password)", email))
		}
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := s.GenerateTokenPair(Claims{
		MerchantID:  id,
		Role:        "merchant",
		Permissions: []string{"read", "write"},
	})
	if err != nil {
		return nil, err
	}

	if s.auditor != nil {
		s.auditor.Log(ctx, id, "login.success", "merchant", id, fmt.Sprintf("successful login for email: %s", email))
	}

	return &AuthResponse{
		TokenPair: *tokenPair,
		User: User{
			MerchantID:  id,
			Role:        "merchant",
			Permissions: []string{"read", "write"},
		},
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	var existingID string
	err := s.db.QueryRow(ctx, `SELECT id FROM merchants WHERE email = $1`, req.Email).Scan(&existingID)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("checking existing email: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), encrypt.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
		}
	}()

	var merchantID string
	secretKey, pubKey, secretHash, pubHash := GenerateSecretAndPublishableKeys()

	encryptedSecret, _ := encrypt.Encrypt([]byte(secretKey))
	encryptedPub, _ := encrypt.Encrypt([]byte(pubKey))

	err = tx.QueryRow(
		ctx,
		`INSERT INTO merchants (name, email, password_hash, secret_key, public_key, status) VALUES ($1, $2, $3, $4, $5, 'active') RETURNING id`,
		req.Name, req.Email, string(hashedPassword), encryptedSecret, encryptedPub,
	).Scan(&merchantID)
	if err != nil {
		return nil, fmt.Errorf("creating merchant: %w", err)
	}

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
	committed = true

	if s.auditor != nil {
		s.auditor.Log(ctx, merchantID, "merchant.created", "merchant", merchantID, fmt.Sprintf("new merchant registered: %s", req.Email))
	}

	tokenPair, err := s.GenerateTokenPair(Claims{
		MerchantID:  merchantID,
		Role:        "merchant",
		Permissions: []string{"read", "write"},
	})
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		TokenPair: *tokenPair,
		User: User{
			MerchantID:  merchantID,
			Role:        "merchant",
			Permissions: []string{"read", "write"},
		},
		SecretKey: secretKey,
		PublicKey: pubKey,
	}, nil
}
