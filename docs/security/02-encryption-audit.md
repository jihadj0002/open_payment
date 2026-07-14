# Encryption Audit

## Data-at-Rest

| Data Store | Encryption | Notes |
|------------|-----------|-------|
| PostgreSQL (RDS) | AES-256 | AWS RDS encryption at rest (KMS) |
| Redis (ElastiCache) | Encryption in transit | Redis AUTH + TLS |
| S3 (backups) | AES-256-SSE | Server-side encryption |

## Data-in-Transit

| Channel | Encryption | Notes |
|---------|-----------|-------|
| API (external) | TLS 1.3 | Terminated at ALB; certs via ACM |
| API (internal) | mTLS | Between services in K8s via Istio sidecar |
| Database | TLS 1.2 | PostgreSQL SSL enforced |
| Redis | TLS | Redis AUTH + TLS |
| Kafka | TLS + SASL | Confluent Kafka with TLS |

## Application-Level Encryption

| Secret | Method | Type |
|--------|--------|------|
| Passwords | bcrypt | One-way hashing |
| API keys | SHA-256 | One-way hashing |
| JWT | HMAC-SHA256 | Symmetric signing |
| Webhook secrets | Random hex | Stored as-is, never logged |
| Card fingerprints | SHA-256 | Deterministic for deduplication |
