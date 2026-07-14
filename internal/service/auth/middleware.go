package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/openpayment/gateway/internal/api"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

func GetClaims(ctx context.Context) *Claims {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}

func AuthMiddleware(authService *AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "invalid authorization header format")
				return
			}
			token := parts[1]

			claims, err := authService.ValidateToken(token)
			if err != nil {
				claims, err = authService.ValidateAPIKey(token)
				if err != nil {
					api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "invalid or expired token")
					return
				}
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
