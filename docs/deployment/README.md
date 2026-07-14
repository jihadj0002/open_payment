# Deployment Guide

## Deployment Methods

| Method | Best For | Complexity | Start Command |
|--------|----------|------------|---------------|
| [Docker Compose](local-development.md#docker-compose-recommended) | Local development, testing | Low | `docker-compose up` |
| [Separate Processes](local-development.md#separate-processes) | Active development, debugging | Medium | `make run` + `npm run dev` |
| [Kubernetes (Minikube)](local-development.md#kubernetes-minikube) | Staging, local K8s testing | High | `kubectl apply -f infra/k8s/` |
| [Production (AWS)](production-deployment.md) | Production deployment | High | `terraform apply` + `kubectl apply` |

## Quick Start (30 seconds)

```bash
git clone <repo>
cd open_payment
docker-compose up
```

Then visit:
- **API:** http://localhost:8080/health
- **Frontend:** http://localhost:3000
- **Mock Processor:** http://localhost:9000/health

## Architecture Overview

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Next.js App    │────▶│  Go API Server   │────▶│  PostgreSQL 16  │
│  (port 3000)    │     │  (port 8080)     │     │  (port 5432)    │
└─────────────────┘     │                  │     └─────────────────┘
                        │  Chi Router      │     ┌─────────────────┐
                        │  10 Services     │────▶│  Redis 7        │
                        │  Auth/Merchant/  │     │  (port 6379)    │
                        │  Payment/...     │     └─────────────────┘
                        └───────┬──────────┘     ┌─────────────────┐
                                │                 │  Kafka          │
                                ▼                 │  (port 9092)    │
                        ┌─────────────────┐      └─────────────────┘
                        │  Mock Processor │
                        │  (port 9000)    │
                        └─────────────────┘
```

## Prerequisites by Method

### Docker Compose
- Docker Engine 24+ and Docker Compose v2

### Separate Processes
- Go 1.22+
- Node.js 20+
- PostgreSQL 16 (local or Docker)
- Redis 7 (local or Docker)

### Kubernetes
- kubectl
- minikube or k3s
- Helm (for Chaos Mesh)

### Production
- AWS CLI configured
- Terraform 1.5+
- kubectl
- Docker Hub account
