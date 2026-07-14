# Risk Register

## Risk Assessment Matrix

| Likelihood | Impact | Rating |
|------------|--------|--------|
| Very Likely (5) | Critical (5) | 25 (Extreme) |
| Likely (4) | High (4) | 16-20 (High) |
| Possible (3) | Medium (3) | 9-15 (Medium) |
| Unlikely (2) | Low (2) | 4-8 (Low) |
| Rare (1) | Minimal (1) | 1-3 (Very Low) |

## Identified Risks

### Technical Risks

| ID | Risk | Likelihood | Impact | Rating | Mitigation |
|----|------|------------|--------|--------|------------|
| R001 | Processor API downtime causes payment failures | 3 | 5 | 15 | Circuit breaker, fallback processor, retry queue |
| R002 | Data loss due to DB corruption | 2 | 5 | 10 | WAL streaming, daily backups, PITR capability |
| R003 | Kafka cluster failure | 2 | 4 | 8 | Multi-AZ brokers, mirror maker, DLQ |
| R004 | Redis cache poisoning | 2 | 3 | 6 | TTL limits, input validation, monitoring |
| R005 | Idempotency bug leads to duplicate charges | 2 | 5 | 10 | Rigorous testing, DB unique constraints, reconciliation |
| R006 | Payment state machine allows invalid transition | 2 | 4 | 8 | 100% test coverage, state validation in DB |
| R007 | SSL/TLS certificate expiry | 3 | 4 | 12 | Auto-renewal via cert-manager, 30-day alert |
| R008 | Secrets leak in code/logs | 2 | 5 | 10 | GitLeaks pre-commit, secret scanning in CI, no secrets in logs |
| R009 | Database connection pool exhaustion | 3 | 4 | 12 | PgBouncer, pool sizing, monitoring |
| R010 | Slow query brings down DB | 3 | 4 | 12 | Statement timeout, slow query logging, query review |

### Security Risks

| ID | Risk | Likelihood | Impact | Rating | Mitigation |
|----|------|------------|--------|--------|------------|
| S001 | SQL injection vulnerability | 2 | 5 | 10 | Parameterized queries only, SAST scanning |
| S002 | Cross-tenant data access (IDOR) | 2 | 5 | 10 | merchant_id scope on all queries, integration tests |
| S003 | API key compromise | 3 | 4 | 12 | Key rotation, IP whitelisting, usage monitoring |
| S004 | DDoS attack on payment API | 3 | 4 | 12 | Cloudflare WAF, auto-scaling, rate limiting |
| S005 | Carding attack (stolen cards testing) | 4 | 3 | 12 | Velocity checks, IP reputation, CAPTCHA |
| S006 | Insider threat (admin access abuse) | 1 | 5 | 5 | Immutable audit logs, MFA, least privilege |
| S007 | PCI DSS non-compliance | 2 | 5 | 10 | Quarterly scans, annual audit, tokenization |
| S008 | Session hijacking | 2 | 4 | 8 | HTTPS-only, secure cookies, short TTL |

### Project Risks

| ID | Risk | Likelihood | Impact | Rating | Mitigation |
|----|------|------------|--------|--------|------------|
| P001 | Key engineer leaves mid-project | 2 | 4 | 8 | Code documentation, knowledge sharing, pair programming |
| P002 | Scope creep delays launch | 4 | 3 | 12 | Strict MVP scope, change control board |
| P003 | Underestimating compliance effort | 3 | 4 | 12 | Compliance expert review early, buffer in schedule |
| P004 | Third-party dependency becomes unsupported | 2 | 3 | 6 | Pin versions, evaluate alternatives |
| P005 | Cloud cost overrun | 3 | 2 | 6 | Budget alerts, right-sizing, reserved instances |
| P006 | Team communication breakdown | 2 | 3 | 6 | Daily standups, weekly syncs, written decisions |
| P007 | Regulatory change mid-project | 1 | 4 | 4 | Monitor regulatory updates, flexible architecture |

### Operational Risks

| ID | Risk | Likelihood | Impact | Rating | Mitigation |
|----|------|------------|--------|--------|------------|
| O001 | AWS region outage | 2 | 5 | 10 | Multi-AZ, backup region strategy, DR plan |
| O002 | Deploy causes production outage | 3 | 4 | 12 | Canary deploy, rollback script, feature flags |
| O003 | Monitoring alerts fatigue | 4 | 2 | 8 | Well-tuned alerts, escalation policies, noise reduction |
| O004 | Insufficient on-call coverage | 2 | 3 | 6 | Rotation schedule, backup on-call, escalation path |
| O005 | Merchant integration bugs discovered post-launch | 3 | 3 | 9 | Sandbox environment, thorough documentation, support team |

## Top 10 Risks (Priority Order)

| Rank | Risk | Score | Owner | Review Frequency |
|------|------|-------|-------|------------------|
| 1 | R001 — Processor API downtime | 15 | Backend Lead | Monthly |
| 2 | R001 — Payment state machine bug | 10 | Backend Lead | Per release |
| 3 | S007 — PCI non-compliance | 10 | Security Lead | Quarterly |
| 4 | R002 — Data loss | 10 | DevOps Lead | Monthly |
| 5 | S003 — API key compromise | 12 | Security Lead | Monthly |
| 6 | P002 — Scope creep | 12 | Product Manager | Per sprint |
| 7 | P003 — Compliance under-estimate | 12 | Product Manager | Monthly |
| 8 | R007 — TLS cert expiry | 12 | DevOps Lead | Weekly auto-check |
| 9 | R009 — Connection pool exhaustion | 12 | DevOps Lead | Monthly |
| 10 | R010 — Slow queries | 12 | Backend Lead | Per release |

## Risk Review Cadence
- **Daily:** Automated monitoring alerts (technical risks)
- **Weekly:** Team standup (project risks, blockers)
- **Monthly:** Risk register review (update likelihood/impact)
- **Quarterly:** Full risk assessment (add new risks, close resolved)
- **Per incident:** Post-mortem identifies new risks
