import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import PaymentsPage from './page'

const mockPaymentsData = {
  data: [
    { id: 'p1', amount: 100, currency: 'USD', status: 'succeeded', payment_method: 'card', created_at: '2026-01-01T00:00:00Z' },
    { id: 'p2', amount: 50, currency: 'BDT', status: 'pending', payment_method: 'wallet', created_at: '2026-01-02T00:00:00Z' },
  ],
  total: 2,
  limit: 10,
  offset: 0,
}

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

describe('PaymentsPage', () => {
  it('renders table with data on success', async () => {
    vi.spyOn(global, 'fetch').mockImplementation((url: string) => {
      if (url.includes('/payments')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: mockPaymentsData }),
        })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
    })
    render(<PaymentsPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      expect(screen.getByText('$100.00')).toBeInTheDocument()
      expect(screen.getByText('USD')).toBeInTheDocument()
      expect(screen.getByText('BDT')).toBeInTheDocument()
    })
  })

  it('shows empty state when no data', async () => {
    vi.spyOn(global, 'fetch').mockImplementation((url: string) => {
      if (url.includes('/payments')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: { data: [], total: 0, limit: 10, offset: 0 } }),
        })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
    })
    render(<PaymentsPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      expect(screen.getByText(/no payments yet/i)).toBeInTheDocument()
    })
  })

  it('renders new payment link', () => {
    render(<PaymentsPage />, { wrapper: createWrapper() })
    expect(screen.getByText(/new payment/i)).toBeInTheDocument()
  })
})
