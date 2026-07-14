# Internal gRPC API Specification

## Overview
Internal services communicate via gRPC over TLS within the Kubernetes cluster. Protobuf v3 is used for schema definition, with backward-compatible evolution.

## Proto Package Structure
```
proto/
├── auth/
│   └── auth.proto
├── merchant/
│   └── merchant.proto
├── payment/
│   └── payment.proto
├── ledger/
│   └── ledger.proto
├── fraud/
│   └── fraud.proto
├── settlement/
│   └── settlement.proto
├── webhook/
│   └── webhook.proto
└── common/
    ├── money.proto
    └── timestamp.proto
```

## Auth Service

### ValidateToken
```protobuf
service AuthService {
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
    rpc ExchangeAPIKey(ExchangeAPIKeyRequest) returns (ExchangeAPIKeyResponse);
    rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
}

message ValidateTokenRequest {
    string token = 1;
}

message ValidateTokenResponse {
    bool valid = 1;
    string merchant_id = 2;
    string user_id = 3;
    string role = 4;
    repeated string permissions = 5;
    int64 expires_at = 6;
}
```

## Merchant Service

### GetMerchantConfig
```protobuf
service MerchantService {
    rpc GetMerchant(GetMerchantRequest) returns (Merchant);
    rpc GetMerchantConfig(GetMerchantConfigRequest) returns (MerchantConfig);
    rpc ValidateAPIKey(ValidateAPIKeyRequest) returns (ValidateAPIKeyResponse);
}

message MerchantConfig {
    string merchant_id = 1;
    int64 fee_rate_bps = 2;       // fee rate in basis points (250 = 2.5%)
    int64 fixed_fee = 3;          // fixed fee in smallest currency unit
    repeated string supported_currencies = 4;
    int64 max_transaction_amount = 5;
    int64 daily_volume_limit = 6;
    bool auto_capture = 7;
    string webhook_secret = 8;
}
```

## Payment Service

### PaymentProcessing
```protobuf
service PaymentService {
    rpc CreatePayment(CreatePaymentRequest) returns (PaymentIntent);
    rpc GetPayment(GetPaymentRequest) returns (PaymentIntent);
    rpc AuthorizePayment(AuthorizePaymentRequest) returns (AuthorizePaymentResponse);
    rpc CapturePayment(CapturePaymentRequest) returns (CapturePaymentResponse);
    rpc RefundPayment(RefundPaymentRequest) returns (RefundResponse);
}

message AuthorizePaymentRequest {
    string payment_id = 1;
    string merchant_id = 2;
    int64 amount = 3;
    string currency = 4;
    PaymentMethodDetails payment_method = 5;
    string idempotency_key = 6;
}

message AuthorizePaymentResponse {
    bool success = 1;
    string processor_transaction_id = 2;
    string authorization_code = 3;
    string processor_response_code = 4;
    string error_message = 5;
}
```

## Ledger Service

### Double-Entry Accounting
```protobuf
service LedgerService {
    rpc CreateLedgerEntry(CreateLedgerEntryRequest) returns (LedgerEntry);
    rpc GetBalance(GetBalanceRequest) returns (Balance);
    rpc GetAccountStatement(GetAccountStatementRequest) returns (AccountStatement);
    rpc Reconcile(ReconcileRequest) returns (ReconcileResponse);
}

message CreateLedgerEntryRequest {
    string transaction_id = 1;
    string merchant_id = 2;
    repeated LedgerLine lines = 3;
    string description = 4;
    int64 timestamp = 5;
    string idempotency_key = 6;
}

message LedgerLine {
    string account_id = 1;
    string account_type = 2;   // "merchant_receivable", "gateway_fee", "tax", "settlement"
    int64 amount = 3;
    string direction = 4;      // "debit" or "credit"
    string currency = 5;
}

message Balance {
    string merchant_id = 1;
    int64 available_amount = 2;
    int64 pending_amount = 3;
    int64 reserve_amount = 4;
    string currency = 5;
}
```

## Fraud Service

### Risk Assessment
```protobuf
service FraudService {
    rpc AssessRisk(AssessRiskRequest) returns (RiskAssessment);
    rpc ReportFraudEvent(ReportFraudEventRequest) returns (ReportFraudEventResponse);
}

message AssessRiskRequest {
    string payment_id = 1;
    string merchant_id = 2;
    int64 amount = 3;
    string currency = 4;
    string customer_ip = 5;
    string customer_email = 6;
    string card_bin = 7;
    string device_fingerprint = 8;
    string customer_phone = 9;
}

message RiskAssessment {
    double score = 1;           // 0.0 (safe) to 1.0 (fraudulent)
    string action = 2;          // "allow", "review", "block"
    repeated string flags = 3;  // ["high_velocity", "vpn_detected", "bin_country_mismatch"]
    string assessment_id = 4;
}
```

## Settlement Service

### Settlement Processing
```protobuf
service SettlementService {
    rpc ProcessSettlement(ProcessSettlementRequest) returns (ProcessSettlementResponse);
    rpc GetSettlementStatus(GetSettlementStatusRequest) returns (SettlementStatus);
}

message ProcessSettlementRequest {
    string merchant_id = 1;
    string date = 2;      // "2026-01-15"
}

message ProcessSettlementResponse {
    string settlement_id = 1;
    int64 total_amount = 2;
    int64 fee_amount = 3;
    int64 net_amount = 4;
    int64 transaction_count = 5;
    string status = 6;    // "processing", "completed", "failed"
}
```

## gRPC Error Handling
Standard gRPC error codes mapped to domain errors:

| gRPC Code | Domain Error | Description |
|-----------|-------------|-------------|
| `INVALID_ARGUMENT` | Validation error | Missing or invalid field |
| `NOT_FOUND` | Resource not found | Payment, merchant, etc. |
| `ALREADY_EXISTS` | Duplicate | Idempotency key already used |
| `FAILED_PRECONDITION` | Invalid state | Payment not capturable |
| `PERMISSION_DENIED` | Authorization | API key invalid, insufficient permissions |
| `UNAVAILABLE` | Service down | Circuit breaker open, dependency down |
| `DEADLINE_EXCEEDED` | Timeout | Upstream call took too long |

## Service Mesh
- All gRPC traffic goes through Istio sidecar proxies
- mTLS enforced between all services
- Retry policies: 2 retries with 200ms backoff
- Circuit breaker: trip after 5 consecutive failures, half-open after 30s
- Timeout: 5s for most calls, 30s for settlement processing
