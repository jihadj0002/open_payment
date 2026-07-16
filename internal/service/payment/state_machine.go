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
	StatusBankPending     = "bank_pending"
)

var transitions = map[string]map[string]bool{
	StatusCreated: {
		StatusPending:         true,
		StatusFailed:          true,
		StatusCanceled:        true,
		StatusWalletInitiated: true,
		StatusBankPending:     true,
	},
	StatusPending: {
		StatusProcessing:      true,
		StatusFailed:          true,
		StatusCanceled:        true,
		StatusWalletInitiated: true,
		StatusBankPending:     true,
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
	StatusBankPending: {
		StatusProcessing: true,
		StatusFailed:     true,
		StatusExpired:    true,
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
