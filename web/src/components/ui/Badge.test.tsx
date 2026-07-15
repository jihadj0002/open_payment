import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { Badge } from './Badge'

describe('Badge', () => {
  it('renders children', () => {
    render(<Badge>Active</Badge>)
    expect(screen.getByText('Active')).toBeInTheDocument()
  })

  it('applies variant class for success', () => {
    render(<Badge variant="success">Success</Badge>)
    expect(screen.getByText('Success').className).toContain('bg-success/10')
  })

  it('applies variant class for danger', () => {
    render(<Badge variant="danger">Danger</Badge>)
    expect(screen.getByText('Danger').className).toContain('bg-danger/10')
  })

  it('applies variant class for warning', () => {
    render(<Badge variant="warning">Warning</Badge>)
    expect(screen.getByText('Warning').className).toContain('bg-warning/10')
  })

  it('applies variant class for info', () => {
    render(<Badge variant="info">Info</Badge>)
    expect(screen.getByText('Info').className).toContain('bg-blue-500/10')
  })

  it('applies variant class for neutral', () => {
    render(<Badge variant="neutral">Neutral</Badge>)
    expect(screen.getByText('Neutral').className).toContain('bg-slate-700')
  })
})
