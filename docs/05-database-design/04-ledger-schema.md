# Ledger Schema

## Core Principle: Double-Entry Accounting

Every financial movement creates **two ledger entries**: a debit from one account and a credit to another. The sum of all debits must always equal the sum of all credits. No account balances are directly updated — they are always calculated as `SUM(credits) - SUM(debits)`.

## ledger_accounts

Each merchant has a fixed set of ledger accounts:

| Account Type | Who Owns | Description |
|-------------|----------|-------------|
| `merchant_receivable` | Merchant | Money owed to merchant (pending settlement) |
| `merchant_available` | Merchant | Settled money available for payout |
| `merchant_reserve` | Merchant | Held for chargebacks, disputes |
| `gateway_fee_receivable` | Gateway | Fees earned but not yet recognized |
| `gateway_fee_revenue` | Gateway | Recognized fee revenue |
| `tax_payable` | Government | VAT/tax collected, to be remitted |
| `processor_payable` | Gateway | Money owed to processor |
| `settlement_clearing` | Gateway | Temporary account during settlement |

```sql
CREATE TYPE account_type AS ENUM (
    'merchant_receivable',
    'merchant_available',
    'merchant_reserve',
    'gateway_fee_receivable',
    'gateway_fee_revenue',
    'tax_payable',
    'processor_payable',
    'settlement_clearing'
);

CREATE TABLE ledger_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    type            account_type NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'BDT',
    version         INTEGER NOT NULL DEFAULT 1,  -- optimistic locking
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (merchant_id, type, currency)
);

CREATE INDEX idx_la_merchant ON ledger_accounts (merchant_id);
```

## ledger_entries

Every entry is immutable. Entries are never updated or deleted — corrections use reversing entries.

```sql
CREATE TYPE entry_direction AS ENUM ('debit', 'credit');
CREATE TYPE entry_type AS ENUM (
    'payment_capture',
    'payment_refund',
    'fee_deduction',
    'tax_deduction',
    'settlement_transfer',
    'payout_initiation',
    'chargeback_debit',
    'chargeback_credit',
    'adjustment',
    'reserve_hold',
    'reserve_release'
);

CREATE TABLE ledger_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL REFERENCES ledger_accounts(id),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    transaction_id  UUID NOT NULL,                    -- references external transaction
    entry_type      entry_type NOT NULL,
    direction       entry_direction NOT NULL,
    amount          BIGINT NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3) NOT NULL,
    balance_before  BIGINT NOT NULL,
    balance_after   BIGINT NOT NULL,
    description     TEXT,
    reference_type  VARCHAR(50),                      -- 'payment_intent', 'refund', 'chargeback'
    reference_id    UUID,                             -- ID of the referenced entity
    idempotency_key VARCHAR(64) UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY HASH (account_id);

CREATE INDEX idx_le_account ON ledger_entries (account_id, created_at DESC);
CREATE INDEX idx_le_merchant ON ledger_entries (merchant_id, created_at DESC);
CREATE INDEX idx_le_transaction ON ledger_entries (transaction_id);
CREATE INDEX idx_le_reference ON ledger_entries (reference_type, reference_id);
CREATE UNIQUE INDEX idx_le_idempotency ON ledger_entries (idempotency_key) WHERE idempotency_key IS NOT NULL;
```

## balances (materialized view)

```sql
CREATE TABLE balances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    account_type    account_type NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'BDT',
    balance         BIGINT NOT NULL DEFAULT 0,
    calculated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (merchant_id, account_type, currency)
);

-- Refresh function
CREATE OR REPLACE FUNCTION refresh_balances(p_merchant_id UUID)
RETURNS VOID AS $$
BEGIN
    DELETE FROM balances WHERE merchant_id = p_merchant_id;

    INSERT INTO balances (merchant_id, account_type, currency, balance, calculated_at)
    SELECT
        la.merchant_id,
        la.type,
        la.currency,
        COALESCE(SUM(
            CASE WHEN le.direction = 'credit' THEN le.amount
                 WHEN le.direction = 'debit' THEN -le.amount
            END
        ), 0) AS balance,
        NOW()
    FROM ledger_accounts la
    LEFT JOIN ledger_entries le ON le.account_id = la.id
    WHERE la.merchant_id = p_merchant_id
    GROUP BY la.merchant_id, la.type, la.currency;
END;
$$ LANGUAGE plpgsql;
```

## Example: Payment Capture ($10.00)

| Account | Direction | Amount | Description |
|---------|-----------|--------|-------------|
| merchant_receivable | Debit | 1000 | 1000 BDT owed to merchant |
| gateway_fee_receivable | Debit | 25 | 2.5% fee |
| tax_payable | Credit | 5 | 0.5% VAT |
| processor_payable | Credit | 970 | Net to merchant after fees |
| settlement_clearing | Credit | 1000 | Total captured amount |

## Example: Settlement Transfer

| Account | Direction | Amount | Description |
|---------|-----------|--------|-------------|
| merchant_receivable | Credit | 1000 | Clear receivable |
| merchant_available | Debit | 975 | Available for payout (after fees) |
| gateway_fee_receivable | Credit | 25 | Recognize fee |
| gateway_fee_revenue | Debit | 25 | Fee revenue recognized |
| processor_payable | Credit | 970 | Pay processor |
| settlement_clearing | Debit | 1000 | Close settlement |

## Example: Refund ($5.00)

| Account | Direction | Amount | Description |
|---------|-----------|--------|-------------|
| merchant_available | Credit | 500 | Reduce available balance |
| merchant_receivable | Debit | 500 | Reversal of receivable |
| gateway_fee_receivable | Credit | 12 | Reverse fee receivable |
| processor_payable | Debit | 488 | Reversal of processor amount |

## key_validation

```sql
-- Validate that all ledger entries balance to zero
SELECT SUM(CASE WHEN direction = 'debit' THEN amount ELSE -amount END) AS net
FROM ledger_entries
WHERE transaction_id = $1;
-- Must return 0 for a balanced transaction.

-- Validate merchant balance
SELECT
    COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount ELSE 0 END), 0) -
    COALESCE(SUM(CASE WHEN direction = 'debit' THEN amount ELSE 0 END), 0) AS calculated_balance
FROM ledger_entries le
JOIN ledger_accounts la ON la.id = le.account_id
WHERE la.merchant_id = $1 AND la.type = 'merchant_available';
```
