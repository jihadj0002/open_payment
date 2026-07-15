# Webhook Backlog Runbook

## Symptoms
- Webhook worker lagging
- `webhook_deliveries` table growing rapidly
- Merchant complaints about missing webhooks
- Monitor: Pending deliveries count > 1000

## Steps

### 1. Assess the Backlog

```sql
-- Check backlog size
SELECT status, COUNT(*) FROM webhook_deliveries
WHERE created_at > NOW() - INTERVAL '1 hour'
GROUP BY status;

-- Check oldest pending delivery
SELECT MIN(created_at) FROM webhook_deliveries WHERE status = 'pending';
```

### 2. Identify Problematic Endpoints

```sql
-- Find endpoints with high failure rates
SELECT w.id, w.url, COUNT(d.id) as total,
  SUM(CASE WHEN d.status = 'failed' THEN 1 ELSE 0 END) as failures
FROM webhook_deliveries d
JOIN webhooks w ON w.id = d.webhook_id
WHERE d.created_at > NOW() - INTERVAL '1 hour'
GROUP BY w.id, w.url
ORDER BY failures DESC;
```

### 3. Rate Limiting

Slow down delivery to failing endpoints:
```bash
# Set webhook rate limit
export WEBHOOK_RATE_LIMIT=10  # deliveries per minute
```

The retry worker will automatically back off using exponential delays:
- 1st retry: 5 seconds
- 2nd retry: 30 seconds
- 3rd retry: 5 minutes
- 4th retry: 30 minutes
- 5th retry: max attempts reached → marked as failed

### 4. Manual Retry

If automatic retry is insufficient:
```bash
# Trigger manual retry for all pending deliveries
# (Restart webhook worker with reduced interval)
docker-compose exec api kill -HUP 1  # triggers graceful restart
```

### 5. Clear Stuck Deliveries

```sql
-- Mark deliveries stuck for > 24 hours as failed
UPDATE webhook_deliveries
SET status = 'failed'
WHERE status = 'pending'
  AND created_at < NOW() - INTERVAL '24 hours'
  AND attempt >= max_attempts;
```

### 6. Merchant Notification

Notify merchants with failing endpoints via the `health_alert` webhook event.
