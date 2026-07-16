package auth

import "time"

type Claims struct {
	MerchantID  string   `json:"merchant_id"`
	UserID      string   `json:"user_id,omitempty"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthResponse struct {
	TokenPair `json:"token_pair"`
	User      User   `json:"user"`
	SecretKey string `json:"secret_key,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
}

type User struct {
	MerchantID  string   `json:"merchant_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
