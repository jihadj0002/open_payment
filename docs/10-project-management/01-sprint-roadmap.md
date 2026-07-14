# Sprint Roadmap

## Phase Overview

| Phase | Duration | Weeks | Team Size | Outcome |
|-------|----------|-------|-----------|---------|
| P1: Planning | 2 weeks | 1-2 | 2 | SRS, ADRs, ERD, Tech stack finalized |
| P2: Design | 2 weeks | 3-4 | 2 | API spec, UI mockups, DB schema, Infra design |
| P3: Infrastructure | 3 weeks | 5-7 | 2 | AWS setup, K8s cluster, CI/CD, Monitoring |
| P4: Core Backend | 8 weeks | 8-15 | 4 | Payment flow, Auth, Merchant, Ledger, Webhooks |
| P5: Fraud + Admin | 3 weeks | 16-18 | 3 | Fraud engine, Admin dashboard, Reporting |
| P6: Frontend | 4 weeks | 19-22 | 2 | Merchant dashboard, Dev portal, Admin UI |
| P7: Testing + Security | 3 weeks | 23-25 | 4 | Load tests, Chaos tests, Pen test, PCI scan |
| P8: Launch | 2 weeks | 26-27 | 4 | Pilot merchants, Monitoring, Documentation |
| P9: Post-Launch | Ongoing | 28+ | 4 | Bug fixes, Performance tuning, Feature additions |

## Detailed Sprint Plan

### P1: Planning (Weeks 1-2)

**Sprint 1 — Requirements (Week 1)**
| Task | Owner | Hours |
|------|-------|-------|
| Define business requirements | Product | 16 |
| Create user stories for all personas | Product | 12 |
| Research compliance requirements (PCI, PSD2, GDPR) | Security | 20 |
| Define non-functional requirements (SLA, latency, throughput) | Tech Lead | 8 |
| Document scope and limitations | Tech Lead | 4 |

**Sprint 2 — Architecture Design (Week 2)**
| Task | Owner | Hours |
|------|-------|-------|
| Create C4 architecture diagrams | Tech Lead | 12 |
| Define microservices boundaries | Tech Lead | 8 |
| Design database ER diagrams | Backend | 12 |
| Design API contracts (OpenAPI spec) | Backend | 16 |
| Security threat modeling (STRIDE) | Security | 8 |
| Initial infrastructure design | DevOps | 8 |

### P2: Design (Weeks 3-4)

**Sprint 3 — Detailed Design (Week 3)**
| Task | Owner | Hours |
|------|-------|-------|
| Complete OpenAPI 3.0 spec for all endpoints | Backend | 24 |
| Design payment state machine | Backend | 8 |
| Design ledger double-entry model | Backend | 12 |
| Design component hierarchy | Frontend | 8 |
| Create UI mockups (Figma) | Design | 24 |

**Sprint 4 — Infrastructure Design (Week 4)**
| Task | Owner | Hours |
|------|-------|-------|
| Terraform module design | DevOps | 16 |
| Kubernetes resource specs | DevOps | 12 |
| CI/CD pipeline design | DevOps | 8 |
| Monitoring dashboard mockups | DevOps | 8 |
| Disaster recovery plan | DevOps | 8 |

### P3: Infrastructure (Weeks 5-7)

**Sprint 5 — Foundation (Week 5)**
| Task | Owner | Hours |
|------|-------|-------|
| AWS account setup + IAM | DevOps | 8 |
| VPC, subnets, security groups | DevOps | 8 |
| EKS cluster creation | DevOps | 12 |
| RDS PostgreSQL setup | DevOps | 8 |

**Sprint 6 — Data Layer (Week 6)**
| Task | Owner | Hours |
|------|-------|-------|
| ElastiCache Redis setup | DevOps | 8 |
| MSK Kafka setup | DevOps | 12 |
| S3 buckets with lifecycle policies | DevOps | 4 |
| Secrets Manager setup | DevOps | 4 |
| DB schema migrations | Backend | 16 |

**Sprint 7 — CI/CD + Monitoring (Week 7)**
| Task | Owner | Hours |
|------|-------|-------|
| GitHub Actions CI/CD pipeline | DevOps | 16 |
| Helm charts for all services | DevOps | 16 |
| Prometheus + Grafana setup | DevOps | 12 |
| Loki + Tempo setup | DevOps | 8 |
| PagerDuty alert configuration | DevOps | 4 |

### P4: Core Backend (Weeks 8-15)

**Sprint 8 — Auth + Merchant (Week 8)**
| Task | Owner | Hours |
|------|-------|-------|
| User registration/login API | Backend | 16 |
| JWT token generation + validation | Backend | 12 |
| MFA (TOTP) implementation | Backend | 8 |
| Merchant CRUD API | Backend | 12 |
| API key generation + validation | Backend | 8 |

**Sprint 9 — Payment Core (Week 9)**
| Task | Owner | Hours |
|------|-------|-------|
| Payment intent creation | Backend | 16 |
| State machine engine | Backend | 12 |
| Card validation (Luhn, BIN check) | Backend | 8 |
| Payment authorization flow | Backend | 16 |
| Idempotency implementation | Backend | 8 |

**Sprint 10 — Capture + Refund (Week 10)**
| Task | Owner | Hours |
|------|-------|-------|
| Payment capture | Backend | 12 |
| Payment void/cancel | Backend | 8 |
| Full refund | Backend | 12 |
| Partial refund | Backend | 8 |
| Refund state machine | Backend | 8 |

**Sprint 11 — Processor Integration (Week 11)**
| Task | Owner | Hours |
|------|-------|-------|
| Processor gateway abstraction | Backend | 12 |
| Card processor adapter (Visa) | Backend | 16 |
| Card processor adapter (Mastercard) | Backend | 8 |
| Mock processor for testing | Backend | 8 |
| Wallet adapter (bKash) | Backend | 12 |

**Sprint 12 — Ledger (Week 12)**
| Task | Owner | Hours |
|------|-------|-------|
| Ledger account model | Backend | 8 |
| Double-entry implementation | Backend | 16 |
| Balance calculation | Backend | 12 |
| Ledger validation + reconciliation | Backend | 12 |

**Sprint 13 — Webhooks (Week 13)**
| Task | Owner | Hours |
|------|-------|-------|
| Webhook delivery engine | Backend | 16 |
| Webhook retry with backoff | Backend | 12 |
| Webhook signature (HMAC) | Backend | 8 |
| Webhook delivery logs | Backend | 8 |
| Webhook endpoint management API | Backend | 8 |

**Sprint 14 — Settlement + Payout (Week 14)**
| Task | Owner | Hours |
|------|-------|-------|
| Settlement batch processing | Backend | 16 |
| Fee calculation engine | Backend | 8 |
| Payout batch creation | Backend | 12 |
| Bank file generation | Backend | 8 |
| Settlement reporting | Backend | 8 |

**Sprint 15 — Customer + Notifications (Week 15)**
| Task | Owner | Hours |
|------|-------|-------|
| Customer CRUD API | Backend | 8 |
| Saved cards + tokenization | Backend | 12 |
| Email notification service | Backend | 8 |
| SMS notification service | Backend | 8 |
| Notification templates | Backend | 8 |

### P5: Fraud + Admin (Weeks 16-18)

**Sprint 16 — Fraud Rules (Week 16)**
| Task | Owner | Hours |
|------|-------|-------|
| Velocity checking engine | Backend | 12 |
| IP reputation + geo checking | Backend | 8 |
| BIN database integration | Backend | 8 |
| Risk scoring engine | Backend | 12 |
| Configurable rule engine | Backend | 12 |

**Sprint 17 — Admin Backend (Week 17)**
| Task | Owner | Hours |
|------|-------|-------|
| Admin merchant management API | Backend | 12 |
| Admin transaction viewing API | Backend | 8 |
| Admin dispute management API | Backend | 12 |
| Admin dashboard aggregation API | Backend | 12 |
| Audit log API | Backend | 8 |

**Sprint 18 — Reporting (Week 18)**
| Task | Owner | Hours |
|------|-------|-------|
| Revenue reporting queries | Backend | 12 |
| CSV export service | Backend | 8 |
| Dashboard aggregate materialized views | Backend | 12 |
| Admin fraud dashboard API | Backend | 8 |

### P6: Frontend (Weeks 19-22)

**Sprint 19 — Design System + Auth (Week 19)**
| Task | Owner | Hours |
|------|-------|-------|
| Design system setup (Tailwind + components) | Frontend | 16 |
| Login/signup pages | Frontend | 12 |
| MFA setup flow | Frontend | 8 |
| Auth middleware + route guards | Frontend | 8 |

**Sprint 20 — Merchant Dashboard (Week 20)**
| Task | Owner | Hours |
|------|-------|-------|
| Dashboard overview page | Frontend | 12 |
| Payments list + detail pages | Frontend | 16 |
| Refunds page | Frontend | 8 |
| Customers page | Frontend | 8 |

**Sprint 21 — Merchant Features (Week 21)**
| Task | Owner | Hours |
|------|-------|-------|
| Balance + payout history | Frontend | 8 |
| API keys management | Frontend | 8 |
| Webhook configuration + logs | Frontend | 12 |
| Reports page | Frontend | 8 |
| Settings page | Frontend | 8 |

**Sprint 22 — Admin + Dev Portal (Week 22)**
| Task | Owner | Hours |
|------|-------|-------|
| Admin dashboard pages | Frontend | 16 |
| Merchant management UI | Frontend | 12 |
| Developer portal (API explorer) | Frontend | 16 |
| Webhook tester | Frontend | 8 |
| Documentation UI | Frontend | 8 |

### P7: Testing + Security (Weeks 23-25)

**Sprint 23 — Testing (Week 23)**
| Task | Owner | Hours |
|------|-------|-------|
| Load testing (k6: 1000 TPS target) | QA | 20 |
| Stress testing (5000 TPS) | QA | 12 |
| Endurance testing (8 hours) | QA | 24 |
| Chaos engineering experiments | QA | 16 |

**Sprint 24 — Security (Week 24)**
| Task | Owner | Hours |
|------|-------|-------|
| Penetration testing | Security | 24 |
| Vulnerability scanning + remediation | Security | 16 |
| PCI DSS SAQ preparation | Compliance | 24 |
| Security audit | Security | 16 |

**Sprint 25 — Bug Fix + Polish (Week 25)**
| Task | Owner | Hours |
|------|-------|-------|
| Fix critical bugs found in testing | All | 40 |
| Performance optimization | Backend | 24 |
| Documentation finalization | Tech Lead | 16 |
| Production readiness checklist | Tech Lead | 8 |

### P8: Launch (Weeks 26-27)

**Sprint 26 — Pilot Preparation (Week 26)**
| Task | Owner | Hours |
|------|-------|-------|
| Pilot merchant onboarding | Product | 24 |
| Sandbox environment launch | DevOps | 8 |
| Monitoring fine-tuning | DevOps | 12 |
| Support runbook preparation | SRE | 12 |

**Sprint 27 — Launch + Stabilize (Week 27)**
| Task | Owner | Hours |
|------|-------|-------|
| Production cutover | DevOps | 16 |
| Launch monitoring (24h watch) | All | 24 |
| Incident response drill | SRE | 4 |
| Post-launch bug fixes | All | 24 |

## Total Effort Estimate

| Phase | Person-Weeks | Cost Estimate (USD) |
|-------|--------------|---------------------|
| P1: Planning | 4 | $12,000 |
| P2: Design | 6 | $18,000 |
| P3: Infrastructure | 8 | $24,000 |
| P4: Core Backend | 32 | $96,000 |
| P5: Fraud + Admin | 9 | $27,000 |
| P6: Frontend | 16 | $48,000 |
| P7: Testing + Security | 12 | $36,000 |
| P8: Launch | 6 | $18,000 |
| **Total** | **93** | **$279,000** |

Assumptions: $75/hour blended rate, 4-person team (Backend x2, Frontend x1, DevOps x1), Product + Security + QA shared.
