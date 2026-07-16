# bKash Integration Guide

Accept bKash mobile wallet payments through your Open Payment Gateway integration.

## Overview

bKash is the leading mobile financial service in Bangladesh. The integration uses the **bKash Tokenized Checkout API** (URL-based flow) where customers are redirected to a bKash-hosted payment page, enter their wallet credentials, and complete payment.

### Flow

1. Merchant creates a Payment Intent via API (`POST /v1/payments`)
2. Customer is redirected to the hosted checkout page (`/checkout/{payment_intent_id}`)
3. Customer selects bKash and clicks Pay
4. Customer is redirected to bKash's secure payment page
5. Customer enters phone number + OTP + PIN on the bKash page
6. bKash calls back our server with the payment status
7. Customer is redirected back to the success/cancel page
8. Merchant receives a webhook (`payment.success` or `payment.failed`)

## Prerequisites

1. Register as a bKash Merchant at the [bKash Merchant Portal](https://www.bkash.com/)
2. Obtain your merchant credentials during onboarding:
   - `BKASH_APP_KEY` — Application key
   - `BKASH_APP_SECRET` — Application secret
   - `BKASH_USERNAME` — API username
   - `BKASH_PASSWORD` — API password
   - `BKASH_BASE_URL` — API base URL (sandbox or production)

## Creating a Payment Intent

```http
POST /v1/payments
Authorization: Bearer <merchant_token>
Content-Type: application/json

{
  "amount": 5000,
  "currency": "BDT",
  "payment_method": "wallet",
  "description": "Order #12345",
  "return_url": "https://merchant.com/order/12345/success",
  "cancel_url": "https://merchant.com/order/12345/cancel",
  "metadata": {
    "customer_phone": "01712345678",
    "provider": "bkash"
  }
}
```

**Response:**

```json
{
  "id": "pi_uuid",
  "amount": 5000,
  "currency": "BDT",
  "status": "created",
  "payment_method": "wallet",
  "client_secret": "pi_uuid_secret_uuid",
  "return_url": "https://merchant.com/order/12345/success",
  "metadata": {
    "customer_phone": "01712345678",
    "provider": "bkash"
  }
}
```

## Checkout Flow

### 1. Get Checkout Session

```http
GET /v1/checkout/{payment_intent_id}
```

Returns the checkout session data for rendering the payment page.

### 2. Initiate Payment

```http
POST /v1/checkout/{payment_intent_id}/pay
Content-Type: application/json

{
  "payment_method": "wallet",
  "provider": "bkash",
  "customer_phone": "01712345678",
  "return_url": "https://merchant.com/success",
  "cancel_url": "https://merchant.com/cancel"
}
```

**Response (redirect to bKash):**

```json
{
  "id": "pi_uuid",
  "status": "wallet_initiated",
  "redirect_url": "https://sandbox.payment.bkash.com/?paymentId=TR0011..."
}
```

### 3. Customer Payment

The customer is redirected to the `redirect_url`. They enter their bKash wallet number, receive an OTP, and confirm with their PIN.

### 4. Callback Handling

After the customer completes payment, bKash redirects back to:
- `/checkout/{payment_intent_id}/success` — on success
- `/checkout/{payment_intent_id}/cancel` — on cancel
- `/checkout/{payment_intent_id}/error` — on failure

If you provided `return_url`/`cancel_url`, the customer will be redirected there instead.

## Webhooks

Merchants receive the following webhook events:

| Event | Trigger |
|-------|---------|
| `payment.success` | Payment completed successfully |
| `payment.failed` | Payment failed or declined |
| `payment.canceled` | Customer cancelled the payment |

### Webhook Payload Example

```json
{
  "id": "evt_uuid",
  "type": "payment.success",
  "created": 1700000000,
  "data": {
    "id": "pi_uuid",
    "amount": 5000,
    "currency": "BDT",
    "status": "succeeded",
    "merchant_id": "merchant_uuid"
  }
}
```

## Testing (Sandbox)

Use the bKash sandbox environment:

| Detail | Value |
|--------|-------|
| Base URL | `https://checkout.sandbox.bka.sh/v1.2.0-beta` |
| Test Wallet | Provided by bKash sandbox |
| OTP | `123456` (sandbox) |
| PIN | `1234` (sandbox) |

## Error Codes

See [bKash Error Codes](/docs/mobile_payment_api_docs/bkash/bkash-error-codes.md) for the full list.

Common errors:
- `2006` — Invalid amount
- `2010` — Invalid OTP
- `2023` — Insufficient balance
- `2056` — Invalid payment state

## Refunds

Refunds are processed through the existing `POST /v1/payments/{id}/refund` endpoint. The gateway calls the bKash Refund Transaction API internally.
