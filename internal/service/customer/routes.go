package customer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterCustomerRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Post("/v1/customers", HandleCreateCustomer(svc))
	r.With(authMW).Get("/v1/customers", HandleListCustomers(svc))
	r.With(authMW).Get("/v1/customers/{id}", HandleGetCustomer(svc))
	r.With(authMW).Patch("/v1/customers/{id}", HandleUpdateCustomer(svc))
	r.With(authMW).Post("/v1/customers/{id}/payment_methods", HandleAttachPaymentMethod(svc))
	r.With(authMW).Get("/v1/customers/{id}/payment_methods", HandleListPaymentMethods(svc))
	r.With(authMW).Delete("/v1/customers/{id}/payment_methods/{pm_id}", HandleDetachPaymentMethod(svc))
}
