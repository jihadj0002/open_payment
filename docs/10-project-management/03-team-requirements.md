# Team Requirements

## Role Descriptions

### Engineering Roles

| Role | Full-Time Equivalent | Key Skills | Phase Required |
|------|---------------------|------------|----------------|
| Tech Lead / Architect | 1 | Go, system design, payment domain, security | P1-P8 |
| Backend Engineer (Senior) | 1 | Go, PostgreSQL, Kafka, REST/gRPC, microservices | P3-P8 |
| Backend Engineer (Mid) | 1 | Go, SQL, REST APIs, testing | P4-P8 |
| Frontend Engineer | 1 | React/Next.js, TypeScript, Tailwind, charts | P6-P8 |
| DevOps/SRE Engineer | 1 | AWS, K8s, Terraform, CI/CD, monitoring | P2-P8 |

### Non-Engineering Roles

| Role | FTE | Key Skills | Phase |
|------|-----|------------|-------|
| Product Manager | 0.5 | Fintech domain, merchant needs, prioritization | P1-P8 |
| UI/UX Designer | 0.5 | Dashboard design, payment UX, component libraries | P2, P6 |
| QA Engineer | 0.5 | Load testing, chaos engineering, automation | P7-P8 |
| Security Engineer | 0.25 | PCI DSS, pen testing, threat modeling | P1, P2, P7 |
| Compliance Officer | 0.25 | PCI DSS, AML/KYC, GDPR, financial regulations | P1, P2, P7, P8 |

## Team Composition by Phase

| Phase | Tech Lead | Backend Sr | Backend Mid | Frontend | DevOps | Product | Designer | QA | Security | Complicance |
|-------|-----------|------------|-------------|----------|--------|---------|----------|----|----------|-------------|
| P1: Planning | 1 | 0.5 | 0 | 0 | 0.5 | 1 | 0 | 0 | 0.5 | 0.5 |
| P2: Design | 1 | 1 | 0 | 0.5 | 1 | 0.5 | 1 | 0 | 0.5 | 0.5 |
| P3: Infra | 0.5 | 0.5 | 0 | 0 | 2 | 0 | 0 | 0 | 0 | 0 |
| P4: Core | 1 | 1 | 1 | 0 | 0.5 | 0.5 | 0 | 0 | 0 | 0 |
| P5: Fraud | 0.5 | 1 | 1 | 0 | 0.5 | 0.5 | 0 | 0 | 0.5 | 0 |
| P6: Frontend | 0.5 | 0.5 | 0.5 | 2 | 0.5 | 0.5 | 0.5 | 0 | 0 | 0 |
| P7: Testing | 0.5 | 1 | 1 | 0.5 | 1 | 0.5 | 0 | 1 | 1 | 1 |
| P8: Launch | 1 | 1 | 1 | 1 | 1 | 1 | 0 | 0.5 | 0.5 | 0.5 |

## Hiring Requirements

### Must-Have Skills (Backend)
- Go: 3+ years production experience
- PostgreSQL: Complex queries, indexing, partitioning
- Kafka: Producers, consumers, schema management
- REST API design: OpenAPI, versioning, idempotency
- Testing: Unit, integration, table-driven tests
- Domain: Familiarity with payment systems (Stripe, etc.)

### Nice-to-Have Skills
- Fintech/payment gateway experience
- Double-entry accounting concepts
- PCI DSS compliance knowledge
- gRPC/protobuf experience
- Kubernetes deployment experience

### Must-Have Skills (Frontend)
- React/Next.js: 3+ years production experience
- TypeScript: Strict mode, generics
- State management: React Query, Zustand
- Testing: Vitest, Playwright
- Tailwind CSS: Component design systems

### Must-Have Skills (DevOps)
- AWS: EKS, RDS, ElastiCache, MSK, S3, IAM
- Kubernetes: Deployments, services, ingress, Helm
- Terraform: Modules, state management, CI/CD
- CI/CD: GitHub Actions
- Monitoring: Prometheus, Grafana, Loki, Tempo

## Team Growth Path

### Month 1-3 (MVP)
- Tech Lead (1)
- Backend Sr (1)
- Backend Mid (1)
- DevOps (1)
- Product (0.5)

### Month 4-6 (Dashboard)
- + Frontend (2)
- + Designer (0.5)

### Month 7-8 (Launch)
- + QA (1)
- + Security (0.5)

## Communication
- **Daily standup:** 15 min, 10am
- **Weekly sync:** 30 min sprint review + planning
- **Slack channels:** #engineering, #infra-alerts, #incident
- **Documentation:** All docs in /docs folder, PR-reviewed
- **Decision log:** ADRs in repository
