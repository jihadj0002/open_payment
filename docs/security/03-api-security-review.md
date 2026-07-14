# API Security Review

## Authentication

- **JWT**: HS256 signed, 15-minute expiry, refresh via 30-day refresh token
- **API Keys**: SHA-256 hashed in DB, prefixed for type detection (`sk_` / `pk_`)
- Both use `Authorization: Bearer <token>` header
- No session cookies, no basic auth

## Authorization

- **Roles**: `merchant`, `admin`, `api_secret`, `api_publishable`
- Admin endpoints enforce `role == "admin"` (403 if not)
- API key scoping: publishable keys only allowed `token:create`

## Rate Limiting

- Not yet implemented at application level
- **Planned**: token-bucket per merchant, configurable tiers (100 / 1000 / 10000 req/min)
- **Headers**: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

## Input Validation

- All handlers validate JSON body
- `amount > 0`, currency whitelisted (BDT, USD)
- Pagination bounds enforced (`per_page` max 100)

## Idempotency

- Write endpoints accept `Idempotency-Key` header (UUID v4)
- Duplicate requests return existing resource without re-processing
- Idempotency keys stored in DB with unique partial index

## CORS

- Wide open during development (`*`)
- Locked down to specific origins in production

## Security Headers (Recommended)

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Content-Security-Policy: default-src 'self'
```

## Open Issues

- [ ] Implement rate limiting middleware with Redis
- [ ] Add request signing (HMAC) for high-value webhook callbacks
- [ ] Lock down CORS to merchant-specific origins
- [ ] Implement audit trail for all admin actions
- [ ] Add API key last-rotation notification
