package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	tokens    int64
	lastRefill time.Time
}

type RateLimiter struct {
	mu        sync.RWMutex
	buckets   map[string]*bucket
	rate      int64
	burst     int64
	refillDur time.Duration
}

func NewRateLimiter(rate, burst int64, refillDur time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets:   make(map[string]*bucket),
		rate:      rate,
		burst:     burst,
		refillDur: refillDur,
	}
}

func (rl *RateLimiter) allow(key string) (bool, int64, int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[key]
	if !ok {
		b = &bucket{
			tokens:    rl.burst,
			lastRefill: time.Now(),
		}
		rl.buckets[key] = b
	}

	now := time.Now()
	elapsed := now.Sub(b.lastRefill)
	refillTokens := int64(elapsed / rl.refillDur) * rl.rate
	if refillTokens > 0 {
		b.tokens = min(b.tokens+refillTokens, rl.burst)
		b.lastRefill = now
	}

	remaining := b.tokens
	if b.tokens <= 0 {
		return false, 0, remaining
	}

	b.tokens--
	return true, b.tokens, rl.burst
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr

			allowed, remaining, limit := limiter.allow(key)
			w.Header().Set("X-RateLimit-Limit", formatInt(limit))
			w.Header().Set("X-RateLimit-Remaining", formatInt(remaining))

			if !allowed {
				w.Header().Set("X-RateLimit-Reset", formatInt(time.Now().Add(time.Second).Unix()))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"type":    "rate_limit",
						"code":    "rate_limit_exceeded",
						"message": "Too many requests",
						"status":  http.StatusTooManyRequests,
					},
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func formatInt(n int64) string {
	if n < 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	if len(buf) == 0 {
		return "0"
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
