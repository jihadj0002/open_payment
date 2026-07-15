# Load Testing with k6

## Prerequisites

- Install [k6](https://k6.io/docs/getting-started/installation/)
- Server running at `http://localhost:8080` (or set `BASE_URL`)

## Running Tests

### Payment Flow Test (Create → Capture → Refund)
```bash
k6 run tests/k6/payment-flow.js
```

### Mixed Workload Test (Payments + Balance + List)
```bash
k6 run tests/k6/mixed-workload.js
```

### Custom Server / Auth
```bash
k6 run -e BASE_URL=http://localhost:8080/v1 -e AUTH_TOKEN=your_token tests/k6/payment-flow.js
```

## Test Scenarios

| Script | VUs | Duration | Description |
|--------|-----|----------|-------------|
| payment-flow.js | 50 ramp → 50 → 0 | 2 min | Login → create payment → capture → refund |
| mixed-workload.js | 100 ramp → 100 → 0 | 5 min | Mix of create, list, get, balance |

## Thresholds

- p95 response time < 500ms
- Error rate < 5% (mixed workload)
- Payment success rate > 99% (payment flow)
