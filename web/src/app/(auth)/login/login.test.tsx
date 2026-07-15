import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import LoginPage from './page'

const mockPush = vi.fn()
const mockLogin = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}))

vi.mock('@/context/AuthContext', () => ({
  useAuth: () => ({
    login: mockLogin,
  }),
}))

describe('LoginPage', () => {
  beforeEach(() => {
    mockPush.mockClear()
    mockLogin.mockClear()
  })

  it('renders email and password fields', () => {
    render(<LoginPage />)
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('shows validation errors for empty fields', async () => {
    const user = userEvent.setup()
    render(<LoginPage />)
    await user.click(screen.getByRole('button', { name: /sign in/i }))
    await waitFor(() => {
      expect(screen.getByText(/invalid email/i)).toBeInTheDocument()
      expect(screen.getByText(/at least 6 characters/i)).toBeInTheDocument()
    })
  })

  it('calls login on valid submit and navigates to dashboard', async () => {
    mockLogin.mockResolvedValue(undefined)
    const user = userEvent.setup()
    render(<LoginPage />)
    await user.type(screen.getByLabelText(/email/i), 'test@test.com')
    await user.type(screen.getByLabelText(/password/i), 'password123')
    await user.click(screen.getByRole('button', { name: /sign in/i }))
    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith('test@test.com', 'password123')
      expect(mockPush).toHaveBeenCalledWith('/merchant/dashboard')
    })
  })

  it('shows error on failed login', async () => {
    mockLogin.mockRejectedValue(new Error('Invalid credentials'))
    const user = userEvent.setup()
    render(<LoginPage />)
    await user.type(screen.getByLabelText(/email/i), 'test@test.com')
    await user.type(screen.getByLabelText(/password/i), 'password123')
    await user.click(screen.getByRole('button', { name: /sign in/i }))
    await waitFor(() => {
      expect(screen.getByText(/invalid credentials/i)).toBeInTheDocument()
    })
  })

  it('renders link to register page', () => {
    render(<LoginPage />)
    expect(screen.getByRole('link', { name: /register/i })).toHaveAttribute('href', '/register')
  })
})
