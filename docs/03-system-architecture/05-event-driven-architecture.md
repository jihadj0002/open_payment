# Event-Driven Architecture (Kafka)

## Topic Catalog

| Topic Name | Partitions | Retention | Compaction | Key | Schema |
|------------|------------|-----------|------------|-----|--------|
| `payment.created` | 6 | 7 days | No | payment_id | Avro |
| `payment.authorized` | 6 | 7 days | No | payment_id | Avro |
| `payment.captured` | 6 | 30 days | No | payment_id | Avro |
| `payment.failed` | 6 | 30 days | No | payment_id | Avro |
| `payment.refunded` | 6 | 30 days | No | payment_id | Avro |
| `payment.voided` | 3 | 7 days | No | payment_id | Avro |
| `payment.chargeback` | 3 | 90 days | No | payment_id | Avro |
| `payment.dispute.resolved` | 3 | 90 days | No | dispute_id | Avro |
| `ledger.entry.created` | 6 | 30 days | No | account_id | Avro |
| `settlement.completed` | 3 | 30 days | No | merchant_id | Avro |
| `payout.initiated` | 3 | 30 days | No | merchant_id | Avro |
| `merchant.updated` | 3 | 7 days | Yes (key=merchant_id) | merchant_id | Avro |
| `fraud.alert` | 3 | 90 days | No | payment_id | Avro |
| `webhook.delivery.attempt` | 6 | 7 days | No | webhook_id | Avro |
| `notification.send` | 3 | 3 days | No | notification_id | Avro |
| `dead.letter` | 3 | 30 days | No | original_topic | Avro |

## Event Schema (Avro) Examples

### payment.created
```json
{
  "type": "record",
  "name": "PaymentCreated",
  "fields": [
    { "name": "event_id", "type": "string" },
    { "name": "event_type", "type": "string", "default": "payment.created" },
    { "name": "timestamp", "type": "long" },
    { "name": "payment_id", "type": "string" },
    { "name": "merchant_id", "type": "string" },
    { "name": "amount", "type": "long" },
    { "name": "currency", "type": "string" },
    { "name": "payment_method", "type": "string" },
    { "name": "idempotency_key", "type": "string" }
  ]
}
```

### payment.authorized
```json
{
  "type": "record",
  "name": "PaymentAuthorized",
  "fields": [
    { "name": "event_id", "type": "string" },
    { "name": "event_type", "type": "string", "default": "payment.authorized" },
    { "name": "timestamp", "type": "long" },
    { "name": "payment_id", "type": "string" },
    { "name": "processor_transaction_id", "type": "string" },
    { "name": "processor_response_code", "type": "string" },
    { "name": "authorization_code", "type": "string" },
    { "name": "amount", "type": "long" },
    { "name": "currency", "type": "string" }
  ]
}
```

## Consumer Groups

| Consumer Group | Topics Consumed | Purpose |
|----------------|-----------------|---------|
| `webhook-service` | payment.*, settlement.*, dispute.* | Deliver webhooks to merchants |
| `notification-service` | payment.* | Send email/SMS notifications |
| `ledger-service` | payment.captured, payment.refunded | Create ledger entries |
| `settlement-service` | payment.captured | Queue for settlement batch |
| `analytics-service` | payment.*, ledger.* | Aggregate for reports |
| `fraud-service` | payment.created | Real-time fraud scoring |
| `audit-service` | payment.*, merchant.* | Immutable audit trail |

## Delivery Semantics
- **Producer:** `acks=all`, `enable.idempotence=true` (exactly-once per partition)
- **Consumer:** `enable.auto.commit=false`, manual commit after processing
- **Retries:** 3 retries with exponential backoff, then dead letter queue

## Dead Letter Queue
Events that fail processing after max retries are sent to `dead.letter` topic with:
- Original event payload
- Error message and stack trace
- Number of retry attempts
- Consumer group that failed

## Schema Registry
- Confluent Schema Registry (or open-source alternative)
- Schema evolution: `BACKWARD` compatibility (can delete fields, add optional fields)
- All services validate schema on produce/consume

## Event Ordering Guarantees
- **Per-payment ordering:** Events for the same `payment_id` land in same partition (keyed by `payment_id`), preserving order
- **Cross-service ordering:** Not guaranteed — services must handle out-of-order events idempotently
