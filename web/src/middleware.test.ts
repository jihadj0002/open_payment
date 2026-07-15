import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockNextResponse = {
  next: vi.fn(() => 'next'),
  redirect: vi.fn(() => 'redirect'),
}

vi.mock('next/server', () => ({
  NextResponse: mockNextResponse,
}))

function mockRequest(url: string, cookieValue?: string) {
  const urlObj = new URL(url, 'http://localhost:3000')
  return {
    url: urlObj.href,
    nextUrl: urlObj,
    cookies: {
      get: vi.fn((name: string) => {
        if (name === 'auth_token' && cookieValue) {
          return { value: cookieValue }
        }
        return undefined
      }),
    },
  }
}

describe('middleware', () => {
  beforeEach(() => {
    vi.resetModules()
    mockNextResponse.next.mockClear()
    mockNextResponse.redirect.mockClear()
  })

  it('redirects to login without auth_token on protected route', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/merchant/dashboard')
    const result = middleware(req)
    expect(result).toBe('redirect')
    expect(mockNextResponse.redirect).toHaveBeenCalledWith(
      expect.objectContaining({ href: 'http://localhost:3000/login' })
    )
  })

  it('allows public routes without token', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/login')
    const result = middleware(req)
    expect(result).toBe('next')
  })

  it('allows root path without token', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/')
    const result = middleware(req)
    expect(result).toBe('next')
  })

  it('allows register path without token', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/register')
    const result = middleware(req)
    expect(result).toBe('next')
  })

  it('allows forgot-password path without token', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/forgot-password')
    const result = middleware(req)
    expect(result).toBe('next')
  })

  it('allows reset-password path without token', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/reset-password')
    const result = middleware(req)
    expect(result).toBe('next')
  })

  it('allows protected routes with auth_token', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/merchant/dashboard', 'valid-token')
    const result = middleware(req)
    expect(result).toBe('next')
  })

  it('allows static assets', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/_next/static/chunk.js')
    const result = middleware(req)
    expect(result).toBe('next')
  })

  it('allows favicon', async () => {
    const { middleware } = await import('./middleware')
    const req = mockRequest('/favicon.ico')
    const result = middleware(req)
    expect(result).toBe('next')
  })
})
