package payment

import "errors"

var (
	ErrInvalidStateTransition = errors.New("invalid payment state transition")
	ErrPaymentNotFound        = errors.New("payment not found")
	ErrPaymentNotCapturable   = errors.New("payment is not in a capturable state")
	ErrPaymentNotRefundable   = errors.New("payment is not in a refundable state")
	ErrPaymentNotVoidable     = errors.New("payment is not in a voidable state")
)

const (
	StatusCreated    = "created"
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusAuthorized = "authorized"
	StatusCaptured   = "captured"
	StatusSucceeded  = "succeeded"
	StatusFailed     = "failed"
	StatusCanceled   = "canceled"
	StatusRefunded   = "refunded"
	StatusExpired    = "expired"
)

var transitions = map[string]map[string]bool{
	StatusCreated: {
		StatusPending:  true,
		StatusFailed:   true,
		StatusCanceled: true,
	},
	StatusPending: {
		StatusProcessing: true,
		StatusFailed:     true,
		StatusCanceled:   true,
	},
	StatusProcessing: {
		StatusAuthorized: true,
		StatusFailed:     true,
	},
	StatusAuthorized: {
		StatusCaptured: true,
		StatusCanceled: true,
		StatusFailed:   true,
	},
	StatusCaptured: {
		StatusSucceeded: true,
		StatusRefunded:  true,
	},
	StatusSucceeded: {},
	StatusFailed:    {},
	StatusCanceled:  {},
	StatusRefunded:  {},
}

func IsValidTransition(from, to string) bool {
	allowed, ok := transitions[from]
	if !ok {
		return false
	}
	return allowed[to]
}
