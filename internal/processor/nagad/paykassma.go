package nagad

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PaykassmaCreatePaymentRequest struct {
	Currency string `json:"currency"`
	UserID   string `json:"user_id,omitempty"`
	UserLabel string `json:"user_label,omitempty"`
}

type PaykassmaCreatePaymentResponse struct {
	Status           string `json:"status"`
	StatusCode       string `json:"statusCode"`
	PaymentRefID     string `json:"paymentRefId"`
	Amount           string `json:"amount"`
	ClientMobileNo   string `json:"clientMobileNo"`
	MerchantMobileNo string `json:"merchantMobileNo"`
	OrderDateTime    string `json:"orderDateTime"`
	Message          string `json:"message,omitempty"`
}

type PaykassmaActivateRequest struct {
	WalletType string `json:"wallet_type"`
	Key1       string `json:"key1"`
	Amount     string `json:"amount"`
}

type PaykassmaActivateResponse struct {
	Status           string `json:"status"`
	StatusCode       string `json:"statusCode"`
	PaymentRefID     string `json:"paymentRefId"`
	TransactionID    string `json:"transactionId,omitempty"`
	IssuerPaymentRefNo string `json:"issuerPaymentRefNo,omitempty"`
	Amount           string `json:"amount"`
	Message          string `json:"message,omitempty"`
}

type PaykassmaPostbackPayload struct {
	MerchantID         string `json:"merchantId"`
	OrderID           string `json:"orderId"`
	PaymentRefID      string `json:"paymentRefId"`
	Amount            string `json:"amount"`
	ClientMobileNo    string `json:"clientMobileNo"`
	TransactionStatus string `json:"transactionStatus"`
	Signature         string `json:"signature"`
}

func (a *Adapter) PaykassmaCreatePayment(ctx context.Context, req PaykassmaCreatePaymentRequest) (*PaykassmaCreatePaymentResponse, error) {
	payload := map[string]interface{}{
		"merchantId": a.creds.MerchantID,
		"currency":   req.Currency,
	}

	if req.UserID != "" {
		payload["userId"] = req.UserID
	}
	if req.UserLabel != "" {
		payload["userLabel"] = req.UserLabel
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal paykassma create payload: %w", err)
	}

	signature, err := a.sign(string(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("sign paykassma create: %w", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(payloadBytes)

	createReq := map[string]string{
		"merchantId": a.creds.MerchantID,
		"payload":    encodedPayload,
		"signature":  signature,
	}

	body, err := json.Marshal(createReq)
	if err != nil {
		return nil, fmt.Errorf("marshal paykassma create request wrapper: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.baseURL+"/api/v1/transaction/create/nagad", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create paykassma request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Merchant-ID", a.creds.MerchantID)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do paykassma create: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read paykassma create response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("paykassma create failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result PaykassmaCreatePaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse paykassma create response: %w", err)
	}

	return &result, nil
}

func (a *Adapter) PaykassmaActivatePayment(ctx context.Context, req PaykassmaActivateRequest) (*PaykassmaActivateResponse, error) {
	payload := map[string]interface{}{
		"merchantId":  a.creds.MerchantID,
		"wallet_type": req.WalletType,
		"key1":        req.Key1,
		"amount":      req.Amount,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal paykassma activate payload: %w", err)
	}

	signature, err := a.sign(string(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("sign paykassma activate: %w", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(payloadBytes)

	activateReq := map[string]string{
		"merchantId": a.creds.MerchantID,
		"payload":    encodedPayload,
		"signature":  signature,
	}

	body, err := json.Marshal(activateReq)
	if err != nil {
		return nil, fmt.Errorf("marshal paykassma activate request wrapper: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.baseURL+"/api/v1/transaction/activate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create activate request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Merchant-ID", a.creds.MerchantID)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do activate: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read activate response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("paykassma activate failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result PaykassmaActivateResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse activate response: %w", err)
	}

	return &result, nil
}

func (a *Adapter) VerifyPaykassmaPostback(payload PaykassmaPostbackPayload) error {
	data := fmt.Sprintf("%s%s%s%s%s",
		payload.MerchantID,
		payload.OrderID,
		payload.PaymentRefID,
		payload.Amount,
		payload.TransactionStatus,
	)
	return a.verifySignature(data, payload.Signature)
}
