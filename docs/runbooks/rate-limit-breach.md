# Rate Limit Breach Runbook

## Symptoms
- 429 Too Many Requests in merchant responses
- Rate limit alerts firing
- Abusive client consuming API resources

## Steps

### 1. Identify the Abusive Client

```bash
# Check rate limit logs
docker-compose logs api | grep "rate limit"

# Find top IPs by request count (requires api_usage_logs)
docker-compose exec postgres psql -U postgres -d openpayment -c "
  SELECT ip_address, COUNT(*) as requests
  FROM api_usage_logs
  WHERE created_at > NOW() - INTERVAL '5 minutes'
  GROUP BY ip_address
  ORDER BY requests DESC
  LIMIT 10;
"
```

### 2. Check Merchant-Specific Usage

```sql
SELECT merchant_id, COUNT(*) as requests
FROM api_usage_logs
WHERE created_at > NOW() - INTERVAL '5 minutes'
GROUP BY merchant_id
ORDER BY requests DESC
LIMIT 10;
```

### 3. Throttle the Abusive Client

**Option A: Block IP at application level**
```bash
# Add to blocklist via admin API
curl -X POST localhost:8080/v1/admin/rate-limits/block \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"ip": "203.0.113.1"}'
```

**Option B: Reduce rate limit for specific merchant**
```bash
# Override merchant rate limit
curl -X PATCH localhost:8080/v1/admin/merchants/{id} \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"rate_limit": 10}'  # 10 requests per minute
```

**Option C: Suspend merchant temporarily**
```bash
curl -X POST localhost:8080/v1/admin/merchants/{id}/suspend \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"reason": "Rate limit abuse"}'
```

### 4. Verify Mitigation

```bash
# Test rate limit headers
curl -I http://localhost:8080/v1/health \
  -H "Authorization: Bearer $TEST_TOKEN"
# Check X-RateLimit-Remaining header
```

### 5. Incident Report

Document the incident with:
- Client IP/merchant ID
- Request rate (req/s)
- Affected endpoints
- Action taken
- Duration of abuse

### 6. Prevention

- Review rate limit thresholds (default: 100 req/s per merchant)
- Consider adding CAPTCHA for suspicious patterns
- Implement IP-based rate limiting for unauthenticated endpoints
- Add alerting for merchants exceeding 80% of their rate limit
