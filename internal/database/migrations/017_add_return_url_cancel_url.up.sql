ALTER TABLE payment_intents
    ADD COLUMN IF NOT EXISTS return_url TEXT,
    ADD COLUMN IF NOT EXISTS cancel_url TEXT,
    ADD COLUMN IF NOT EXISTS client_secret VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_payment_intents_client_secret ON payment_intents(client_secret);
