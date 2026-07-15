package settlement

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterSettlementRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Post("/settlements", HandleTriggerSettlement(svc))
	r.With(authMW).Get("/settlements", HandleListSettlements(svc))
	r.With(authMW).Get("/settlements/{id}", HandleGetSettlement(svc))
	r.With(authMW).Get("/settlements/report", HandleSettlementReport(svc))
}
