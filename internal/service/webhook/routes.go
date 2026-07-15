package webhook

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterWebhookRoutes(r chi.Router, svc *Service, authMW func(http.Handler) http.Handler) {
	r.With(authMW).Post("/webhook_endpoints", HandleCreateEndpoint(svc))
	r.With(authMW).Get("/webhook_endpoints", HandleListEndpoints(svc))
	r.With(authMW).Delete("/webhook_endpoints/{id}", HandleDeleteEndpoint(svc))
	r.With(authMW).Post("/webhook_endpoints/{id}/rotate-secret", HandleRotateSecret(svc))
	r.With(authMW).Get("/webhook_endpoints/{id}/events", HandleListEvents(svc))
	r.With(authMW).Post("/webhook_endpoints/{id}/events/{eventId}/replay", HandleReplayEvent(svc))
	r.With(authMW).Get("/webhook_endpoints/{id}/health", HandleWebhookHealth(svc))
}
