# Frontend Implementation

## Tech Stack
- **Framework:** Next.js 14+ (App Router)
- **Language:** TypeScript (strict mode)
- **Styling:** Tailwind CSS 3.4
- **State:** TanStack React Query v5 + Zustand
- **Forms:** React Hook Form + Zod
- **Charts:** Apache ECharts
- **HTTP:** fetch (native) with custom wrapper
- **Testing:** Vitest + Playwright

## Project Structure

```
frontend/
├── public/
│   ├── icons/
│   └── images/
│
├── src/
│   ├── app/
│   │   ├── (auth)/
│   │   │   ├── login/
│   │   │   │   └── page.tsx
│   │   │   ├── register/
│   │   │   │   └── page.tsx
│   │   │   └── layout.tsx
│   │   ├── (dashboard)/
│   │   │   ├── merchant/
│   │   │   │   ├── layout.tsx
│   │   │   │   ├── page.tsx
│   │   │   │   ├── payments/
│   │   │   │   ├── refunds/
│   │   │   │   └── ...
│   │   │   └── admin/
│   │   │       └── ...
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── not-found.tsx
│   │
│   ├── components/
│   │   ├── ui/                    # Primitive components
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── DataTable.tsx
│   │   │   ├── Badge.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── Toast.tsx
│   │   │   ├── Input.tsx
│   │   │   └── Select.tsx
│   │   ├── charts/                # Chart components
│   │   │   ├── RevenueChart.tsx
│   │   │   └── PaymentStatusPie.tsx
│   │   ├── layout/                # Layout components
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Header.tsx
│   │   │   ├── Breadcrumbs.tsx
│   │   │   └── PageContainer.tsx
│   │   └── payments/              # Feature-specific
│   │       ├── PaymentTable.tsx
│   │       ├── PaymentDetail.tsx
│   │       ├── PaymentTimeline.tsx
│   │       ├── RefundModal.tsx
│   │       └── CreatePaymentForm.tsx
│   │
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts         # Base API client
│   │   │   ├── payments.ts       # Payment API functions
│   │   │   ├── merchants.ts
│   │   │   ├── customers.ts
│   │   │   └── webhooks.ts
│   │   ├── hooks/
│   │   │   ├── usePayments.ts    # React Query hooks
│   │   │   ├── useCustomers.ts
│   │   │   ├── useBalance.ts
│   │   │   └── useWebhooks.ts
│   │   ├── stores/
│   │   │   ├── auth-store.ts     # Zustand auth store
│   │   │   └── ui-store.ts       # UI state (sidebar, theme)
│   │   ├── utils/
│   │   │   ├── format.ts         # Currency, date formatters
│   │   │   ├── cn.ts             # classnames helper
│   │   │   └── validators.ts     # Zod schemas
│   │   └── types/
│   │       ├── payment.ts
│   │       ├── merchant.ts
│   │       └── api.ts
│   │
│   ├── providers/
│   │   ├── QueryProvider.tsx     # React Query provider
│   │   └── AuthProvider.tsx      # Auth context
│   │
│   └── styles/
│       └── globals.css           # Tailwind imports, base styles
│
├── middleware.ts                  # Auth middleware
├── next.config.js
├── tailwind.config.ts
├── tsconfig.json
├── vitest.config.ts
├── playwright.config.ts
└── package.json
```

## Key Implementation Patterns

### API Client
```typescript
// src/lib/api/client.ts
class ApiClient {
  private baseUrl: string;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL!;
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    // Token from auth store (dashboard) or API key (developer portal)
    const token = useAuthStore.getState().accessToken;
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return headers;
  }

  private async request<T>(
    method: string,
    path: string,
    data?: unknown,
    options?: { idempotencyKey?: string }
  ): Promise<T> {
    const headers = this.getHeaders();

    if (options?.idempotencyKey) {
      headers['Idempotency-Key'] = options.idempotencyKey;
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers,
      body: data ? JSON.stringify(data) : undefined,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new ApiError(response.status, error);
    }

    return response.json();
  }

  get<T>(path: string, params?: Record<string, string>) {
    const query = params ? '?' + new URLSearchParams(params).toString() : '';
    return this.request<T>('GET', `${path}${query}`);
  }

  post<T>(path: string, data?: unknown, options?: { idempotencyKey?: string }) {
    return this.request<T>('POST', path, data, options);
  }

  patch<T>(path: string, data?: unknown) {
    return this.request<T>('PATCH', path, data);
  }

  delete<T>(path: string) {
    return this.request<T>('DELETE', path);
  }
}

export const api = new ApiClient();
```

### Data Fetching Hook
```typescript
// src/lib/hooks/usePayments.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api/client';
import { PaymentIntent, CreatePaymentRequest, PaymentFilters } from '@/lib/types/payment';

export const paymentKeys = {
  all: ['payments'] as const,
  list: (filters: PaymentFilters) => ['payments', 'list', filters] as const,
  detail: (id: string) => ['payments', 'detail', id] as const,
};

export function usePayments(filters: PaymentFilters) {
  return useQuery({
    queryKey: paymentKeys.list(filters),
    queryFn: () => api.get<PaginatedResponse<PaymentIntent>>('/v1/payments', filters as any),
    placeholderData: keepPreviousData,
  });
}

export function usePayment(id: string) {
  return useQuery({
    queryKey: paymentKeys.detail(id),
    queryFn: () => api.get<PaymentIntent>(`/v1/payments/${id}`),
    enabled: !!id,
  });
}

export function useCreatePayment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePaymentRequest) =>
      api.post('/v1/payments', data, {
        idempotencyKey: crypto.randomUUID(),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: paymentKeys.all });
    },
  });
}

export function useCapturePayment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      api.post(`/v1/payments/${id}/capture`),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: paymentKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: paymentKeys.all });
    },
  });
}
```

### Server Component Example
```typescript
// src/app/(dashboard)/merchant/payments/page.tsx
export default async function PaymentsPage() {
  return (
    <PageContainer
      title="Payments"
      description="View and manage your payments"
      action={<CreatePaymentButton />}
    >
      <PaymentFilters />
      <Suspense fallback={<PaymentTableSkeleton />}>
        <PaymentTableClient />
      </Suspense>
    </PageContainer>
  );
}
```

### Client Component Example
```typescript
'use client';

export function RefundModal({ payment, onClose }: Props) {
  const refundMutation = useRefundPayment();
  const form = useForm({
    resolver: zodResolver(refundSchema),
    defaultValues: {
      amount: payment.amount,
      reason: 'customer_request',
    },
  });

  const onSubmit = (data: RefundForm) => {
    refundMutation.mutate(
      { paymentId: payment.id, ...data },
      {
        onSuccess: () => {
          toast.success('Refund processed');
          onClose();
        },
        onError: (err) => {
          toast.error(err.message);
        },
      }
    );
  };

  return (
    <Modal open onClose={onClose}>
      <Modal.Header>Issue Refund</Modal.Header>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Modal.Body>
          <Input
            label="Amount (BDT)"
            type="number"
            {...form.register('amount', { valueAsNumber: true })}
            error={form.formState.errors.amount?.message}
          />
          <Select
            label="Reason"
            {...form.register('reason')}
            options={[
              { value: 'customer_request', label: 'Customer Request' },
              { value: 'duplicate', label: 'Duplicate' },
            ]}
          />
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="danger" loading={refundMutation.isPending}>
            Process Refund
          </Button>
        </Modal.Footer>
      </form>
    </Modal>
  );
}
```

## Testing

### Unit Tests (Vitest)
```typescript
// __tests__/components/RefundModal.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { RefundModal } from '@/components/payments/RefundModal';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

describe('RefundModal', () => {
  it('shows refund amount prefilled', () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <RefundModal
          payment={{ id: 'pi_123', amount: 1000, currency: 'BDT' }}
          onClose={vi.fn()}
        />
      </QueryClientProvider>
    );

    expect(screen.getByDisplayValue('1000')).toBeInTheDocument();
  });
});
```

### E2E Tests (Playwright)
```typescript
// e2e/payment-flow.spec.ts
test('merchant can create a payment', async ({ page }) => {
  // Login
  await page.goto('/login');
  await page.fill('[name="email"]', 'test@merchant.com');
  await page.fill('[name="password"]', 'password123');
  await page.click('button[type="submit"]');

  // Navigate to payments
  await page.click('text=Payments');
  await page.click('text=Create Payment');

  // Fill form
  await page.fill('[name="amount"]', '500');
  await page.click('text=Create');

  // Verify payment created
  await expect(page.locator('text=BDT 500')).toBeVisible();
  await expect(page.locator('text=succeeded')).toBeVisible();
});
```
