# Open Payment Gateway — Task Board

**Last Updated:** 2026-07-14
**Maintained by:** Orchestrator Agent

---

## Active Tasks

### Phase 0: Bootstrap ✅ (Complete)
| Task ID | Title | Status | Assignee |
|---------|-------|--------|----------|
| TASK-P0-001 | Go module + project structure | DONE | backend-engineer |
| TASK-P0-002 | Config + middleware stack | DONE | backend-engineer |
| TASK-P0-003 | Next.js frontend scaffold | DONE | frontend-engineer |
| TASK-P0-004 | Docker Compose + env | DONE | devops-engineer |
| TASK-P0-005 | Mock payment processor | DONE | backend-engineer |
| TASK-P0-006 | Seed data + migrations | DONE | backend-engineer |

### Phase 1: Core Backend Infrastructure (Current)
| Task ID | Title | Priority | Status | Assignee | Dependencies |
|---------|-------|----------|--------|----------|--------------|
| TASK-IMP-001 | Database connection layer + base repository | CRITICAL | DONE | backend-engineer | TASK-P0-001 |
| TASK-IMP-002 | Auth service (JWT, API key middleware, login/register) | CRITICAL | DONE | backend-engineer | TASK-IMP-001 |
| TASK-IMP-003 | Merchant service (CRUD, onboarding, API key management) | CRITICAL | DONE | backend-engineer | TASK-IMP-002 |
| TASK-IMP-004 | Payment service (create, process, capture, refund, void) | CRITICAL | DONE | backend-engineer | TASK-IMP-003 |
| TASK-IMP-005 | Payment state machine | CRITICAL | DONE | backend-engineer | TASK-IMP-004 |
| TASK-IMP-006 | Customer service + payment methods | HIGH | DONE | backend-engineer | TASK-IMP-004 |
| TASK-IMP-007 | Webhook system (register, deliver, retry, sign) | HIGH | DONE | backend-engineer | TASK-IMP-004 |
| TASK-IMP-008 | Balance/Ledger service (double-entry) | HIGH | DONE | backend-engineer | TASK-IMP-004 |
| TASK-IMP-009 | Fraud detection service (rule engine) | HIGH | DONE | backend-engineer | TASK-IMP-004 |
| TASK-IMP-010 | Settlement service | MEDIUM | TODO | backend-engineer | TASK-IMP-008 |

### Phase 2: Admin API
| Task ID | Title | Priority | Status | Assignee | Dependencies |
|---------|-------|----------|--------|----------|--------------|
| TASK-ADM-001 | Admin merchant management endpoints | HIGH | TODO | backend-engineer | TASK-IMP-003 |
| TASK-ADM-002 | Admin transaction/dispute endpoints | HIGH | TODO | backend-engineer | TASK-IMP-004 |
| TASK-ADM-003 | Fee configuration system | MEDIUM | TODO | backend-engineer | TASK-IMP-003 |
| TASK-ADM-004 | System config endpoints | MEDIUM | TODO | backend-engineer | — |
| TASK-ADM-005 | Audit log system | MEDIUM | TODO | backend-engineer | — |

### Phase 3: Frontend
| Task ID | Title | Priority | Status | Assignee | Dependencies |
|---------|-------|----------|--------|----------|--------------|
| TASK-FE-001 | Login/register pages with API integration | HIGH | TODO | frontend-engineer | TASK-IMP-002 |
| TASK-FE-002 | Merchant dashboard (stats, transactions table) | HIGH | TODO | frontend-engineer | TASK-IMP-004 |
| TASK-FE-003 | Payment form component | HIGH | TODO | frontend-engineer | TASK-IMP-004 |
| TASK-FE-004 | API keys management page | MEDIUM | TODO | frontend-engineer | TASK-IMP-003 |
| TASK-FE-005 | Webhook configuration page | MEDIUM | TODO | frontend-engineer | TASK-IMP-007 |

### Phase 4: Infrastructure & DevOps
| Task ID | Title | Priority | Status | Assignee | Dependencies |
|---------|-------|----------|--------|----------|--------------|
| TASK-DEV-001 | CI/CD pipeline (GitHub Actions) | HIGH | TODO | devops-engineer | TASK-P0-001 |
| TASK-DEV-002 | Terraform modules (VPC, RDS, EKS, ElastiCache) | HIGH | TODO | devops-engineer | — |
| TASK-DEV-003 | Monitoring stack (Prometheus + Grafana dashboards) | HIGH | TODO | devops-engineer | — |
| TASK-DEV-004 | K8s manifests for staging + production | HIGH | TODO | devops-engineer | — |
| TASK-DEV-005 | Log aggregation (Loki/ELK) | MEDIUM | TODO | devops-engineer | — |

### Phase 5: Testing & QA
| Task ID | Title | Priority | Status | Assignee | Dependencies |
|---------|-------|----------|--------|----------|--------------|
| TASK-QA-001 | Unit tests for all services | CRITICAL | TODO | qa-engineer | TASK-IMP-001 through 010 |
| TASK-QA-002 | Integration tests (API test suite) | HIGH | TODO | qa-engineer | TASK-IMP-001 through 010 |
| TASK-QA-003 | Load testing (k6 scripts) | MEDIUM | TODO | qa-engineer | TASK-IMP-001 through 010 |
| TASK-QA-004 | Chaos engineering experiments | LOW | TODO | qa-engineer | TASK-IMP-001 through 010 |

### Phase 6: Security & Compliance
| Task ID | Title | Priority | Status | Assignee | Dependencies |
|---------|-------|----------|--------|----------|--------------|
| TASK-SEC-001 | PCI DSS self-assessment questionnaire | CRITICAL | TODO | compliance-officer | — |
| TASK-SEC-002 | Encryption audit (data-at-rest, data-in-transit) | CRITICAL | TODO | security-engineer | — |
| TASK-SEC-003 | API security review (rate limiting, auth, CORS) | HIGH | TODO | security-engineer | — |
| TASK-SEC-004 | Dependency vulnerability scan | HIGH | TODO | security-engineer | — |

---

## Legend
| Status | Meaning |
|--------|---------|
| BACKLOG | Not yet prioritized |
| TODO | Ready to be worked on |
| IN_PROGRESS | Currently being worked on |
| REVIEW | Ready for review |
| DONE | Completed |
| BLOCKED | Cannot proceed due to dependency |

## Quick Stats
- **Total Tasks:** 37
- **IN_PROGRESS:** 1
- **TODO:** 25
- **DONE:** 11
