package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/service/auth"
)

func HandleCreateEndpoint(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var req CreateEndpointRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		endpoint, err := svc.CreateEndpoint(r.Context(), claims.MerchantID, req)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "validation_error", "invalid_endpoint", err.Error())
			return
		}

		api.RespondJSON(w, http.StatusCreated, endpoint)
	}
}

func HandleListEndpoints(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		endpoints, err := svc.ListEndpoints(r.Context(), claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(endpoints))
		for i, v := range endpoints {
			data[i] = v
		}

		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleDeleteEndpoint(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "endpoint ID is required")
			return
		}

		if err := svc.DeleteEndpoint(r.Context(), id, claims.MerchantID); err != nil {
			if err.Error() == "endpoint not found" {
				api.RespondError(w, http.StatusNotFound, "not_found", "endpoint_not_found", "endpoint not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleRotateSecret(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "endpoint ID is required")
			return
		}

		result, err := svc.RotateSecret(r.Context(), id, claims.MerchantID)
		if err != nil {
			if err.Error() == "endpoint not found" {
				api.RespondError(w, http.StatusNotFound, "not_found", "endpoint_not_found", "endpoint not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, result)
	}
}
