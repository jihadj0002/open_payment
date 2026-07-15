-- This migration notes that merchant secret_key and public_key should be encrypted
-- at the application level. The ENCRYPTION_KEY env var must be set before encrypting.
-- The actual encryption is handled by the application layer during registration.
-- For existing records, keys are encrypted on next read/write cycle.
-- This is a NO-OP SQL migration; the application handles encryption.
SELECT 1;
