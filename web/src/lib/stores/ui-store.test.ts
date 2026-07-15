import { describe, it, expect, beforeEach } from 'vitest'
import { useUIStore } from './ui-store'

describe('ui-store', () => {
  beforeEach(() => {
    useUIStore.setState({ sidebarOpen: true })
  })

  it('initializes with sidebarOpen true', () => {
    const state = useUIStore.getState()
    expect(state.sidebarOpen).toBe(true)
  })

  it('toggleSidebar flips sidebarOpen', () => {
    useUIStore.getState().toggleSidebar()
    expect(useUIStore.getState().sidebarOpen).toBe(false)
    useUIStore.getState().toggleSidebar()
    expect(useUIStore.getState().sidebarOpen).toBe(true)
  })
})
