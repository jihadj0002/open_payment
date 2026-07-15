package settlement

import (
	"context"
	"fmt"
	"time"

	"github.com/openpayment/gateway/internal/service/ledger"
)

type Service struct {
	repo      *Repository
	ledgerSvc *ledger.Service
}

func NewService(repo *Repository, ledgerSvc *ledger.Service) *Service {
	return &Service{repo: repo, ledgerSvc: ledgerSvc}
}

func (s *Service) TriggerSettlement(ctx context.Context, merchantID string, req SettlementRequest) (*Settlement, error) {
	if req.Currency == "" {
		req.Currency = "BDT"
	}

	amount, err := s.repo.GetUnsettledVolume(ctx, merchantID, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("get unsettled volume: %w", err)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("nothing to settle")
	}

	now := time.Now()
	fee := amount * 5 / 1000
	netAmount := amount - fee

	settlement := &Settlement{
		MerchantID:  merchantID,
		Amount:      amount,
		Currency:    req.Currency,
		Status:      StatusPending,
		Fee:         fee,
		NetAmount:   netAmount,
		PeriodStart: now.Add(-7 * 24 * time.Hour),
		PeriodEnd:   now,
	}

	if err := s.repo.CreateSettlement(ctx, settlement); err != nil {
		return nil, fmt.Errorf("create settlement: %w", err)
	}

	if err := s.ledgerSvc.RecordSettlement(ctx, merchantID, req.Currency, netAmount); err != nil {
		return nil, fmt.Errorf("record settlement in ledger: %w", err)
	}

	if err := s.repo.UpdateSettlementStatus(ctx, settlement.ID, StatusCompleted); err != nil {
		return nil, fmt.Errorf("update settlement status: %w", err)
	}

	nowTime := time.Now()
	settlement.Status = StatusCompleted
	settlement.CompletedAt = &nowTime

	return settlement, nil
}

func (s *Service) GetSettlement(ctx context.Context, id, merchantID string) (*Settlement, error) {
	settlement, err := s.repo.GetSettlement(ctx, id, merchantID)
	if err != nil {
		return nil, err
	}
	if settlement == nil {
		return nil, fmt.Errorf("settlement not found")
	}
	return settlement, nil
}

func (s *Service) ListSettlements(ctx context.Context, merchantID string, page, perPage int) ([]Settlement, error) {
	offset := (page - 1) * perPage
	return s.repo.ListSettlements(ctx, merchantID, perPage, offset)
}

func (s *Service) GetSettlementReport(ctx context.Context, merchantID, currency string, from, to time.Time) (*SettlementReport, error) {
	items, err := s.repo.GetSettlementReport(ctx, merchantID, currency, from, to)
	if err != nil {
		return nil, err
	}

	report := &SettlementReport{
		From:  from.Format("2006-01-02"),
		To:    to.Format("2006-01-02"),
		Items: items,
	}

	for _, item := range items {
		report.Totals.TotalVolume += item.TotalVolume
		report.Totals.TotalFees += item.TotalFees
		report.Totals.TotalNet += item.TotalNet
		report.Totals.TransactionCount += item.TransactionCount
	}

	return report, nil
}
