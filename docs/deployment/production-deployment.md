# Production Deployment Guide

## Architecture (AWS)

```
Internet
    │
    ▼
CloudFront / Route53
    │
    ▼
ALB (HTTPS, WAF)
    │
    ├──▶ EKS Cluster (openpayment namespace)
    │       ├── openpayment-api (3-10 replicas)
    │       ├── Prometheus + Grafana
    │       └── Loki + Promtail
    │
    ├──▶ RDS PostgreSQL (Multi-AZ)
    ├──▶ ElastiCache Redis
    └──▶ MSK Kafka (optional)
```

## Prerequisites

- AWS account with admin access
- AWS CLI configured (`aws configure`)
- Terraform 1.5+
- Docker Hub account
- GitHub repository with Actions enabled

## Step 1: Infrastructure Provisioning

### Configure Terraform Backend

```bash
cd infra/terraform

# Create S3 bucket for Terraform state (one-time)
aws s3 mb s3://openpayment-terraform-state --region ap-southeast-1

# Review variables
cat variables.tf
```

### Set Required Variables

```bash
export TF_VAR_db_password="$(openssl rand -base64 32)"
export TF_VAR_environment="production"
```

### Deploy Infrastructure

```bash
terraform init
terraform plan -out=tfplan
terraform apply tfplan
```

This creates:
- **VPC** with 3 public + 3 private subnets across AZs
- **RDS PostgreSQL** (db.r6g.large, Multi-AZ, 30-day backups)
- **ElastiCache Redis** (cache.r6g.large, 7-day backups)
- **EKS Cluster** (Kubernetes 1.28, managed node group with 2-6 t3.large nodes)

### Capture Outputs

```bash
terraform output
# vpc_id = "vpc-xxxxx"
# eks_cluster_name = "openpayment-eks-xxxxx"
# rds_endpoint = "openpayment-rds.xxxxx.ap-southeast-1.rds.amazonaws.com"
# redis_endpoint = "openpayment-redis.xxxxx.ap-southeast-1.cache.amazonaws.com"
```

## Step 2: Configure Kubernetes

### Update kubeconfig

```bash
aws eks update-kubeconfig --name $(terraform output -raw eks_cluster_name) --region ap-southeast-1
```

### Deploy Infrastructure Components

```bash
# Create namespace
kubectl apply -f infra/k8s/namespace.yaml

# Create ConfigMap
kubectl apply -f infra/k8s/configmap.yaml

# Create Secrets (update with actual values)
# Edit infra/k8s/secret.yaml with RDS endpoint, Redis endpoint, JWT secret
kubectl apply -f infra/k8s/secret.yaml

# Deploy application
kubectl apply -f infra/k8s/deployment.yaml
kubectl apply -f infra/k8s/service.yaml
kubectl apply -f infra/k8s/hpa.yaml

# Configure ingress
kubectl apply -f infra/k8s/ingress.yaml
```

### Verify Deployment

```bash
# Check pods
kubectl get pods -n openpayment -w

# Check services
kubectl get svc -n openpayment

# Check ingress (get ALB DNS)
kubectl get ingress -n openpayment -w
```

## Step 3: Configure CI/CD

### GitHub Secrets

Add these to your GitHub repository (Settings → Secrets and variables → Actions):

| Secret | Value |
|--------|-------|
| `DOCKER_USERNAME` | Docker Hub username |
| `DOCKER_PASSWORD` | Docker Hub password or token |
| `AWS_ACCESS_KEY_ID` | IAM user access key (deploy permissions) |
| `AWS_SECRET_ACCESS_KEY` | IAM user secret key |
| `KUBE_CONFIG_DATA` | Base64-encoded kubeconfig |

### CI Pipeline

The CI pipeline (`ci.yml`) runs on every push/PR to `main`/`develop`:
1. Go lint + vet
2. Unit tests (with PostgreSQL service)
3. Go build
4. Frontend lint + build

### CD Pipeline

The CD pipeline (`cd.yml`) runs on push to `main`:
1. Build and push Docker images to Docker Hub
2. (Future: automatic K8s rollout via ArgoCD or GitHub Actions)

### Security Scan

The security workflow (`security.yml`) runs weekly:
- Trivy filesystem vulnerability scan
- Gosec Go security analysis
- Results uploaded as SARIF to GitHub Security tab

## Step 4: Monitoring & Alerting

### Prometheus + Grafana

```bash
# Deploy Prometheus stack (via helm or manifests)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

# Apply custom alerts
kubectl apply -f infra/monitoring/prometheus/alerts.yml

# Import Grafana dashboard
# Open Grafana (port-forward: kubectl port-forward svc/prometheus-grafana 3000:80 -n monitoring)
# Import infra/monitoring/grafana/dashboards/api-performance.json
```

### Loki + Promtail

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki -n logging --create-namespace
helm install promtail grafana/promtail -n logging

# Add Loki datasource in Grafana using infra/logging/grafana/datasources/loki.yml
```

## Step 5: DNS & SSL

```bash
# Configure Route53
# Create A record pointing to ALB DNS name

# SSL is automatic via AWS Certificate Manager (ACM)
# The ALB ingress annotation handles HTTPS termination
```

## Step 6: Post-Deployment Verification

```bash
# Full API smoke test
curl -X POST https://api.openpayment.com/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Prod Test","email":"test@prod.com","password":"Test123!"}'

# Health check
curl https://api.openpayment.com/v1/health

# Check Grafana dashboards
# - API Performance dashboard shows request rate, latency, errors
# - Check Loki logs for any errors
```

## Scaling Considerations

| Component | Scaling Strategy | Initial Size |
|-----------|-----------------|--------------|
| API Server | HPA (CPU > 70%) | 3 pods, max 10 |
| PostgreSQL | RDS instance upgrade, read replicas | db.r6g.large |
| Redis | ElastiCache cluster mode | cache.r6g.large |
| Kafka | Add partitions, brokers | 3 brokers (MSK) |

## Disaster Recovery

| Scenario | RTO | RPO | Recovery Steps |
|----------|-----|-----|----------------|
| Pod failure | < 30s | N/A | K8s auto-restarts |
| Node failure | < 5min | N/A | EKS replaces node |
| AZ outage | < 15min | < 5min | Multi-AZ RDS failover |
| Region outage | < 2hr | < 1hr | DR region with Terraform |
| Data corruption | < 4hr | < 24hr | RDS point-in-time restore |
