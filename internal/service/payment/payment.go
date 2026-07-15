package payment

import "context"

type PaymentRepository interface {
	CreatePaymentIntent(ctx context.Context, pi *PaymentIntent) error
	GetPaymentIntent(ctx context.Context, id, merchantID string) (*PaymentIntent, error)
	UpdatePaymentIntentStatus(ctx context.Context, id, status string) error
	UpdatePaymentIntentCapture(ctx context.Context, id string, amountCapturable, amountReceived int64, status string) error
	ListPaymentIntents(ctx context.Context, merchantID string, limit, offset int) ([]PaymentIntent, error)
	GetByIdempotencyKey(ctx context.Context, key, merchantID string) (*PaymentIntent, error)
	GetTransactionByIdempotencyKey(ctx context.Context, key, merchantID string) (*Transaction, error)
	CreateTransaction(ctx context.Context, tx *Transaction) error
	ListStatusHistory(ctx context.Context, paymentIntentID string) ([]StatusHistoryEntry, error)
}

type Processor interface {
	ProcessCard(req CardRequest) (*ProcessorResponse, error)
}

type WebhookService interface {
	DispatchEvent(ctx context.Context, merchantID, eventType string, data interface{})
}

type LedgerService interface {
	RecordPayment(ctx context.Context, merchantID, transactionID, currency string, amount, fee int64) error
	RecordRefund(ctx context.Context, merchantID, transactionID, currency string, amount int64) error
}
