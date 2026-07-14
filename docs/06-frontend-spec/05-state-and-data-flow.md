# State Management and Data Flow

## Tech Stack
- **Server State:** TanStack React Query (v5)
- **Client State:** Zustand
- **Form State:** React Hook Form + Zod
- **URL State:** Next.js searchParams + useRouter

## State Categories

| Category | Tool | Example Data |
|----------|------|-------------|
| Server state (API data) | React Query | Payments list, customer data, balance |
| Client state (UI) | Zustand | Sidebar open/closed, selected filters |
| Form state | React Hook Form | Create payment form, refund form |
| URL state | Next.js router | Page, sort, filter params |
| Auth state | Zustand + React Query | JWT token, user profile, permissions |

## React Query Configuration

### Query Client Setup
```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,        // 30 seconds before refetch
      gcTime: 5 * 60_000,       // 5 minutes in cache
      retry: 2,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,                 // Don't retry mutations
    },
  },
});
```

### Query Keys Pattern
```typescript
export const queryKeys = {
  payments: {
    all: ['payments'] as const,
    list: (filters: PaymentFilters) => ['payments', 'list', filters] as const,
    detail: (id: string) => ['payments', 'detail', id] as const,
  },
  customers: {
    all: ['customers'] as const,
    list: (filters: CustomerFilters) => ['customers', 'list', filters] as const,
    detail: (id: string) => ['customers', 'detail', id] as const,
  },
  balance: {
    all: ['balance'] as const,
    transactions: (filters: BalanceTxFilters) => ['balance', 'transactions', filters] as const,
  },
  webhooks: {
    all: ['webhooks'] as const,
    endpoints: ['webhooks', 'endpoints'] as const,
    logs: (filters: WebhookLogFilters) => ['webhooks', 'logs', filters] as const,
  },
  merchant: {
    profile: ['merchant', 'profile'] as const,
    users: ['merchant', 'users'] as const,
    apiKeys: ['merchant', 'api-keys'] as const,
  },
};
```

### Data Fetching Hook Example
```typescript
export function usePayments(filters: PaymentFilters) {
  return useQuery({
    queryKey: queryKeys.payments.list(filters),
    queryFn: () => api.getPayments(filters),
    placeholderData: keepPreviousData,  // Keep old data while fetching next page
  });
}

export function usePaymentDetail(id: string) {
  return useQuery({
    queryKey: queryKeys.payments.detail(id),
    queryFn: () => api.getPayment(id),
    enabled: !!id,
  });
}
```

### Mutation + Cache Update Pattern
```typescript
export function useRefundPayment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: RefundRequest) => api.createRefund(data),
    onSuccess: (refund) => {
      // Invalidate affected queries
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.all });
      queryClient.invalidateQueries({ queryKey: queryKeys.balance.all });
      toast.success('Refund processed successfully');
    },
    onError: (error: ApiError) => {
      toast.error(error.message);
    },
  });
}
```

## Zustand Stores

### Auth Store
```typescript
interface AuthState {
  user: User | null;
  merchant: Merchant | null;
  permissions: string[];
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      merchant: null,
      permissions: [],
      isAuthenticated: false,
      login: async (credentials) => {
        const response = await api.login(credentials);
        set({
          user: response.user,
          merchant: response.merchant,
          permissions: response.permissions,
          isAuthenticated: true,
        });
      },
      logout: () => {
        set({ user: null, merchant: null, permissions: [], isAuthenticated: false });
        queryClient.clear();
      },
      refreshToken: async () => {
        const response = await api.refreshToken();
        set({ user: response.user, permissions: response.permissions });
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        merchant: state.merchant,
        permissions: state.permissions,
      }),
    }
  )
);
```

### UI Store
```typescript
interface UIState {
  sidebarOpen: boolean;
  theme: 'light' | 'dark' | 'system';
  toggleSidebar: () => void;
  setTheme: (theme: 'light' | 'dark' | 'system') => void;
}

export const useUIStore = create<UIState>()((set) => ({
  sidebarOpen: true,
  theme: 'system',
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  setTheme: (theme) => set({ theme }),
}));
```

## API Layer
```typescript
// lib/api/client.ts
class ApiClient {
  private baseUrl: string;
  private apiKey: string | null = null;
  private jwtToken: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  setApiKey(key: string) { this.apiKey = key; }
  setJWT(token: string) { this.jwtToken = token; }

  private async request<T>(method: string, path: string, data?: unknown): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (this.apiKey) headers['Authorization'] = `Bearer ${this.apiKey}`;
    if (this.jwtToken) headers['Authorization'] = `Bearer ${this.jwtToken}`;

    // Add idempotency key for mutations
    if (['POST', 'PATCH', 'DELETE'].includes(method)) {
      headers['Idempotency-Key'] = crypto.randomUUID();
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers,
      body: data ? JSON.stringify(data) : undefined,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new ApiError(response.status, error.error || {});
    }

    return response.json();
  }

  get<T>(path: string, params?: Record<string, string>) {
    const query = params ? '?' + new URLSearchParams(params).toString() : '';
    return this.request<T>('GET', `${path}${query}`);
  }

  post<T>(path: string, data?: unknown) {
    return this.request<T>('POST', path, data);
  }

  patch<T>(path: string, data?: unknown) {
    return this.request<T>('PATCH', path, data);
  }

  delete<T>(path: string) {
    return this.request<T>('DELETE', path);
  }
}

export const api = new ApiClient(process.env.NEXT_PUBLIC_API_URL!);
```

## Optimistic Updates
```typescript
export function useCapturePayment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (paymentId: string) => api.post(`/payments/${paymentId}/capture`),
    onMutate: async (paymentId) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: queryKeys.payments.detail(paymentId) });
      // Snapshot previous value
      const previous = queryClient.getQueryData(queryKeys.payments.detail(paymentId));
      // Optimistically update
      queryClient.setQueryData(queryKeys.payments.detail(paymentId), (old: PaymentIntent) => ({
        ...old,
        status: 'succeeded',
        amount_received: old.amount,
      }));
      return { previous };
    },
    onError: (err, paymentId, context) => {
      // Rollback on error
      queryClient.setQueryData(queryKeys.payments.detail(paymentId), context?.previous);
    },
    onSettled: (data, error, paymentId) => {
      // Refetch to ensure correct state
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.detail(paymentId) });
    },
  });
}
```

## Error Handling
```typescript
export class ApiError extends Error {
  constructor(
    public status: number,
    public error: { type: string; code: string; message: string; param?: string }
  ) {
    super(error.message);
    this.name = 'ApiError';
  }

  get isValidationError() { return this.status === 400; }
  get isAuthError() { return this.status === 401; }
  get isRateLimit() { return this.status === 429; }
  get isServerError() { return this.status >= 500; }
}
```
