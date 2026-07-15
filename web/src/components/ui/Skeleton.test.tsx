import { render } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { Skeleton } from './Skeleton'

describe('Skeleton', () => {
  it('renders with text variant by default', () => {
    const { container } = render(<Skeleton />)
    const el = container.firstChild as HTMLElement
    expect(el.className).toContain('h-4')
    expect(el.className).toContain('animate-pulse')
  })

  it('applies card variant', () => {
    const { container } = render(<Skeleton variant="card" />)
    const el = container.firstChild as HTMLElement
    expect(el.className).toContain('h-24')
  })

  it('applies table-row variant', () => {
    const { container } = render(<Skeleton variant="table-row" />)
    const el = container.firstChild as HTMLElement
    expect(el.className).toContain('h-12')
  })

  it('applies custom className', () => {
    const { container } = render(<Skeleton className="extra" />)
    const el = container.firstChild as HTMLElement
    expect(el.className).toContain('extra')
  })
})
