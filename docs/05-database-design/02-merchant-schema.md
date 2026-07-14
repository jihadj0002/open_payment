# Merchant Schema

## merchants
```sql
CREATE TABLE merchants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_name   VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    phone           VARCHAR(20) NOT NULL,
    website         VARCHAR(500),
    description     TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'active', 'suspended', 'terminated')),
    verification_status VARCHAR(20) NOT NULL DEFAULT 'unverified'
                    CHECK (verification_status IN ('unverified', 'pending', 'verified', 'rejected')),
    country         VARCHAR(2) NOT NULL DEFAULT 'BD',
    currency        VARCHAR(3) NOT NULL DEFAULT 'BDT',
    timezone        VARCHAR(50) NOT NULL DEFAULT 'Asia/Dhaka',
    kyc_documents   JSONB,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_merchants_status ON merchants (status);
CREATE INDEX idx_merchants_email ON merchants (email);
CREATE INDEX idx_merchants_created_at ON merchants (created_at DESC);
```

## merchant_users
```sql
CREATE TABLE merchant_users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    email           VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'member'
                    CHECK (role IN ('owner', 'admin', 'developer', 'analyst', 'member')),
    mfa_enabled     BOOLEAN NOT NULL DEFAULT false,
    mfa_secret      VARCHAR(100),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (merchant_id, email)
);

CREATE INDEX idx_merchant_users_merchant ON merchant_users (merchant_id);
```

## merchant_settings
```sql
CREATE TABLE merchant_settings (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id             UUID NOT NULL UNIQUE REFERENCES merchants(id) ON DELETE CASCADE,
    auto_capture            BOOLEAN NOT NULL DEFAULT true,
    statement_descriptor    VARCHAR(22),
    send_email_receipts     BOOLEAN NOT NULL DEFAULT false,
    receipt_email_from      VARCHAR(255),
    ip_whitelist            TEXT[],
    webhook_url             VARCHAR(500),
    webhook_secret          VARCHAR(64),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## api_keys
```sql
CREATE TYPE api_key_mode AS ENUM ('live', 'test');

CREATE TABLE api_keys (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    key_prefix      VARCHAR(8) NOT NULL,
    key_hash        VARCHAR(64) NOT NULL,
    key_last4       VARCHAR(4) NOT NULL,
    mode            api_key_mode NOT NULL DEFAULT 'live',
    permissions     TEXT[] NOT NULL DEFAULT '{}',
    allowed_ips     TEXT[],
    expires_at      TIMESTAMPTZ,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    last_used_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_api_keys_hash ON api_keys (key_hash);
CREATE INDEX idx_api_keys_merchant ON api_keys (merchant_id);
```

## webhook_configs
```sql
CREATE TABLE webhook_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    url             VARCHAR(500) NOT NULL,
    events          TEXT[] NOT NULL DEFAULT '{}',
    secret          VARCHAR(64) NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    description     VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (merchant_id, url)
);

CREATE INDEX idx_webhook_configs_merchant ON webhook_configs (merchant_id);
```

## fee_configs
```sql
CREATE TABLE fee_configs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id         UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name                VARCHAR(100) NOT NULL DEFAULT 'Default',
    rate_bps            INTEGER NOT NULL DEFAULT 250,       -- 250 = 2.5%
    fixed_fee           BIGINT NOT NULL DEFAULT 300,         -- in smallest currency unit
    cross_border_rate_bps INTEGER DEFAULT 300,              -- 300 = 3.0%
    monthly_fee         BIGINT NOT NULL DEFAULT 0,
    min_monthly_volume  BIGINT DEFAULT 0,
    is_active           BOOLEAN NOT NULL DEFAULT true,
    effective_from      DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to        DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fee_configs_merchant ON fee_configs (merchant_id);
```

## kyc_documents
```sql
CREATE TABLE kyc_documents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    document_type   VARCHAR(50) NOT NULL,
    s3_key          VARCHAR(500) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'verified', 'rejected')),
    rejection_reason TEXT,
    uploaded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    verified_at     TIMESTAMPTZ,
    verified_by     UUID REFERENCES admin_users(id)
);

CREATE INDEX idx_kyc_merchant ON kyc_documents (merchant_id);
```
