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
	StatusCreated         = "created"
	StatusPending         = "pending"
	StatusProcessing      = "processing"
	StatusAuthorized      = "authorized"
	StatusCaptured        = "captured"
	StatusSucceeded       = "succeeded"
	StatusFailed          = "failed"
	StatusCanceled        = "canceled"
	StatusRefunded        = "refunded"
	StatusExpired         = "expired"
	StatusWalletInitiated = "wallet_initiated"
)

var transitions = map[string]map[string]bool{
	StatusCreated: {
		StatusPending:         true,
		StatusFailed:          true,
		StatusCanceled:        true,
		StatusWalletInitiated: true,
	},
	StatusPending: {
		StatusProcessing:      true,
		StatusFailed:          true,
		StatusCanceled:        true,
		StatusWalletInitiated: true,
	},
	StatusProcessing: {
		StatusAuthorized: true,
		StatusFailed:     true,
		StatusCaptured:   true,
		StatusSucceeded:  true,
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
	StatusWalletInitiated: {
		StatusProcessing: true,
		StatusFailed:     true,
		StatusCanceled:   true,
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
