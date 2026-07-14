# Business Requirements

## Market Need
Online payment volume globally exceeds \$10T annually. Merchants in emerging markets (South Asia, Southeast Asia, Africa) need local payment method support (bKash, Nagad, UPI, etc.) combined with international card processing — a gap many global gateways don't fill.

## Revenue Model
| Revenue Stream | Description | Expected Margin |
|----------------|-------------|-----------------|
| Transaction Fee | Per-transaction percentage (2.5% + \$0.30 typical) | Primary revenue |
| Monthly Subscription | Tiered plans based on volume (Starter, Business, Enterprise) | Recurring |
| Cross-border Fee | Additional 1–2% for international transactions | High margin |
| Payout Fee | Fixed fee per settlement payout | Low volume, stable |
| Value-added Services | Fraud prevention, analytics, invoicing | Upsell |

## Target Merchant Segments
| Segment | Volume/Month | Needs | Priority |
|---------|-------------|-------|----------|
| Small e-commerce | < 1,000 transactions | Simple checkout, low fees | High |
| SaaS platforms | 1,000–50,000 | Subscriptions, webhooks, API-first | High |
| Marketplaces | 10,000–500,000 | Split payments, seller onboarding | Medium |
| Enterprise | 500,000+ | Multi-currency, dedicated support, SLA | Medium |

## Success Metrics (First Year)
- **Merchants onboarded:** 500+ active merchants
- **Transaction volume:** \$5M+ processed monthly
- **Payment success rate:** >92%
- **Uptime:** 99.9%+
- **Average API response:** <200ms p95
- **Fraud loss rate:** <0.5% of volume
- **Developer NPS:** >40

## Competitive Landscape
| Competitor | Strengths | Weaknesses |
|------------|-----------|------------|
| Stripe | Developer experience, global reach | Expensive for high volume, limited local methods |
| SSLCommerz | Bangladesh market leader | Limited international support |
| Paypal | Brand trust, buyer protection | High fees, poor API |
| Adyen | Enterprise, single platform | Complex onboarding, minimum volumes |
| Razorpay | India market, feature-rich | India-only |

## Monetization Strategy (Phased)
- **Phase 1 (MVP):** Transaction fee only — attract merchants with simple pricing
- **Phase 2:** Introduce tiered subscriptions + cross-border fees
- **Phase 3:** Value-added services (fraud analytics, invoicing, reporting)
