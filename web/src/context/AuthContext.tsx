'use client'
import { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react'
import { api } from '@/lib/api'

interface MerchantProfile {
  id: string
  name: string
  email: string
  webhook_url?: string
  status: string
  created_at: string
  updated_at: string
}

interface User {
  merchant_id: string
  role: string
  permissions: string[]
}

interface AuthContextType {
  user: User | null
  token: string | null
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string) => Promise<void>
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const storedToken = localStorage.getItem('auth_token')
    if (storedToken) {
      setToken(storedToken)
      api.get<MerchantProfile>('/merchants/profile')
        .then((res) => {
          const m = res.data
          setUser({
            merchant_id: m.id,
            role: 'merchant',
            permissions: ['read', 'write'],
          })
        })
        .catch(() => {
          localStorage.removeItem('auth_token')
          setToken(null)
        })
        .finally(() => setIsLoading(false))
    } else {
      setIsLoading(false)
    }
  }, [])

  const login = useCallback(async (email: string, password: string) => {
    const res = await api.post<{ access_token: string; user: User }>('/auth/login', { email, password })
    localStorage.setItem('auth_token', res.data.access_token)
    setToken(res.data.access_token)
    setUser(res.data.user)
  }, [])

  const register = useCallback(async (name: string, email: string, password: string) => {
    const res = await api.post<{ access_token: string; user: User }>('/auth/register', { name, email, password })
    localStorage.setItem('auth_token', res.data.access_token)
    setToken(res.data.access_token)
    setUser(res.data.user)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem('auth_token')
    setToken(null)
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider value={{ user, token, login, register, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}
