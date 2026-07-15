import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import DashboardPage from './page'

const mockBalanceData = {
  available: 1000,
  pending: 200,
  reserve: 100,
  merchant_id: 'm1',
  currency: 'USD',
}

const mockPaymentsData = [
  { id: 'p1', amount: 50, currency: 'USD', status: 'succeeded', payment_method: 'card', created_at: '2026-01-01T00:00:00Z' },
  { id: 'p2', amount: 25, currency: 'USD', status: 'pending', payment_method: 'wallet', created_at: '2026-01-02T00:00:00Z' },
]

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

describe('DashboardPage', () => {
  beforeEach(() => {
    vi.spyOn(global, 'fetch').mockRestore?.()
  })

  it('renders skeleton while loading', () => {
    vi.spyOn(global, 'fetch').mockImplementation(() => new Promise(() => {}))
    render(<DashboardPage />, { wrapper: createWrapper() })
    const skeletons = document.querySelectorAll('.animate-pulse')
    expect(skeletons.length).toBeGreaterThan(0)
  })

  it('renders balance and payments on success', async () => {
    vi.spyOn(global, 'fetch').mockImplementation((url: string) => {
      if (url.includes('/balance')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: mockBalanceData }),
        })
      }
      if (url.includes('/payments')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: mockPaymentsData }),
        })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
    })
    render(<DashboardPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      expect(screen.getByText('Total Volume')).toBeInTheDocument()
      expect(screen.getByText('$1,200.00')).toBeInTheDocument()
      expect(screen.getByText('Available Balance')).toBeInTheDocument()
      expect(screen.getByText('$1,000.00')).toBeInTheDocument()
    })
  })

  it('shows error state on failure', async () => {
    vi.spyOn(global, 'fetch').mockRejectedValue(new Error('Network error'))
    render(<DashboardPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      expect(screen.getByText(/failed to load dashboard data/i)).toBeInTheDocument()
    })
  })

  it('shows create payment link when no transactions and zero balance', async () => {
    vi.spyOn(global, 'fetch').mockImplementation((url: string) => {
      if (url.includes('/balance')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: { ...mockBalanceData, available: 0, pending: 0 } }),
        })
      }
      if (url.includes('/payments')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: [] }),
        })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
    })
    render(<DashboardPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      expect(screen.getByText(/create your first payment/i)).toBeInTheDocument()
    })
  })
})
