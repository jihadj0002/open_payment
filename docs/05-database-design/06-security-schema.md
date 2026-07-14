# Security Schema

## admin_users
```sql
CREATE TABLE admin_users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'operator'
                    CHECK (role IN ('super_admin', 'operator', 'compliance', 'support', 'readonly')),
    mfa_enabled     BOOLEAN NOT NULL DEFAULT true,
    mfa_secret      VARCHAR(100),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## roles
```sql
CREATE TABLE roles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(50) NOT NULL UNIQUE,
    description     TEXT,
    is_system       BOOLEAN NOT NULL DEFAULT false,    -- system roles can't be deleted
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed roles
INSERT INTO roles (name, description, is_system) VALUES
    ('super_admin', 'Full system access', true),
    ('operator', 'Day-to-day operations', true),
    ('compliance', 'KYC/AML reviews', true),
    ('support', 'Merchant support', true),
    ('readonly', 'View-only access', true);
```

## permissions
```sql
CREATE TABLE permissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource        VARCHAR(50) NOT NULL,     -- 'merchant', 'payment', 'settlement', etc.
    action          VARCHAR(50) NOT NULL,     -- 'create', 'read', 'update', 'delete', 'approve'
    description     TEXT,
    UNIQUE (resource, action)
);
```

## role_permissions
```sql
CREATE TABLE role_permissions (
    role_id         UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id   UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);
```

## sessions
```sql
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,                    -- references merchant_users or admin_users
    user_type       VARCHAR(20) NOT NULL              -- 'merchant' or 'admin'
                    CHECK (user_type IN ('merchant', 'admin')),
    refresh_token   VARCHAR(255) NOT NULL UNIQUE,
    access_token_jti VARCHAR(64) NOT NULL UNIQUE,
    ip_address      INET,
    user_agent      TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON sessions (user_id);
CREATE INDEX idx_sessions_refresh ON sessions (refresh_token);
CREATE INDEX idx_sessions_active ON sessions (user_id) WHERE is_active = true;
```

## audit_logs (immutable)
```sql
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID NOT NULL,
    actor_type      VARCHAR(20) NOT NULL              -- 'merchant', 'admin', 'system'
                    CHECK (actor_type IN ('merchant', 'admin', 'system')),
    action          VARCHAR(100) NOT NULL,            -- 'merchant.create', 'payment.refund', etc.
    resource_type   VARCHAR(50) NOT NULL,
    resource_id     UUID,
    details         JSONB,                            -- before/after values, additional context
    ip_address      INET,
    user_agent      TEXT,
    request_id      VARCHAR(64),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_al_actor ON audit_logs (actor_id, created_at DESC);
CREATE INDEX idx_al_action ON audit_logs (action, created_at DESC);
CREATE INDEX idx_al_resource ON audit_logs (resource_type, resource_id);
CREATE INDEX idx_al_created_at ON audit_logs (created_at DESC);
```

## security_logs
```sql
CREATE TYPE security_event_type AS ENUM (
    'login_success', 'login_failure', 'mfa_success', 'mfa_failure',
    'password_change', 'api_key_created', 'api_key_revoked',
    'rate_limit_exceeded', 'suspicious_ip', 'brute_force_detected',
    'permission_denied', 'session_expired', 'admin_action'
);

CREATE TABLE security_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type      security_event_type NOT NULL,
    actor_id        UUID,
    actor_type      VARCHAR(20),
    ip_address      INET,
    user_agent      TEXT,
    resource        VARCHAR(100),
    details         JSONB,
    severity        VARCHAR(10) NOT NULL DEFAULT 'info'
                    CHECK (severity IN ('info', 'warning', 'critical')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_sl_event ON security_logs (event_type, created_at DESC);
CREATE INDEX idx_sl_severity ON security_logs (severity, created_at DESC);
CREATE INDEX idx_sl_actor ON security_logs (actor_id, created_at DESC);
```

## login_history
```sql
CREATE TABLE login_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    user_type       VARCHAR(20) NOT NULL,
    ip_address      INET NOT NULL,
    user_agent      TEXT,
    success         BOOLEAN NOT NULL,
    failure_reason  VARCHAR(100),
    mfa_used        BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lh_user ON login_history (user_id, created_at DESC);
CREATE INDEX idx_lh_ip ON login_history (ip_address);
```

## api_key_history
```sql
CREATE TABLE api_key_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID NOT NULL REFERENCES merchants(id),
    key_prefix      VARCHAR(8) NOT NULL,
    action          VARCHAR(20) NOT NULL              -- 'created', 'revoked', 'rotated'
                    CHECK (action IN ('created', 'revoked', 'rotated')),
    actor_id        UUID NOT NULL,
    ip_address      INET,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_akh_merchant ON api_key_history (merchant_id, created_at DESC);
```

## fraud_rules
```sql
CREATE TABLE fraud_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID REFERENCES merchants(id),    -- NULL = global rules
    name            VARCHAR(100) NOT NULL,
    rule_type       VARCHAR(50) NOT NULL,              -- 'velocity', 'ip_country', 'amount', 'bin', etc.
    conditions      JSONB NOT NULL,                    -- rule conditions
    action          VARCHAR(20) NOT NULL               -- 'block', 'review', 'allow'
                    CHECK (action IN ('block', 'review', 'allow')),
    priority        INTEGER NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fr_merchant ON fraud_rules (merchant_id);
CREATE INDEX idx_fr_active ON fraud_rules (is_active);
```
