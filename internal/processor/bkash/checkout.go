package bkash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CreatePaymentRequest struct {
	Mode                   string `json:"mode"`
	PayerReference         string `json:"payerReference"`
	CallbackURL            string `json:"callbackURL"`
	Amount                 string `json:"amount"`
	Currency               string `json:"currency"`
	Intent                 string `json:"intent"`
	MerchantInvoiceNumber  string `json:"merchantInvoiceNumber"`
}

type CreatePaymentResponse struct {
	StatusCode            string `json:"statusCode"`
	StatusMessage         string `json:"statusMessage"`
	PaymentID             string `json:"paymentID"`
	BkashURL              string `json:"bkashURL"`
	CallbackURL           string `json:"callbackURL"`
	SuccessCallbackURL    string `json:"successCallbackURL"`
	FailureCallbackURL    string `json:"failureCallbackURL"`
	CancelledCallbackURL  string `json:"cancelledCallbackURL"`
	Amount                string `json:"amount"`
	Currency              string `json:"currency"`
	Intent                string `json:"intent"`
	TransactionStatus     string `json:"transactionStatus"`
	PaymentCreateTime     string `json:"paymentCreateTime"`
	MerchantInvoiceNumber string `json:"merchantInvoiceNumber"`
}

type ExecutePaymentResponse struct {
	StatusCode            string `json:"statusCode"`
	StatusMessage         string `json:"statusMessage"`
	PaymentID             string `json:"paymentID"`
	PayerReference        string `json:"payerReference"`
	CustomerMsisdn        string `json:"customerMsisdn"`
	TrxID                 string `json:"trxID"`
	Amount                string `json:"amount"`
	TransactionStatus     string `json:"transactionStatus"`
	PaymentExecuteTime    string `json:"paymentExecuteTime"`
	Currency              string `json:"currency"`
	Intent                string `json:"intent"`
	MerchantInvoiceNumber string `json:"merchantInvoiceNumber"`
	AgreementID           string `json:"agreementID,omitempty"`
}

type QueryPaymentResponse struct {
	StatusCode            string `json:"statusCode"`
	StatusMessage         string `json:"statusMessage"`
	PaymentID             string `json:"paymentID"`
	Amount                string `json:"amount"`
	Currency              string `json:"currency"`
	TransactionStatus     string `json:"transactionStatus"`
	TrxID                 string `json:"trxID,omitempty"`
	CustomerMsisdn        string `json:"customerMsisdn,omitempty"`
	MerchantInvoiceNumber string `json:"merchantInvoiceNumber"`
}

type Adapter struct {
	tokenManager *TokenManager
	appKey       string
	baseURL      string
	httpClient   *http.Client
	callbackBase string
}

func NewAdapter(tokenManager *TokenManager, appKey, baseURL, callbackBase string) *Adapter {
	return &Adapter{
		tokenManager: tokenManager,
		appKey:       appKey,
		baseURL:      strings.TrimRight(baseURL, "/"),
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		callbackBase: strings.TrimRight(callbackBase, "/"),
	}
}

func (a *Adapter) CreateCheckoutPayment(ctx context.Context, amount int64, merchantRef, customerPhone string) (*CreatePaymentResponse, error) {
	return a.CreateCheckoutPaymentWithIntent(ctx, amount, merchantRef, customerPhone, "sale")
}

func (a *Adapter) CreateCheckoutPaymentWithIntent(ctx context.Context, amount int64, merchantRef, customerPhone, intent string) (*CreatePaymentResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	callbackURL := a.callbackBase + "/v1/bkash/callback"

	req := CreatePaymentRequest{
		Mode:                  "0011",
		PayerReference:        customerPhone,
		CallbackURL:           callbackURL,
		Amount:                fmt.Sprintf("%.2f", float64(amount)/100),
		Currency:              "BDT",
		Intent:                intent,
		MerchantInvoiceNumber: merchantRef,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal create payment: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/tokenized/checkout/create", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do create payment: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read create payment response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bKash create payment failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result CreatePaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse create payment response: %w", err)
	}

	if result.StatusCode != "0000" {
		return nil, fmt.Errorf("bKash create payment error [%s]: %s", result.StatusCode, result.StatusMessage)
	}

	return &result, nil
}

func (a *Adapter) ExecutePayment(ctx context.Context, paymentID string) (*ExecutePaymentResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	body := fmt.Sprintf(`{"paymentID":"%s"}`, paymentID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/tokenized/checkout/execute", strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create execute request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do execute payment: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read execute response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bKash execute payment failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result ExecutePaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse execute response: %w", err)
	}

	if result.StatusCode != "0000" && result.StatusCode != "" {
		return nil, fmt.Errorf("bKash execute payment error [%s]: %s", result.StatusCode, result.StatusMessage)
	}

	return &result, nil
}

func (a *Adapter) QueryPayment(ctx context.Context, paymentID string) (*QueryPaymentResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	body := fmt.Sprintf(`{"paymentID":"%s"}`, paymentID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/tokenized/checkout/status", strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create status request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do status request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read status response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bKash query payment failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result QueryPaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse status response: %w", err)
	}

	return &result, nil
}

func (a *Adapter) CaptureAuthorizedPayment(ctx context.Context, paymentID string) (*ExecutePaymentResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.baseURL+"/checkout/payment/execute/"+paymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("create capture request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do capture request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read capture response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bKash capture failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result ExecutePaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse capture response: %w", err)
	}

	return &result, nil
}

func (a *Adapter) VoidAuthorizedPayment(ctx context.Context, paymentID string) error {
	token := a.tokenManager.GetToken()
	if token == "" {
		return fmt.Errorf("bKash token not available")
	}

	body := fmt.Sprintf(`{"paymentID":"%s"}`, paymentID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.baseURL+"/tokenized/checkout/execute", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("create void request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do void request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bKash void failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
