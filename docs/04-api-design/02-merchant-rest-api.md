# Merchant REST API Specification

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/service/merchant/routes.go`, `internal/service/payment/routes.go`, `internal/service/auth/routes.go`, `internal/service/customer/routes.go`, `internal/service/ledger/routes.go`, `internal/service/webhook/routes.go`, `internal/service/settlement/routes.go`, `internal/service/fraud/routes.go`, `internal/service/admin/routes.go`

## Payments

### POST /v1/payments — Create Payment Intent

Creates a payment intent for authorization and capture.

**Request:**
```json
{
  "amount": 1000,
  "currency": "BDT",
  "payment_method": "card",
  "payment_method_data": {
    "card": {
      "number": "4111111111111111",
      "exp_month": 12,
      "exp_year": 2027,
      "cvc": "123"
    }
  },
  "customer_id": "cus_abc123",
  "description": "Order #1234",
  "metadata": {
    "order_id": "ORD-1234"
  },
  "confirm": true,
  "return_url": "https://example.com/checkout/success",
  "capture_method": "automatic"
}
```

**Response (201 Created):**
```json
{
  "id": "pi_abc123",
  "object": "payment_intent",
  "amount": 1000,
  "amount_capturable": 1000,
  "amount_received": 0,
  "currency": "BDT",
  "status": "requires_capture",
  "customer_id": "cus_abc123",
  "description": "Order #1234",
  "metadata": {
    "order_id": "ORD-1234"
  },
  "capture_method": "automatic",
  "payment_method": "card",
  "payment_method_details": {
    "card": {
      "last4": "1111",
      "brand": "visa",
      "exp_month": 12,
      "exp_year": 2027
    }
  },
  "created": 1735689600,
  "client_secret": "pi_abc123_secret_xyz789"
}
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| amount | integer | Yes | Amount in smallest currency unit (paise/satangs) |
| currency | string | Yes | ISO 4217 currency code (BDT, USD, etc.) |
| payment_method | string | Yes | `card`, `wallet`, `bank_transfer` |
| payment_method_data | object | Conditional | Required if no saved payment method |
| customer_id | string | No | Reference to saved customer |
| description | string | No | Max 255 characters |
| metadata | object | No | Key-value pairs (max 20 keys) |
| confirm | boolean | No | Whether to confirm immediately (default false) |
| return_url | string | Conditional | Required if redirect needed (3DS, wallet redirect) |
| capture_method | enum | No | `automatic` (default), `manual` |
| idempotency_key | header | Recommended | Prevents duplicate payment intents |

### GET /v1/payments/:id — Retrieve Payment Intent

**Response:**
```json
{
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
  },
  "charges": {
    "data": [{
      "id": "ch_abc123",
      "amount": 1000,
      "status": "succeeded",
      "created": 1735689600
    }]
  },
  "created": 1735689600
}
```

### GET /v1/payments — List Payments

Query parameters:
- `customer`: Filter by customer ID
- `status`: `requires_payment_method`, `requires_confirmation`, `requires_capture`, `succeeded`, `failed`, `canceled`
- `created[gt]`, `created[gte]`, `created[lt]`, `created[lte]`: Date range filters
- `cursor`, `limit`: Pagination

### POST /v1/payments/:id/capture — Capture Payment

Captures an authorized payment. Required if `capture_method: manual`.

**Request:**
```json
{
  "amount_to_capture": 1000
}
```

**Response (200 OK):** Returns the updated Payment Intent with `status: succeeded`.

### POST /v1/payments/:id/void — Void Payment

Voids an uncaptured authorization.

**Response (200 OK):** Payment Intent with `status: canceled`.

### POST /v1/payments/:id/cancel — Cancel Payment

> **Note:** Use `POST /v1/payments/{id}/void` instead. Cancel is not implemented separately.

## Refunds

> **Note:** Refunds are not a separate resource. Use `POST /v1/payments/{id}/refund` to refund a captured payment.

### POST /v1/payments/{id}/refund — Refund Payment

**Request:**
```json
{
  "amount": 500,
  "reason": "customer_request",
  "metadata": {
    "return_reference": "RET-001"
  }
}
```

**Response (200 OK):** Returns the updated Payment Intent.

## Auth Endpoints

> **Code Ref:** `internal/service/auth/routes.go`

### POST /v1/auth/login — Login

**Request:** `{"email": "...", "password": "..."}`

**Response:** JWT access token + merchant profile.

### POST /v1/auth/register — Register Merchant

**Request:** `{"name": "...", "email": "...", "password": "..."}`

**Response:** JWT access token + merchant profile + API keys (`secret_key`, `public_key`).

### POST /v1/auth/refresh — Refresh Token

**Request:** `{"refresh_token": "..."}`

**Response:** New access + refresh tokens.

### POST /v1/auth/forgot-password — Forgot Password

**Request:** `{"email": "..."}`

**Response:** 200 OK (email sent if account exists).

### POST /v1/auth/reset-password — Reset Password

**Request:** `{"token": "...", "password": "..."}`

**Response:** 200 OK (password updated).

### GET /v1/auth/me — Current User

**Response:** Merchant profile data.

## Customers

> **Code Ref:** `internal/service/customer/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/customers` | Create customer |
| GET | `/v1/customers` | List customers |
| GET | `/v1/customers/{id}` | Get customer |
| PATCH | `/v1/customers/{id}` | Update customer |
| POST | `/v1/customers/{id}/payment_methods` | Attach payment method |
| GET | `/v1/customers/{id}/payment_methods` | List saved payment methods |
| DELETE | `/v1/customers/{id}/payment_methods/{pm_id}` | Detach payment method |

## Balance / Ledger

> **Code Ref:** `internal/service/ledger/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/balance` | Retrieve merchant balance |
| GET | `/v1/balance/transactions` | List balance transactions |

## Webhook Endpoints

> **Code Ref:** `internal/service/webhook/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/webhook_endpoints` | Create webhook endpoint |
| GET | `/v1/webhook_endpoints` | List webhook endpoints |
| DELETE | `/v1/webhook_endpoints/{id}` | Delete webhook endpoint |
| POST | `/v1/webhook_endpoints/{id}/rotate-secret` | Rotate webhook signing secret |
| GET | `/v1/webhook_endpoints/{id}/events` | List event delivery history |
| POST | `/v1/webhook_endpoints/{id}/events/{eventId}/replay` | Replay a webhook event |
| GET | `/v1/webhook_endpoints/{id}/health` | Webhook endpoint health status |

### POST /v1/webhook_endpoints — Create Webhook Endpoint

**Request:**
```json
{
  "url": "https://example.com/webhook",
  "enabled_events": [
    "payment.success",
    "payment.failed",
    "refund.completed"
  ],
  "description": "Production webhook"
}
```

## Tokens

> **TODO: Not implemented** — Planned for v1.1.

### POST /v1/tokens — Create Token (Card)

Use publishable key for this endpoint.

**Request:**
```json
{
  "card": {
    "number": "4111111111111111",
    "exp_month": 12,
    "exp_year": 2027,
    "cvc": "123"
  }
}
```

**Response:**
```json
{
  "id": "tok_abc123",
  "object": "token",
  "card": {
    "last4": "1111",
    "brand": "visa",
    "exp_month": 12,
    "exp_year": 2027
  },
  "created": 1735689600
}
```

## Subscriptions (v1.1+)

> **TODO: Not implemented** — Planned for v1.1.

### POST /v1/subscriptions — Create Subscription

**Request:**
```json
{
  "customer_id": "cus_abc123",
  "items": [{
    "price_data": {
      "currency": "BDT",
      "product": "prod_basic_plan",
      "unit_amount": 50000,
      "interval": "month"
    },
    "quantity": 1
  }],
  "trial_period_days": 14
}
```

## Merchant Profile & API Keys

> **Code Ref:** `internal/service/merchant/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/merchants/profile` | Get merchant profile |
| PATCH | `/v1/merchants/profile` | Update merchant profile |
| GET | `/v1/merchants/api_keys` | List API keys |
| POST | `/v1/merchants/api_keys` | Create API key |
| DELETE | `/v1/merchants/api_keys/{id}` | Revoke API key |

## Settlements

> **Code Ref:** `internal/service/settlement/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/settlements` | Trigger settlement |
| GET | `/v1/settlements` | List settlements |
| GET | `/v1/settlements/{id}` | Get settlement details |
| GET | `/v1/settlements/report` | Settlement report with filters |

## Fraud Detection

> **Code Ref:** `internal/service/fraud/routes.go`, `internal/service/admin/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/merchants/fraud_config` | Get fraud detection config |
| PUT | `/v1/merchants/fraud_config` | Update fraud detection config |
| GET | `/v1/merchants/fraud/rules` | List fraud rules |
| POST | `/v1/merchants/fraud/rules` | Create fraud rule |
| PATCH | `/v1/merchants/fraud/rules/{id}` | Update fraud rule |
| DELETE | `/v1/merchants/fraud/rules/{id}` | Delete fraud rule |

## Merchant Onboarding

> **Code Ref:** `internal/service/admin/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/merchants/onboarding/status` | Get onboarding status |
| POST | `/v1/merchants/onboarding` | Submit onboarding data |
| POST | `/v1/merchants/onboarding/documents` | Upload onboarding document |

## Merchant Stats

> **Code Ref:** `internal/service/admin/routes.go`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/merchants/{id}/stats` | Get merchant usage stats |
