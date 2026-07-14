# Monitoring and Observability

## Stack Components

| Component | Tool | Purpose |
|-----------|------|---------|
| Metrics | Prometheus + Grafana | Time-series data collection, dashboards |
| Logging | Loki + Promtail | Centralized log aggregation |
| Tracing | OpenTelemetry + Tempo | Distributed tracing across services |
| Error Tracking | Sentry | Real-time error monitoring |
| Uptime | Cloudflare + StatusPage | External monitoring, status page |
| Alerting | PagerDuty + AlertManager | On-call notifications |
| APM | Grafana + Tempo | Performance monitoring |

## Metrics Collection

### Prometheus Metrics per Service

Each service exposes `/metrics` endpoint with:

```go
// Payment Service metrics
var (
    PaymentCreated = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "payment_created_total",
        Help: "Total number of payment intents created",
    }, []string{"merchant_id", "currency", "payment_method"})

    PaymentStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "payment_status_current",
        Help: "Current payments by status",
    }, []string{"status"})

    PaymentLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "payment_processing_duration_seconds",
        Help:    "Payment processing latency in seconds",
        Buckets: prometheus.DefBuckets,
    }, []string{"operation", "status"})

    WebhookDeliveryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "webhook_delivery_duration_seconds",
        Help:    "Webhook delivery latency",
        Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
    }, []string{"status"})
)
```

### Key Metrics by Category

| Category | Metric | Alert Threshold |
|----------|--------|-----------------|
| API | http_requests_total (by status, path) | 5xx > 0.1% |
| API | http_request_duration_seconds (p50, p95, p99) | p99 > 1s |
| Payments | payment_created_total | — |
| Payments | payment_processing_duration_seconds | p95 > 5s |
| Payments | payment_failure_rate | > 5% |
| Ledger | ledger_entries_total | — |
| Ledger | ledger_reconciliation_diff | != 0 (critical) |
| Webhook | webhook_delivery_duration_seconds | p99 > 10s |
| Webhook | webhook_failure_total | > 1% |
| System | cpu_usage, memory_usage | > 80% |
| System | go_goroutines | — |
| Queue | kafka_consumer_lag | > 1000 |
| Queue | kafka_messages_total | — |
| DB | pg_connections, pg_replication_lag | > 100 conns, > 10MB lag |
| Redis | redis_hit_ratio | < 85% |

### Prometheus Recording Rules
```yaml
groups:
  - name: payment_recording.rules
    rules:
      - record: payment:success_rate:5m
        expr: |
          sum(rate(payment_created_total{status="succeeded"}[5m]))
          /
          sum(rate(payment_created_total[5m]))
          * 100

      - record: api:p99_latency:5m
        expr: |
          histogram_quantile(0.99,
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path))
```

## Grafana Dashboards

### 1. Payment Operations Dashboard
- **Row 1:** Payment success rate (gauge), Volume (stat), Active merchants
- **Row 2:** Payment latency (heatmap), Success rate over time (graph)
- **Row 3:** Payments by status (pie chart), Payments by method (bar chart)
- **Row 4:** Recent failed payments (table), Webhook delivery status

### 2. System Health Dashboard
- **Row 1:** Service status (stat table, green/red)
- **Row 2:** CPU/Memory per service (graph), Request rate per service
- **Row 3:** Error rate by service, Error rate by HTTP status
- **Row 4:** Database connections, Replication lag, Redis hit ratio

### 3. Business KPI Dashboard
- **Row 1:** Revenue today, Revenue this month, Pending settlements
- **Row 2:** Revenue trend (7d, 30d), Refund rate, Chargeback rate
- **Row 3:** Top merchants by volume, Payment method distribution
- **Row 4:** Average transaction value, New merchants this month

### 4. Security Dashboard
- **Row 1:** Failed login attempts, Rate limit hits, API auth failures
- **Row 2:** Fraud alerts history, Blocked transactions, Flagged IPs
- **Row 3:** Audit log volume, Admin actions timeline

## Alert Rules

### Critical Alerts (PagerDuty)
```yaml
groups:
  - name: critical
    rules:
      - alert: PaymentServiceDown
        expr: up{job="payment-service"} == 0
        for: 1m
        labels: { severity: critical }
        annotations:
          summary: "Payment service is down"

      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.001
        for: 2m
        labels: { severity: critical }
        annotations:
          summary: "API error rate > 0.1%"

      - alert: HighPaymentFailure
        expr: payment_failure_rate > 0.05
        for: 5m
        labels: { severity: critical }
        annotations:
          summary: "Payment failure rate > 5%"

      - alert: LedgerImbalance
        expr: ledger_reconciliation_diff != 0
        labels: { severity: critical }
        annotations:
          summary: "Ledger imbalance detected"

      - alert: DatabaseReplicationLag
        expr: pg_replication_lag > 10485760  # 10MB
        for: 5m
        labels: { severity: critical }
```

### Warning Alerts (Slack)
```yaml
      - alert: HighLatency
        expr: api:p99_latency:5m > 1
        for: 5m
        labels: { severity: warning }

      - alert: ConsumerLag
        expr: kafka_consumer_lag > 1000
        for: 5m
        labels: { severity: warning }

      - alert: WebhookFailure
        expr: rate(webhook_failure_total[5m]) > 0.01
        for: 5m
        labels: { severity: warning }
```

## Logging Strategy

### Log Format (Structured JSON)
```json
{
  "timestamp": "2026-01-15T10:30:00Z",
  "level": "info",
  "service": "payment-service",
  "request_id": "req_abc123",
  "payment_id": "pi_abc123",
  "merchant_id": "m_abc123",
  "message": "Payment authorized successfully",
  "duration_ms": 245,
  "processor_response": "00",
  "trace_id": "abcdef123456"
}
```

### Log Levels
| Level | Usage |
|-------|-------|
| `debug` | Detailed debug info (disabled in prod) |
| `info` | Normal operations (payment created, captured) |
| `warn` | Unexpected but handled (retry, slow request) |
| `error` | Error requiring investigation (payment failed, DB error) |
| `fatal` | Service can't continue (will trigger crash) |

### What NOT to Log
- Card numbers (PAN), CVV, full track data
- API secret keys
- Passwords, password hashes
- JWT tokens
- Personal data beyond what's necessary
- Internal IPs in production logs

## Distributed Tracing

### Trace Sampling
| Strategy | Description | Rate |
|----------|-------------|------|
| Always sample | Payment processing flows | 100% |
| Head-based | General API requests | 10% |
| Error-based | Failed requests | 100% |
| Slow requests | >1s duration | 100% |

### Tracing Context Propagation
```go
// Each service extracts/injects trace context
func PaymentHandler(w http.ResponseWriter, r *http.Request) {
    // Extract trace context from incoming request
    ctx := otel.GetTextMapPropagator().Extract(r.Context(),
        propagation.HeaderCarrier(r.Header))

    // Create span
    ctx, span := tracer.Start(ctx, "payment.create")
    defer span.End()

    // Pass context to downstream calls
    resp, err := fraudService.AssessRisk(ctx, req)
}
```

### Key Spans
| Span | Parent | What It Measures |
|------|--------|-----------------|
| HTTP request | (root) | Full API request duration |
| Auth check | HTTP request | JWT/API key validation |
| Payment create | HTTP request | Payment intent creation |
| Fraud check | Payment create | Risk assessment latency |
| Processor auth | Payment create | Authorization with card network |
| Ledger entry | Payment create/capture | Ledger write latency |
| Webhook delivery | (separate) | Webhook HTTP call |
