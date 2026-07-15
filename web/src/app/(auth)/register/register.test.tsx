import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import RegisterPage from './page'

const mockPush = vi.fn()
const mockRegister = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}))

vi.mock('@/context/AuthContext', () => ({
  useAuth: () => ({
    register: mockRegister,
  }),
}))

describe('RegisterPage', () => {
  beforeEach(() => {
    mockPush.mockClear()
    mockRegister.mockClear()
  })

  it('renders all form fields', () => {
    render(<RegisterPage />)
    expect(screen.getByLabelText(/full name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/^password/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument()
  })

  it('shows validation errors for empty fields', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)
    await user.click(screen.getByRole('button', { name: /create account/i }))
    await waitFor(() => {
      expect(screen.getByText(/at least 2 characters/i)).toBeInTheDocument()
      expect(screen.getByText(/invalid email/i)).toBeInTheDocument()
    })
  })

  it('shows error when passwords do not match', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)
    await user.type(screen.getByLabelText(/full name/i), 'Test User')
    await user.type(screen.getByLabelText(/email/i), 'test@test.com')
    await user.type(screen.getAllByLabelText(/password/i)[0], 'password1')
    await user.type(screen.getAllByLabelText(/password/i)[1], 'password2')
    await user.click(screen.getByRole('button', { name: /create account/i }))
    await waitFor(() => {
      expect(screen.getByText(/passwords do not match/i)).toBeInTheDocument()
    })
  })

  it('calls register on valid submit', async () => {
    mockRegister.mockResolvedValue(undefined)
    const user = userEvent.setup()
    render(<RegisterPage />)
    await user.type(screen.getByLabelText(/full name/i), 'Test User')
    await user.type(screen.getByLabelText(/email/i), 'test@test.com')
    await user.type(screen.getAllByLabelText(/password/i)[0], 'password123')
    await user.type(screen.getAllByLabelText(/password/i)[1], 'password123')
    await user.click(screen.getByRole('button', { name: /create account/i }))
    await waitFor(() => {
      expect(mockRegister).toHaveBeenCalledWith('Test User', 'test@test.com', 'password123')
    })
  })

  it('shows error on failed registration', async () => {
    mockRegister.mockRejectedValue(new Error('Email already taken'))
    const user = userEvent.setup()
    render(<RegisterPage />)
    await user.type(screen.getByLabelText(/full name/i), 'Test User')
    await user.type(screen.getByLabelText(/email/i), 'taken@test.com')
    await user.type(screen.getAllByLabelText(/password/i)[0], 'password123')
    await user.type(screen.getAllByLabelText(/password/i)[1], 'password123')
    await user.click(screen.getByRole('button', { name: /create account/i }))
    await waitFor(() => {
      expect(screen.getByText(/email already taken/i)).toBeInTheDocument()
    })
  })

  it('renders link to login page', () => {
    render(<RegisterPage />)
    expect(screen.getByRole('link', { name: /sign in/i })).toHaveAttribute('href', '/login')
  })
})
