# Data Flow Diagrams

## 1. Standard Card Payment Flow

```
Customer              Merchant            Gateway           Processor          Card Network
   │                     │                  │                  │                   │
   │  1. Browse & add    │                  │                  │                   │
   │────────────────────▶│                  │                  │                   │
   │                     │                  │                  │                   │
   │  2. Checkout        │                  │                  │                   │
   │────────────────────▶│                  │                  │                   │
   │                     │  3. POST /v1/payments               │                   │
   │                     │─────────────────▶│                  │                   │
   │                     │                  │  4. Auth check   │                   │
   │                     │                  │  5. Fraud check  │                   │
   │                     │                  │  6. Create PI    │                   │
   │                     │<─── return url ──│                  │                   │
   │                     │                  │                  │                   │
   │  7. Redirect to     │                  │                  │                   │
   │  checkout page      │                  │                  │                   │
   │◀────────────────────│                  │                  │                   │
   │                     │                  │                  │                   │
   │  8. Enter card      │                  │                  │                   │
   │  details + submit   │                  │                  │                   │
   │─────────────────────┼─────────────────▶│                  │                   │
   │                     │                  │  9. Authorize    │                   │
   │                     │                  │─────────────────▶│                  │
   │                     │                  │                  │  10. Auth req      │
   │                     │                  │                  │──────────────────▶│
   │                     │                  │                  │                   │
   │                     │                  │                  │  11. Auth resp     │
   │                     │                  │                  │◀──────────────────│
   │                     │                  │  12. Auth result │                   │
   │                     │                  │◀─────────────────│                   │
   │                     │                  │                  │                   │
   │  13. Show success   │                  │  14. Kafka:      │                   │
   │◀────────────────────│                  │  payment.auth    │                   │
   │                     │                  │  15. Webhook:    │                   │
   │                     │                  │  payment.success │                   │
   │                     │◀── webhook ──────│                  │                   │
   │                     │                  │                  │                   │
```

## 2. Capture and Settlement Flow

```
Merchant              Gateway            Ledger             Settlement          Bank
   │                     │                  │                  │                  │
   │  1. POST /capture   │                  │                  │                  │
   │────────────────────▶│                  │                  │                  │
   │                     │  2. Process      │                  │                  │
   │                     │     capture      │                  │                  │
   │                     │  3. gRPC: create │                  │                  │
   │                     │     debit/credit │                  │                  │
   │                     │─────────────────▶│                  │                  │
   │                     │                  │  4. Update       │                  │
   │                     │                  │     balances     │                  │
   │                     │◀── OK ───────────│                  │                  │
   │                     │                  │                  │                  │
   │  5. Return success  │                  │                  │                  │
   │◀────────────────────│                  │                  │                  │
   │                     │                  │                  │                  │
   │                     │  6. Kafka:       │                  │                  │
   │                     │  payment.captured│                  │                  │
   │                     │──────────────────┼─────────────────▶│                  │
   │                     │                  │                  │                  │
   │                     │                  │   [Daily Batch]  │                  │
   │                     │                  │  7. Aggregate    │                  │
   │                     │                  │     all captures │                  │
   │                     │                  │  8. Calculate    │                  │
   │                     │                  │     fees + tax   │                  │
   │                     │                  │  9. Create       │                  │
   │                     │                  │     payout batch │                  │
   │                     │                  │  10. Initiate    │                  │
   │                     │                  │     bank tx      │─────────────────▶│
   │                     │                  │                  │  11. Confirm     │
   │                     │                  │                  │◀─────────────────│
   │                     │                  │  12. Mark settled│                  │
   │                     │                  │  13. Webhook:    │                  │
   │                     │                  │     settlement   │                  │
   │                     │◀── webhook ──────│                  │                  │
```

## 3. Refund Flow

```
Customer              Merchant            Gateway            Ledger            Processor
   │                     │                  │                  │                  │
   │  1. Request refund  │                  │                  │                  │
   │────────────────────▶│                  │                  │                  │
   │                     │  2. POST /refund │                  │                  │
   │                     │─────────────────▶│                  │                  │
   │                     │                  │  3. Validate     │                  │
   │                     │                  │     payment state│                  │
   │                     │                  │     (must be     │                  │
   │                     │                  │     captured)    │                  │
   │                     │                  │  4. Send refund  │                  │
   │                     │                  │     to processor──────────────────▶│
   │                     │                  │                  │  5. Process     │
   │                     │                  │                  │     refund      │
   │                     │                  │                  │◀── confirm ─────│
   │                     │                  │  6. gRPC: create │                  │
   │                     │                  │     ledger entry──────────────────▶│
   │                     │                  │                  │                  │
   │                     │                  │  7. Kafka:       │                  │
   │                     │                  │     payment.refunded               │
   │  8. Webhook:        │                  │                  │                  │
   │  refund.completed   │◀── webhook ──────│                  │                  │
   │◀────────────────────│                  │                  │                  │
   │  9. Notify customer │                  │                  │                  │
   │────────────────────▶│                  │                  │                  │
```

## 4. Chargeback Flow

```
Card Network          Gateway            Ledger            Merchant            Admin
   │                     │                  │                  │                  │
   │  1. Chargeback      │                  │                  │                  │
   │     notification    │                  │                  │                  │
   │────────────────────▶│                  │                  │                  │
   │                     │  2. Create       │                  │                  │
   │                     │     dispute      │                  │                  │
   │                     │  3. Reserve      │                  │                  │
   │                     │     amount       │                  │                  │
   │                     │─────────────────▶│                  │                  │
   │                     │  4. Webhook:     │                  │                  │
   │                     │  chargeback      │                  │                  │
   │                     │  .created        │─────────────────▶│                  │
   │                     │                  │                  │                  │
   │                     │                  │  5. Submit       │                  │
   │                     │                  │     evidence     │                  │
   │                     │                  │◀─────────────────│                  │
   │                     │  6. Forward      │                  │                  │
   │                     │     evidence     │                  │                  │
   │                     │     to network   │                  │                  │
   │◀────────────────────│                  │                  │                  │
   │                     │                  │                  │                  │
   │  7. Result          │                  │                  │                  │
   │  (won/lost)         │                  │                  │                  │
   │────────────────────▶│  8. Release/     │                  │                  │
   │                     │     deduct       │                  │                  │
   │                     │     reserve      │─────────────────▶│                  │
   │                     │  9. Webhook      │                  │                  │
   │                     │  dispute.        │─────────────────▶│                  │
   │                     │  resolved        │                  │                  │
```

## 5. Authentication Flow

```
Client               API Gateway         Auth Service         Redis              DB
  │                       │                  │                  │                 │
  │ 1. POST /auth/login   │                  │                  │                 │
  │──────────────────────▶│─────────────────▶│                  │                 │
  │                       │                  │ 2. Validate      │                 │
  │                       │                  │    credentials   │────────────────▶│
  │                       │                  │ 3. Check MFA     │                 │
  │                       │                  │    (if enabled)  │                 │
  │                       │                  │ 4. Generate JWT  │                 │
  │                       │                  │    (access +     │                 │
  │                       │                  │     refresh)     │                 │
  │                       │                  │ 5. Store session─────────────────▶│
  │                       │                  │ 6. Store refresh │                 │
  │                       │                  │─────────────────▶│                 │
  │                       │◀── tokens ───────│                  │                 │
  │◀──── tokens ──────────│                  │                  │                 │
  │                       │                  │                  │                 │
  │ 7. GET /api/payments  │                  │                  │                 │
  │ (Authorization: Bearer│                  │                  │                 │
  │  <JWT>)               │                  │                  │                 │
  │──────────────────────▶│ 8. Validate JWT  │                  │                 │
  │                       │─────────────────▶│                  │                 │
  │                       │                  │ 9. Check        │                 │
  │                       │                  │    blacklist     │─────────────────▶│
  │                       │                  │ 10. Verify role  │                 │
  │                       │◀── OK ───────────│                  │                 │
  │                       │ 11. Route to     │                  │                 │
  │                       │     Merchant Svc │                  │                 │
  │◀──── response ────────│                  │                  │                 │
```
