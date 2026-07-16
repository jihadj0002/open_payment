package processor

const (
	ProviderBKash = "bkash"
	ProviderNagad = "nagad"
)

type CardRequest struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	CVV            string `json:"cvv"`
	Token          string `json:"token,omitempty"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	MerchantRef    string `json:"merchant_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

type BankRequest struct {
	AccountNumber  string `json:"account_number"`
	RoutingNumber  string `json:"routing_number"`
	AccountType    string `json:"account_type"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	MerchantRef    string `json:"merchant_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

type WalletRequest struct {
	WalletID       string `json:"wallet_id"`
	Provider       string `json:"provider"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	MerchantRef    string `json:"merchant_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

type WalletInitResponse struct {
	Success       bool   `json:"success"`
	RedirectURL   string `json:"redirect_url"`
	PaymentRef    string `json:"payment_ref"`
	ProviderRef   string `json:"provider_ref"`
	Status        string `json:"status"`
	Message       string `json:"message,omitempty"`
}

type ProcessorResponse struct {
	Success      bool   `json:"success"`
	ProcessorRef string `json:"processor_ref"`
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	Fee          int64  `json:"fee"`
	ProcessedAt  string `json:"processed_at"`
}
