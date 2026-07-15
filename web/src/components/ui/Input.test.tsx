import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { Input } from './Input'

describe('Input', () => {
  it('renders label', () => {
    render(<Input id="email" label="Email" />)
    expect(screen.getByLabelText('Email')).toBeInTheDocument()
  })

  it('renders error message', () => {
    render(<Input id="email" label="Email" error="Required" />)
    expect(screen.getByText('Required')).toBeInTheDocument()
  })

  it('calls onChange handler', async () => {
    const handleChange = vi.fn()
    const user = userEvent.setup()
    render(<Input id="name" onChange={handleChange} />)
    await user.type(screen.getByRole('textbox'), 'a')
    expect(handleChange).toHaveBeenCalled()
  })

  it('forwards ref', () => {
    const ref = { current: null }
    render(<Input id="test" ref={ref} />)
    expect(ref.current).toBeInstanceOf(HTMLInputElement)
  })

  it('applies error styling', () => {
    render(<Input id="test" error="Error" />)
    const input = screen.getByRole('textbox')
    expect(input.className).toContain('border-danger')
  })
})
