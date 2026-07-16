# Nagad Integration Guide

Accept Nagad mobile wallet payments through your Open Payment Gateway integration.

## Overview

Nagad is a state-owned digital wallet service in Bangladesh. The integration uses the **Nagad Merchant API** where customers are redirected to a Nagad-hosted payment page, authenticate with their wallet, and complete payment.

### Flow

1. Merchant creates a Payment Intent via API (`POST /v1/payments`)
2. Customer is redirected to the hosted checkout page (`/checkout/{payment_intent_id}`)
3. Customer selects Nagad and clicks Pay
4. Customer is redirected to Nagad's secure payment page
5. Customer enters their Nagad wallet credentials on the Nagad page
6. Nagad calls back our server with the payment status
7. Customer is redirected back to the success/cancel page
8. Merchant receives a webhook (`payment.success` or `payment.failed`)

## Prerequisites

1. Register as a Nagad Merchant at the [Nagad Merchant Portal](https://www.nagad.com.bd/)
2. Obtain your merchant credentials:
   - `NAGAD_MERCHANT_ID` — Your unique merchant ID
   - `NAGAD_MERCHANT_PRIVATE_KEY` — Your RSA private key for signing requests
   - `NAGAD_PG_PUBLIC_KEY` — Nagad Payment Gateway's RSA public key for verification
   - `NAGAD_BASE_URL` — API base URL (sandbox or production)
3. Whitelist your server IP address with Nagad
4. Configure your callback URL with Nagad support

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
    "provider": "nagad"
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
    "provider": "nagad"
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
  "provider": "nagad",
  "customer_phone": "01712345678",
  "return_url": "https://merchant.com/success",
  "cancel_url": "https://merchant.com/cancel"
}
```

### 3. Customer Payment

The customer is redirected to the Nagad payment page. They authenticate with their Nagad wallet and confirm the payment.

### 4. Callback Handling

After the customer completes payment, Nagad sends a POST callback to our server. The customer is then redirected back to:
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

## Security

- All requests to Nagad are signed using RSA-SHA256 with your merchant private key
- Nagad callbacks are verified using the Nagad PG public key
- Always use HTTPS in production
- Never expose your merchant private key in client-side code

## Nagad Response Fields

| Field | Description |
|-------|-------------|
| `merchantId` | Your merchant ID |
| `orderId` | Our order/payment reference |
| `paymentRefId` | Nagad payment reference ID |
| `amount` | Transaction amount |
| `clientMobileNo` | Customer's mobile number |
| `issuerPaymentRefNo` | Issuer payment reference number |
| `status` | Transaction status (`Success`, `Cancel`, `Pending`) |

## Testing (Sandbox)

| Detail | Value |
|--------|-------|
| Mode | Set `isSandbox: true` in credentials |
| Test Wallet | Provided by Nagad sandbox |
| Amounts | Use minimal amounts for testing |

## Refunds

Refunds are processed through the existing `POST /v1/payments/{id}/refund` endpoint.
