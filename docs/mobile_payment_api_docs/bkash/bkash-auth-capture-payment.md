# bKash - Authorization & Capture (Auth & Capture)

Pre-authorize an amount and capture it later upon service confirmation. Compatible with the Checkout (URL Based) flow.

## Create Payment (Authorization)

Creates a payment with `intent: "authorization"` to block the amount without charging immediately.

**Request URL:** `{base_URL}/checkout/payment/create`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Content-Type | application/json |
| Accept | application/json |
| Authorization | id_token from Grant/Refresh Token API |
| X-App-Key | App key from onboarding |

### Request Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| amount | Yes | string | Amount to authorize |
| currency | Yes | string | Currency (BDT) |
| intent | Yes | string | Must be `"authorization"` |
| merchantInvoiceNumber | Yes | string | Unique invoice number (max 255 chars) |
| merchantAssociationInfo | No | string | TLV format for aggregators |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code |
| paymentID | string | Payment ID (expires after 24 hours) |
| bkashURL | string | URL to redirect customer |
| callbackURL | string | Base callback URL |
| successCallbackURL | string | Success callback |
| failureCallbackURL | string | Failure callback |
| cancelledCallbackURL | string | Cancelled callback |
| amount | string | Authorized amount |
| currency | string | Currency (BDT) |
| intent | string | `authorization` |
| transactionStatus | string | Status (`Initiated`) |
| paymentCreateTime | string | Creation timestamp |
| merchantInvoiceNumber | string | Invoice number |

### Sample Request

```json
POST /checkout/payment/create HTTP/1.1
Host: {base_URL}
Content-Type: application/json
Accept: application/json
authorization: id_token
x-app-key: x-app-key

{
  "mode": "0011",
  "payerReference": "01723888888",
  "callbackURL": "yourDomain.com",
  "merchantAssociationInfo": "MI05MID54RF09123456One",
  "amount": "500",
  "currency": "BDT",
  "intent": "authorization",
  "merchantInvoiceNumber": "Inv0124"
}
```

### Sample Response

```json
{
  "statusCode": "0000",
  "statusMessage": "Successful",
  "paymentID": "TR0011ON1565154754797",
  "callbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout",
  "successCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=success&signature=cm8HBfl65A",
  "failureCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=failure&signature=cm8HBfl65A",
  "cancelledCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=cancel&signature=cm8HBfl65A",
  "amount": "500",
  "intent": "authorization",
  "currency": "BDT",
  "paymentCreateTime": "2019-08-07T11:12:34:978 GMT+0600",
  "transactionStatus": "Initiated",
  "merchantInvoiceNumber": "Inv0124"
}
```

## Execute Payment (Capture)

Captures the pre-authorized amount after service confirmation.

**Request URL:** `{base_URL}/checkout/payment/execute/{paymentID}`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Content-Type | application/json |
| Accept | application/json |
| Authorization | id_token from Grant/Refresh Token API |
| X-App-Key | App key from onboarding |

### Request Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| paymentID | Yes | string (path) | PaymentID from Create Payment response |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code |
| paymentID | string | Payment ID |
| customerMsisdn | string | Customer MSISDN |
| payerReference | string | Payer reference |
| trxID | string | Transaction ID |
| transactionStatus | string | Final status (`Completed`) |
| amount | string | Captured amount |
| currency | string | Currency (BDT) |
| intent | string | `authorization` |
| merchantInvoiceNumber | string | Invoice number |

### Sample Request

```json
POST /checkout/payment/execute/TR000139pVZqh169096716**** HTTP/1.1
Host: {base_URL}
Accept: application/json
authorization: id_token
x-app-key: x-app-key
```

### Sample Response

```json
{
  "statusCode": "0000",
  "statusMessage": "Successful",
  "paymentID": "TR000139pVZqh169096716****",
  "agreementID": "TokenizedMerchant02************",
  "payerReference": "01619777***",
  "customerMsisdn": "01619777***",
  "trxID": "AH260BY2HW",
  "amount": "500",
  "transactionStatus": "Completed",
  "paymentExecuteTime": "2023-08-02T15:08:15:744 GMT+0600",
  "currency": "BDT",
  "intent": "authorization",
  "merchantInvoiceNumber": "Inv0124"
}
```
