import { render, screen, act, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { Toaster, toast } from './Toast'

describe('Toast', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders success toast', () => {
    render(<Toaster />)
    act(() => { toast.success('Success message') })
    expect(screen.getByText('Success message')).toBeInTheDocument()
  })

  it('renders error toast', () => {
    render(<Toaster />)
    act(() => { toast.error('Error message') })
    expect(screen.getByText('Error message')).toBeInTheDocument()
  })

  it('renders info toast', () => {
    render(<Toaster />)
    act(() => { toast.info('Info message') })
    expect(screen.getByText('Info message')).toBeInTheDocument()
  })

  it('auto-dismisses after duration', () => {
    render(<Toaster />)
    act(() => { toast.success('Auto dismiss') })
    expect(screen.getByText('Auto dismiss')).toBeInTheDocument()
    act(() => { vi.advanceTimersByTime(4000) })
    expect(screen.queryByText('Auto dismiss')).not.toBeInTheDocument()
  })

  it('removes toast on close button click', () => {
    render(<Toaster />)
    act(() => { toast.success('Dismiss me') })
    expect(screen.getByText('Dismiss me')).toBeInTheDocument()
    const container = screen.getByText('Dismiss me').closest('[class*="flex"]')
    const closeBtn = container?.querySelector('button')
    act(() => {
      if (closeBtn) fireEvent.click(closeBtn)
    })
    expect(screen.queryByText('Dismiss me')).not.toBeInTheDocument()
  })
})
