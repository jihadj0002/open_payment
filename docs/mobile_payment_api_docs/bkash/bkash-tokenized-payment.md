# bKash - Tokenized Checkout: Create & Execute Payment

Create and execute payments using a pre-existing agreement (tokenized flow). Customers pay using only their wallet PIN without re-entering credentials.

## Create Payment (Tokenized)

Creates a payment request using an existing agreement.

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
| mode | Yes | string | Must be `"0001"` for tokenized payment |
| payerReference | Yes | string | Reference value (max 255 chars) |
| callbackURL | Yes | string | Base URL for callback URLs |
| agreementID | Yes | string | Agreement ID from Execute Agreement response |
| amount | Yes | string | Payment amount |
| currency | Yes | string | Currency (BDT only) |
| intent | Yes | string | Must be `"sale"` for tokenized checkout |
| merchantInvoiceNumber | Yes | string | Unique invoice number (max 255 chars) |
| merchantAssociationInfo | No | string | TLV format for aggregators/system integrators |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code |
| statusMessage | string | Status description |
| paymentID | string | bKash payment ID (expires after 24 hours) |
| agreementID | string | Agreement ID used |
| bkashURL | string | URL to redirect customer for PIN entry |
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
  "agreementID": "TokenizedMerchant01L3IKB6H1565072174986",
  "mode": "0001",
  "payerReference": "01723888888",
  "callbackURL": "yourDomain.com",
  "merchantAssociationInfo": "MI05MID54RF09123456One",
  "amount": "12",
  "currency": "BDT",
  "intent": "sale",
  "merchantInvoiceNumber": "Inv0124"
}
```

### Sample Response

```json
{
  "statusCode": "0000",
  "statusMessage": "Successful",
  "paymentID": "TR0001VK1565072365492",
  "callbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout",
  "successCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=success&signature=cm8HBfl65A",
  "failureCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=failure&signature=cm8HBfl65A",
  "cancelledCallbackURL": "https://yourdomain.com/callback?version=v1.2.0-beta&product=tokenized-checkout&paymentID=TR0011dQPHnuY1720518383420&status=cancel&signature=cm8HBfl65A",
  "amount": "500",
  "intent": "sale",
  "currency": "BDT",
  "agreementID": "TokenizedMerchant01L3IKB6H1565072174986",
  "paymentCreateTime": "2019-08-06T12:19:25:593 GMT+0600",
  "transactionStatus": "Initiated",
  "merchantInvoiceNumber": "Inv0124"
}
```

## Execute Payment (Tokenized)

Finalizes a tokenized payment request after the customer enters their PIN.

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
| statusMessage | string | Status description |
| paymentID | string | Payment ID |
| agreementID | string | Agreement ID |
| customerMsisdn | string | Customer MSISDN |
| payerReference | string | Payer reference |
| trxID | string | Transaction ID for this payment |
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
  "paymentID": "TR0001VK1565072365492"
}
```

### Sample Response

```json
{
  "statusCode": "0000",
  "statusMessage": "Successful",
  "paymentID": "TR0001VK1565072365492",
  "agreementID": "TokenizedMerchant01L3IKB6H1565072174986",
  "payerReference": "01770618575",
  "customerMsisdn": "01770618575",
  "trxID": "6H6201QDIY",
  "amount": "12",
  "transactionStatus": "Completed",
  "paymentExecuteTime": "2019-08-06T12:22:41:428 GMT+0600",
  "currency": "BDT",
  "intent": "sale",
  "merchantInvoiceNumber": "TestForOnmobile"
}
```
