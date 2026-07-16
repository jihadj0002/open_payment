# bKash Payment Gateway - Overview

bKash Online Payment Gateway provides well-structured and secure payment APIs for accepting payments from bKash customer accounts on Web or Mobile.

## Solutions

| Product | Description | Use Case |
|---------|-------------|----------|
| **Checkout** | Embedded payment experience with bKash hosted page | Standard e-commerce payments |
| **Tokenized Checkout** | Token-based payment to remember payer | Recurring customers, saved agreements |
| **Auth & Capture** | Pre-authorization to block amount and deduct later | Service confirmation workflows |
| **Instant Payout (B2C)** | Send funds directly to beneficiaries in real time | Disbursements, commissions, refunds |
| **Webhook (IPN)** | Real-time payment notifications | Server-side transaction status updates |

## API Version

- Current version: **v1.2.0-beta** (Inferno Dragon)
- Base URL (sandbox): `https://checkout.sandbox.bka.sh/v1.2.0-beta`
- Base URL (production): Provided during onboarding

## Integration Flows

1. **Checkout (URL Based)** - Create payment → Customer redirected to bKash page → Customer authorizes → Execute payment
2. **Tokenized Checkout** - Create agreement → Customer authorizes → Execute agreement → Create payment using agreement → Customer pays with PIN only → Execute payment
3. **Auth & Capture** - Create payment (authorization intent) → Customer authorizes → Execute authorization → Capture later
4. **B2C Payout** - Single disbursement to beneficiary wallet

## Authentication

All API calls require:
- `Authorization` header: ID token from Grant Token API
- `X-APP-Key` header: Application key provided during onboarding

## Environment

- Sandbox: `checkout.sandbox.bka.sh`
- Production: URL provided during merchant onboarding
- Currency: BDT only
