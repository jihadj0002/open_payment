# Open Payment Gateway — Task Board

**Last Updated:** 2026-07-15
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

**Status:** TODO | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Build `/merchant/payments/[id]` page with payment timeline, amount breakdown, capture/refund/void actions, transaction list, webhook delivery logs.

---

### TASK-FE-PAGE-002: Refunds page

**Status:** TODO | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Build `/merchant/refunds` page with table, filters, and initiate refund modal.

---

### TASK-FE-PAGE-003: Customers page (list + detail)

**Status:** TODO | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Build `/merchant/customers` list page and `/merchant/customers/[id]` detail page with payment history and saved payment methods.

---

### TASK-FE-PAGE-004: Balance page

**Status:** TODO | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Build `/merchant/balance` page with balance breakdown (available/pending/reserve), transaction history, and payout history.

---

### TASK-FE-PAGE-005: Settings page (profile, users, security)

**Status:** TODO | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Build `/merchant/settings` with profile editing, team user management, bank accounts, and security settings (password change, MFA).

---

### TASK-FE-PAGE-006: Webhook Logs page

**Status:** TODO | **Priority:** MEDIUM | **Assignee:** frontend-engineer

**Description:**
Build `/merchant/webhooks/logs` page showing delivery attempt history with status, response codes, and payload preview.

---

### TASK-FE-PAGE-007: Reports page

**Status:** TODO | **Priority:** MEDIUM | **Assignee:** frontend-engineer

**Description:**
Build `/merchant/reports` with revenue, transaction, and settlement reports with date range filtering and CSV/PDF export.

---

### TASK-FE-PAGE-008: Route restructuring to /merchant/* paths

**Status:** TODO | **Priority:** MEDIUM | **Assignee:** frontend-engineer

**Description:**
Move existing pages from flat routes (`/dashboard`, `/payments`, `/api-keys`, `/webhooks`) into `(dashboard)/merchant/` route group. Create proper layout hierarchy with `(auth)` and `(dashboard)/merchant/` layouts.

---

### TASK-FE-PAGE-009: Forgot/Reset password pages

**Status:** TODO | **Priority:** LOW | **Assignee:** frontend-engineer

**Description:**
Build `/forgot-password` and `/reset-password` pages with email-based reset flow.

---

## Quick Stats
| Phase | Tasks | Status |
|-------|-------|--------|
| Phase 1: Fix Existing | 4 | 0/4 DONE |
| Phase 2: New Pages | 9 | 0/9 DONE |
