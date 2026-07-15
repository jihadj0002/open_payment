CREATE TABLE IF NOT EXISTS disputes (
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

CREATE INDEX IF NOT EXISTS idx_disputes_payment_intent ON disputes(payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_disputes_merchant ON disputes(merchant_id);
CREATE INDEX IF NOT EXISTS idx_disputes_status ON disputes(status);
