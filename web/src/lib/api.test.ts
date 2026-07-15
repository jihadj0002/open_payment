import { describe, it, expect, vi, beforeEach } from 'vitest'
import { api, APIError } from './api'

const mockFetch = vi.fn()
global.fetch = mockFetch

beforeEach(() => {
  mockFetch.mockReset()
  localStorage.clear()
})

describe('api', () => {
  it('attaches Bearer token from localStorage', async () => {
    localStorage.setItem('auth_token', 'my-token')
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'ok' }),
    })
    await api.get('/test')
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/test'),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer my-token',
        }),
      })
    )
  })

  it('includes Idempotency-Key for POST', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'ok' }),
    })
    await api.post('/test', { foo: 'bar' })
    const call = mockFetch.mock.calls[0]
    const headers = call[1].headers
    expect(headers['Idempotency-Key']).toBeDefined()
    expect(headers['Idempotency-Key']).toBeTruthy()
  })

  it('includes Idempotency-Key for PATCH', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'ok' }),
    })
    await api.patch('/test', { foo: 'bar' })
    const call = mockFetch.mock.calls[0]
    expect(call[1].headers['Idempotency-Key']).toBeDefined()
  })

  it('includes Idempotency-Key for DELETE', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'ok' }),
    })
    await api.delete('/test')
    const call = mockFetch.mock.calls[0]
    expect(call[1].headers['Idempotency-Key']).toBeDefined()
  })

  it('throws APIError on non-ok response', async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: () => Promise.resolve({ error: 'BadRequest', message: 'Bad request', statusCode: 400 }),
    })
    await expect(api.get('/test')).rejects.toMatchObject({
      error: 'BadRequest',
      message: 'Bad request',
      statusCode: 400,
    })
  })

  it('clears token and redirects on 401', async () => {
    const { getComputedStyle } = window
    delete (window as any).location
    window.location = { href: '' } as any
    Object.defineProperty(window, 'location', { value: { href: '' }, writable: true })

    localStorage.setItem('auth_token', 'expired')
    mockFetch.mockResolvedValue({
      ok: false,
      status: 401,
      json: () => Promise.resolve({ error: 'Unauthorized', message: 'Invalid token', statusCode: 401 }),
    })
    await expect(api.get('/protected')).rejects.toThrow()
    expect(localStorage.getItem('auth_token')).toBeNull()
  })

  it('calls api.get with GET method', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'ok' }),
    })
    const result = await api.get('/test')
    expect(result).toEqual({ data: 'ok' })
    expect(mockFetch.mock.calls[0][1].method).toBe('GET')
  })

  it('calls api.post with POST method', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'created' }),
    })
    const result = await api.post('/test', { foo: 'bar' })
    expect(result).toEqual({ data: 'created' })
    expect(mockFetch.mock.calls[0][1].method).toBe('POST')
  })

  it('calls api.patch with PATCH method', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'updated' }),
    })
    const result = await api.patch('/test', { foo: 'bar' })
    expect(result).toEqual({ data: 'updated' })
    expect(mockFetch.mock.calls[0][1].method).toBe('PATCH')
  })

  it('calls api.delete with DELETE method', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'deleted' }),
    })
    const result = await api.delete('/test')
    expect(result).toEqual({ data: 'deleted' })
    expect(mockFetch.mock.calls[0][1].method).toBe('DELETE')
  })
})
