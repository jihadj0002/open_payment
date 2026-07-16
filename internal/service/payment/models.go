package payment

import (
	"encoding/json"
	"time"
)

type PaymentIntent struct {
	ID               string            `json:"id"`
	MerchantID       string            `json:"merchant_id"`
	CustomerID       *string           `json:"customer_id,omitempty"`
	Amount           int64             `json:"amount"`
	AmountCapturable int64             `json:"amount_capturable"`
	AmountReceived   int64             `json:"amount_received"`
	Currency         string            `json:"currency"`
	Status           string            `json:"status"`
	CaptureMethod    string            `json:"capture_method"`
	PaymentMethod    string            `json:"payment_method,omitempty"`
	WalletProvider   string            `json:"wallet_provider,omitempty"`
	IdempotencyKey   *string           `json:"idempotency_key,omitempty"`
	Description      *string           `json:"description,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	ClientSecret     string            `json:"client_secret"`
	ReturnURL        *string           `json:"return_url,omitempty"`
	CancelURL        *string           `json:"cancel_url,omitempty"`
	RedirectURL      *string           `json:"redirect_url,omitempty"`
	ProviderRef      *string           `json:"provider_ref,omitempty"`
	ErrorCode        *string           `json:"error_code,omitempty"`
	ErrorMessage     *string           `json:"error_message,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	StatusHistory    []StatusHistoryEntry `json:"status_history,omitempty"`
}

type Transaction struct {
	ID                string          `json:"id"`
	PaymentIntentID   string          `json:"payment_intent_id"`
	MerchantID        string          `json:"merchant_id"`
	Type              string          `json:"type"`
	Status            string          `json:"status"`
	Amount            int64           `json:"amount"`
	Currency          string          `json:"currency"`
	ProcessorRef      *string         `json:"processor_ref,omitempty"`
	ProcessorResponse json.RawMessage `json:"processor_response,omitempty"`
	Fee               int64           `json:"fee,omitempty"`
	NetAmount         int64           `json:"net_amount,omitempty"`
	IdempotencyKey    *string         `json:"idempotency_key,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}

type StatusHistoryEntry struct {
	ID              string    `json:"id"`
	PaymentIntentID string    `json:"payment_intent_id"`
	OldStatus       *string   `json:"old_status,omitempty"`
	NewStatus       string    `json:"new_status"`
	ChangedBy       string    `json:"changed_by"`
	Reason          *string   `json:"reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreatePaymentRequest struct {
	Amount         int64             `json:"amount"`
	Currency       string            `json:"currency"`
	PaymentMethod  string            `json:"payment_method"`
	CustomerID     *string           `json:"customer_id,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	Confirm        bool              `json:"confirm,omitempty"`
	CaptureMethod  string            `json:"capture_method,omitempty"`
	ReturnURL      *string           `json:"return_url,omitempty"`
	CancelURL      *string           `json:"cancel_url,omitempty"`
	IdempotencyKey *string           `json:"idempotency_key,omitempty"`
}

type CaptureRequest struct {
	AmountToCapture *int64 `json:"amount_to_capture,omitempty"`
}
