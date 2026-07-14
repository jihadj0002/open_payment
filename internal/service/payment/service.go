package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	pkgErr "github.com/openpayment/gateway/internal/pkg/errors"
)

var supportedCurrencies = map[string]bool{
	"BDT": true,
	"USD": true,
}

type Service struct {
	repo      PaymentRepository
	processor Processor
}

func NewService(repo PaymentRepository, processor Processor) *Service {
	return &Service{
		repo:      repo,
		processor: processor,
	}
}

func (s *Service) CreatePayment(ctx context.Context, merchantID string, req CreatePaymentRequest) (*PaymentIntent, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if !supportedCurrencies[req.Currency] {
		return nil, fmt.Errorf("unsupported currency: %s", req.Currency)
	}

	if req.CaptureMethod == "" {
		req.CaptureMethod = "automatic"
	}

	if req.IdempotencyKey != nil {
		existing, err := s.repo.GetByIdempotencyKey(ctx, *req.IdempotencyKey, merchantID)
		if err != nil && !errors.Is(err, pkgErr.ErrNotFound) {
			return nil, fmt.Errorf("check idempotency: %w", err)
		}
		if existing != nil {
			return existing, nil
		}
	}

	initialStatus := StatusCreated
	if req.Confirm {
		initialStatus = StatusPending
	}

	pi := &PaymentIntent{
		MerchantID:       merchantID,
		CustomerID:       req.CustomerID,
		Amount:           req.Amount,
		AmountCapturable: req.Amount,
		Currency:         req.Currency,
		Status:           initialStatus,
		CaptureMethod:    req.CaptureMethod,
		PaymentMethod:    req.PaymentMethod,
		IdempotencyKey:   req.IdempotencyKey,
		Description:      req.Description,
		Metadata:         req.Metadata,
		ReturnURL:        req.ReturnURL,
	}

	if err := s.repo.CreatePaymentIntent(ctx, pi); err != nil {
		return nil, fmt.Errorf("create payment intent: %w", err)
	}

	if req.Confirm && req.PaymentMethod == "card" {
		if err := s.ProcessPayment(ctx, pi); err != nil {
			return nil, fmt.Errorf("process payment: %w", err)
		}
	}

	return pi, nil
}

func (s *Service) GetPayment(ctx context.Context, id, merchantID string) (*PaymentIntent, error) {
	return s.repo.GetPaymentIntent(ctx, id, merchantID)
}

func (s *Service) ProcessPayment(ctx context.Context, pi *PaymentIntent) error {
	if !IsValidTransition(pi.Status, StatusProcessing) {
		return ErrInvalidStateTransition
	}

	pi.Status = StatusProcessing
	if err := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusProcessing); err != nil {
		return fmt.Errorf("update to processing: %w", err)
	}

	cardReq := CardRequest{
		Amount:         pi.Amount,
		Currency:       pi.Currency,
		MerchantRef:    pi.ID,
		IdempotencyKey: pi.ID,
	}

	procResp, err := s.processor.ProcessCard(cardReq)
	if err != nil {
		pi.Status = StatusFailed
		errMsg := err.Error()
		pi.ErrorMessage = &errMsg
		if updateErr := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed); updateErr != nil {
			return fmt.Errorf("update to failed: %w", updateErr)
		}
		return fmt.Errorf("processor error: %w", err)
	}

	respBytes, _ := json.Marshal(procResp)

	if !procResp.Success {
		pi.Status = StatusFailed
		pi.ErrorMessage = &procResp.Message
		if updateErr := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed); updateErr != nil {
			return fmt.Errorf("update to failed: %w", updateErr)
		}
		return fmt.Errorf("processor declined: %s", procResp.Message)
	}

	nextStatus := StatusAuthorized
	amountCapturable := pi.Amount
	amountReceived := int64(0)
	if pi.CaptureMethod == "automatic" {
		nextStatus = StatusCaptured
		amountCapturable = 0
		amountReceived = pi.Amount
	}

	pi.Status = nextStatus
	pi.AmountCapturable = amountCapturable
	pi.AmountReceived = amountReceived

	if err := s.repo.UpdatePaymentIntentCapture(ctx, pi.ID, amountCapturable, amountReceived, nextStatus); err != nil {
		return fmt.Errorf("update after process: %w", err)
	}

	tx := &Transaction{
		PaymentIntentID:   pi.ID,
		MerchantID:        pi.MerchantID,
		Type:              "capture",
		Status:            "succeeded",
		Amount:            pi.Amount,
		Currency:          pi.Currency,
		ProcessorRef:      &procResp.ProcessorRef,
		ProcessorResponse: respBytes,
		Fee:               procResp.Fee,
		NetAmount:         pi.Amount - procResp.Fee,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (s *Service) CapturePayment(ctx context.Context, id, merchantID string, amount *int64) (*PaymentIntent, error) {
	pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
	if err != nil {
		return nil, err
	}

	if pi.Status != StatusAuthorized {
		return nil, ErrPaymentNotCapturable
	}

	captureAmount := pi.AmountCapturable
	if amount != nil {
		if *amount > pi.AmountCapturable {
			return nil, fmt.Errorf("capture amount exceeds capturable amount")
		}
		captureAmount = *amount
	}

	newCapturable := pi.AmountCapturable - captureAmount
	newReceived := pi.AmountReceived + captureAmount
	newStatus := StatusCaptured
	if newCapturable == 0 {
		newStatus = StatusSucceeded
	}

	if err := s.repo.UpdatePaymentIntentCapture(ctx, id, newCapturable, newReceived, newStatus); err != nil {
		return nil, fmt.Errorf("update capture: %w", err)
	}

	tx := &Transaction{
		PaymentIntentID: id,
		MerchantID:      merchantID,
		Type:            "capture",
		Status:          "succeeded",
		Amount:          captureAmount,
		Currency:        pi.Currency,
		Fee:             0,
		NetAmount:       captureAmount,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create capture transaction: %w", err)
	}

	pi.Status = newStatus
	pi.AmountCapturable = newCapturable
	pi.AmountReceived = newReceived

	return pi, nil
}

func (s *Service) RefundPayment(ctx context.Context, id, merchantID string, amount int64, reason string) (*Transaction, error) {
	pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
	if err != nil {
		return nil, err
	}

	if pi.Status != StatusCaptured && pi.Status != StatusSucceeded {
		return nil, ErrPaymentNotRefundable
	}

	refundAmount := amount
	if refundAmount <= 0 {
		refundAmount = pi.AmountReceived
	}
	if refundAmount > pi.AmountReceived {
		return nil, fmt.Errorf("refund amount exceeds received amount")
	}

	newReceived := pi.AmountReceived - refundAmount
	newStatus := pi.Status
	if newReceived == 0 {
		newStatus = StatusRefunded
	}

	tx := &Transaction{
		PaymentIntentID: id,
		MerchantID:      merchantID,
		Type:            "refund",
		Status:          "succeeded",
		Amount:          refundAmount,
		Currency:        pi.Currency,
		NetAmount:       -refundAmount,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create refund transaction: %w", err)
	}

	if newStatus == StatusRefunded {
		if err := s.repo.UpdatePaymentIntentStatus(ctx, id, StatusRefunded); err != nil {
			return nil, fmt.Errorf("update to refunded: %w", err)
		}
	}

	pi.AmountReceived = newReceived
	pi.AmountCapturable = int64(math.Max(0, float64(pi.AmountCapturable-refundAmount)))

	return tx, nil
}

func (s *Service) VoidPayment(ctx context.Context, id, merchantID string) (*PaymentIntent, error) {
	pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
	if err != nil {
		return nil, err
	}

	if pi.Status != StatusAuthorized && pi.Status != StatusPending {
		return nil, ErrPaymentNotVoidable
	}

	if err := s.repo.UpdatePaymentIntentStatus(ctx, id, StatusCanceled); err != nil {
		return nil, fmt.Errorf("void payment: %w", err)
	}

	tx := &Transaction{
		PaymentIntentID: id,
		MerchantID:      merchantID,
		Type:            "void",
		Status:          "succeeded",
		Amount:          0,
		Currency:        pi.Currency,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create void transaction: %w", err)
	}

	pi.Status = StatusCanceled

	return pi, nil
}

func (s *Service) ListPayments(ctx context.Context, merchantID string, limit, offset int) ([]PaymentIntent, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListPaymentIntents(ctx, merchantID, limit, offset)
}
