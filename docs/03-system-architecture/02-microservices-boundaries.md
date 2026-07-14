# Microservices Boundaries

## Service Inventory

| Service | Language | Responsibility | Owns Data | Depends On |
|---------|----------|----------------|-----------|------------|
| API Gateway | Go (Envoy) | Routing, auth, rate limit, TLS termination | None (stateless) | Auth Service |
| Auth Service | Go | Authentication, OAuth2, JWT, MFA, sessions | users, sessions, roles, permissions | PostgreSQL, Redis |
| Merchant Service | Go | Merchant onboarding, KYC, API keys, webhook config | merchants, merchant_users, api_keys, webhooks | PostgreSQL, S3 (docs) |
| Customer Service | Go | Customer profiles, saved payment methods, tokens | customers, saved_cards, tokens | PostgreSQL |
| Payment Service | Go | Payment intents, capture, refund, void, state machine | payment_intents, transactions, refunds | PostgreSQL, Kafka, Fraud Service, Processor Gateway |
| Fraud Service | Go + Python | Risk checks, velocity, ML scoring, rules engine | risk_rules, blacklist, risk_scores | PostgreSQL, Redis |
| Processor Gateway | Go | Adapter to card networks, wallet APIs, bank APIs | processor_logs | Kafka |
| Ledger Service | Go | Double-entry accounting, balance, reconciliation | ledger_accounts, ledger_entries, balances | PostgreSQL |
| Settlement Service | Go | Batch settlement, payout calculation, bank file generation | settlements, payout_batches | PostgreSQL, Kafka, S3 |
| Webhook Service | Go | Webhook delivery, retry, signature, logging | webhook_deliveries, webhook_logs | PostgreSQL, Kafka, Redis |
| Notification Service | Go | Email, SMS push notifications | notification_templates, notification_logs | Kafka, SES, SMS provider |
| Admin Service | Go | Admin dashboard APIs, merchant management, reporting | admin_users, audit_logs | PostgreSQL |
| Analytics Service | Python (optional) | Reporting, dashboards, data aggregation | materialized_views | PostgreSQL, ClickHouse |

## Service Boundaries Diagram

```
                    ┌─────────────────────┐
                    │    API Gateway       │
                    │  (Envoy/NGINX)       │
                    └─────────┬───────────┘
                              │
          ┌───────────────────┼───────────────────┐
          │                   │                   │
          ▼                   ▼                   ▼
   ┌──────────┐       ┌──────────────┐    ┌──────────────┐
   │   Auth   │       │   Merchant   │    │   Customer   │
   │ Service  │       │   Service    │    │   Service    │
   └────┬─────┘       └──────┬───────┘    └──────┬───────┘
        │                    │                    │
        └────────────────────┼────────────────────┘
                             │
                             ▼
                    ┌────────────────────┐
                    │   Payment Service  │
                    │  (Orchestrator)    │
                    └──┬──────┬──────┬───┘
                       │      │      │
          ┌────────────┘      │      └────────────┐
          ▼                   ▼                   ▼
   ┌──────────┐       ┌──────────────┐    ┌──────────────┐
   │  Fraud   │       │  Processor   │    │   Ledger     │
   │ Service  │       │  Gateway     │    │   Service    │
   └──────────┘       └──────────────┘    └──────────────┘
                                                 │
                                                 ▼
                                          ┌──────────────┐
                                          │  Settlement  │
                                          │   Service    │
                                          └──────────────┘
```

## Data Ownership Rules

| Data Entity | Owned By | Can Be Read By |
|-------------|----------|-----------------|
| Merchant profile | Merchant Service | Admin Service, Payment Service (limited) |
| API keys | Merchant Service | None (hashed only) |
| Customer data | Customer Service | Payment Service (by reference) |
| Payment intents | Payment Service | Merchant Service (by reference), Ledger Service |
| Transactions | Payment Service | Ledger Service, Settlement Service |
| Ledger entries | Ledger Service | Settlement Service, Reporting |
| Sessions | Auth Service | Auth Service only |
| Audit logs | Admin Service | Admin Service |

## Inter-Service Communication Contracts

Each service exposes both REST (external-facing) and gRPC (internal-facing):
- **External:** REST/JSON via API Gateway (port 443)
- **Internal:** gRPC/protobuf (port 9xxx, cluster-local)
- **Events:** Kafka topics (asynchronous)

### Kafka Topics
```
payment.created
payment.authorized
payment.captured
payment.failed
payment.refunded
payment.voided
payment.chargeback
ledger.entry.created
settlement.completed
payout.initiated
merchant.updated
fraud.alert
```
