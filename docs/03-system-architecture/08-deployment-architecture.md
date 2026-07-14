# Deployment Architecture

## Kubernetes Cluster Design

```
┌────────────────────────────────────────────────────────────────────┐
│                      EKS Cluster (ap-south-1)                      │
│                                                                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │  System NS   │  │  Payment NS  │  │  Admin NS    │            │
│  │  - ingress   │  │  - payment   │  │  - admin-svc │            │
│  │  - cert-man  │  │  - merchant  │  │  - reporting  │            │
│  │  - monitoring│  │  - fraud     │  │  - settlement │            │
│  │  - logging   │  │  - auth      │  │              │            │
│  └──────────────┘  │  - ledger    │  └──────────────┘            │
│                     │  - webhook   │                               │
│                     │  - notif     │                               │
│                     └──────────────┘                               │
└────────────────────────────────────────────────────────────────────┘
```

## Namespace Structure

| Namespace | Purpose | Resources |
|-----------|---------|-----------|
| `system` | Cluster infrastructure | ingress-nginx, cert-manager, external-dns |
| `monitoring` | Observability | prometheus, grafana, loki, tempo, opentelemetry-collector |
| `payment` | Core payment services | auth, merchant, customer, payment, fraud, processor-gateway |
| `ledger` | Financial services | ledger, settlement |
| `webhook` | Webhook delivery | webhook-service |
| `notification` | Notifications | notification-service |
| `admin` | Admin and reporting | admin-service, reporting-service |
| `frontend` | Dashboard frontends | merchant-dashboard, admin-dashboard, developer-portal |
| `data` | Data infrastructure | kafka, redis, opensearch |
| `ci-cd` | Build and deploy | jenkins-agent, argo-cd |

## Resource Allocation per Service (Production)

| Service | CPU Request | CPU Limit | Memory Request | Memory Limit | Replicas (min/max) |
|---------|-------------|-----------|----------------|--------------|---------------------|
| API Gateway | 1 | 2 | 2Gi | 4Gi | 3 / 10 |
| Auth | 0.5 | 1 | 512Mi | 1Gi | 2 / 6 |
| Merchant | 0.5 | 1 | 512Mi | 1Gi | 2 / 6 |
| Customer | 0.5 | 1 | 512Mi | 1Gi | 2 / 4 |
| Payment | 1 | 2 | 1Gi | 2Gi | 3 / 10 |
| Fraud | 1 | 2 | 1Gi | 2Gi | 2 / 6 |
| Processor Gateway | 0.5 | 1 | 512Mi | 1Gi | 2 / 8 |
| Ledger | 0.5 | 1 | 512Mi | 1Gi | 2 / 6 |
| Settlement | 0.5 | 1 | 1Gi | 2Gi | 1 / 3 |
| Webhook | 1 | 2 | 1Gi | 2Gi | 2 / 8 |
| Notification | 0.3 | 0.5 | 256Mi | 512Mi | 1 / 3 |
| Admin | 0.5 | 1 | 512Mi | 1Gi | 1 / 3 |

## Node Groups

| Node Group | Instance Type | Min Size | Max Size | Use Case |
|------------|---------------|----------|----------|----------|
| `system` | t3.medium | 2 | 4 | System components, monitoring |
| `services` | m6i.large | 3 | 20 | Core microservices |
| `data` | r6i.large | 3 | 6 | Kafka, Redis (if self-hosted) |
| `batch` | m6i.large | 1 | 5 | Settlement, reporting, analytics jobs |

## Deployment Strategy

### Blue-Green Deployment
```
1. Deploy new version (green) alongside current (blue)
2. Run smoke tests against green
3. Switch load balancer from blue to green
4. Monitor for 5 minutes
5. If healthy, terminate blue
6. If errors, switch back to blue immediately
```

### Canary Deployment (for payment service)
```
1. Deploy new version
2. Route 5% of traffic to new version (2 min)
3. Route 25% -> 2 min
4. Route 50% -> 2 min
5. Route 100%
6. If error rate > 0.1% at any step, auto-rollback
```

## Ingress Configuration
- **External:** AWS ALB via ingress-nginx
- **Internal:** Cluster-local gRPC (no ingress)
- **TLS:** cert-manager with Let's Encrypt (production) + AWS Certificate Manager
- **WAF:** AWS WAF rules in front of ALB

## CI/CD Pipeline (GitHub Actions)

```
Git push → main branch
    │
    ▼
[1] Lint + Format Check (golangci-lint, prettier)
    │
    ▼
[2] Unit Tests (go test ./...)
    │
    ▼
[3] Build Docker Images (docker build)
    │
    ▼
[4] Security Scan (Trivy, Snyk)
    │
    ▼
[5] Push to ECR (AWS Elastic Container Registry)
    │
    ▼
[6] Deploy to Staging (helm upgrade --install)
    │
    ▼
[7] Integration Tests (against staging)
    │
    ▼
[8] Deploy to Production (canary or blue-green)
    │
    ▼
[9] Smoke Tests (production)
    │
    ▼
[10] Slack Notification (deployment success/failure)
```

## Health Check Endpoints
Every service exposes:
- `GET /health` — Returns 200 if service is healthy
- `GET /ready` — Returns 200 if service is ready to accept traffic
- `GET /metrics` — Prometheus metrics

## Backup and DR
See `08-infrastructure/05-disaster-recovery.md` for detailed procedures.
