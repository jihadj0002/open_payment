# Merchant Schema

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/database/migrations/001_initial_schema.up.sql`, `003_add_merchant_password.up.sql`, `010_encrypt_merchant_keys.up.sql`, `016_password_reset_tokens.up.sql`

## merchants (migration 001 + 003 + 010)

The actual `merchants` table is simpler than the spec. The following columns **exist** in the migration:

```sql
CREATE TABLE merchants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,       -- spec uses business_name
    email           VARCHAR(255) NOT NULL UNIQUE,
    secret_key      VARCHAR(64) NOT NULL UNIQUE,  -- encrypted at app level
    public_key      VARCHAR(64) NOT NULL UNIQUE,  -- encrypted at app level
    webhook_url     TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Added in migration 003:
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
```

> **Differences from spec:** Column is `name` not `business_name`. No `phone`, `website`, `description`, `verification_status`, `country`, `currency`, `timezone`, `kyc_documents` columns exist. `secret_key` and `public_key` are stored (encrypted at application level via migration 010). Keys added via `password_hash` column.

## merchant_users
> **TODO: Not implemented** — No table exists in migrations.

## merchant_settings
> **TODO: Not implemented** — No table exists in migrations.

## api_keys (migration 001)

```sql
CREATE TABLE api_keys (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    key_prefix      VARCHAR(8) NOT NULL,
    key_hash        VARCHAR(64) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    permissions     JSONB NOT NULL DEFAULT '["read"]',
    last_used_at    TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

> **Differences from spec:** No `key_last4`, `mode` enum, `allowed_ips`, or `is_active` columns. Uses `status` instead. Permissions are JSONB array.

## webhooks (migration 001 + 012)

```sql
CREATE TABLE webhooks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    event               VARCHAR(50) NOT NULL,
    url                 TEXT NOT NULL,
    secret              VARCHAR(64) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Added in migration 012 (webhook secret rotation):
ALTER TABLE webhooks ADD COLUMN IF NOT EXISTS previous_secret VARCHAR(64);
ALTER TABLE webhooks ADD COLUMN IF NOT EXISTS previous_secret_expires_at TIMESTAMPTZ;
```

## webhook_deliveries (migration 001)

```sql
CREATE TABLE webhook_deliveries (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id          UUID NOT NULL REFERENCES webhooks(id),
    event               VARCHAR(50) NOT NULL,
    payload             JSONB NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',
    attempt             INT NOT NULL DEFAULT 0,
    max_attempts        INT NOT NULL DEFAULT 3,
    response_code       INT,
    response_body       TEXT,
    next_attempt_at     TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## fee_configs (migration 007 — admin schema)

```sql
CREATE TABLE fee_configs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    rate                NUMERIC(5,2) NOT NULL DEFAULT 2.2,
    fixed_fee           BIGINT NOT NULL DEFAULT 0,
    cross_border_rate   NUMERIC(5,2),
    monthly_fee         BIGINT NOT NULL DEFAULT 0,
    min_monthly_volume  BIGINT NOT NULL DEFAULT 0,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

> **Note:** `fee_configs` is system-level (not per-merchant) in current implementation. No `merchant_id` column exists.

## kyc_documents
> **TODO: Not implemented** — No table exists in migrations.

## password_reset_tokens (migration 016)

```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
