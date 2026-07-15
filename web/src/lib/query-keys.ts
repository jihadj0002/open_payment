export const queryKeys = {
  payments: {
    all: ['payments'] as const,
    list: (params?: Record<string, string>) => ['payments', 'list', params] as const,
    detail: (id: string) => ['payments', 'detail', id] as const,
  },
  balance: {
    all: ['balance'] as const,
    transactions: (params?: Record<string, string>) => ['balance', 'transactions', params] as const,
  },
  settlements: {
    list: (params?: Record<string, string>) => ['settlements', 'list', params] as const,
  },
  apiKeys: {
    all: ['api-keys'] as const,
  },
  webhooks: {
    all: ['webhooks'] as const,
  },
  customers: {
    list: () => ['customers', 'list'] as const,
    detail: (id: string) => ['customers', 'detail', id] as const,
  },
  merchant: {
    profile: ['merchant', 'profile'] as const,
  },
  reports: {
    all: ['reports'] as const,
  },
  fraud: {
    rules: ['fraud', 'rules'] as const,
  },
  onboarding: {
    status: ['onboarding', 'status'] as const,
  },
}
