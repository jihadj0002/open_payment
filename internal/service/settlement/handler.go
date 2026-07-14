package settlement

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/service/auth"
)

func HandleTriggerSettlement(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var req SettlementRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		settlement, err := svc.TriggerSettlement(r.Context(), claims.MerchantID, req)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "settlement_error", "settlement_failed", err.Error())
			return
		}

		api.RespondJSON(w, http.StatusCreated, settlement)
	}
}

func HandleGetSettlement(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "settlement ID is required")
			return
		}

		settlement, err := svc.GetSettlement(r.Context(), id, claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusNotFound, "not_found", "settlement_not_found", "settlement not found")
			return
		}

		api.RespondJSON(w, http.StatusOK, settlement)
	}
}

func HandleListSettlements(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage < 1 || perPage > 100 {
			perPage = 20
		}

		settlements, err := svc.ListSettlements(r.Context(), claims.MerchantID, page, perPage)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(settlements))
		for i, s := range settlements {
			data[i] = s
		}

		api.RespondJSON(w, http.StatusOK, data)
	}
}
