package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "status"},
	)

	PaymentStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_status_total",
			Help: "Total number of payments by status",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ErrorCount)
	prometheus.MustRegister(PaymentStatus)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		status := strconv.Itoa(ww.Status())
		duration := time.Since(start).Seconds()

		RequestCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

		if ww.Status() >= 400 {
			ErrorCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		}
	})
}

func RecordPaymentStatus(status string) {
	PaymentStatus.WithLabelValues(status).Inc()
}
