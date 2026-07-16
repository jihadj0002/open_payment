export interface APIResponse<T> {
  data: T
  message?: string
}

export interface APIError {
  error: string
  message: string
  statusCode: number
}

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

let isRefreshing = false
let refreshPromise: Promise<string | null> | null = null

function decodeJWT(token: string): Record<string, unknown> | null {
  try {
    const payload = token.split('.')[1]
    return JSON.parse(atob(payload))
  } catch {
    return null
  }
}

function isTokenExpiringSoon(): boolean {
  if (typeof window === 'undefined') return false
  const token = localStorage.getItem('auth_token')
  if (!token) return false
  const decoded = decodeJWT(token)
  if (!decoded || typeof decoded.exp !== 'number') return false
  return decoded.exp * 1000 - Date.now() < 120000
}

async function refreshAccessToken(): Promise<string | null> {
  if (typeof window === 'undefined') return null

  if (isRefreshing && refreshPromise) {
    return refreshPromise
  }

  const refreshToken = localStorage.getItem('auth_refresh_token')
  if (!refreshToken) {
    clearAuthAndRedirect()
    return null
  }

  isRefreshing = true
  refreshPromise = (async () => {
    try {
      const apiPath = '/v1/auth/refresh'
      const response = await fetch(`${BASE_URL}${apiPath}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      })

      if (!response.ok) {
        clearAuthAndRedirect()
        return null
      }

      const result: APIResponse<{ access_token: string; refresh_token: string }> = await response.json()
      const { access_token, refresh_token: newRefreshToken } = result.data

      localStorage.setItem('auth_token', access_token)
      if (newRefreshToken) {
        localStorage.setItem('auth_refresh_token', newRefreshToken)
      }
      return access_token
    } catch {
      clearAuthAndRedirect()
      return null
    } finally {
      isRefreshing = false
      refreshPromise = null
    }
  })()

  return refreshPromise
}

function clearAuthAndRedirect() {
  if (typeof window === 'undefined') return
  localStorage.removeItem('auth_token')
  localStorage.removeItem('auth_refresh_token')
  window.location.href = '/login'
}

async function request<T>(
  method: string,
  path: string,
  body?: unknown,
  timeoutMs = 30000
): Promise<APIResponse<T>> {
  if (typeof window !== 'undefined') {
    if (isTokenExpiringSoon()) {
      const newToken = await refreshAccessToken()
      if (newToken) {
        return request<T>(method, path, body, timeoutMs)
      }
    }
  }

  const token =
    typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  if (['POST', 'PATCH', 'DELETE'].includes(method)) {
    headers['Idempotency-Key'] = crypto.randomUUID()
  }

  const apiPath = path.startsWith('/v1') ? path : `/v1${path}`
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs)

  try {
    const response = await fetch(`${BASE_URL}${apiPath}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
      signal: controller.signal,
    })

    if (!response.ok) {
      if (response.status === 401) {
        const newToken = await refreshAccessToken()
        if (newToken) {
          headers['Authorization'] = `Bearer ${newToken}`
          const retryResponse = await fetch(`${BASE_URL}${apiPath}`, {
            method,
            headers,
            body: body ? JSON.stringify(body) : undefined,
            signal: controller.signal,
          })
          if (retryResponse.ok) {
            return retryResponse.json()
          }
        }
        clearAuthAndRedirect()
        throw {
          error: 'AuthError',
          message: 'Session expired. Please login again.',
          statusCode: 401,
        } satisfies APIError
      }

      if (response.status === 429) {
        throw {
          error: 'RateLimitError',
          message: 'Too many requests. Please wait a moment and try again.',
          statusCode: 429,
        } satisfies APIError
      }

      const error: APIError = await response.json().catch(() => ({
        error: 'UnknownError',
        message: 'An unexpected error occurred',
        statusCode: response.status,
      }))
      throw error
    }

    return response.json()
  } catch (err: unknown) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw {
        error: 'TimeoutError',
        message: 'Request timed out. Please check your connection and try again.',
        statusCode: 408,
      } satisfies APIError
    }

    if (err instanceof TypeError && err.message === 'Failed to fetch') {
      throw {
        error: 'NetworkError',
        message: 'Unable to connect to the server. Please check your internet connection.',
        statusCode: 0,
      } satisfies APIError
    }

    throw err
  } finally {
    clearTimeout(timeoutId)
  }
}

export const api = {
  get: <T>(path: string) => request<T>('GET', path),
  post: <T>(path: string, body: unknown) => request<T>('POST', path, body),
  patch: <T>(path: string, body: unknown) => request<T>('PATCH', path, body),
  delete: <T>(path: string) => request<T>('DELETE', path),
}
