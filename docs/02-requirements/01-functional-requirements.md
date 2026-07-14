# Functional Requirements

## FR-01: Merchant Onboarding
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-01.1 | Merchant submits registration form (business name, email, phone, address) | P0 | |
| FR-01.2 | System sends email verification link | P0 | |
| FR-01.3 | Merchant completes KYC (upload business documents, owner ID, proof of address) | P0 | Document storage in S3 |
| FR-01.4 | Admin reviews and approves/rejects merchant | P0 | Approval triggers API key generation |
| FR-01.5 | Merchant receives API keys (publishable + secret) after approval | P0 | Secret key shown only once |
| FR-01.6 | Merchant can reset API keys | P1 | Old keys invalidated immediately |
| FR-01.7 | Merchant can configure webhook URLs and events to subscribe to | P0 | |
| FR-01.8 | Merchant can add multiple users with role-based permissions | P1 | |

## FR-02: Payment Processing
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-02.1 | Merchant creates a Payment Intent via API | P0 | Idempotent |
| FR-02.2 | Customer is redirected to checkout or payment is processed server-side | P0 | |
| FR-02.3 | System validates card details (Luhn, expiry, CVV) | P0 | |
| FR-02.4 | System performs basic fraud checks before authorization | P0 | Velocity, IP reputation, BIN check |
| FR-02.5 | System authorizes payment with processor (card network/wallet API) | P0 | |
| FR-02.6 | On success, payment moves to `authorized` state | P0 | |
| FR-02.7 | Merchant can capture an authorized payment (full amount) | P0 | |
| FR-02.8 | Merchant can void an uncaptured authorization | P0 | |
| FR-02.9 | System settles captured payments daily (batch) | P0 | |
| FR-02.10 | Merchant receives webhook for each state change | P0 | |

## FR-03: Refunds
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-03.1 | Merchant can initiate full refund on captured payment | P0 | |
| FR-03.2 | Merchant can initiate partial refund | P0 | |
| FR-03.3 | System processes refund through processor | P0 | |
| FR-03.4 | Refunded amount is deducted from merchant balance | P0 | |
| FR-03.5 | Multiple partial refunds allowed until captured amount exhausted | P1 | |
| FR-03.6 | Refund history available in dashboard and API | P0 | |

## FR-04: Disputes and Chargebacks
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-04.1 | System receives chargeback notification from processor | P1 | |
| FR-04.2 | Merchant is notified of chargeback via webhook and dashboard | P1 | |
| FR-04.3 | Merchant can submit evidence documents to contest chargeback | P1 | |
| FR-04.4 | Admin tracks dispute lifecycle and resolution | P2 | |
| FR-04.5 | Chargeback amount is deducted from merchant balance | P1 | |

## FR-05: Balance and Payouts
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-05.1 | Merchant can view current balance and transaction history | P0 | |
| FR-05.2 | System calculates available balance (settled - fees - reserves - chargebacks) | P0 | |
| FR-05.3 | System initiates payout to merchant bank account on schedule (daily/weekly/monthly) | P0 | |
| FR-05.4 | Merchant can view payout history and download statements | P1 | |
| FR-05.5 | System supports multiple payout methods (bank transfer, wallet) | P2 | |

## FR-06: Webhooks
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-06.1 | System sends webhook events to merchant-configured endpoints | P0 | |
| FR-06.2 | Each webhook payload is signed with HMAC-SHA256 | P0 | |
| FR-06.3 | Webhook includes event type, timestamp, idempotency key | P0 | |
| FR-06.4 | System retries failed webhook delivery (exponential backoff, max 3 retries) | P0 | |
| FR-06.5 | Failed webhooks are logged and visible in merchant dashboard | P0 | |
| FR-06.6 | Merchant can manually replay webhooks | P1 | |
| FR-06.7 | Webhook delivery time < 5 seconds p99 | P0 | |

## FR-07: Admin Operations
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-07.1 | Admin can view all merchants and their status | P0 | |
| FR-07.2 | Admin can approve, suspend, or terminate merchant accounts | P0 | |
| FR-07.3 | Admin can view all transactions across merchants | P1 | |
| FR-07.4 | Admin can manually trigger settlement | P2 | |
| FR-07.5 | Admin can view audit logs | P0 | |
| FR-07.6 | Admin can configure system-wide fee structures | P1 | |
| FR-07.7 | Admin can manage currencies, countries, and supported banks | P2 | |

## FR-08: Reporting
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-08.1 | Merchant can view daily/monthly revenue summary | P1 | |
| FR-08.2 | Merchant can export transactions as CSV | P1 | |
| FR-08.3 | Merchant can view settlement reports | P1 | |
| FR-08.4 | Admin can view platform-wide revenue and growth metrics | P2 | |
| FR-08.5 | Admin can generate tax reports per merchant | P2 | |

## FR-09: API Management
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-09.1 | API supports idempotency via Idempotency-Key header | P0 | |
| FR-09.2 | API supports pagination via cursor-based or offset-limit | P0 | |
| FR-09.3 | API returns consistent error response format | P0 | |
| FR-09.4 | API rate limits per API key (configurable per merchant tier) | P0 | |
| FR-09.5 | API supports request signing (HMAC) | P0 | |
| FR-09.6 | Sandbox API mirrors production API exactly | P0 | |

## FR-10: Customer Management
| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-10.1 | Merchant can create customer profiles via API | P1 | |
| FR-10.2 | Merchant can save payment methods to customer profiles | P1 | Tokenized |
| FR-10.3 | Merchant can charge saved payment methods without re-entering details | P1 | |
| FR-10.4 | Customer payment history viewable by merchant | P1 | |
