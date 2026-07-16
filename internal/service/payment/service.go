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
	repo            PaymentRepository
	processor       Processor
	walletProviders map[string]WalletProvider
	bankProvider    BankProvider
	webhookSvc      WebhookService
	ledgerSvc       LedgerService
}

func NewService(repo PaymentRepository, processor Processor) *Service {
	return &Service{
		repo:            repo,
		processor:       processor,
		walletProviders: make(map[string]WalletProvider),
	}
}

func (s *Service) WithWalletProvider(p WalletProvider) *Service {
	s.walletProviders[p.GetProviderName()] = p
	return s
}

func (s *Service) WithBankProvider(p BankProvider) *Service {
	s.bankProvider = p
	return s
}

func (s *Service) GetWalletProvider(name string) WalletProvider {
	return s.walletProviders[name]
}

func (s *Service) WithWebhook(svc WebhookService) *Service {
	s.webhookSvc = svc
	return s
}

func (s *Service) WithLedger(svc LedgerService) *Service {
	s.ledgerSvc = svc
	return s
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

func (s *Service) GetPaymentByID(ctx context.Context, id string) (*PaymentIntent, error) {
	pi, err := s.repo.GetPaymentIntent(ctx, id, "")
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (s *Service) ProcessPayment(ctx context.Context, pi *PaymentIntent) error {
	if !IsValidTransition(pi.Status, StatusProcessing) {
		return ErrInvalidStateTransition
	}

	pi.Status = StatusProcessing
	if err := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusProcessing); err != nil {
		return fmt.Errorf("update to processing: %w", err)
	}

	switch pi.PaymentMethod {
	case "card":
		return s.processCardPayment(ctx, pi)
	case "wallet":
		return s.processWalletPayment(ctx, pi)
	case "bank_transfer":
		return s.processBankTransfer(ctx, pi)
	default:
		pi.Status = StatusFailed
		errMsg := fmt.Sprintf("unsupported payment method: %s", pi.PaymentMethod)
		pi.ErrorMessage = &errMsg
		if updateErr := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed); updateErr != nil {
			return fmt.Errorf("update to failed: %w", updateErr)
		}
		return fmt.Errorf("unsupported payment method: %s", pi.PaymentMethod)
	}
}

func (s *Service) processCardPayment(ctx context.Context, pi *PaymentIntent) error {
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

	return s.finalizePayment(ctx, pi, procResp)
}

func (s *Service) processWalletPayment(ctx context.Context, pi *PaymentIntent) error {
	if len(s.walletProviders) == 0 {
		pi.Status = StatusFailed
		errMsg := "no wallet providers configured"
		pi.ErrorMessage = &errMsg
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("no wallet providers configured")
	}

	providerName := pi.WalletProvider
	if providerName == "" {
		providerName = "bkash"
	}

	provider, ok := s.walletProviders[providerName]
	if !ok {
		pi.Status = StatusFailed
		errMsg := fmt.Sprintf("wallet provider %q not found", providerName)
		pi.ErrorMessage = &errMsg
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("wallet provider %q not found", providerName)
	}

	customerPhone := ""
	if pi.Metadata != nil {
		customerPhone = pi.Metadata["customer_phone"]
	}

	returnURL := ""
	cancelURL := ""
	if pi.ReturnURL != nil {
		returnURL = *pi.ReturnURL
	}

	initReq := WalletInitRequest{
		Amount:        pi.Amount,
		Currency:      pi.Currency,
		MerchantRef:   pi.ID,
		CustomerPhone: customerPhone,
		ReturnURL:     returnURL,
		CancelURL:     cancelURL,
	}

	initResp, err := provider.InitiatePayment(ctx, initReq)
	if err != nil {
		pi.Status = StatusFailed
		errMsg := err.Error()
		pi.ErrorMessage = &errMsg
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("%s init failed: %w", providerName, err)
	}

	if !initResp.Success {
		pi.Status = StatusFailed
		pi.ErrorMessage = &initResp.Message
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("%s init declined: %s", providerName, initResp.Message)
	}

	pi.ProviderRef = &initResp.ProviderRef
	pi.RedirectURL = &initResp.RedirectURL
	s.repo.UpdatePaymentIntentProvider(ctx, pi.ID, initResp.ProviderRef, initResp.RedirectURL)

	if err := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusWalletInitiated); err != nil {
		return fmt.Errorf("update to wallet_initiated: %w", err)
	}

	pi.Status = StatusWalletInitiated

	return nil
}

func (s *Service) processBankTransfer(ctx context.Context, pi *PaymentIntent) error {
	if s.bankProvider == nil {
		pi.Status = StatusFailed
		errMsg := "bank transfer provider not configured"
		pi.ErrorMessage = &errMsg
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("bank transfer provider not configured")
	}

	initReq := BankInitRequest{
		Amount:      pi.Amount,
		Currency:    pi.Currency,
		MerchantRef: pi.ID,
	}

	initResp, err := s.bankProvider.InitiatePayment(ctx, initReq)
	if err != nil {
		pi.Status = StatusFailed
		errMsg := err.Error()
		pi.ErrorMessage = &errMsg
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("bank init failed: %w", err)
	}

	if !initResp.Success {
		pi.Status = StatusFailed
		pi.ErrorMessage = &initResp.Message
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("bank init declined: %s", initResp.Message)
	}

	pi.RedirectURL = &initResp.RedirectURL
	s.repo.UpdatePaymentIntentRedirect(ctx, pi.ID, initResp.RedirectURL)
	pi.Status = StatusWalletInitiated
	return nil
}

func (s *Service) finalizePayment(ctx context.Context, pi *PaymentIntent, procResp *ProcessorResponse) error {
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

	s.dispatchWebhook(ctx, pi.MerchantID, "payment.success", pi)
	s.recordLedgerPayment(ctx, pi.MerchantID, tx.ID, pi.Currency, pi.Amount, procResp.Fee)

	return nil
}

func (s *Service) HandleWalletCallback(ctx context.Context, paymentID string, status string, providerRef string) error {
	pi, err := s.repo.GetPaymentIntent(ctx, paymentID, "")
	if err != nil {
		return fmt.Errorf("get payment intent: %w", err)
	}

	if pi.Status != StatusWalletInitiated {
		return fmt.Errorf("payment is not in wallet_initiated state (current: %s)", pi.Status)
	}

	if status == "cancel" {
		pi.Status = StatusCanceled
		if err := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusCanceled); err != nil {
			return fmt.Errorf("update to canceled: %w", err)
		}
		return nil
	}

	if err := s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusProcessing); err != nil {
		return fmt.Errorf("update to processing: %w", err)
	}
	pi.Status = StatusProcessing

	providerName := pi.WalletProvider
	if providerName == "" {
		providerName = "bkash"
	}

	provider, ok := s.walletProviders[providerName]
	if !ok {
		return fmt.Errorf("wallet provider %q not found", providerName)
	}

	procResp, err := provider.ExecutePayment(ctx, providerRef)
	if err != nil {
		pi.Status = StatusFailed
		errMsg := err.Error()
		pi.ErrorMessage = &errMsg
		s.repo.UpdatePaymentIntentStatus(ctx, pi.ID, StatusFailed)
		return fmt.Errorf("execute payment: %w", err)
	}

	return s.finalizePayment(ctx, pi, procResp)
}

func (s *Service) dispatchWebhook(ctx context.Context, merchantID, eventType string, data interface{}) {
	if s.webhookSvc != nil {
		s.webhookSvc.DispatchEvent(ctx, merchantID, eventType, data)
	}
}

func (s *Service) recordLedgerPayment(ctx context.Context, merchantID, transactionID, currency string, amount, fee int64) {
	if s.ledgerSvc != nil {
		if err := s.ledgerSvc.RecordPayment(ctx, merchantID, transactionID, currency, amount, fee); err != nil {
			_ = err
		}
	}
}

func (s *Service) recordLedgerRefund(ctx context.Context, merchantID, transactionID, currency string, amount int64) {
	if s.ledgerSvc != nil {
		if err := s.ledgerSvc.RecordRefund(ctx, merchantID, transactionID, currency, amount); err != nil {
			_ = err
		}
	}
}

func (s *Service) CapturePayment(ctx context.Context, id, merchantID string, amount *int64, idempotencyKey *string) (*PaymentIntent, error) {
	if idempotencyKey != nil {
		existing, err := s.checkTransactionIdempotency(ctx, *idempotencyKey, merchantID, "capture")
		if err != nil && !errors.Is(err, pkgErr.ErrNotFound) {
			return nil, fmt.Errorf("check idempotency: %w", err)
		}
		if existing != nil {
			pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
			if err != nil {
				return nil, err
			}
			return pi, nil
		}
	}

	pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
	if err != nil {
		return nil, err
	}

	if !IsValidTransition(pi.Status, StatusCaptured) {
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
		IdempotencyKey:  idempotencyKey,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create capture transaction: %w", err)
	}

	pi.Status = newStatus
	pi.AmountCapturable = newCapturable
	pi.AmountReceived = newReceived

	s.dispatchWebhook(ctx, pi.MerchantID, "payment.success", pi)
	s.recordLedgerPayment(ctx, pi.MerchantID, tx.ID, pi.Currency, captureAmount, 0)

	return pi, nil
}

func (s *Service) checkTransactionIdempotency(ctx context.Context, idempotencyKey, merchantID string, txType string) (*Transaction, error) {
	if idempotencyKey == "" {
		return nil, nil
	}
	existing, err := s.repo.GetTransactionByIdempotencyKey(ctx, idempotencyKey, merchantID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Type == txType {
		return existing, nil
	}
	return nil, nil
}

func (s *Service) RefundPayment(ctx context.Context, id, merchantID string, amount int64, reason string, idempotencyKey *string) (*Transaction, error) {
	if idempotencyKey != nil {
		existing, err := s.checkTransactionIdempotency(ctx, *idempotencyKey, merchantID, "refund")
		if err != nil && !errors.Is(err, pkgErr.ErrNotFound) {
			return nil, fmt.Errorf("check idempotency: %w", err)
		}
		if existing != nil {
			return existing, nil
		}
	}

	pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
	if err != nil {
		return nil, err
	}

	if !IsValidTransition(pi.Status, StatusRefunded) {
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
		IdempotencyKey:  idempotencyKey,
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

	s.dispatchWebhook(ctx, pi.MerchantID, "refund.completed", tx)
	s.recordLedgerRefund(ctx, pi.MerchantID, tx.ID, pi.Currency, refundAmount)

	return tx, nil
}

func (s *Service) VoidPayment(ctx context.Context, id, merchantID string, idempotencyKey *string) (*PaymentIntent, error) {
	if idempotencyKey != nil {
		existing, err := s.checkTransactionIdempotency(ctx, *idempotencyKey, merchantID, "void")
		if err != nil && !errors.Is(err, pkgErr.ErrNotFound) {
			return nil, fmt.Errorf("check idempotency: %w", err)
		}
		if existing != nil {
			pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
			if err != nil {
				return nil, err
			}
			return pi, nil
		}
	}

	pi, err := s.repo.GetPaymentIntent(ctx, id, merchantID)
	if err != nil {
		return nil, err
	}

	if !IsValidTransition(pi.Status, StatusCanceled) {
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
		IdempotencyKey:  idempotencyKey,
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
