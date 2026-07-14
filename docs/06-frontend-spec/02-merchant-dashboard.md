# Merchant Dashboard

## Layout Structure

```
┌────────────────────────────────────────────────────────┐
│  Header: Logo | Search | Notifications | Profile ▼    │
├────────┬───────────────────────────────────────────────┤
│        │                                               │
│ Sidebar│           Main Content Area                    │
│        │                                               │
│  📊    │  ┌─────────────────────────────────────────┐  │
│  Overview  │  Stat Cards Row                         │  │
│  💳      │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐    │  │
│  Payments │  │Revenue│ │Pending│ │Refund│ │Chrgbk│    │  │
│  🔄      │  └──────┘ └──────┘ └──────┘ └──────┘    │  │
│  Refunds  │                                         │  │
│  👤      │  ┌─────────────────────────────────────┐  │  │
│  Customers│  │ Recent Transactions                  │  │  │
│  🔑      │  │ ┌────┬────────┬──────┬──────────┬──┐│  │  │
│  API Keys │  │ │ID  │Amount  │Status│ Date     │  ││  │  │
│  🔔      │  │ ├────┼────────┼──────┼──────────┼──┤│  │  │
│  Webhooks │  │ │... │BDT 500 │✅    │15 Jan    │  ││  │  │
│  📄      │  │ └────┴────────┴──────┴──────────┴──┘│  │  │
│  Reports  │  └─────────────────────────────────────┘  │  │
│  ⚙️       │                                         │  │
│  Settings │  ┌─────────────────────────────────────┐  │  │
│          │  │ Chart: Revenue (7 days)              │  │  │
│          │  └─────────────────────────────────────┘  │  │
└──────────┴───────────────────────────────────────────┘
```

## Pages

### 1. Overview (Dashboard Home)
- **Stats Cards:** Today's revenue, Monthly revenue, Pending amount, Refunds this month, Chargebacks
- **Chart:** Revenue trend (7 days / 30 days / 12 months)
- **Recent Transactions:** Last 10 transactions with status
- **Quick Actions:** Create payment, View balance, API keys

### 2. Payments
- **Table:** All payments with columns: ID, Amount, Currency, Status, Customer, Date, Actions
- **Filters:** Date range, Status, Payment method, Amount range
- **Search:** By payment ID, customer email, order ID
- **Detail View:** Click to expand/modal with full payment details, transaction timeline, refund option
- **Bulk Actions:** Export selected as CSV

### 3. Payment Detail
- **Header:** Payment ID, Status badge, Amount
- **Timeline:** Visual state machine showing each transition with timestamps
- **Details:** Amount breakdown, payment method, customer info, metadata
- **Transactions:** List of authorization, capture, refund transactions
- **Actions:** Capture (if authorized), Refund (if captured), Void (if authorized)
- **Webhook Logs:** Delivery status for this payment's events

### 4. Refunds
- **Table:** All refunds with payment link, amount, status, reason, date
- **Filters:** Date range, Status, Reason
- **Initiate Refund:** Modal with amount (defaults to payment remaining), reason selector

### 5. Customers
- **Table:** All customers with name, email, phone, total spent, payment count, last payment
- **Detail View:** Customer profile, saved payment methods, transaction history
- **Create Customer:** Form with name, email, phone

### 6. Balance
- **Balance Breakdown:** Available, Pending, Reserve amounts by currency
- **Transaction History:** All balance-affecting events (captures, refunds, fees, settlements)
- **Payout History:** Past payouts with amounts, dates, status

### 7. API Keys
- **Keys Table:** Name, prefix, mode (live/test), created date, last used, status
- **Create Key:** Form with name, permissions checkboxes, IP whitelist (optional)
- **Reveal Secret:** Modal warning, shows secret once
- **Revoke Key:** Confirmation modal
- **Security Notice:** Reminder to not share secret keys

### 8. Webhooks
- **Endpoints Table:** URL, events subscribed, status, last delivery
- **Add Endpoint:** Form with URL, event checkboxes
- **Webhook Logs:** Delivery attempt history with status, response code, payload preview
- **Test Webhook:** Button to send a test event
- **Replay:** Retry failed webhooks

### 9. Reports
- **Revenue Report:** Daily/weekly/monthly/date range with chart + table
- **Transaction Report:** All transactions filtered by date range, downloadable as CSV
- **Settlement Report:** Settlement batches with amounts, fees, net
- **Export:** CSV, PDF generation

### 10. Settings
- **Profile:** Business name, email, phone, address
- **Users:** Team members list, invite user, role management
- **Bank Accounts:** Add/update payout bank details
- **Notifications:** Email notification preferences
- **Security:** Password change, MFA setup, active sessions

## Navigation & Permissions

| Nav Item | Owner | Admin | Developer | Analyst | Member |
|----------|-------|-------|-----------|---------|--------|
| Overview | ✅ | ✅ | ✅ | ✅ | ✅ |
| Payments | ✅ | ✅ | ✅ | ✅ | ✅ |
| Refunds | ✅ | ✅ | ✅ | ✅ | - |
| Customers | ✅ | ✅ | ✅ | ✅ | ✅ |
| Balance | ✅ | ✅ | - | ✅ | - |
| API Keys | ✅ | ✅ | ✅ | - | - |
| Webhooks | ✅ | ✅ | ✅ | - | - |
| Reports | ✅ | ✅ | - | ✅ | - |
| Settings | ✅ | ✅ | - | - | - |

## State Handling Per Page

| Page | Loading State | Empty State | Error State |
|------|---------------|-------------|-------------|
| Overview | Skeleton cards + shimmer | First-time setup prompt | "Failed to load data" + retry |
| Payments | Skeleton table | "No payments yet. Try creating one." | Error toast + retry button |
| Payment Detail | Spinner + skeleton | (always loads) | "Payment not found" |
| Refunds | Skeleton table | "No refunds yet" | Error toast |
| Customers | Skeleton table | "Add your first customer" CTA | Error toast |
| Balance | Skeleton stat cards | "BDT 0.00 available" | Error toast |
| API Keys | Skeleton table | "Generate your first API key" CTA | Error toast |
| Webhooks | Skeleton table | "Configure your first webhook" CTA | Error toast |
| Settings | Spinner per section | Default form values | Error per failed section |
