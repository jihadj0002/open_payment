# Security Schema

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/database/migrations/007_admin_schema.up.sql`, `013_api_usage_logs.up.sql`

## audit_logs (migration 007 — implemented)

```sql
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID NOT NULL,
    action          VARCHAR(50) NOT NULL,
    resource_type   VARCHAR(50) NOT NULL,
    resource_id     VARCHAR(100),
    details         TEXT,
    ip_address      VARCHAR(45),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

> **Differences from spec:** No `actor_type`, no partitioning, `details` is TEXT, no `request_id`, `ip_address` is VARCHAR(45) not INET.

## system_config (migration 007 — implemented)

```sql
CREATE TABLE system_config (
    id                     VARCHAR(20) PRIMARY KEY DEFAULT 'default',
    supported_currencies   JSONB NOT NULL DEFAULT '["BDT","USD"]',
    supported_countries    JSONB NOT NULL DEFAULT '["BD"]',
    max_transaction_amount BIGINT NOT NULL DEFAULT 10000000,
    maintenance_mode       BOOLEAN NOT NULL DEFAULT false,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## api_usage_logs (migration 013 — implemented)

```sql
CREATE TABLE api_usage_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id VARCHAR(255) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## admin_users
> **TODO: Not implemented** — No table exists in migrations. Admin auth is handled via JWT role checking.

## roles
> **TODO: Not implemented** — No table exists in migrations.

## permissions
> **TODO: Not implemented** — No table exists in migrations.

## role_permissions
> **TODO: Not implemented** — No table exists in migrations.

## sessions
> **TODO: Not implemented** — No table exists in migrations. Session management is handled via JWT tokens.

## security_logs
> **TODO: Not implemented** — No table exists in migrations.

## login_history
> **TODO: Not implemented** — No table exists in migrations.

## api_key_history
> **TODO: Not implemented** — No table exists in migrations.

## fraud_rules
> **TODO: Not implemented** — No `fraud_rules` table in migrations. Fraud rules are managed at application level via `fraud_configs` table (migration 005) and the fraud rule CRUD endpoints.
