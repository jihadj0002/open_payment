package wallet

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MockPayment struct {
	ID            string
	Amount        int64
	Currency      string
	MerchantRef   string
	CustomerPhone string
	Status        string
	RedirectURL   string
	ProviderRef   string
	CreatedAt     time.Time
}

type MockAdapter struct {
	publicURL  string
	mu         sync.RWMutex
	payments   map[string]*MockPayment
	httpClient *http.Client
}

func NewMockAdapter(publicURL string) *MockAdapter {
	return &MockAdapter{
		publicURL:  publicURL,
		payments:   make(map[string]*MockPayment),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (m *MockAdapter) InitiatePayment(ctx context.Context, amount int64, currency, merchantRef, customerPhone, returnURL, cancelURL string) (redirectURL string, paymentRef string, providerRef string, err error) {
	paymentID := "mock_" + uuid.New().String()
	providerRef = "MOCK_REF_" + uuid.New().String()[:8]

	mp := &MockPayment{
		ID:            paymentID,
		Amount:        amount,
		Currency:      currency,
		MerchantRef:   merchantRef,
		CustomerPhone: customerPhone,
		Status:        "initiated",
		ProviderRef:   providerRef,
		CreatedAt:     time.Now(),
	}

	redirectURL = fmt.Sprintf("%s/mock-checkout/%s", m.publicURL, paymentID)

	m.mu.Lock()
	m.payments[paymentID] = mp
	m.mu.Unlock()

	return redirectURL, paymentID, providerRef, nil
}

func (m *MockAdapter) ExecutePayment(ctx context.Context, providerRef string) (success bool, processorRef string, status string, message string, err error) {
	m.mu.RLock()
	var foundID string
	for id, p := range m.payments {
		if p.ProviderRef == providerRef {
			foundID = id
			break
		}
	}
	m.mu.RUnlock()

	if foundID == "" {
		return false, "", "failed", "mock payment not found", fmt.Errorf("mock payment not found for ref: %s", providerRef)
	}

	m.mu.Lock()
	if p, ok := m.payments[foundID]; ok {
		p.Status = "completed"
	}
	m.mu.Unlock()

	return true, "MOCK_TXN_" + uuid.New().String()[:8], "succeeded", "mock payment completed", nil
}

func (m *MockAdapter) GetPaymentStatus(paymentID string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if p, ok := m.payments[paymentID]; ok {
		return p.Status
	}
	return "not_found"
}
