package webhook

import "time"

type Endpoint struct {
	ID         string    `json:"id"`
	MerchantID string    `json:"merchant_id"`
	Event      string    `json:"event"`
	URL        string    `json:"url"`
	Secret     string    `json:"-"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateEndpointRequest struct {
	URL   string `json:"url"`
	Event string `json:"event"`
}

type Delivery struct {
	ID           string     `json:"id"`
	WebhookID    string     `json:"webhook_id"`
	Event        string     `json:"event"`
	Payload      []byte     `json:"-"`
	Status       string     `json:"status"`
	Attempt      int        `json:"attempt"`
	MaxAttempts  int        `json:"max_attempts"`
	ResponseCode *int       `json:"response_code,omitempty"`
	NextAttemptAt *time.Time `json:"next_attempt_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type WebhookEvent struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Created int64       `json:"created"`
	Data    interface{} `json:"data"`
}

const (
	EventPaymentSuccess  = "payment.success"
	EventPaymentFailed   = "payment.failed"
	EventPaymentPending  = "payment.pending"
	EventRefundCompleted = "refund.completed"
	EventChargebackCreated = "chargeback.created"
)

var validEvents = map[string]bool{
	EventPaymentSuccess:    true,
	EventPaymentFailed:     true,
	EventPaymentPending:    true,
	EventRefundCompleted:   true,
	EventChargebackCreated: true,
}
