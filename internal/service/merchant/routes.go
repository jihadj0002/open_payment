package merchant

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterMerchantRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Get("/merchants/profile", HandleGetProfile(svc))
	r.With(authMW).Patch("/merchants/profile", HandleUpdateProfile(svc))
	r.With(authMW).Get("/merchants/api_keys", HandleListAPIKeys(svc))
	r.With(authMW).Post("/merchants/api_keys", HandleCreateAPIKey(svc))
	r.With(authMW).Delete("/merchants/api_keys/{id}", HandleRevokeAPIKey(svc))
}
