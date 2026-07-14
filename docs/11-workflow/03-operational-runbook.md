# Operational Runbook

## On-Call Responsibilities

### Primary On-Call
- **Rotation:** 1 week (Mon 9am → Mon 9am)
- **Hours:** 24/7 for SEV-1, business hours for others
- **Response times:** SEV-1 (5 min), SEV-2 (15 min), SEV-3 (1 hour)
- **Tools:** Laptop with internet, PagerDuty app, Slack, AWS console access

### Secondary On-Call
- **Rotation:** Same schedule as primary (offset)
- **Responsibilities:** Backup for SEV-1 escalation, assist during incidents

## Incident Response

### SEV-1: Service Down / Revenue Impact

#### Detection
- PagerDuty alert (phone call + push notification)
- Slack #incident channel ping

#### Triage (5 min)
```
1. ACKNOWLEDGE the alert in PagerDuty
2. JOIN the #incident Slack channel
3. ASSESS: Is this a real issue? (check Grafana, error logs)
4. DECLARE severity: "SEV-1 declared — Payment API 5xx > 5%"
5. ANNOUNCE in #incident: "Investigating payment API errors"
```

#### Common Scenarios

**Scenario A: Payment API returning 5xx**
```bash
# 1. Check if it's a recent deploy
kubectl rollout history deployment/payment-service -n payment

# 2. Rollback if recent (within 15 min)
kubectl rollout undo deployment/payment-service -n payment

# 3. Check pod logs
kubectl logs -n payment -l app=payment-service --tail=100

# 4. Check database
psql -h $DB_HOST -c "SELECT count(*), state FROM pg_stat_activity GROUP BY state;"

# 5. Check Redis
redis-cli -h $REDIS_HOST ping

# 6. Check Kafka consumer lag
kafka-consumer-groups --bootstrap-server $KAFKA_BROKER \
  --group payment-service --describe
```

**Scenario B: Database connection pool exhausted**
```bash
# 1. Check current connections
psql -h $DB_HOST -c "SELECT count(*) FROM pg_stat_activity;"

# 2. Find idle connections
psql -h $DB_HOST -c "SELECT pid, state, query FROM pg_stat_activity WHERE state = 'idle';"

# 3. Kill idle connections (if critical)
psql -h $DB_HOST -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE state = 'idle' AND pid <> pg_backend_pid();"

# 4. Scale PgBouncer or increase pool size
kubectl scale deployment/pgbouncer --replicas=3 -n data
```

**Scenario C: Kafka consumer lag**
```bash
# 1. Check lag
kafka-consumer-groups --bootstrap-server $KAFKA_BROKER \
  --group payment-service --describe

# 2. Check consumer pod health
kubectl logs -n payment -l app=kafka-consumer --tail=50

# 3. Restart consumer if stuck
kubectl rollout restart deployment/kafka-consumer -n payment
```

### SEV-2: Partial Degradation

#### Examples
- Webhook delivery delayed > 30 seconds
- Dashboard slow (> 5s load time)
- admin dashboard errors for specific merchants

#### Response
```
1. Acknowledge alert
2. Investigate within 15 minutes
3. Fix or mitigate within 1 hour
4. Post update to #incident every 30 minutes
```

### SEV-3: Minor Issue

#### Examples
- UI formatting bug
- Report export failing
- Stale dashboard data

#### Response
```
1. Create JIRA ticket
2. Fix during normal business hours
3. No page required
```

## Common Runbook Scripts

### Health Check
```bash
#!/bin/bash
# scripts/health-check.sh
ENV=${1:-production}

echo "=== Health Check: $ENV ==="

# API health
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" https://api.$ENV.openpayment.com/health)
echo "API Health: $HTTP_CODE"

# Check each service
for svc in payment merchant auth ledger webhook; do
    STATUS=$(kubectl get pods -n payment -l app=$svc-service -o jsonpath='{.items[*].status.phase}')
    echo "$svc-service: $STATUS"
done

# Database
DB_STATUS=$(psql -h $DB_HOST -c "SELECT 1" 2>&1)
echo "Database: ${DB_STATUS:0:50}..."

# Kafka
KAFKA_STATUS=$(kafka-broker-api-versions --bootstrap-server $KAFKA_BROKER 2>&1 | head -1)
echo "Kafka: $KAFKA_STATUS"
```

### Rollback Script
```bash
#!/bin/bash
# scripts/rollback.sh
SERVICE=${1:?Service name required}
REVISION=${2:-0}  # 0 = previous revision

kubectl rollout undo deployment/$SERVICE-service -n payment \
  --to-revision=$REVISION

echo "Rolling back $SERVICE to revision $REVISION..."
kubectl rollout status deployment/$SERVICE-service -n payment --timeout=5m
```

### Scale Service
```bash
#!/bin/bash
# scripts/scale.sh
SERVICE=${1:?Service name}
REPLICAS=${2:?Number of replicas}

kubectl scale deployment/$SERVICE-service -n payment --replicas=$REPLICAS
echo "Scaled $SERVICE to $REPLICAS replicas"
```

## Maintenance Procedures

### Database Migration
```bash
# 1. Backup first
pg_dump -h $DB_HOST -d paymentdb > /tmp/pre-migration-backup.sql

# 2. Apply migration (one connection)
kubectl exec -n payment deployment/payment-service -- \
  ./migrate -path /migrations -database $DATABASE_URL up 1

# 3. Verify
kubectl exec -n payment deployment/payment-service -- \
  ./migrate -path /migrations -database $DATABASE_URL version

# 4. If failed:
kubectl exec -n payment deployment/payment-service -- \
  ./migrate -path /migrations -database $DATABASE_URL down 1
```

### Certificate Rotation
```yaml
# TLS certificates are auto-renewed via cert-manager
# Verify:
kubectl get certificate -n system
kubectl describe certificate payment-api-tls -n system

# Manual renewal if needed:
kubectl delete certificate payment-api-tls -n system
# cert-manager will auto-recreate with new cert
```

### Secrets Rotation
```bash
# Rotate DB password
aws secretsmanager update-secret \
  --secret-id /payment/production/db-password \
  --secret-string "NewP@ssw0rd!"

# Update in Kubernetes (reload pods)
kubectl patch deployment payment-service -n payment \
  --patch '{"spec":{"template":{"metadata":{"annotations":{"secret-version":"2"}}}}}'
```

## Post-Mortem Template

```markdown
# Incident Post-Mortem: INC-XXX

## Summary
- **Date:** 2026-01-15
- **Duration:** 45 minutes (14:30 - 15:15 UTC)
- **Severity:** SEV-1
- **Impact:** Payment API 502 errors for 0.5% of requests
- **Detection:** PagerDuty alert (error rate > 1%)

## Timeline
| Time (UTC) | Event |
|------------|-------|
| 14:30 | Alert triggered: payment API error rate > 1% |
| 14:31 | Engineer acknowledged |
| 14:33 | Identified: DB connection pool exhausted |
| 14:35 | Executed: Kill idle connections |
| 14:38 | Error rate returned to normal |
| 14:40 | Root cause: Application connection leak |
| 15:00 | Hotfix deployed: close idle connections |
| 15:15 | Monitoring confirmed stable |

## Root Cause
The payment service was not closing idle database connections after long-running
reporting queries. Over time, connections accumulated and exhausted the pool.

## Resolution
- Killed idle connections (immediate mitigation)
- Deployed fix to properly close connections with `defer rows.Close()`
- Added connection pool monitoring alert

## Action Items
- [x] Add `defer rows.Close()` to payment repository queries
- [x] Add Grafana panel for DB connection pool usage
- [ ] Add unit test that simulates connection exhaustion
- [ ] Review all repository methods for proper connection cleanup

## Lessons Learned
1. Add connection pool monitoring sooner
2. Add integration tests that verify connection cleanup
3. Review all database code for connection leaks before next deploy
```

## On-Call Handoff Template
```markdown
# On-Call Handoff — Week 12

## Current Incidents
- INC-045: Dashboard report export failing (SEV-3, in progress — assignee: Frontend)
- INC-046: Webhook delivery delayed for merchant X (SEV-3, investigating)

## Active Alerts
- [ ] None — all alerts resolved

## Known Issues
- Staging DB has high CPU on Monday mornings (scheduled reporting)
- Merchant Y has high refund rate (likely fraud, compliance notified)

## Deployment Schedule
- Wednesday: Payment service v2.8.1 (bug fix)
- Friday: Merchant service v1.4.0 (new feature)

## Tips for Next Week
- The scheduled reporting job runs at 8am UTC daily — expect brief CPU spike
- Merchant Y is high-risk — any chargeback activity needs immediate review
- PagerDuty escalation: Primary → Secondary → Engineering Manager
```
