package ledger

import (
	"net/http"
	"strconv"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/service/auth"
)

func HandleGetBalance(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		currency := r.URL.Query().Get("currency")
		if currency == "" {
			currency = "BDT"
		}

		balance, err := svc.GetBalance(r.Context(), claims.MerchantID, currency)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, balance)
	}
}

func HandleGetTransactions(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		currency := r.URL.Query().Get("currency")
		if currency == "" {
			currency = "BDT"
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage < 1 || perPage > 100 {
			perPage = 20
		}

		offset := (page - 1) * perPage

		txs, err := svc.GetBalanceTransactions(r.Context(), claims.MerchantID, currency, perPage, offset)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(txs))
		for i, tx := range txs {
			data[i] = tx
		}

		api.RespondPaginated(w, data, len(txs), page, perPage)
	}
}
