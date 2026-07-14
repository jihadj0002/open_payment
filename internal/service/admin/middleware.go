package admin

import (
	"net/http"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/service/auth"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil || claims.Role != "admin" {
			api.RespondError(w, http.StatusForbidden, "auth_error", "forbidden", "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
