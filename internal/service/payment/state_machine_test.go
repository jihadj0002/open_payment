//go:build unit

package payment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidTransition_Valid(t *testing.T) {
	validTransitions := []struct {
		from string
		to   string
	}{
		{StatusCreated, StatusPending},
		{StatusCreated, StatusFailed},
		{StatusCreated, StatusCanceled},
		{StatusPending, StatusProcessing},
		{StatusPending, StatusFailed},
		{StatusPending, StatusCanceled},
		{StatusProcessing, StatusAuthorized},
		{StatusProcessing, StatusCaptured},
		{StatusProcessing, StatusSucceeded},
		{StatusProcessing, StatusFailed},
		{StatusAuthorized, StatusCaptured},
		{StatusAuthorized, StatusCanceled},
		{StatusAuthorized, StatusFailed},
		{StatusCaptured, StatusSucceeded},
		{StatusCaptured, StatusRefunded},
		{StatusCreated, StatusWalletInitiated},
		{StatusPending, StatusWalletInitiated},
		{StatusWalletInitiated, StatusProcessing},
		{StatusWalletInitiated, StatusFailed},
		{StatusWalletInitiated, StatusCanceled},
		{StatusCreated, StatusBankPending},
		{StatusPending, StatusBankPending},
		{StatusBankPending, StatusProcessing},
		{StatusBankPending, StatusFailed},
		{StatusBankPending, StatusExpired},
	}

	for _, tt := range validTransitions {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			assert.True(t, IsValidTransition(tt.from, tt.to),
				"expected %s -> %s to be valid", tt.from, tt.to)
		})
	}
}

func TestIsValidTransition_Invalid(t *testing.T) {
	invalidTransitions := []struct {
		from string
		to   string
	}{
		{StatusCreated, StatusProcessing},
		{StatusCreated, StatusAuthorized},
		{StatusCreated, StatusCaptured},
		{StatusCreated, StatusSucceeded},
		{StatusCreated, StatusRefunded},
		{StatusPending, StatusCreated},
		{StatusPending, StatusAuthorized},
		{StatusPending, StatusCaptured},
		{StatusPending, StatusSucceeded},
		{StatusProcessing, StatusCreated},
		{StatusProcessing, StatusPending},
		{StatusProcessing, StatusCanceled},
		{StatusProcessing, StatusRefunded},
		{StatusAuthorized, StatusCreated},
		{StatusAuthorized, StatusPending},
		{StatusAuthorized, StatusProcessing},
		{StatusAuthorized, StatusSucceeded},
		{StatusAuthorized, StatusRefunded},
		{StatusCaptured, StatusCreated},
		{StatusCaptured, StatusPending},
		{StatusCaptured, StatusProcessing},
		{StatusCaptured, StatusAuthorized},
		{StatusCaptured, StatusCanceled},
		{StatusCaptured, StatusFailed},
	}

	for _, tt := range invalidTransitions {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			assert.False(t, IsValidTransition(tt.from, tt.to),
				"expected %s -> %s to be invalid", tt.from, tt.to)
		})
	}
}

func TestIsValidTransition_Terminal(t *testing.T) {
	terminalStates := []string{StatusSucceeded, StatusFailed, StatusCanceled, StatusRefunded}
	allStates := []string{
		StatusCreated, StatusPending, StatusProcessing,
		StatusAuthorized, StatusCaptured, StatusSucceeded,
		StatusFailed, StatusCanceled, StatusRefunded, StatusExpired,
		StatusWalletInitiated, StatusBankPending,
	}

	for _, terminal := range terminalStates {
		t.Run(terminal+" rejects all", func(t *testing.T) {
			for _, target := range allStates {
				assert.False(t, IsValidTransition(terminal, target),
					"expected terminal state %s -> %s to be invalid", terminal, target)
			}
		})
	}
}

func TestIsValidTransition_UnknownState(t *testing.T) {
	assert.False(t, IsValidTransition("unknown", StatusCreated))
	assert.False(t, IsValidTransition(StatusCreated, "unknown"))
	assert.False(t, IsValidTransition("unknown", "unknown"))
}

func TestIsValidTransition_RefundedGuard(t *testing.T) {
	assert.True(t, IsValidTransition(StatusCaptured, StatusRefunded))
	assert.True(t, IsValidTransition(StatusCaptured, StatusSucceeded))
	assert.False(t, IsValidTransition(StatusRefunded, StatusCaptured))
	assert.False(t, IsValidTransition(StatusSucceeded, StatusCaptured))
}

func TestIsValidTransition_AuthorizedTransitions(t *testing.T) {
	assert.True(t, IsValidTransition(StatusAuthorized, StatusCaptured))
	assert.True(t, IsValidTransition(StatusAuthorized, StatusCanceled))
	assert.True(t, IsValidTransition(StatusAuthorized, StatusFailed))
	assert.False(t, IsValidTransition(StatusAuthorized, StatusRefunded))
	assert.False(t, IsValidTransition(StatusAuthorized, StatusSucceeded))
}
