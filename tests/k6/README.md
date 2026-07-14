# k6 Load Testing Scripts

## Prerequisites

- [k6](https://k6.io/docs/getting-started/installation/) installed (`brew install k6`, `apt install k6`, etc.)
- The payment gateway API running locally or on a staging environment

## Scripts

| Script | Type | Description |
|--------|------|-------------|
| `smoke-test.js` | Smoke | Minimal load (1 VU, 30s) — validates basic functionality |
| `load-test.js` | Load | Gradual ramp to 100 VUs over 12min — tests normal traffic |
| `stress-test.js` | Stress | Ramp to 300 VUs over 12min — finds breaking point |
| `soak-test.js` | Soak | 100 VUs sustained for 60min — detects memory/resource leaks |

## How to Run

```bash
# Smoke test (default: localhost:8080)
k6 run tests/k6/smoke-test.js

# Load test
k6 run tests/k6/load-test.js

# Stress test
k6 run tests/k6/stress-test.js

# Soak test
k6 run tests/k6/soak-test.js
```

## Configuration

Override the target URL or API key via environment variables:

```bash
# Custom target URL
BASE_URL=http://staging:8080 k6 run tests/k6/load-test.js

# Custom API key
API_KEY=sk_live_xxxxxxxxxxxxxx k6 run tests/k6/smoke-test.js
```

## Thresholds

| Script | p(95) | p(99) | Failure Rate |
|--------|-------|-------|--------------|
| Smoke  | <500ms | —    | <1%          |
| Load   | <500ms | <1500ms | <1%        |
| Stress | <1000ms | <3000ms | <5%        |
| Soak   | <600ms | —    | <2%          |

If a threshold is crossed, k6 exits with a non-zero code and marks the test as failed.
