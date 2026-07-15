package customer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterCustomerRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Post("/customers", HandleCreateCustomer(svc))
	r.With(authMW).Get("/customers", HandleListCustomers(svc))
	r.With(authMW).Get("/customers/{id}", HandleGetCustomer(svc))
	r.With(authMW).Patch("/customers/{id}", HandleUpdateCustomer(svc))
	r.With(authMW).Post("/customers/{id}/payment_methods", HandleAttachPaymentMethod(svc))
	r.With(authMW).Get("/customers/{id}/payment_methods", HandleListPaymentMethods(svc))
	r.With(authMW).Delete("/customers/{id}/payment_methods/{pm_id}", HandleDetachPaymentMethod(svc))
}
