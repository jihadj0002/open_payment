package webhook

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/service/auth"
)

type EventListItem struct {
	ID                string     `json:"id"`
	Event             string     `json:"event"`
	Status            string     `json:"status"`
	AttemptCount      int        `json:"attempt_count"`
	LastAttemptedAt   *time.Time `json:"last_attempted_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type HealthStatus struct {
	WebhookID           string   `json:"webhook_id"`
	SuccessRate         float64  `json:"success_rate"`
	TotalAttempts       int      `json:"total_attempts"`
	LastSuccessAt       *time.Time `json:"last_success_at,omitempty"`
	LastFailureAt       *time.Time `json:"last_failure_at,omitempty"`
	ConsecutiveFailures int      `json:"consecutive_failures"`
	IsActive            bool     `json:"is_active"`
}

func HandleListEvents(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "webhook endpoint ID is required")
			return
		}

		deliveries, err := svc.repo.ListDeliveriesByWebhook(r.Context(), id, claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		events := make([]EventListItem, 0, len(deliveries))
		seen := make(map[string]bool)
		for _, d := range deliveries {
			if seen[d.ID] {
				continue
			}
			seen[d.ID] = true
			events = append(events, EventListItem{
				ID:              d.ID,
				Event:           d.Event,
				Status:          d.Status,
				AttemptCount:    d.Attempt,
				LastAttemptedAt: d.NextAttemptAt,
				CreatedAt:       d.CreatedAt,
			})
		}

		data := make([]interface{}, len(events))
		for i, v := range events {
			data[i] = v
		}
		api.RespondJSON(w, http.StatusOK, data)
	}
}

func HandleReplayEvent(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		endpointID := chi.URLParam(r, "id")
		eventID := chi.URLParam(r, "eventId")
		if endpointID == "" || eventID == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "endpoint ID and event ID are required")
			return
		}

		delivery, err := svc.repo.GetDeliveryByID(r.Context(), eventID)
		if err != nil {
			api.RespondError(w, http.StatusNotFound, "not_found", "event_not_found", "event not found")
			return
		}

		endpoint, err := svc.repo.GetEndpoint(r.Context(), endpointID, claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusNotFound, "not_found", "endpoint_not_found", "webhook endpoint not found")
			return
		}

		go func() {
			ctx := context.Background()
			svc.sendWebhook(ctx, endpoint, delivery.Event, nil)
		}()

		api.RespondJSON(w, http.StatusOK, map[string]string{
			"status":  "replayed",
			"event":   delivery.Event,
			"message": "event re-dispatched",
		})
	}
}

func HandleWebhookHealth(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "webhook endpoint ID is required")
			return
		}

		endpoint, err := svc.repo.GetEndpoint(r.Context(), id, claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusNotFound, "not_found", "endpoint_not_found", "webhook endpoint not found")
			return
		}

		deliveries, err := svc.repo.ListDeliveriesByWebhook(r.Context(), id, claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		total := len(deliveries)
		successCount := 0
		consecutiveFailures := 0
		var lastSuccessAt *time.Time
		var lastFailureAt *time.Time

		for _, d := range deliveries {
			if d.Status == "delivered" {
				successCount++
				t := d.CreatedAt
				if lastSuccessAt == nil || t.After(*lastSuccessAt) {
					lastSuccessAt = &t
				}
			} else {
				t := d.CreatedAt
				if lastFailureAt == nil || t.After(*lastFailureAt) {
					lastFailureAt = &t
				}
			}
		}

		for i := len(deliveries) - 1; i >= 0; i-- {
			if deliveries[i].Status != "delivered" {
				consecutiveFailures++
			} else {
				break
			}
		}

		successRate := 0.0
		if total > 0 {
			successRate = float64(successCount) / float64(total) * 100
		}

		isActive := endpoint.Status == "active"
		if consecutiveFailures >= 10 {
			isActive = false
		}

		api.RespondJSON(w, http.StatusOK, HealthStatus{
			WebhookID:           id,
			SuccessRate:         successRate,
			TotalAttempts:       total,
			LastSuccessAt:       lastSuccessAt,
			LastFailureAt:       lastFailureAt,
			ConsecutiveFailures: consecutiveFailures,
			IsActive:            isActive,
		})
	}
}

func HandleAutoDisableUnhealthy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "monitoring_active"})
	}
}
