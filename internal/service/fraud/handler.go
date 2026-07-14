package fraud

import (
	"encoding/json"
	"net/http"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/service/auth"
)

func HandleGetFraudConfig(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		cfg, err := svc.GetFraudConfig(r.Context(), claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, cfg)
	}
}

func HandleUpdateFraudConfig(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var cfg FraudConfig
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_json", "invalid request body")
			return
		}

		if err := svc.UpdateFraudConfig(r.Context(), claims.MerchantID, &cfg); err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, cfg)
	}
}
