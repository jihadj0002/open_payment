package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/service/merchant"
)

func RegisterAdminRoutes(r chi.Router, svc *Service, statsSvc *StatsService, merchantSvc *merchant.Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Group(func(r chi.Router) {
		r.Use(AdminOnly)

		r.Get("/admin/merchants", HandleListMerchants(svc))
		r.Get("/admin/merchants/{id}", HandleGetMerchant(svc))
		r.Patch("/admin/merchants/{id}", HandleUpdateMerchant(svc, merchantSvc))
		r.Post("/admin/merchants/{id}/approve", HandleApproveMerchant(svc))
		r.Post("/admin/merchants/{id}/suspend", HandleSuspendMerchant(svc))
		r.Post("/admin/merchants/{id}/terminate", HandleTerminateMerchant(svc))
		r.Post("/admin/merchants/{id}/reset-api-keys", HandleResetMerchantAPIKeys(merchantSvc))
		r.Get("/admin/merchants/{id}/stats", HandleGetMerchantStats(statsSvc))

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

	r.With(authMW).Group(func(r chi.Router) {
		r.Get("/merchants/{id}/stats", HandleGetMerchantStats(statsSvc))
		r.Get("/merchants/onboarding/status", HandleGetMerchantOnboardingStatus(svc))
		r.Post("/merchants/onboarding", HandleSubmitOnboarding(svc))
		r.Post("/merchants/onboarding/documents", HandleUploadOnboardingDocument(svc))
	})

	r.With(authMW).Group(func(r chi.Router) {
		r.Get("/merchants/fraud/rules", HandleListFraudRules(svc))
		r.Post("/merchants/fraud/rules", HandleCreateFraudRule(svc))
		r.Patch("/merchants/fraud/rules/{id}", HandleUpdateFraudRule(svc))
		r.Delete("/merchants/fraud/rules/{id}", HandleDeleteFraudRule(svc))
	})
}
