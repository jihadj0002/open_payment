CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS merchants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    secret_key      VARCHAR(64) NOT NULL UNIQUE,
    public_key      VARCHAR(64) NOT NULL UNIQUE,
    webhook_url     TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    email           VARCHAR(255),
    phone           VARCHAR(20),
    name            VARCHAR(255),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payment_intents (
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

CREATE TABLE IF NOT EXISTS transactions (
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

CREATE TABLE IF NOT EXISTS ledger_entries (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id      UUID NOT NULL REFERENCES transactions(id),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    entry_type          VARCHAR(20) NOT NULL,
    amount              BIGINT NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'BDT',
    balance_before      BIGINT NOT NULL DEFAULT 0,
    balance_after       BIGINT NOT NULL DEFAULT 0,
    description         TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhooks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    event               VARCHAR(50) NOT NULL,
    url                 TEXT NOT NULL,
    secret              VARCHAR(64) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhook_deliveries (
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

CREATE TABLE IF NOT EXISTS api_keys (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    key_prefix          VARCHAR(8) NOT NULL,
    key_hash            VARCHAR(64) NOT NULL UNIQUE,
    name                VARCHAR(255) NOT NULL,
    permissions         JSONB NOT NULL DEFAULT '["read"]',
    last_used_at        TIMESTAMPTZ,
    expires_at          TIMESTAMPTZ,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_intents_merchant ON payment_intents(merchant_id);
CREATE INDEX IF NOT EXISTS idx_payment_intents_status ON payment_intents(status);
CREATE INDEX IF NOT EXISTS idx_payment_intents_idempotency ON payment_intents(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_transactions_payment_intent ON transactions(payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_transactions_merchant ON transactions(merchant_id);
CREATE INDEX IF NOT EXISTS idx_ledger_merchant ON ledger_entries(merchant_id);
CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);

-- Seed data
INSERT INTO merchants (id, name, email, secret_key, public_key, webhook_url, status)
VALUES
    ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'Demo Merchant', 'merchant@demo.com',
     'sk_test_demo_secret_key_1234567890', 'pk_test_demo_public_key_1234567890',
     'https://webhook.demo.com/hooks', 'active'),
    ('b2c3d4e5-f6a7-8901-bcde-f12345678901', 'Test Shop', 'admin@testshop.com',
     'sk_test_testshop_secret_abcdef123456', 'pk_test_testshop_public_abcdef123456',
     'https://testshop.com/webhook', 'active')
ON CONFLICT (email) DO NOTHING;

INSERT INTO customers (merchant_id, email, phone, name)
SELECT m.id, 'customer1@example.com', '+8801712345678', 'Rahim Ahmed'
FROM merchants m WHERE m.email = 'merchant@demo.com'
AND NOT EXISTS (SELECT 1 FROM customers WHERE email = 'customer1@example.com');

INSERT INTO customers (merchant_id, email, phone, name)
SELECT m.id, 'customer2@example.com', '+8801812345678', 'Karim Hasan'
FROM merchants m WHERE m.email = 'merchant@demo.com'
AND NOT EXISTS (SELECT 1 FROM customers WHERE email = 'customer2@example.com');

INSERT INTO api_keys (merchant_id, key_prefix, key_hash, name, permissions)
SELECT m.id, 'sk_test_', 'hashed_test_key_1', 'Test Key', '["read","write"]'
FROM merchants m WHERE m.email = 'merchant@demo.com'
AND NOT EXISTS (SELECT 1 FROM api_keys WHERE key_hash = 'hashed_test_key_1');
