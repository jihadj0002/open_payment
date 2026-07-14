# Webhooks

## Overview
Webhooks notify merchants of asynchronous events (payment success, failure, refunds, etc.). The gateway sends HTTP POST requests to the merchant's configured endpoint with a signed JSON payload.

## Webhook Events

### payment.success
Sent when a payment is successfully captured (or automatically captured after authorization).

```json
{
  "id": "evt_abc123",
  "type": "payment.success",
  "created": 1735689600,
  "data": {
    "object": {
      "id": "pi_abc123",
      "object": "payment_intent",
      "amount": 1000,
      "currency": "BDT",
      "status": "succeeded",
      "customer_id": "cus_abc123",
      "description": "Order #1234",
      "metadata": { "order_id": "ORD-1234" },
      "payment_method": "card",
      "payment_method_details": {
        "card": { "last4": "1111", "brand": "visa" }
      }
    }
  }
}
```

### payment.failed
Sent when a payment authorization fails.

```json
{
  "id": "evt_def456",
  "type": "payment.failed",
  "created": 1735689601,
  "data": {
    "object": {
      "id": "pi_def456",
      "object": "payment_intent",
      "amount": 1000,
      "currency": "BDT",
      "status": "failed",
      "last_payment_error": {
        "code": "card_declined",
        "message": "Your card was declined. Please use a different card.",
        "decline_code": "do_not_honor"
      }
    }
  }
}
```

### payment.pending
Sent when a payment requires additional action (3DS, redirect).

### refund.completed
Sent when a refund is successfully processed.

```json
{
  "id": "evt_ghi789",
  "type": "refund.completed",
  "created": 1735689602,
  "data": {
    "object": {
      "id": "re_abc123",
      "object": "refund",
      "amount": 500,
      "currency": "BDT",
      "payment_id": "pi_abc123",
      "status": "succeeded",
      "reason": "customer_request"
    }
  }
}
```

### refund.failed
Sent when a refund cannot be processed.

### chargeback.created
Sent when a chargeback is initiated by the cardholder's bank.

```json
{
  "id": "evt_jkl012",
  "type": "chargeback.created",
  "created": 1735689603,
  "data": {
    "object": {
      "id": "cb_abc123",
      "payment_id": "pi_abc123",
      "amount": 1000,
      "currency": "BDT",
      "reason": "fraudulent",
      "status": "needs_response",
      "respond_by": 1736294400
    }
  }
}
```

### chargeback.resolved
Sent when a chargeback is resolved (won or lost).

### settlement.completed
Sent when a settlement batch completes and merchant balance is updated.

### payout.sent
Sent when a payout is initiated to the merchant's bank account.

## Webhook Signature Verification

Every webhook request includes the following headers:

| Header | Description |
|--------|-------------|
| `X-Webhook-ID` | Unique webhook delivery ID |
| `X-Webhook-Timestamp` | Unix timestamp (seconds) of payload generation |
| `X-Webhook-Signature` | HMAC-SHA256 hex digest |

### Signature Generation
```
signature = HMAC-SHA256(
    secret = webhook_secret,
    message = timestamp + "." + payload_body
)
```

### Signature Verification (Node.js Example)
```javascript
const crypto = require('crypto');

const secret = 'whsec_abc123';
const payload = JSON.stringify(req.body);
const timestamp = req.headers['x-webhook-timestamp'];
const signature = req.headers['x-webhook-signature'];

const expected = crypto
    .createHmac('sha256', secret)
    .update(timestamp + '.' + payload)
    .digest('hex');

if (signature !== expected) {
    throw new Error('Invalid webhook signature');
}

// Optional: Check timestamp is within 5 minutes
const age = Math.floor(Date.now() / 1000) - parseInt(timestamp);
if (age > 300) {
    throw new Error('Webhook too old');
}
```

### Signature Verification (Go Example)
```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "time"
)

func VerifyWebhook(secret []byte, payload []byte, timestamp string, signature string) bool {
    mac := hmac.New(sha256.New, secret)
    mac.Write([]byte(timestamp))
    mac.Write([]byte("."))
    mac.Write(payload)
    expected := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expected))
}
```

## Retry Policy

| Attempt | Delay | Total Time |
|---------|-------|------------|
| 1 | Immediate | 0s |
| 2 | 5 seconds | 5s |
| 3 | 30 seconds | 35s |
| 4 | 5 minutes | 5m 35s |
| 5 | 30 minutes | 35m 35s |
| Max | — | 35m 35s |

- If all retries fail, the webhook is marked as `failed` in the dashboard
- Merchants can manually replay failed webhooks from the dashboard
- Webhooks that fail repeatedly may be auto-disabled with a notification to the merchant

## Merchant Response

The merchant's endpoint must return a `2xx` HTTP status within 10 seconds to acknowledge receipt. Any non-2xx response (including `3xx`, `4xx`, `5xx`) triggers a retry.

## Best Practices
1. **Verify signature** before trusting the payload
2. **Respond quickly** — return 200 immediately, process asynchronously
3. **Idempotency** — use `id` field to deduplicate (same webhook may be sent multiple times)
4. **Timeout** — handle cases where the webhook endpoint is slow or down

## Webhook Delivery Logs
Merchants can view webhook delivery logs in the dashboard:
- Delivery status (success/failed/pending)
- HTTP response code
- Response body (first 1KB)
- Timestamps of each attempt
