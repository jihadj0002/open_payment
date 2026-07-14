# Entity Relationship Diagrams

## Core Domain

```
┌──────────────┐       ┌──────────────────┐
│   Merchant   │1────N│  MerchantUser     │
├──────────────┤       ├──────────────────┤
│ id (PK)      │       │ id (PK)          │
│ business_name│       │ merchant_id (FK) │
│ email        │       │ email            │
│ phone        │       │ password_hash    │
│ status       │       │ role             │
│ verification │       │ mfa_enabled      │
│ country      │       │ created_at       │
│ currency     │       └──────────────────┘
│ created_at   │
└──────┬───────┘
       │
       │1
       │
       ├──────────────┐
       │              │
       │1             │1
       ▼              ▼
┌──────────────┐ ┌──────────────────┐
│  APIKey      │ │  WebhookConfig   │
├──────────────┤ ├──────────────────┤
│ id (PK)      │ │ id (PK)          │
│ merchant_id  │ │ merchant_id (FK) │
│ name         │ │ url              │
│ key_prefix   │ │ events           │
│ key_hash     │ │ secret           │
│ key_last4    │ │ is_active        │
│ mode         │ │ created_at       │
│ expires_at   │ └──────────────────┘
│ is_active    │
│ created_at   │
└──────────────┘
```

## Payment Domain

```
┌──────────────────┐       ┌──────────────────┐
│  Customer        │1────N│  PaymentIntent   │
├──────────────────┤       ├──────────────────┤
│ id (PK)          │       │ id (PK)          │
│ merchant_id (FK) │       │ merchant_id (FK) │
│ email            │       │ customer_id (FK) │
│ name             │       │ amount           │
│ phone            │       │ currency         │
│ created_at       │       │ status           │
└──────────────────┘       │ capture_method   │
                           │ idempotency_key  │
      1                    │ description      │
       │                   │ metadata         │
       │                   │ client_secret    │
       │                   │ created_at       │
       │                   │ updated_at       │
       │                   └────────┬─────────┘
       │                            │
       │                            │1
       │                            │
       │                    ┌───────┴──────────┐
       │                    │                  │
       ▼                    ▼                  ▼
┌──────────────────┐ ┌──────────────┐ ┌──────────────┐
│  SavedCard       │ │  Transaction │ │   Refund     │
├──────────────────┤ ├──────────────┤ ├──────────────┤
│ id (PK)          │ │ id (PK)      │ │ id (PK)      │
│ customer_id (FK) │ │ payment_id   │ │ payment_id   │
│ token            │ │ type         │ │ amount       │
│ last4            │ │ amount       │ │ currency     │
│ brand            │ │ currency     │ │ status       │
│ exp_month        │ │ status       │ │ reason       │
│ exp_year         │ │ processor_id │ │ processor_id │
│ is_default       │ │ created_at   │ │ created_at   │
│ created_at       │ └──────────────┘ └──────────────┘
└──────────────────┘
```

## Ledger Domain

```
┌──────────────────┐       ┌──────────────────┐
│  LedgerAccount   │1────N│  LedgerEntry     │
├──────────────────┤       ├──────────────────┤
│ id (PK)          │       │ id (PK)          │
│ merchant_id (FK) │       │ account_id (FK)  │
│ type             │       │ transaction_id   │
│ currency         │       │ amount           │
│ created_at       │       │ direction        │
└──────────────────┘       │ entry_type       │
                           │ description      │
                           │ reference_type   │
                           │ reference_id     │
                           │ created_at       │
                           └──────────────────┘
```

## Settlement Domain

```
┌──────────────────┐       ┌──────────────────┐
│  SettlementBatch │1────N│  SettlementLine  │
├──────────────────┤       ├──────────────────┤
│ id (PK)          │       │ id (PK)          │
│ merchant_id (FK) │       │ batch_id (FK)    │
│ period_start     │       │ payment_id (FK)  │
│ period_end       │       │ amount           │
│ total_amount     │       │ fee_amount       │
│ fee_amount       │       │ tax_amount       │
│ tax_amount       │       │ net_amount       │
│ net_amount       │       │ created_at       │
│ status           │       └──────────────────┘
│ settled_at       │
└──────────────────┘
       │
       │1
       ▼
┌──────────────────┐
│  PayoutBatch     │
├──────────────────┤
│ id (PK)          │
│ merchant_id (FK) │
│ amount           │
│ fee_amount       │
│ net_amount       │
│ bank_account     │
│ status           │
│ initiated_at     │
│ completed_at     │
└──────────────────┘
```

## Security Domain

```
┌──────────────┐       ┌──────────────────┐
│  AdminUser   │1────N│  AuditLog        │
├──────────────┤       ├──────────────────┤
│ id (PK)      │       │ id (PK)          │
│ email        │       │ actor_id         │
│ role         │       │ actor_type       │
│ mfa_enabled  │       │ action           │
│ created_at   │       │ resource_type    │
└──────────────┘       │ resource_id      │
                       │ details          │
┌──────────────┐       │ ip_address       │
│  FraudEvent  │       │ user_agent       │
├──────────────┤       │ created_at       │
│ id (PK)      │       └──────────────────┘
│ payment_id   │
│ rule_triggered│
│ risk_score   │
│ action_taken │
│ created_at   │
└──────────────┘
```

## Key Relationships Summary
- Merchant 1:N → MerchantUser, APIKey, WebhookConfig, FeeConfig, Customer, PaymentIntent
- Customer 1:N → PaymentIntent, SavedCard
- PaymentIntent 1:N → Transaction, Refund
- PaymentIntent 1:1 → Dispute (optional)
- LedgerAccount 1:N → LedgerEntry
- Merchant 1:N → SettlementBatch 1:N → SettlementLine
- SettlementBatch 1:1 → PayoutBatch (optional)
- AdminUser 1:N → AuditLog
- PaymentIntent 1:N → FraudEvent
