# CI/CD Pipeline

## GitHub Actions Workflows

### Main Pipeline (Build + Test + Deploy)

```yaml
name: Build and Deploy
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  AWS_REGION: ap-south-1
  ECR_REPOSITORY: openpayment
  K8S_NAMESPACE: payment

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  test:
    needs: [lint]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run unit tests
        run: go test ./... -race -coverprofile=coverage.out -covermode=atomic
      - name: Upload coverage
        uses: codecov/codecov-action@v4

  build-and-scan:
    needs: [test]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ env.AWS_REGION }}

      - name: Login to ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2

      - name: Build Docker image
        run: |
          docker build -t ${{ steps.login-ecr.outputs.registry }}/${{ env.ECR_REPOSITORY }}:${{ github.sha }} .

      - name: Scan for vulnerabilities (Trivy)
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: '${{ steps.login-ecr.outputs.registry }}/${{ env.ECR_REPOSITORY }}:${{ github.sha }}'
          format: 'sarif'
          output: 'trivy-results.sarif'
          severity: 'CRITICAL,HIGH'

      - name: Upload Trivy results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'

      - name: Push image to ECR
        run: |
          docker push ${{ steps.login-ecr.outputs.registry }}/${{ env.ECR_REPOSITORY }}:${{ github.sha }}
          docker tag ${{ steps.login-ecr.outputs.registry }}/${{ env.ECR_REPOSITORY }}:${{ github.sha }} \
            ${{ steps.login-ecr.outputs.registry }}/${{ env.ECR_REPOSITORY }}:latest
          docker push ${{ steps.login-ecr.outputs.registry }}/${{ env.ECR_REPOSITORY }}:latest

  deploy-staging:
    needs: [build-and-scan]
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: azure/setup-helm@v4

      - name: Deploy to staging
        run: |
          helm upgrade --install payment-service ./helm/payment-service \
            --namespace staging \
            --values ./helm/payment-service/values-staging.yaml \
            --set image.tag=${{ github.sha }} \
            --wait --timeout 5m

      - name: Run smoke tests
        run: |
          ./scripts/smoke-test.sh staging

  deploy-production:
    needs: [build-and-scan]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4
      - uses: azure/setup-helm@v4

      - name: Canary deploy (5% → 25% → 100%)
        run: |
          # Deploy canary
          helm upgrade --install payment-service ./helm/payment-service \
            --namespace production \
            --values ./helm/payment-service/values-prod.yaml \
            --set image.tag=${{ github.sha }} \
            --set canary.enabled=true \
            --set canary.weight=5 \
            --wait --timeout 5m

          # Monitor for 2 minutes
          sleep 120
          ./scripts/check-health.sh production

          # Increase to 25%
          helm upgrade --install payment-service ./helm/payment-service \
            --set canary.weight=25 \
            --reuse-values
          sleep 120
          ./scripts/check-health.sh production

          # Full rollout
          helm upgrade --install payment-service ./helm/payment-service \
            --set canary.enabled=false \
            --reuse-values
          sleep 120
          ./scripts/check-health.sh production

      - name: Slack notification
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: 'Payment service deployed to production (${{ github.sha }})'
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}

  security-scan-scheduled:
    if: github.event_name == 'schedule'
    runs-on: ubuntu-latest
    steps:
      - name: Weekly security scan
        run: |
          # Full dependency scan
          npm audit --audit-level=high
          go list -m all | nancy sleuth
          trivy fs --severity CRITICAL,HIGH .
```

## Pipeline Stages Summary

```
PR (develop/main)
  └── Lint + Test ──▶ Build Image ──▶ Security Scan ──▶ (results posted to PR)

Push to develop
  └── Lint + Test ──▶ Build + Scan ──▶ Deploy to Staging ──▶ Smoke Tests

Push to main
  └── Lint + Test ──▶ Build + Scan ──▶ Canary Deploy (5%→25%→100%)
       └── Health Check at each step ──▶ Slack Notification

Scheduled (weekly)
  └── Full Security Scan ──▶ Report to Security Channel
```

## Smoke Test Script
```bash
#!/bin/bash
# scripts/smoke-test.sh
ENV=$1

echo "Running smoke tests against ${ENV}..."

# Health check
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" https://api-${ENV}.openpayment.com/health)
if [ "$HEALTH" != "200" ]; then
    echo "Health check failed: $HEALTH"
    exit 1
fi

# Create payment
PAYMENT=$(curl -s -X POST https://api-${ENV}.openpayment.com/v1/payments \
    -H "Authorization: Bearer ${TEST_API_KEY}" \
    -H "Content-Type: application/json" \
    -d '{"amount": 100, "currency": "BDT", "payment_method": "card", ...}')
echo "Payment created: $PAYMENT"

# Verify payment was created
PAYMENT_ID=$(echo $PAYMENT | jq -r '.id')
STATUS=$(curl -s https://api-${ENV}.openpayment.com/v1/payments/$PAYMENT_ID \
    -H "Authorization: Bearer ${TEST_API_KEY}" | jq -r '.status')

if [ "$STATUS" != "succeeded" ]; then
    echo "Payment creation failed: $STATUS"
    exit 1
fi

echo "All smoke tests passed!"
```
