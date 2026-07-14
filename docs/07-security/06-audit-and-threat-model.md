# Audit and Threat Model

## Immutable Audit Logging

### Design Principles
1. **Append-only:** Logs are never modified or deleted
2. **Tamper-evident:** Hash chain links entries together
3. **Complete:** Every state-changing action is logged
4. **Searchable:** Structured JSON fields for querying

### Audit Log Schema
```sql
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID NOT NULL,
    actor_type      VARCHAR(20) NOT NULL,          -- 'merchant', 'admin', 'system'
    action          VARCHAR(100) NOT NULL,          -- 'merchant.create', 'payment.refund'
    resource_type   VARCHAR(50) NOT NULL,
    resource_id     UUID,
    details         JSONB,                         -- before/after values
    ip_address      INET,
    user_agent      TEXT,
    request_id      VARCHAR(64),
    previous_hash   VARCHAR(64),                   -- SHA-256 of previous log entry
    hash            VARCHAR(64),                   -- SHA-256 of this entry
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Hash Chain Implementation
```go
func CreateAuditLog(ctx context.Context, entry AuditLogEntry) error {
    tx, err := db.BeginTx(ctx, nil)

    // Get previous log's hash
    var prevHash string
    tx.QueryRow(ctx, `
        SELECT hash FROM audit_logs
        ORDER BY created_at DESC, id DESC
        LIMIT 1 FOR UPDATE
    `).Scan(&prevHash)

    // Calculate hash of this entry
    entry.PreviousHash = prevHash
    entry.Hash = calculateHash(entry)

    // Insert
    tx.Exec(ctx, `INSERT INTO audit_logs ...`, ...)
    return tx.Commit()
}

func calculateHash(entry AuditLogEntry) string {
    data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
        entry.ID, entry.ActorID, entry.Action, entry.ResourceType,
        entry.ResourceID, entry.Details, entry.CreatedAt, entry.PreviousHash, entry.RequestID)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

### What Gets Logged
| Action | Actor | Details Captured |
|--------|-------|------------------|
| Merchant create | Admin | Business name, email, country |
| Merchant approve/reject | Admin | Previous status → new status, reason |
| Payment create | Merchant | Amount, currency, customer, payment method |
| Payment capture | Merchant | Amount captured |
| Payment refund | Merchant | Amount, reason |
| API key create | Merchant | Key name, permissions, IP |
| API key revoke | Merchant | Key prefix |
| User login | User | IP, user agent, success/failure |
| User password change | User | — |
| MFA enable/disable | User | — |
| Fee config change | Admin | Before/after values |
| Settlement trigger | System/Admin | Batch ID, date |

## Threat Model (STRIDE)

### Assets
| Asset | Sensitivity | Location |
|-------|-------------|----------|
| Cardholder data (tokenized) | High | Database (tokens), Processor (PAN) |
| API secret keys | Critical | Database (hashed), User (one-time display) |
| Merchant credentials | High | Database (hashed) |
| JWT signing keys | Critical | KMS |
| Transaction data | High | Database |
| Ledger entries | Critical | Database |
| Bank account details | High | Database (encrypted) |
| KYC documents | High | S3 (encrypted) |

### Threats by Component

#### API Gateway
| Threat | Mitigation |
|--------|------------|
| Spoofing (fake merchant) | API key + HMAC signature verification |
| Tampering (modify request) | HMAC signature validation, TLS in transit |
| Repudiation | All requests logged in audit trail |
| Information disclosure | TLS 1.3, response sanitization |
| DoS | Rate limiting, WAF, DDoS protection |
| Elevation of privilege | RBAC middleware, permission checks |

#### Payment Service
| Threat | Mitigation |
|--------|------------|
| IDOR (access other merchant's payment) | merchant_id check on every query |
| Tampering (modify amount/status) | Idempotency + state machine validation |
| Race condition | Optimistic locking with version column |
| Replay attack | Idempotency keys, nonce validation |
| SQL injection | Parameterized queries |

#### Ledger Service
| Threat | Mitigation |
|--------|------------|
| Invalid ledger entry (unbalanced) | Constraint: sum(debits) = sum(credits) |
| Duplicate entry | Idempotency key |
| Unauthorized balance change | Service-to-service gRPC auth |
| Data loss | Append-only design, replication |

#### Database
| Threat | Mitigation |
|--------|------------|
| Unauthorized access | Network isolation, IAM auth, encrypted connections |
| Data exfiltration | Encryption at rest, access logging |
| Data loss | WAL streaming, daily backups, PITR |
| SQL injection | Parameterized queries only |

#### Infrastructure
| Threat | Mitigation |
|--------|------------|
| Container escape | Run as non-root, read-only FS, seccomp |
| Network sniffing | mTLS between all services |
| Secrets exposure | Vault/Secrets Manager, never in code |
| Compromised image | Vulnerability scanning, signed images |

## Incident Response Plan

### Severity Levels
| Level | Definition | Response Time | Example |
|-------|------------|---------------|---------|
| SEV-1 | Service down, revenue impact | 5 min | Payment API 5xx > 1% |
| SEV-2 | Partial degradation | 15 min | Dashboard slow, webhook delay |
| SEV-3 | Minor issue, non-revenue | 1 hour | Stale data in reports |
| SEV-4 | Non-critical | 24 hours | UI bug, doc typo |

### Response Playbook
```
1. DETECT
   - Alert from monitoring (PagerDuty)
   - User report
   - Automated anomaly detection

2. TRIAGE (5 min)
   - Confirm severity
   - Notify on-call engineer
   - Create incident channel (#incident-xxx)

3. MITIGATE (15 min)
   - Rollback if recent deploy
   - Scale up if capacity issue
   - Block malicious traffic if attack
   - Failover if region issue

4. RESOLVE
   - Apply fix
   - Verify fix
   - Close incident

5. POST-MORTEM (72 hours)
   - Root cause analysis
   - Timeline
   - Action items
   - Share with team
```

### On-Call Rotation
- **Primary:** 1 week rotation (8am–8pm on-call)
- **Secondary:** Always available for SEV-1 escalation
- **Escalation:** Primary → Secondary → Engineering Manager
- **Handoff:** Weekly handoff document + recorded status

### Communication
- **Internal:** Slack #incident channel
- **Status page:** StatusPage.io for merchant-facing status
- **Post-mortem:** Shared in team meeting + wiki
