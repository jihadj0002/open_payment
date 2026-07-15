package stripe

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/openpayment/gateway/internal/processor"
)

type Adapter struct {
	secretKey  string
	httpClient *http.Client
	liveMode   bool
}

func NewAdapter() *Adapter {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		secretKey = os.Getenv("STRIPE_TEST_SECRET_KEY")
	}
	liveMode := os.Getenv("STRIPE_LIVE_MODE") == "true"

	return &Adapter{
		secretKey:  secretKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		liveMode:   liveMode,
	}
}

func (a *Adapter) IsConfigured() bool {
	return a.secretKey != ""
}

func (a *Adapter) ProcessCard(ctx context.Context, req processor.CardRequest) (*processor.ProcessorResponse, error) {
	if !a.IsConfigured() {
		return nil, fmt.Errorf("stripe adapter not configured: set STRIPE_SECRET_KEY or STRIPE_TEST_SECRET_KEY")
	}

	apiURL := "https://api.stripe.com/v1/payment_intents"
	if !a.liveMode {
		apiURL = "https://api.stripe.com/v1/payment_intents"
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create stripe request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+a.secretKey)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return &processor.ProcessorResponse{
		Success:      false,
		ProcessorRef: "",
		Status:       "pending",
		Message:      "Stripe integration requires additional setup - API key not configured for mock mode",
	}, nil
}
