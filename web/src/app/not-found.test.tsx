import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import NotFound from './not-found'

describe('NotFound', () => {
  it('renders 404 heading', () => {
    render(<NotFound />)
    expect(screen.getByText('404')).toBeInTheDocument()
  })

  it('renders go home link pointing to /', () => {
    render(<NotFound />)
    const link = screen.getByRole('link', { name: /go to dashboard/i })
    expect(link).toHaveAttribute('href', '/')
  })
})
