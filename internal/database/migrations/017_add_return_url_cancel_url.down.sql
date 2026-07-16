DROP INDEX IF EXISTS idx_payment_intents_client_secret;

ALTER TABLE payment_intents
    DROP COLUMN IF EXISTS client_secret,
    DROP COLUMN IF EXISTS cancel_url,
    DROP COLUMN IF EXISTS return_url;
