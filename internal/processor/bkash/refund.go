package bkash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RefundRequest struct {
	PaymentID    string `json:"paymentId"`
	TrxID        string `json:"trxId"`
	RefundAmount string `json:"refundAmount"`
	SKU          string `json:"sku,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

type RefundResponse struct {
	OriginalTrxID            string `json:"originalTrxId"`
	RefundTrxID              string `json:"refundTrxId"`
	RefundTransactionStatus  string `json:"refundTransactionStatus"`
	OriginalTrxAmount        string `json:"originalTrxAmount"`
	RefundAmount             string `json:"refundAmount"`
	Currency                 string `json:"currency"`
	CompletedTime            string `json:"completedTime"`
	SKU                      string `json:"sku,omitempty"`
	Reason                   string `json:"reason,omitempty"`
}

func (a *Adapter) RefundPayment(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal refund request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/v2/tokenized-checkout/refund/payment/transaction", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create refund request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do refund request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read refund response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bKash refund failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result RefundResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse refund response: %w", err)
	}

	return &result, nil
}

type RefundStatusRequest struct {
	PaymentID string `json:"paymentId"`
	TrxID     string `json:"trxId"`
}

type RefundTransactionItem struct {
	RefundTrxID             string `json:"refundTrxId"`
	RefundTransactionStatus string `json:"refundTransactionStatus"`
	RefundAmount            string `json:"refundAmount"`
	CompletedTime           string `json:"completedTime"`
}

type RefundStatusResponse struct {
	OriginalTrxID            string                  `json:"originalTrxId"`
	OriginalTrxAmount        string                  `json:"originalTrxAmount"`
	OriginalTrxCompletedTime string                  `json:"originalTrxCompletedTime"`
	RefundTransactions       []RefundTransactionItem `json:"refundTransactions"`
}

func (a *Adapter) QueryRefundStatus(ctx context.Context, req RefundStatusRequest) (*RefundStatusResponse, error) {
	token := a.tokenManager.GetToken()
	if token == "" {
		return nil, fmt.Errorf("bKash token not available")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal refund status request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/v2/tokenized-checkout/refund/payment/status", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create refund status request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("X-App-Key", a.appKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do refund status request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read refund status response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bKash refund status failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result RefundStatusResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse refund status response: %w", err)
	}

	return &result, nil
}
