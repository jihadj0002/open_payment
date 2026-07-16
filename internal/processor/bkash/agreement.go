package bkash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateAgreementRequest struct {
	Mode           string `json:"mode"`
	PayerReference string `json:"payerReference"`
	CallbackURL    string `json:"callbackURL"`
	Amount         string `json:"amount,omitempty"`
	Currency       string `json:"currency,omitempty"`
	Intent         string `json:"intent,omitempty"`
}

type CreateAgreementResponse struct {
	StatusCode          string `json:"statusCode"`
	StatusMessage       string `json:"statusMessage"`
	PaymentID           string `json:"paymentID"`
	BkashURL            string `json:"bkashURL"`
	CallbackURL         string `json:"callbackURL"`
	SuccessCallbackURL  string `json:"successCallbackURL"`
	FailureCallbackURL  string `json:"failureCallbackURL"`
	CancelledCallbackURL string `json:"cancelledCallbackURL"`
	PayerReference      string `json:"payerReference"`
	AgreementStatus     string `json:"agreementStatus"`
	AgreementCreateTime string `json:"agreementCreateTime"`
}

type ExecuteAgreementResponse struct {
	StatusCode           string `json:"statusCode"`
	StatusMessage        string `json:"statusMessage"`
	PaymentID            string `json:"paymentID"`
	AgreementID          string `json:"agreementID"`
	PayerReference       string `json:"payerReference"`
	CustomerMsisdn       string `json:"customerMsisdn"`
	AgreementExecuteTime string `json:"agreementExecuteTime"`
	AgreementStatus      string `json:"agreementStatus"`
}

func (a *Adapter) CreateAgreement(ctx context.Context, payerReference, callbackURL string) (*CreateAgreementResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	req := CreateAgreementRequest{
		Mode:           "0000",
		PayerReference: payerReference,
		CallbackURL:    callbackURL,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal create agreement: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/tokenized/checkout/create", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create agreement request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do create agreement: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read create agreement response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bKash create agreement failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result CreateAgreementResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse create agreement response: %w", err)
	}

	if result.StatusCode != "0000" {
		return nil, fmt.Errorf("bKash create agreement error [%s]: %s", result.StatusCode, result.StatusMessage)
	}

	return &result, nil
}

func (a *Adapter) ExecuteAgreement(ctx context.Context, paymentID string) (*ExecuteAgreementResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	body := []byte(fmt.Sprintf(`{"paymentID":"%s"}`, paymentID))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/tokenized/checkout/execute", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create execute agreement request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do execute agreement: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read execute agreement response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bKash execute agreement failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result ExecuteAgreementResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse execute agreement response: %w", err)
	}

	if result.StatusCode != "0000" && result.StatusCode != "" {
		return nil, fmt.Errorf("bKash execute agreement error [%s]: %s", result.StatusCode, result.StatusMessage)
	}

	return &result, nil
}
