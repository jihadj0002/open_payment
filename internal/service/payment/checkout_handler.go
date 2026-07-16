package payment

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
)

func HandleGetCheckoutSession(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "payment ID is required")
			return
		}

		pi, err := svc.GetPaymentByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "payment_not_found", "payment not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		resp := map[string]interface{}{
			"id":              pi.ID,
			"amount":          pi.Amount,
			"currency":        pi.Currency,
			"status":          pi.Status,
			"payment_method":  pi.PaymentMethod,
			"description":     pi.Description,
			"return_url":      pi.ReturnURL,
			"cancel_url":      pi.CancelURL,
			"redirect_url":    pi.RedirectURL,
			"client_secret":   pi.ClientSecret,
			"metadata":        pi.Metadata,
		}

		api.RespondJSON(w, http.StatusOK, resp)
	}
}

type initiateCheckoutRequest struct {
	PaymentMethod string `json:"payment_method"`
	Provider      string `json:"provider,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	ReturnURL     string `json:"return_url,omitempty"`
	CancelURL     string `json:"cancel_url,omitempty"`
}

func HandleInitiateCheckout(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "payment ID is required")
			return
		}

		var req initiateCheckoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		pi, err := svc.GetPaymentByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "payment_not_found", "payment not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		if pi.Status != StatusCreated && pi.Status != StatusPending {
			api.RespondError(w, http.StatusBadRequest, "invalid_state", "payment_not_initiable",
				"payment is not in an initiable state")
			return
		}

		if req.ReturnURL != "" {
			pi.ReturnURL = &req.ReturnURL
		}
		if req.CancelURL != "" {
			pi.CancelURL = &req.CancelURL
		}
		if req.CustomerPhone != "" {
			if pi.Metadata == nil {
				pi.Metadata = make(map[string]string)
			}
			pi.Metadata["customer_phone"] = req.CustomerPhone
		}

		pi.PaymentMethod = req.PaymentMethod
		pi.WalletProvider = req.Provider

		if err := svc.ProcessPayment(r.Context(), pi); err != nil {
			api.RespondError(w, http.StatusBadRequest, "payment_error", "process_failed", err.Error())
			return
		}

		resp := map[string]interface{}{
			"id":            pi.ID,
			"status":        pi.Status,
			"redirect_url":  pi.RedirectURL,
			"client_secret": pi.ClientSecret,
		}

		api.RespondJSON(w, http.StatusOK, resp)
	}
}

func HandleCheckoutSuccess(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "payment ID is required")
			return
		}

		pi, err := svc.GetPaymentByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondJSON(w, http.StatusOK, map[string]interface{}{
					"success": false,
					"message": "Payment not found",
				})
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"success":   pi.Status == StatusSucceeded || pi.Status == StatusCaptured,
			"status":    pi.Status,
			"amount":    pi.Amount,
			"currency":  pi.Currency,
			"returnURL": pi.ReturnURL,
		})
	}
}
