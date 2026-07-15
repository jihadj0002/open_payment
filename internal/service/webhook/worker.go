package webhook

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type RetryWorker struct {
	svc      *Service
	interval time.Duration
}

func NewRetryWorker(svc *Service, interval time.Duration) *RetryWorker {
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &RetryWorker{
		svc:      svc,
		interval: interval,
	}
}

func (w *RetryWorker) Start(ctx context.Context) {
	log.Info().Dur("interval", w.interval).Msg("webhook retry worker started")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debug().Msg("webhook retry worker: checking pending deliveries")
			w.svc.RetryPendingDeliveries(ctx)
		case <-ctx.Done():
			log.Info().Msg("webhook retry worker stopped")
			return
		}
	}
}
