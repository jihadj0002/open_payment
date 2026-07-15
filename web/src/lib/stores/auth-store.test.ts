import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuthStore } from './auth-store'
import { api } from '@/lib/api'

const mockLoginResponse = {
  data: {
    token_pair: { access_token: 'test-token' },
    user: { merchant_id: 'm1', role: 'merchant', permissions: ['read'] },
  },
}

const mockRegisterResponse = {
  data: {
    token_pair: { access_token: 'reg-token' },
    user: { merchant_id: 'm2', role: 'merchant', permissions: ['read'] },
  },
}

const mockProfileResponse = {
  data: { id: 'm1', name: 'Test', email: 'test@test.com' },
}

describe('auth-store', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: true,
    })
    localStorage.clear()
  })

  it('login sets user, token, isAuthenticated', async () => {
    vi.spyOn(api, 'post').mockResolvedValue(mockLoginResponse)
    await useAuthStore.getState().login('test@test.com', 'password')
    const state = useAuthStore.getState()
    expect(state.token).toBe('test-token')
    expect(state.isAuthenticated).toBe(true)
    expect(state.user?.merchant_id).toBe('m1')
    expect(localStorage.getItem('auth_token')).toBe('test-token')
  })

  it('register sets user, token, isAuthenticated', async () => {
    vi.spyOn(api, 'post').mockResolvedValue(mockRegisterResponse)
    await useAuthStore.getState().register('Test', 'test@test.com', 'password')
    const state = useAuthStore.getState()
    expect(state.token).toBe('reg-token')
    expect(state.isAuthenticated).toBe(true)
    expect(state.user?.merchant_id).toBe('m2')
    expect(localStorage.getItem('auth_token')).toBe('reg-token')
  })

  it('logout clears user, token, isAuthenticated and removes localStorage', () => {
    useAuthStore.setState({
      user: { merchant_id: 'm1', role: 'merchant', permissions: ['read'] },
      token: 'test-token',
      isAuthenticated: true,
    })
    localStorage.setItem('auth_token', 'test-token')
    useAuthStore.getState().logout()
    const state = useAuthStore.getState()
    expect(state.user).toBeNull()
    expect(state.token).toBeNull()
    expect(state.isAuthenticated).toBe(false)
    expect(localStorage.getItem('auth_token')).toBeNull()
  })

  it('fetchProfile populates user on success', async () => {
    vi.spyOn(api, 'get').mockResolvedValue(mockProfileResponse)
    useAuthStore.setState({ isLoading: true })
    await useAuthStore.getState().fetchProfile()
    const state = useAuthStore.getState()
    expect(state.user?.merchant_id).toBe('m1')
    expect(state.isLoading).toBe(false)
  })

  it('fetchProfile clears state on error', async () => {
    vi.spyOn(api, 'get').mockRejectedValue(new Error('fail'))
    useAuthStore.setState({ token: 'bad-token', isLoading: true })
    localStorage.setItem('auth_token', 'bad-token')
    await useAuthStore.getState().fetchProfile()
    const state = useAuthStore.getState()
    expect(state.user).toBeNull()
    expect(state.token).toBeNull()
    expect(state.isLoading).toBe(false)
  })
})
