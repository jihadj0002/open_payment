# Error Handling

## Error Taxonomy

| Error Type | HTTP Status | Description |
|------------|-------------|-------------|
| `invalid_request_error` | 400 | Malformed request, missing fields, validation failures |
| `authentication_error` | 401 | Invalid or missing API key / JWT |
| `permission_error` | 403 | API key or user lacks required permissions |
| `not_found_error` | 404 | Requested resource doesn't exist |
| `conflict_error` | 409 | Idempotency conflict, resource state conflict |
| `rate_limit_error` | 429 | Too many requests |
| `api_error` | 500 | Unexpected server error |

## Error Response Format

```json
{
  "error": {
    "type": "invalid_request_error",
    "code": "missing_required_field",
    "message": "The field 'amount' is required",
    "param": "amount",
    "status": 400,
    "details": {
      "field": "amount",
      "reason": "required"
    }
  },
  "request_id": "req_abc123"
}
```

## Error Codes Reference

### Validation Errors (400)

| Code | Message | Param |
|------|---------|-------|
| `missing_required_field` | "The field '{field}' is required" | Field name |
| `invalid_amount` | "Amount must be a positive integer" | amount |
| `amount_too_small` | "Amount must be at least 1 {currency}" | amount |
| `amount_too_large` | "Amount cannot exceed {max}" | amount |
| `invalid_currency` | "Unsupported currency: {currency}" | currency |
| `invalid_card_number` | "Invalid card number" | payment_method_data.card.number |
| `invalid_expiry_month` | "Invalid expiration month" | payment_method_data.card.exp_month |
| `invalid_expiry_year` | "Invalid expiration year" | payment_method_data.card.exp_year |
| `invalid_cvc` | "Invalid CVC code" | payment_method_data.card.cvc |
| `card_expired` | "Card is expired" | payment_method_data.card |
| `unsupported_card_brand` | "Card brand {brand} is not supported" | payment_method_data.card |
| `invalid_customer_id` | "Customer not found" | customer_id |
| `invalid_payment_method` | "Unsupported payment method" | payment_method |
| `invalid_metadata` | "Metadata key '{key}' exceeds 40 characters" | metadata |
| `invalid_metadata_value` | "Metadata value exceeds 500 characters" | metadata |
| `invalid_return_url` | "Invalid URL format" | return_url |
| `invalid_idempotency_key` | "Idempotency key must be a UUID v4" | Idempotency-Key header |

### Authentication Errors (401)

| Code | Message | Notes |
|------|---------|-------|
| `invalid_api_key` | "Invalid API key" | Key doesn't exist or is revoked |
| `expired_api_key` | "API key has expired" | Key rotated or expired |
| `invalid_jwt` | "Invalid or expired JWT" | Token expired or malformed |
| `invalid_signature` | "Request signature verification failed" | HMAC mismatch |
| `expired_timestamp` | "Request timestamp is too old" | Timestamp > 5 minutes from now |

### Permission Errors (403)

| Code | Message |
|------|---------|
| `insufficient_permissions` | "API key does not have permission for this action" |
| `merchant_suspended` | "Merchant account is suspended" |
| `merchant_terminated` | "Merchant account has been terminated" |
| `ip_not_whitelisted` | "IP address not in allowlist" |

### Not Found Errors (404)

| Code | Message |
|------|---------|
| `payment_not_found` | "Payment intent not found" |
| `refund_not_found` | "Refund not found" |
| `customer_not_found` | "Customer not found" |
| `webhook_not_found` | "Webhook endpoint not found" |
| `subscription_not_found` | "Subscription not found" |

### Conflict Errors (409)

| Code | Message | Resolution |
|------|---------|------------|
| `idempotency_key_used` | "Idempotency key already used for a different request" | Use unique key |
| `payment_already_captured` | "Payment has already been captured" | Refund instead |
| `payment_not_capturable` | "Payment is not in a capturable state" | Wait for authorization |
| `payment_not_refundable` | "Payment is not in a refundable state" | Must be captured first |
| `refund_amount_exceeded` | "Refund amount exceeds remaining capturable amount" | Lower refund amount |
| `duplicate_webhook_url` | "A webhook endpoint with this URL already exists" | Use different URL |

### Rate Limit Errors (429)

| Code | Message |
|------|---------|
| `rate_limit_exceeded` | "Too many requests. Please retry after {seconds} seconds" |
| `webhook_rate_limit_exceeded` | "Too many webhook registrations. Please wait and try again" |

### Server Errors (500)

| Code | Message |
|------|---------|
| `internal_error` | "An unexpected error occurred" |
| `processor_error` | "Payment processor returned an error" |
| `processor_timeout` | "Payment processor did not respond in time" |
| `ledger_error` | "Ledger service error — transaction may not be recorded" |
| `fraud_service_error` | "Fraud detection service unavailable" |
| `webhook_delivery_error` | "Webhook delivery failed after max retries" |

## Idempotency

### How It Works
1. Client generates a UUID v4 for each write request
2. Passes it as `Idempotency-Key` header
3. Server checks if key was already used:
   - **Not used:** Process request, cache response for 24 hours
   - **Same key, same request:** Return cached response
   - **Same key, different request:** Return 409 conflict
4. Cache is stored in Redis with 24h TTL

### Idempotency Scope
| Operation | Idempotency Scope | Notes |
|-----------|-------------------|-------|
| Create Payment | Per key | Same key = same payment returned |
| Capture | Per payment + key | Only 1 capture per payment + key |
| Refund | Per payment + key | Only 1 refund per payment + key |
| Create Customer | Per key | Same key = same customer returned |

## Best Practices for Clients

### Always use idempotency keys on write operations
```go
import "github.com/google/uuid"

func (c *Client) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*PaymentIntent, error) {
    idempotencyKey := uuid.New().String()
    return c.post(ctx, "/v1/payments", req, withIdempotencyKey(idempotencyKey))
}
```

### Handle idempotency conflict (409)
```python
try:
    payment = gateway.Payment.create(amount=1000)
except gateway.IdempotencyError as e:
    # Key was already used for a DIFFERENT request
    # Generate a new key and retry
    payment = gateway.Payment.create(
        amount=1000,
        idempotency_key=str(uuid.uuid4())
    )
```

### Exponential backoff for rate limits
```javascript
async function retryWithBackoff(fn, maxRetries = 3) {
    for (let i = 0; i < maxRetries; i++) {
        try {
            return await fn();
        } catch (err) {
            if (err.status === 429) {
                const wait = Math.pow(2, i) * 1000;
                await new Promise(r => setTimeout(r, wait));
                continue;
            }
            throw err;
        }
    }
    throw new Error('Max retries exceeded');
}
```
