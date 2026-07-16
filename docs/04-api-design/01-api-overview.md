# API Overview

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/api/router.go`, `internal/service/*/routes.go`

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
> **TODO: Not implemented** — This is a planned feature for v1.1.

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
| `X-Signature` | TODO v1.1 | HMAC signature (not implemented) |
| `X-Timestamp` | TODO v1.1 | Unix timestamp in ms (not implemented) |
| `X-Nonce` | TODO v1.1 | Unique request identifier (not implemented) |
| `Accept-Language` | No | `en`, `bn` for localization |

### Response Headers

| Header | Description |
|--------|-------------|
| `X-Request-Id` | Unique request identifier (via `chimw.RequestID`) |
| `X-RateLimit-Limit` | Rate limit ceiling for the current endpoint |
| `X-RateLimit-Remaining` | Number of requests remaining in the current window |
| `X-API-Version` | API version (`1`) |

## Pagination

> **Note:** Currently uses **offset-based pagination** (`page`/`per_page`). Cursor-based pagination is planned for v1.1.

### Offset-based pagination (current)
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

### Cursor-based pagination (planned for v1.1)
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
| Auth | `/v1/auth` | None | Login, register, refresh, forgot/reset password |
| Payments | `/v1/payments` | Bearer JWT | Create, list, get, capture, refund, void |
| Customers | `/v1/customers` | Bearer JWT | Customer profiles, saved payment methods |
| Balance / Ledger | `/v1/balance` | Bearer JWT | View balance, transactions |
| Webhook Endpoints | `/v1/webhook_endpoints` | Bearer JWT | Configure webhook endpoints, logs, replay |
| Settlements | `/v1/settlements` | Bearer JWT | Trigger/List/Get settlements, reports |
| Merchant | `/v1/merchants` | Bearer JWT | Profile, API keys, onboarding, fraud rules |
| Admin | `/v1/admin` | Bearer JWT + Admin role | Merchant management, system config, disputes |
| Fraud Config | `/v1/merchants/fraud_config` | Bearer JWT | Fraud detection settings |
| Subscriptions | `/v1/subscriptions` | Secret key | **TODO: Not implemented** |
| Tokens | `/v1/tokens` | Publishable key | **TODO: Not implemented** |
