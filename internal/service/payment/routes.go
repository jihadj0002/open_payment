package payment

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterPaymentRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Post("/payments", HandleCreatePayment(svc))
	r.With(authMW).Get("/payments", HandleListPayments(svc))
	r.With(authMW).Get("/payments/{id}", HandleGetPayment(svc))
	r.With(authMW).Post("/payments/{id}/capture", HandleCapturePayment(svc))
	r.With(authMW).Post("/payments/{id}/refund", HandleRefundPayment(svc))
	r.With(authMW).Post("/payments/{id}/void", HandleVoidPayment(svc))
}
