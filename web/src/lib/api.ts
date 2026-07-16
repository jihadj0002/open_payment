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

async function request<T>(
  method: string,
  path: string,
  body?: unknown,
  timeoutMs = 30000
): Promise<APIResponse<T>> {
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
        localStorage.removeItem('auth_token')
        if (typeof window !== 'undefined') {
          window.location.href = '/login'
        }
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
