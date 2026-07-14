package ledger

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterLedgerRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Get("/balance", HandleGetBalance(svc))
	r.With(authMW).Get("/balance/transactions", HandleGetTransactions(svc))
}
