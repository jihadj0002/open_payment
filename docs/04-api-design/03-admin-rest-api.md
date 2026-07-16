# Admin REST API Specification

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `internal/service/admin/routes.go`

## Authentication
All admin endpoints require JWT authentication with `admin` role.
Headers: `Authorization: Bearer <jwt_token>`

## Merchants

### GET /v1/admin/merchants — List Merchants

Query parameters:
- `status`: `pending`, `active`, `suspended`, `terminated`
- `verification_status`: `unverified`, `pending`, `verified`, `rejected`
- `q`: Search by business name, email
- `cursor`, `limit`: Pagination

### GET /v1/admin/merchants/:id — Retrieve Merchant

**Response:**
```json
{
  "id": "m_abc123",
  "business_name": "Acme Ltd",
  "email": "owner@acme.com",
  "phone": "+8801700000000",
  "status": "active",
  "verification_status": "verified",
  "country": "BD",
  "currency": "BDT",
  "fee_config": {
    "rate": 2.5,
    "fixed_fee": 300,
    "cross_border_rate": 3.0
  },
  "volume_this_month": 1500000,
  "created_at": "2025-01-15T10:00:00Z"
}
```

### PATCH /v1/admin/merchants/:id — Update Merchant

### POST /v1/admin/merchants/:id/approve — Approve Merchant

### POST /v1/admin/merchants/:id/suspend — Suspend Merchant

**Request:**
```json
{
  "reason": "suspicious_activity",
  "notify_merchant": true
}
```

### POST /v1/admin/merchants/:id/terminate — Terminate Merchant

### POST /v1/admin/merchants/:id/reset-api-keys — Reset Merchant API Keys

Regenerates API keys for a merchant. Their existing keys are revoked and new keys generated.

### GET /v1/admin/merchants/:id/stats — Get Merchant Stats

Returns aggregated usage statistics for a specific merchant (daily/weekly/monthly).

### POST /v1/admin/merchants — Create Merchant (Admin)

> **TODO: Not implemented** — Planned for future release.

### GET /v1/admin/merchants/:id/kyc — Get KYC Documents

> **TODO: Not implemented** — KYC document management endpoints pending.

**Response (spec):**
```json
{
  "documents": [
    {
      "type": "business_registration",
      "status": "verified",
      "url": "https://s3.amazonaws.com/...",
      "uploaded_at": "2025-01-15T10:30:00Z",
      "verified_at": "2025-01-16T09:00:00Z"
    }
  ]
}
```

### POST /v1/admin/merchants/:id/kyc/approve — Approve KYC

> **TODO: Not implemented**

### POST /v1/admin/merchants/:id/kyc/reject — Reject KYC

> **TODO: Not implemented**

## Transactions

### GET /v1/admin/transactions — List All Transactions

Query parameters:
- `merchant_id`: Filter by merchant
- `status`: Payment status
- `payment_method`: `card`, `wallet`, etc.
- `amount[gt]`, `amount[gte]`, `amount[lt]`, `amount[lte]`
- `created[gt]`, `created[gte]`, `created[lt]`, `created[lte]`

### GET /v1/admin/transactions/:id — Retrieve Transaction

## Settlements

> **Note:** Settlements are managed via `GET /v1/settlements` and `POST /v1/settlements` (merchant-facing). Admin-specific settlement endpoints are not yet separated.

## SLA Monitoring

### GET /v1/admin/sla — Get SLA Metrics

Returns p50/p95/p99 latency, availability, and error rate metrics.

> **Code Ref:** TASK-PROD-P2-009 — SLA Service implemented.

## Disputes

### GET /v1/admin/disputes — List Disputes

### GET /v1/admin/disputes/:id — Retrieve Dispute

### POST /v1/admin/disputes/:id/resolve — Resolve Dispute

**Request:**
```json
{
  "resolution": "merchant_won",
  "notes": "Evidence submitted was sufficient"
}
```

## Fee Configuration

### GET /v1/admin/fee_configs — List Fee Configs

### POST /v1/admin/fee_configs — Create Fee Config

**Request:**
```json
{
  "name": "Standard",
  "rate": 2.5,
  "fixed_fee": 300,
  "cross_border_rate": 3.0,
  "monthly_fee": 0,
  "min_monthly_volume": 0
}
```

### PATCH /v1/admin/fee_configs/:id — Update Fee Config

## System Configuration

### GET /v1/admin/config — Get System Config

### PATCH /v1/admin/config — Update System Config

```json
{
  "supported_currencies": ["BDT", "USD"],
  "supported_countries": ["BD", "US"],
  "supported_banks": ["DBBL", "BRAC", "City"],
  "maintenance_mode": false,
  "max_transaction_amount": 5000000
}
```

## Audit Logs

### GET /v1/admin/audit_logs — List Audit Logs

Query parameters:
- `actor_id`: Filter by user
- `action`: `merchant.approve`, `payment.refund`, etc.
- `resource_type`: `merchant`, `payment`, `user`
- `created[gt]`, `created[gte]`

**Response:**
```json
{
  "data": [{
    "id": "al_abc123",
    "actor_id": "u_admin1",
    "action": "merchant.approve",
    "resource_type": "merchant",
    "resource_id": "m_abc123",
    "details": { "previous_status": "pending", "new_status": "active" },
    "ip_address": "203.0.113.1",
    "created_at": "2025-01-16T09:00:00Z"
  }]
}
```
