# Kubernetes Setup

## Cluster Configuration

### EKS Cluster
```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: openpayment-prod
  region: ap-south-1
  version: "1.29"

vpc:
  cidr: "10.0.0.0/16"
  subnets:
    private:
      ap-south-1a: { cidr: "10.0.2.0/24" }
      ap-south-1b: { cidr: "10.0.3.0/24" }
      ap-south-1c: { cidr: "10.0.4.0/24" }
    public:
      ap-south-1a: { cidr: "10.0.1.0/24" }
  nat:
    gateway: Single

nodeGroups:
  - name: system
    instanceType: t3.medium
    desiredCapacity: 2
    labels:
      role: system
    taints:
      - key: "CriticalAddonsOnly"
        value: "true"
        effect: "NoSchedule"

  - name: services
    instanceType: m6i.large
    desiredCapacity: 3
    minSize: 3
    maxSize: 20
    labels:
      role: services
    tags:
      k8s.io/cluster-autoscaler/enabled: "true"

  - name: data
    instanceType: r6i.large
    desiredCapacity: 3
    labels:
      role: data

  - name: batch
    instanceType: m6i.large
    desiredCapacity: 1
    minSize: 1
    maxSize: 5
    labels:
      role: batch
    taints:
      - key: "batch"
        value: "true"
        effect: "NoSchedule"

iam:
  withOIDC: true
  serviceAccounts:
    - metadata:
        name: payment-service
      attachPolicy:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Action:
              - "kms:Decrypt"
              - "kms:GenerateDataKey"
            Resource: "arn:aws:kms:ap-south-1:xxx:key/xxx"
```

## Core Add-ons

| Add-on | Purpose | Installation |
|--------|---------|--------------|
| ingress-nginx | Ingress controller | Helm chart |
| cert-manager | TLS certificate management | Helm chart |
| external-dns | DNS record management | Helm chart |
| cluster-autoscaler | Node auto-scaling | Helm chart |
| metrics-server | Resource metrics | Helm chart |
| Prometheus Stack | Monitoring (kube-prometheus) | Helm chart |
| Loki | Log aggregation | Helm chart |
| Tempo | Distributed tracing | Helm chart |
| OpenTelemetry Collector | Trace/metric collection | Helm chart |
| Istio | Service mesh (future) | Helm chart |
| Falco | Runtime security | Helm chart |
| Goldilocks | Resource recommendations | Helm chart |

## Resource Quotas

### Per Namespace
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: payment-ns-quota
  namespace: payment
spec:
  hard:
    requests.cpu: "20"
    requests.memory: "40Gi"
    limits.cpu: "40"
    limits.memory: "80Gi"
    persistentvolumeclaims: "10"
    pods: "50"
```

### Horizontal Pod Autoscaler (Payment Service)
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: payment-service-hpa
  namespace: payment
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: payment-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 75
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
    scaleUp:
      stabilizationWindowSeconds: 60
```

## Pod Disruption Budgets
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: payment-service-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: payment-service
```

## Helm Chart Structure

```
helm/
├── charts/                    # Local dependencies
├── payment-service/
│   ├── Chart.yaml
│   ├── values.yaml
│   ├── values-prod.yaml
│   ├── values-staging.yaml
│   └── templates/
│       ├── deployment.yaml
│       ├── service.yaml
│       ├── hpa.yaml
│       ├── pdb.yaml
│       ├── serviceaccount.yaml
│       ├── configmap.yaml
│       ├── secret.yaml        # (references external secrets)
│       ├── networkpolicy.yaml
│       └── servicemonitor.yaml
├── merchant-service/
│   └── ...
└── global/
    └── values.yaml            # Shared values across all services
```
