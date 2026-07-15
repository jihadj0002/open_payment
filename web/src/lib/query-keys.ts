export const queryKeys = {
  payments: {
    all: ['payments'] as const,
    list: (params?: Record<string, string>) => ['payments', 'list', params] as const,
    detail: (id: string) => ['payments', 'detail', id] as const,
  },
  balance: {
    all: ['balance'] as const,
  },
  apiKeys: {
    all: ['api-keys'] as const,
  },
  webhooks: {
    all: ['webhooks'] as const,
  },
  merchant: {
    profile: ['merchant', 'profile'] as const,
  },
}
