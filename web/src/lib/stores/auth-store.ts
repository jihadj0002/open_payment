import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { api } from '@/lib/api'

function setCookie(name: string, value: string, days = 7) {
  const expires = new Date(Date.now() + days * 864e5).toUTCString()
  document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Lax`
}

function setSessionCookie(name: string, value: string) {
  document.cookie = `${name}=${encodeURIComponent(value)}; path=/; SameSite=Lax`
}

function removeCookie(name: string) {
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; SameSite=Lax`
}

interface User {
  merchant_id: string
  role: string
  permissions: string[]
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string) => Promise<void>
  logout: () => void
  setUser: (user: User | null) => void
  setToken: (token: string | null) => void
  setLoading: (loading: boolean) => void
  fetchProfile: () => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (email: string, password: string) => {
        const res = await api.post<{ token_pair: { access_token: string; refresh_token: string }; user: User }>('/auth/login', { email, password })
        const token = res.data.token_pair.access_token
        const refreshToken = res.data.token_pair.refresh_token
        localStorage.setItem('auth_token', token)
        localStorage.setItem('auth_refresh_token', refreshToken)
        setCookie('auth_token', token)
        setSessionCookie('just_logged_in', '1')
        set({ user: res.data.user, token, isAuthenticated: true })
      },

      register: async (name: string, email: string, password: string) => {
        const res = await api.post<{ token_pair: { access_token: string; refresh_token: string }; user: User }>('/auth/register', { name, email, password })
        const token = res.data.token_pair.access_token
        const refreshToken = res.data.token_pair.refresh_token
        localStorage.setItem('auth_token', token)
        localStorage.setItem('auth_refresh_token', refreshToken)
        setCookie('auth_token', token)
        setSessionCookie('just_logged_in', '1')
        set({ user: res.data.user, token, isAuthenticated: true })
      },

      logout: () => {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_refresh_token')
        removeCookie('auth_token')
        set({ user: null, token: null, isAuthenticated: false })
      },

      setUser: (user) => set({ user, isAuthenticated: !!user }),
      setToken: (token) => set({ token }),
      setLoading: (isLoading) => set({ isLoading }),

      fetchProfile: async () => {
        try {
          const res = await api.get<{ id: string; name: string; email: string }>('/merchants/profile')
          set({
            user: {
              merchant_id: res.data.id,
              role: 'merchant',
              permissions: ['read', 'write'],
            },
            isLoading: false,
          })
        } catch {
          localStorage.removeItem('auth_token')
          set({ user: null, token: null, isLoading: false })
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
