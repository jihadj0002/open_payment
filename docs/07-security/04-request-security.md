# Request Security

## HMAC Request Signing

### Signature Generation Algorithm
```
signature = HEX(HMAC-SHA256(
    key    = secret_api_key,
    data   = HTTP_METHOD + "\n"
           + URI_PATH + "\n"
           + QUERY_STRING + "\n"
           + CONTENT_TYPE + "\n"
           + TIMESTAMP + "\n"
           + NONCE + "\n"
           + REQUEST_BODY
))
```

### Implementation (Go)
```go
func GenerateSignature(secret string, method, path, query, contentType string,
    timestamp int64, nonce string, body []byte) string {

    mac := hmac.New(sha256.New, []byte(secret))

    mac.Write([]byte(method))
    mac.Write([]byte("\n"))
    mac.Write([]byte(path))
    mac.Write([]byte("\n"))
    mac.Write([]byte(query))
    mac.Write([]byte("\n"))
    mac.Write([]byte(contentType))
    mac.Write([]byte("\n"))
    mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
    mac.Write([]byte("\n"))
    mac.Write([]byte(nonce))
    mac.Write([]byte("\n"))
    mac.Write(body)

    return hex.EncodeToString(mac.Sum(nil))
}
```

### Verification (Server)
```go
func VerifySignature(r *http.Request, secret string) error {
    provided := r.Header.Get("X-Signature")
    timestamp := r.Header.Get("X-Timestamp")
    nonce := r.Header.Get("X-Nonce")

    if provided == "" || timestamp == "" || nonce == "" {
        return ErrMissingSignature
    }

    // Check timestamp freshness (max 5 minutes)
    ts, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil || time.Since(time.UnixMilli(ts)) > 5*time.Minute {
        return ErrExpiredTimestamp
    }

    // Check nonce uniqueness (prevent replay)
    if redis.Exists("nonce:" + nonce) {
        return ErrReplayedRequest
    }
    redis.Set("nonce:" + nonce, true, 5*time.Minute)

    // Read body
    body, err := io.ReadAll(r.Body)
    r.Body = io.NopCloser(bytes.NewBuffer(body))

    // Generate expected signature
    expected := GenerateSignature(secret,
        r.Method, r.URL.Path, r.URL.RawQuery,
        r.Header.Get("Content-Type"), ts, nonce, body)

    // Constant-time comparison
    if !hmac.Equal([]byte(provided), []byte(expected)) {
        return ErrInvalidSignature
    }

    return nil
}
```

## Replay Prevention

### Nonce Uniqueness
- Each request must include a unique `X-Nonce` (UUID v4)
- Used nonces stored in Redis with TTL = 5 minutes + timestamp drift margin
- Prevents same signed request from being replayed

### Timestamp Validation
- Server checks `X-Timestamp` is within ±5 minutes of server time
- Prevents captured requests from being replayed much later

## CSRF Protection
- **State-changing requests:** Require `X-CSRF-Token` header (not cookie)
- **SameSite:** Cookies set to `SameSite=Lax` for dashboard
- **Origin check:** Server validates `Origin` header matches allowed origins

## CORS Configuration
```go
allowedOrigins = []string{
    "https://dashboard.openpayment.com",
    "https://sandbox.openpayment.com",
}

func CORSConfig() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type", "Idempotency-Key",
                            "X-Signature", "X-Timestamp", "X-Nonce"},
        ExposeHeaders:    []string{"X-Request-Id", "X-RateLimit-Remaining"},
        AllowCredentials: false,  // API keys, not cookies
        MaxAge:           3600,
    })
}
```

## Input Validation

### Server-Side Validation
```go
type CreatePaymentRequest struct {
    Amount      int64  `json:"amount" binding:"required,min=1,max=999999999"`
    Currency    string `json:"currency" binding:"required,len=3,oneof=BDT USD"`
    PaymentMethod string `json:"payment_method" binding:"required,oneof=card wallet bank_transfer"`

    // Card validation
    CardNumber  string `json:"card_number" binding:"omitempty,creditcard"`
    ExpMonth    int    `json:"exp_month" binding:"omitempty,min=1,max=12"`
    ExpYear     int    `json:"exp_year" binding:"omitempty,min=2026,max=2040"`
    CVV         string `json:"cvv" binding:"omitempty,len=3|len=4"`
}
```

### Validation Rules
| Field | Rules |
|-------|-------|
| Amount | Positive integer, min 1, max per merchant config |
| Currency | ISO 4217, supported by platform |
| Card number | Luhn algorithm, BIN prefix check |
| Email | RFC 5322, max 255 chars |
| Phone | E.164 format, valid prefix |
| URL | HTTPS only, valid format, max 500 chars |
| Metadata keys | Max 40 chars, alphanumeric + underscore |
| Metadata values | Max 500 chars |
| Idempotency key | UUID v4 format |

## SQL Injection Prevention
- **ALWAYS use parameterized queries** (never string interpolation)
- **ORM:** Use `sqlx` or `pgx` named parameters
- **Dynamic filtering:** Build WHERE clauses programmatically with parameterized values

```go
// SAFE: parameterized
db.Query("SELECT * FROM payments WHERE merchant_id = $1 AND status = $2", merchantID, status)

// DANGEROUS: string interpolation
db.Query(fmt.Sprintf("SELECT * FROM payments WHERE merchant_id = '%s'", userInput))
```

## Additional Protections
- **Rate limiting:** Per API key, per endpoint, per IP
- **Request size limit:** 1MB max body size
- **Header size limit:** 8KB max header size
- **Timeout:** 30s for all API requests
- **IP allowlisting:** Optional per API key (merchant-configured)
