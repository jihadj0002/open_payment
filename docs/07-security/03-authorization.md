# Authorization

## RBAC Model

### Roles (Merchant Dashboard)

| Role | Description | Default Permissions |
|------|-------------|---------------------|
| `owner` | Full access to merchant account | All permissions |
| `admin` | Manage payments, refunds, settings | All except API key management, user management |
| `developer` | API keys, webhooks, logs | API key management, read-only payments |
| `analyst` | Reports, read-only data | Read-only access to payments, customers, reports |
| `member` | Basic view access | Read-only payments, customers |

### Roles (Admin Dashboard)

| Role | Description | Default Permissions |
|------|-------------|---------------------|
| `super_admin` | Full system access | All permissions |
| `operator` | Day-to-day operations | Merchant management, transactions, disputes |
| `compliance` | KYC/AML reviews | Merchant verification, audit logs |
| `support` | Merchant support | Read-only merchant info, support tickets |
| `readonly` | View-only access | Read-only everything |

### Permission Matrix

| Resource | Action | owner | admin | developer | analyst | member |
|----------|--------|-------|-------|-----------|---------|--------|
| payments | read | ✅ | ✅ | ✅ | ✅ | ✅ |
| payments | write | ✅ | ✅ | ❌ | ❌ | ❌ |
| refunds | read | ✅ | ✅ | ✅ | ✅ | ✅ |
| refunds | write | ✅ | ✅ | ❌ | ❌ | ❌ |
| customers | read | ✅ | ✅ | ✅ | ✅ | ✅ |
| customers | write | ✅ | ✅ | ❌ | ❌ | ❌ |
| api_keys | read | ✅ | ❌ | ✅ | ❌ | ❌ |
| api_keys | write | ✅ | ❌ | ✅ | ❌ | ❌ |
| webhooks | read | ✅ | ✅ | ✅ | ❌ | ❌ |
| webhooks | write | ✅ | ❌ | ✅ | ❌ | ❌ |
| balance | read | ✅ | ✅ | ❌ | ✅ | ❌ |
| reports | read | ✅ | ✅ | ❌ | ✅ | ❌ |
| settings | read | ✅ | ✅ | ❌ | ❌ | ❌ |
| settings | write | ✅ | ❌ | ❌ | ❌ | ❌ |
| users | read | ✅ | ❌ | ❌ | ❌ | ❌ |
| users | write | ✅ | ❌ | ❌ | ❌ | ❌ |

## Permission Check Implementation

### Go Middleware
```go
func RequirePermission(resource string, action string) middleware.Func {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := GetClaims(r.Context())
            if !hasPermission(claims.Permissions, resource, action) {
                respondError(w, 403, "permission_error",
                    "insufficient_permissions",
                    "You do not have permission to perform this action")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func hasPermission(userPerms []string, resource, action string) bool {
    // Check explicit permission
    if contains(userPerms, fmt.Sprintf("%s:%s", resource, action)) {
        return true
    }
    // Admin override (admin:*)
    if contains(userPerms, "admin:*") {
        return true
    }
    return false
}
```

## Data Isolation

### Multi-tenant Data Access
Every query MUST include merchant_id filter:
```go
// Correct: scoped to merchant
func (r *PaymentRepository) FindByID(ctx context.Context, merchantID, paymentID uuid.UUID) (*PaymentIntent, error) {
    return r.db.Get(ctx, `
        SELECT * FROM payment_intents
        WHERE id = $1 AND merchant_id = $2
    `, paymentID, merchantID)
}

// WRONG: no merchant scope (would return any merchant's payment)
func (r *PaymentRepository) FindByID(paymentID uuid.UUID) (*PaymentIntent, error) { ... }
```

### Admin Access
Admin users can view any merchant's data:
```go
func (r *PaymentRepository) AdminFindByID(ctx context.Context, paymentID uuid.UUID) (*PaymentIntent, error) {
    return r.db.Get(ctx, `
        SELECT pi.*, m.business_name as merchant_name
        FROM payment_intents pi
        JOIN merchants m ON m.id = pi.merchant_id
        WHERE pi.id = $1
    `, paymentID)
}
```

## Authorization Logic Summary

```
Request
  │
  ▼
1. Authenticate (who is this?)
  │
  ▼
2. Authorize (are they allowed?)
  │  ├── Check role-based permissions
  │  ├── Check resource ownership (merchant scope)
  │  └── Check additional constraints (IP whitelist, rate limit)
  │
  ▼
3. Audit (log the action)
```
