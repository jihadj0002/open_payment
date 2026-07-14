# Service Implementation Details

## Payment Service

### Key Types
```go
type PaymentIntent struct {
    ID              uuid.UUID       `json:"id"`
    MerchantID      uuid.UUID       `json:"merchant_id"`
    CustomerID      *uuid.UUID      `json:"customer_id,omitempty"`
    Amount          int64           `json:"amount"`
    AmountCapturable int64          `json:"amount_capturable"`
    AmountReceived  int64           `json:"amount_received"`
    Currency        string          `json:"currency"`
    Status          PaymentStatus   `json:"status"`
    CaptureMethod   CaptureMethod   `json:"capture_method"`
    PaymentMethod   string          `json:"payment_method,omitempty"`
    Description     string          `json:"description,omitempty"`
    Metadata        map[string]string `json:"metadata,omitempty"`
    CreatedAt       time.Time       `json:"created"`
    UpdatedAt       time.Time       `json:"updated"`
}

type PaymentStatus string
const (
    PaymentStatusCreated       PaymentStatus = "created"
    PaymentStatusPending       PaymentStatus = "pending"
    PaymentStatusProcessing    PaymentStatus = "processing"
    PaymentStatusAuthorized    PaymentStatus = "authorized"
    PaymentStatusCaptured      PaymentStatus = "captured"
    PaymentStatusSucceeded     PaymentStatus = "succeeded"
    PaymentStatusFailed        PaymentStatus = "failed"
    PaymentStatusCanceled      PaymentStatus = "canceled"
    PaymentStatusRefunded      PaymentStatus = "refunded"
)
```

### Core Logic Flow
```go
type PaymentService struct {
    repo        PaymentRepository
    processor   ProcessorClient
    fraud       FraudService
    ledger      LedgerService
    producer    EventProducer
    stateMachine *StateMachine
}

func (s *PaymentService) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentIntent, error) {
    payment := &PaymentIntent{
        ID:             uuid.New(),
        MerchantID:     ctx.Value(CtxMerchantID).(uuid.UUID),
        Amount:         req.Amount,
        AmountCapturable: req.Amount,
        Currency:       req.Currency,
        Status:         PaymentStatusCreated,
        CaptureMethod:  CaptureMethod(cfg.CaptureMethod),
        Description:    req.Description,
        CreatedAt:      time.Now(),
    }

    // Idempotency check
    if req.IdempotencyKey != "" {
        existing, err := s.repo.FindByIdempotencyKey(ctx, req.IdempotencyKey)
        if err == nil {
            return existing, nil
        }
        if err != ErrNotFound {
            return nil, err
        }
    }

    // Validate
    if err := validatePayment(payment); err != nil {
        return nil, err
    }

    // Save initial state
    if err := s.repo.Create(ctx, payment); err != nil {
        return nil, fmt.Errorf("creating payment: %w", err)
    }

    // Produce event
    s.producer.Produce(ctx, "payment.created", payment)

    if req.Confirm {
        return s.ConfirmPayment(ctx, payment.ID)
    }

    return payment, nil
}

func (s *PaymentService) ConfirmPayment(ctx context.Context, id uuid.UUID) (*PaymentIntent, error) {
    payment, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, ErrPaymentNotFound
    }

    if err := s.stateMachine.Transition(payment.Status, PaymentStatusProcessing); err != nil {
        return nil, err
    }

    // Fraud check
    risk, err := s.fraud.AssessRisk(ctx, payment)
    if err != nil {
        s.logger.Error("fraud check failed", "payment_id", id, "error", err)
    }
    if risk.Action == "block" {
        payment.Status = PaymentStatusFailed
        payment.ErrorMessage = "Transaction blocked by fraud detection"
        s.repo.Update(ctx, payment)
        s.producer.Produce(ctx, "payment.failed", payment)
        return payment, nil
    }

    // Authorize with processor
    auth, err := s.processor.Authorize(ctx, payment)
    if err != nil {
        payment.Status = PaymentStatusFailed
        payment.ErrorMessage = err.Error()
        s.repo.Update(ctx, payment)
        s.producer.Produce(ctx, "payment.failed", payment)
        return payment, err
    }

    payment.Status = PaymentStatusAuthorized
    payment.ProcessorTransactionID = auth.TransactionID
    s.repo.Update(ctx, payment)
    s.producer.Produce(ctx, "payment.authorized", payment)

    if payment.CaptureMethod == CaptureMethodAutomatic {
        return s.CapturePayment(ctx, id)
    }

    return payment, nil
}
```

## Ledger Service

### Double-Entry Implementation
```go
type LedgerService struct {
    repo LedgerRepository
}

func (s *LedgerService) CreateCaptureEntries(ctx context.Context, tx TransactionEntry) error {
    // Each capture creates balanced ledger entries
    total := tx.Amount
    fee := calculateFee(total, tx.FeeRate)
    tax := calculateTax(fee, tx.TaxRate)
    net := total - fee - tax

    entries := []LedgerEntry{
        {
            AccountID:  tx.MerchantReceivableAccountID,
            Direction:  "debit",
            Amount:     total,
            EntryType:  "payment_capture",
        },
        {
            AccountID:  tx.GatewayFeeAccountID,
            Direction:  "debit",
            Amount:     fee,
            EntryType:  "fee_deduction",
        },
        {
            AccountID:  tx.TaxPayableAccountID,
            Direction:  "credit",
            Amount:     tax,
            EntryType:  "tax_deduction",
        },
        {
            AccountID:  tx.ProcessorPayableAccountID,
            Direction:  "credit",
            Amount:     net,
            EntryType:  "payment_capture",
        },
    }

    // Verify balance
    var sum int64
    for _, e := range entries {
        if e.Direction == "debit" {
            sum += e.Amount
        } else {
            sum -= e.Amount
        }
    }
    if sum != 0 {
        return ErrUnbalancedEntry
    }

    return s.repo.BulkCreate(ctx, tx.ID, entries)
}
```

## State Machine

```go
type StateMachine struct {
    transitions map[PaymentStatus][]PaymentStatus
}

func NewStateMachine() *StateMachine {
    return &StateMachine{
        transitions: map[PaymentStatus][]PaymentStatus{
            PaymentStatusCreated:    {PaymentStatusPending, PaymentStatusFailed, PaymentStatusCanceled},
            PaymentStatusPending:    {PaymentStatusProcessing, PaymentStatusFailed, PaymentStatusCanceled},
            PaymentStatusProcessing: {PaymentStatusAuthorized, PaymentStatusFailed},
            PaymentStatusAuthorized: {PaymentStatusCaptured, PaymentStatusFailed, PaymentStatusCanceled},
            PaymentStatusCaptured:   {PaymentStatusSucceeded, PaymentStatusRefunded, PaymentStatusFailed},
            PaymentStatusSucceeded:  {PaymentStatusRefunded, PaymentStatusPartiallyRefunded},
        },
    }
}

func (sm *StateMachine) Transition(from, to PaymentStatus) error {
    allowed, ok := sm.transitions[from]
    if !ok {
        return fmt.Errorf("%w: %s → %s", ErrNoTransitionsDefined, from, to)
    }
    for _, s := range allowed {
        if s == to {
            return nil
        }
    }
    return fmt.Errorf("%w: %s → %s", ErrInvalidTransition, from, to)
}
```

## Test Approach

```go
// Unit tests: mock all dependencies
// Integration tests: use testcontainers for Postgres/Kafka/Redis
// E2E tests: deploy full stack, test via API

func TestPaymentFlow_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    ctx := context.Background()

    // Start test containers
    postgres, _ := testcontainers.StartPostgres(ctx)
    kafka, _ := testcontainers.StartKafka(ctx)
    redis, _ := testcontainers.StartRedis(ctx)

    defer postgres.Terminate(ctx)
    defer kafka.Terminate(ctx)
    defer redis.Terminate(ctx)

    // Run migrations
    db := database.NewPostgres(postgres.ConnectionString())
    database.RunMigrations(db)

    // Create service with real dependencies
    svc := payment.NewService(
        payment.NewRepository(db),
        processor.NewMockClient(),
        fraud.NewMockClient(),
        ledger.NewService(ledger.NewRepository(db)),
        messaging.NewProducer(kafka.Brokers()),
    )

    // Test complete flow
    payment, err := svc.CreatePayment(ctx, CreatePaymentRequest{
        Amount:   1000,
        Currency: "BDT",
        PaymentMethod: "card",
        Confirm:  true,
    })
    require.NoError(t, err)
    require.Equal(t, PaymentStatusAuthorized, payment.Status)
}
```
