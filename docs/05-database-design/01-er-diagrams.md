# Entity Relationship Diagrams

> **Status:** ✅ Updated 2026-07-16
> **Note:** These ER diagrams represent the **target architecture**. Some entities shown are **not yet migrated** — see individual schema docs for implementation status.

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
- Merchant 1:N → MerchantUser (**TODO**), APIKey ✅, Webhook ✅, FeeConfig ✅, Customer ✅, PaymentIntent ✅
- Customer 1:N → PaymentIntent ✅, SavedPaymentMethod ✅
- PaymentIntent 1:N → Transaction ✅, StatusHistory ✅
- PaymentIntent 1:1 → Dispute ✅ (optional)
- LedgerAccount 1:N → LedgerEntry ⚠️ (simplified — no LedgerAccount table)
- Merchant 1:N → SettlementBatch 1:N → SettlementLine (**TODO** — single settlements table)
- SettlementBatch 1:1 → PayoutBatch (**TODO**)
- AdminUser 1:N → AuditLog (**TODO** — admin_users not implemented)
- PaymentIntent 1:N → FraudCheck ✅ | FraudEvent (**TODO**)

## Entity Implementation Status
| Entity | Status | Migration |
|--------|--------|-----------|
| Merchant | ✅ | 001 |
| Customer | ✅ | 001 |
| PaymentIntent | ✅ | 001 |
| Transaction | ✅ | 001 |
| APIKey | ✅ | 001 |
| Webhook | ✅ | 001 |
| WebhookDelivery | ✅ | 001 |
| LedgerEntry | ✅ | 001 |
| PaymentMethod | ✅ | 004 |
| FraudCheck | ✅ | 005 |
| FraudConfig | ✅ | 005 |
| Settlement | ✅ | 006 |
| FeeConfig | ✅ | 007 |
| SystemConfig | ✅ | 007 |
| AuditLog | ✅ | 007 |
| Dispute | ✅ | 009 |
| StatusHistory | ✅ | 011 |
| ApiUsageLog | ✅ | 013 |
| SavedPaymentMethod | ✅ | 014 |
| PasswordResetToken | ✅ | 016 |
| MerchantUser | ❌ TODO | — |
| MerchantSettings | ❌ TODO | — |
| KycDocument | ❌ TODO | — |
| AdminUser | ❌ TODO | — |
| Session | ❌ TODO | — |
| Role/Permission | ❌ TODO | — |
| SecurityLog | ❌ TODO | — |
| LoginHistory | ❌ TODO | — |
| FraudRule | ❌ TODO | — |
| Refund (standalone) | ❌ TODO (uses transactions) | — |
| Chargeback (standalone) | ❌ TODO (uses disputes) | — |
| SettlementBatch/Line | ❌ TODO | — |
| PayoutBatch | ❌ TODO | — |
| Token | ❌ TODO | — |
| DeviceFingerprint | ❌ TODO | — |
| RiskScore | ❌ TODO | — |
