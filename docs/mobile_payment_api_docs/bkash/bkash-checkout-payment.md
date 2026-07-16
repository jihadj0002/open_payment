# bKash - Checkout (URL Based): Create & Execute Payment

Standard checkout flow where customers enter wallet and PIN on a bKash-hosted payment page without requiring a pre-existing agreement.

## Create Payment (Checkout)

Creates a payment request for URL-based checkout.

**Request URL:** `{base_URL}/tokenized/checkout/create`

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
| mode | Yes | string | Must be `"0011"` for URL-based checkout |
| payerReference | Yes | string | Reference value (max 255 chars). Wallet number pre-populates bKash page |
| callbackURL | Yes | string | Base URL for callback URLs |
| amount | Yes | string | Payment amount |
| currency | Yes | string | Currency (BDT only) |
| intent | Yes | string | Must be `"sale"` |
| merchantInvoiceNumber | Yes | string | Unique invoice number (max 255 chars) |
| merchantAssociationInfo | No | string | TLV format for aggregators/system integrators |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code |
| paymentID | string | Payment ID (expires after 24 hours) |
| bkashURL | string | URL to redirect customer for wallet + PIN entry |
| callbackURL | string | Base callback URL |
| successCallbackURL | string | Success callback |
| failureCallbackURL | string | Failure callback |
| cancelledCallbackURL | string | Cancelled callback |
| amount | string | Transaction amount |
| currency | string | Currency (BDT) |
| intent | string | Payment intent |
| transactionStatus | string | Status (`Initiated`) |
| paymentCreateTime | string | Creation timestamp |
| merchantInvoiceNumber | string | Invoice number |

### Sample Request

```json
POST /tokenized/checkout/create HTTP/1.1
Host: {base_URL}
Content-Type: application/json
Accept: application/json
authorization: id_token
x-app-key: x-app-key

{
  "mode": "0011",
  "payerReference": "01723888888",
  "callbackURL": "yourdomain.com",
  "merchantAssociationInfo": "MI05MID54RF09123456One",
  "amount": "500",
  "currency": "BDT",
  "intent": "sale",
  "merchantInvoiceNumber": "Inv0124"
}
```

### Sample Response

```json
{
  "paymentID": "TR0011dQPHnuY1720518383420",
  "bkashURL": "https://sandbox.payment.bkash.com/?paymentId=TR0011dQPHnuY1720518383420&hash=...&mode=0011&apiVersion=v1.2.0-beta/",
  "callbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout",
  "successCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=success&signature=cm8HBfl65A",
  "failureCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=failure&signature=cm8HBfl65A",
  "cancelledCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=cancel&signature=cm8HBfl65A",
  "amount": "500",
  "intent": "sale",
  "currency": "BDT",
  "paymentCreateTime": "2024-07-09T15:46:23:420 GMT+0600",
  "transactionStatus": "Initiated",
  "merchantInvoiceNumber": "Inv0124",
  "statusCode": "0000",
  "statusMessage": "Successful"
}
```

## Execute Payment (Checkout)

Finalizes a checkout payment request.

**Request URL:** `{base_URL}/tokenized/checkout/execute`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Accept | application/json |
| Authorization | id_token from Grant/Refresh Token API |
| X-App-Key | App key from onboarding |

### Request Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| paymentID | Yes | string | PaymentID from Create Payment response |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code |
| paymentID | string | Payment ID |
| customerMsisdn | string | Customer MSISDN |
| payerReference | string | Payer reference |
| trxID | string | Transaction ID |
| transactionStatus | string | Final status (`Completed`) |
| amount | string | Transaction amount |
| currency | string | Currency (BDT) |
| intent | string | Payment intent |
| merchantInvoiceNumber | string | Invoice number |

### Sample Request

```json
POST /tokenized/checkout/execute HTTP/1.1
Host: {base_URL}
Accept: application/json
authorization: id_token
x-app-key: x-app-key

{
  "paymentID": "TR0011ON1565154754797"
}
```

### Sample Response

```json
{
  "statusCode": "0000",
  "statusMessage": "Successful",
  "paymentID": "TR0011ON1565154754797",
  "payerReference": "01770618575",
  "customerMsisdn": "01770618575",
  "trxID": "6H7801QFYM",
  "amount": "15",
  "transactionStatus": "Completed",
  "paymentExecuteTime": "2019-08-07T11:15:56:336 GMT+0600",
  "currency": "BDT",
  "intent": "sale",
  "merchantInvoiceNumber": "MER1231"
}
```
