# First PR to Production Walkthrough

## Scenario: Adding Partial Refund Support

This walkthrough demonstrates the complete development lifecycle — from a feature request to production deployment.

### Step 1: Ticket Creation
```
JIRA: OP-456 — Add partial refund support
Type: Feature
Priority: High
Description:
  As a merchant, I want to issue partial refunds on captured payments
  so that I can handle partial order returns.

Acceptance Criteria:
  1. POST /v1/refunds accepts partial amount (less than payment total)
  2. Multiple partial refunds allowed until captured amount exhausted
  3. Refund total cannot exceed captured amount
  4. Webhook event sent for each partial refund
  5. Ledger entries created for partial refund
  6. Merchant balance updated correctly
```

### Step 2: Design (Tech Lead Review)
```
Design Doc:
  - Modify POST /v1/refunds to accept optional `amount` field
  - When `amount` omitted, full refund (backward compatible)
  - Validate: amount > 0, amount <= remaining_capturable
  - Track `amount_refunded` on payment_intent
  - One refund record per partial refund
  - Ledger entries: reverse pro-rata portion of fees and taxes
  - Webhook: refund.completed with partial amount

Reviewed by: Tech Lead ✅
Date: 2026-01-10
```

### Step 3: Branch Creation
```bash
git checkout develop
git pull origin develop
git checkout -b feature/OP-456-partial-refund
```

### Step 4: Backend Implementation

#### 4a. Update Models
```go
// internal/service/payment/models.go

// Add to PaymentIntent:
type PaymentIntent struct {
    // ... existing fields
    AmountRefunded int64 `json:"amount_refunded" db:"amount_refunded"`
}
```

#### 4b. Database Migration
```sql
-- internal/database/migrations/000005_add_amount_refunded.up.sql
ALTER TABLE payment_intents ADD COLUMN amount_refunded BIGINT NOT NULL DEFAULT 0;

-- Add check constraint
ALTER TABLE payment_intents ADD CONSTRAINT check_refund_not_exceed_captured
    CHECK (amount_refunded + COALESCE((
        SELECT COALESCE(SUM(amount), 0) FROM refunds WHERE payment_id = id
    ), 0) <= amount_received);
```

#### 4c. Update Service Logic
```go
// internal/service/payment/service.go

func (s *PaymentService) CreateRefund(ctx context.Context, req CreateRefundRequest) (*Refund, error) {
    payment, err := s.repo.FindPaymentByID(ctx, req.PaymentID)
    if err != nil {
        return nil, ErrPaymentNotFound
    }

    // Validate payment is captured
    if payment.Status != PaymentStatusCaptured && payment.Status != PaymentStatusSucceeded {
        return nil, ErrPaymentNotRefundable
    }

    // Determine refund amount
    refundAmount := req.Amount
    if refundAmount == 0 {
        refundAmount = payment.AmountReceived - payment.AmountRefunded
    }

    // Validate amount
    if refundAmount <= 0 {
        return nil, ErrInvalidRefundAmount
    }
    remaining := payment.AmountReceived - payment.AmountRefunded
    if refundAmount > remaining {
        return nil, fmt.Errorf("%w: remaining %d, requested %d", ErrRefundAmountExceeded, remaining, refundAmount)
    }

    // Process refund with processor
    processorRefund, err := s.processor.Refund(ctx, payment, refundAmount)
    if err != nil {
        return nil, fmt.Errorf("processor refund failed: %w", err)
    }

    // Create refund record
    refund := &Refund{
        ID:        uuid.New(),
        PaymentID: payment.ID,
        Amount:    refundAmount,
        Status:    RefundStatusSucceeded,
        Reason:    req.Reason,
    }
    if err := s.repo.CreateRefund(ctx, refund); err != nil {
        return nil, err
    }

    // Update payment amount_refunded
    if err := s.repo.IncrementRefunded(ctx, payment.ID, refundAmount); err != nil {
        return nil, err
    }

    // Create ledger entries (pro-rata fee reversal)
    if err := s.ledger.CreateRefundEntries(ctx, payment, refundAmount); err != nil {
        return nil, err
    }

    // Determine new payment status
    newRefunded := payment.AmountRefunded + refundAmount
    newStatus := PaymentStatusPartiallyRefunded
    if newRefunded >= payment.AmountReceived {
        newStatus = PaymentStatusRefunded
    }
    s.repo.UpdatePaymentStatus(ctx, payment.ID, newStatus)

    // Produce event
    s.producer.Produce(ctx, "payment.refunded", refund)

    return refund, nil
}
```

#### 4d. Update API Handler
```go
// internal/service/payment/handler.go

type CreateRefundRequest struct {
    PaymentID string `json:"payment_id" binding:"required,uuid"`
    Amount    int64  `json:"amount" binding:"omitempty,min=1"`
    Reason    string `json:"reason" binding:"omitempty,oneof=customer_request duplicate fraudulent other"`
}
```

#### 4e. Write Unit Tests
```go
// internal/service/payment/service_test.go

func TestCreateRefund_Partial_Success(t *testing.T) {
    // Setup
    payment := &PaymentIntent{
        ID:             uuid.New(),
        Amount:         1000,
        AmountReceived: 1000,
        AmountRefunded: 0,
        Status:         PaymentStatusSucceeded,
    }

    mockRepo := new(MockPaymentRepository)
    mockRepo.On("FindPaymentByID", mock.Anything, payment.ID).Return(payment, nil)
    mockRepo.On("CreateRefund", mock.Anything, mock.Anything).Return(nil)
    mockRepo.On("IncrementRefunded", mock.Anything, payment.ID, int64(300)).Return(nil)
    mockRepo.On("UpdatePaymentStatus", mock.Anything, payment.ID, PaymentStatusPartiallyRefunded).Return(nil)

    mockProcessor := new(MockProcessorClient)
    mockProcessor.On("Refund", mock.Anything, payment, int64(300)).Return(&ProcessorRefund{Status: "succeeded"}, nil)

    mockLedger := new(MockLedgerService)
    mockLedger.On("CreateRefundEntries", mock.Anything, payment, int64(300)).Return(nil)

    // Execute
    svc := NewPaymentService(mockRepo, mockProcessor, nil, mockLedger, nil)
    refund, err := svc.CreateRefund(context.Background(), CreateRefundRequest{
        PaymentID: payment.ID.String(),
        Amount:    300,
        Reason:    "customer_request",
    })

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, int64(300), refund.Amount)
    assert.Equal(t, RefundStatusSucceeded, refund.Status)
    assert.Equal(t, PaymentStatusPartiallyRefunded, payment.Status)
    mockRepo.AssertExpectations(t)
}

func TestCreateRefund_ExceedsRemaining_ReturnsError(t *testing.T) {
    payment := &PaymentIntent{
        ID:             uuid.New(),
        Amount:         1000,
        AmountReceived: 1000,
        AmountRefunded: 800,
        Status:         PaymentStatusPartiallyRefunded,
    }

    mockRepo := new(MockPaymentRepository)
    mockRepo.On("FindPaymentByID", mock.Anything, payment.ID).Return(payment, nil)

    svc := NewPaymentService(mockRepo, nil, nil, nil, nil)
    _, err := svc.CreateRefund(context.Background(), CreateRefundRequest{
        PaymentID: payment.ID.String(),
        Amount:    300,
    })

    assert.ErrorIs(t, err, ErrRefundAmountExceeded)
}
```

### Step 5: Run Tests
```bash
go test ./internal/service/payment/... -v -run TestCreateRefund -race
=== RUN   TestCreateRefund_Partial_Success
=== RUN   TestCreateRefund_ExceedsRemaining_ReturnsError
--- PASS: TestCreateRefund_Partial_Success (0.01s)
--- PASS: TestCreateRefund_ExceedsRemaining_ReturnsError (0.01s)
PASS
ok  	github.com/openpayment/gateway/internal/service/payment	2.345s
```

### Step 6: Commit
```bash
git add internal/service/payment/
git add internal/database/migrations/000005_add_amount_refunded.up.sql
git add internal/database/migrations/000005_add_amount_refunded.down.sql
git commit -m "feat(payment): add partial refund support

- Add amount field to POST /v1/refunds (optional, defaults to full)
- Add amount_refunded column to payment_intents
- Implement pro-rata fee reversal in ledger
- Track payment status: partially_refunded vs refunded
- Unit tests for partial refund scenarios

Closes OP-456"
```

### Step 7: Push and Create PR
```bash
git push origin feature/OP-456-partial-refund
gh pr create \
  --title "feat(payment): add partial refund support" \
  --body "Implements OP-456 — Partial Refund Support

Changes:
- Add optional amount parameter to refund endpoint
- Track amount_refunded on payment_intents
- Add partially_refunded payment status
- Ledger pro-rata fee reversal
- Unit tests covering all scenarios" \
  --reviewer tech-lead \
  --base develop
```

### Step 8: CI Pipeline Runs
```
PR #42 opened → CI triggers:
  [1] Lint: ✅ passed (30s)
  [2] Unit Tests: ✅ passed (45 tests, 2.3s)
  [3] Integration Tests: ✅ passed (12 tests, 18s)
  [4] Build: ✅ passed (Docker image built)
  [5] Security Scan: ✅ passed (0 critical, 0 high)
```

### Step 9: Code Review
```
Reviewer: Tech Lead
Comments:
  - [minor] Consider adding validation for refund amount
    minimum (1 BDT) ✓ (done)
  - [question] Should we allow refund of 0 amount?
    → No, removed that possibility
  - LGTM ✅

Approved by: Tech Lead
Date: 2026-01-11 14:30
```

### Step 10: Merge to Develop
```bash
gh pr merge --squash
git checkout develop
git pull origin develop
```

### Step 11: Deploy to Staging
```
CI detects merge to develop:
  - Builds image
  - Runs full test suite
  - Deploys to staging via Helm upgrade
  - Runs smoke tests against staging
  All passed ✅
  
Staging URL: https://staging.api.openpayment.com/v1/refunds
```

### Step 12: QA Verification
```
QA Engineer tests on staging:
  ✅ Full refund (no amount) — works as before
  ✅ Partial refund (amount=300) — refund processed
  ✅ Multiple partial refunds — correct remaining calculation
  ✅ Exceed remaining — returns 400 error
  ✅ Refund on uncaptured payment — returns 400 error
  ✅ Webhook delivered with correct amount
  ✅ Ledger shows correct entries
  ✅ Balance updated correctly

QA Sign-off: ✅
```

### Step 13: Merge to Main
```bash
git checkout main
git pull origin main
git merge develop
git push origin main
```

### Step 14: Production Canary Deploy
```
CI detects merge to main:
  - Builds production image
  - Deploys canary (5% traffic → 2 min)
    → Health check: ✅, Error rate: ✅, Latency: ✅
  - Increases to 25% traffic → 2 min
    → Health check: ✅, Error rate: ✅, Latency: ✅
  - Full rollout (100% traffic)
    → All metrics nominal
  - Slack notification: "Deployed: Payment Partial Refund"
```

### Step 15: Post-Deploy Monitoring
```
Monitoring watch (2 hours):
  - Error rate: 0.02% (normal baseline)
  - Refund API p95: 120ms (within target)
  - Payment success rate: 98.2% (unchanged)
  - No PagerDuty alerts
  - Grafana dashboard: refund count increasing as expected
```

### Step 16: Close Ticket
```
JIRA: OP-456
Status: Done ✅
Resolution: Deployed to production 2026-01-12 10:30 UTC
Release notes: "Partial refunds now supported — use amount parameter"
```

## Total Time: 2 days (design 0.5d, implementation 1d, testing 0.5d)
## Lines Changed: +245 / -12 across 8 files
## Releases: Staging → 1 hour, Production → 1 day after QA
