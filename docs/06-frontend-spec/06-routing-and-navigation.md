# Routing and Navigation

## Next.js App Router Structure

```
app/
├── (auth)/
│   ├── login/
│   │   └── page.tsx
│   ├── register/
│   │   └── page.tsx
│   ├── forgot-password/
│   │   └── page.tsx
│   ├── reset-password/
│   │   └── page.tsx
│   └── layout.tsx                    (centered card layout)
│
├── (dashboard)/
│   ├── merchant/
│   │   ├── page.tsx                  (Overview)
│   │   ├── layout.tsx                (sidebar + header)
│   │   ├── payments/
│   │   │   ├── page.tsx              (list)
│   │   │   └── [id]/
│   │   │       └── page.tsx          (detail)
│   │   ├── refunds/
│   │   │   └── page.tsx
│   │   ├── customers/
│   │   │   ├── page.tsx              (list)
│   │   │   └── [id]/
│   │   │       └── page.tsx          (detail)
│   │   ├── balance/
│   │   │   └── page.tsx
│   │   ├── api-keys/
│   │   │   └── page.tsx
│   │   ├── webhooks/
│   │   │   ├── page.tsx              (endpoints list)
│   │   │   └── logs/
│   │   │       └── page.tsx          (delivery logs)
│   │   ├── reports/
│   │   │   ├── page.tsx              (report list)
│   │   │   └── [type]/
│   │   │       └── page.tsx          (specific report)
│   │   └── settings/
│   │       ├── page.tsx              (profile)
│   │       ├── users/
│   │       │   └── page.tsx
│   │       └── security/
│   │           └── page.tsx
│   │
│   ├── admin/
│   │   ├── page.tsx                  (Admin Dashboard)
│   │   ├── layout.tsx                (admin sidebar)
│   │   ├── merchants/
│   │   │   ├── page.tsx              (list)
│   │   │   └── [id]/
│   │   │       └── page.tsx          (detail)
│   │   ├── transactions/
│   │   │   └── page.tsx
│   │   ├── disputes/
│   │   │   ├── page.tsx
│   │   │   └── [id]/
│   │   │       └── page.tsx
│   │   ├── fraud/
│   │   │   └── page.tsx
│   │   ├── settlement/
│   │   │   └── page.tsx
│   │   ├── reports/
│   │   │   └── page.tsx
│   │   ├── config/
│   │   │   └── page.tsx
│   │   └── audit-logs/
│   │       └── page.tsx
│   │
│   └── developer/
│       ├── page.tsx                  (API Keys — quick access)
│       ├── docs/
│       │   └── page.tsx              (API reference)
│       ├── explorer/
│       │   └── page.tsx              (API explorer)
│       ├── webhook-tester/
│       │   └── page.tsx
│       └── logs/
│           └── page.tsx
│
├── checkout/
│   └── [payment_intent_client_secret]/
│       └── page.tsx                  (Hosted checkout page)
│
├── api/                              (API routes — BFF layer)
│   ├── auth/
│   │   ├── login/route.ts
│   │   ├── logout/route.ts
│   │   └── refresh/route.ts
│   └── proxy/                        (proxy to backend APIs)
│       └── [...path]/route.ts
│
├── layout.tsx                        (root layout: fonts, providers)
├── page.tsx                          (landing/redirect)
└── not-found.tsx
```

## Route Guards

### Auth Guard
```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const publicRoutes = ['/login', '/register', '/forgot-password', '/api/auth/login'];

export function middleware(request: NextRequest) {
  const token = request.cookies.get('access_token')?.value;
  const isPublic = publicRoutes.some((route) =>
    request.nextUrl.pathname.startsWith(route)
  );

  if (!token && !isPublic) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
};
```

### Role Guard (per page)
```typescript
// In each layout or page:
export default function MerchantPaymentsPage() {
  const { user } = useAuthStore();

  if (!user?.permissions.includes('payments:read')) {
    return <Unauthorized />;
  }

  return <PaymentsPage />;
}
```

## Navigation Components

### Sidebar
```typescript
interface SidebarItem {
  href: string;
  label: string;
  icon: React.ComponentType;
  permissions?: string[];         // Required permissions to see
  children?: SidebarItem[];       // Collapsible sub-items
}

const merchantNav: SidebarItem[] = [
  { href: '/merchant', label: 'Overview', icon: LayoutDashboard, permissions: ['dashboard:read'] },
  { href: '/merchant/payments', label: 'Payments', icon: CreditCard, permissions: ['payments:read'] },
  { href: '/merchant/refunds', label: 'Refunds', icon: RotateCcw, permissions: ['refunds:read'] },
  { href: '/merchant/customers', label: 'Customers', icon: Users, permissions: ['customers:read'] },
  { href: '/merchant/balance', label: 'Balance', icon: Wallet, permissions: ['balance:read'] },
  { href: '/merchant/api-keys', label: 'API Keys', icon: Key, permissions: ['api_keys:read'] },
  { href: '/merchant/webhooks', label: 'Webhooks', icon: Webhook, permissions: ['webhooks:read'] },
  { href: '/merchant/reports', label: 'Reports', icon: BarChart3, permissions: ['reports:read'] },
  { href: '/merchant/settings', label: 'Settings', icon: Settings, permissions: ['settings:read'] },
];
```

### Breadcrumbs
```typescript
// Auto-generated from path segments
<Breadcrumbs>
  <BreadcrumbItem href="/merchant">Dashboard</BreadcrumbItem>
  <BreadcrumbItem href="/merchant/payments">Payments</BreadcrumbItem>
  <BreadcrumbItem>pi_abc123</BreadcrumbItem>
</Breadcrumbs>
```

## Hosted Checkout Page
- Accessed via `/{client_secret}` (short URL)
- Server-side renders payment form with merchant branding
- Supports: Card input (Stripe Elements-like), Wallet redirect, Bank transfer instructions
- States:
  - **Loading:** Skeleton payment form
  - **Ready:** Payment form with card input
  - **Processing:** Spinner on submit button
  - **Success:** Green checkmark, confirmation message
  - **Failed:** Error message with retry option
  - **Expired:** Session expired message
