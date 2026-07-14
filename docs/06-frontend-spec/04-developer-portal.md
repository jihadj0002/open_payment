# Developer Portal

## Layout
Similar to Merchant Dashboard but with Developer-focused navigation:

| Nav Item | Purpose |
|----------|---------|
| API Keys | Create/manage API keys |
| Documentation | Interactive API reference |
| API Explorer | Test endpoints live |
| Webhook Tester | Send test webhooks |
| Logs | API request/response logs |
| SDKs | Download SDKs and examples |

## Pages

### 1. API Keys
- Same as Merchant Dashboard API Keys section
- Clear visual distinction between **Publishable Key** (pk_*) and **Secret Key** (sk_*)
- Mode toggle: Test mode (sandbox) vs Live mode
- Permission scopes per key

### 2. Documentation
- **Sidebar navigation** by API resource:
  - Payments (Create, Retrieve, List, Capture, Refund, Void)
  - Customers
  - Refunds
  - Webhooks
  - Balance
  - Tokens
  - Errors
- Each endpoint shows:
  - HTTP method + path
  - Description
  - Request parameters table
  - Request example (cURL, Python, Go, JavaScript, PHP)
  - Response example
  - Error scenarios
- Search bar for documentation
- OpenAPI spec download link

### 3. API Explorer
- Interactive API testing tool (like Postman in-browser)
- Endpoint selector dropdown
- Request builder:
  - Method + URL pre-filled
  - Headers (Authorization auto-filled with user's selected key)
  - Body editor with JSON validation + syntax highlighting
- Response viewer:
  - Status code + body with syntax highlighting
  - Response time
  - Request ID
- History of recent requests
- **Sandbox mode only** — cannot be used with live keys

### 4. Webhook Tester
- **Send Test Event:**
  - Event type dropdown (payment.success, payment.failed, refund.completed)
  - Target URL input
  - Custom payload editor (optional)
  - Signature secret display
- **Delivery Log:** History of sent test webhooks with response
- **Endpoint Validator:** Test a configured webhook endpoint returns proper response

### 5. Logs
- **API Request Logs:** All API requests made by the merchant
  - Columns: Timestamp, Method, Path, Status, Response time, IP
  - Click to expand: Full request/response headers and body
  - Filters: Method, Status code, Date range
  - Search: By request ID, path
- **Webhook Delivery Logs:** All webhook delivery attempts
  - Columns: Timestamp, Event type, URL, Status, Response code, Attempt #
  - Click to expand: Payload, response body
  - Filters: Status, Event type, Date range

### 6. SDKs & Libraries
- **Available SDKs:**
  - Go: `go get github.com/openpayment/go-sdk`
  - Python: `pip install openpayment`
  - JavaScript: `npm install openpayment`
  - PHP: `composer require openpayment/php-sdk`
  - Ruby: `gem install openpayment`
- Each SDK shows:
  - Installation command
  - Quick start example
  - Link to GitHub repo
  - API reference
- **Sample Projects:** Links to example e-commerce integrations

### 7. Quick Start Guide
- Step-by-step: "Accept your first payment in 5 minutes"
  1. Get API keys from sandbox
  2. Install SDK
  3. Create a payment
  4. Handle webhook
  5. Go live
- Code snippets for each step in multiple languages

## Developer Experience Principles
- **Zero setup time:** Sandbox keys generated on signup
- **Consistent SDKs:** Same method signatures across languages
- **Clear error messages:** Every error includes a link to relevant docs
- **Rate limit headers:** Exposed in response for programmatic handling
- **Idempotency explained:** Examples show correct usage
