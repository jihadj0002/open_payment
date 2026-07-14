# API Overview

## Base URLs

| Environment | Base URL |
|-------------|----------|
| Production | `https://api.openpayment.com/v1` |
| Sandbox | `https://sandbox.api.openpayment.com/v1` |
| Development | `https://dev.api.openpayment.com/v1` |

## API Versioning
- Version prefix: `/v1/` in URL path
- Breaking changes trigger `/v2/`
- Non-breaking changes (new fields, new endpoints) don't require version bump
- Deprecated versions supported for 12 months minimum

## Authentication Methods

### API Key Authentication (Merchant API)
```
Header: Authorization: Bearer sk_live_xxxxxxxxxxxxx
```
- Secret keys start with `sk_live_` (production) or `sk_test_` (sandbox)
- Publishable keys start with `pk_live_` / `pk_test_`
- Secret key used for all server-side API calls

### JWT Authentication (Dashboard)
```
Header: Authorization: Bearer <jwt_token>
```
- Obtained via OAuth2 login flow
- Access token expires in 15 minutes
- Refresh token expires in 30 days
- JWT contains: `{ merchant_id, user_id, role, permissions }`

### HMAC Request Signing (Optional, recommended)
```
Header: X-Signature: <hex-encoded HMAC-SHA256>
Header: X-Timestamp: <unix_timestamp_ms>
Header: X-Nonce: <random_uuid>
```
- Signing key: `HMAC-SHA256(secret_key, method + path + body + timestamp + nonce)`
- Timestamp must be within 5 minutes of server time
- Each nonce can only be used once (prevents replay)

## Common Headers

| Header | Required | Description |
|--------|----------|-------------|
| `Authorization` | Yes | Bearer token or API key |
| `Content-Type` | Yes | `application/json` |
| `Idempotency-Key` | For write ops | UUID v4, unique per request |
| `X-Signature` | Optional | HMAC signature |
| `X-Timestamp` | With signature | Unix timestamp in ms |
| `X-Nonce` | With signature | Unique request identifier |
| `Accept-Language` | No | `en`, `bn` for localization |

## Pagination

### Cursor-based pagination (default)
```
GET /v1/transactions?cursor=txn_abc123&limit=50
```
Response:
```json
{
  "data": [...],
  "has_more": true,
  "next_cursor": "txn_xyz789"
}
```

### Offset-based pagination (for dashboards)
```
GET /v1/transactions?page=1&per_page=25
```
Response:
```json
{
  "data": [...],
  "total": 1042,
  "page": 1,
  "per_page": 25,
  "total_pages": 42
}
```

## Response Format

### Success
```json
{
  "id": "pi_abc123",
  "object": "payment_intent",
  "amount": 1000,
  "currency": "BDT",
  "status": "succeeded",
  "created": 1735689600,
  ...
}
```

### Error
```json
{
  "error": {
    "type": "invalid_request_error",
    "code": "amount_too_small",
    "message": "Amount must be at least 1 BDT",
    "param": "amount",
    "status": 400
  },
  "request_id": "req_abc123"
}
```

## Rate Limiting

| Tier | Rate Limit | Headers Returned |
|------|------------|------------------|
| Free | 100 req/min | `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset` |
| Business | 1,000 req/min | Same |
| Enterprise | 10,000 req/min | Same |

When exceeded:
```
Status: 429 Too Many Requests
Retry-After: 30
```

## API Groups

| Group | Base Path | Auth | Description |
|-------|-----------|------|-------------|
| Payments | `/v1/payments` | Secret key | Create, capture, refund, void |
| Customers | `/v1/customers` | Secret key | Customer profiles, saved methods |
| Subscriptions | `/v1/subscriptions` | Secret key | Recurring billing |
| Balance | `/v1/balance` | Secret key | View balance, transactions |
| Webhooks | `/v1/webhooks` | Secret key | Configure webhook endpoints |
| Tokens | `/v1/tokens` | Publishable key | Card tokenization |
| Admin | `/v1/admin` | JWT (admin role) | Merchant management, system config |
| Auth | `/v1/auth` | None | Login, MFA, token refresh |
