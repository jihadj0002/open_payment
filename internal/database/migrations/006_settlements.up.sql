CREATE TABLE IF NOT EXISTS settlements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    amount          BIGINT NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'BDT',
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    fee             BIGINT NOT NULL DEFAULT 0,
    net_amount      BIGINT NOT NULL DEFAULT 0,
    payout_ref      VARCHAR(255),
    period_start    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    period_end      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_settlements_merchant ON settlements(merchant_id);
CREATE INDEX IF NOT EXISTS idx_settlements_status ON settlements(status);
