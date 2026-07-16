# Customer Schema

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/database/migrations/001_initial_schema.up.sql`, `004_add_payment_methods.up.sql`, `014_saved_payment_methods.up.sql`

## customers
```sql
CREATE TABLE customers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    email           VARCHAR(255),
    name            VARCHAR(255),
    phone           VARCHAR(20),
    description     TEXT,
    metadata        JSONB DEFAULT '{}',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (merchant_id, email)
);

CREATE INDEX idx_customers_merchant ON customers (merchant_id);
CREATE INDEX idx_customers_email ON customers (merchant_id, email);
CREATE INDEX idx_customers_phone ON customers (merchant_id, phone);
```

## saved_cards (tokenized card data)
```sql
CREATE TABLE saved_cards (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id     UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    token           VARCHAR(100) NOT NULL,        -- processor/vault token
    last4           VARCHAR(4) NOT NULL,
    brand           VARCHAR(20) NOT NULL,
    exp_month       INTEGER NOT NULL,
    exp_year        INTEGER NOT NULL,
    cardholder_name VARCHAR(255),
    fingerprint     VARCHAR(64),                  -- for deduplication
    is_default      BOOLEAN NOT NULL DEFAULT false,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sc_customer ON saved_cards (customer_id);
CREATE INDEX idx_sc_merchant ON saved_cards (merchant_id);
CREATE INDEX idx_sc_fingerprint ON saved_cards (fingerprint);
```

## tokens (one-time use tokens for payment)
```sql
CREATE TYPE token_type AS ENUM ('card', 'wallet', 'bank_account');

CREATE TABLE tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    type            token_type NOT NULL,
    -- Card token data
    card_last4      VARCHAR(4),
    card_brand      VARCHAR(20),
    card_exp_month  INTEGER,
    card_exp_year   INTEGER,
    -- Wallet token data
    wallet_type     VARCHAR(20),
    wallet_phone    VARCHAR(20),
    -- Processor token (created after tokenization)
    processor_token VARCHAR(100),
    -- Metadata
    is_used         BOOLEAN NOT NULL DEFAULT false,
    expires_at      TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '24 hours'),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tokens_merchant ON tokens (merchant_id);
CREATE INDEX idx_tokens_unused ON tokens (merchant_id) WHERE is_used = false;
```

## device_fingerprints
```sql
CREATE TABLE device_fingerprints (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fingerprint_id  VARCHAR(100) NOT NULL UNIQUE,
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    customer_id     UUID REFERENCES customers(id),
    data            JSONB,                    -- browser info, screen res, fonts, etc.
    ip_address      INET,
    user_agent      TEXT,
    risk_score      REAL DEFAULT 0.0,
    first_seen_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_df_fingerprint ON device_fingerprints (fingerprint_id);
CREATE INDEX idx_df_customer ON device_fingerprints (customer_id);
```

## risk_scores
```sql
CREATE TABLE risk_scores (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id      UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    score           REAL NOT NULL,                     -- 0.0 to 1.0
    action          VARCHAR(20) NOT NULL                -- 'allow', 'review', 'block'
                    CHECK (action IN ('allow', 'review', 'block')),
    flags           TEXT[] DEFAULT '{}',
    rules_triggered TEXT[] DEFAULT '{}',
    ml_model_name   VARCHAR(50),
    ml_model_version VARCHAR(20),
    assessment_id   VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rs_payment ON risk_scores (payment_id);
CREATE INDEX idx_rs_merchant ON risk_scores (merchant_id, created_at DESC);
```
