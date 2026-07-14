# Step-by-Step Build Process

This is the day-by-day execution plan from zero to production. Each step is a concrete action with deliverables.

## Week 1-2: Foundation

### Day 1-2: Project Setup
```
[ ] 1. Create Go module: `go mod init github.com/openpayment/gateway`
[ ] 2. Set up project directory structure (cmd, internal, pkg)
[ ] 3. Create main.go with health check endpoint
[ ] 4. Set up Dockerfile with multi-stage build
[ ] 5. Run `go mod tidy`
[ ] 6. Verify: `curl localhost:8080/health` returns 200

Deliverables: Go module, main.go, Dockerfile, health check endpoint
```

### Day 3-4: Configuration
```
[ ] 1. Create config.go with env-var based configuration
[ ] 2. Implement config struct: Port, DB, Redis, Kafka, JWT, LogLevel
[ ] 3. Add godotenv support for dev mode
[ ] 4. Create config.example.env file
[ ] 5. Write config unit tests

Deliverables: Config loading, example env file
```

### Day 5: Logging + Middleware
```
[ ] 1. Set up zerolog/slog structured logger
[ ] 2. Implement request ID middleware (X-Request-Id)
[ ] 3. Implement request logging middleware
[ ] 4. Implement panic recovery middleware
[ ] 5. Implement CORS middleware
[ ] 6. Verify: each request has request_id in log

Deliverables: Logger, middleware stack (request_id, logging, recovery, CORS)
```

### Day 6-7: Database Setup
```
[ ] 1. Set up PostgreSQL connection (pgx pool)
[ ] 2. Install golang-migrate
[ ] 3. Create first migration: 000001_create_merchants.up.sql
[ ] 4. Create initial schema (merchants, merchant_users)
[ ] 5. Write migration runner
[ ] 6. Write repository base (transaction support)
[ ] 7. Verify: migration runs, tables created

Deliverables: DB connection, migrations, base repository
```

### Day 8-10: CI/CD Infrastructure
```
[ ] 1. Create .github/workflows/ci.yml
[ ] 2. Add lint step (golangci-lint)
[ ] 3. Add test step (go test ./...)
[ ] 4. Add build step (docker build)
[ ] 5. Add security scan (Trivy)
[ ] 6. Configure ECR repository
[ ] 7. Write Helm chart for the app
[ ] 8. Deploy to staging cluster
[ ] 9. Verify: PR creates build, merge deploys to staging

Deliverables: CI/CD pipeline, Helm chart, staging deployment
```

## Week 3-4: Auth + Merchant Service

### Day 11-13: Auth Service
```
[ ] 1. Create auth service files (service.go, handler.go, repository.go)
[ ] 2. Implement password hashing (bcrypt, cost=12)
[ ] 3. Implement user registration endpoint: POST /v1/auth/register
[ ] 4. Implement login endpoint: POST /v1/auth/login
[ ] 5. Implement JWT generation (access + refresh tokens)
[ ] 6. Implement JWT validation middleware
[ ] 7. Implement MFA setup (TOTP)
[ ] 8. Implement MFA verification during login
[ ] 9. Write auth unit tests

Deliverables: Registration, login, MFA, JWT middleware
```

### Day 14-16: API Key Management
```
[ ] 1. Implement API key generation (sk_live_ / sk_test_)
[ ] 2. Implement key hashing (SHA-256) for storage
[ ] 3. Implement key lookup by hash
[ ] 4. Implement API key validation middleware
[ ] 5. Implement CRUD endpoints for API keys
[ ] 6. Implement IP whitelist check
[ ] 7. Write API key tests

Deliverables: API key endpoints, validation middleware
```

### Day 17-19: Merchant Service
```
[ ] 1. Implement merchant CRUD endpoints
[ ] 2. Implement merchant status workflow (pending → active → suspended)
[ ] 3. Implement KYC document upload (S3 storage)
[ ] 4. Implement merchant settings CRUD
[ ] 5. Implement webhook configuration endpoints
[ ] 6. Write merchant service tests

Deliverables: Merchant management, KYC, webhook config
```

### Day 20: Rate Limiting
```
[ ] 1. Implement Redis-based sliding window rate limiter
[ ] 2. Add rate limit config per API key
[ ] 3. Implement X-RateLimit-* headers
[ ] 4. Implement 429 response on exceeded limit
[ ] 5. Write rate limiter tests

Deliverables: Rate limiter middleware
```

## Week 5-7: Payment Core

### Day 21-24: Payment Intent
```
[ ] 1. Create payment state machine engine
[ ] 2. Implement all state transitions with validation
[ ] 3. Write 100% test coverage for state machine
[ ] 4. Implement POST /v1/payments (create)
[ ] 5. Implement GET /v1/payments/:id (retrieve)
[ ] 6. Implement GET /v1/payments (list with filters)
[ ] 7. Implement idempotency on payment creation
[ ] 8. Implement client_secret generation
[ ] 9. Write payment intent tests

Deliverables: Payment CRUD, state machine, idempotency
```

### Day 25-27: Card Processing
```
[ ] 1. Implement card number validation (Luhn algorithm)
[ ] 2. Implement BIN lookup (card brand, issuing bank, country)
[ ] 3. Implement card expiration and CVC validation
[ ] 4. Create processor client interface
[ ] 5. Implement mock processor client
[ ] 6. Implement authorization flow: POST /v1/payments/:id/confirm
[ ] 7. Implement capture: POST /v1/payments/:id/capture
[ ] 8. Implement void: POST /v1/payments/:id/void
[ ] 9. Write card processing tests

Deliverables: Card validation, processor integration, auth/capture
```

### Day 28-30: Refunds
```
[ ] 1. Implement POST /v1/refunds (full refund)
[ ] 2. Implement partial refund logic
[ ] 3. Implement GET /v1/refunds (list)
[ ] 4. Implement GET /v1/refunds/:id (retrieve)
[ ] 5. Implement refund state machine
[ ] 6. Implement refund idempotency
[ ] 7. Write refund tests

Deliverables: Full/partial refund, refund history
```

### Day 31-33: Webhook Engine
```
[ ] 1. Implement webhook delivery via Kafka
[ ] 2. Implement webhook HTTP POST with HMAC signature
[ ] 3. Implement retry with exponential backoff (max 5)
[ ] 4. Implement dead letter queue for failed deliveries
[ ] 5. Implement webhook delivery logging
[ ] 6. Implement webhook endpoint CRUD API
[ ] 7. Write webhook tests

Deliverables: Webhook delivery engine, retry, logging
```

### Day 34-35: Kafka Integration
```
[ ] 1. Set up Kafka producer
[ ] 2. Define all event topics (payment.created, payment.authorized, etc.)
[ ] 3. Set up Kafka consumer groups
[ ] 4. Implement event schema validation
[ ] 5. Write event producer/consumer integration tests

Deliverables: Full event streaming integration
```

## Week 8-9: Ledger + Settlement

### Day 36-39: Ledger
```
[ ] 1. Implement ledger account model (setup on merchant creation)
[ ] 2. Implement double-entry creation
[ ] 3. Implement balance calculation from entries
[ ] 4. Implement balance validation (debits = credits)
[ ] 5. Implement GET /v1/balance endpoint
[ ] 6. Implement fee and tax calculation
[ ] 7. Write 100% coverage ledger tests

Deliverables: Double-entry ledger, balance API
```

### Day 40-42: Settlement
```
[ ] 1. Implement daily settlement batch processing (cron job)
[ ] 2. Implement settlement line items per transaction
[ ] 3. Implement payout batch creation
[ ] 4. Implement bank transfer file generation
[ ] 5. Implement settlement webhook events
[ ] 6. Write settlement tests

Deliverables: Settlement batch, payout, bank files
```

## Week 10-11: Fraud + Admin

### Day 43-45: Fraud Service
```
[ ] 1. Implement IP reputation check (maxmind GeoIP2)
[ ] 2. Implement velocity check (transactions per time window)
[ ] 3. Implement BIN country mismatch detection
[ ] 4. Implement VPN/TOR detection
[ ] 5. Implement configurable rule engine
[ ] 6. Implement risk scoring (weighted rule scores)
[ ] 7. Integrate fraud check into payment flow
[ ] 8. Write fraud tests

Deliverables: Fraud rule engine, risk scoring
```

### Day 46-48: Admin Backend
```
[ ] 1. Implement admin merchant management APIs
[ ] 2. Implement admin transaction viewer
[ ] 3. Implement dispute management
[ ] 4. Implement audit log middleware
[ ] 5. Implement system configuration API
[ ] 6. Implement fee/tax configuration
[ ] 7. Write admin API tests

Deliverables: Admin API endpoints
```

## Week 12-14: Frontend

### Day 49-51: Design System + Auth
```
[ ] 1. Initialize Next.js project
[ ] 2. Set up Tailwind config with design tokens
[ ] 3. Build UI components (Button, Card, Table, Modal, Badge, Input)
[ ] 4. Set up React Query client provider
[ ] 5. Set up Zustand stores (auth, ui)
[ ] 6. Build login page
[ ] 7. Build MFA setup page
[ ] 8. Implement auth middleware (Next.js middleware)

Deliverables: Design system, login, MFA
```

### Day 52-55: Merchant Dashboard
```
[ ] 1. Build sidebar + header layout
[ ] 2. Build dashboard overview (stat cards, charts)
[ ] 3. Build payments list (table, filters, search, pagination)
[ ] 4. Build payment detail (timeline, actions, refund button)
[ ] 5. Build refund initiation modal
[ ] 6. Build customers list + detail
[ ] 7. Build balance page
[ ] 8. Build API keys page
[ ] 9. Build webhooks page
[ ] 10. Build reports page
[ ] 11. Build settings page

Deliverables: Full merchant dashboard
```

### Day 56-58: Admin + Dev Portal
```
[ ] 1. Build admin dashboard overview
[ ] 2. Build merchant management pages
[ ] 3. Build KYC review interface
[ ] 4. Build dispute management
[ ] 5. Build fraud dashboard
[ ] 6. Build API explorer (interactive docs)
[ ] 7. Build webhook tester
[ ] 8. Build API documentation page

Deliverables: Admin dashboard, developer portal
```

## Week 15-16: Testing + Launch

### Day 59-61: Load + Security Testing
```
[ ] 1. Run k6 load tests (100/1000/5000 TPS)
[ ] 2. Fix bottlenecks found in load tests
[ ] 3. Run chaos experiments (pod kill, network partition)
[ ] 4. Run external pen test
[ ] 5. Fix all critical/high findings
[ ] 6. Complete PCI SAQ D questionnaire

Deliverables: Load test results, pen test report, PCI SAQ
```

### Day 62-63: Production Launch
```
[ ] 1. Deploy to production cluster
[ ] 2. Run smoke tests against production
[ ] 3. Verify monitoring and alerts
[ ] 4. Activate sandbox environment
[ ] 5. Onboard pilot merchants
[ ] 6. Begin 24h monitoring rotation
[ ] 7. Verify: merchants can create live payments

Deliverables: LIVE payment gateway processing real transactions
```

## Post-Launch (Week 17+)

### Day 64-66: Stabilization
```
[ ] 1. Fix production bugs as discovered
[ ] 2. Tune monitoring alerts
[ ] 3. Add missing test coverage
[ ] 4. Write post-launch retrospective
[ ] 5. Plan v1.1 features (subscriptions, mobile app, multi-currency)

Deliverables: Stable production system, retrospective, v1.1 plan
```
