//go:build unit

package payment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreatePaymentIntent(ctx context.Context, pi *PaymentIntent) error {
	args := m.Called(ctx, pi)
	return args.Error(0)
}

func (m *mockRepo) GetPaymentIntent(ctx context.Context, id, merchantID string) (*PaymentIntent, error) {
	args := m.Called(ctx, id, merchantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PaymentIntent), args.Error(1)
}

func (m *mockRepo) UpdatePaymentIntentStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *mockRepo) UpdatePaymentIntentCapture(ctx context.Context, id string, amountCapturable, amountReceived int64, status string) error {
	args := m.Called(ctx, id, amountCapturable, amountReceived, status)
	return args.Error(0)
}

func (m *mockRepo) ListPaymentIntents(ctx context.Context, merchantID string, limit, offset int) ([]PaymentIntent, error) {
	args := m.Called(ctx, merchantID, limit, offset)
	return args.Get(0).([]PaymentIntent), args.Error(1)
}

func (m *mockRepo) GetByIdempotencyKey(ctx context.Context, key, merchantID string) (*PaymentIntent, error) {
	args := m.Called(ctx, key, merchantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PaymentIntent), args.Error(1)
}

func (m *mockRepo) CreateTransaction(ctx context.Context, tx *Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) ProcessCard(req CardRequest) (*ProcessorResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ProcessorResponse), args.Error(1)
}

func TestCreatePayment_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"
	req := CreatePaymentRequest{
		Amount:    1000,
		Currency:  "USD",
		PaymentMethod: "card",
	}

	repo.On("CreatePaymentIntent", mock.Anything, mock.MatchedBy(func(pi *PaymentIntent) bool {
		return pi.Amount == 1000 && pi.Currency == "USD" && pi.MerchantID == merchantID
	})).Return(nil).Once()

	pi, err := svc.CreatePayment(context.Background(), merchantID, req)
	assert.NoError(t, err)
	assert.NotNil(t, pi)
	assert.Equal(t, int64(1000), pi.Amount)
	assert.Equal(t, "USD", pi.Currency)
	assert.Equal(t, StatusCreated, pi.Status)
	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
}

func TestCreatePayment_InvalidAmount(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"

	tests := []struct {
		name   string
		amount int64
	}{
		{"zero", 0},
		{"negative", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreatePaymentRequest{
				Amount:    tt.amount,
				Currency:  "USD",
				PaymentMethod: "card",
			}
			pi, err := svc.CreatePayment(context.Background(), merchantID, req)
			assert.Error(t, err)
			assert.Nil(t, pi)
			assert.Contains(t, err.Error(), "amount must be greater than 0")
		})
	}
}

func TestCreatePayment_Idempotency(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"
	key := "idem_key_1"
	existing := &PaymentIntent{
		ID:         "pi_1",
		MerchantID: merchantID,
		Amount:     1000,
		Currency:   "USD",
		Status:     StatusCreated,
	}

	repo.On("GetByIdempotencyKey", mock.Anything, key, merchantID).Return(existing, nil).Once()

	req := CreatePaymentRequest{
		Amount:         1000,
		Currency:       "USD",
		PaymentMethod:  "card",
		IdempotencyKey: &key,
	}

	pi, err := svc.CreatePayment(context.Background(), merchantID, req)
	assert.NoError(t, err)
	assert.Equal(t, existing, pi)
	repo.AssertExpectations(t)
}

func TestCreatePayment_IdempotencyNotFound(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"
	key := "idem_key_new"

	repo.On("GetByIdempotencyKey", mock.Anything, key, merchantID).Return(nil, pkgErr.ErrNotFound).Once()
	repo.On("CreatePaymentIntent", mock.Anything, mock.Anything).Return(nil).Once()

	req := CreatePaymentRequest{
		Amount:         1000,
		Currency:       "USD",
		PaymentMethod:  "card",
		IdempotencyKey: &key,
	}

	pi, err := svc.CreatePayment(context.Background(), merchantID, req)
	assert.NoError(t, err)
	assert.NotNil(t, pi)
	assert.Equal(t, int64(1000), pi.Amount)
	repo.AssertExpectations(t)
}

func TestCreatePayment_UnsupportedCurrency(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	req := CreatePaymentRequest{
		Amount:    1000,
		Currency:  "EUR",
		PaymentMethod: "card",
	}

	pi, err := svc.CreatePayment(context.Background(), "merch_1", req)
	assert.Error(t, err)
	assert.Nil(t, pi)
	assert.Contains(t, err.Error(), "unsupported currency")
}

func TestCapturePayment_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"
	pi := &PaymentIntent{
		ID:               "pi_1",
		MerchantID:       merchantID,
		Amount:           1000,
		AmountCapturable: 1000,
		AmountReceived:   0,
		Currency:         "USD",
		Status:           StatusAuthorized,
		CaptureMethod:    "manual",
	}

	repo.On("GetPaymentIntent", mock.Anything, "pi_1", merchantID).Return(pi, nil).Once()
	repo.On("UpdatePaymentIntentCapture", mock.Anything, "pi_1", int64(0), int64(1000), StatusSucceeded).Return(nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == "pi_1" && tx.Type == "capture" && tx.Amount == int64(1000)
	})).Return(nil).Once()

	result, err := svc.CapturePayment(context.Background(), "pi_1", merchantID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, StatusSucceeded, result.Status)
	assert.Equal(t, int64(0), result.AmountCapturable)
	assert.Equal(t, int64(1000), result.AmountReceived)
	repo.AssertExpectations(t)
}

func TestCapturePayment_InvalidState(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"

	tests := []struct {
		name   string
		status string
	}{
		{"created", StatusCreated},
		{"pending", StatusPending},
		{"captured", StatusCaptured},
		{"succeeded", StatusSucceeded},
		{"failed", StatusFailed},
		{"canceled", StatusCanceled},
		{"refunded", StatusRefunded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi := &PaymentIntent{
				ID:     "pi_1",
				MerchantID: merchantID,
				Amount: 1000,
				Status: tt.status,
			}
			repo.On("GetPaymentIntent", mock.Anything, "pi_1", merchantID).Return(pi, nil).Once()

			result, err := svc.CapturePayment(context.Background(), "pi_1", merchantID, nil)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrPaymentNotCapturable)
			assert.Nil(t, result)
		})
	}
}

func TestRefundPayment_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"
	pi := &PaymentIntent{
		ID:             "pi_1",
		MerchantID:     merchantID,
		Amount:         1000,
		AmountReceived: 1000,
		Currency:       "USD",
		Status:         StatusCaptured,
	}

	repo.On("GetPaymentIntent", mock.Anything, "pi_1", merchantID).Return(pi, nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == "pi_1" && tx.Type == "refund" && tx.Amount == int64(1000)
	})).Return(nil).Once()
	repo.On("UpdatePaymentIntentStatus", mock.Anything, "pi_1", StatusRefunded).Return(nil).Once()

	tx, err := svc.RefundPayment(context.Background(), "pi_1", merchantID, 0, "customer request")
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "refund", tx.Type)
	assert.Equal(t, "succeeded", tx.Status)
	repo.AssertExpectations(t)
}

func TestRefundPayment_InvalidState(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"

	tests := []struct {
		name   string
		status string
	}{
		{"created", StatusCreated},
		{"pending", StatusPending},
		{"authorized", StatusAuthorized},
		{"failed", StatusFailed},
		{"canceled", StatusCanceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi := &PaymentIntent{
				ID:     "pi_1",
				MerchantID: merchantID,
				Amount: 1000,
				Status: tt.status,
			}
			repo.On("GetPaymentIntent", mock.Anything, "pi_1", merchantID).Return(pi, nil).Once()

			tx, err := svc.RefundPayment(context.Background(), "pi_1", merchantID, 100, "test")
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrPaymentNotRefundable)
			assert.Nil(t, tx)
		})
	}
}

func TestVoidPayment_Success(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"

	tests := []struct {
		name   string
		status string
	}{
		{"authorized", StatusAuthorized},
		{"pending", StatusPending},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi := &PaymentIntent{
				ID:         "pi_1",
				MerchantID: merchantID,
				Amount:     1000,
				Currency:   "USD",
				Status:     tt.status,
			}

			repo.On("GetPaymentIntent", mock.Anything, "pi_1", merchantID).Return(pi, nil).Once()
			repo.On("UpdatePaymentIntentStatus", mock.Anything, "pi_1", StatusCanceled).Return(nil).Once()
			repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
				return tx.PaymentIntentID == "pi_1" && tx.Type == "void"
			})).Return(nil).Once()

			result, err := svc.VoidPayment(context.Background(), "pi_1", merchantID)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, StatusCanceled, result.Status)
		})
	}
}

func TestVoidPayment_InvalidState(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"

	tests := []struct {
		name   string
		status string
	}{
		{"created", StatusCreated},
		{"captured", StatusCaptured},
		{"succeeded", StatusSucceeded},
		{"failed", StatusFailed},
		{"canceled", StatusCanceled},
		{"refunded", StatusRefunded},
		{"processing", StatusProcessing},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi := &PaymentIntent{
				ID:     "pi_1",
				MerchantID: merchantID,
				Amount: 1000,
				Status: tt.status,
			}
			repo.On("GetPaymentIntent", mock.Anything, "pi_1", merchantID).Return(pi, nil).Once()

			result, err := svc.VoidPayment(context.Background(), "pi_1", merchantID)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrPaymentNotVoidable)
			assert.Nil(t, result)
		})
	}
}

func TestProcessPayment_ProcessorSuccess(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"
	pi := &PaymentIntent{
		ID:            "pi_1",
		MerchantID:    merchantID,
		Amount:        1000,
		AmountCapturable: 1000,
		Currency:      "USD",
		Status:        StatusPending,
		CaptureMethod: "automatic",
	}

	procResp := &ProcessorResponse{
		Success:      true,
		ProcessorRef: "ref_123",
		Status:       "completed",
		Fee:          30,
		ProcessedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	repo.On("UpdatePaymentIntentStatus", mock.Anything, "pi_1", StatusProcessing).Return(nil).Once()
	proc.On("ProcessCard", mock.Anything).Return(procResp, nil).Once()
	repo.On("UpdatePaymentIntentCapture", mock.Anything, "pi_1", int64(0), int64(1000), StatusCaptured).Return(nil).Once()
	repo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
		return tx.PaymentIntentID == "pi_1" && tx.Type == "capture" && tx.Status == "succeeded" && tx.Fee == 30
	})).Return(nil).Once()

	err := svc.ProcessPayment(context.Background(), pi)
	assert.NoError(t, err)
	assert.Equal(t, StatusCaptured, pi.Status)
	assert.Equal(t, int64(0), pi.AmountCapturable)
	assert.Equal(t, int64(1000), pi.AmountReceived)
	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
}

func TestProcessPayment_ProcessorFailure(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	merchantID := "merch_1"
	pi := &PaymentIntent{
		ID:         "pi_1",
		MerchantID: merchantID,
		Amount:     1000,
		Currency:   "USD",
		Status:     StatusPending,
	}

	procResp := &ProcessorResponse{
		Success: false,
		Message: "card declined",
		Status:  "failed",
	}

	repo.On("UpdatePaymentIntentStatus", mock.Anything, "pi_1", StatusProcessing).Return(nil).Once()
	proc.On("ProcessCard", mock.Anything).Return(procResp, nil).Once()
	repo.On("UpdatePaymentIntentStatus", mock.Anything, "pi_1", StatusFailed).Return(nil).Once()

	err := svc.ProcessPayment(context.Background(), pi)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "processor declined")
	assert.Equal(t, StatusFailed, pi.Status)
	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
}

func TestProcessPayment_ProcessorServerError(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	pi := &PaymentIntent{
		ID:         "pi_2",
		MerchantID: "merch_1",
		Amount:     2000,
		Currency:   "USD",
		Status:     StatusPending,
	}

	repo.On("UpdatePaymentIntentStatus", mock.Anything, "pi_2", StatusProcessing).Return(nil).Once()
	proc.On("ProcessCard", mock.Anything).Return(nil, errors.New("connection refused")).Once()
	repo.On("UpdatePaymentIntentStatus", mock.Anything, "pi_2", StatusFailed).Return(nil).Once()

	err := svc.ProcessPayment(context.Background(), pi)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "processor error")
	assert.Equal(t, StatusFailed, pi.Status)
	repo.AssertExpectations(t)
	proc.AssertExpectations(t)
}

func TestProcessPayment_InvalidStateTransition(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	pi := &PaymentIntent{
		ID:     "pi_1",
		Status: StatusCreated,
	}

	err := svc.ProcessPayment(context.Background(), pi)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStateTransition)
}

func TestListPayments_DefaultPagination(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	proc := new(mockProcessor)
	svc := NewService(repo, proc)

	repo.On("ListPaymentIntents", mock.Anything, "merch_1", 10, 0).Return([]PaymentIntent{}, nil).Once()

	pis, err := svc.ListPayments(context.Background(), "merch_1", 0, -1)
	assert.NoError(t, err)
	assert.NotNil(t, pis)
	repo.AssertExpectations(t)
}
