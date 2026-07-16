package payment

import "context"

type PaymentRepository interface {
	CreatePaymentIntent(ctx context.Context, pi *PaymentIntent) error
	GetPaymentIntent(ctx context.Context, id, merchantID string) (*PaymentIntent, error)
	UpdatePaymentIntentStatus(ctx context.Context, id, status string) error
	UpdatePaymentIntentCapture(ctx context.Context, id string, amountCapturable, amountReceived int64, status string) error
	UpdatePaymentIntentProvider(ctx context.Context, id, providerRef, redirectURL string) error
	UpdatePaymentIntentRedirect(ctx context.Context, id, redirectURL string) error
	ListPaymentIntents(ctx context.Context, merchantID string, limit, offset int) ([]PaymentIntent, error)
	GetByIdempotencyKey(ctx context.Context, key, merchantID string) (*PaymentIntent, error)
	GetTransactionByIdempotencyKey(ctx context.Context, key, merchantID string) (*Transaction, error)
	GetPaymentIntentByProviderRef(ctx context.Context, providerRef string) (*PaymentIntent, error)
	CreateTransaction(ctx context.Context, tx *Transaction) error
	ListStatusHistory(ctx context.Context, paymentIntentID string) ([]StatusHistoryEntry, error)
}

type Processor interface {
	ProcessCard(req CardRequest) (*ProcessorResponse, error)
}

type WalletProvider interface {
	InitiatePayment(ctx context.Context, req WalletInitRequest) (*WalletInitResponse, error)
	ExecutePayment(ctx context.Context, paymentRef string) (*ProcessorResponse, error)
	GetProviderName() string
}

type WalletInitRequest struct {
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	MerchantRef   string `json:"merchant_ref"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	ReturnURL     string `json:"return_url"`
	CancelURL     string `json:"cancel_url"`
}

type WalletInitResponse struct {
	Success     bool   `json:"success"`
	RedirectURL string `json:"redirect_url"`
	PaymentRef  string `json:"payment_ref"`
	ProviderRef string `json:"provider_ref"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
}

type BankProvider interface {
	InitiatePayment(ctx context.Context, req BankInitRequest) (*BankInitResponse, error)
	GetProviderName() string
}

type BankInitRequest struct {
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	MerchantRef string `json:"merchant_ref"`
}

type BankInitResponse struct {
	Success     bool   `json:"success"`
	RedirectURL string `json:"redirect_url"`
	PaymentRef  string `json:"payment_ref"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
}

type WebhookService interface {
	DispatchEvent(ctx context.Context, merchantID, eventType string, data interface{})
}

type LedgerService interface {
	RecordPayment(ctx context.Context, merchantID, transactionID, currency string, amount, fee int64) error
	RecordRefund(ctx context.Context, merchantID, transactionID, currency string, amount int64) error
}
