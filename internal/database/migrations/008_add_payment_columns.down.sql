ALTER TABLE payment_intents
    DROP COLUMN IF EXISTS amount_capturable,
    DROP COLUMN IF EXISTS amount_received,
    DROP COLUMN IF EXISTS capture_method;

ALTER TABLE transactions
    DROP COLUMN IF EXISTS idempotency_key;
