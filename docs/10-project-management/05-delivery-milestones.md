# Delivery Milestones

## Milestone Overview

| Milestone | Date (Target) | Phase | Criteria | Stakeholder Sign-off |
|-----------|--------------|-------|----------|---------------------|
| M1: Architecture Sign-off | Week 2 | P1 | SRS, ADRs, ERD, API spec, Threat model | Tech Lead + Security |
| M2: Infrastructure Ready | Week 7 | P3 | K8s cluster, CI/CD, Monitoring, Staging env | DevOps + Tech Lead |
| M3: Core Payment Flow | Week 11 | P4 | Payment create → authorize → capture → webhook | Tech Lead + QA |
| M4: Ledger + Settlement | Week 14 | P4 | Ledger entries, balance, settlement batch, payout | Tech Lead + Finance |
| M5: Merchant Dashboard | Week 21 | P6 | Merchant dashboard functional (payments, refunds, reports) | Product + Design |
| M6: Admin + Dev Portal | Week 22 | P6 | Admin dashboard, API explorer, webhook tester | Product + Tech Lead |
| M7: Security + Compliance | Week 24 | P7 | Pen test passed, PCI SAQ submitted, vulns fixed | Security + Compliance |
| M8: Pilot Launch | Week 26 | P8 | 5-10 merchants processing live transactions | Product + CEO |
| M9: Production Launch | Week 27 | P8 | Public availability, monitoring, support active | All |

## M1: Architecture Sign-off

### Criteria
- [ ] All C4 diagrams reviewed and approved
- [ ] OpenAPI 3.0 spec covering all MVP endpoints
- [ ] Database ER diagrams reviewed
- [ ] Payment state machine defined
- [ ] Kafka topics and event schemas defined
- [ ] Security threat model completed
- [ ] Infrastructure design reviewed
- [ ] Tech stack finalized with ADRs

### Sign-off
- **Tech Lead:** ___
- **Security Lead:** ___
- **Date:** ___

## M2: Infrastructure Ready

### Criteria
- [ ] AWS VPC, subnets, security groups provisioned
- [ ] EKS cluster running with worker nodes
- [ ] RDS PostgreSQL (Multi-AZ) deployed
- [ ] ElastiCache Redis cluster deployed
- [ ] MSK Kafka cluster deployed
- [ ] S3 buckets with lifecycle policies
- [ ] Secrets Manager configured
- [ ] CI/CD pipeline building and deploying
- [ ] Prometheus + Grafana stack operational
- [ ] Loki log aggregation working
- [ ] Staging environment deployed
- [ ] Helm charts for core services

### Sign-off
- **DevOps:** ___
- **Tech Lead:** ___
- **Date:** ___

## M3: Core Payment Flow

### Criteria
- [ ] Payment intent creation via API
- [ ] Card validation (Luhn, BIN, CVC, expiry)
- [ ] Fraud checks (basic velocity, IP)
- [ ] Payment authorization with processor
- [ ] Payment capture
- [ ] Payment void/cancel
- [ ] Full and partial refund
- [ ] Webhook delivery for all payment events
- [ ] Idempotency working
- [ ] State machine prevents invalid transitions
- [ ] Integration tests passing for full flow
- [ ] API documentation for payment endpoints

### Sign-off
- **Tech Lead:** ___
- **QA:** ___
- **Date:** ___

## M4: Ledger + Settlement

### Criteria
- [ ] Double-entry ledger on every capture/refund
- [ ] Balance calculation accurate
- [ ] Daily settlement batch processing
- [ ] Fee and tax calculation correct
- [ ] Payout batch creation
- [ ] Bank transfer file generation
- [ ] Balance validation (debits = credits)
- [ ] Webhook for settlement events
- [ ] 100% test coverage on ledger
- [ ] Reconciliation check passes

### Sign-off
- **Tech Lead:** ___
- **Finance/Compliance:** ___
- **Date:** ___

## M5: Merchant Dashboard

### Criteria
- [ ] Login/register flow complete
- [ ] MFA setup and verification
- [ ] Dashboard overview with stats charts
- [ ] Payment list with filters and search
- [ ] Payment detail with timeline
- [ ] Refund initiation
- [ ] Customer list and detail
- [ ] Balance and payout history
- [ ] API keys management
- [ ] Webhook configuration and logs
- [ ] Reports with export
- [ ] Settings (profile, users, security)
- [ ] Role-based access control
- [ ] Responsive design (desktop + tablet)
- [ ] Error and empty states handled

### Sign-off
- **Product:** ___
- **Design:** ___
- **Date:** ___

## M6: Admin + Dev Portal

### Criteria
- [ ] Admin dashboard with platform stats
- [ ] Merchant management (list, detail, approve, suspend)
- [ ] KYC document review interface
- [ ] Transaction viewer (all merchants)
- [ ] Dispute management
- [ ] Fraud dashboard (rules, alerts)
- [ ] Settlement monitoring
- [ ] Fee configuration
- [ ] Audit log viewer
- [ ] API explorer (interactive docs)
- [ ] Webhook tester tool
- [ ] API reference documentation
- [ ] Quick start guide

### Sign-off
- **Product:** ___
- **Tech Lead:** ___
- **Date:** ___

## M7: Security + Compliance

### Criteria
- [ ] External penetration test passed (no critical/high findings)
- [ ] All high-severity vulnerabilities fixed
- [ ] PCI DSS SAQ D completed and submitted
- [ ] Dependency vulnerabilities resolved
- [ ] Secret scanning passes
- [ ] Container scanning passes
- [ ] Rate limiting validated
- [ ] Authentication/authorization audited
- [ ] Audit log completeness verified
- [ ] Incident response plan validated

### Sign-off
- **Security Lead:** ___
- **Compliance:** ___
- **Date:** ___

## M8: Pilot Launch

### Criteria
- [ ] Sandbox environment live
- [ ] Pilot merchants onboarded (5+)
- [ ] Monitoring and alerting tuned
- [ ] Operational runbook written
- [ ] Support channels established
- [ ] Incident response tested
- [ ] Pilot merchants processing live transactions
- [ ] No SEV-1 incidents in first week
- [ ] Payment success rate > 90%
- [ ] Merchant feedback collected

### Sign-off
- **Product:** ___
- **CEO:** ___
- **Date:** ___

## M9: Production Launch

### Criteria
- [ ] Public signup open
- [ ] Documentation website live
- [ ] SDKs published (Go, Python, JS)
- [ ] Status page active
- [ ] 24/7 monitoring and on-call rotation
- [ ] Support ticketing system active
- [ ] SLA published
- [ ] All launch-blocking bugs resolved
- [ ] 72h post-launch monitoring completed
- [ ] Retrospective scheduled

### Sign-off
- **Tech Lead:** ___
- **CEO:** ___
- **Date:** ___
