# Code Quality Standards

## Linters and Formatters

### Go (golangci-lint)
```yaml
# .golangci.yml
linters:
  enable:
    - errcheck           # Check returned errors
    - gosimple           # Simplify code
    - govet              # Suspicious constructs
    - ineffassign        # Ineffective assignments
    - staticcheck        # Static analysis
    - typecheck          # Type checking
    - unused             # Unused code
    - gosec              # Security issues
    - revive             # Style checker
    - gofmt              # Formatting
    - gocyclo            # Cyclomatic complexity
    - misspell           # Spelling mistakes
    - prealloc           # Slice preallocation hints
    - bodyclose          # HTTP response body check
    - noctx              # Context propagation check
    - sqlclosecheck      # SQL rows close check

linters-settings:
  gocyclo:
    min-complexity: 15
  errcheck:
    check-type-assertions: true
  misspell:
    locale: US

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
        - gosec
```

### TypeScript (ESLint + Prettier)
```json
// .eslintrc.json
{
  "extends": [
    "next/core-web-vitals",
    "plugin:@typescript-eslint/recommended",
    "prettier"
  ],
  "rules": {
    "@typescript-eslint/no-unused-vars": ["error", { "argsIgnorePattern": "^_" }],
    "@typescript-eslint/explicit-function-return-type": "warn",
    "no-console": ["warn", { "allow": ["warn", "error"] }],
    "react-hooks/exhaustive-deps": "error"
  }
}
```

### Pre-commit Hooks
```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json

  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.0
    hooks:
      - id: go-fmt
      - id: go-vet
      - id: go-lint
      - id: go-imports

  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.18.0
    hooks:
      - id: gitleaks
```

## Code Review Checklist

### Architecture & Design
- [ ] Does the code follow the established patterns?
- [ ] Are cross-service dependencies minimized?
- [ ] Is the change backward-compatible (API, DB schema)?
- [ ] Are there any idempotency concerns addressed?
- [ ] Is error handling consistent with our patterns?

### Security
- [ ] Are all inputs validated?
- [ ] Is the merchant_id scope enforced on all queries?
- [ ] Are secrets/config never hardcoded?
- [ ] Are there no SQL injection possibilities?
- [ ] Are sensitive fields not logged?
- [ ] Is rate limiting applied?

### Performance
- [ ] Are N+1 queries avoided?
- [ ] Are appropriate indexes used?
- [ ] Is pagination used for list endpoints?
- [ ] Are database queries using proper filters?
- [ ] Caching considered for read-heavy paths?

### Testing
- [ ] Unit tests for new logic?
- [ ] Integration tests for service interaction?
- [ ] Edge cases covered (empty, invalid, boundary)?
- [ ] Error paths tested?
- [ ] Test names describe what's being tested?

### Code Style
- [ ] No commented-out code?
- [ ] No TODO without ticket reference?
- [ ] Functions are small (<30 lines)?
- [ ] Meaningful variable names?
- [ ] Exported functions have Go doc comments?

## Branch Naming Convention
```
feature/OP-123-add-refund-api
bugfix/OP-456-fix-payment-state-transition
hotfix/OP-789-critical-security-fix
chore/OP-012-update-dependencies
```

## Commit Message Convention
```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `perf`, `security`
Scopes: `payment`, `ledger`, `auth`, `merchant`, `webhook`, `fraud`, `infra`, `ui`

```
feat(payment): add partial refund support

Allows merchants to refund a portion of a captured payment.
Implements OP-456.

Closes OP-456
```

## PR Template
```markdown
## Description
[Brief description of the change]

## Related Issue
Closes #[issue_number]

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation
- [ ] Refactor

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] E2E tests pass
- [ ] Manual testing completed

## Security Considerations
- [ ] Input validation added
- [ ] Authorization scoped
- [ ] No sensitive data logged

## Checklist
- [ ] Code follows project style
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No new TODOs without ticket reference
```
