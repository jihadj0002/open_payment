# Merchant REST API Specification

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

Cancels a payment that hasn't been captured yet.

## Refunds

### POST /v1/refunds — Create Refund

**Request:**
```json
{
  "payment_id": "pi_abc123",
  "amount": 500,
  "reason": "customer_request",
  "metadata": {
    "return_reference": "RET-001"
  }
}
```

**Response (201 Created):**
```json
{
  "id": "re_abc123",
  "object": "refund",
  "amount": 500,
  "currency": "BDT",
  "payment_id": "pi_abc123",
  "status": "succeeded",
  "reason": "customer_request",
  "metadata": { "return_reference": "RET-001" },
  "created": 1735689600
}
```

### GET /v1/refunds/:id — Retrieve Refund

### GET /v1/refunds — List Refunds

## Customers

### POST /v1/customers — Create Customer

**Request:**
```json
{
  "email": "customer@example.com",
  "name": "John Doe",
  "phone": "+8801700000000",
  "metadata": {
    "internal_id": "USR-001"
  }
}
```

**Response (201 Created):**
```json
{
  "id": "cus_abc123",
  "object": "customer",
  "email": "customer@example.com",
  "name": "John Doe",
  "phone": "+8801700000000",
  "metadata": { "internal_id": "USR-001" },
  "created": 1735689600
}
```

### POST /v1/customers/:id/payment_methods — Attach Payment Method

**Request:**
```json
{
  "payment_method": "card",
  "payment_method_data": {
    "card": {
      "number": "4111111111111111",
      "exp_month": 12,
      "exp_year": 2027,
      "cvc": "123"
    }
  },
  "set_as_default": true
}
```

### GET /v1/customers/:id/payment_methods — List Saved Methods

## Balance

### GET /v1/balance — Retrieve Balance

**Response:**
```json
{
  "object": "balance",
  "available": [
    { "amount": 150000, "currency": "BDT" }
  ],
  "pending": [
    { "amount": 25000, "currency": "BDT" }
  ],
  "reserve": [
    { "amount": 5000, "currency": "BDT" }
  ]
}
```

### GET /v1/balance/transactions — List Balance Transactions

## Webhook Endpoints

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

### GET /v1/webhook_endpoints — List Webhook Endpoints

### DELETE /v1/webhook_endpoints/:id — Delete Webhook Endpoint

## Tokens

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
