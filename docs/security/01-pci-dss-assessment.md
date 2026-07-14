# PCI DSS Self-Assessment Questionnaire Reference

## Merchant Level

**Level 4** (< 20,000 transactions/year)

## SAQ Type

**SAQ A** — Card-not-present merchants; all cardholder data handled by third-party processor.

## Requirements Coverage

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 1 | Install and maintain firewall | PLANNED | AWS Security Groups + K8s Network Policies |
| 2 | Secure passwords | COMPLIANT | bcrypt for merchant passwords, SHA-256 for API keys |
| 3 | Protect stored cardholder data | PARTIAL | Only last4 + fingerprint stored; full PAN never touches our servers (tokenized via mock processor) |
| 4 | Encrypt transmission | COMPLIANT | HTTPS enforced in production (ALB termination) |
| 5 | Anti-virus | N/A | Linux containers — not applicable |
| 6 | Secure systems | IN PROGRESS | Regular security updates in Docker builds |
| 7 | Access control | COMPLIANT | JWT + API key auth, role-based permissions |
| 8 | Unique IDs | COMPLIANT | bcrypt passwords, unique API keys per merchant |
| 9 | Physical security | N/A | AWS data centers |
| 10 | Logging | PARTIAL | Audit log table exists; needs SIEM integration |
| 11 | Regular testing | PLANNED | Weekly Trivy scans via GitHub Actions |
| 12 | Policy | IN PROGRESS | This document |

## Security Controls Implemented

- **Password hashing**: bcrypt with cost factor 10
- **API key hashing**: SHA-256 (one-way; key shown only once at creation)
- **JWT signing**: HMAC-SHA256 with configurable secret
- **Card data**: Only last4 digits stored; full card tokenized via processor; SHA-256 fingerprint for deduplication
- **No plaintext secrets in code**: all secrets injected via environment variables

## Gaps to Address

- [ ] Implement key rotation for JWT secret
- [ ] Add rate limiting middleware
- [ ] Set up SIEM for audit log analysis
- [ ] Implement WAF (AWS WAF on ALB)
- [ ] Penetration testing before going live
