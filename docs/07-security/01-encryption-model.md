# Encryption Model

## Encryption in Transit

### TLS Configuration
- **Minimum version:** TLS 1.2
- **Preferred version:** TLS 1.3
- **Certificate:** AWS Certificate Manager (ACM) for public-facing ALBs
- **mTLS:** Internal service-to-service communication via Istio sidecar
- **Cipher suites (TLS 1.2):** ECDHE-ECDSA-AES128-GCM-SHA256, ECDHE-RSA-AES128-GCM-SHA256
- **Cipher suites (TLS 1.3):** TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384

### TLS Termination Points
```
Internet ──TLS 1.3──▶ Cloudflare ──TLS 1.3──▶ ALB ──TLS 1.2/1.3──▶ Service Pod
                                                              │
                                                    Istio Sidecar (mTLS)
                                                              │
                                                    ┌─────────┴─────────┐
                                                    │                   │
                                              Service A           Service B
                                              (mTLS)              (mTLS)
```

### Certificate Rotation
- ACM certificates auto-renew
- Istio mTLS certificates rotated every 24 hours
- Custom CA for internal certificates

## Encryption at Rest

### Database Encryption
- **RDS PostgreSQL:** AES-256 encryption at rest (AWS-managed keys)
- **ElastiCache Redis:** Encryption at rest enabled
- **MSK Kafka:** Encryption at rest enabled
- **S3:** Server-side encryption with AES-256 (SSE-S3) for KYC docs, SSE-KMS for sensitive data

### Application-Level Encryption

#### Sensitive Data Encryption
```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
)

type Encryptor struct {
    key []byte  // 256-bit key from KMS
}

func (e *Encryptor) Encrypt(plaintext []byte) (string, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }

    aead, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, aead.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(encoded string) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, err
    }

    aead, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := aead.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return aead.Open(nil, nonce, ciphertext, nil)
}
```

#### What Gets Encrypted at Application Level
| Data Field | Encryption | Key | Notes |
|------------|-----------|-----|-------|
| Card number (PAN) | Not stored | — | Replaced with processor token immediately |
| CVV | Not stored | — | Never persisted, passed through to processor |
| API secret keys | AES-256-GCM | Service key | Stored hashed (bcrypt) in DB, shown once |
| Processor credentials | AES-256-GCM | KMS key | Encrypted at service level |
| Webhook secrets | AES-256-GCM | Service key | Stored encrypted, decrypted for signing |
| PII (phone, email) | Not encrypted | — | Protected by access controls, not encryption |

## Key Management (AWS KMS)

### Key Hierarchy
```
AWS KMS Customer Master Key (CMK)
        │
        ├── Data Key (for S3 SSE-C)
        ├── Data Key (for RDS TDE)
        ├── Service Encryption Key (per service)
        │       ├── Payment Service Key
        │       ├── Auth Service Key
        │       └── Merchant Service Key
        └── API Signing Key
```

### Key Rotation
| Key Type | Rotation Period | Method |
|----------|----------------|--------|
| KMS CMK | Annually | AWS automatic rotation |
| Service encryption keys | Every 90 days | Generate new, re-encrypt on read |
| API signing keys | Every 180 days | Merchant manually rotates |
| TLS certificates | 13 months | ACM auto-renewal |

### HSM Integration
- AWS CloudHSM for PCI-scoped operations
- HSM stores: Card network private keys, Master encryption keys
- All cryptographic operations within HSM boundary

## Tokenization

### Card Tokenization Flow
```
1. Customer enters card number on checkout page
2. Frontend sends card to POST /v1/tokens (using publishable key)
3. Payment Service forwards to processor's tokenization endpoint
4. Processor returns a network token (Visa, Mastercard)
5. Gateway stores only: token, last4, brand, exp_month, exp_year
6. PAN is never stored on gateway infrastructure
```

### Token Storage
```sql
-- Only non-sensitive card data is stored
CREATE TABLE saved_cards (
    id UUID PRIMARY KEY,
    token VARCHAR(100) NOT NULL,     -- processor-issued token
    last4 VARCHAR(4),
    brand VARCHAR(20),
    exp_month INTEGER,
    exp_year INTEGER
    -- NO card_number, NO CVV
);
```
