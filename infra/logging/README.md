# Logging Stack

The Open Payment Gateway uses **Loki + Promtail + Grafana** for log aggregation.

## Components

- **Loki** — horizontally-scalable, highly-available log aggregation system (port 3100).
- **Promtail** — log collector that scrapes container logs and pushes them to Loki (port 9080).
- **Grafana** — visualization frontend, pre-configured with a Loki datasource.

## Access

Port-forward Loki to query it directly:
```bash
kubectl port-forward svc/loki 3100:3100
```

Then add a Loki datasource in Grafana at `http://loki:3100` under **Configuration > Data Sources**.

## Querying

Use LogQL in Grafana Explore. Example queries:

- All logs from the API app:
  ```
  {app="openpayment-api"}
  ```

- Filter by namespace and pod:
  ```
  {namespace="production", pod="openpayment-api-*"}
  ```

- Count error logs per 5m:
  ```
  sum(rate({app="openpayment-api"} |= "error" [5m])) by (level)
  ```

## Retention

- Samples older than 7 days (168h) are rejected.
- Logs are retained for 14 days (336h).
