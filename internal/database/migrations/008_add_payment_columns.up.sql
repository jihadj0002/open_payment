ALTER TABLE payment_intents
    ADD COLUMN IF NOT EXISTS amount_capturable BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS amount_received BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS capture_method VARCHAR(20) NOT NULL DEFAULT 'automatic';

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(64);

CREATE INDEX IF NOT EXISTS idx_transactions_idempotency ON transactions(idempotency_key);
