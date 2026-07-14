package payment

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

func HandleCreatePayment(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var req CreatePaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		pi, err := svc.CreatePayment(r.Context(), claims.MerchantID, req)
		if err != nil {
			if errors.Is(err, pkgErr.ErrConflict) {
				api.RespondError(w, http.StatusConflict, "conflict", "idempotency_conflict", err.Error())
				return
			}
			api.RespondError(w, http.StatusBadRequest, "validation_error", "invalid_payment", err.Error())
			return
		}

		api.RespondJSON(w, http.StatusCreated, pi)
	}
}

func HandleGetPayment(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "payment ID is required")
			return
		}

		pi, err := svc.GetPayment(r.Context(), id, claims.MerchantID)
		if err != nil {
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "payment_not_found", "payment not found")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, pi)
	}
}

func HandleListPayments(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		intents, err := svc.ListPayments(r.Context(), claims.MerchantID, limit, offset)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(intents))
		for i, v := range intents {
			data[i] = v
		}

		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleCapturePayment(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "payment ID is required")
			return
		}

		var req CaptureRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		pi, err := svc.CapturePayment(r.Context(), id, claims.MerchantID, req.AmountToCapture)
		if err != nil {
			if errors.Is(err, ErrPaymentNotCapturable) {
				api.RespondError(w, http.StatusBadRequest, "invalid_state", "not_capturable", err.Error())
				return
			}
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "payment_not_found", "payment not found")
				return
			}
			api.RespondError(w, http.StatusBadRequest, "validation_error", "capture_failed", err.Error())
			return
		}

		api.RespondJSON(w, http.StatusOK, pi)
	}
}

func HandleRefundPayment(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "payment ID is required")
			return
		}

		var req struct {
			Amount int64  `json:"amount"`
			Reason string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		tx, err := svc.RefundPayment(r.Context(), id, claims.MerchantID, req.Amount, req.Reason)
		if err != nil {
			if errors.Is(err, ErrPaymentNotRefundable) {
				api.RespondError(w, http.StatusBadRequest, "invalid_state", "not_refundable", err.Error())
				return
			}
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "payment_not_found", "payment not found")
				return
			}
			api.RespondError(w, http.StatusBadRequest, "validation_error", "refund_failed", err.Error())
			return
		}

		api.RespondJSON(w, http.StatusOK, tx)
	}
}

func HandleVoidPayment(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "payment ID is required")
			return
		}

		pi, err := svc.VoidPayment(r.Context(), id, claims.MerchantID)
		if err != nil {
			if errors.Is(err, ErrPaymentNotVoidable) {
				api.RespondError(w, http.StatusBadRequest, "invalid_state", "not_voidable", err.Error())
				return
			}
			if errors.Is(err, pkgErr.ErrNotFound) {
				api.RespondError(w, http.StatusNotFound, "not_found", "payment_not_found", "payment not found")
				return
			}
			api.RespondError(w, http.StatusBadRequest, "validation_error", "void_failed", err.Error())
			return
		}

		api.RespondJSON(w, http.StatusOK, pi)
	}
}
