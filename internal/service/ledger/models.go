package ledger

import "time"

type Entry struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	MerchantID    string    `json:"merchant_id"`
	EntryType     string    `json:"entry_type"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type Balance struct {
	MerchantID string `json:"merchant_id"`
	Currency   string `json:"currency"`
	Available  int64  `json:"available"`
	Pending    int64  `json:"pending"`
	Reserve    int64  `json:"reserve"`
}

type BalanceTransaction struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created"`
}

const (
	EntryTypePaymentIn  = "payment_in"
	EntryTypeRefundOut  = "refund_out"
	EntryTypeFee        = "fee"
	EntryTypeSettlement = "settlement"
	EntryTypePayout     = "payout"
	EntryTypeChargeback = "chargeback"
	EntryTypeAdjustment = "adjustment"
)
