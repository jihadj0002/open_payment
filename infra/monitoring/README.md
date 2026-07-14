# Monitoring Stack

This directory contains configuration for the Open Payment Gateway monitoring stack, consisting of **Prometheus**, **Grafana**, and **AlertManager**.

## Components

| Component      | Purpose                                    | Port  |
| -------------- | ------------------------------------------ | ----- |
| **Prometheus** | Metrics collection, storage, and alerting  | 9090  |
| **Grafana**    | Dashboards and visualization               | 3000  |
| **AlertManager** | Alert routing (Slack, PagerDuty)         | 9093  |

## Quick Start

1. Ensure Prometheus and Grafana are running (e.g. via Docker Compose or Kubernetes).
2. Prometheus scrapes the `openpayment-api` service at `/metrics`.
3. Grafana connects to the Prometheus datasource at `http://prometheus:9090`.
4. Import the dashboard from `grafana/dashboards/api-performance.json` (UID: `openpayment-api-performance`).

## Access

| Service      | URL                          |
| ------------ | ---------------------------- |
| Prometheus   | http://localhost:9090         |
| Grafana      | http://localhost:3000         |
| AlertManager | http://localhost:9093         |

Default Grafana credentials: `admin` / `admin` (change on first login).

## Dashboard

The **Open Payment Gateway — API Performance** dashboard includes:

- **Request Rate** — QPS over time
- **Latency (p50 / p95 / p99)** — percentile latency graph
- **Error Rate** — percentage of 5xx responses
- **Active Payments** — gauge showing in-flight payments
- **CPU / Memory** — per-pod resource usage
- **HTTP Status Breakdown** — stacked bar chart (2xx, 4xx, 5xx)

## Alerting Rules

Defined in `prometheus/alerts.yml`:

| Alert           | Condition                                     | Severity |
| --------------- | --------------------------------------------- | -------- |
| HighErrorRate   | > 5% of responses are 5xx for 2 minutes       | critical |
| HighLatency     | P95 latency > 0.5s for 2 minutes              | warning  |
| PodDown         | API pod is unresponsive for 1 minute           | critical |
| HighCPUUsage    | CPU usage > 80% for 5 minutes                 | warning  |

Alerts are routed via AlertManager to Slack (all alerts) and PagerDuty (critical only).

## k6 Integration

For load testing with k6, metrics can be pushed to a Prometheus remote-write endpoint:

```bash
k6 run --out output-prometheus-remote script.js
```

See `k6-metrics/k6-prometheus.yml` for configuration reference.

## Configuration

- **Prometheus**: `prometheus/prometheus.yml`
- **Alert Rules**: `prometheus/alerts.yml`
- **AlertManager**: `alertmanager/config.yml`
- **Grafana Datasource**: `grafana/datasources/prometheus.yml`
- **Grafana Dashboard**: `grafana/dashboards/api-performance.json`
