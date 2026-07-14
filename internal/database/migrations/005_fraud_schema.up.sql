CREATE TABLE IF NOT EXISTS fraud_checks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_intent_id   UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id         UUID NOT NULL REFERENCES merchants(id),
    score               INT NOT NULL DEFAULT 0,
    threshold           INT NOT NULL DEFAULT 50,
    verdict             VARCHAR(20) NOT NULL DEFAULT 'pass',
    flags               TEXT[] DEFAULT '{}',
    checked_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fraud_configs (
    merchant_id         UUID PRIMARY KEY REFERENCES merchants(id),
    max_amount          BIGINT NOT NULL DEFAULT 1000000,
    block_vpn           BOOLEAN NOT NULL DEFAULT false,
    max_ip_count        INT NOT NULL DEFAULT 10,
    enabled             BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fraud_checks_payment ON fraud_checks(payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_fraud_checks_merchant ON fraud_checks(merchant_id);
