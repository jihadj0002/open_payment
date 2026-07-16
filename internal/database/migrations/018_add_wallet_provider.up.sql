ALTER TABLE payment_intents
    ADD COLUMN IF NOT EXISTS wallet_provider VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_payment_intents_provider_ref ON payment_intents(provider_ref);
