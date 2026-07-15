ALTER TABLE webhooks ADD COLUMN IF NOT EXISTS previous_secret VARCHAR(64);
ALTER TABLE webhooks ADD COLUMN IF NOT EXISTS previous_secret_expires_at TIMESTAMPTZ;
