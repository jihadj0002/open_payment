# Scaling Strategy

## Scaling Dimensions

| Dimension | Current Capacity | Target | Strategy |
|-----------|-----------------|--------|----------|
| API requests | 10,000 req/s | 50,000 req/s | Horizontal pod scaling |
| Transaction throughput | 1,000 TPS | 5,000 TPS | Kafka partitioning + DB scaling |
| Data volume | 10M txns/month | 100M txns/month | Partitioning + archival |
| Merchants | 1,000 | 10,000 | Read replicas for dashboard |
| Concurrent users | 500 dashboard | 5,000 dashboard | CDN + caching + read replicas |
| Webhook delivery | 10,000/hour | 100,000/hour | Dedicated webhook workers |

## Horizontal Pod Autoscaling (HPA)

### Service Autoscaling Configuration
```yaml
# Payment Service
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: payment-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: payment-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Pods
    pods:
      metric:
        name: payment_processing_duration_seconds
      target:
        type: AverageValue
        averageValue: 500  # ms
```

### Autoscaling Triggers per Service

| Service | Min | Max | CPU Trigger | Custom Metric Trigger |
|---------|-----|-----|-------------|----------------------|
| API Gateway | 3 | 10 | 70% | Request rate > 1000/s |
| Payment | 3 | 20 | 70% | Payment processing > 500ms |
| Merchant | 2 | 6 | 70% | — |
| Auth | 2 | 6 | 70% | Login rate > 100/min |
| Webhook | 2 | 8 | 60% | Webhook queue depth > 100 |
| Fraud | 2 | 6 | 70% | — |
| Ledger | 2 | 6 | 70% | Queue lag > 1000 |
| Notification | 1 | 3 | 70% | — |

## Database Scaling

### Read Replicas
```sql
-- Application-level read/write splitting
type DBClient struct {
    writer *sql.DB    // Primary (writes + reads needing consistency)
    reader *sql.DB    // Replica (read-only, eventually consistent)
}

func (c *DBClient) GetPayment(ctx context.Context, id string) (*Payment, error) {
    // Payment creation/update reads from writer (strong consistency)
    // Dashboard/list reads from reader (eventual consistency)
    return c.reader.GetContext(ctx, "SELECT * FROM payment_intents WHERE id = $1", id)
}
```

### Connection Pool Sizing
| Service | Max Connections | Pool Size | Notes |
|---------|----------------|------------|-------|
| Payment | 25 | 10 | High throughput |
| Merchant | 10 | 5 | Low throughput |
| Ledger | 15 | 8 | Financial integrity |
| Webhook | 20 | 10 | Batch reading |
| Admin | 10 | 5 | Dashboard queries |

### Scaling Trigger Matrix

| Metric | Low Load | Medium Load | High Load | Action |
|--------|----------|-------------|-----------|--------|
| CPU < 50% | ✅ | — | — | Reduce pods |
| CPU > 70% | — | ⚠️ | — | Add pods |
| CPU > 85% | — | — | 🚨 | Add pods + scale DB |
| Queue lag < 100 | ✅ | — | — | Reduce consumers |
| Queue lag > 1000 | — | ⚠️ | — | Add consumers |
| Queue lag > 10000 | — | — | 🚨 | Scale Kafka partitions |
| DB connections < 50 | ✅ | — | — | Normal |
| DB connections > 100 | — | ⚠️ | — | Scale readers |
| DB CPU > 80% | — | — | 🚨 | Add replicas or shard |

## Database Sharding (When Needed)

### Sharding Key: merchant_id
```sql
-- Shard allocation table
CREATE TABLE shard_map (
    merchant_id UUID PRIMARY KEY,
    shard_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Application-level shard routing
func getShardForMerchant(merchantID uuid.UUID) *sql.DB {
    shardID := hash(merchantID) % numShards
    return shards[shardID]
}
```

### Sharding Trigger Conditions
- Database CPU consistently > 80%
- Write throughput > 5,000 TPS
- Total data > 5TB
- Connection pool exhausted

## Caching Strategy for Scale

### Multi-Layer Cache
```
Request
  │
  ▼
┌──────────────┐
│  CDN Cache   │  (Cloudflare: static assets, API responses for GET with Cache-Control)
│  TTL: 5-60s  │
└──────┬───────┘
       │ (miss)
       ▼
┌──────────────┐
│  Redis Cache  │  (Application-level: sessions, idempotency, rate limits, merchant config)
│  TTL: 5s-24h  │
└──────┬───────┘
       │ (miss)
       ▼
┌──────────────┐
│  PostgreSQL   │  (Source of truth)
└──────────────┘
```

### Cache Warming
```go
// Warm cache for frequently accessed data
func warmMerchantCache(merchantID uuid.UUID) {
    config, _ := db.GetMerchantConfig(merchantID)
    redis.Set("merchant:"+merchantID+":config", config, 5*time.Minute)

    rateLimit, _ := db.GetRateLimitConfig(merchantID)
    redis.Set("merchant:"+merchantID+":rate_limit", rateLimit, 5*time.Minute)
}
```

## Kafka Scaling

### Partition Count Scaling
- Initial: 6 partitions per topic
- Scale trigger: consumer lag > 10,000 for > 5 minutes
- Partition count formula: `max(6, max_consumers * 2)`
- Re-partitioning requires topic recreation or new topic with data migration

### Consumer Group Scaling
```yaml
# Add more consumers within the same consumer group
# Each partition is consumed by exactly one consumer
# More partitions needed for more parallelism
webhook-service:
  concurrency: 10  # goroutines per consumer
  maxPollRecords: 100
  fetchMaxBytes: 52428800  # 50MB
```

## Cost Optimization at Scale

| Strategy | Savings | Effort |
|----------|---------|--------|
| Spot instances (batch + non-critical) | 60-70% | Low |
| Reserved instances (database + data nodes) | 30-40% | Low |
| Auto-scaling down to zero (batch, non-prod) | Variable | Low |
| S3 lifecycle (archive old data) | 80% on storage | Low |
| Right-sizing (Goldilocks recommendations) | 20-30% | Medium |
| Read replica count adjustment | 10-20% | Low |
| Kafka retention reduction (older → S3) | 40% on Kafka | Medium |
