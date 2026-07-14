package customer

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
	"github.com/openpayment/gateway/internal/service/auth"
)

func HandleCreateCustomer(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var req CreateCustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		customer, err := svc.CreateCustomer(r.Context(), claims.MerchantID, req)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusCreated, customer)
	}
}

func HandleGetCustomer(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "customer ID is required")
			return
		}

		customer, err := svc.GetCustomer(r.Context(), id, claims.MerchantID)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "customer_not_found", "customer not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, customer)
	}
}

func HandleListCustomers(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

		customers, err := svc.ListCustomers(r.Context(), claims.MerchantID, page, perPage)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(customers))
		for i, v := range customers {
			data[i] = v
		}

		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleUpdateCustomer(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "customer ID is required")
			return
		}

		var req UpdateCustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		customer, err := svc.UpdateCustomer(r.Context(), id, claims.MerchantID, req)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "customer_not_found", "customer not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, customer)
	}
}

func HandleAttachPaymentMethod(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		customerID := chi.URLParam(r, "id")
		if customerID == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "customer ID is required")
			return
		}

		var req AttachPaymentMethodRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		pm, err := svc.AttachPaymentMethod(r.Context(), claims.MerchantID, customerID, req)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "customer_not_found", "customer not found")
				return
			}
			if errors.Is(err, ErrInvalidCardNumber) || errors.Is(err, ErrInvalidExpiry) || errors.Is(err, ErrInvalidWallet) {
				api.RespondError(w, http.StatusBadRequest, "validation_error", "invalid_payment_method", err.Error())
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusCreated, pm)
	}
}

func HandleListPaymentMethods(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		customerID := chi.URLParam(r, "id")
		if customerID == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "customer ID is required")
			return
		}

		methods, err := svc.ListPaymentMethods(r.Context(), customerID, claims.MerchantID)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "customer_not_found", "customer not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(methods))
		for i, v := range methods {
			data[i] = v
		}

		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleDetachPaymentMethod(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		customerID := chi.URLParam(r, "id")
		pmID := chi.URLParam(r, "pm_id")
		if customerID == "" || pmID == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "customer ID and payment method ID are required")
			return
		}

		if err := svc.DetachPaymentMethod(r.Context(), customerID, claims.MerchantID, pmID); err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "payment_method_not_found", "payment method not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
