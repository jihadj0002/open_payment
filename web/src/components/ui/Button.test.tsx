import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { Button } from './Button'

describe('Button', () => {
  it('renders children', () => {
    render(<Button>Click me</Button>)
    expect(screen.getByRole('button', { name: /click me/i })).toBeInTheDocument()
  })

  it('applies variant classes', () => {
    render(<Button variant="danger">Delete</Button>)
    const btn = screen.getByRole('button')
    expect(btn.className).toContain('bg-danger')
  })

  it('applies size classes', () => {
    render(<Button size="lg">Large</Button>)
    const btn = screen.getByRole('button')
    expect(btn.className).toContain('h-12')
  })

  it('handles disabled state', () => {
    render(<Button disabled>Disabled</Button>)
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('calls onClick handler', async () => {
    const handleClick = vi.fn()
    const user = userEvent.setup()
    render(<Button onClick={handleClick}>Click</Button>)
    await user.click(screen.getByRole('button'))
    expect(handleClick).toHaveBeenCalledTimes(1)
  })

  it('forwards ref', () => {
    const ref = { current: null }
    render(<Button ref={ref}>Ref</Button>)
    expect(ref.current).toBeInstanceOf(HTMLButtonElement)
  })
})
