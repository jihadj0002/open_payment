package merchant

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
	"github.com/openpayment/gateway/internal/service/auth"
)

func HandleGetProfile(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		merchant, err := svc.GetProfile(r.Context(), claims.MerchantID)
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

func HandleUpdateProfile(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var req UpdateMerchantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if req.Name != nil && *req.Name == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_name", "name cannot be empty")
			return
		}

		merchant, err := svc.UpdateProfile(r.Context(), claims.MerchantID, req)
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

func HandleListAPIKeys(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		keys, err := svc.ListAPIKeys(r.Context(), claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, keys)
	}
}

func HandleCreateAPIKey(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var req CreateAPIKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if req.Name == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_name", "name is required")
			return
		}

		key, err := svc.CreateAPIKey(r.Context(), claims.MerchantID, req)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusCreated, key)
	}
}

func HandleRevokeAPIKey(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		keyID := chi.URLParam(r, "id")
		if keyID == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "API key ID is required")
			return
		}

		if err := svc.RevokeAPIKey(r.Context(), keyID, claims.MerchantID); err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "api_key_not_found", "API key not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
