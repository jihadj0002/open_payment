package settlement

import "time"

type Settlement struct {
    ID          string     `json:"id"`
    MerchantID  string     `json:"merchant_id"`
    Amount      int64      `json:"amount"`
    Currency    string     `json:"currency"`
    Status      string     `json:"status"`
    Fee         int64      `json:"fee"`
    NetAmount   int64      `json:"net_amount"`
    PayoutRef   *string    `json:"payout_ref,omitempty"`
    PeriodStart time.Time  `json:"period_start"`
    PeriodEnd   time.Time  `json:"period_end"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}

type SettlementRequest struct {
    Currency string `json:"currency"`
}

type SettlementReportItem struct {
	Currency         string `json:"currency"`
	TotalVolume      int64  `json:"total_volume"`
	TotalFees        int64  `json:"total_fees"`
	TotalNet         int64  `json:"total_net"`
	TransactionCount int    `json:"transaction_count"`
}

type SettlementReport struct {
	From   string                `json:"from"`
	To     string                `json:"to"`
	Items  []SettlementReportItem `json:"items"`
	Totals struct {
		TotalVolume      int64 `json:"total_volume"`
		TotalFees        int64 `json:"total_fees"`
		TotalNet         int64 `json:"total_net"`
		TransactionCount int   `json:"transaction_count"`
	} `json:"totals"`
}

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)
