
A payment gateway has to solve four major problems:

Process payments reliably.
Protect customer money and data.
Integrate with merchants easily.
Comply with financial regulations.

Below is a software engineering perspective of what a complete payment gateway requires.

1. High Level System Architecture
                Customer
                    │
          Merchant Website/App
                    │
      Merchant SDK/API (REST/GraphQL)
                    │
          API Gateway + WAF
                    │
────────────────────────────────────────
        Authentication Service
────────────────────────────────────────
                    │
         Payment Orchestrator
                    │
     ┌──────────────┼──────────────┐
     │              │              │
 Card Processor   Wallets      Bank APIs
     │              │              │
 Visa           bKash         Bank Transfer
 Mastercard     Nagad         ACH
 Amex           PayPal
────────────────────────────────────────
                    │
        Settlement Engine
                    │
      Ledger & Accounting Service
                    │
      Merchant Wallet/Balance
                    │
        Payout Service
                    │
 Merchant Bank Account
2. Core Services

Instead of one monolithic backend, divide the system into services.

API Gateway

Responsible for

Rate limiting
Authentication
API Keys
JWT validation
IP filtering
Load balancing
Request logging
WAF
Merchant Service

Stores

Merchant

id

business_name

email

phone

status

verification_status

api_key

secret_key

webhook_secret

currency

country

timezone

created_at

Features

Merchant onboarding

Business verification

KYC

Dashboard

API management

Webhook management

Roles

Multiple users

API usage

Customer Service

Stores

Customer Profiles

Saved Cards

Payment History

Tokens

Device Fingerprints

Risk Score

Payment Service

Handles

Create Payment

Capture

Cancel

Refund

Partial Refund

Void

Split Payment

Installment

Recurring

Subscription

Transaction Service

Keeps immutable transaction history.

Never modify rows.

Only append.

Example

Payment Created

Authorized

Captured

Refunded

Failed

Expired

Chargeback

Disputed
Settlement Service

Calculates

Merchant Earnings

Gateway Fee

Tax

Commission

Reserve Balance

Settlement Date

Bank Transfer

Ledger Service

This is one of the most important services.

Use Double Entry Accounting.

Example

Customer Wallet

Merchant Wallet

Gateway Wallet

Bank Wallet

Reserve Wallet

Each transaction creates

Debit

Credit

No direct balance editing.

Fraud Detection Service

Checks

Velocity attacks

Card testing

IP reputation

VPN

TOR

Country mismatch

Impossible travel

BIN country mismatch

Device fingerprint

Bot detection

Email risk

Disposable email

Chargeback history

ML Risk Score

Notification Service

Email

SMS

Push

Webhook

Slack

Discord

Reporting Service

Merchant Reports

Revenue

Refunds

Settlement

Tax

Fees

CSV Export

PDF

Analytics

Webhook Service

Merchants receive

payment.success

payment.failed

payment.pending

refund.completed

chargeback.created

subscription.renewed

Webhook retries

Signature verification

Idempotency

Logs

3. Payment Flow

Customer

↓

Merchant creates payment

↓

Gateway validates merchant

↓

Create Payment Intent

↓

Risk Analysis

↓

3D Secure

↓

Authorization

↓

Capture

↓

Settlement

↓

Merchant Balance

↓

Payout

4. Database Design
Merchant
Merchant
MerchantUser
MerchantSettings
MerchantAPIKeys
MerchantWebhook
MerchantFees
MerchantBalance
MerchantSettlement
Payment
PaymentIntent
Transaction
Refund
Dispute
Chargeback
PaymentMethod
SavedCard
Accounting
LedgerAccount
LedgerEntry
Balance
Settlement
Invoice
Fee
Tax
Security
AuditLog

SecurityLog

LoginHistory

APIKeyHistory

Device

Session

Permission
5. API Design

Merchant API

POST /payments

POST /refund

POST /capture

POST /void

GET /payment/{id}

GET /transactions

GET /balance

POST /webhooks

POST /tokenize

POST /customer

POST /subscription
6. Merchant Dashboard

Overview

Revenue

Today's Revenue

Monthly Revenue

Pending

Refunds

Chargebacks

Settlement

Balance

API Keys

Webhook

Developers

Logs

Transactions

Customers

Analytics

Reports

Settings

Roles

Users

Support

7. Admin Dashboard

Merchant Approval

KYC

Freeze Merchant

Refund Management

Dispute Resolution

Fraud Detection

Logs

Monitoring

Settlement

Reports

API Monitoring

Webhook Logs

Configuration

Currencies

Countries

Banks

Fees

Taxes

Support Tickets

Risk Dashboard

8. Security Architecture

Most important section.

Encryption

AES-256

TLS 1.3

RSA

ECDSA

HSM support

Secrets Manager

Never store plaintext secrets.

Authentication

OAuth2

JWT

API Keys

Refresh Tokens

MFA

Passkeys

Session Management

Authorization

RBAC

ABAC

Permission Matrix

Least Privilege

API Security

Rate Limiting

Nonce

Timestamp

Request Signature

HMAC

Replay Prevention

CSRF

CORS

Input Validation

Infrastructure Security

Private Network

Zero Trust

VPN

Firewall

WAF

DDoS Protection

IDS

IPS

Secrets Rotation

Logging

Every request

Every payment

Every refund

Every login

Every API Key

Every admin action

Immutable Audit Logs

9. Compliance

Without these, many regions won't allow you to operate.

PCI DSS Level 1

EMVCo

PSD2 (Europe)

Strong Customer Authentication (SCA)

SOC 2 Type II

ISO 27001

GDPR (Europe)

CCPA (California)

AML (Anti-Money Laundering)

KYC (Know Your Customer)

Sanctions screening

OFAC checks (where applicable)

Data retention policies

Privacy consent management

10. Fraud Prevention

Device Fingerprinting

Behavior Analysis

Velocity Checks

Machine Learning

Risk Scoring

Card BIN Database

Proxy Detection

Geo Detection

Email Reputation

Phone Reputation

Behavioral Biometrics

Transaction Limits

Blacklist

Whitelist

Chargeback Prediction

11. Reliability

Load Balancer

Auto Scaling

Redis Cache

Message Queue

Kafka

RabbitMQ

Retry Queue

Dead Letter Queue

Circuit Breaker

Bulkhead

Timeouts

Health Checks

Horizontal Scaling

Database Replication

Automatic Failover

Backups

Disaster Recovery

12. Monitoring

Prometheus

Grafana

OpenTelemetry

Distributed Tracing

ELK Stack

Sentry

PagerDuty

Metrics

Payment Success Rate

Latency

API Errors

Settlement Delay

Webhook Failure Rate

CPU

Memory

DB

Redis

Queue

13. Security Monitoring

Detect

Credential Stuffing

Brute Force

SQL Injection

XSS

SSRF

RCE

Bot Attack

DDoS

Anomalies

Impossible Login

Suspicious Merchant

Multiple Failed Cards

14. Performance Goals

API Response

<200 ms for most read requests

Payment Initialization

<500 ms

Webhook Delivery

<5 seconds (with retries)

Availability

99.99%+

RPO

<5 minutes

RTO

<30 minutes

15. Payment State Machine
Created

↓

Pending

↓

Authorized

↓

Captured

↓

Settled

↓

Paid Out

Other branches

Pending

↓

Failed

↓

Cancelled

↓

Expired

↓

Refunded

↓

Chargeback

↓

Disputed

Never allow invalid transitions (for example, refunding a payment that was never captured).

16. Technology Stack Example

Frontend

React or Next.js (merchant/admin dashboards)
Tailwind CSS
TypeScript

Backend

Go, Java (Spring Boot), Rust, or ASP.NET Core for core payment processing
Python or Node.js for auxiliary services (analytics, notifications)

API

REST with optional GraphQL for dashboard queries
OpenAPI/Swagger documentation

Databases

PostgreSQL (transactional data)
Redis (caching, rate limiting)
Elasticsearch/OpenSearch (search and logs)

Messaging

Kafka or RabbitMQ

Storage

S3-compatible object storage for documents and reports

Infrastructure

Kubernetes
Docker
NGINX or Envoy
Terraform
GitHub Actions or GitLab CI/CD

Observability

Prometheus
Grafana
OpenTelemetry
Loki or ELK
17. Development Roadmap

A practical way to build this is in phases:

MVP
Merchant onboarding
API keys
Payment intents
Transaction tracking
Webhooks
Basic dashboard
Production
Ledger
Refunds
Settlements
Payouts
Audit logging
RBAC
Monitoring
Basic fraud detection
Enterprise
Multi-currency
Multi-region
Split payments
Subscriptions
Tokenization
Disputes and chargebacks
Advanced fraud engine
High availability across regions
Global Scale
Payment routing across multiple processors
Smart retry logic
AI-assisted fraud detection
Marketplace support
Embedded finance features
Developer ecosystem (SDKs, plugins, sandbox environments)
A critical architectural principle

One of the biggest differences between a hobby payment system and a production gateway like Stripe or SSLCommerz is the internal ledger. Every movement of money should be represented as immutable, balanced ledger entries rather than updating account balances directly. Combined with idempotent APIs, append-only transaction logs, comprehensive audit trails, strong encryption, and rigorous compliance controls, this provides the foundation for correctness, traceability, and resilience at scale.

Building a gateway to this standard is a significant engineering effort involving backend, security, DevOps, compliance, risk, and finance expertise. The reward is a platform that can safely handle high transaction volumes while remaining auditable, secure, and reliable.