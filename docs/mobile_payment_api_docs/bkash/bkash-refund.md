# bKash - Refund Transaction & Refund Status

Process full or partial refunds on completed transactions. Supports multiple partial refunds (up to 10 times) until the full amount is refunded.

## Refund Transaction

Process a payment reversal for a transaction.

**Request URL:** `{base_URL}/v2/tokenized-checkout/refund/payment/transaction`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Accept | application/json |
| Content-Type | application/json |
| Authorization | id_token from Grant/Refresh Token API |
| X-App-Key | App key from onboarding |

### Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| paymentId | string | Yes | Payment ID from Create Payment response |
| trxId | string | Yes | Transaction ID from Execute Payment response |
| refundAmount | string | Yes | Amount to refund (max 2 decimal places, e.g. 25.20) |
| sku | string | No | Product/service info (max 255 chars) |
| reason | string | No | Refund reason (max 255 chars) |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| originalTrxId | string | Original transaction ID |
| refundTrxId | string | Refund transaction ID |
| refundTransactionStatus | string | Refund status (`Completed`) |
| originalTrxAmount | string | Original transaction amount |
| refundAmount | string | Amount refunded |
| currency | string | Currency (BDT) |
| completedTime | string | Refund completion timestamp |
| sku | string | SKU info |
| reason | string | Refund reason |

### Sample Request

```json
POST /v2/tokenized-checkout/refund/payment/transaction HTTP/1.1
Host: {base_URL}
Content-Type: application/json
Accept: application/json
authorization: id_token
x-app-key: x-app-key

{
  "paymentId": "TR0001xt7mXxG1718274354990",
  "trxId": "BFD90JRLST",
  "refundAmount": "1",
  "sku": "test",
  "reason": "test"
}
```

### Sample Response

```json
{
  "originalTrxId": "BFD90JRLST",
  "refundTrxId": "BFD90JRMH9",
  "refundTransactionStatus": "Completed",
  "originalTrxAmount": "4.59",
  "refundAmount": "1.00",
  "currency": "BDT",
  "completedTime": "2024-06-13T16:27:25:422 GMT+0600",
  "sku": "test",
  "reason": "test"
}
```

## Refund Status

Check the status of a refund transaction.

**Request URL:** `{base_URL}/v2/tokenized-checkout/refund/payment/status`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Accept | application/json |
| Authorization | id_token from Grant/Refresh Token API |
| X-App-Key | App key from onboarding |

### Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| paymentId | string | Yes | Payment ID from Create Payment |
| trxId | string | Yes | Transaction ID from Execute Payment |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| originalTrxId | string | Original transaction ID |
| originalTrxAmount | string | Original transaction amount |
| originalTrxCompletedTime | string | Original transaction completed time |
| refundTransactions[] | array | List of refund records |
| refundTrxId | string | Refund transaction ID (per item) |
| refundTransactionStatus | string | Refund status (per item) |
| refundAmount | string | Refunded amount (per item) |
| completedTime | string | Refund completed time (per item) |

### Sample Request

```json
POST /v2/tokenized-checkout/refund/payment/status HTTP/1.1
Host: {base_URL}
Content-Type: application/json
Accept: application/json
authorization: id_token
x-app-key: x-app-key

{
  "paymentId": "TR0001xt7mXxG1718274354990",
  "trxId": "BFD90JRLST"
}
```

### Sample Response

```json
{
  "originalTrxId": "BFD90JRLST",
  "originalTrxAmount": "4.59",
  "originalTrxCompletedTime": "2024-06-13T16:26:41:486 GMT+0600",
  "refundTransactions": [
    {
      "refundTrxId": "BFD90JRMH9",
      "refundTransactionStatus": "Completed",
      "refundAmount": "1.00",
      "completedTime": "2024-06-13T16:27:24:000"
    },
    {
      "refundTrxId": "BFD90JRMK7",
      "refundTransactionStatus": "Completed",
      "refundAmount": "2.00",
      "completedTime": "2024-06-17T18:27:24:000"
    }
  ]
}
```

## Refund Error Codes

| Error Message | Code |
|---------------|------|
| The merchant is not permitted to initiate this transaction | 2082 |
| The identity of the debit or credit party is in a state which prohibits the execution of this transaction | 2081 |
| The identity is not permitted to initiate this transaction | 2080 |
| Invalid app Token | 2079 |
| Invalid Reason | 2078 |
| Invalid TrxID | 2077 |
| Reason Character Limit Exceeded | 2076 |
| SKU Character Limit Exceeded | 2075 |
| The transaction cannot be reversed | 2074 |
| Invalid SKU | 2073 |
| Refund amount not valid | 2072 |
| Refund after %s days not allowed | 2071 |
| Transaction not yet completed | 2127 |
| Insufficient Balance | 2023 |

## Notes

- Multiple partial refunds are supported (up to 10 times)
- If no response within 30 seconds, call Refund Status API to get the actual status
- Default timeout for bKash APIs: 30 seconds
- Only `Completed` status means a successful refund
