package fraud

import "time"

type FraudCheck struct {
	ID              string    `json:"id"`
	PaymentIntentID string    `json:"payment_intent_id"`
	MerchantID      string    `json:"merchant_id"`
	Score           int       `json:"score"`
	Threshold       int       `json:"threshold"`
	Verdict         string    `json:"verdict"`
	Flags           []string  `json:"flags"`
	CheckedAt       time.Time `json:"checked_at"`
}

type RuleResult struct {
	RuleName string `json:"rule_name"`
	Score    int    `json:"score"`
	Reason   string `json:"reason,omitempty"`
	Flagged  bool   `json:"flagged"`
}

type FraudConfig struct {
	MerchantID string `json:"merchant_id"`
	MaxAmount  int64  `json:"max_amount"`
	BlockVPN   bool   `json:"block_vpn"`
	MaxIpCount int    `json:"max_ip_count"`
	Enabled    bool   `json:"enabled"`
}

const (
	VerdictPass   = "pass"
	VerdictReview = "review"
	VerdictBlock  = "block"
)
