# Payment Schema

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/database/migrations/001_initial_schema.up.sql`, `008_add_payment_columns.up.sql`, `009_disputes.up.sql`, `011_status_history.up.sql`, `014_saved_payment_methods.up.sql`

## payment_intents (migration 001 + 008)

```sql
CREATE TABLE payment_intents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    customer_id         UUID REFERENCES customers(id),
    amount              BIGINT NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'BDT',
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',
    idempotency_key     VARCHAR(64) UNIQUE,
    description         TEXT,
    metadata            JSONB DEFAULT '{}',
    failure_reason      TEXT,
    processor           VARCHAR(20),
    processor_ref       VARCHAR(255),
    live_mode           BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Added in migration 008:
ALTER TABLE payment_intents
    ADD COLUMN IF NOT EXISTS amount_capturable BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS amount_received BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS capture_method VARCHAR(20) NOT NULL DEFAULT 'automatic';
```

> **Differences from spec:** No `payment_status` enum, `capture_method` enum, `client_secret`, `return_url`, `cancel_url`, `error_code`, `error_message` columns. No `PARTITION BY RANGE`. Status is `VARCHAR(20)`.

## transactions (migration 001 + 008)

```sql
CREATE TABLE transactions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_intent_id   UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    type                VARCHAR(20) NOT NULL,
    amount              BIGINT NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'BDT',
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',
    processor_ref       VARCHAR(255),
    processor_response  JSONB DEFAULT '{}',
    fee                 BIGINT DEFAULT 0,
    net_amount          BIGINT DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Added in migration 008:
ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(64);
```

> **Differences from spec:** Column is `payment_intent_id` not `payment_id`. No type/status enums, no `authorization_code`/`failure_code`/`failure_message` columns. No `PARTITION BY RANGE`.

## refunds
> **TODO: Not implemented** — No separate `refunds` table. Refunds are recorded as transactions with `type = 'refund'` in the transactions table.

## disputes (migration 009)

```sql
CREATE TABLE disputes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_intent_id   UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    amount              BIGINT NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'BDT',
    reason              VARCHAR(255) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'needs_response',
    evidence_due_by     TIMESTAMPTZ,
    evidence_submitted_at TIMESTAMPTZ,
    metadata            JSONB DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## status_history (migration 011)

```sql
CREATE TABLE status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_intent_id UUID NOT NULL REFERENCES payment_intents(id) ON DELETE CASCADE,
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    changed_by VARCHAR(255) NOT NULL DEFAULT 'system',
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## payment_methods (migration 004)

```sql
CREATE TABLE payment_methods (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    customer_id     UUID NOT NULL REFERENCES customers(id),
    type            VARCHAR(20) NOT NULL,
    token           VARCHAR(255),
    last4           VARCHAR(4),
    brand           VARCHAR(20),
    exp_month       INT,
    exp_year        INT,
    cardholder_name VARCHAR(255),
    fingerprint     VARCHAR(64),
    wallet_type     VARCHAR(50),
    wallet_phone    VARCHAR(20),
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## saved_payment_methods (migration 014)

```sql
CREATE TABLE saved_payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id VARCHAR(255) NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    merchant_id VARCHAR(255) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('card', 'wallet')),
    last4 VARCHAR(4) NOT NULL,
    brand VARCHAR(50) NOT NULL DEFAULT '',
    exp_month INTEGER,
    exp_year INTEGER,
    is_default BOOLEAN NOT NULL DEFAULT false,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(customer_id, token)
);
```

## chargebacks
> **TODO: Not implemented** — No separate `chargebacks` table. Chargeback data is tracked via the `disputes` table.
```sql
CREATE TYPE dispute_status AS ENUM (
    'needs_response', 'under_review', 'won', 'lost', 'closed'
);

CREATE TYPE dispute_reason AS ENUM (
    'fraudulent', 'duplicate', 'subscription_canceled',
    'product_not_received', 'product_unacceptable', 'credit_not_processed',
    'general'
);

CREATE TABLE disputes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id      UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    amount          BIGINT NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    status          dispute_status NOT NULL DEFAULT 'needs_response',
    reason          dispute_reason NOT NULL,
    evidence        JSONB,
    respond_by      TIMESTAMPTZ NOT NULL,
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_disputes_merchant ON disputes (merchant_id);
CREATE INDEX idx_disputes_status ON disputes (status);
```

## chargebacks
```sql
CREATE TABLE chargebacks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dispute_id      UUID REFERENCES disputes(id),
    payment_id      UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    amount          BIGINT NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    reason_code     VARCHAR(50),
    case_id         VARCHAR(100),
    status          VARCHAR(20) NOT NULL DEFAULT 'open',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at     TIMESTAMPTZ
);

CREATE INDEX idx_cb_payment ON chargebacks (payment_id);
```

## payment_methods (normalized card data)
```sql
CREATE TABLE payment_methods (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    customer_id     UUID REFERENCES customers(id),
    type            VARCHAR(20) NOT NULL CHECK (type IN ('card', 'wallet', 'bank_account')),
    -- Card-specific fields
    token           VARCHAR(100),               -- network token or gateway token
    last4           VARCHAR(4),
    brand           VARCHAR(20),                -- visa, mastercard, amex
    exp_month       INTEGER,
    exp_year        INTEGER,
    cardholder_name VARCHAR(255),
    fingerprint     VARCHAR(64),                -- unique card fingerprint (for dedup)
    -- Wallet-specific
    wallet_type     VARCHAR(20),                -- bKash, Nagad, PayPal
    wallet_phone    VARCHAR(20),
    -- Bank-specific
    bank_account_number VARCHAR(50),
    bank_routing_number VARCHAR(50),
    bank_name       VARCHAR(100),
    -- Metadata
    is_default      BOOLEAN NOT NULL DEFAULT false,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pm_merchant ON payment_methods (merchant_id);
CREATE INDEX idx_pm_customer ON payment_methods (customer_id);
CREATE INDEX idx_pm_fingerprint ON payment_methods (fingerprint);
```
