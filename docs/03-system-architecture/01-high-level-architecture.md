# High-Level System Architecture (C4 Model)

## Level 1: System Context Diagram

```
[Customer] --(pays via)--> [Merchant Website/App]
                                |
                     (uses)     |     (uses)
                        v       v         v
                   [Open Payment Gateway System]
                                |
                    +-----------+-----------+
                    |                       |
            [Card Networks]          [Banking System]
            (Visa, Mastercard)    (ACH, Bank Transfers)
```

## Level 2: Container Diagram

```
                          Internet
                             |
                      Cloudflare WAF + CDN
                             |
                   API Gateway (Envoy/NGINX)
                             |
         +-------------------+-------------------+
         |                                       |
  [Auth Service]                         [Rate Limiter]
  (OAuth2, JWT, MFA)                     (Redis-based)
         |                                       |
         +-------------------+-------------------+
                             |
              +--------------+--------------+
              |                             |
   [Merchant API Service]          [Admin API Service]
   (REST port 8080)                (REST port 8081)
              |                             |
              +--------------+--------------+
                             |
                    [Payment Orchestrator]
                             |
         +-------------------+-------------------+
         |                   |                   |
   [Fraud Engine]    [Processor Gateway]   [Ledger Service]
   (checks, ML)      (Vis, Mastercard)    (double-entry)
         |                   |                   |
         +-------------------+-------------------+
                             |
                       [Kafka Cluster]
                             |
         +-------------------+-------------------+
         |                   |                   |
   [Webhook Service]   [Notification Svc]  [Settlement Svc]
         |                   |                   |
         +-------------------+-------------------+
                             |
                    [Data Storage Layer]
         +-------------------+-------------------+
         |                   |                   |
   [PostgreSQL]         [Redis]            [S3/ES]
   (Primary DB)     (Cache, Sessions)    (Logs, Docs)
```

## Level 3: Component Diagram (Payment Service)

```
[Payment Service]
  |
  +-- HTTP Handler (REST API layer)
  |     +-- POST /v1/payments
  |     +-- GET /v1/payments/{id}
  |     +-- POST /v1/payments/{id}/capture
  |     +-- POST /v1/payments/{id}/refund
  |
  +-- Service Layer (Business logic)
  |     +-- PaymentIntentService
  |     |     +-- CreatePaymentIntent()
  |     |     +-- AuthorizePayment()
  |     |     +-- CapturePayment()
  |     |     +-- RefundPayment()
  |     +-- FraudCheckService (calls Fraud Engine)
  |     +-- ProcessorService (calls Processor Gateway)
  |     +-- LedgerService (calls Ledger via gRPC)
  |
  +-- Repository Layer (Data access)
  |     +-- PaymentIntentRepository
  |     +-- TransactionRepository
  |     +-- CustomerRepository
  |
  +-- Event Producer (Kafka)
        +-- payment.created
        +-- payment.authorized
        +-- payment.captured
        +-- payment.refunded
        +-- payment.failed
```

## Network Architecture

```
Internet
    |
Cloudflare (WAF, DDoS protection, TLS termination)
    |
AWS Route53 (DNS)
    |
ALB (Application Load Balancer) — public subnet
    |
+---[Public Subnet]---+
|   Bastion Host      |
|   NAT Gateway       |
+---------------------+
    |
+---[Private Subnet: App Tier]---+
|   EKS Worker Nodes           |
|   - Payment Service Pods     |
|   - Auth Service Pods        |
|   - Admin Service Pods       |
|   - Fraud Service Pods       |
+-------------------------------+
    |
+---[Private Subnet: Data Tier]---+
|   RDS PostgreSQL (Multi-AZ)     |
|   ElastiCache Redis Cluster     |
|   MSK Kafka Brokers             |
|   OpenSearch Cluster            |
+-----------------------------------+
    |
+---[Private Subnet: Storage]---+
|   S3 VPC Endpoint             |
+--------------------------------+
```

## Service Communication Patterns

| Pattern | Protocols | Use Cases |
|---------|-----------|-----------|
| Synchronous | HTTP/REST, gRPC | API requests, internal service calls |
| Asynchronous | Kafka | Payment events, notifications, analytics |
| Batch | Scheduled jobs | Settlement, reconciliation, report generation |

## Data Flow: Payment Lifecycle

```
1. Merchant POST /v1/payments  ──HTTP──> API Gateway
2. API Gateway ──> Auth Service (validate API key, rate limit)
3. API Gateway ──> Payment Service (create payment intent)
4. Payment Service ──> Fraud Service (risk check sync)
5. Payment Service ──> Kafka (payment.created event)
6. Payment Service <returns redirect URL to merchant>
7. Customer submits card ──> Payment Service (from merchant)
8. Payment Service ──> Processor Gateway (authorize)
9. Processor Gateway <--> Card Network (Visa/MC)
10. Payment Service ──> Kafka (payment.authorized / failed)
11. [Later] Merchant POST /v1/payments/{id}/capture
12. Payment Service ──> Kafka (payment.captured)
13. Payment Service ──> Ledger Service (gRPC: create entries)
14. Ledger Service ──> Kafka (ledger.entry.created)
15. Settlement Service (batch, daily) ──> process captures
16. Settlement Service ──> Payout to merchant bank account
```
