import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { Modal } from './Modal'

describe('Modal', () => {
  it('renders when open', () => {
    render(<Modal open onClose={vi.fn()} title="Test Modal">Content</Modal>)
    expect(screen.getByText('Test Modal')).toBeInTheDocument()
    expect(screen.getByText('Content')).toBeInTheDocument()
  })

  it('does not render when closed', () => {
    render(<Modal open={false} onClose={vi.fn()} title="Test">Content</Modal>)
    expect(screen.queryByText('Content')).not.toBeInTheDocument()
  })

  it('closes on backdrop click', async () => {
    const handleClose = vi.fn()
    const user = userEvent.setup()
    render(<Modal open onClose={handleClose} title="Test">Content</Modal>)
    await user.click(screen.getByRole('button', { name: /close/i }))
    expect(handleClose).toHaveBeenCalled()
  })

  it('calls onClose on Escape key', async () => {
    const handleClose = vi.fn()
    const user = userEvent.setup()
    render(<Modal open onClose={handleClose} title="Test">Content</Modal>)
    await user.keyboard('{Escape}')
    expect(handleClose).toHaveBeenCalled()
  })

  it('renders footer', () => {
    render(<Modal open onClose={vi.fn()} title="Test" footer={<button>Save</button>}>Content</Modal>)
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
  })
})
