package card

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/openpayment/gateway/internal/processor"
)

type Adapter struct {
	baseURL    string
	httpClient *http.Client
}

func NewAdapter(baseURL string) *Adapter {
	return &Adapter{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *Adapter) ProcessCard(ctx context.Context, req processor.CardRequest) (*processor.ProcessorResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal card request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/process/card", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var procResp processor.ProcessorResponse
	if err := json.NewDecoder(resp.Body).Decode(&procResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &procResp, nil
}

func (a *Adapter) SupportsNetwork(network string) bool {
	switch strings.ToLower(network) {
	case "visa", "mastercard", "amex", "discover":
		return true
	default:
		return false
	}
}
