DROP INDEX IF EXISTS idx_payment_intents_provider_ref;
ALTER TABLE payment_intents DROP COLUMN IF EXISTS wallet_provider;
