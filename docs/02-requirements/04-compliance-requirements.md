# Compliance Requirements

## PCI DSS (Payment Card Industry Data Security Standard)

### Scope
All systems that store, process, or transmit cardholder data are in scope.

### Applicable Requirements (v4.0)

| Requirement | How We Meet It | Evidence |
|-------------|----------------|----------|
| **1.** Install and maintain network security controls | Private VPC, security groups, network policies | Terraform config, network diagrams |
| **2.** Apply secure configurations | CIS benchmarked AMIs, hardened containers | Security hardening docs |
| **3.** Protect stored cardholder data | Tokenization — never store PAN after auth; AES-256 for tokens | Tokenization architecture |
| **4.** Protect cardholder data in transit | TLS 1.3 for all external + internal communication | TLS config, certificate management |
| **5.** Protect all systems from malware | Antivirus on bastion hosts, container scanning in CI | Trivy scan results |
| **6.** Develop and maintain secure systems | Secure coding guidelines, code review, dependency scanning | CI/CD pipeline |
| **7.** Restrict access by need-to-know | RBAC + ABAC, least privilege IAM roles | IAM policies, permission matrix |
| **8.** Identify and authenticate users | MFA for admin, strong password policy, session timeouts | Auth service config |
| **9.** Restrict physical access | Cloud infrastructure (AWS data centers) | AWS compliance reports |
| **10.** Log and monitor all access | Immutable audit logs, centralized logging, SIEM alerts | Audit log schema, Loki |
| **11.** Test security systems regularly | Weekly vuln scan, quarterly pen test, annual ASV scan | Scan reports |
| **12.** Support information security policy | Documented policies, employee training, incident response plan | Policy docs |

### Card Data Handling Rules
- **Never store** full PAN, CVV, or magnetic stripe data after authorization
- **Tokenize** card data immediately after processor returns token
- **Mask** PAN in logs: show only last 4 digits
- **Encrypt** any stored tokens with AES-256
- **Key rotation** every 90 days for encryption keys
- **Separate** card data network from corporate network

## PSD2 / Strong Customer Authentication (SCA)

### Requirements
- Two-factor authentication for electronic payments > €30
- Exemptions: low-value (<€30), recurring (same amount same merchant), trusted beneficiaries, corporate payments
- Dynamic linking: amount and payee linked to authentication

### Implementation
- Integrate 3D Secure 2.0 (EMV 3DS) for card transactions
- Support exemption requests with proper reason codes
- Maintain authentication data for audit (3 years)

## GDPR (General Data Protection Regulation)

| Requirement | Implementation |
|-------------|---------------|
| Lawful basis for processing | Contractual necessity (payment processing) |
| Consent management | Record consent for marketing, data sharing |
| Data minimization | Only collect data needed for payment |
| Right to access | API for merchants to export customer data |
| Right to erasure | Delete customer data within 30 days of request |
| Data portability | Export customer data in JSON/CSV |
| Breach notification | Notify DPA within 72 hours |
| DPO appointment | Designate Data Protection Officer |
| Data retention | Financial data: 7 years; non-financial: 2 years after last activity |
| Cross-border transfer | Standard Contractual Clauses for data leaving Bangladesh |

## AML / KYC (Anti-Money Laundering / Know Your Customer)

### Merchant KYC Requirements
- Business registration certificate
- Owner/government-issued ID (passport, national ID)
- Proof of business address (utility bill, bank statement)
- Ultimate Beneficial Owner (UBO) declaration for >25% ownership
- Sanctions screening against OFAC, UN, EU sanctions lists

### Transaction Monitoring
- Flag transactions > \$10,000 (single or aggregated in 24h)
- Flag rapid successive transactions (velocity)
- Flag transactions from high-risk jurisdictions
- Suspicious Activity Report (SAR) filing process

## SOC 2 Type II

### Trust Services Criteria
| Criteria | Implementation |
|----------|---------------|
| **Security** | Firewall, WAF, IDS/IPS, access controls, encryption |
| **Availability** | 99.99% uptime, DR plan, auto-scaling |
| **Processing Integrity** | Double-entry ledger, reconciliation, idempotency |
| **Confidentiality** | Encryption, access controls, data classification |
| **Privacy** | GDPR compliance, consent management, data retention |

## Regional Compliance
- **Bangladesh:** Bangladesh Bank guidelines for payment gateway operators, IT security framework
- **India (future):** RBI guidelines for payment aggregators, data localization
- **EU (future):** PSD2, GDPR as above
- **USA (future):** State-level money transmitter licenses (MTL), OFAC sanctions

## Compliance Calendar
| Activity | Frequency | Owner |
|----------|-----------|-------|
| PCI DSS scan | Quarterly + after infrastructure changes | Security team |
| Penetration test | Quarterly | Third-party |
| SOC 2 audit | Annually | External auditor |
| KYC review (new merchants) | Per onboarding | Compliance team |
| KYC refresh (existing) | Annually | Compliance team |
| Sanctions list update | Weekly | Automated |
| Data retention cleanup | Monthly | Automated |
| Incident response drill | Quarterly | SRE + Security |
| Employee security training | Bi-annually | HR + Security |
