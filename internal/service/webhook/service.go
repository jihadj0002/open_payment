package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	maxWebhookAttempts = 5
	httpTimeout        = 10 * time.Second
)

var retryDelays = []time.Duration{
	0,
	5 * time.Second,
	30 * time.Second,
	5 * time.Minute,
	30 * time.Minute,
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEndpoint(ctx context.Context, merchantID string, req CreateEndpointRequest) (*Endpoint, error) {
	if req.URL == "" {
		return nil, errors.New("url is required")
	}
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		return nil, errors.New("url must start with http:// or https://")
	}
	if !validEvents[req.Event] {
		return nil, fmt.Errorf("invalid event: %s", req.Event)
	}

	e := &Endpoint{
		MerchantID: merchantID,
		Event:      req.Event,
		URL:        req.URL,
		Secret:     generateSecret(),
		Status:     "active",
	}

	if err := s.repo.CreateEndpoint(ctx, e); err != nil {
		return nil, fmt.Errorf("create endpoint: %w", err)
	}

	return e, nil
}

func (s *Service) ListEndpoints(ctx context.Context, merchantID string) ([]Endpoint, error) {
	return s.repo.ListEndpoints(ctx, merchantID)
}

func (s *Service) DeleteEndpoint(ctx context.Context, id, merchantID string) error {
	_, err := s.repo.GetEndpoint(ctx, id, merchantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("endpoint not found")
		}
		return err
	}
	return s.repo.DeleteEndpoint(ctx, id, merchantID)
}

func (s *Service) RotateSecret(ctx context.Context, id, merchantID string) (*RotateSecretResponse, error) {
	endpoint, err := s.repo.GetEndpoint(ctx, id, merchantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("endpoint not found")
		}
		return nil, err
	}

	newSecret := generateSecret()
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	if err := s.repo.RotateSecret(ctx, id, endpoint.Secret, newSecret, expiresAt); err != nil {
		return nil, fmt.Errorf("rotate secret: %w", err)
	}

	return &RotateSecretResponse{
		EndpointID: id,
		NewSecret:  newSecret,
		Message:    "new secret generated; previous secret remains valid for 24 hours",
	}, nil
}

func (s *Service) DispatchEvent(ctx context.Context, merchantID, eventType string, data interface{}) {
	endpoints, err := s.repo.GetEndpointsForEvent(ctx, merchantID, eventType)
	if err != nil || len(endpoints) == 0 {
		return
	}

	for _, ep := range endpoints {
		ep := ep
		s.sendWebhook(ctx, &ep, eventType, data)
	}
}

func (s *Service) sendWebhook(ctx context.Context, endpoint *Endpoint, eventType string, data interface{}) {
	eventID := uuid.New().String()
	now := time.Now()

	payload := WebhookEvent{
		ID:      eventID,
		Type:    eventType,
		Created: now.Unix(),
		Data:    data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	timestamp := fmt.Sprintf("%d", now.Unix())
	sig := computeSignature(endpoint.Secret, timestamp, body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.URL, bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-ID", eventID)
	req.Header.Set("X-Webhook-Timestamp", timestamp)
	req.Header.Set("X-Webhook-Signature", sig)

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)

	delivery := &Delivery{
		WebhookID:   endpoint.ID,
		Event:       eventType,
		Payload:     body,
		Attempt:     1,
		MaxAttempts: maxWebhookAttempts,
	}

	if err != nil {
		delivery.Status = "failed"
		s.createDelivery(ctx, delivery)
		return
	}
	defer resp.Body.Close()

	respBodyBytes, _ := io.ReadAll(resp.Body)
	code := resp.StatusCode
	delivery.ResponseCode = &code

	_ = respBodyBytes

	if code >= 200 && code < 300 {
		delivery.Status = "delivered"
		delivery.NextAttemptAt = nil
	} else {
		delivery.Status = "pending"
		delay := retryDelay(1)
		if delay > 0 {
			next := now.Add(delay)
			delivery.NextAttemptAt = &next
		}
	}

	s.createDelivery(ctx, delivery)
}

func (s *Service) createDelivery(ctx context.Context, d *Delivery) {
	if err := s.repo.CreateDelivery(ctx, d); err != nil {
		return
	}
}

func (s *Service) RetryPendingDeliveries(ctx context.Context) {
	deliveries, err := s.repo.GetPendingDeliveries(ctx, 50)
	if err != nil || len(deliveries) == 0 {
		return
	}

	for _, d := range deliveries {
		d := d
		s.retryDelivery(ctx, &d)
	}
}

func (s *Service) retryDelivery(ctx context.Context, d *Delivery) {
	endpoint, err := s.repo.GetDeliveryWebhook(ctx, d.WebhookID)
	if err != nil || endpoint == nil || endpoint.Status != "active" {
		d.Status = "failed"
		s.repo.UpdateDelivery(ctx, d)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.URL, bytes.NewReader(d.Payload))
	if err != nil {
		return
	}

	var payloadData WebhookEvent
	if err := json.Unmarshal(d.Payload, &payloadData); err != nil {
		return
	}

	timestamp := fmt.Sprintf("%d", payloadData.Created)
	sig := computeSignature(endpoint.Secret, timestamp, d.Payload)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-ID", payloadData.ID)
	req.Header.Set("X-Webhook-Timestamp", timestamp)
	req.Header.Set("X-Webhook-Signature", sig)

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)

	now := time.Now()

	d.Attempt++

	if err != nil {
		if d.Attempt >= d.MaxAttempts {
			d.Status = "failed"
			d.NextAttemptAt = nil
		} else {
			d.Status = "pending"
			delay := retryDelay(d.Attempt)
			next := now.Add(delay)
			d.NextAttemptAt = &next
		}
		s.repo.UpdateDelivery(ctx, d)
		return
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	d.ResponseCode = &code

	if code >= 200 && code < 300 {
		d.Status = "delivered"
		d.NextAttemptAt = nil
	} else {
		if d.Attempt >= d.MaxAttempts {
			d.Status = "failed"
			d.NextAttemptAt = nil
		} else {
			d.Status = "pending"
			delay := retryDelay(d.Attempt)
			next := now.Add(delay)
			d.NextAttemptAt = &next
		}
	}

	s.repo.UpdateDelivery(ctx, d)
}

func retryDelay(attempt int) time.Duration {
	idx := attempt - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(retryDelays) {
		return retryDelays[len(retryDelays)-1]
	}
	return retryDelays[idx]
}

func computeSignature(secret, timestamp string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func generateSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
