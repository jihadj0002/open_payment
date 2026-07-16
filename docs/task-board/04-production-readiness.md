# Open Payment Gateway — Production Readiness Tasks

**Last Updated:** 2026-07-16
**Maintained by:** Orchestrator Agent

---

## Phase 11: Production Readiness Remaining Fixes (from Audit)

**Status:** 0/10 DONE

---

### PROD-001: Fix CORS — add Railway app domain + make origins configurable via env var

**Status:** PENDING | **Priority:** CRITICAL | **Assignee:** backend-engineer

**Description:**
Add `openpaymentweb-production.up.railway.app` to the allowed CORS origins. Make CORS origins fully configurable via an environment variable (e.g. `CORS_ALLOWED_ORIGINS`) so they can be set per environment without code changes.

**Acceptance Criteria:**
- [ ] `openpaymentweb-production.up.railway.app` added to allowed origins in `internal/api/router.go`
- [ ] CORS origins loaded from `config.Config.CORSAllowedOrigins` (comma-separated env var)
- [ ] Dev default includes `http://localhost:3000`
- [ ] Wildcard `*` NOT used in production
- [ ] CORS headers correctly echo the requesting origin when allowed

---

### PROD-002: Fix payments list response format mismatch

**Status:** PENDING | **Priority:** CRITICAL | **Assignee:** backend-engineer

**Description:**
Backend endpoint `GET /v1/payments` returns response wrapped in `{data: [...], pagination: {...}}` but the frontend's `query-keys.ts` factory and API consumption expect `{payments: [...], pagination: {...}}`. Align the response format so both sides agree.

**Acceptance Criteria:**
- [ ] Backend returns `{payments: [...], pagination: {...}}` (preferred — change backend, not frontend)
- [ ] OR frontend updated to consume `{data: [...], pagination: {...}}` and all existing pages work
- [ ] Response change does NOT break settlements, customers, or other list endpoints
- [ ] Integration test verifies response shape

---

### PROD-003: Fix settlements response format (backend returns bare array, frontend expects paginated object)

**Status:** PENDING | **Priority:** CRITICAL | **Assignee:** backend-engineer

**Description:**
Backend `GET /v1/settlements` returns a bare JSON array `[{...}, {...}]` but the frontend expects a paginated object `{settlements: [...], pagination: {...}}`. Standardize the response format.

**Acceptance Criteria:**
- [ ] Backend wraps settlements response in `{settlements: [...], pagination: {...}}`
- [ ] Pagination metadata included (page, limit, total, total_pages)
- [ ] Frontend settlements page updated if needed to match
- [ ] No regression to other paginated endpoints

---

### PROD-004: Add Next.js Dockerfile for frontend standalone deployment

**Status:** PENDING | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Create a `Dockerfile` in the frontend directory using Next.js standalone output mode. This enables lightweight container images for production deployment.

**Acceptance Criteria:**
- [ ] Dockerfile created at `frontend/Dockerfile` (or `web/Dockerfile` depending on directory structure)
- [ ] Uses multi-stage build (deps → build → production)
- [ ] `output: 'standalone'` configured in `next.config.js`/`next.config.mjs`
- [ ] `.dockerignore` excludes `node_modules`, `.next`, `git`, `tests`
- [ ] Image copies `public/`, `.next/static/`, and `standalone` output
- [ ] Port exposed (default 3000), health check optional
- [ ] Build verified (`docker build` succeeds)

---

### PROD-005: Fix inverted Redis health check in router.go

**Status:** PENDING | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
The Redis health check in `internal/api/router.go` has an inverted condition — it reports Redis as healthy when it's actually down, and vice versa. Fix the boolean logic.

**Acceptance Criteria:**
- [ ] Locate the Redis ping/health check in `router.go` (likely in the health endpoint handler)
- [ ] Fix the inverted condition so healthy Redis returns healthy and unhealthy Redis returns unhealthy
- [ ] Verify with unit test or manual test

---

### PROD-006: Guard init.sql seed data for production only

**Status:** PENDING | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
The `init.sql` migration seeds test/development data (demo merchants, fake payments). This must NOT run in production. Add an environment guard or create a separate seed script.

**Acceptance Criteria:**
- [ ] Seed data in `init.sql` wrapped with `-- IF NOT PRODUCTION` guard (e.g. `APP_ENV != 'production'`)
- [ ] OR seed data moved to a separate `seed.sql` that is only executed in dev/test
- [ ] Migration runner skips seed data when `APP_ENV=production`
- [ ] Production deployment does not contain test merchants or dummy payments

---

### PROD-007: Add frontend Dockerfile

**Status:** PENDING | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
Create a production Dockerfile for the frontend application. This may be the same as PROD-004 or a separate Dockerfile depending on project structure.

**Acceptance Criteria:**
- [ ] Dockerfile created in the appropriate frontend directory
- [ ] Multi-stage build with production optimizations
- [ ] `.dockerignore` configured
- [ ] Image builds successfully

---

### PROD-008: Fix password reset — add email sending integration

**Status:** PENDING | **Priority:** HIGH | **Assignee:** backend-engineer

**Description:**
The forgot/reset password flow creates reset tokens but never actually sends an email. Add SMTP configuration and email sending. At minimum, log the reset link clearly to console so developers can test.

**Acceptance Criteria:**
- [ ] SMTP config added to `config.Config` (host, port, username, password, from address)
- [ ] Email service created in `internal/pkg/email/` with `SendResetPasswordEmail(to, token)` method
- [ ] If SMTP is not configured, log the reset URL to console with a clear `[RESET PASSWORD]` prefix
- [ ] If SMTP is configured, send the actual email
- [ ] Forgot-password endpoint calls email service
- [ ] Reset token expiry enforced (existing migration 016 has expiry)

---

### PROD-009: Add token refresh loop in frontend

**Status:** PENDING | **Priority:** HIGH | **Assignee:** frontend-engineer

**Description:**
JWT tokens expire. The frontend must periodically refresh the token before expiry to keep the user session alive. Implement a silent token refresh loop.

**Acceptance Criteria:**
- [ ] Token refresh endpoint `POST /v1/auth/refresh` consumed by frontend
- [ ] Refresh interval calculated based on token expiry (e.g. refresh at 80% of TTL)
- [ ] Silent refresh — no UI interruption, no redirect
- [ ] If refresh fails (e.g. refresh token expired), redirect to login
- [ ] Token stored in Zustand auth store (already persisted)
- [ ] Refresh loop starts on app load when authenticated
- [ ] Refresh loop stops on logout

---

### PROD-010: Update docs to match actual code

**Status:** PENDING | **Priority:** MEDIUM | **Assignee:** orchestrator-agent

**Description:**
Audit all documentation against the actual codebase state. Update docs to reflect changes made in this phase.

**Acceptance Criteria:**
- [ ] API docs updated for any endpoint changes from PROD-002, PROD-003
- [ ] Deployment docs updated for Dockerfile changes (PROD-004, PROD-007)
- [ ] CORS configuration documented (PROD-001)
- [ ] Email/SMTP config documented (PROD-008)
- [ ] Token refresh flow documented (PROD-009)
- [ ] All doc changes cross-referenced against code

---

## Quick Stats
| Task | Priority | Status |
|------|----------|--------|
| PROD-001 | CRITICAL | PENDING |
| PROD-002 | CRITICAL | PENDING |
| PROD-003 | CRITICAL | PENDING |
| PROD-004 | HIGH | PENDING |
| PROD-005 | HIGH | PENDING |
| PROD-006 | HIGH | PENDING |
| PROD-007 | HIGH | PENDING |
| PROD-008 | HIGH | PENDING |
| PROD-009 | HIGH | PENDING |
| PROD-010 | MEDIUM | PENDING |
