# Go Modular Monolith Structure

## Project Structure

```
openpayment/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point, dependency injection
│
├── internal/
│   ├── api/
│   │   ├── router.go                # HTTP router setup
│   │   ├── middleware/
│   │   │   ├── auth.go              # JWT/API key validation
│   │   │   ├── rate_limiter.go      # Rate limiting
│   │   │   ├── logging.go           # Request logging
│   │   │   ├── recovery.go          # Panic recovery
│   │   │   ├── cors.go              # CORS
│   │   │   ├── idempotency.go       # Idempotency check
│   │   │   └── request_id.go        # X-Request-Id injection
│   │   └── response.go              # Standard response helpers
│   │
│   ├── service/                     # Business logic layers
│   │   ├── auth/
│   │   │   ├── service.go           # AuthService interface + impl
│   │   │   ├── service_test.go
│   │   │   ├── handler.go           # HTTP handlers
│   │   │   ├── handler_test.go
│   │   │   ├── repository.go        # DB operations
│   │   │   ├── repository_test.go
│   │   │   └── models.go            # Domain types
│   │   ├── merchant/
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   ├── repository.go
│   │   │   └── models.go
│   │   ├── payment/
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   ├── repository.go
│   │   │   ├── models.go
│   │   │   └── state_machine.go     # Payment state transitions
│   │   ├── customer/
│   │   ├── ledger/
│   │   ├── settlement/
│   │   ├── fraud/
│   │   ├── webhook/
│   │   └── notification/
│   │
│   ├── processor/                   # External payment processor adapters
│   │   ├── client.go                # Interface
│   │   ├── card/
│   │   │   └── processor.go         # Card network integration
│   │   ├── wallet/
│   │   │   └── bkash.go             # bKash integration
│   │   └── bank/
│   │       └── transfer.go          # Bank transfer integration
│   │
│   ├── messaging/
│   │   ├── producer.go              # Kafka producer
│   │   └── consumer.go              # Kafka consumer setup
│   │
│   ├── config/
│   │   ├── config.go                # Configuration loading
│   │   └── config_test.go
│   │
│   ├── database/
│   │   ├── postgres.go              # PostgreSQL connection
│   │   ├── migrations/              # SQL migration files
│   │   │   ├── 000001_create_merchants.up.sql
│   │   │   ├── 000001_create_merchants.down.sql
│   │   │   ├── 000002_create_payments.up.sql
│   │   │   └── ...
│   │   └── repository.go            # Base repository
│   │
│   ├── cache/
│   │   ├── redis.go                 # Redis client
│   │   └── rate_limiter.go          # Rate limit implementation
│   │
│   └── pkg/
│       ├── logger/
│       │   └── logger.go            # Structured logging
│       ├── validator/
│       │   └── validator.go         # Input validation
│       ├── encrypt/
│       │   └── encrypt.go           # AES encryption
│       └── errors/
│           └── errors.go            # Domain error types
│
├── pkg/                             # Public shared packages
│   ├── money/
│   │   └── money.go                 # Money type (amount, currency)
│   ├── pagination/
│   │   └── pagination.go
│   └── testutil/
│       └── fixtures.go              # Test helpers
│
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
├── .golangci.yml
└── README.md
```

## Dependency Injection Wire

```go
// cmd/server/main.go
func main() {
    cfg := config.Load()
    logger := logger.New(cfg.LogLevel)

    db := database.NewPostgres(cfg.Database)
    cache := cache.NewRedis(cfg.Redis)
    kafkaProducer := messaging.NewProducer(cfg.Kafka)

    // Repositories
    merchantRepo := merchant.NewRepository(db)
    paymentRepo := payment.NewRepository(db)
    customerRepo := customer.NewRepository(db)
    ledgerRepo := ledger.NewRepository(db)

    // Services
    authSvc := auth.NewService(merchantRepo, cache, cfg.JWT)
    paymentSvc := payment.NewService(paymentRepo, ledgerRepo, kafkaProducer)
    merchantSvc := merchant.NewService(merchantRepo)
    ledgerSvc := ledger.NewService(ledgerRepo)

    // Router
    router := api.NewRouter(logger)
    router.Use(middleware.RequestID())
    router.Use(middleware.Logging(logger))
    router.Use(middleware.Recovery(logger))
    router.Use(middleware.CORS(cfg.CORS))

    // Auth routes
    authHandler := auth.NewHandler(authSvc)
    authHandler.Register(router)

    // Merchant routes (authenticated)
    merchantHandler := merchant.NewHandler(merchantSvc)
    merchantHandler.Register(router.Group("/v1", middleware.Auth(authSvc)))

    // Payment routes
    paymentHandler := payment.NewHandler(paymentSvc)
    paymentHandler.Register(router.Group("/v1", middleware.Auth(authSvc)))

    // Start server
    srv := &http.Server{Addr: ":" + cfg.Port, Handler: router}
    srv.ListenAndServe()
}
```

## Coding Standards

### Naming
| Convention | Example |
|------------|---------|
| Interface names | `PaymentService`, `MerchantRepository` |
| Interface suffix: `-er` for single method | `Logger`, `Validator` |
| Package names lowercase | `payment`, `merchant`, `ledger` |
| File names snake_case | `state_machine.go`, `rate_limiter.go` |
| Constants ALL_CAPS | `StatusCreated`, `MaxRetries` |
| Error vars `Err` prefix | `ErrPaymentNotFound`, `ErrInvalidState` |

### Error Handling
```go
// Use domain-specific error types
var (
    ErrPaymentNotFound    = errors.New("payment not found")
    ErrInvalidState      = errors.New("payment is not in a valid state for this operation")
    ErrInsufficientFunds = errors.New("insufficient funds")
)

// Wrap errors with context
if err := repo.Save(ctx, payment); err != nil {
    return fmt.Errorf("saving payment %s: %w", payment.ID, err)
}

// Check error types
if errors.Is(err, ErrPaymentNotFound) {
    return api.NotFound("payment_not_found", "Payment not found")
}
```

### Testing
```go
// Table-driven tests
func TestCreatePayment(t *testing.T) {
    tests := []struct {
        name    string
        req     CreatePaymentRequest
        want    PaymentIntent
        wantErr error
    }{
        {
            name: "valid card payment",
            req: CreatePaymentRequest{
                Amount:   1000,
                Currency: "BDT",
                PaymentMethod: "card",
            },
            want: PaymentIntent{
                Status: "created",
                Amount: 1000,
            },
            wantErr: nil,
        },
        {
            name: "invalid amount",
            req:  CreatePaymentRequest{Amount: -100},
            wantErr: ErrInvalidAmount,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := NewPaymentService(mockRepo, mockLedger, mockProducer)
            got, err := svc.CreatePayment(context.Background(), tt.req)

            if !errors.Is(err, tt.wantErr) {
                t.Errorf("got err %v, want %v", err, tt.wantErr)
            }
            if err == nil && !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %+v, want %+v", got, tt.want)
            }
        })
    }
}
```

## Makefile
```makefile
.PHONY: build test lint run migrate

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

test-verbose:
	go test ./... -v -race

lint:
	golangci-lint run ./...

run:
	go run ./cmd/server

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

docker-build:
	docker build -t openpayment:latest .

coverage-html:
	go tool cover -html=coverage.out

# CI-friendly
ci: lint test build
```
