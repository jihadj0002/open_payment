# Phase 10: bKash & Nagad Mobile Wallet Integration

**Started:** 2026-07-16
**Completed:** 2026-07-16
**Status:** COMPLETE

---

## Task Overview

Integrate bKash and Nagad mobile wallet payment systems into the payment gateway.

---

## Task List

### TASK-WALLET-001: Database migration — add return_url, cancel_url, client_secret columns

**Status:** DONE | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
The `PaymentIntent` model already had `ReturnURL`, `CancelURL`, and `ClientSecret` fields but the database schema didn't have these columns. Created a migration to add them.

**Files:**
- `internal/database/migrations/017_add_return_url_cancel_url.up.sql`
- `internal/database/migrations/017_add_return_url_cancel_url.down.sql`

---

### TASK-WALLET-002: bKash processor adapter — token management

**Status:** DONE | **Priority:** CRITICAL | **Assignee:** backend-engineer

**Description:**
Implemented the bKash token management layer: Grant Token API call with auto-refresh loop.

**Files:**
- `internal/processor/bkash/token.go` — `TokenManager` with grant/refresh/auto-refresh loop
- Config via env vars: `BKASH_APP_KEY`, `BKASH_APP_SECRET`, `BKASH_USERNAME`, `BKASH_PASSWORD`, `BKASH_BASE_URL`

---

### TASK-WALLET-003: bKash processor adapter — checkout payment

**Status:** DONE | **Priority:** CRITICAL | **Assignee:** backend-engineer

**Description:**
Implemented the bKash Checkout (URL-based) payment flow: Create Payment, Execute Payment, Query Payment.

**Files:**
- `internal/processor/bkash/checkout.go` — `Adapter` with Create/Execute/Query APIs
- `internal/processor/bkash/refund.go` — `RefundPayment` via bKash Refund API

---

### TASK-WALLET-004: bKash webhook/callback handler

**Status:** DONE | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
Implemented server-side handlers for bKash callbacks (GET with query params) and AWS SNS IPN webhooks.

**Files:**
- `internal/service/payment/callback_handler.go` — `HandleBkashCallback`, `HandleBkashWebhook`, `HandleBkashCallbackURL`

---

### TASK-WALLET-005: Nagad processor adapter — payment flow

**Status:** DONE | **Priority:** CRITICAL | **Assignee:** backend-engineer

**Description:**
Implemented the Nagad payment flow with RSA-signed initialization requests and completion.

**Files:**
- `internal/processor/nagad/nagad.go` — `Adapter` with Initialize/Complete payment, RSA signing

---

### TASK-WALLET-006: Nagad callback handler

**Status:** DONE | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
Implemented server-side handler for Nagad payment notifications (POST callbacks with transaction status).

**Files:**
- `internal/service/payment/callback_handler.go` — `HandleNagadCallback`

---

### TASK-WALLET-007: Extend payment service — wallet/bank_transfer routing

**Status:** DONE | **Priority:** CRITICAL | **Assignee:** backend-engineer

**Description:**
Updated `ProcessPayment` to handle `wallet` (bKash/Nagad) and `bank_transfer` payment methods. Added `WalletProvider` and `BankProvider` interfaces.

**Files:**
- `internal/service/payment/service.go` — `processWalletPayment`, `processBankTransfer`, `HandleWalletCallback`
- `internal/service/payment/payment.go` — `WalletProvider`, `BankProvider`, `WalletInitRequest/Response` interfaces
- `internal/service/payment/wallet_provider.go` — `BkashWalletProvider`, `NagadWalletProvider` implementations

---

### TASK-WALLET-008: Public checkout routes

**Status:** DONE | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
Created public API endpoints for checkout sessions, payment initiation, and success checking.

**Files:**
- `internal/service/payment/checkout_handler.go` — `HandleGetCheckoutSession`, `HandleInitiateCheckout`, `HandleCheckoutSuccess`
- `internal/service/payment/routes.go` — `RegisterCheckoutRoutes`, `RegisterBkashRoutes`, `RegisterNagadRoutes`
- `internal/service/payment/callback_handler.go` — Full callback handling for both providers

---

### TASK-WALLET-009: Hosted checkout frontend page

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Created customer-facing checkout page at `/checkout/[payment_intent_id]` with bKash and Nagad payment options.

**Files:**
- `web/src/app/checkout/[payment_intent_id]/page.tsx` — Full checkout page with bKash/Nagad selection, phone input, pay button

---

### TASK-WALLET-010: Payment success/failure frontend pages

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer

**Description:**
Created styled success, cancel, and error pages for the checkout flow.

**Files:**
- `web/src/app/checkout/[payment_intent_id]/success/page.tsx`
- `web/src/app/checkout/[payment_intent_id]/cancel/page.tsx`
- `web/src/app/checkout/[payment_intent_id]/error/page.tsx`

---

### TASK-WALLET-011: Wire adapters into main.go

**Status:** DONE | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
Updated `cmd/server/main.go` to initialize bKash token manager, bKash adapter, Nagad adapter, and register all checkout/callback routes.

**Files:**
- `cmd/server/main.go` — bKash/Nagad initialization, route registration
- `.env.example` — Added bKash/Nagad/PUBLIC_URL env vars

---

### TASK-WALLET-012: Merchant integration documentation

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** compliance-officer

**Description:**
Created merchant-facing documentation for bKash and Nagad integration with API examples.

**Files:**
- `docs/merchant/bkash-integration.md`
- `docs/merchant/nagad-integration.md`

---

### TASK-WALLET-013: Update state machine — add wallet_initiated status

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** backend-engineer

**Description:**
Added `wallet_initiated` status to the payment state machine for tracking the wallet payment flow.

**Files:**
- `internal/service/payment/state_machine.go` — Added `StatusWalletInitiated` and transitions
- `internal/service/payment/state_machine_test.go` — Updated tests

---

### TASK-WALLET-014: Integration tests for bKash/Nagad flows

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** qa-engineer

**Description:**
Updated all existing unit tests to match new code. All tests pass.

**Files:**
- `internal/service/payment/service_test.go` — Added `PaymentMethod:"card"` to test PaymentIntents
- `internal/service/payment/state_machine_test.go` — Added wallet_initiated transitions
- `internal/service/payment/integration_test.go` — Added `PaymentMethod:"card"` to PaymentIntents

---

## Files Created/Modified

### New Files Created (15)

1. `internal/database/migrations/017_add_return_url_cancel_url.up.sql`
2. `internal/database/migrations/017_add_return_url_cancel_url.down.sql`
3. `internal/processor/bkash/token.go`
4. `internal/processor/bkash/checkout.go`
5. `internal/processor/bkash/refund.go`
6. `internal/processor/nagad/nagad.go`
7. `internal/service/payment/wallet_provider.go`
8. `internal/service/payment/checkout_handler.go`
9. `internal/service/payment/callback_handler.go`
10. `web/src/app/checkout/[payment_intent_id]/page.tsx`
11. `web/src/app/checkout/[payment_intent_id]/success/page.tsx`
12. `web/src/app/checkout/[payment_intent_id]/cancel/page.tsx`
13. `web/src/app/checkout/[payment_intent_id]/error/page.tsx`
14. `docs/merchant/bkash-integration.md`
15. `docs/merchant/nagad-integration.md`
16. `docs/task-board/03-phase10-mobile-wallet-integration.md`

### Existing Files Modified (9)

1. `internal/processor/types.go` — Added `WalletInitResponse`, provider constants
2. `internal/service/payment/service.go` — Added wallet/bank routing, HandleWalletCallback, WalletProvider/BankProvider DI
3. `internal/service/payment/payment.go` — Extended Processor interface, added WalletProvider/BankProvider interfaces
4. `internal/service/payment/models.go` — Added CancelURL, RedirectURL, ProviderRef to PaymentIntent
5. `internal/service/payment/repository.go` — Added UpdatePaymentIntentProvider, UpdatePaymentIntentRedirect, updated queries for new columns
6. `internal/service/payment/routes.go` — Added RegisterCheckoutRoutes, RegisterBkashRoutes, RegisterNagadRoutes
7. `internal/service/payment/state_machine.go` — Added StatusWalletInitiated
8. `cmd/server/main.go` — Wired bKash/Nagad adapters, registered all new routes
9. `web/src/middleware.ts` — Added `/checkout` to public routes
10. `.env.example` — Added bKash, Nagad, and PUBLIC_URL env vars

---

## Verification

- `go build ./...` — ✅ Compiles
- `go test -tags=unit ./internal/service/payment/...` — ✅ All unit tests pass
- `npm run build` (frontend) — ✅ All pages compile, checkout routes visible
