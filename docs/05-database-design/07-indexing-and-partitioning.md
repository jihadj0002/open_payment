# Indexing and Partitioning

> **Status:** ✅ Updated 2026-07-16
> **Note:** Partitioning is **not yet implemented**. The current schema uses regular tables with indexes. Partitioning (by time for `payment_intents`, `transactions`, `ledger_entries`, `audit_logs`) is planned for high-volume production use.

## Partitioning Strategy

### Partition by Time (payment_intents, transactions, refunds, ledger_entries)

```sql
-- Monthly partitioning for payment_intents
CREATE TABLE payment_intents_2026_01
    PARTITION OF payment_intents
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE payment_intents_2026_02
    PARTITION OF payment_intents
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

-- Automate partition creation with pg_partman
CREATE EXTENSION IF NOT EXISTS pg_partman;
SELECT partman.create_parent(
    p_parent_table := 'public.payment_intents',
    p_control := 'created_at',
    p_type := 'native',
    p_interval := '1 month',
    p_premake := 3
);
```

### Partition by Hash (ledger_entries for write distribution)

```sql
CREATE TABLE ledger_entries_0
    PARTITION OF ledger_entries
    FOR VALUES WITH (MODULUS 8, REMAINDER 0);
CREATE TABLE ledger_entries_1
    PARTITION OF ledger_entries
    FOR VALUES WITH (MODULUS 8, REMAINDER 1);
-- ... through 7
```

### Partition by Time (audit_logs, security_logs)

```sql
SELECT partman.create_parent(
    p_parent_table := 'public.audit_logs',
    p_control := 'created_at',
    p_type := 'native',
    p_interval := '1 month',
    p_premake := 3
);
```

## Indexing Strategy

### High-Write Tables (Optimize for INSERT + point lookups)

| Table | Index | Reason | Type |
|-------|-------|--------|------|
| payment_intents | `(merchant_id, created_at DESC)` | Merchant transaction list (most common query) | B-tree |
| payment_intents | `(idempotency_key)` WHERE idempotency_key IS NOT NULL | Idempotency lookup | Unique B-tree |
| payment_intents | `(status, created_at)` | Batch settlement processing (find capturable) | Partial B-tree |
| transactions | `(payment_id)` | Transaction history for a payment | B-tree |
| transactions | `(merchant_id, created_at DESC)` | Merchant transaction history | B-tree |
| refunds | `(payment_id)` | Refund lookup by payment | B-tree |
| ledger_entries | `(account_id, created_at DESC)` | Balance calculation | B-tree |
| ledger_entries | `(reference_type, reference_id)` | Find entries for a specific transaction | B-tree |

### High-Read Tables (Optimize for lookups)

| Table | Index | Reason | Type |
|-------|-------|--------|------|
| merchants | `(email)` | Login/lookup | Unique B-tree |
| merchants | `(status)` | Admin filtering | B-tree |
| api_keys | `(key_hash)` | API key lookup (authentication path) | Unique B-tree |
| customers | `(merchant_id, email)` | Customer lookup by merchant | Unique B-tree |
| sessions | `(refresh_token)` | Token refresh | Unique B-tree |

### Covering Indexes

```sql
-- Covering index for merchant transaction list (avoids table heap lookups)
CREATE INDEX idx_txn_merchant_covering
    ON transactions (merchant_id, created_at DESC)
    INCLUDE (id, payment_id, type, amount, currency, status);

-- Covering index for payment status check (most frequent read)
CREATE INDEX idx_pi_status_covering
    ON payment_intents (id)
    INCLUDE (merchant_id, status, amount, currency);
```

### Partial Indexes

```sql
-- Only index active/processing payments (a small fraction of total)
CREATE INDEX idx_pi_active
    ON payment_intents (merchant_id, created_at)
    WHERE status IN ('created', 'pending', 'processing', 'authorized', 'captured');

-- Only index failed webhook deliveries (for retry logic)
CREATE INDEX idx_wh_failed
    ON webhook_deliveries (merchant_id, created_at)
    WHERE status = 'failed' AND retry_count < 5;
```

### Concurrent Index Creation
For production migrations:
```sql
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pi_merchant_created
    ON payment_intents (merchant_id, created_at DESC);
```

## Query Patterns and Optimization

### Most Frequent Queries

| Query | Table | Expected Frequency | Optimization |
|-------|-------|--------------------|--------------|
| Get payment by ID | payment_intents | Very high | PK lookup |
| List merchant payments | payment_intents | High | Index (merchant_id, created_at) |
| Check balance | ledger_entries | High | Index (account_id, created_at) |
| Find capturable payments | payment_intents | Medium (batch) | Partial index on status |
| Validate API key | api_keys | Very high | Unique index on key_hash |
| Find customer by email | customers | Medium | Unique index (merchant_id, email) |
| List webhook logs | webhook_deliveries | Medium | Index (merchant_id, created_at) |
| Audit trail queries | audit_logs | Low | Index (actor_id, created_at) |

### Query Optimization Tips

1. **ALWAYS filter by merchant_id** in merchant-facing queries (partition elimination)
2. **Use LIMIT** for list queries (usually show 10–50 items)
3. **Avoid SELECT *** — fetch only needed columns
4. **Use cursors** for pagination on large tables (cursor-based > offset-based)
5. **Materialized views** for dashboard aggregates (refresh every 5 minutes)

## Migration Strategy

### Schema Changes

| Change Type | Strategy | Example |
|-------------|----------|---------|
| Add column | SAFE: ALTER TABLE ADD COLUMN with default | Adding metadata columns |
| Add index | SAFE: CREATE INDEX CONCURRENTLY | New query optimization |
| Remove column | SAFE: Mark as deprecated, remove in next release | Removing unused fields |
| Rename column | SAFE: Add new column, migrate data, drop old | Renaming `status` to `payment_status` |
| Change column type | RISKY: Add new column, backfill, drop old | Changing amount from INT to BIGINT |
| Split table | RISKY: Use views for backward compat | Extracting address from merchant table |
| Merge tables | RISKY: Create new table, migrate, drop old | -- |

### Migration Tool
- **golang-migrate** or **Flyway** for Go services
- Migrations are SQL files in `migrations/` directory
- Naming: `20260101_001_create_merchants.up.sql` / `20260101_001_create_merchants.down.sql`
- All migrations are tested in CI pipeline
- Rolling deployments: migration runs before new code deploys

## Data Archival

### Procedure (Monthly)
```sql
-- 1. Create archive table
CREATE TABLE payment_intents_archive_2025_12 (LIKE payment_intents INCLUDING ALL);

-- 2. Move old partitions
INSERT INTO payment_intents_archive_2025_12
SELECT * FROM payment_intents_2025_12;

-- 3. Detach and drop old partition
ALTER TABLE payment_intents DETACH PARTITION payment_intents_2025_12;
DROP TABLE payment_intents_2025_12;

-- 4. Dump to S3
-- pg_dump --table=payment_intents_archive_2025_12 | gzip | aws s3 cp - s3://bucket/archives/

-- 5. Drop archive table from PG (optional, can keep for 12 months hot)
```

### Retention Schedule
| Data | Hot (PG) | Warm (PG Archive) | Cold (S3) | Total Retention |
|------|----------|--------------------|------------|-----------------|
| payment_intents | 3 months | 12 months | 7 years | 7 years |
| transactions | 3 months | 12 months | 7 years | 7 years |
| refunds | 3 months | 12 months | 7 years | 7 years |
| ledger_entries | 3 months | 12 months | 7 years | 7 years |
| audit_logs | 6 months | 24 months | 7 years | 7 years |
| webhook_logs | 7 days | 30 days | 90 days | 90 days |
| sessions | Until expiry | None | None | Session TTL |
| password_reset_tokens | 24 hours | None | None | TTL |
