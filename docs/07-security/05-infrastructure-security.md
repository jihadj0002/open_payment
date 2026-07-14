# Infrastructure Security

## Network Security

### VPC Architecture
```
┌─────────────────────────────────────────┐
│  VPC (10.0.0.0/16)                      │
│                                         │
│  ┌──────────────────────────────────┐   │
│  │ Public Subnet (10.0.1.0/24)      │   │
│  │ - ALB (Internet-facing)          │   │
│  │ - NAT Gateway                    │   │
│  │ - Bastion Host (restricted SG)   │   │
│  └──────────────────────────────────┘   │
│                                         │
│  ┌──────────────────────────────────┐   │
│  │ Private Subnet: App (10.0.2.0/24)│   │
│  │ - EKS Worker Nodes               │   │
│  │ - All service pods               │   │
│  │ (No direct internet access)      │   │
│  └──────────────────────────────────┘   │
│                                         │
│  ┌──────────────────────────────────┐   │
│  │ Private Subnet: Data (10.0.3.0/24)│  │
│  │ - RDS PostgreSQL                  │   │
│  │ - ElastiCache Redis               │   │
│  │ - MSK Kafka Brokers               │   │
│  │ (No internet, no app access)      │   │
│  └──────────────────────────────────┘   │
│                                         │
│  ┌──────────────────────────────────┐   │
│  │ Private Subnet: Storage          │   │
│  │ (10.0.4.0/24)                    │   │
│  │ - VPC Endpoint to S3             │   │
│  └──────────────────────────────────┘   │
│                                         │
│  Security Groups:                       │
│  - ALB: 443 from Internet (Cloudflare) │
│  - App: 8080 from ALB, 9090 from Prom  │
│  - Data: 5432 from App, 6379 from App  │
└─────────────────────────────────────────┘
```

### Network Policies (Kubernetes)
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: payment-service-policy
  namespace: payment
spec:
  podSelector:
    matchLabels:
      app: payment-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: system
      podSelector:
        matchLabels:
          app: api-gateway
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: fraud-service
    - podSelector:
        matchLabels:
          app: ledger-service
    - podSelector:
        matchLabels:
          app: processor-gateway
    ports:
    - protocol: TCP
      port: 9090  # gRPC
```

## WAF (Web Application Firewall)

### Cloudflare WAF Rules
- **Rate limiting:** 1000 req/s per IP
- **Block known attack patterns:** SQLi, XSS, RFI, LFI, SSRF
- **Block TOR exit nodes** (for payment API)
- **Geoblocking:** Restrict admin dashboard to specific countries (optional)
- **Bot management:** Challenge or block known bots
- **API shield:** Protect API endpoints from abuse

### AWS WAF (on ALB)
- **IP reputation lists:** AWS Managed Rules + third-party
- **SQL injection:** AWS Managed SQLi rule group
- **XSS:** AWS Managed XSS rule group
- **Size constraints:** Block requests > 1MB
- **Rate-based rule:** 2000 req/5min per IP

## Secrets Management

### HashiCorp Vault / AWS Secrets Manager

| Secret | Storage | Rotation |
|--------|---------|----------|
| Database credentials | AWS Secrets Manager | 30 days |
| API signing keys | AWS Secrets Manager | Manual (merchant) |
| Processor API keys | Vault (encrypted with KMS) | 90 days |
| Encryption keys | AWS KMS + HSM | Annual |
| JWT signing keys | AWS KMS | 90 days |
| TLS certificates | ACM | Auto-renew |

### Secrets in Code
```go
// NEVER do this:
var StripeSecret = "sk_live_abc123..."

// ALWAYS do this:
func GetProcessorSecret(ctx context.Context, processor string) (string, error) {
    result, err := secretsManager.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: aws.String("payment/" + processor + "/api_key"),
    })
    if err != nil {
        return "", err
    }
    return *result.SecretString, nil
}
```

## Container Security
- **Base images:** Distroless or Alpine, minimal footprint
- **Vulnerability scanning:** Trivy in CI/CD + scheduled scans
- **Run as non-root:** Containers run with UID 10001
- **Read-only filesystem:** For stateless services
- **Security context:**
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 10001
  capabilities:
    drop: ["ALL"]
  readOnlyRootFilesystem: true
```

## DDoS Protection
- **Cloudflare:** L3/L4/L7 DDoS protection
- **AWS Shield Advanced:** Additional DDoS protection
- **Auto-scaling:** Handle traffic spikes
- **Rate limiting:** Per-IP and per-key limits
- **Load shedding:** Drop low-priority requests under extreme load

## IDS/IPS
- **AWS GuardDuty:** Network anomaly detection
- **Falco:** Runtime security monitoring on Kubernetes
- **WAF logs:** Sent to Security Hub for analysis
- **Anomaly detection:** ML-based on API traffic patterns

## Compliance Automation
- **PCI DSS scope:** Documented, minimized, and monitored
- **Security scanning:** Weekly full scan + CI/CD per-build
- **Configuration auditing:** AWS Config + custom rules
- **IAM analysis:** Access Analyzer for unused permissions
