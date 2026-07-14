# Caching Strategy (Redis)

## Redis Cluster Topology

```
           ┌───────────────────┐
           │  Redis Cluster    │
           │  3 masters        │
           │  3 replicas       │
           │  (ElastiCache)    │
           └────────┬──────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
   ┌────▼────┐ ┌────▼────┐ ┌────▼────┐
   │Shard 1  │ │Shard 2  │ │Shard 3  │
   │m6g.large│ │m6g.large│ │m6g.large│
   └─────────┘ └─────────┘ └─────────┘
```

## Cache Usage Inventory

| Use Case | Key Pattern | TTL | Data Type | Eviction | Priority |
|----------|-------------|-----|-----------|----------|----------|
| API Rate Limiter | `ratelimit:{api_key}:{endpoint}` | 1s–60s | Sorted Set (sliding window) | No eviction (volatile) | Critical |
| Idempotency | `idempotency:{key}` | 24h | String (response JSON) | No eviction | Critical |
| Session Store | `session:{token}` | 7d | Hash (user_id, role, expiry) | allkeys-lru | Critical |
| JWT Blacklist | `jwt_blacklist:{jti}` | Until JWT expiry | String (revoked_at) | No eviction | Critical |
| Merchant Config | `merchant:{id}:config` | 5min | Hash (fee rates, webhook URLs) | allkeys-lru | High |
| Rate Limit Config | `ratelimit_config:{api_key}` | 5min | Hash (max_reqs, window) | allkeys-lru | High |
| Processor Auth Token | `processor:{name}:token` | Token expiry - 60s | String (access_token) | allkeys-lru | High |
| Fraud: IP Reputation | `fraud:ip:{address}` | 1h | Hash (score, country, isp) | volatile-lru | Medium |
| Fraud: BIN Data | `fraud:bin:{bin}` | 24h | Hash (bank, country, type) | allkeys-lru | Medium |
| Dashboard Aggregates | `dashboard:{merchant}:{metric}` | 5min | String (JSON serialized) | allkeys-lfu | Medium |
| OTP Store | `otp:{phone_or_email}` | 5min | String (code, attempts) | No eviction | High |
| Distributed Locks | `lock:{resource}` | 30s (auto-extend) | String (lock_holder) | No eviction | Critical |
| Webhook Rate Limit | `webhook_ratelimit:{merchant_id}` | 1s | Counter | No eviction | High |
| Payment Status Cache | `payment:{id}:status` | 30s | String (status JSON) | allkeys-lru | Medium |

## Cache Patterns

### 1. Rate Limiting (Sliding Window)
```python
# For each request:
key = f"ratelimit:{api_key}:{endpoint}:{current_minute}"
current = INCR(key)
EXPIRE(key, 61)  # 1 min + 1s buffer

if current > max_allowed:
    return 429 Too Many Requests
```

### 2. Idempotency
```python
key = f"idempotency:{idempotency_key}"
result = GET(key)
if result:
    return result  # Return cached response

# Process the request...
# On success:
SET(key, response_json, EX=86400)
```

### 3. Distributed Lock (Redlock)
```python
# Ensure only one instance processes a batch job
lock_key = f"lock:settlement:{date}"
lock_value = uuid4()
acquired = SET(lock_key, lock_value, NX=True, PX=30000)
if not acquired:
    return "Another instance is processing"

try:
    process_settlement(date)
finally:
    # Release only if we own the lock
    if GET(lock_key) == lock_value:
        DEL(lock_key)
```

## Failure Handling
- **Redis down:** Services degrade gracefully
  - Rate limiting falls back to local in-memory counter (approximate)
  - Idempotency falls back to database query (slower but safe)
  - Sessions fall back to database lookup
  - Locks: operations that require locks check database instead
- **Read failures:** Return stale cache data if available, else query DB
- **Write failures:** Log error, continue without cache (DB is source of truth)

## Memory Sizing
- Estimated: 10GB total (sessions + rate limiter + idempotency + fraud data)
- ElastiCache m6g.large nodes: 3 shards × 8GB = 24GB total (with headroom)
- Maxmemory policy: `allkeys-lru` for cache data
- Critical keys (idempotency, locks) use `NOEVICTION` via separate Redis instance or on same with careful key management

## Monitoring
- **Cache hit ratio:** Target >90% for session, >95% for idempotency
- **Eviction rate:** Alert if >1000 keys/sec evicted
- **Memory usage:** Alert at >80% of maxmemory
- **Latency:** Alert if p99 > 5ms
