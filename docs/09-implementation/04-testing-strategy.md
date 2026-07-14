# Testing Strategy

## Test Pyramid

```
        ╱╲
       ╱  ╲          E2E Tests (5%)
      ╱    ╲         Playwright, k6, Litmus
     ╱      ╲
    ╱────────╲
   ╱          ╲     Integration Tests (25%)
  ╱            ╲    testcontainers (PG, Kafka, Redis)
 ╱──────────────╲
╱                  ╲ Unit Tests (70%)
╱   Go test + mock  ╲ + Frontend unit (Vitest)
╲────────────────────╲
```

## Unit Tests

### Backend (Go)
```go
// Test patterns
func TestCreatePayment_ValidRequest_ReturnsPayment(t *testing.T) {
    mockRepo := new(MockPaymentRepository)
    mockProcessor := new(MockProcessorClient)
    mockFraud := new(MockFraudService)
    mockLedger := new(MockLedgerService)
    mockProducer := new(MockEventProducer)

    svc := NewPaymentService(mockRepo, mockProcessor, mockFraud, mockLedger, mockProducer)

    req := CreatePaymentRequest{Amount: 1000, Currency: "BDT", Confirm: true}

    // Set expectations
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    mockFraud.On("AssessRisk", mock.Anything, mock.Anything).Return(RiskAssessment{Action: "allow"}, nil)
    mockProcessor.On("Authorize", mock.Anything, mock.Anything).Return(AuthorizationResult{Status: "approved"}, nil)
    mockProducer.On("Produce", mock.Anything, mock.Anything, mock.Anything).Return(nil)

    payment, err := svc.CreatePayment(context.Background(), req)

    assert.NoError(t, err)
    assert.Equal(t, PaymentStatusAuthorized, payment.Status)
    mockRepo.AssertExpectations(t)
}
```

### Coverage Targets
| Module | Coverage Target | Notes |
|--------|----------------|-------|
| Core services (payment, ledger) | >90% | Critical business logic |
| Supporting services (merchant, customer) | >80% | Standard CRUD |
| Infrastructure (DB, cache, messaging) | >70% | Hard to mock fully |
| State machine | 100% | All transitions tested |
| Frontend components | >80% | React Testing Library |
| Frontend hooks | >90% | Data fetching, mutations |

## Integration Tests

### Test Containers
```go
func TestPaymentIntegration(t *testing.T) {
    ctx := context.Background()

    // Start real dependencies
    postgres, err := testcontainers.StartPostgres(ctx)
    require.NoError(t, err)
    defer postgres.Terminate(ctx)

    kafka, err := testcontainers.StartKafka(ctx)
    require.NoError(t, err)
    defer kafka.Terminate(ctx)

    redis, err := testcontainers.StartRedis(ctx)
    require.NoError(t, err)
    defer redis.Terminate(ctx)

    // Run migrations
    db := connectDB(postgres.ConnectionString())
    runMigrations(db)

    // Initialize service with real infra
    svc := NewPaymentService(
        NewPaymentRepository(db),
        NewMockProcessor(),
        NewMockFraud(),
        NewLedgerService(NewLedgerRepository(db)),
        NewKafkaProducer(kafka.Brokers()),
    )

    // Create merchant
    merchant := createTestMerchant(db)

    // Test full payment cycle
    payment, err := svc.CreatePayment(ctx, CreatePaymentRequest{
        MerchantID: merchant.ID,
        Amount:     1000,
        Currency:   "BDT",
        Confirm:    true,
    })
    require.NoError(t, err)
    assert.Equal(t, PaymentStatusAuthorized, payment.Status)

    // Verify event was produced
    msg, err := consumeMessage(kafka.Brokers(), "payment.authorized")
    require.NoError(t, err)
    assert.Contains(t, string(msg.Value), payment.ID.String())

    // Verify ledger entries
    entries, err := getLedgerEntries(db, payment.ID)
    require.NoError(t, err)
    assert.Len(t, entries, 4)  // debit, fee, tax, net
}
```

## E2E Tests

### API E2E (k6 Load Tests)
```javascript
// k6/payment-flow.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up
    { duration: '5m', target: 100 },  // Stay
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

export default function () {
  const payload = JSON.stringify({
    amount: 1000,
    currency: 'BDT',
    payment_method: 'card',
    confirm: true,
  });

  const res = http.post('https://api.openpayment.com/v1/payments', payload, {
    headers: {
      'Authorization': `Bearer ${__ENV.API_KEY}`,
      'Content-Type': 'application/json',
      'Idempotency-Key': `${__VU}-${__ITER}`,
    },
  });

  check(res, {
    'status is 201 or 200': (r) => r.status === 201 || r.status === 200,
    'payment status is success or authorized': (r) => {
      const body = JSON.parse(r.body);
      return ['authorized', 'succeeded'].includes(body.status);
    },
  });
}
```

## Chaos Engineering (Litmus)

### Experiments
```yaml
# litmus/experiments/pod-delete.yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: payment-service-chaos
spec:
  appinfo:
    appns: payment
    applabel: app=payment-service
    appkind: deployment
  experiments:
  - name: pod-delete
    spec:
      duration: 60s
      chaosServiceAccount: litmus-admin
      probe:
      - name: "check-api-health"
        type: httpProbe
        httpProbe/inputs:
          url: "https://api.openpayment.com/health"
          expectedStatusCode: 200
```

### Chaos Test Scenarios
| Scenario | Tool | Frequency | Success Criteria |
|----------|------|-----------|------------------|
| Pod crash | Litmus | Monthly | Service recovers in <30s, no data loss |
| Network partition | Litmus | Quarterly | Service degrades gracefully, no double-charge |
| CPU spike | Litmus | Quarterly | Auto-scaling triggers, no dropped requests |
| Database failover | Manual DR | Quarterly | RTO < 5min, RPO < 1min |
| Kafka outage | Litmus | Quarterly | Events queued, processed after recovery |
| Redis failure | Litmus | Monthly | Fallback to DB, no incorrect state |

## Security Testing

### CI/CD Security Scans
```yaml
# Already covered in CI/CD pipeline:
- Dependency scan:    Trivy, `go list -m all`
- Container scan:     Trivy Docker image scan
- SAST:              golangci-lint security checks
- Secret detection:  GitLeaks in CI
```

### Manual Security Testing
| Test Type | Frequency | Tester |
|-----------|-----------|--------|
| Penetration test | Quarterly | Third-party |
| Vulnerability assessment | Monthly | Internal security |
| Code security review | Per feature | Security champion |
| PCI compliance scan | Quarterly | ASV |
| Social engineering | Annually | Third-party |

## Performance Benchmarks

### CI Performance Regression
```yaml
# .github/benchmarks/payment-benchmark.yml
name: Benchmark Payments
on:
  pull_request:
    paths: ['internal/service/payment/**']

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run benchmarks
        run: go test -bench=. -benchmem ./internal/service/payment/...
      - name: Compare with main
        uses: benchmark-action/github-action-benchmark@v1
```

### Target Benchmarks
| Operation | Current Target | Regression Threshold |
|-----------|---------------|---------------------|
| Create payment (card) | <300ms | +20% |
| Capture payment | <200ms | +20% |
| Refund payment | <200ms | +20% |
| List payments (page) | <100ms | +20% |
| Fraud check | <100ms | +20% |
| Ledger entry creation | <50ms | +20% |
| Webhook delivery | <2s p95 | +25% |
