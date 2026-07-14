# Non-Functional Requirements

## NFR-01: Performance

| ID | Requirement | Target | Measurement |
|----|------------|--------|-------------|
| NFR-01.1 | API response time (read) | <100ms p50, <200ms p95, <500ms p99 | Prometheus + Grafana |
| NFR-01.2 | API response time (write) | <200ms p50, <500ms p95, <1000ms p99 | Prometheus + Grafana |
| NFR-01.3 | Payment intent creation | <500ms p95 | APM tracing |
| NFR-01.4 | Webhook delivery | <5s p99 from event → merchant receives | Webhook latency metric |
| NFR-01.5 | Dashboard page load | <2s p95 | RUM / synthetic monitoring |
| NFR-01.6 | Report generation | <30s for 100K transactions | Batch job monitoring |
| NFR-01.7 | Settlement batch processing | <10m for 100K transactions | Batch job monitoring |

## NFR-02: Scalability

| ID | Requirement | Target | Notes |
|----|------------|--------|-------|
| NFR-02.1 | Concurrent API requests | 10,000 req/s initial, 50,000 req/s target | Horizontal scaling via K8s HPA |
| NFR-02.2 | Transaction throughput | 1,000 TPS initial, 5,000 TPS target | Kafka partitions + PG replicas |
| NFR-02.3 | Merchant capacity | 10,000 merchants | Single PG cluster |
| NFR-02.4 | Data growth | 10M transactions/month | Partition by month |
| NFR-02.5 | Auto-scaling trigger | CPU > 70% or memory > 75% | HPA + cluster autoscaler |

## NFR-03: Availability

| ID | Requirement | Target | Measurement |
|----|------------|--------|-------------|
| NFR-03.1 | System uptime (monthly) | 99.99% | Uptime check every 30s |
| NFR-03.2 | Planned maintenance window | <4h/month | Notified 7 days in advance |
| NFR-03.3 | Recovery Point Objective (RPO) | <5 minutes | WAL streaming to standby |
| NFR-03.4 | Recovery Time Objective (RTO) | <30 minutes | Automated failover |
| NFR-03.5 | Degraded mode | Read-only if payment processor down | Circuit breaker |

## NFR-04: Security

| ID | Requirement | Target | Notes |
|----|------------|--------|-------|
| NFR-04.1 | Encryption in transit | TLS 1.2 minimum, TLS 1.3 preferred | All external + internal comms |
| NFR-04.2 | Encryption at rest | AES-256 | Database, backups, S3 objects |
| NFR-04.3 | Key management | AWS KMS + HSM | Automatic key rotation |
| NFR-04.4 | API authentication | HMAC-SHA256 | Request signing |
| NFR-04.5 | Rate limiting | 100 req/s per API key | Configurable per tier |
| NFR-04.6 | Vulnerability scanning | Weekly | CI/CD + scheduled |
| NFR-04.7 | Penetration testing | Quarterly | Third-party |
| NFR-04.8 | PCI DSS compliance | Level 1 | Annual audit |
| NFR-04.9 | Audit logging | All financial + admin actions | Immutable, append-only |

## NFR-05: Reliability

| ID | Requirement | Target | Notes |
|----|------------|--------|-------|
| NFR-05.1 | Error rate (5xx) | <0.01% of requests | |
| NFR-05.2 | Idempotency guarantee | Exactly-once for all writes | Idempotency-Key required |
| NFR-05.3 | Duplicate payment prevention | 100% | Idempotency + unique constraints |
| NFR-05.4 | Data integrity | Double-entry ledger ensures sum(credits) = sum(debits) | Automatic reconciliation |
| NFR-05.5 | Graceful degradation | Payment processing continues if dashboard/analytics down | Service isolation |

## NFR-06: Maintainability

| ID | Requirement | Target | Notes |
|----|------------|--------|-------|
| NFR-06.1 | Test coverage (core services) | >80% | Go test + integration tests |
| NFR-06.2 | Test coverage (ledger) | >90% | Critical financial logic |
| NFR-06.3 | API documentation coverage | 100% of endpoints | OpenAPI 3.0 spec |
| NFR-06.4 | Code documentation | All exported Go types/functions documented | Go doc comments |
| NFR-06.5 | Deployment frequency | Daily during development | CI/CD with automated tests |
| NFR-06.6 | Rollback capability | One-click rollback to any previous version | Kubernetes rolling update |

## NFR-07: Compliance

| ID | Requirement | Target | Notes |
|----|------------|--------|-------|
| NFR-07.1 | PCI DSS | Level 1 | SAQ D for service providers |
| NFR-07.2 | PSD2 / SCA | Strong Customer Authentication | 3D Secure 2.0 |
| NFR-07.3 | GDPR | Full compliance | Data retention, deletion, consent |
| NFR-07.4 | AML / KYC | Merchant verification | Document collection + screening |
| NFR-07.5 | Data retention | 7 years for financial data | Archival policy |
| NFR-07.6 | Right to erasure | Customer data deletable within 30 days | GDPR compliance |

## NFR-08: Observability

| ID | Requirement | Target | Notes |
|----|------------|--------|-------|
| NFR-08.1 | Metrics collection | All services expose Prometheus metrics | /metrics endpoint |
| NFR-08.2 | Distributed tracing | 100% of payment flows | OpenTelemetry → Jaeger/Tempo |
| NFR-08.3 | Centralized logging | All services log to stdout → Loki | Structured JSON logs |
| NFR-08.4 | Alerting | Critical alerts within 1 minute | PagerDuty / OpsGenie |
| NFR-08.5 | Dashboards | Payment success rate, latency, error rate, volume | Grafana |
