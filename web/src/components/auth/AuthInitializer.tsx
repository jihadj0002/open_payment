'use client'
import { useEffect } from 'react'
import { useAuthStore } from '@/lib/stores/auth-store'

function setCookie(name: string, value: string, days = 7) {
  const expires = new Date(Date.now() + days * 864e5).toUTCString()
  document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Lax`
}

export function AuthInitializer({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((s) => s.token)
  const setToken = useAuthStore((s) => s.setToken)
  const setLoading = useAuthStore((s) => s.setLoading)
  const fetchProfile = useAuthStore((s) => s.fetchProfile)

  useEffect(() => {
    const storedToken = localStorage.getItem('auth_token')
    if (storedToken) {
      setCookie('auth_token', storedToken)
      setToken(storedToken)
      fetchProfile()
    } else {
      setLoading(false)
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return <>{children}</>
}
