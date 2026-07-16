# bKash - Create Payment (Sale or Authorize)

**Endpoint:** `POST https://checkout.sandbox.bka.sh/v1.2.0-beta/checkout/payment/create`

Creates a payment for sale or authorization using the bKash checkout API.

---

## Headers

| Header | Type | Required | Description |
|--------|------|----------|-------------|
| `X-APP-Key` | string | Yes | Application key |

---

## Body Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `amount` | string | Yes | Transaction amount |
| `currency` | string | Yes | Currency (e.g. `BDT`) |
| `intent` | string | Yes | Payment intent (`sale` or `authorization`) |
| `merchantInvoiceNumber` | string | Yes | Unique invoice number from merchant |
| `merchantAssociationInfo` | string | No | Merchant association information |

---

## Response

### 200 - Success

Returns the created payment details.

---

## Example Request

```shell
curl -X POST https://checkout.sandbox.bka.sh/v1.2.0-beta/checkout/payment/create \
  -H "X-APP-Key: <your-app-key>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "100",
    "currency": "BDT",
    "intent": "sale",
    "merchantInvoiceNumber": "INV-12345"
  }'
```

---

## Notes

- Use the sandbox base URL (`checkout.sandbox.bka.sh`) for testing.
- For production, replace with the production endpoint provided by bKash.
- The `intent` field determines whether the payment is captured immediately (`sale`) or only authorized (`authorization`) for later capture.
