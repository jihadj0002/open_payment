# bKash - Instant Payout (B2C)

Send funds directly to beneficiaries in real time. Designed for instant disbursements to channel partners, not batch processing.

## Overview

The B2C Payout solution enables real-time fund transfer to beneficiaries. This is a single-disbursement API (not batch). For bulk disbursements, consult your bKash account manager.

## Workflow

1. Merchant initiates a payout request to a beneficiary wallet
2. bKash processes the transfer in real time
3. Beneficiary receives funds instantly

## API Endpoint

**Request URL:** `{base_URL}/b2c/payment`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Content-Type | application/json |
| Accept | application/json |
| Authorization | id_token from Grant/Refresh Token API |
| X-App-Key | App key from onboarding |

### Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| amount | string | Yes | Amount to disburse |
| currency | string | Yes | Currency (BDT) |
| merchantInvoiceNumber | string | Yes | Unique invoice number |
| receiverMSISDN | string | Yes | Beneficiary wallet number |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| statusCode | string | Status code |
| statusMessage | string | Status description |
| trxID | string | Transaction ID |
| amount | string | Disbursed amount |
| currency | string | Currency (BDT) |

## Use Cases

- Commission payouts to agents
- Refund disbursements
- Incentive payments
- Partner settlements
