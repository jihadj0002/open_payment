ALTER TABLE settlements DROP COLUMN IF EXISTS currency;
ALTER TABLE settlements DROP COLUMN IF EXISTS fee_amount;
ALTER TABLE settlements DROP COLUMN IF EXISTS net_amount;
ALTER TABLE settlements DROP COLUMN IF EXISTS transaction_count;
ALTER TABLE settlements DROP COLUMN IF EXISTS metadata;
