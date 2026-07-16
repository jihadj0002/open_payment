# bKash - Webhooks (Instant Payment Notification)

Real-time payment notifications sent via AWS SNS (Simple Notification Service). Webhooks allow merchants to subscribe to payment events and receive HTTP POST payloads when transactions occur.

## Overview

bKash automatically notifies merchants of payment events via webhooks. Notifications are sent for payments received through:
- bKash Checkout or Direct products
- bKash Menu (*247#)
- bKash App
- bKash QR scan

## Configuration

1. **Set up a listener URL** - Create a server endpoint that accepts HTTP POST requests
2. **Share the listener URL** - Provide it to the bKash technical support team during onboarding
3. **Subscription confirmation** - Receive a subscription confirmation payload; acknowledge it to activate

## Payload Format

All notifications arrive as AWS SNS messages wrapped in JSON with a `Message` field containing the actual transaction data.

### Subscription Confirmation

```json
POST / HTTP/1.1
x-amz-sns-message-type: SubscriptionConfirmation
x-amz-sns-message-id: 165545c9-2a5c-472c-8df2-7ff2be2b3b1b
x-amz-sns-topic-arn: arn:aws:sns:us-west-2:123456789012:MyTopic
Content-Type: text/plain; charset=UTF-8

{
  "Type": "SubscriptionConfirmation",
  "MessageId": "165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
  "Token": "...",
  "TopicArn": "arn:aws:sns:us-west-2:123456789012:MyTopic",
  "Message": "You have chosen to subscribe...",
  "SubscribeURL": "https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&...",
  "Timestamp": "2012-04-26T20:45:04.751Z",
  "SignatureVersion": "1",
  "Signature": "...",
  "SigningCertURL": "https://sns.us-west-2.amazonaws.com/..."
}
```

### Payment Notification

```json
{
  "Type": "Notification",
  "MessageId": "20d48143-6af4-571d-b7cb-d211e6a2ac69",
  "TopicArn": "arn:aws:sns:ap-southeast-1:354285753755:bpt_01823072645",
  "Message": "{
    \"dateTime\": \"20180419122246\",
    \"debitMSISDN\": \"8801700000001\",
    \"creditOrganizationName\": \"Org 01\",
    \"creditShortCode\": \"01929918***\",
    \"trxID\": \"4J420ANOXC\",
    \"transactionStatus\": \"Completed\",
    \"transactionType\": \"1003\",
    \"amount\": \"100\",
    \"currency\": \"BDT\",
    \"transactionReference\": \"Test_Payment\",
    \"merchantInvoiceNumber\": \"orderId1233\"
  }",
  "Timestamp": "2018-04-19T12:22:46.236Z",
  "SignatureVersion": "1",
  "Signature": "...",
  "SigningCertURL": "https://sns.ap-southeast-1.amazonaws.com/...",
  "UnsubscribeURL": "https://sns.ap-southeast-1.amazonaws.com/..."
}
```

### Coupon Payment Notification

Includes extra fields for coupon-associated transactions:

```json
{
  "Message": "{
    \"amount\": \"90\",
    \"couponAmount\": \"10\",
    \"merchantShareAmount\": \"5\",
    \"saleAmount\": \"100.00\",
    \"currency\": \"BDT\",
    \"trxID\": \"4J420ANOXC\",
    \"transactionStatus\": \"Completed\",
    \"transactionType\": \"10002294\",
    ...
  }"
}
```

## Message Fields

| Field | Description |
|-------|-------------|
| dateTime | Transaction date/time |
| debitMSISDN | Customer wallet number |
| creditOrganizationName | Merchant organization name |
| creditShortCode | Merchant short code |
| trxID | Transaction ID |
| transactionStatus | Status (`Completed`) |
| transactionType | Channel type code |
| amount | Transaction amount (after coupon) |
| saleAmount | Full sale amount (with coupon) |
| couponAmount | Coupon discount amount |
| merchantShareAmount | Merchant's coupon contribution |
| currency | Currency (BDT) |
| transactionReference | User-provided reference |
| merchantInvoiceNumber | Merchant invoice number |

## Transaction Types

| Code | Channel |
|------|---------|
| 10002294 | Payment via API |
| 10003126 | Payment via QR |
| 10002175 | Payment via USSD |
| 10002809 | Redeem Voucher |
| 10002264 | M2M Transfer via API |
| 10003209 | M2M Transfer via QR |
| 10002177 | M2M Transfer via USSD |
| 10003476 | Payment via Bank |
| 10003237 | B2B Collection Wallet to Merchant Plus |
| 10003236 | Distributor to B2B Collection Wallet |
| 10004036 | DSO to Merchant Plus-B2BC via API |

## Notes

- Merchants receive notifications only for successfully completed payments
- For additional data (e.g., references), use the Search Transaction Details API
- Verify webhook signatures using AWS SNS signing certificates
