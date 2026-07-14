# Payment State Machine

## States and Transitions

```
                    ┌──────────┐
                    │  Created  │
                    └────┬─────┘
                         │
                    ┌────▼─────┐
              ┌────▶│  Pending  │◀────────────┐
              │     └────┬─────┘              │
              │          │                    │
              │     ┌────▼─────┐              │
              │     │Processing│              │
              │     └────┬─────┘              │
              │          │                    │
         ┌────┴────┐     │              ┌─────┴──────┐
         │  Failed │◀────┼──────────────│  Cancelled  │
         └─────────┘     │              └────────────┘
                         ▼
                    ┌───────────┐
                    │Authorized │
                    └─────┬─────┘
                          │
                    ┌─────▼──────┐
               ┌────│  Captured  │
               │    └─────┬──────┘
               │          │
               │    ┌─────▼──────┐
               │    │  Settled   │
               │    └─────┬──────┘
               │          │
          ┌────┴────┐ ┌───▼───────┐
          │ Refunded│ │  Paid Out │
          └─────────┘ └───────────┘

          ┌──────────┐     ┌───────────┐
          │Chargeback│────▶│  Disputed │
          └──────────┘     └─────┬─────┘
                          ┌──────┴──────┐
                          │  Won/Lost   │
                          └─────────────┘
```

## Transition Rules

| From | To | Allowed? | Conditions | Trigger |
|------|----|----------|------------|---------|
| Created | Pending | Yes | Valid request, passed initial validation | API POST /v1/payments |
| Created | Failed | Yes | Validation failure (bad card, invalid amount) | Validation error |
| Created | Cancelled | Yes | Merchant cancels before processing | API POST /v1/payments/{id}/cancel |
| Pending | Processing | Yes | Fraud check passed, forwarded to processor | Internal |
| Pending | Failed | Yes | Fraud check failed | Fraud engine |
| Pending | Cancelled | Yes | Merchant cancels during pending | API |
| Processing | Authorized | Yes | Processor approved | Processor response |
| Processing | Failed | Yes | Processor declined (insufficient funds, etc.) | Processor response |
| Authorized | Captured | Yes | Merchant captures within auth window (7 days) | API POST /v1/payments/{id}/capture |
| Authorized | Cancelled | Yes | Void authorization before capture | API POST /v1/payments/{id}/void |
| Authorized | Failed | Yes | Authorization expires (7 days) | Expiry scheduler |
| Captured | Settled | Yes | Settlement batch runs | Batch job (daily) |
| Captured | Refunded | Yes | Full refund processed | API POST /v1/payments/{id}/refund |
| Captured | Chargeback | Yes | Cardholder disputes with bank | Processor notification |
| Settled | Paid Out | Yes | Payout batch initiated | Batch job (daily/weekly) |
| Settled | Refunded | Yes | Refund after settlement (deducts from balance) | API (if within balance) |
| Settled | Chargeback | Yes | Chargeback received after settlement | Processor notification |
| Paid Out | Refunded | Yes | Refund after payout (negative on next payout) | API |
| Paid Out | Chargeback | Yes | Chargeback after payout | Processor notification |
| Refunded | — | Terminal | — | — |
| Chargeback | Disputed | Yes | Merchant contests chargeback | Merchant submits evidence |
| Disputed | Won | Yes | Card network rules in merchant's favor | Network response |
| Disputed | Lost | Yes | Card network rules against merchant | Network response |

## Invalid Transitions (Explicitly Forbidden)

| From | To | Reason |
|------|----|--------|
| Created | Captured | Must authorize first |
| Pending | Captured | Must authorize first |
| Authorized | Refunded | Must capture first |
| Refunded | Captured | Cannot re-capture a refunded payment |
| Cancelled | Authorized | Cannot revive a cancelled payment |
| Failed | Authorized | Transaction already failed |
| Captured | Cancelled | Already captured, can't void; use refund |
| Settled | Cancelled | Use refund instead |

## Idempotency Key Handling Per State

| Operation | Idempotency Behavior |
|-----------|---------------------|
| Create Payment | If same idempotency key, return existing Payment Intent (idempotent) |
| Capture | If same key, return success without re-processing |
| Refund | If same key, return existing refund without re-processing |
| Void | If same key, return success without re-processing |
| Webhook | Each webhook has unique ID; merchant handler should deduplicate |

## Expiry and Cleanup
- **Authorization expiry:** Auto-void after 7 days (card network limit)
- **Pending expiry:** Auto-cancel after 1 hour if no customer action
- **Unsettled capture:** Mark as flagged for admin review after 30 days
- **Dispute response window:** 21 days (Visa), 30 days (Mastercard)
