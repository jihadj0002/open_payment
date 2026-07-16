CREATE TABLE IF NOT EXISTS saved_payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('card', 'wallet')),
    last4 VARCHAR(4) NOT NULL,
    brand VARCHAR(50) NOT NULL DEFAULT '',
    exp_month INTEGER,
    exp_year INTEGER,
    is_default BOOLEAN NOT NULL DEFAULT false,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(customer_id, token)
);

CREATE INDEX idx_saved_payment_methods_customer ON saved_payment_methods(customer_id);
CREATE INDEX idx_saved_payment_methods_merchant ON saved_payment_methods(merchant_id);
