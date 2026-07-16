package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type APIUsageLogger interface {
	LogAPIUsage(ctx context.Context, merchantID, endpoint, method string, statusCode int, durationMs int, ipAddress, userAgent string) error
}

func NewAPIUsageMiddleware(logger APIUsageLogger, merchantIDFromCtx func(ctx context.Context) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			durationMs := int(time.Since(start).Milliseconds())
			statusCode := ww.Status()

			merchantID := merchantIDFromCtx(r.Context())

			if merchantID != "" {
				logger.LogAPIUsage(r.Context(), merchantID, r.URL.Path, r.Method, statusCode, durationMs, r.RemoteAddr, r.UserAgent())
			}
		})
	}
}
