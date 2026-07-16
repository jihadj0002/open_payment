ALTER TABLE payment_intents
    ADD COLUMN IF NOT EXISTS wallet_provider VARCHAR(50),
    ADD COLUMN IF NOT EXISTS provider_ref VARCHAR(255),
    ADD COLUMN IF NOT EXISTS redirect_url TEXT,
    ADD COLUMN IF NOT EXISTS amount_capturable BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS amount_received BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS capture_method VARCHAR(20) NOT NULL DEFAULT 'automatic';

CREATE INDEX IF NOT EXISTS idx_payment_intents_provider_ref ON payment_intents(provider_ref);
