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
  body?: unknown
): Promise<APIResponse<T>> {
  const token =
    typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem('auth_token')
      if (typeof window !== 'undefined') {
        window.location.href = '/login'
      }
    }
    const error: APIError = await response.json().catch(() => ({
      error: 'UnknownError',
      message: 'An unexpected error occurred',
      statusCode: response.status,
    }))
    throw error
  }

  return response.json()
}

export const api = {
  get: <T>(path: string) => request<T>('GET', path),
  post: <T>(path: string, body: unknown) => request<T>('POST', path, body),
  patch: <T>(path: string, body: unknown) => request<T>('PATCH', path, body),
  delete: <T>(path: string) => request<T>('DELETE', path),
}
