# Mobile App Specification

> **Status:** ❌ **TODO: Not implemented** — No mobile app exists yet. This doc is a spec/roadmap.

## Tech Stack
- **Framework:** Flutter 3.x
- **State Management:** Riverpod
- **Networking:** Dio + Retrofit
- **Local Storage:** Hive + flutter_secure_storage
- **Charts:** fl_chart

## Scope (MVP)
Merchant mobile app focused on:
- Real-time payment notifications
- Dashboard overview
- Transaction history view (read-only)
- Refund initiation (with confirmation)
- Push notifications for payment events

## Screens

### 1. Authentication
- **Login screen:** Email + password, "Remember me" toggle, biometric unlock option
- **PIN/Biometric:** Optional quick access PIN or fingerprint/face unlock
- **Session management:** Auto-logout after inactivity (configurable)

### 2. Dashboard
- **Stats cards:** Today's revenue, Pending amount, Successful payments count
- **Recent transactions:** Last 5 payments with status indicators
- **Quick actions:** View all transactions, View balance
- **Refresh:** Pull-to-refresh, auto-refresh every 30s when app is in foreground

### 3. Transactions
- **List view:** Infinite scroll, paginated API calls
- **Filters:** Date range, Status, Payment method
- **Detail view:** Amount, status, payment method, customer, date, transaction ID
- **Share receipt:** Generate and share payment confirmation via OS share sheet

### 4. Notifications
- **Push notifications:** payment.success, payment.failed, refund.completed, chargeback.created
- **Notification list:** In-app notification center with read/unread status
- **Tap to navigate:** Tap notification → open relevant transaction detail

### 5. Profile
- **Merchant info:** Business name, email, phone
- **App settings:** Notification preferences, PIN/biometric config, theme (light/dark)
- **Support:** Contact support, FAQ links

## Navigation

```
Bottom Navigation:
┌──────────┬────────────┬────────────┬──────────┐
│  📊      │  💳        │  🔔        │  👤      │
│ Dashboard│ Transactions│ Notifications│ Profile│
└──────────┴────────────┴────────────┴──────────┘

Dashboard → Transaction list → Transaction detail
Notifications → Transaction detail
Transactions → Filters → Transaction detail
```

## Push Notifications

### Setup
- **Firebase Cloud Messaging (FCM)** for Android
- **Apple Push Notification Service (APNS)** for iOS (via FCM proxy)
- Device token registered on login, deregistered on logout

### Notification Payload
```json
{
  "notification": {
    "title": "Payment Received",
    "body": "BDT 500.00 from customer@example.com"
  },
  "data": {
    "type": "payment.success",
    "payment_id": "pi_abc123",
    "amount": "500",
    "currency": "BDT",
    "merchant_id": "m_abc123"
  }
}
```

### Notification Types
| Event | Title | Body |
|-------|-------|------|
| payment.success | Payment Received | "BDT {amount} from {customer}" |
| payment.failed | Payment Failed | "BDT {amount} — {reason}" |
| refund.completed | Refund Processed | "BDT {amount} refunded to {customer}" |
| chargeback.created | Chargeback Alert | "BDT {amount} — dispute from {customer}" |
| settlement.completed | Settlement Complete | "BDT {amount} settled to your account" |

## State Handling

| Screen | Loading | Empty | Error | Offline |
|--------|---------|-------|-------|---------|
| Dashboard | Skeleton cards | "No data yet" | Error with retry | Show cached data + "Last updated: X min ago" |
| Transactions | Shimmer list | "No transactions" | Error with retry | Show cached data |
| Notifications | Skeleton list | "No notifications" | Error with retry | Show cached data |
| Profile | Spinner | — | Toast error | Editable fields disabled |

## Caching Strategy
- **Dashboard stats:** Cache for 5 minutes, show cached + stale indicator
- **Transaction list:** Cache last page, infinite scroll pre-fetches next
- **Transaction detail:** Cache for 2 minutes
- **Auth token:** Secure storage, refresh silently when expired

## Release Plan
| Phase | Version | Features |
|-------|---------|----------|
| MVP | v1.0 | Auth, Dashboard, Transactions, Notifications, Profile |
| v1.1 | v1.1 | Refund initiation from mobile, Push notification preferences |
| v1.2 | v1.2 | Biometric auth, Dark mode, Multiple currencies support |
| v2.0 | v2.0 | QR code payments, Customer messaging, Analytics charts |

## Platforms
- **Minimum OS:** Android 8.0 (API 26), iOS 15.0
- **Target:** Android 14, iOS 17
- **Distribution:** Google Play Store, Apple App Store
