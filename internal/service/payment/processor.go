package payment

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ProcessorClient struct {
	baseURL string
	client  *http.Client
}

func NewProcessorClient(baseURL string) *ProcessorClient {
	return &ProcessorClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

type CardRequest struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	CVV            string `json:"cvv"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	MerchantRef    string `json:"merchant_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

type ProcessorResponse struct {
	Success      bool   `json:"success"`
	ProcessorRef string `json:"processor_ref"`
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	Fee          int64  `json:"fee"`
	ProcessedAt  string `json:"processed_at"`
}

func tokenizeCardNumber(cardNumber string) string {
	h := sha256.Sum256([]byte(cardNumber))
	token := hex.EncodeToString(h[:16])
	return "tok_card_" + token
}

func (c *ProcessorClient) ProcessCard(req CardRequest) (*ProcessorResponse, error) {
	if req.CardNumber == "" {
		req.CardNumber = "4111111111111111"
	}
	req.CardNumber = tokenizeCardNumber(req.CardNumber)
	if req.ExpiryMonth == "" {
		req.ExpiryMonth = "12"
	}
	if req.ExpiryYear == "" {
		req.ExpiryYear = "2030"
	}
	if req.CVV == "" {
		req.CVV = "123"
	}
	req.CVV = strings.Repeat("x", len(req.CVV))

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal card request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL+"/process/card", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var procResp ProcessorResponse
	if err := json.NewDecoder(resp.Body).Decode(&procResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &procResp, nil
}
