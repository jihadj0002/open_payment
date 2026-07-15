package requestid

import (
	"context"

	chimw "github.com/go-chi/chi/v5/middleware"
)

type contextKey string

const reqIDKey contextKey = "request_id"

func FromContext(ctx context.Context) string {
	if id := chimw.GetReqID(ctx); id != "" {
		return id
	}
	return ""
}

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, reqIDKey, id)
}
