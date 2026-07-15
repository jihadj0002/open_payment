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

## Quick Stats
| Phase | Tasks | Status |
|-------|-------|--------|
| Phase 1: Fix Existing | 4 | 4/4 DONE |
| Phase 2: New Pages | 9 | 9/9 DONE |
| Phase 3: Testing | 8 | 8/8 DONE |
