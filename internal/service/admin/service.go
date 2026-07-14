package admin

import (
	"context"

	"github.com/openpayment/gateway/internal/service/merchant"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListMerchants(ctx context.Context, status, query string, page, perPage int) ([]MerchantListItem, error) {
	offset := (page - 1) * perPage
	return s.repo.ListMerchants(ctx, status, query, perPage, offset)
}

func (s *Service) GetMerchantDetail(ctx context.Context, id string) (*merchant.Merchant, error) {
	return s.repo.GetMerchantDetail(ctx, id)
}

func (s *Service) ApproveMerchant(ctx context.Context, id string) error {
	return s.repo.ApproveMerchant(ctx, id)
}

func (s *Service) SuspendMerchant(ctx context.Context, id, reason string) error {
	return s.repo.SuspendMerchant(ctx, id, reason)
}

func (s *Service) TerminateMerchant(ctx context.Context, id string) error {
	return s.repo.TerminateMerchant(ctx, id)
}

func (s *Service) ListTransactions(ctx context.Context, merchantID, status, paymentMethod string, page, perPage int) ([]TransactionListItem, error) {
	offset := (page - 1) * perPage
	return s.repo.ListTransactions(ctx, merchantID, status, paymentMethod, perPage, offset)
}

func (s *Service) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *Service) ListDisputes(ctx context.Context, status string, page, perPage int) ([]Dispute, error) {
	offset := (page - 1) * perPage
	return s.repo.ListDisputes(ctx, status, perPage, offset)
}

func (s *Service) GetDispute(ctx context.Context, id string) (*Dispute, error) {
	return s.repo.GetDispute(ctx, id)
}

func (s *Service) ResolveDispute(ctx context.Context, id, resolution, notes string) error {
	return s.repo.ResolveDispute(ctx, id, resolution, notes)
}

func (s *Service) ListFeeConfigs(ctx context.Context) ([]FeeConfig, error) {
	return s.repo.ListFeeConfigs(ctx)
}

func (s *Service) GetFeeConfig(ctx context.Context, id string) (*FeeConfig, error) {
	return s.repo.GetFeeConfig(ctx, id)
}

func (s *Service) CreateFeeConfig(ctx context.Context, f *FeeConfig) error {
	return s.repo.CreateFeeConfig(ctx, f)
}

func (s *Service) UpdateFeeConfig(ctx context.Context, id string, f *FeeConfig) error {
	return s.repo.UpdateFeeConfig(ctx, id, f)
}

func (s *Service) GetSystemConfig(ctx context.Context) (*SystemConfig, error) {
	return s.repo.GetSystemConfig(ctx)
}

func (s *Service) UpdateSystemConfig(ctx context.Context, cfg *SystemConfig) error {
	return s.repo.UpdateSystemConfig(ctx, cfg)
}

func (s *Service) ListAuditLogs(ctx context.Context, actorID, action, resourceType string, page, perPage int) ([]AuditLogEntry, error) {
	offset := (page - 1) * perPage
	return s.repo.ListAuditLogs(ctx, actorID, action, resourceType, perPage, offset)
}

func (s *Service) CreateAuditLog(ctx context.Context, e *AuditLogEntry) error {
	return s.repo.CreateAuditLog(ctx, e)
}
