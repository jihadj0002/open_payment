# Nagad REST API Integration

This guide covers the REST API integration for Nagad Payment Gateway, including the Paykassma aggregation flow and direct Nagad API endpoints.

## Prerequisites

1. Register on the Nagad merchant portal
2. Obtain your `merchantID`, `merchantPrivateKey`, and `pgPublicKey`
3. Whitelist your server IP address with Nagad
4. Configure your callback URL for payment notifications

## Payment Flow (via Paykassma)

1. Your platform calls **Create Payment** to get Nagad wallet details
2. Customer receives Nagad wallet details and pays in the Nagad app
3. Customer enters the 10-character transaction ID in the payment window
4. Your platform sends an **Activation** request to confirm payment
5. Paykassma sends a signed postback to your server
6. Funds are credited to the user's balance

## API Endpoints

### Create Payment

| Detail | Value |
|--------|-------|
| **Endpoint** | `POST /api/v1/transaction/create/nagad` |
| **Purpose** | Get Nagad wallet details for payment |
| **Currency** | `BDT` |

**Request parameters:**
- `currency`: `BDT`
- User label identifying the customer

### Activate Payment

| Detail | Value |
|--------|-------|
| **Endpoint** | `POST /api/v1/transaction/activate` |
| **Purpose** | Confirm payment with transaction ID |

**Request parameters:**
- `wallet_type`: `nagad`
- `key1`: Customer's 10-character transaction ID
- `amount`: Transaction amount

### Deposit Postback

| Detail | Value |
|--------|-------|
| **Endpoint** | Your callback URL |
| **Method** | `POST` |
| **Purpose** | Receive signed transaction status |

Your callback must respond with `{"status":"ok"}` upon receipt.

### Withdrawal (Payout)

| Detail | Value |
|--------|-------|
| **Endpoint** | `POST /v2/withdrawal/create` |
| **Purpose** | Payout to Nagad account |

**Request parameters:**
- `payment_system`: `nagad`

## Integration Options

### 1. Payment Plugin (iFrame)

Use an embeddable iFrame for quick integration:
- `wallet_type`: `nagad`
- `currency_code`: `BDT`
- Available languages: `en`, `bn`

### 2. REST API

Full control over the create/activate payment flow. Recommended for custom platforms.

### 3. Postback

Configure callback URLs via Paykassma support for receiving transaction status notifications.

## Account Types Supported

| Account Type | Description |
|-------------|-------------|
| Personal | Send Money |
| Agent | Cash Out |
| Merchant | Make Payment |

## Sandbox Testing

1. Set `isSandbox: true` in your credentials
2. Use test `merchantID` and `publicKey` provided by Nagad
3. Simulate all scenarios:
   - Successful payment
   - Cancellation
   - Insufficient funds
4. Verify signatures and callback handling

## Production Deployment

1. Change to production credentials
2. Set `isSandbox: false`
3. Ensure SSL/HTTPS is properly configured
4. Test the complete flow with real transactions (minimal amounts)
5. Monitor callback URLs for postback reliability

## Security Considerations

- Always use HTTPS for all API calls
- Verify postback signatures from Nagad/Paykassma
- Never expose merchant credentials in client-side code
- Store transaction details in your database for reconciliation
- Whitelist Nagad IP addresses on your server
