# Authentication

## Authentication Methods

| Method | Use Case | Security Level |
|--------|----------|----------------|
| API Key (secret key) | Server-to-server API calls | High (with HMAC signing) |
| JWT (access + refresh tokens) | Dashboard users | High (short-lived + MFA) |
| OAuth2 + OIDC | Third-party integrations | Medium (delegated) |
| MFA (TOTP/Passkeys) | Admin dashboard | Very High |
| Session cookie | Merchant dashboard | Medium (HttpOnly + Secure) |

## API Key Authentication

### Key Generation
```go
func GenerateAPIKey() (string, string, error) {
    // Generate 32 random bytes
    keyBytes := make([]byte, 32)
    if _, err := rand.Read(keyBytes); err != nil {
        return "", "", err
    }

    // Encode to hex
    rawKey := hex.EncodeToString(keyBytes)

    // Prefix with type indicator
    apiKey := "sk_live_" + rawKey  // or "sk_test_" for sandbox

    // Hash for storage
    hash := sha256.Sum256([]byte(apiKey))
    hashedKey := hex.EncodeToString(hash[:])

    // Prefix for display (first 8 chars)
    prefix := apiKey[:8]

    // Last 4 chars
    last4 := apiKey[len(apiKey)-4:]

    return apiKey, hashedKey, nil
    // apiKey is returned ONCE to the user, never stored
    // hashedKey is stored in database
}
```

### Key Validation Path
```
1. Client sends: Authorization: Bearer sk_live_abc...
2. API Gateway extracts key, hashes with SHA-256
3. Looks up key_hash in database
4. Returns merchant_id, permissions, rate limit tier
5. Request proceeds with merchant context
```

## JWT Authentication

### Token Structure
```
Access Token (15 min TTL):
Header:  { "alg": "RS256", "kid": "key-id-1" }
Payload: {
  "sub": "user_uuid",
  "merchant_id": "m_abc123",
  "role": "owner",
  "permissions": ["payments:write", "refunds:write"],
  "iat": 1735689600,
  "exp": 1735690500,
  "jti": "unique_token_id"
}

Refresh Token (30 day TTL):
- Opaque random string (or structured JWT)
- Stored in database (hashed)
- One-time use (rotation on each refresh)
```

### Token Flow
```
┌─────────┐          ┌──────────┐          ┌──────────┐
│  Client │          │   Auth   │          │   Redis  │
│         │          │ Service  │          │          │
│  POST   │          │          │          │          │
│ /login  │─────────▶│          │          │          │
│         │          │  Verify  │          │          │
│         │          │  creds   │          │          │
│         │          │  Check   │          │          │
│         │          │  MFA     │          │          │
│         │          │  Create  │          │          │
│         │          │  tokens  │          │          │
│         │          │  Store   │─────────▶│          │
│         │          │  refresh │          │          │
│         │◀─────────│          │          │          │
│  Access │  tokens  │          │          │          │
│  + Ref  │          │          │          │          │
│         │          │          │          │          │
│  GET    │          │          │          │          │
│ /api/   │─────────▶│  Verify  │          │          │
│ payments│  (JWT)   │  JWT sig │          │          │
│         │          │  Check   │          │          │
│         │          │  black   │─────────▶│          │
│         │          │  list    │          │          │
│         │◀─────────│          │          │          │
│ 200 OK  │  (valid) │          │          │          │
└─────────┘          └──────────┘          └──────────┘
```

### JWT Blacklist
- On logout / password change / MFA disable: JWT ID (jti) added to Redis blacklist
- Blacklist TTL = JWT expiry
- Middleware checks blacklist before accepting any JWT

### Token Refresh
```go
func RefreshTokens(refreshToken string) (*Tokens, error) {
    hash := sha256.Sum256([]byte(refreshToken))
    session, err := db.GetSessionByRefreshHash(hash)
    if err != nil || session.ExpiresAt < time.Now() {
        return nil, ErrInvalidRefreshToken
    }

    // Rotate: invalidate old refresh, issue new pair
    db.DeleteSession(session.ID)
    oldAccessJTI := session.AccessTokenJTI
    redis.Set("jwt_blacklist:"+oldAccessJTI, time.Now(), TTL)

    return GenerateTokens(session.UserID, session.MerchantID, session.Role)
}
```

## MFA (Multi-Factor Authentication)

### TOTP (Time-based One-Time Password)
- Standard TOTP (RFC 6238)
- 30-second window
- 6-digit codes
- Recovery codes: 10 one-time use codes generated on setup

### Setup Flow
```
1. Admin enables MFA from settings
2. Server generates TOTP secret
3. QR code displayed (otpauth:// URI)
4. Admin scans with authenticator app
5. Admin enters code to verify
6. MFA marked as enabled, recovery codes shown
```

### Login with MFA
```
POST /auth/login
{ "email": "...", "password": "..." }
→ 200 { "requires_mfa": true, "mfa_token": "temp_token" }

POST /auth/login/mfa
{ "mfa_token": "temp_token", "code": "123456" }
→ 200 { "access_token": "...", "refresh_token": "..." }
```

## Password Policy
| Requirement | Value |
|-------------|-------|
| Minimum length | 12 characters |
| Complexity | At least 1 uppercase, 1 lowercase, 1 number, 1 special |
| History | No reuse of last 5 passwords |
| Max age | 90 days |
| Lockout | 5 attempts → 15 min lockout |
| Rate limit | 10 attempts per email per hour |

## Session Management
- **Active sessions:** Viewable in Settings → Security
- **Revoke session:** Invalidate specific refresh token
- **Revoke all:** Invalidate all sessions (password change, MFA toggle)
- **Idle timeout:** 30 minutes (dashboard), 2 hours (API)
- **Absolute timeout:** 12 hours (re-login required)
