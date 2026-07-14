# Database Architecture

## PostgreSQL Cluster Topology

```
           ┌──────────────────────┐
           │  Primary (writer)     │
           │  r6g.2xlarge         │
           │  8 vCPU, 64GB RAM    │
           │  1000 GB gp3 SSD     │
           └──────┬───────────────┘
                  │
        ┌─────────┼─────────┐
        │         │         │
   ┌────▼───┐ ┌───▼────┐ ┌───▼────┐
   │Replica1│ │Replica2│ │Replica3│
   │(reader)│ │(reader)│ │(reader)│
   │r6g.xl  │ │r6g.xl  │ │r6g.xl  │
   └────────┘ └────────┘ └────────┘
```

## Database per Service (Schema-per-Service)

| Service | Database Name | Tables |
|---------|---------------|--------|
| Auth | `auth_db` | users, sessions, roles, permissions, login_history |
| Merchant | `merchant_db` | merchants, merchant_users, merchant_settings, api_keys, webhook_configs, fee_configs |
| Customer | `customer_db` | customers, saved_cards, tokens, device_fingerprints |
| Payment | `payment_db` | payment_intents, transactions, refunds, disputes, chargebacks, payment_methods |
| Ledger | `ledger_db` | ledger_accounts, ledger_entries, balances, reconciliation_logs |
| Settlement | `settlement_db` | settlements, settlement_lines, payout_batches, payout_transactions |
| Webhook | `webhook_db` | webhook_deliveries, webhook_logs, webhook_retries |
| Fraud | `fraud_db` | risk_rules, risk_scores, blacklist, whitelist, fraud_events |
| Admin | `admin_db` | admin_users, audit_logs, support_tickets, system_config |

## Connection Pooling
- **PgBouncer** (or RDS Proxy) between services and PostgreSQL
- Pool mode: transaction-level pooling
- Max connections per pool: 25 (services), 50 (admin tools)

## Partitioning Strategy

### Partition by Month
```sql
-- payment_intents partitioned by created_at
CREATE TABLE payment_intents (
    id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
) PARTITION BY RANGE (created_at);

CREATE TABLE payment_intents_2026_01
    PARTITION OF payment_intents
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE payment_intents_2026_02
    PARTITION OF payment_intents
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
```

### Partition by Merchant ID (Hash)
```sql
-- ledger_entries partitioned by hash of account_id
CREATE TABLE ledger_entries (
    id UUID NOT NULL,
    account_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    entry_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
) PARTITION BY HASH (account_id);

CREATE TABLE ledger_entries_0
    PARTITION OF ledger_entries
    FOR VALUES WITH (MODULUS 4, REMAINDER 0);
```

## Indexing Strategy

### Critical Indexes
| Table | Index | Type | Reason |
|-------|-------|------|--------|
| payment_intents | (merchant_id, created_at DESC) | B-tree | Merchant transaction list |
| payment_intents | (idempotency_key) | Unique B-tree | Idempotency guarantee |
| transactions | (payment_id) | B-tree | Lookup by payment |
| refunds | (payment_id) | B-tree | Refund history for a payment |
| ledger_entries | (account_id, created_at DESC) | B-tree | Balance history |
| ledger_entries | (transaction_id) | B-tree | Source transaction lookup |
| api_keys | (key_hash) | Unique B-tree | API key lookup |
| webhook_deliveries | (merchant_id, created_at DESC) | B-tree | Webhook log for merchant |
| audit_logs | (actor_id, created_at DESC) | B-tree | Audit trail queries |

### Partial Indexes
```sql
-- Only index active payments
CREATE INDEX idx_active_payments ON payment_intents (merchant_id, created_at)
    WHERE status NOT IN ('settled', 'refunded', 'failed', 'cancelled');
```

## Data Archival
- **Active data:** Kept in partitioned tables for 12 months
- **Historical data:** Moved to archive tables (same schema) monthly
- **Archive storage:** S3 via pg_dump + pg_restore
- **Access to archives:** Via dedicated archive API or direct DB restore

## Backup Strategy
| Backup Type | Frequency | Retention | Notes |
|-------------|-----------|-----------|-------|
| WAL streaming | Continuous | 7 days | To S3 via pg_receivewal |
| Full backup | Daily | 30 days | pg_dump to S3 |
| Monthly snapshot | Monthly | 7 years | Compliance requirement |
| Point-in-time recovery | — | 7 days | WAL + base backup |

## Read Replica Usage
- Report queries run on replicas
- Dashboard queries (non-real-time) run on replicas
- Analytics/materialized views refresh from replicas
- Always route writes to primary via application logic
