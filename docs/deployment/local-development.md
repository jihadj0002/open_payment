# Local Development Guide

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/config/config.go`

## Method 1: Docker Compose (Recommended)

### Prerequisites
- Docker Engine 24+ with Compose v2 plugin

### Steps

```bash
# 1. Clone the repository
git clone <repo-url>
cd open_payment

# 2. Start all services
docker-compose up --build

# 3. Verify everything is running
curl http://localhost:8080/health
# → {"data":{"status":"ok","service":"open-payment-gateway"}}

curl http://localhost:9000/health
# → {"name":"mock-processor","type":"payment","healthy":true}
```

### What starts

| Service | Container | Port | Purpose |
|---------|-----------|------|---------|
| API Server | openpayment-app | 8080 | Go HTTP API |
| PostgreSQL | openpayment-postgres | 5432 | Primary database |
| Redis | openpayment-redis | 6379 | Caching, rate limiting |
| Kafka | openpayment-kafka | 9092 | Event bus |
| Zookeeper | openpayment-zookeeper | 2181 | Kafka coordinator |
| Mock Processor | openpayment-mock-processor | 9000 | Payment simulator |

### Test the API

```bash
# Health check
curl localhost:8080/health

# Register a merchant
curl -X POST localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Test Merchant","email":"test@example.com","password":"password123"}'
# → Returns JWT tokens

# Create a payment (use the access_token from register)
curl -X POST localhost:8080/payments \
  -H 'Authorization: Bearer <access_token>' \
  -H 'Content-Type: application/json' \
  -d '{"amount":1000,"currency":"BDT","payment_method":"card","confirm":true}'

# Check balance
curl localhost:8080/balance \
  -H 'Authorization: Bearer <access_token>'
```

---

## Method 2: Separate Processes

Use this when actively developing and need faster rebuilds.

### Prerequisites
- Go 1.22+, Node.js 20+, PostgreSQL 16, Redis 7

### Step 1: Start Infrastructure

```bash
# Start PostgreSQL
docker run -d --name pg \
  -e POSTGRES_DB=paymentdb \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Start Redis
docker run -d --name redis \
  -p 6379:6379 \
  redis:7-alpine

# Start Kafka (optional — only if working on event-driven features)
docker run -d --name zookeeper \
  -p 2181:2181 \
  confluentinc/cp-zookeeper:7.5.0

docker run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_BROKER_ID=1 \
  -e KAFKA_ZOOKEEPER_CONNECT=localhost:2181 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  confluentinc/cp-kafka:7.5.0
```

### Step 2: Start Mock Processor

```bash
cd mock-processor
go run . &
# Runs on port 9000
```

### Step 3: Start Backend

```bash
# Terminal 1
cd open_payment
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/paymentdb?sslmode=disable
export REDIS_URL=redis://localhost:6379/0
export JWT_SECRET=<your-64-char-hex-secret>
make run
# Runs on port 8080
```

Or with hot-reload (requires `air`):
```bash
make dev
```

### Step 4: Start Frontend

```bash
# Terminal 2
cd web
npm install
npm run dev
# Runs on port 3000
```

### Step 5: Run Tests

```bash
# Unit tests
make test

# With verbose output
go test ./... -v -count=1 -timeout 60s
```

---

## Method 3: Kubernetes (Minikube)

### Prerequisites
- minikube, kubectl, Docker

### Steps

```bash
# 1. Start minikube
minikube start --cpus 4 --memory 8192

# 2. Build images into minikube's Docker
eval $(minikube docker-env)
docker build -t openpayment/gateway:latest .
cd mock-processor && docker build -t openpayment/mock-processor:latest .

# 3. Install dependencies (Postgres + Redis via Helm)
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgres bitnami/postgresql \
  --set auth.database=paymentdb \
  --set auth.username=postgres \
  --set auth.password=postgres
helm install redis bitnami/redis --set auth.enabled=false

# 4. Deploy the app
kubectl apply -f infra/k8s/namespace.yaml
kubectl apply -f infra/k8s/configmap.yaml
# Update secret.yaml with correct DB and Redis URLs
kubectl apply -f infra/k8s/secret.yaml
kubectl apply -f infra/k8s/deployment.yaml
kubectl apply -f infra/k8s/service.yaml
kubectl apply -f infra/k8s/hpa.yaml

# 5. Access the API
kubectl port-forward svc/openpayment-api 8080:8080 -n openpayment
curl localhost:8080/health
```

---

## Troubleshooting

### Database Connection Failures
```bash
# Verify PostgreSQL is running
docker ps | grep postgres

# Test connection
psql -h localhost -U postgres -d paymentdb -c "SELECT 1"

# Check logs
docker logs openpayment-postgres
```

### Migration Errors
```bash
# Migrations run automatically on startup
# Check startup logs:
docker logs openpayment-app

# To re-run migrations, restart the app:
docker-compose restart app
```

### "Port already in use"
```bash
# Find what's using the port
lsof -i :8080

# Change port in .env:
# PORT=8081
```

### Go Build Errors
```bash
# Clean and rebuild
make clean && make build

# Update dependencies
go mod tidy
```
