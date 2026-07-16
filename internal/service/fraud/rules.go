package fraud

import (
	"context"
	"fmt"
	"strings"
)

type FraudCheckRequest struct {
	PaymentIntentID string
	MerchantID      string
	Amount          int64
	Currency        string
	CustomerID      string
	PaymentMethod   string
	IdempotencyKey  string
	IPAddress       string
}

type Rule interface {
	Name() string
	Evaluate(ctx context.Context, req *FraudCheckRequest) *RuleResult
}

type AmountThresholdRule struct{}

func (r *AmountThresholdRule) Name() string { return "amount_threshold" }

func (r *AmountThresholdRule) Evaluate(ctx context.Context, req *FraudCheckRequest) *RuleResult {
	cfg, ok := ctx.Value(fraudConfigKey).(*FraudConfig)
	if !ok || cfg == nil {
		return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
	}
	if req.Amount > cfg.MaxAmount {
		return &RuleResult{
			RuleName: r.Name(),
			Score:    30,
			Reason:   fmt.Sprintf("amount %d exceeds max %d", req.Amount, cfg.MaxAmount),
			Flagged:  true,
		}
	}
	return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
}

type HighVelocityRule struct{}

func (r *HighVelocityRule) Name() string { return "high_velocity" }

func (r *HighVelocityRule) Evaluate(ctx context.Context, req *FraudCheckRequest) *RuleResult {
	count, ok := ctx.Value(velocityCountKey).(int)
	if !ok {
		return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
	}
	if count > 10 {
		return &RuleResult{
			RuleName: r.Name(),
			Score:    25,
			Reason:   fmt.Sprintf("%d transactions from customer in last hour", count),
			Flagged:  true,
		}
	}
	if count > 5 {
		return &RuleResult{
			RuleName: r.Name(),
			Score:    10,
			Reason:   fmt.Sprintf("%d transactions from customer in last hour", count),
			Flagged:  true,
		}
	}
	return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
}

type NewCustomerRule struct{}

func (r *NewCustomerRule) Name() string { return "new_customer" }

func (r *NewCustomerRule) Evaluate(ctx context.Context, req *FraudCheckRequest) *RuleResult {
	isNew, ok := ctx.Value(isNewCustomerKey).(bool)
	if !ok || !isNew {
		return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
	}
	return &RuleResult{
		RuleName: r.Name(),
		Score:    15,
		Reason:   "customer created less than 1 hour ago",
		Flagged:  true,
	}
}

type CardBINCheckRule struct{}

func (r *CardBINCheckRule) Name() string { return "card_bin_check" }

var highRiskBINs = []string{"400000", "411111", "444444", "401288"}

func (r *CardBINCheckRule) Evaluate(ctx context.Context, req *FraudCheckRequest) *RuleResult {
	bin := extractBIN(req.PaymentMethod)
	if bin == "" {
		return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
	}
	for _, prefix := range highRiskBINs {
		if strings.HasPrefix(bin, prefix) {
			return &RuleResult{
				RuleName: r.Name(),
				Score:    10,
				Reason:   fmt.Sprintf("card BIN %s matches high-risk range", bin),
				Flagged:  true,
			}
		}
	}
	return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
}

func extractBIN(paymentMethod string) string {
	cleaned := strings.NewReplacer(" ", "", "-", "").Replace(paymentMethod)
	if len(cleaned) >= 6 {
		return cleaned[:6]
	}
	return cleaned
}

type IdempotencyReuseRule struct{}

func (r *IdempotencyReuseRule) Name() string { return "idempotency_reuse" }

func (r *IdempotencyReuseRule) Evaluate(ctx context.Context, req *FraudCheckRequest) *RuleResult {
	isReused, ok := ctx.Value(idempotencyReusedKey).(bool)
	if !ok || !isReused {
		return &RuleResult{RuleName: r.Name(), Score: 0, Flagged: false}
	}
	return &RuleResult{
		RuleName: r.Name(),
		Score:    20,
		Reason:   "idempotency key was previously used",
		Flagged:  true,
	}
}

type contextKey string

const (
	fraudConfigKey       contextKey = "fraud_config"
	velocityCountKey     contextKey = "velocity_count"
	isNewCustomerKey     contextKey = "is_new_customer"
	idempotencyReusedKey contextKey = "idempotency_reused"
)

type RuleEngine struct {
	rules []Rule
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: []Rule{
			&AmountThresholdRule{},
			&HighVelocityRule{},
			&NewCustomerRule{},
			&CardBINCheckRule{},
			&IdempotencyReuseRule{},
		},
	}
}

func (e *RuleEngine) Evaluate(ctx context.Context, req *FraudCheckRequest) (*FraudCheck, error) {
	if req == nil {
		return nil, fmt.Errorf("fraud check request is nil")
	}

	totalScore := 0
	threshold := 50
	var flags []string

	for _, rule := range e.rules {
		result := rule.Evaluate(ctx, req)
		if result.Flagged {
			totalScore += result.Score
			flags = append(flags, result.RuleName)
		}
	}

	verdict := determineVerdict(totalScore, threshold)

	return &FraudCheck{
		Score:     totalScore,
		Threshold: threshold,
		Verdict:   verdict,
		Flags:     flags,
	}, nil
}

func determineVerdict(score, threshold int) string {
	if score >= threshold {
		return VerdictBlock
	}
	if score >= 21 {
		return VerdictReview
	}
	return VerdictPass
}
