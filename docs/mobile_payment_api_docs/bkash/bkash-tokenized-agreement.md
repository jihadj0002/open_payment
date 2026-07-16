# bKash - Tokenized Checkout: Create & Execute Agreement

Tokenized Checkout allows merchants to create an agreement with a customer's wallet so the customer can pay using only their wallet PIN on subsequent purchases.

## Create Agreement

Creates an agreement request for a customer wallet.

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
| mode | Yes | string | Must be `"0000"` for agreement creation |
| payerReference | Yes | string | Reference value (max 255 chars). Wallet number pre-populates bKash page |
| callbackURL | Yes | string | Base URL for generating success/failure/cancelled callback URLs |
| amount | Yes | string | Payment amount |
| currency | Yes | string | Currency (BDT only) |
| intent | Yes | string | Must be `"Sale"` for tokenized checkout |
| merchantInvoiceNumber | No | string | Unique invoice number (max 255 chars) |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code (`0000` = success) |
| statusMessage | string | Status description |
| paymentID | string | bKash payment ID (expires after 24 hours) |
| bkashURL | string | URL to redirect customer for wallet entry |
| callbackURL | string | Base callback URL |
| successCallbackURL | string | Success callback URL |
| failureCallbackURL | string | Failure callback URL |
| cancelledCallbackURL | string | Cancelled callback URL |
| payerReference | string | Payer reference value |
| agreementStatus | string | Agreement status (Initiated) |
| agreementCreateTime | string | Agreement creation timestamp |

### Sample Request

```json
POST /tokenized/checkout/create HTTP/1.1
Host: {base_URL}
Content-Type: application/json
Accept: application/json
authorization: id_token
x-app-key: x-app-key

{
  "mode": "0000",
  "callbackURL": "yourdomain.com",
  "payerReference": "0173499999"
}
```

### Sample Response

```json
{
  "statusCode": "0000",
  "statusMessage": "Successful",
  "paymentID": "TR00008C1565071974689",
  "bkashURL": "https://sandbox.payment.bkash.com/redirect/tokenized/?paymentID=TR00008C1565071974689&hash=...&mode=0000&apiVersion=v1.2.0-beta",
  "callbackURL": "yourdomain.com",
  "successCallbackURL": "yourdomain.com?paymentID=TR00008C1565071974689&status=success",
  "failureCallbackURL": "yourdomain.com?paymentID=TR00008C1565071974689&status=failure",
  "cancelledCallbackURL": "yourdomain.com?paymentID=TR00008C1565071974689&status=cancel",
  "payerReference": "01770618575",
  "agreementStatus": "Initiated",
  "agreementCreateTime": "2019-08-06T12:12:54:728 GMT+0600"
}
```

## Execute Agreement

Finalizes an agreement creation request after the customer enters their credentials on the bKash page.

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
| paymentID | Yes | string | PaymentID from Create Agreement response |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code |
| statusMessage | string | Status description |
| paymentID | string | Payment ID |
| agreementID | string | Generated agreement ID (store at merchant end) |
| customerMsisdn | string | Customer's MSISDN |
| payerReference | string | Payer reference value |
| agreementExecuteTime | string | Execution timestamp |
| agreementStatus | string | Status (Completed) |

### Sample Request

```json
POST /tokenized/checkout/execute HTTP/1.1
Host: {base_URL}
Accept: application/json
authorization: id_token
x-app-key: x-app-key

{
  "paymentID": "TR00008C1565071974689"
}
```

### Sample Response

```json
{
  "statusCode": "0000",
  "statusMessage": "Successful",
  "paymentID": "TR00008C1565071974689",
  "agreementID": "TokenizedMerchant01L3IKB6H1565072174986",
  "payerReference": "01770618575",
  "agreementExecuteTime": "2019-08-06T12:16:14:985 GMT+0600",
  "agreementStatus": "Completed",
  "customerMsisdn": "01770618575"
}
```

## Important Notes

- A payment ID expires after **24 hours** if not executed
- A payment ID is valid for **one execution only**
- Payment IDs can only be used for query and search after use
- Store the `agreementID` at the merchant end for future payments
