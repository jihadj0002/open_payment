# Chaos Engineering Experiments

These experiments use [Chaos Mesh](https://chaos-mesh.org/) to test the payment gateway's resilience.

## Prerequisites
- Chaos Mesh installed in the cluster: `helm install chaos-mesh chaos-mesh/chaos-mesh -n chaos-mesh`
- Namespace: `openpayment`

## Experiments

| Experiment | Type | Description |
|------------|------|-------------|
| pod-kill.yaml | PodChaos | Kill 1 random API pod, verify recovery |
| network-delay.yaml | NetworkChaos | Add 200ms latency to payment service |
| cpu-stress.yaml | StressChaos | Stress 1 pod to 80% CPU |

## Run an experiment
```bash
kubectl apply -f tests/chaos/pod-kill.yaml
# Observe recovery:
kubectl get pods -n openpayment -w
kubectl delete -f tests/chaos/pod-kill.yaml
```
