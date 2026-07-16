package payment

import (
	"context"
	"fmt"

	"github.com/openpayment/gateway/internal/processor/bkash"
	"github.com/openpayment/gateway/internal/processor/nagad"
	"github.com/openpayment/gateway/internal/processor/wallet"
)

type BkashWalletProvider struct {
	adapter *bkash.Adapter
}

func NewBkashWalletProvider(adapter *bkash.Adapter) *BkashWalletProvider {
	return &BkashWalletProvider{adapter: adapter}
}

func (p *BkashWalletProvider) GetProviderName() string {
	return "bkash"
}

func (p *BkashWalletProvider) InitiatePayment(ctx context.Context, req WalletInitRequest) (*WalletInitResponse, error) {
	customerPhone := req.CustomerPhone
	if customerPhone == "" {
		customerPhone = "01" 
	}

	createResp, err := p.adapter.CreateCheckoutPayment(ctx, req.Amount, req.MerchantRef, customerPhone)
	if err != nil {
		return nil, fmt.Errorf("bkash create payment: %w", err)
	}

	return &WalletInitResponse{
		Success:     true,
		RedirectURL: createResp.BkashURL,
		PaymentRef:  req.MerchantRef,
		ProviderRef: createResp.PaymentID,
		Status:      "initiated",
	}, nil
}

func (p *BkashWalletProvider) ExecutePayment(ctx context.Context, paymentRef string) (*ProcessorResponse, error) {
	executeResp, err := p.adapter.ExecutePayment(ctx, paymentRef)
	if err != nil {
		return nil, fmt.Errorf("bkash execute payment: %w", err)
	}

	success := executeResp.TransactionStatus == "Completed"
	status := "succeeded"
	if !success {
		status = "failed"
	}

	return &ProcessorResponse{
		Success:      success,
		ProcessorRef: executeResp.TrxID,
		Status:       status,
		Message:      executeResp.StatusMessage,
		Fee:          0,
		ProcessedAt:  executeResp.PaymentExecuteTime,
	}, nil
}

type NagadWalletProvider struct {
	adapter *nagad.Adapter
}

func NewNagadWalletProvider(adapter *nagad.Adapter) *NagadWalletProvider {
	return &NagadWalletProvider{adapter: adapter}
}

func (p *NagadWalletProvider) GetProviderName() string {
	return "nagad"
}

func (p *NagadWalletProvider) InitiatePayment(ctx context.Context, req WalletInitRequest) (*WalletInitResponse, error) {
	nagadReq := nagad.PaymentRequest{
		Amount:        fmt.Sprintf("%.2f", float64(req.Amount)/100),
		Currency:      "BDT",
		CustomerMobile: req.CustomerPhone,
	}

	initResp, err := p.adapter.InitializePayment(ctx, nagadReq)
	if err != nil {
		return nil, fmt.Errorf("nagad init payment: %w", err)
	}

	redirectURL := ""
	if initResp.CallbackURL != "" {
		redirectURL = initResp.CallbackURL
	}

	return &WalletInitResponse{
		Success:     true,
		RedirectURL: redirectURL,
		PaymentRef:  req.MerchantRef,
		ProviderRef: initResp.PaymentRefID,
		Status:      "initiated",
	}, nil
}

func (p *NagadWalletProvider) ExecutePayment(ctx context.Context, paymentRef string) (*ProcessorResponse, error) {
	completeResp, err := p.adapter.CompletePayment(ctx, paymentRef)
	if err != nil {
		return nil, fmt.Errorf("nagad complete payment: %w", err)
	}

	success := completeResp.Status == "Success" || completeResp.StatusCode == "Success"
	status := "succeeded"
	if !success {
		status = "failed"
	}

	processorRef := completeResp.IssuerPaymentRefNo
	if processorRef == "" {
		processorRef = completeResp.PaymentRefID
	}

	return &ProcessorResponse{
		Success:      success,
		ProcessorRef: processorRef,
		Status:       status,
		Message:      completeResp.Status,
		Fee:          0,
		ProcessedAt:  completeResp.IssuerPaymentDateTime,
	}, nil
}

type MockWalletProvider struct {
	adapter *wallet.MockAdapter
}

func NewMockWalletProvider(adapter *wallet.MockAdapter) *MockWalletProvider {
	return &MockWalletProvider{adapter: adapter}
}

func (p *MockWalletProvider) GetProviderName() string {
	return "mock"
}

func (p *MockWalletProvider) InitiatePayment(ctx context.Context, req WalletInitRequest) (*WalletInitResponse, error) {
	redirectURL, paymentRef, providerRef, err := p.adapter.InitiatePayment(ctx, req.Amount, req.Currency, req.MerchantRef, req.CustomerPhone, req.ReturnURL, req.CancelURL)
	if err != nil {
		return nil, fmt.Errorf("mock init payment: %w", err)
	}

	return &WalletInitResponse{
		Success:     true,
		RedirectURL: redirectURL,
		PaymentRef:  paymentRef,
		ProviderRef: providerRef,
		Status:      "initiated",
	}, nil
}

func (p *MockWalletProvider) ExecutePayment(ctx context.Context, paymentRef string) (*ProcessorResponse, error) {
	success, processorRef, status, message, err := p.adapter.ExecutePayment(ctx, paymentRef)
	if err != nil {
		return nil, fmt.Errorf("mock execute payment: %w", err)
	}

	return &ProcessorResponse{
		Success:      success,
		ProcessorRef: processorRef,
		Status:       status,
		Message:      message,
		Fee:          0,
		ProcessedAt:  "",
	}, nil
}
