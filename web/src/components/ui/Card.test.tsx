import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { Card } from './Card'

describe('Card', () => {
  it('renders children', () => {
    render(<Card>Content</Card>)
    expect(screen.getByText('Content')).toBeInTheDocument()
  })

  it('renders header', () => {
    render(<Card header={<h2>Header</h2>}>Body</Card>)
    expect(screen.getByText('Header')).toBeInTheDocument()
  })

  it('renders footer', () => {
    render(<Card footer={<div>Footer</div>}>Body</Card>)
    expect(screen.getByText('Footer')).toBeInTheDocument()
  })

  it('applies custom className', () => {
    render(<Card className="custom-class">Content</Card>)
    expect(screen.getByText('Content').parentElement?.className).toContain('custom-class')
  })
})
