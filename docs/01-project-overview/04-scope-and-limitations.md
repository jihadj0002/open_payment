# Scope and Limitations

## In-Scope (v1.0 MVP)

### Payment Methods
- Credit/debit cards (Visa, Mastercard, Amex)
- Digital wallets (bKash, Nagad)
- Bank transfers (manual/ACH)

### Core Features
- Merchant onboarding and verification (KYC)
- API key generation and management
- Payment intent creation, authorization, capture
- Refunds (full and partial)
- Transaction history and search
- Webhook notifications (payment.success, payment.failed, refund.completed)
- Merchant dashboard (read-only transaction views, balance overview)
- Admin dashboard (merchant approval, basic dispute management)
- Sandbox environment for developers

### Security
- PCI DSS Level 1 compliance
- TLS 1.3 encryption in transit
- AES-256 encryption at rest
- Tokenization of card data
- API authentication (HMAC + API keys)
- Rate limiting and basic DDoS protection

### Infrastructure
- Kubernetes cluster (single region)
- PostgreSQL with read replicas
- Redis for caching and rate limiting
- Kafka for event streaming
- Prometheus + Grafana monitoring

## Out-of-Scope (Post-MVP)

### Payment Methods (v1.1+)
- Google Pay, Apple Pay
- PayPal
- UPI, Paytm
- Cryptocurrencies
- BNPL (Buy Now Pay Later)

### Features (v1.1+)
- Subscription/recurring billing engine
- Split payments and marketplace payouts
- Smart payment routing
- Installment plans
- Invoice generation and sending
- Multi-currency settlement
- Advanced fraud ML engine

### Infrastructure (v1.1+)
- Multi-region deployment
- Global CDN for webhook delivery
- Cross-region DR
- Sharding

## Known Limitations (v1.0)
- Single currency support (BDT initially)
- Single region deployment (AWS ap-south-1)
- Manual settlement cycles (daily batch)
- Basic fraud rules only (no ML in MVP)
- No mobile SDK — REST API only
- No GraphQL — REST only
- No real-time analytics dashboard (batch-updated)

## Assumptions
- PostgreSQL handles transactional load without sharding in v1
- Kafka cluster of 3 brokers sufficient for MVP throughput
- Single AWS region provides acceptable latency for target merchants
- Team of 4–6 engineers can deliver MVP in 12–16 weeks
- Merchants have basic technical ability to integrate REST APIs
