# Admin Dashboard

## Layout

```
┌────────────────────────────────────────────────────────┐
│  Header: OpenPayment Admin | Search | Notif | Profile  │
├────────┬───────────────────────────────────────────────┤
│        │                                               │
│ Admin  │           Main Content Area                    │
│ Sidebar│                                               │
│        │  ┌─────────────────────────────────────────┐  │
│ 🏠     │  │ Overview Stats                          │  │
│ Dashboard  │  Total Merchants | Active | Pending KYC    │  │
│ 🏪     │  │ Volume Today | Success Rate | Fraud %    │  │
│ Merchants│ └─────────────────────────────────────────┘  │
│ 📊     │                                               │
│ Transctn│  ┌─────────────────────────────────────────┐  │
│ ⚖️      │  │ Pending Approvals                       │  │
│ Disputes│  │ ┌────┬──────────┬───────┬──────────┬──┐│  │
│ 🎭     │  │ │Date│Business  │Country│ Docs     │  ││  │
│ Fraud   │  │ ├────┼──────────┼───────┼──────────┼──┤│  │
│ 🏦     │  │ │... │Acme Ltd  │BD     │ 3/3      │✅││  │
│ Settelemt│ │ └────┴──────────┴───────┴──────────┴──┘│  │
│ 📄     │  └─────────────────────────────────────────┘  │
│ Reports │                                               │
│ ⚙️      │  ┌─────────────────────────────────────────┐  │
│ Config  │  │ System Health                           │  │
│ 📋     │  │ 🟢 Payment Service  | 99.99% uptime     │  │
│ Audit   │  │ 🟢 Auth Service     | 99.99% uptime     │  │
│         │  │ 🟢 Ledger Service   | 99.99% uptime     │  │
└─────────┴──└─────────────────────────────────────────┘──┘
```

## Pages

### 1. Dashboard (Overview)
- **Stats:** Total merchants (active/pending/suspended), Total volume (today/month), Success rate, Fraud rate, Pending disputes
- **Pending Approvals:** List of merchants awaiting KYC review with quick approve/reject
- **System Health:** Status indicators for each service, uptime %, recent incidents
- **Recent Events:** Security alerts, failed transactions, fraud alerts
- **Quick Actions:** Approve merchants, View disputes, Run settlement

### 2. Merchants
- **Table:** All merchants with business name, email, status, verification status, volume, date
- **Filters:** Status, Verification status, Country, Date range
- **Search:** By business name, email, merchant ID
- **Bulk Actions:** Approve selected, Suspend selected, Export CSV
- **Detail View:** Full merchant profile, KYC documents, transaction volume graph, fee config, webhook logs, user management

### 3. Merchant Detail
- **Tabs:**
  - **Profile:** Business info, status, verification status
  - **KYC Documents:** List with view/download, approve/reject buttons
  - **Transactions:** All merchant transactions with filters
  - **Balance:** Current balance, settlement history, payout history
  - **Settlements:** Settlement batches with status and amounts
  - **Users:** Team members with roles
  - **API Keys:** All keys for this merchant
  - **Webhooks:** All webhook endpoints and delivery logs
  - **Activity Log:** Audit trail for this merchant

### 4. Transactions
- **Table:** All platform transactions with merchant, amount, status, payment method, date
- **Filters:** Merchant, Status, Amount range, Date range, Payment method
- **Export:** CSV export with date range

### 5. Disputes
- **Table:** All disputes with merchant, amount, reason, status, response deadline
- **Filters:** Status, Reason, Merchant
- **Detail View:** Dispute timeline, payment info, evidence submissions, resolution actions

### 6. Fraud Dashboard
- **Stats:** Total blocked today, Total flagged for review, False positive rate, Fraud rate trend
- **Rules Engine:** List of active rules with enable/disable toggles, priority, action
- **Recent Alerts:** List of flagged transactions with merchant, amount, flags, action taken
- **Rule Builder:** Create/edit rules with condition builder UI

### 7. Settlement
- **Batches Table:** Date, merchant count, total amount, fees, status, processed at
- **Manual Trigger:** Force settlement for a merchant or date range
- **Pending Batches:** Queued but not yet processed

### 8. Reports
- **Platform Summary:** Total volume, revenue, merchant growth (daily/monthly)
- **Revenue Report:** Gateway fee revenue over time
- **Merchant Report:** Per-merchant volume, fees, refund rate
- **Export Options:** CSV, PDF

### 9. Configuration
- **Currencies:** Enable/disable supported currencies
- **Countries:** Enable/disable supported countries
- **Banks:** Manage supported bank list
- **Fee Templates:** Create/edit fee structures
- **Fraud Rules:** Global rule configuration
- **System Config:** Maintenance mode toggle, max transaction amounts

### 10. Audit Logs
- **Table:** Timestamp, Actor, Action, Resource, IP, User Agent
- **Filters:** Actor, Action type, Date range, Resource type
- **Search:** Free text search on details
- **Export:** JSON/CSV export

### 11. Support Tickets (Future)
- **Table:** Ticket ID, Merchant, Subject, Status, Priority, Created
- **Detail View:** Conversation thread, merchant info, related transactions
