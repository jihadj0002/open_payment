ALTER TABLE webhooks DROP COLUMN IF EXISTS previous_secret;
ALTER TABLE webhooks DROP COLUMN IF EXISTS previous_secret_expires_at;
