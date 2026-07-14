# Disaster Recovery

## Recovery Objectives

| Metric | Target | Measurement |
|--------|--------|-------------|
| RPO (Recovery Point Objective) | <5 minutes | Maximum data loss in an incident |
| RTO (Recovery Time Objective) | <30 minutes | Time to restore service |
| MTTD (Mean Time to Detect) | <5 minutes | Automated alerting |
| MTTR (Mean Time to Recover) | <1 hour | Time from alert to recovery |

## Backup Strategy

### PostgreSQL Backups
```yaml
# Backup schedule
schedule:
  - type: continuous_wal
    description: "WAL streaming to S3"
    tool: pg_receivewal / WAL-G
    retention: 7 days
    target: "s3://openpayment-backups-db/wal/"

  - type: full
    description: "Daily full backup"
    schedule: "0 3 * * *"  # Daily at 3 AM
    tool: pg_dump / WAL-G
    retention: 30 days
    target: "s3://openpayment-backups-db/full/"

  - type: monthly
    description: "Monthly snapshot for compliance"
    schedule: "0 4 1 * *"  # 1st of each month
    tool: pg_dump
    retention: 7 years
    target: "s3://openpayment-backups-db/monthly/"
```

### Point-in-Time Recovery (PITR)
```bash
#!/bin/bash
# scripts/restore-pitr.sh
TIMESTAMP=$1  # e.g., "2026-01-15 14:30:00 UTC"

# Restore base backup
wal-g backup-fetch /var/lib/postgresql/16/main LATEST

# Replay WAL to specific timestamp
wal-g wal-fetch --target-dir /var/lib/postgresql/16/main \
    --target-timestamp "$(date -d "$TIMESTAMP" +%s)"

# Start PostgreSQL
pg_ctl start -D /var/lib/postgresql/16/main
```

## High Availability Architecture

### Multi-AZ Deployment (Single Region)
```
Availability Zone A          Availability Zone B          Availability Zone C
    │                              │                              │
┌────┴─────┐                 ┌────┴─────┐                 ┌────┴─────┐
│ Primary  │                 │  Standby  │                 │  Standby  │
│ RDS PG   │◀───Sync Repl───▶│  (RO)     │◀───Async Repl──▶│  (RO)     │
└──────────┘                 └──────────┘                 └──────────┘
    │                              │
    ├──────────────────────────────┤
    │                              │
┌───┴────────┐              ┌──────┴──────┐
│ MSK Kafka  │              │ MSK Kafka   │
│ Broker 1   │              │ Broker 2    │
└────────────┘              └─────────────┘
```

### Failover Procedure

#### Database Failover
```bash
#!/bin/bash
# scripts/failover-db.sh

echo "1. Detecting primary failure..."

# Check if primary is reachable
if ! pg_isready -h $PRIMARY_HOST; then
    echo "Primary is down. Initiating failover..."

    # 2. Promote the most advanced replica
    aws rds promote-read-replica \
        --db-instance-identifier $REPLICA_IDENTIFIER

    # 3. Wait for promotion
    aws rds wait db-instance-available \
        --db-instance-identifier $REPLICA_IDENTIFIER

    # 4. Update DNS/connection string
    NEW_ENDPOINT=$(aws rds describe-db-instances \
        --db-instance-identifier $REPLICA_IDENTIFIER \
        --query 'DBInstances[0].Endpoint.Address' \
        --output text)

    # 5. Update service configs
    kubectl patch configmap payment-config \
        -n payment \
        --patch "{\"data\":{\"DB_HOST\": \"$NEW_ENDPOINT\"}}"

    # 6. Roll pods to pick up new config
    kubectl rollout restart deployment -n payment

    echo "Failover complete. New primary: $NEW_ENDPOINT"
fi
```

#### Kubernetes Failover
```yaml
# Pod Anti-Affinity ensures pods spread across AZs
spec:
  affinity:
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          labelSelector:
            matchExpressions:
            - key: app
              operator: In
              values:
              - payment-service
          topologyKey: topology.kubernetes.io/zone
```

## Disaster Scenarios

| Scenario | Impact | Recovery Procedure | RTO |
|----------|--------|-------------------|-----|
| Single pod crash | Minor | Kubernetes reschedules | <1 min |
| Worker node failure | Partial | Pods rescheduled | <2 min |
| AZ failure | Major | Cross-AZ failover | <5 min |
| Region failure | Critical | DR region activation | <30 min |
| Data corruption | Critical | PITR restore | <30 min |
| Ransomware | Critical | Isolate, restore from clean backup | <4 hours |
| Cloud provider outage | Critical | Multi-cloud contingency | Manual |

## DR Region Strategy

### Phase 1 (MVP)
- Single region (ap-south-1)
- Multi-AZ for resilience
- Backups to S3 (cross-region copy enabled)
- Manual DR runbook

### Phase 2 (Production)
- Same region, multiple AZs
- Automated failover for database
- Read replicas in different AZs
- DR drills quarterly

### Phase 3 (Enterprise)
- Active-passive multi-region
- Second region (ap-southeast-1) in standby
- Kafka cross-region replication (MirrorMaker)
- DNS-based traffic switching (Route53)
- Full DR automation

## DR Testing Schedule

| Test | Frequency | Type | Success Criteria |
|-----|-----------|------|------------------|
| Pod reschedule | Weekly | Automated | Pod restarts in <30s |
| Node drain | Monthly | Automated | Pods migrate gracefully |
| DB failover | Quarterly | Automated | RTO <5 min, RPO <1 min |
| AZ failure | Quarterly | Game day | Service continues in remaining AZs |
| PITR restore | Quarterly | Manual | Data restored to point in time |
| Full region failover | Annually | Game day | Full service in DR region <30 min |
