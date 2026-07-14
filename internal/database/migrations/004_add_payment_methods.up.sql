CREATE TABLE IF NOT EXISTS payment_methods (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    customer_id     UUID NOT NULL REFERENCES customers(id),
    type            VARCHAR(20) NOT NULL,
    token           VARCHAR(255),
    last4           VARCHAR(4),
    brand           VARCHAR(20),
    exp_month       INT,
    exp_year        INT,
    cardholder_name VARCHAR(255),
    fingerprint     VARCHAR(64),
    wallet_type     VARCHAR(50),
    wallet_phone    VARCHAR(20),
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_methods_customer ON payment_methods(customer_id);
CREATE INDEX idx_payment_methods_merchant ON payment_methods(merchant_id);
CREATE INDEX idx_payment_methods_fingerprint ON payment_methods(fingerprint);
