package fraud

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterFraudRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Get("/merchants/fraud_config", HandleGetFraudConfig(svc))
	r.With(authMW).Put("/merchants/fraud_config", HandleUpdateFraudConfig(svc))
}
