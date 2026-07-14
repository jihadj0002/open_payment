package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterAdminRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Group(func(r chi.Router) {
		r.Use(AdminOnly)

		r.Get("/admin/merchants", HandleListMerchants(svc))
		r.Get("/admin/merchants/{id}", HandleGetMerchant(svc))
		r.Post("/admin/merchants/{id}/approve", HandleApproveMerchant(svc))
		r.Post("/admin/merchants/{id}/suspend", HandleSuspendMerchant(svc))
		r.Post("/admin/merchants/{id}/terminate", HandleTerminateMerchant(svc))

		r.Get("/admin/transactions", HandleListTransactions(svc))
		r.Get("/admin/transactions/{id}", HandleGetTransaction(svc))

		r.Get("/admin/disputes", HandleListDisputes(svc))
		r.Get("/admin/disputes/{id}", HandleGetDispute(svc))
		r.Post("/admin/disputes/{id}/resolve", HandleResolveDispute(svc))

		r.Get("/admin/fee_configs", HandleListFeeConfigs(svc))
		r.Post("/admin/fee_configs", HandleCreateFeeConfig(svc))
		r.Patch("/admin/fee_configs/{id}", HandleUpdateFeeConfig(svc))

		r.Get("/admin/config", HandleGetSystemConfig(svc))
		r.Patch("/admin/config", HandleUpdateSystemConfig(svc))

		r.Get("/admin/audit_logs", HandleListAuditLogs(svc))
	})
}
