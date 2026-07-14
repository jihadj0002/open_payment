CREATE TABLE IF NOT EXISTS fee_configs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    rate                NUMERIC(5,2) NOT NULL DEFAULT 2.2,
    fixed_fee           BIGINT NOT NULL DEFAULT 0,
    cross_border_rate   NUMERIC(5,2),
    monthly_fee         BIGINT NOT NULL DEFAULT 0,
    min_monthly_volume  BIGINT NOT NULL DEFAULT 0,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS system_config (
    id                     VARCHAR(20) PRIMARY KEY DEFAULT 'default',
    supported_currencies   JSONB NOT NULL DEFAULT '["BDT","USD"]',
    supported_countries    JSONB NOT NULL DEFAULT '["BD"]',
    max_transaction_amount BIGINT NOT NULL DEFAULT 10000000,
    maintenance_mode       BOOLEAN NOT NULL DEFAULT false,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO system_config (id) VALUES ('default') ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID NOT NULL,
    action          VARCHAR(50) NOT NULL,
    resource_type   VARCHAR(50) NOT NULL,
    resource_id     VARCHAR(100),
    details         TEXT,
    ip_address      VARCHAR(45),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_actor ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);
