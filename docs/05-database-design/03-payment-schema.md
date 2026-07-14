# Payment Schema

## payment_intents
```sql
CREATE TYPE payment_status AS ENUM (
    'created', 'pending', 'processing', 'authorized',
    'captured', 'succeeded', 'failed', 'canceled',
    'expired', 'refunded', 'partially_refunded'
);

CREATE TYPE capture_method AS ENUM ('automatic', 'manual');

CREATE TABLE payment_intents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    customer_id         UUID REFERENCES customers(id),
    amount              BIGINT NOT NULL CHECK (amount > 0),
    amount_capturable   BIGINT NOT NULL,
    amount_received     BIGINT NOT NULL DEFAULT 0,
    currency            VARCHAR(3) NOT NULL DEFAULT 'BDT',
    status              payment_status NOT NULL DEFAULT 'created',
    capture_method      capture_method NOT NULL DEFAULT 'automatic',
    payment_method      VARCHAR(20),
    idempotency_key     VARCHAR(64) UNIQUE,
    description         VARCHAR(255),
    metadata            JSONB DEFAULT '{}',
    client_secret       VARCHAR(64) NOT NULL,
    return_url          VARCHAR(500),
    cancel_url          VARCHAR(500),
    error_code          VARCHAR(50),
    error_message       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_pi_merchant_status ON payment_intents (merchant_id, status);
CREATE INDEX idx_pi_customer ON payment_intents (customer_id);
CREATE INDEX idx_pi_created_at ON payment_intents (created_at DESC);
CREATE UNIQUE INDEX idx_pi_idempotency ON payment_intents (idempotency_key) WHERE idempotency_key IS NOT NULL;
```

## transactions
```sql
CREATE TYPE transaction_type AS ENUM (
    'authorization', 'capture', 'refund', 'void',
    'settlement', 'chargeback', 'adjustment'
);

CREATE TYPE transaction_status AS ENUM (
    'pending', 'succeeded', 'failed', 'reversed'
);

CREATE TABLE transactions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id          UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    type                transaction_type NOT NULL,
    status              transaction_status NOT NULL DEFAULT 'pending',
    amount              BIGINT NOT NULL,
    currency            VARCHAR(3) NOT NULL,
    processor_id        VARCHAR(100),
    processor_response  JSONB,
    authorization_code  VARCHAR(50),
    failure_code        VARCHAR(50),
    failure_message     TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_txn_payment ON transactions (payment_id);
CREATE INDEX idx_txn_merchant ON transactions (merchant_id, created_at DESC);
CREATE INDEX idx_txn_processor ON transactions (processor_id);
```

## refunds
```sql
CREATE TYPE refund_status AS ENUM ('pending', 'succeeded', 'failed');

CREATE TABLE refunds (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id      UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    amount          BIGINT NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3) NOT NULL,
    status          refund_status NOT NULL DEFAULT 'pending',
    reason          VARCHAR(50) CHECK (reason IN ('customer_request', 'duplicate', 'fraudulent', 'other')),
    metadata        JSONB DEFAULT '{}',
    processor_id    VARCHAR(100),
    failure_reason  TEXT,
    idempotency_key VARCHAR(64),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refunds_payment ON refunds (payment_id);
CREATE INDEX idx_refunds_merchant ON refunds (merchant_id, created_at DESC);
CREATE UNIQUE INDEX idx_refunds_idempotency ON refunds (idempotency_key) WHERE idempotency_key IS NOT NULL;
```

## disputes
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
