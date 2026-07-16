package bkash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PayoutRequest struct {
	Amount                string `json:"amount"`
	Currency              string `json:"currency"`
	MerchantInvoiceNumber string `json:"merchantInvoiceNumber"`
	ReceiverMSISDN        string `json:"receiverMSISDN"`
}

type PayoutResponse struct {
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
	TrxID         string `json:"trxID"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
}

func (a *Adapter) SendPayout(ctx context.Context, req PayoutRequest) (*PayoutResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal payout request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/b2c/payment", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create payout request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do payout request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read payout response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bKash payout failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result PayoutResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse payout response: %w", err)
	}

	if result.StatusCode != "0000" {
		return nil, fmt.Errorf("bKash payout error [%s]: %s", result.StatusCode, result.StatusMessage)
	}

	return &result, nil
}
