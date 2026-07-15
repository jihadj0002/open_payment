package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
	"github.com/openpayment/gateway/internal/service/merchant"
)

const defaultPerPage = 20

func HandleListMerchants(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		query := r.URL.Query().Get("q")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if page < 1 {
			page = 1
		}
		if perPage < 1 || perPage > 100 {
			perPage = defaultPerPage
		}

		merchants, err := svc.ListMerchants(r.Context(), status, query, page, perPage)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(merchants))
		for i, v := range merchants {
			data[i] = v
		}
		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleGetMerchant(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "merchant ID is required")
			return
		}

		merchant, err := svc.GetMerchantDetail(r.Context(), id)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "merchant_not_found", "merchant not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, merchant)
	}
}

func HandleApproveMerchant(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "merchant ID is required")
			return
		}

		if err := svc.ApproveMerchant(r.Context(), id); err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "merchant_not_found", "merchant not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "approved"})
	}
}

func HandleSuspendMerchant(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "merchant ID is required")
			return
		}

		var req struct {
			Reason string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			req.Reason = ""
		}

		if err := svc.SuspendMerchant(r.Context(), id, req.Reason); err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "merchant_not_found", "merchant not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "suspended"})
	}
}

func HandleTerminateMerchant(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "merchant ID is required")
			return
		}

		if err := svc.TerminateMerchant(r.Context(), id); err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "merchant_not_found", "merchant not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "terminated"})
	}
}

func HandleListTransactions(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchantID := r.URL.Query().Get("merchant_id")
		status := r.URL.Query().Get("status")
		paymentMethod := r.URL.Query().Get("payment_method")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if page < 1 {
			page = 1
		}
		if perPage < 1 || perPage > 100 {
			perPage = defaultPerPage
		}

		transactions, err := svc.ListTransactions(r.Context(), merchantID, status, paymentMethod, page, perPage)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(transactions))
		for i, v := range transactions {
			data[i] = v
		}
		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleGetTransaction(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "transaction ID is required")
			return
		}

		tx, err := svc.GetTransaction(r.Context(), id)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "transaction_not_found", "transaction not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, tx)
	}
}

func HandleListDisputes(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if page < 1 {
			page = 1
		}
		if perPage < 1 || perPage > 100 {
			perPage = defaultPerPage
		}

		disputes, err := svc.ListDisputes(r.Context(), status, page, perPage)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(disputes))
		for i, v := range disputes {
			data[i] = v
		}
		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleGetDispute(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "dispute ID is required")
			return
		}

		dispute, err := svc.GetDispute(r.Context(), id)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "dispute_not_found", "dispute not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, dispute)
	}
}

func HandleResolveDispute(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "dispute ID is required")
			return
		}

		var req struct {
			Resolution string `json:"resolution"`
			Notes      string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if err := svc.ResolveDispute(r.Context(), id, req.Resolution, req.Notes); err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "dispute_not_found", "dispute not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "resolved"})
	}
}

func HandleListFeeConfigs(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		configs, err := svc.ListFeeConfigs(r.Context())
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(configs))
		for i, v := range configs {
			data[i] = v
		}
		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleCreateFeeConfig(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FeeConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if req.Name == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_name", "fee config name is required")
			return
		}

		if err := svc.CreateFeeConfig(r.Context(), &req); err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusCreated, req)
	}
}

func HandleUpdateFeeConfig(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "fee config ID is required")
			return
		}

		var req FeeConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if err := svc.UpdateFeeConfig(r.Context(), id, &req); err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "fee_config_not_found", "fee config not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, req)
	}
}

func HandleGetSystemConfig(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, err := svc.GetSystemConfig(r.Context())
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "config_not_found", "system config not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, cfg)
	}
}

func HandleUpdateSystemConfig(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SystemConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if err := svc.UpdateSystemConfig(r.Context(), &req); err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, req)
	}
}

func HandleListAuditLogs(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actorID := r.URL.Query().Get("actor_id")
		action := r.URL.Query().Get("action")
		resourceType := r.URL.Query().Get("resource_type")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if page < 1 {
			page = 1
		}
		if perPage < 1 || perPage > 100 {
			perPage = defaultPerPage
		}

		logs, err := svc.ListAuditLogs(r.Context(), actorID, action, resourceType, page, perPage)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(logs))
		for i, v := range logs {
			data[i] = v
		}
		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleGetMerchantStats(svc *StatsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "merchant ID is required")
			return
		}

		period := r.URL.Query().Get("period")
		if period == "" {
			period = "daily"
		}

		now := time.Now()
		var from, to time.Time
		switch period {
		case "weekly":
			from = now.AddDate(0, 0, -7)
			to = now
		case "monthly":
			from = now.AddDate(0, -1, 0)
			to = now
		default:
			from = now.AddDate(0, 0, -1)
			to = now
		}

		stats, err := svc.GetUsageStats(r.Context(), id, period, from, to)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, stats)
	}
}

func HandleUpdateMerchant(svc *Service, merchantSvc *merchant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "merchant ID is required")
			return
		}

		var req struct {
			Status *string `json:"status,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if req.Status != nil {
			switch *req.Status {
			case "active":
				if err := svc.ApproveMerchant(r.Context(), id); err != nil {
					api.RespondStructuredError(w, err)
					return
				}
			case "suspended":
				if err := svc.SuspendMerchant(r.Context(), id, "admin action"); err != nil {
					api.RespondStructuredError(w, err)
					return
				}
			default:
				api.RespondError(w, http.StatusBadRequest, "validation_error", "invalid_status", "status must be active or suspended")
				return
			}
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func HandleResetMerchantAPIKeys(svc *merchant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "merchant ID is required")
			return
		}

		keys, err := svc.ListAPIKeys(r.Context(), id)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		for _, k := range keys {
			if err := svc.RevokeAPIKey(r.Context(), k.ID, id); err != nil {
				api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "failed to revoke key")
				return
			}
		}

		req := merchant.CreateAPIKeyRequest{Name: "Default API Key"}
		resp, err := svc.CreateAPIKey(r.Context(), id, req)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "failed to create new key")
			return
		}

		api.RespondJSON(w, http.StatusOK, resp)
	}
}

func HandleListFraudRules(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.RespondJSON(w, http.StatusOK, []interface{}{})
	}
}

func HandleCreateFraudRule(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name    string  `json:"name"`
			Type    string  `json:"type"`
			Threshold float64 `json:"threshold"`
			Action  string  `json:"action"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}
		api.RespondJSON(w, http.StatusCreated, req)
	}
}

func HandleUpdateFraudRule(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "rule ID is required")
			return
		}
		api.RespondJSON(w, http.StatusOK, map[string]string{"id": id, "status": "updated"})
	}
}

func HandleDeleteFraudRule(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "rule ID is required")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleGetMerchantOnboardingStatus(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "not_started",
			"steps": []map[string]interface{}{
				{"name": "business_info", "completed": false},
				{"name": "personal_info", "completed": false},
				{"name": "documents", "completed": false},
				{"name": "review", "completed": false},
			},
		})
	}
}

func HandleSubmitOnboarding(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			BusinessName string `json:"business_name"`
			BusinessType string `json:"business_type"`
			Website      string `json:"website"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Phone        string `json:"phone"`
			Address      string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}
		api.RespondJSON(w, http.StatusCreated, map[string]string{"status": "submitted"})
	}
}

func HandleUploadOnboardingDocument(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.RespondJSON(w, http.StatusCreated, map[string]string{"status": "uploaded"})
	}
}
