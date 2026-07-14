package fraud

import (
	"context"
	"fmt"
)

type Service struct {
	repo   *Repository
	engine *RuleEngine
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo:   repo,
		engine: NewRuleEngine(),
	}
}

func (s *Service) AssessRisk(ctx context.Context, merchantID string, req *FraudCheckRequest) (*FraudCheck, error) {
	cfg, err := s.repo.GetFraudConfig(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("get fraud config: %w", err)
	}

	if !cfg.Enabled {
		return &FraudCheck{
			PaymentIntentID: req.PaymentIntentID,
			MerchantID:      merchantID,
			Score:           0,
			Threshold:       50,
			Verdict:         VerdictPass,
			Flags:           []string{},
		}, nil
	}

	evalCtx := context.WithValue(ctx, fraudConfigKey, cfg)

	if req.CustomerID != "" {
		count, err := s.repo.GetRecentTransactionsByCustomer(ctx, req.CustomerID, merchantID, 60)
		if err != nil {
			return nil, fmt.Errorf("get velocity count: %w", err)
		}
		evalCtx = context.WithValue(evalCtx, velocityCountKey, count)
	}

	check, err := s.engine.Evaluate(evalCtx, req)
	if err != nil {
		return nil, fmt.Errorf("rule engine evaluation: %w", err)
	}

	check.PaymentIntentID = req.PaymentIntentID
	check.MerchantID = merchantID

	if err := s.repo.SaveFraudCheck(ctx, check); err != nil {
		return nil, fmt.Errorf("save fraud check: %w", err)
	}

	return check, nil
}

func (s *Service) GetFraudConfig(ctx context.Context, merchantID string) (*FraudConfig, error) {
	return s.repo.GetFraudConfig(ctx, merchantID)
}

func (s *Service) UpdateFraudConfig(ctx context.Context, merchantID string, cfg *FraudConfig) error {
	cfg.MerchantID = merchantID
	return s.repo.UpsertFraudConfig(ctx, cfg)
}
