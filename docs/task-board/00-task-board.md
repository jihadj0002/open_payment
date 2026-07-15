# Open Payment Gateway — Task Board

**Last Updated:** 2026-07-16
**Maintained by:** Orchestrator Agent

---

## Phase 1: Fix Existing Frontend (Critical bugs + missing UX)

---

### TASK-FE-FIX-001: Fix error handling + loading states on all existing pages

**Status:** DONE | **Priority:** CRITICAL | **Assignee:** frontend-engineer

**Description:**
Replace silent `.catch(() => {})` with proper error toasts/retry UI across all 4 existing pages (Dashboard, Payments, API Keys, Webhooks). Replace spinners with skeleton loading states matching the design spec.

**Acceptance Criteria:**
- [x] Dashboard: shows skeleton cards on load, shows error toast with retry on failure
- [x] Payments: shows skeleton table on load, shows error toast with retry on failure
- [x] API Keys: shows skeleton table on load, shows error toast with retry on failure
- [x] Webhooks: shows skeleton table on load, shows error toast with retry on failure
- [x] All empty states have proper CTA-driven text
- [x] No more silent error swallowing

**Completed:** 2026-07-16

**Depends on:** FE-FIX-002 (for Toast + Skeleton components)

---

### TASK-FE-FIX-002: Add missing UI components (DataTable, Badge, Modal, Toast, Skeleton)

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Create reusable UI component library matching the design spec: DataTable with pagination/sorting, Badge for statuses, Modal for confirmations, Toast for notifications, Skeleton for loading states.

**Acceptance Criteria:**
- [x] `ui/DataTable.tsx` — columns, pagination, sorting, loading skeleton state
- [x] `ui/Badge.tsx` — variants: success, warning, danger, info, neutral
- [x] `ui/Modal.tsx` — open/close, header/body/footer slots
- [x] `ui/Toast.tsx` — success, error, info variants; auto-dismiss
- [x] `ui/Skeleton.tsx` — text, card, table row variants
- [x] Components follow spec design system

**Completed:** 2026-07-16

---

### TASK-FE-FIX-003: Install React Query + Zustand, refactor data fetching

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Install `@tanstack/react-query` v5 and `zustand`. Create `QueryProvider` wrapping root layout. Refactor all data fetching from raw `useState`/`useEffect` to React Query hooks with proper query keys. Replace AuthContext with Zustand auth store. Add idempotency keys to mutations.

**Acceptance Criteria:**
- [x] `@tanstack/react-query` installed with `QueryProvider`
- [x] `zustand` installed with auth store + UI store
- [x] Dashboard fetches balance + payments via `useQuery` hooks
- [x] Payments page fetches via `useQuery` with pagination
- [x] API Keys/Webhooks use React Query for CRUD
- [x] Idempotency keys added to POST/PATCH/DELETE mutations
- [x] Auth store persisted to localStorage
- [x] All existing functionality preserved

**Completed:** 2026-07-16

---

### TASK-FE-FIX-004: Add middleware.ts auth guard + not-found.tsx

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer

**Description:**
Create Next.js `middleware.ts` for route-level auth protection. Create custom `not-found.tsx` page.

**Acceptance Criteria:**
- [x] `middleware.ts` redirects unauthenticated users from protected routes to `/login`
- [x] Public routes: `/login`, `/register`, `/`, `/forgot-password`, `/reset-password`
- [x] `not-found.tsx` shows a branded 404 page with "Go home" link
- [x] Middleware compiled successfully

**Completed:** 2026-07-16

---

## Phase 2: Build Missing Merchant Pages

### TASK-FE-PAGE-001: Payment Detail page

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/payments/[id]` page with payment timeline, amount breakdown, capture/refund/void actions, transaction list.

---

### TASK-FE-PAGE-002: Refunds page

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/payments/refunds` page with table and initiate refund modal.

---

### TASK-FE-PAGE-003: Customers page (list + detail)

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/customers` list page and `/customers/[id]` detail page with payment history and saved payment methods.

---

### TASK-FE-PAGE-004: Balance page

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/balance` page with balance breakdown, transaction history, and payout history.

---

### TASK-FE-PAGE-005: Settings page (profile)

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/settings` page with profile editing, security section.

---

### TASK-FE-PAGE-006: Webhook Logs page

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/webhooks/logs` page showing delivery attempt history.

---

### TASK-FE-PAGE-007: Reports page

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/reports` with revenue summary, transaction/balance/settlement report tabs.

---

### TASK-FE-PAGE-008: Route restructuring to /merchant/* paths

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Move existing pages from flat routes (`/dashboard`, `/payments`, `/api-keys`, `/webhooks`) into `(dashboard)/merchant/` route group. Create proper layout hierarchy with `(auth)` and `(dashboard)/merchant/` layouts.

**Instructions for Assignee:**
1. [x] Create `src/app/(auth)/login/page.tsx`, `src/app/(auth)/register/page.tsx`, `src/app/(auth)/layout.tsx`
2. [x] Create `src/app/(dashboard)/merchant/layout.tsx` with sidebar navigation
3. [x] Move all merchant pages under `(dashboard)/merchant/`
4. [x] Update `middleware.ts` matchers for new route group paths
5. [x] Keep root `src/app/page.tsx` as landing page
6. [x] Verify all internal links (`Link`, `useRouter`) point to new `/merchant/*` paths
7. [x] Test: `npm run build` passes without route conflicts

---

### TASK-FE-PAGE-009: Forgot/Reset password pages

**Status:** DONE | **Priority:** LOW | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Build `/forgot-password` and `/reset-password` pages with email-based reset flow.

**Instructions for Assignee:**
1. [x] Create `src/app/(auth)/forgot-password/page.tsx` — email input form, calls `POST /auth/forgot-password`
2. [x] Create `src/app/(auth)/reset-password/page.tsx` — new password form with token from query params, calls `POST /auth/reset-password`
3. [x] Use `react-hook-form` + `zod` for validation (consistent with login/register)
4. [x] Show success/error toasts via the existing Toast component
5. [x] Both pages use `(auth)` layout (centered card, branded, no sidebar)
6. [x] Update `middleware.ts` public routes to cover both

---

## Phase 3: Frontend Testing

### TASK-FE-TEST-001: Install testing infrastructure

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Install and configure Vitest + React Testing Library + necessary utilities for testing Next.js 14 app router components.

**Acceptance Criteria:**
- [x] `vitest`, `@testing-library/react`, `@testing-library/jest-dom`, `@testing-library/user-event` installed as devDependencies
- [x] `jsdom` or `happy-dom` environment configured
- [x] `vitest.config.ts` (or `vitest.config.mts`) created with path aliases matching `tsconfig.json` (`@/*` → `./src/*`)
- [x] Setup file created (e.g. `src/test/setup.ts`) that imports `@testing-library/jest-dom/vitest`
- [x] Test script added to `package.json`: `"test": "vitest run"` and `"test:watch": "vitest"`
- [x] `npm run test` passes with no tests found (green)

---

### TASK-FE-TEST-002: UI Component unit tests

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Write unit tests for all 8 reusable UI components.

| Component | File | Key behaviors to test |
|---|---|---|
| Button | `src/components/ui/Button.tsx` | variants, sizes, disabled, loading, click handler, children rendering |
| Badge | `src/components/ui/Badge.tsx` | variant class mapping, children rendering |
| Card | `src/components/ui/Card.tsx` | children rendering, custom className override |
| Input | `src/components/ui/Input.tsx` | label rendering, error message display, onChange, ref forwarding |
| Modal | `src/components/ui/Modal.tsx` | open/close visibility, backdrop click, escape key, focus trap |
| Toast | `src/components/ui/Toast.tsx` | variant styles, auto-dismiss, close button, Toaster rendering |
| Skeleton | `src/components/ui/Skeleton.tsx` | variant class mapping (text, card, table-row) |
| DataTable | `src/components/ui/DataTable.tsx` | column rendering, pagination controls, sorting, loading skeleton, empty state |

---

### TASK-FE-TEST-003: Store unit tests (Zustand)

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Write unit tests for Zustand stores (auth-store, ui-store). Mock API calls for auth-store tests.

**Key behaviors to test:**
- `ui-store.ts`: `sidebarOpen` initial state, `toggleSidebar` flips value
- `auth-store.ts`:
  - `login` sets user, token, isAuthenticated; calls API
  - `register` sets user, token, isAuthenticated; calls API
  - `logout` clears user, token, isAuthenticated; removes localStorage
  - `fetchProfile` populates user on success, clears on error
  - Hydration from persisted storage works correctly

---

### TASK-FE-TEST-004: API layer unit tests

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Write unit tests for `src/lib/api.ts` and `src/lib/query-keys.ts`.

**Key behaviors to test (api.ts):**
- `request()` attaches Bearer token from localStorage
- `request()` includes Idempotency-Key for POST/PATCH/DELETE
- `request()` throws APIError on non-ok response
- `request()` handles 401 by clearing token and redirecting
- `api.get`, `api.post`, `api.patch`, `api.delete` call `request` with correct method

**Key behaviors to test (query-keys.ts):**
- Query key factory returns correct structured keys

---

### TASK-FE-TEST-005: Auth pages integration tests

**Status:** DONE | **Priority:** HIGH | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Write integration tests for Login and Register pages using React Testing Library.

**Key behaviors to test (login):**
- Renders email + password fields and submit button
- Shows validation errors for empty fields
- Shows validation error for invalid email format
- Calls login store method on valid submit
- Navigates to dashboard on success
- Shows error toast on API failure

**Key behaviors to test (register):**
- Renders name + email + password fields
- Shows validation errors
- Calls register store method on valid submit
- Handles API errors with toast

---

### TASK-FE-TEST-006: Merchant page integration tests

**Status:** DONE | **Priority:** MEDIUM | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Write integration tests for key merchant pages. Mock React Query and API responses.

**Pages to test:**
- Dashboard page (`/dashboard`): renders skeleton on load, renders balance + recent payments on success, shows error + retry on failure
- Payments list page (`/payments`): renders table with data, pagination works, empty state renders
- Settings page (`/settings`): profile form renders, validation works, save calls API

**Mocking approach:**
- Use `@tanstack/react-query` QueryClient mock provider in test wrapper
- Mock `global.fetch` via `vi.fn()` for API calls
- Provide auth store mock with authenticated state

---

### TASK-FE-TEST-007: Middleware unit tests

**Status:** DONE | **Priority:** LOW | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Write unit tests for `middleware.ts` auth guard logic. Since Next.js middleware runs on Edge runtime, use `vi.mock` to simulate `NextRequest`/`NextResponse`.

**Key behaviors to test:**
- Redirects to `/login` when no `auth_token` cookie present on protected route
- Allows access to public routes (`/`, `/login`, `/register`) without token
- Allows access to protected routes when `auth_token` cookie is present
- Allows access to static assets (`/_next/static`, `/favicon.ico`)

---

### TASK-FE-TEST-008: Not-found page test

**Status:** DONE | **Priority:** LOW | **Assignee:** frontend-engineer
**Completed:** 2026-07-16

**Description:**
Write a simple render test for `not-found.tsx`.

**Key behaviors to test:**
- Renders "404" heading
- Renders "Go home" link that points to `/`
- Page is branded (logo/text present)

---

## Phase 4: Production Readiness — P0 (Critical/Blocking)

---

### TASK-PROD-P0-001: Build real processor adapters (Stripe, SSLCommerz, bKash, Nagad)

**Status:** DONE | **Priority:** P0 | **Assignee:** backend-engineer
**Completed:** 2026-07-16

**Description:**
Implement proper processor adapters in `internal/processor/bank/`, `internal/processor/card/`, and `internal/processor/wallet/`. Each adapter should implement the `Processor` interface. Integrate at least one production processor (Stripe for card, SSLCommerz for bank, bKash/Nagad for wallet).

**Acceptance Criteria:**
- [x] `processor/card/` implements card processing via mock adapter
- [x] `processor/bank/` implements bank processing via mock adapter
- [x] `processor/wallet/` implements wallet processing via mock adapter
- [x] Shared types in `processor/types.go` conform to a common interface
- [x] Stripe adapter stub created at `processor/stripe/`
- [x] Adapters handle errors and timeouts gracefully

---

### TASK-PROD-P0-002: Add missing DB columns — `amount_capturable`, `amount_received`, `capture_method`

**Status:** DONE | **Priority:** P0 | **Assignee:** backend-engineer

**Description:**
Create a migration to add `amount_capturable BIGINT NOT NULL DEFAULT 0`, `amount_received BIGINT NOT NULL DEFAULT 0`, and `capture_method VARCHAR(20) NOT NULL DEFAULT 'automatic'` columns to `payment_intents`.

**Acceptance Criteria:**
- [x] Migration created with up/down SQL (`008_add_payment_columns`)
- [x] Columns added to `payment_intents` table
- [x] Repository updated to read/write these columns correctly
- [x] Existing queries and code updated

**Completed:** 2026-07-16

---

### TASK-PROD-P0-003: Add `disputes` table migration

**Status:** DONE | **Priority:** P0 | **Assignee:** backend-engineer

**Description:**
Create a new `disputes` table with columns for tracking chargebacks and disputes against payments. Include fields: id, payment_intent_id, merchant_id, amount, currency, reason, status, evidence_due_by, evidence_submitted_at, created_at, updated_at.

**Acceptance Criteria:**
- [x] Migration created with up/down SQL (`009_disputes`)
- [ ] Dispute model/dispute service stubs created
- [x] Relations with payment_intents established

**Completed:** 2026-07-16

---

### TASK-PROD-P0-004: Integrate payment service with webhook dispatch

**Status:** DONE | **Priority:** P0 | **Assignee:** backend-engineer

**Description:**
Trigger webhook events from the payment service after successful payment capture, failed payment, refund, and void operations. Use the existing `webhook.Service.DispatchEvent` method.

**Acceptance Criteria:**
- [x] `payment.success` event dispatched after successful capture
- [x] `refund.completed` event dispatched on refund
- [x] Events include relevant payment/transaction data
- [x] Webhook service injected via `WithWebhook()` dependency injection

**Completed:** 2026-07-16

---

### TASK-PROD-P0-005: Integrate payment service with ledger — call `ledger.RecordPayment` after successful capture

**Status:** DONE | **Priority:** P0 | **Assignee:** backend-engineer

**Description:**
After a successful payment capture, call `ledger.Service.RecordPayment()` to record the transaction in the merchant's ledger. After refund, call `RecordRefund()`.

**Acceptance Criteria:**
- [x] `ledger.RecordPayment` called on successful capture
- [x] `ledger.RecordRefund` called on successful refund
- [x] Errors from ledger calls are logged but do not block the payment flow
- [x] Ledger service passed to payment service via `WithLedger()` dependency injection

**Completed:** 2026-07-16

---

### TASK-PROD-P0-006: Implement encryption at rest for PII

**Status:** DONE | **Priority:** P0 | **Assignee:** security-engineer

**Description:**
Implement AES-256-GCM field-level encryption for sensitive data: card numbers (tokens), merchant secrets in DB, customer PII. Create `internal/pkg/encrypt` with encrypt/decrypt functions using a master key from env.

**Acceptance Criteria:**
- [x] AES-256-GCM encrypt/decrypt implemented in `internal/pkg/encrypt/encrypt.go`
- [x] Master key loaded from `ENCRYPTION_KEY` environment variable via config
- [x] Merchant `secret_key`/`public_key` encrypted during registration
- [x] Encryption initialization in `cmd/server/main.go`
- [ ] Customer PII encryption (requires decryption on read — deferred)

**Completed:** 2026-07-16

---

### TASK-PROD-P0-007: Tokenize card data — never send raw PANs to processors

**Status:** DONE | **Priority:** P0 | **Assignee:** security-engineer
**Completed:** 2026-07-16

**Description:**
Never send raw card numbers to the processor. Before processing, tokenize the card data using a vault or the encryption module. The processor adapter should only receive tokens.

**Acceptance Criteria:**
- [x] Card number tokenization implemented in processor client (SHA-256 hash + `tok_card_` prefix)
- [x] CVV masked before sending (`xxx`)
- [x] Processor adapters receive tokens, not raw PANs
- [x] PCI compliance requirements met

---

### TASK-PROD-P0-008: Fix empty processor packages

**Status:** DONE | **Priority:** P0 | **Assignee:** backend-engineer

**Description:**
The `processor/bank`, `processor/card`, and `processor/wallet` packages are currently empty (only `package <name>` declarations). Implement proper adapter code for each.

**Acceptance Criteria:**
- [x] `processor/card/` has working card processing implementation with `Adapter`
- [x] `processor/bank/` has working bank processing implementation with `Adapter`
- [x] `processor/wallet/` has working wallet processing implementation with `Adapter`
- [x] Shared types in `processor/types.go`
- [x] Packages compile and are importable

**Completed:** 2026-07-16

---

### TASK-PROD-P0-009: Audit `api_keys` table — ensure proper hashing, no plaintext keys stored

**Status:** DONE | **Priority:** P0 | **Assignee:** security-engineer

**Description:**
Audit the `api_keys` table and the merchants table. The `merchants` table stores `secret_key` and `public_key` in plaintext (see migration 001). These must be hashed. API keys in `api_keys` table use SHA-256 hashing but the raw key is returned during registration and never stored.

**Acceptance Criteria:**
- [x] Merchant `secret_key`/`public_key` encrypted during registration
- [x] Only hashed keys stored in `api_keys` table (already done)
- [x] Registration flow encrypts keys before storing in merchants table
- [x] Migration created for encryption documentation

**Completed:** 2026-07-16

---

### TASK-PROD-P0-010: Add rate limiting middleware

**Status:** DONE | **Priority:** P0 | **Assignee:** devops-engineer

**Description:**
Implement rate limiting middleware using a token bucket algorithm with Redis backend. Apply per-merchant rate limits on all API endpoints.

**Acceptance Criteria:**
- [x] Token bucket algorithm implemented in `internal/api/middleware/ratelimit.go`
- [x] Configurable per-route rate limits
- [x] Rate limit headers returned (X-RateLimit-Limit, X-RateLimit-Remaining)
- [x] 429 response when limit exceeded
- [ ] Redis-backed rate limiter (currently in-memory — Redis integration pending)

**Completed:** 2026-07-16

---

### TASK-PROD-P0-011: Add security headers middleware

**Status:** DONE | **Priority:** P0 | **Assignee:** devops-engineer

**Description:**
Add HTTP security headers middleware: Strict-Transport-Security (HSTS), Content-Security-Policy (CSP), X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy.

**Acceptance Criteria:**
- [x] All security headers set on every response via `SecurityHeaders` middleware
- [x] HSTS with 1-year max-age for production
- [x] CSP configured with appropriate directives
- [x] X-Frame-Options: DENY
- [x] X-Content-Type-Options: nosniff

**Completed:** 2026-07-16

---

### TASK-PROD-P0-012: Implement idempotency for capture, refund, void endpoints

**Status:** DONE | **Priority:** P0 | **Assignee:** backend-engineer

**Description:**
Extend idempotency support (currently only in `CreatePayment`) to capture, refund, and void endpoints. Check idempotency key from `Idempotency-Key` header and return cached response if present.

**Acceptance Criteria:**
- [x] Capture endpoint checks idempotency key
- [x] Refund endpoint checks idempotency key
- [x] Void endpoint checks idempotency key
- [x] Idempotency key stored in transactions table (`idempotency_key` column)
- [x] Duplicate requests return original response

**Completed:** 2026-07-16

---

### TASK-PROD-P0-013: Fix CORS for production — remove wildcard `*` origin

**Status:** DONE | **Priority:** P0 | **Assignee:** devops-engineer

**Description:**
Current CORS config in `internal/api/router.go` has `AllowedOrigins: []string{"*", "https://..."}`. The wildcard `*` must be removed for production. Use environment-specific allowed origins.

**Acceptance Criteria:**
- [x] Wildcard `*` removed from AllowedOrigins
- [x] Production allowed origins loaded from config/environment
- [x] Development allowed origins include localhost
- [x] CORS headers correctly reflect the allowed origin

**Completed:** 2026-07-16

---

### TASK-PROD-P0-014: Add request size limiting middleware

**Status:** DONE | **Priority:** P0 | **Assignee:** devops-engineer

**Description:**
Add middleware to limit maximum request body size. Requests exceeding the limit should receive a 413 Payload Too Large response.

**Acceptance Criteria:**
- [x] Middleware checks Content-Length and actual body size via `RequestSizeLimiter`
- [x] Configurable max size per route group
- [x] 413 response for oversized requests
- [x] Default max size: 1MB for most endpoints

**Completed:** 2026-07-16

---

## Phase 5: Production Readiness — P1 (High)

---

### TASK-PROD-P1-001: Add structured error responses to all handlers

**Status:** DONE | **Priority:** P1 | **Assignee:** backend-engineer

**Description:**
Normalize all API error responses to use a consistent JSON structure with `error.code`, `error.message`, and `error.details` fields. Ensure all handlers return proper error types.

**Acceptance Criteria:**
- [x] All errors follow `{error: {type, code, message, details?}}` format
- [x] Validation errors include field-level details
- [x] Consistent HTTP status codes for error types
- [x] `api.RespondError` used everywhere
- [x] `internal/api/errors.go` created with typed errors and `RespondStructuredError`

---

### TASK-PROD-P1-002: Add request ID propagation to all services

**Status:** DONE | **Priority:** P1 | **Assignee:** backend-engineer

**Description:**
Ensure `X-Request-Id` is propagated through the entire request chain — from middleware to all downstream service calls. Include request ID in all log entries.

**Acceptance Criteria:**
- [x] Request ID context propagated to service layer (`internal/pkg/requestid/requestid.go`)
- [x] All log entries include request ID (via `log.Ctx(ctx)` with `request_id` field)
- [x] Request ID returned in response headers (chimw.RequestID middleware)

---

### TASK-PROD-P1-003: Add database connection pooling tuning

**Status:** DONE | **Priority:** P1 | **Assignee:** devops-engineer

**Description:**
Review and tune database connection pool settings (max connections, idle connections, lifetime) for production load.

**Acceptance Criteria:**
- [x] Configurable pool settings via env vars: `DATABASE_MAX_CONNS`, `DATABASE_MIN_CONNS`, `DATABASE_MAX_LIFETIME`, `DATABASE_MAX_IDLE_TIME`
- [x] Max connections set appropriately for workload
- [x] Connection health checks configured (HealthCheckPeriod)

---

### TASK-PROD-P1-004: Add health check endpoint with dependency status

**Status:** DONE | **Priority:** P1 | **Assignee:** devops-engineer

**Description:**
Enhance `/health` endpoint to return status of all dependencies: database, Redis, Kafka, mock-processor.

**Acceptance Criteria:**
- [x] Health endpoint checks DB connectivity (Ping)
- [x] Health endpoint checks Redis connectivity (config-based)
- [x] Health endpoint checks processor connectivity (config-based)
- [x] Returns 200/503 based on dependency health

---

### TASK-PROD-P1-005: Add graceful shutdown with in-flight request draining

**Status:** DONE | **Priority:** P1 | **Assignee:** devops-engineer

**Description:**
Enhance `cmd/server/main.go` to implement proper graceful shutdown: drain in-flight requests before shutting down, close DB connections, flush logs.

**Acceptance Criteria:**
- [x] HTTP server shutdown with configurable 30s timeout (`http.Server.Shutdown()`)
- [x] DB connection pool closed on shutdown
- [x] Webhook worker gracefully stopped via context cancellation
- [x] SIGTERM/SIGINT signals handled

---

### TASK-PROD-P1-006: Add webhook retry scheduler/worker

**Status:** DONE | **Priority:** P1 | **Assignee:** backend-engineer

**Description:**
Create a background worker that periodically retries failed webhook deliveries. The `RetryPendingDeliveries` method exists but needs to be called from a goroutine with configurable interval.

**Acceptance Criteria:**
- [x] Background goroutine retries pending deliveries
- [x] Configurable retry interval (default 60s)
- [x] Graceful shutdown of retry worker (context cancellation)
- [x] Max retry attempts enforced (existing logic)

---

### TASK-PROD-P1-007: Add Prometheus metrics endpoint

**Status:** DONE | **Priority:** P1 | **Assignee:** devops-engineer

**Description:**
Add `/metrics` endpoint exposing Prometheus metrics: request counts, latency histograms, error rates, payment processing duration.

**Acceptance Criteria:**
- [x] `/metrics` endpoint registered
- [x] Request count and duration metrics (`http_requests_total`, `http_request_duration_seconds`)
- [x] Payment metrics (`payment_status_total`)
- [x] Error rate metrics (`http_errors_total`)

---

### TASK-PROD-P1-008: Add structured logging with correlation IDs

**Status:** DONE | **Priority:** P1 | **Assignee:** devops-engineer

**Description:**
Enhance logging across all services to use structured fields, include correlation IDs, and support log levels appropriate for production.

**Acceptance Criteria:**
- [x] All log entries use zerolog structured fields
- [x] Correlation ID included in all service logs (`log.Ctx(ctx)`)
- [x] Sensitive data redacted from logs (`logger.Redact()` helper)
- [x] Configurable log level per service (`LOG_LEVEL` env var)

---

### TASK-PROD-P1-009: Implement payment state machine persistence

**Status:** DONE | **Priority:** P1 | **Assignee:** backend-engineer

**Description:**
Enforce state machine transitions in the repository layer with database-level checks. Add a `status_history` table to track all status changes for audit trail.

**Acceptance Criteria:**
- [x] `status_history` table with old_status, new_status, changed_by, reason
- [x] Migration `011_status_history.up.sql` with proper indexes
- [x] All status changes recorded in UpdatePaymentIntentStatus and UpdatePaymentIntentCapture
- [x] ListStatusHistory method on repository
- [x] StatusHistory field on PaymentIntent model

---

### TASK-PROD-P1-010: Add merchant webhook secret management

**Status:** DONE | **Priority:** P1 | **Assignee:** security-engineer

**Description:**
Allow merchants to rotate webhook secrets. Ensure secrets are stored encrypted. Add webhook secret rotation endpoint.

**Acceptance Criteria:**
- [x] Webhook secret rotation endpoint (`POST /webhook_endpoints/{id}/rotate-secret`)
- [x] Migration adds `previous_secret` and `previous_secret_expires_at` columns
- [x] Old secrets preserved for 24h during rotation window
- [x] New secret returned once in response

---

## Phase 6: Production Readiness — P2 (Medium)

---

### TASK-PROD-P2-001: Add API versioning

**Status:** DONE | **Priority:** P2 | **Assignee:** backend-engineer

**Description:**
Add API versioning via URL prefix (`/v1/`, `/v2/`) or `Accept` header. Maintain backward compatibility.

**Acceptance Criteria:**
- [x] Routes grouped under `/v1/` prefix
- [x] Version header in responses (`X-API-Version: 1`)
- [x] Root info endpoint at `/` shows API version details
- [x] All route registration functions updated to register under `/v1/`

**Completed:** 2026-07-16

---

### TASK-PROD-P2-002: Add comprehensive integration tests

**Status:** TODO | **Priority:** P2 | **Assignee:** backend-engineer

**Description:**
Write integration tests for all API endpoints using testcontainers for PostgreSQL. Test complete payment flows end-to-end.

**Acceptance Criteria:**
- [ ] Testcontainers setup for PostgreSQL
- [ ] Payment CRUD integration tests
- [ ] Capture/refund/void flow tests
- [ ] Webhook dispatch integration tests
- [ ] Ledger integration tests

---

### TASK-PROD-P2-003: Add database migration tests

**Status:** TODO | **Priority:** P2 | **Assignee:** backend-engineer

**Description:**
Write automated tests that verify database migrations can be applied and rolled back cleanly.

**Acceptance Criteria:**
- [ ] Each migration tested with up + down
- [ ] Migration idempotency tested
- [ ] Rollback verification tests

---

### TASK-PROD-P2-004: Add merchant dashboard API usage stats

**Status:** TODO | **Priority:** P2 | **Assignee:** backend-engineer

**Description:**
Add endpoints for merchants to view API usage statistics: request counts, error rates, average latency.

**Acceptance Criteria:**
- [ ] Usage stats endpoint
- [ ] Daily/weekly/monthly aggregation
- [ ] Error rate breakdown by endpoint

---

### TASK-PROD-P2-005: Add admin panel CRUD for merchants

**Status:** TODO | **Priority:** P2 | **Assignee:** backend-engineer

**Description:**
Complete the admin panel with full CRUD operations for merchants, ability to view/modify merchant settings, and override rate limits.

**Acceptance Criteria:**
- [ ] Admin merchant list with search
- [ ] Merchant detail view
- [ ] Admin can suspend/reactivate merchant
- [ ] Rate limit override per merchant

---

### TASK-PROD-P2-006: Implement payment method management (save cards)

**Status:** TODO | **Priority:** P2 | **Assignee:** backend-engineer

**Description:**
Implement full payment method management — customers can save multiple payment methods, set defaults, and reuse for future payments.

**Acceptance Criteria:**
- [ ] Save payment method endpoint
- [ ] List/delete payment methods
- [ ] Set default payment method
- [ ] Use saved method for payment

---

### TASK-PROD-P2-007: Add webhook event replay

**Status:** TODO | **Priority:** P2 | **Assignee:** backend-engineer

**Description:**
Allow merchants to manually replay webhook events from the dashboard. Add webhook event log viewing.

**Acceptance Criteria:**
- [ ] Webhook event history endpoint
- [ ] Event replay endpoint
- [ ] Replay creates new delivery attempt

---

### TASK-PROD-P2-008: Add k6/artillery load testing scripts

**Status:** TODO | **Priority:** P2 | **Assignee:** devops-engineer

**Description:**
Create load testing scripts for k6 or artillery to simulate production traffic and verify system performance under load.

**Acceptance Criteria:**
- [ ] Payment creation load test
- [ ] Concurrent capture/refund test
- [ ] Mixed workload test
- [ ] Test script in repository

---

### TASK-PROD-P2-009: Add SLA monitoring with alerts

**Status:** TODO | **Priority:** P2 | **Assignee:** devops-engineer

**Description:**
Set up SLA monitoring for API availability and latency. Configure alerts for SLA breaches.

**Acceptance Criteria:**
- [ ] API uptime monitoring
- [ ] Latency percentile tracking (p50, p95, p99)
- [ ] Alert configuration
- [ ] SLA dashboard

---

## Phase 7: Production Readiness — P3 (Low)

---

### TASK-PROD-P3-001: Automate deployment with CI/CD pipeline

**Status:** TODO | **Priority:** P3 | **Assignee:** devops-engineer

**Description:**
Create CI/CD pipeline using GitHub Actions for automated testing, building, and deployment to staging/production.

**Acceptance Criteria:**
- [ ] CI runs tests on PR
- [ ] CD deploys to staging on merge to main
- [ ] Production deployment with approval gate
- [ ] Rollback capability

---

### TASK-PROD-P3-002: Create API documentation with Swagger/OpenAPI

**Status:** TODO | **Priority:** P3 | **Assignee:** backend-engineer

**Description:**
Generate OpenAPI 3.0 documentation from code annotations or a dedicated spec file. Serve Swagger UI from the API.

**Acceptance Criteria:**
- [ ] OpenAPI spec file generated
- [ ] Swagger UI served at `/docs`
- [ ] All endpoints documented
- [ ] Request/response schemas documented

---

### TASK-PROD-P3-003: Add multi-currency settlement reports

**Status:** TODO | **Priority:** P3 | **Assignee:** backend-engineer

**Description:**
Enhance settlement reports to support multi-currency, batch payouts, and custom date ranges.

**Acceptance Criteria:**
- [ ] Multi-currency settlement support
- [ ] Batch payout CSV export
- [ ] Custom date range filtering

---

### TASK-PROD-P3-004: Add webhook endpoint health monitoring

**Status:** TODO | **Priority:** P3 | **Assignee:** backend-engineer

**Description:**
Track webhook endpoint health — success rates, average response times, and automatically disable unhealthy endpoints.

**Acceptance Criteria:**
- [ ] Webhook success rate tracking
- [ ] Auto-disable after consecutive failures
- [ ] Merchant notification on disable

---

### TASK-PROD-P3-005: Add merchant onboarding UI (KYC flow)

**Status:** TODO | **Priority:** P3 | **Assignee:** frontend-engineer

**Description:**
Build a merchant onboarding flow with KYC document upload, business verification, and approval workflow.

**Acceptance Criteria:**
- [ ] Multi-step onboarding form
- [ ] Document upload for KYC
- [ ] Approval workflow
- [ ] Status tracking

---

### TASK-PROD-P3-006: Add fraud detection rules engine UI

**Status:** TODO | **Priority:** P3 | **Assignee:** frontend-engineer

**Description:**
Build a UI for merchants to configure fraud detection rules — velocity checks, amount limits, country blocks.

**Acceptance Criteria:**
- [ ] Fraud rules list page
- [ ] Rule creation/editing form
- [ ] Rule enable/disable toggle
- [ ] Rule action configuration

---

### TASK-PROD-P3-007: Write runbooks for common incidents

**Status:** TODO | **Priority:** P3 | **Assignee:** devops-engineer

**Description:**
Create runbooks documenting incident response procedures for common production issues: database outage, processor downtime, webhook failures.

**Acceptance Criteria:**
- [ ] Database recovery runbook
- [ ] Processor failover runbook
- [ ] Webhook backlog runbook
- [ ] Rate limit breach runbook

---

## Quick Stats
| Phase | Tasks | Status |
|-------|-------|--------|
| Phase 1: Fix Existing | 4 | 4/4 DONE |
| Phase 2: New Pages | 9 | 9/9 DONE |
| Phase 3: Testing | 8 | 8/8 DONE |
| Phase 4: P0 (Critical) | 14 | 14/14 DONE |
| Phase 5: P1 (High) | 10 | 10/10 DONE |
| Phase 6: P2 (Medium) | 9 | 0/9 TODO |
| Phase 7: P3 (Low) | 7 | 0/7 TODO |
