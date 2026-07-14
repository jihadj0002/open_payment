# Performance Targets and Benchmarks

## API Performance Targets

| Endpoint | p50 | p95 | p99 | Throughput |
|----------|-----|-----|-----|------------|
| GET /v1/payments/:id | <30ms | <80ms | <200ms | 1000 req/s |
| GET /v1/payments | <50ms | <150ms | <300ms | 500 req/s |
| POST /v1/payments (card, auto-capture) | <300ms | <800ms | <1500ms | 200 req/s |
| POST /v1/payments/:id/capture | <100ms | <300ms | <500ms | 200 req/s |
| POST /v1/payments/:id/refund | <100ms | <300ms | <500ms | 200 req/s |
| GET /v1/balance | <30ms | <80ms | <150ms | 300 req/s |
| GET /v1/webhook_logs | <50ms | <150ms | <300ms | 200 req/s |
| POST /v1/webhook_endpoints | <50ms | <100ms | <200ms | 50 req/s |
| POST /v1/auth/login | <200ms | <500ms | <1000ms | 100 req/s |

## Infrastructure Performance Targets

| Component | Metric | Target | Criticality |
|-----------|--------|--------|-------------|
| PostgreSQL | Query time (simple lookup) | <5ms | Critical |
| PostgreSQL | Query time (list with filter) | <50ms | High |
| PostgreSQL | Write throughput | 2000 TPS | Critical |
| Redis | GET latency | <1ms | Critical |
| Redis | SET latency | <2ms | Critical |
| Redis | Hit ratio | >90% | High |
| Kafka | Produce latency (acks=all) | <10ms | Critical |
| Kafka | Consumer lag | <1000 | High |
| Kafka | Throughput | 5000 msg/s | High |

## Performance Budgets

### Frontend
| Metric | Budget | Measurement |
|--------|--------|-------------|
| First Contentful Paint (FCP) | <1.5s | Lighthouse |
| Largest Contentful Paint (LCP) | <2.5s | Lighthouse |
| Time to Interactive (TTI) | <3.5s | Lighthouse |
| JavaScript bundle size | <300KB | webpack-bundle-analyzer |
| API response to render | <500ms | React DevTools |
| Dashboard page load (p95) | <2s | RUM |

### Backend
| Metric | Budget | Measurement |
|--------|--------|-------------|
| API response time (p95) | <500ms | Prometheus |
| Payment init to auth (p95) | <5s | Distributed tracing |
| Webhook delivery (p99) | <5s | Prometheus |
| Batch settlement (100K txns) | <10min | Job metrics |
| Report generation (1M txns) | <30s | Job metrics |

## Benchmark Test Scenarios

### Scenario 1: Normal Load
- **Target:** 100 req/s sustained
- **Duration:** 15 minutes
- **Acceptance criteria:**
  - p95 latency < 200ms for reads
  - p95 latency < 500ms for writes
  - Error rate < 0.1%
  - Zero payment processing failures

### Scenario 2: Peak Load
- **Target:** 1000 req/s sustained
- **Duration:** 30 minutes
- **Acceptance criteria:**
  - p95 latency < 500ms for reads
  - p95 latency < 1000ms for writes
  - Error rate < 0.5%
  - Payment success rate > 95%

### Scenario 3: Stress Test
- **Target:** 5000 req/s for 5 minutes
- **Duration:** 5 minutes
- **Acceptance criteria:**
  - No data loss
  - Graceful degradation (429 responses)
  - Full recovery within 2 minutes after load drops
  - No crash loops or restart cycles

### Scenario 4: Endurance Test
- **Target:** 200 req/s
- **Duration:** 8 hours
- **Acceptance criteria:**
  - No memory leak (memory stable over 8h)
  - No connection pool exhaustion
  - No gradual latency increase
  - All batches/message processing completes

## Benchmark Tools

### k6 Script for Payment API
```javascript
import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const paymentLatency = new Trend('payment_latency');
const paymentErrors = new Rate('payment_errors');

export const options = {
  thresholds: {
    payment_latency: ['p(95)<500'],
    payment_errors: ['rate<0.01'],
    http_req_duration: ['p(95)<200'],
  },
};

export default function () {
  group('create payment', () => {
    const res = http.post(
      'https://api.openpayment.com/v1/payments',
      JSON.stringify({
        amount: 1000,
        currency: 'BDT',
        payment_method: 'card',
        confirm: true,
      }),
      {
        headers: {
          'Authorization': `Bearer ${__ENV.API_KEY}`,
          'Content-Type': 'application/json',
          'Idempotency-Key': `${__VU}-${Date.now()}`,
        },
      }
    );

    paymentLatency.add(res.timings.duration);
    paymentErrors.add(res.status !== 201 && res.status !== 200);

    check(res, {
      'payment created': (r) => r.status === 201,
      'has payment id': (r) => JSON.parse(r.body).id !== undefined,
    });
  });
}
```

### Go Benchmark
```go
func BenchmarkCreatePayment(b *testing.B) {
    svc := setupBenchmarkService()

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, err := svc.CreatePayment(context.Background(), CreatePaymentRequest{
                Amount:   1000,
                Currency: "BDT",
                PaymentMethod: "card",
            })
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

## Performance Optimization Techniques

| Technique | Expected Improvement | Applied To | Effort |
|-----------|---------------------|------------|--------|
| Connection pooling | 5-10x | Database, Redis, HTTP | Low |
| Query optimization (indexes) | 10-100x | Database queries | Medium |
| Caching (Redis) | 50-100x | Read-heavy endpoints | Medium |
| Pagination | Linear vs O(n) | List endpoints | Low |
| Batch processing | 10x | Settlement, reports | Medium |
| Async processing (Kafka) | Non-blocking | Webhooks, notifications | Medium |
| Connection reuse (keep-alive) | 2-3x | External API calls | Low |
| Compression (gzip) | 3-5x | API responses | Low |

## Performance Monitoring

### What to Monitor
- **Real-time:** Grafana dashboards for all key metrics
- **Alerts:** PagerDuty for performance degradation
- **APM:** Distributed traces for slow requests
- **RUM:** Real User Monitoring for frontend
- **Synthetic:** k6 running every 5 minutes from multiple regions
- **Regression:** Benchmark comparison in CI per PR
