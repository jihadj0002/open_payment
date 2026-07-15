# Database Recovery Runbook

## Database Failover

### Symptoms
- API returns 503 Service Unavailable
- Health check shows database unhealthy
- Error logs: `connection refused` or `pq: SSL is not enabled`

### Steps

1. **Verify the outage**
   ```bash
   docker-compose exec postgres pg_isready -U postgres
   ```

2. **Check PostgreSQL logs**
   ```bash
   docker-compose logs postgres --tail=100
   ```

3. **Failover to replica if available**
   ```bash
   # Promote replica to primary
   docker-compose exec postgres-replica pg_ctl promote
   ```

4. **Update DATABASE_URL in environment to point to new primary**
   ```bash
   # Update .env or Railway env vars
   DATABASE_URL=postgres://user:pass@new-primary:5432/openpayment?sslmode=require
   ```

5. **Restart the API**
   ```bash
   docker-compose restart api
   ```

6. **Verify recovery**
   ```bash
   curl http://localhost:8080/health
   ```

## Restore from Backup

1. **List available backups**
   ```bash
   ls -la /backups/
   ```

2. **Restore database**
   ```bash
   dropdb -U postgres openpayment
   createdb -U postgres openpayment
   pg_restore -U postgres -d openpayment /backups/openpayment_$(date +%Y%m%d).dump
   ```

3. **Re-run migrations**
   ```bash
   # This is automatic — restart the API
   docker-compose restart api
   ```

## Investigate Slow Queries

1. **Enable query logging**
   ```sql
   SET log_min_duration_statement = 500;  -- log queries > 500ms
   ```

2. **Find slow queries**
   ```sql
   SELECT query, calls, total_time / calls AS avg_time_ms
   FROM pg_stat_statements
   ORDER BY avg_time_ms DESC
   LIMIT 10;
   ```

3. **Check for missing indexes**
   ```sql
   SELECT schemaname, tablename, indexname, idx_scan
   FROM pg_stat_user_indexes
   ORDER BY idx_scan ASC;
   ```

4. **Analyze query plans**
   ```sql
   EXPLAIN ANALYZE SELECT * FROM payment_intents WHERE ...
   ```

5. **Common fixes**
   - Add missing indexes
   - VACUUM ANALYZE
   - Increase `work_mem` for sort-heavy queries
   - Add connection pooling (PgBouncer)
