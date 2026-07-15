import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import SettingsPage from './page'

const mockProfile = {
  id: 'm1',
  name: 'Test Merchant',
  email: 'merchant@test.com',
  webhook_url: '',
  status: 'active',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
}

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

describe('SettingsPage', () => {
  it('renders profile form on success', async () => {
    vi.spyOn(global, 'fetch').mockImplementation((url: string) => {
      if (url.includes('/merchants/profile')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: mockProfile }),
        })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
    })
    render(<SettingsPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      expect(screen.getByText(/settings/i)).toBeInTheDocument()
      expect(screen.getByDisplayValue('Test Merchant')).toBeInTheDocument()
      expect(screen.getByText('merchant@test.com')).toBeInTheDocument()
    })
  })

  it('shows loading skeleton initially', () => {
    vi.spyOn(global, 'fetch').mockImplementation(() => new Promise(() => {}))
    render(<SettingsPage />, { wrapper: createWrapper() })
    const skeletons = document.querySelectorAll('.animate-pulse')
    expect(skeletons.length).toBeGreaterThan(0)
  })

  it('disables save button when name unchanged', async () => {
    vi.spyOn(global, 'fetch').mockImplementation((url: string) => {
      if (url.includes('/merchants/profile')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: mockProfile }),
        })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
    })
    render(<SettingsPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      const saveBtn = screen.getByRole('button', { name: /save changes/i })
      expect(saveBtn).toBeDisabled()
    })
  })

  it('enables save button when name changed', async () => {
    vi.spyOn(global, 'fetch').mockImplementation((url: string) => {
      if (url.includes('/merchants/profile')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ data: mockProfile }),
        })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
    })
    const user = userEvent.setup()
    render(<SettingsPage />, { wrapper: createWrapper() })
    await waitFor(() => {
      expect(screen.getByDisplayValue('Test Merchant')).toBeInTheDocument()
    })
    const input = screen.getByDisplayValue('Test Merchant')
    await user.clear(input)
    await user.type(input, 'Updated Merchant')
    await waitFor(() => {
      const saveBtn = screen.getByRole('button', { name: /save changes/i })
      expect(saveBtn).not.toBeDisabled()
    })
  })
})
