//go:build integration

package payment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWebhookSvc struct {
	mock.Mock
}

func (m *mockWebhookSvc) DispatchEvent(ctx context.Context, merchantID, eventType string, data interface{}) {
	m.Called(ctx, merchantID, eventType, data)
}

func (m *mockWebhookSvc) SendWebhook(ctx context.Context, endpointID, payload string, secret string) error {
	args := m.Called(ctx, endpointID, payload, secret)
	return args.Error(0)
}

func (m *mockWebhookSvc) RetryPendingDeliveries(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type mockLedgerSvc struct {
	mock.Mock
}

func (m *mockLedgerSvc) RecordPayment(ctx context.Context, merchantID, transactionID, currency string, amount, fee int64) error {
	args := m.Called(ctx, merchantID, transactionID, currency, amount, fee)
	return args.Error(0)
}

func (m *mockLedgerSvc) RecordRefund(ctx context.Context, merchantID, transactionID, currency string, amount int64) error {
	args := m.Called(ctx, merchantID, transactionID, currency, amount)
	return args.Error(0)
}

func (m *mockLedgerSvc) GetBalance(ctx context.Context, merchantID string) (int64, error) {
	args := m.Called(ctx, merchantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockLedgerSvc) GetTransactions(ctx context.Context, merchantID string, limit, offset int) ([]interface{}, error) {
	args := m.Called(ctx, merchantID, limit, offset)
	return args.Get(0).([]interface{}), args.Error(1)
}

func TestIntegration_CreatePaymentFlow(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration"
	req := CreatePaymentRequest{
		Amount:    2000,
		Currency:  "USD",
		PaymentMethod: "card",
	}

	repo.On("CreatePaymentIntent", mock.Anything, mock.MatchedBy(func(pi *PaymentIntent) bool {
		return pi.Amount == 2000 && pi.Currency == "USD" && pi.MerchantID == merchantID
	})).Return(nil).Once()

	pi, err := svc.CreatePayment(context.Background(), merchantID, req)
	assert.NoError(t, err)
	assert.NotNil(t, pi)
	assert.Equal(t, int64(2000), pi.Amount)
	assert.Equal(t, "USD", pi.Currency)
	assert.Equal(t, StatusCreated, pi.Status)

	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
	webhook.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestIntegration_CreateThenCapturePayment(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration"
	paymentID := "pi_integration_" + uuid.New().String()

	pi := &PaymentIntent{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           5000,
		AmountCapturable: 5000,
		AmountReceived:   0,
		Currency:         "USD",
		Status:           StatusAuthorized,
		CaptureMethod:    "manual",
	}

	repo.On("GetPaymentIntent", mock.Anything, paymentID, merchantID).Return(pi, nil).Once()
	repo.On("UpdatePaymentIntentCapture", mock.Anything, paymentID, int64(0), int64(5000), StatusSucceeded).Return(nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == paymentID && tx.Type == "capture" && tx.Amount == int64(5000)
	})).Return(nil).Once()

	webhook.On("DispatchEvent", mock.Anything, merchantID, "payment.success", mock.Anything).Once()
	ledger.On("RecordPayment", mock.Anything, merchantID, mock.Anything, "USD", int64(5000), int64(0)).Return(nil).Once()

	result, err := svc.CapturePayment(context.Background(), paymentID, merchantID, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, StatusSucceeded, result.Status)
	assert.Equal(t, int64(0), result.AmountCapturable)
	assert.Equal(t, int64(5000), result.AmountReceived)

	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
	webhook.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestIntegration_CreateCaptureRefundFlow(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration"
	paymentID := "pi_integration_refund"

	pi := &PaymentIntent{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           3000,
		AmountCapturable: 0,
		AmountReceived:   3000,
		Currency:         "USD",
		Status:           StatusSucceeded,
	}

	repo.On("GetPaymentIntent", mock.Anything, paymentID, merchantID).Return(pi, nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == paymentID && tx.Type == "refund" && tx.Amount == int64(3000)
	})).Return(nil).Once()
	repo.On("UpdatePaymentIntentStatus", mock.Anything, paymentID, StatusRefunded).Return(nil).Once()

	webhook.On("DispatchEvent", mock.Anything, merchantID, "refund.completed", mock.Anything).Once()
	ledger.On("RecordRefund", mock.Anything, merchantID, mock.Anything, "USD", int64(3000)).Return(nil).Once()

	tx, err := svc.RefundPayment(context.Background(), paymentID, merchantID, 3000, "test refund", nil)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "refund", tx.Type)
	assert.Equal(t, "succeeded", tx.Status)
	assert.Equal(t, int64(3000), tx.Amount)

	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
	webhook.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestIntegration_CreateThenVoidPayment(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration"
	paymentID := "pi_integration_void"

	pi := &PaymentIntent{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           1500,
		AmountCapturable: 1500,
		AmountReceived:   0,
		Currency:         "USD",
		Status:           StatusAuthorized,
	}

	repo.On("GetPaymentIntent", mock.Anything, paymentID, merchantID).Return(pi, nil).Once()
	repo.On("UpdatePaymentIntentStatus", mock.Anything, paymentID, StatusCanceled).Return(nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == paymentID && tx.Type == "void" && tx.Status == "succeeded"
	})).Return(nil).Once()

	result, err := svc.VoidPayment(context.Background(), paymentID, merchantID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, StatusCanceled, result.Status)

	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
	webhook.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestIntegration_FullLifecycleFlow(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration_lifecycle"
	paymentID := "pi_lifecycle_1"

	pi := &PaymentIntent{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           10000,
		AmountCapturable: 10000,
		AmountReceived:   0,
		Currency:         "USD",
		Status:           StatusAuthorized,
		CaptureMethod:    "manual",
	}

	repo.On("GetPaymentIntent", mock.Anything, paymentID, merchantID).Return(pi, nil).Twice()
	repo.On("UpdatePaymentIntentCapture", mock.Anything, paymentID, int64(0), int64(10000), StatusSucceeded).Return(nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == paymentID && tx.Type == "capture"
	})).Return(nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == paymentID && tx.Type == "refund"
	})).Return(nil).Once()
	repo.On("UpdatePaymentIntentStatus", mock.Anything, paymentID, StatusRefunded).Return(nil).Once()

	webhook.On("DispatchEvent", mock.Anything, merchantID, "payment.success", mock.Anything).Once()
	webhook.On("DispatchEvent", mock.Anything, merchantID, "refund.completed", mock.Anything).Once()
	ledger.On("RecordPayment", mock.Anything, merchantID, mock.Anything, "USD", int64(10000), int64(0)).Return(nil).Once()
	ledger.On("RecordRefund", mock.Anything, merchantID, mock.Anything, "USD", int64(10000)).Return(nil).Once()

	result, err := svc.CapturePayment(context.Background(), paymentID, merchantID, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, StatusSucceeded, result.Status)

	tx, err := svc.RefundPayment(context.Background(), paymentID, merchantID, 10000, "full refund", nil)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "refund", tx.Type)
	assert.Equal(t, "succeeded", tx.Status)

	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
	webhook.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestIntegration_AutoCapturePaymentFlow(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration_auto"
	paymentID := "pi_auto_1"

	pi := &PaymentIntent{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           7500,
		AmountCapturable: 7500,
		Currency:         "USD",
		Status:           StatusPending,
		CaptureMethod:    "automatic",
	}

	procResp := &ProcessorResponse{
		Success:      true,
		ProcessorRef: "proc_ref_auto_1",
		Status:       "completed",
		Fee:          150,
		ProcessedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	repo.On("UpdatePaymentIntentStatus", mock.Anything, paymentID, StatusProcessing).Return(nil).Once()
	proc.On("ProcessCard", mock.Anything).Return(procResp, nil).Once()
	repo.On("UpdatePaymentIntentCapture", mock.Anything, paymentID, int64(0), int64(7500), StatusCaptured).Return(nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == paymentID && tx.Type == "capture" && tx.Fee == 150
	})).Return(nil).Once()

	webhook.On("DispatchEvent", mock.Anything, merchantID, "payment.success", mock.Anything).Once()
	ledger.On("RecordPayment", mock.Anything, merchantID, mock.Anything, "USD", int64(7500), int64(150)).Return(nil).Once()

	err := svc.ProcessPayment(context.Background(), pi)
	assert.NoError(t, err)
	assert.Equal(t, StatusCaptured, pi.Status)

	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
	webhook.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestIntegration_PaymentFailureFlow(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration_fail"
	paymentID := "pi_fail_1"

	pi := &PaymentIntent{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           5000,
		AmountCapturable: 5000,
		Currency:         "USD",
		Status:           StatusPending,
		CaptureMethod:    "automatic",
	}

	repo.On("UpdatePaymentIntentStatus", mock.Anything, paymentID, StatusProcessing).Return(nil).Once()
	proc.On("ProcessCard", mock.Anything).Return(nil, errors.New("processor failure")).Once()
	repo.On("UpdatePaymentIntentStatus", mock.Anything, paymentID, StatusFailed).Return(nil).Once()

	err := svc.ProcessPayment(context.Background(), pi)
	assert.Error(t, err)
	assert.Equal(t, StatusFailed, pi.Status)

	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
	webhook.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestIntegration_ConcurrentRequests(t *testing.T) {
	repo := new(mockRepo)
	proc := new(mockProcessor)
	webhook := new(mockWebhookSvc)
	ledger := new(mockLedgerSvc)
	svc := NewService(repo, proc).WithWebhook(webhook).WithLedger(ledger)

	merchantID := "merch_integration_concurrent"
	paymentID := "pi_concurrent_1"

	pi := &PaymentIntent{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           2000,
		AmountCapturable: 2000,
		AmountReceived:   0,
		Currency:         "USD",
		Status:           StatusAuthorized,
		CaptureMethod:    "manual",
	}

	repo.On("GetPaymentIntent", mock.Anything, paymentID, merchantID).Return(pi, nil)
	repo.On("UpdatePaymentIntentCapture", mock.Anything, paymentID, int64(0), int64(2000), StatusSucceeded).Return(nil)
	repo.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil)
	webhook.On("DispatchEvent", mock.Anything, merchantID, "payment.success", mock.Anything)
	ledger.On("RecordPayment", mock.Anything, merchantID, mock.Anything, "USD", int64(2000), int64(0)).Return(nil)

	t.Run("capture", func(t *testing.T) {
		result, err := svc.CapturePayment(context.Background(), paymentID, merchantID, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, StatusSucceeded, result.Status)
	})
}
