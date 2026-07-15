# Processor Downtime Runbook

## Symptoms
- Payment creation fails with `processor error`
- Error logs: `connection refused` to mock-processor or Stripe
- Monitoring alerts for high payment failure rate

## Steps

### 1. Identify the Downstream Processor

Check logs to identify which processor is down:
```bash
docker-compose logs api --tail=50 | grep processor
```

### 2. Processor Failover

**If using Stripe:**
- Check [status.stripe.com](https://status.stripe.com)
- If degraded, switch to backup processor (e.g., SSLCommerz for BDT)
- Update `PROCESSOR_BACKEND` env var

**If using Mock Processor (dev):**
```bash
docker-compose restart mock-processor
```

### 3. Queue Management

If processor will be down for > 5 minutes:
```bash
# Enable maintenance mode to queue payments
curl -X PATCH localhost:8080/v1/admin/config \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"maintenance_mode": true}'
```

Payments will fail fast with 503 — clients should retry with idempotency keys.

### 4. Merchant Communication

```bash
# Send mass notification via webhook events
curl -X POST localhost:8080/v1/admin/notifications \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"type": "processor_outage", "message": "Card processing is temporarily unavailable"}'
```

### 5. Recovery Verification

```bash
# Test processor health
curl http://localhost:8080/health
# Check processor status field

# Run a test payment
curl -X POST localhost:8080/v1/payments \
  -H "Authorization: Bearer $TEST_TOKEN" \
  -d '{"amount": 100, "currency": "USD", "payment_method": "card", "confirm": true}'
```

### 6. Process Queued Payments

After recovery, merchants should retry failed payments with the same idempotency key.
