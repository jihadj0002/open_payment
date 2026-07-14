'use client'

import { useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/context/AuthContext'
import { LayoutDashboard, CreditCard, Users, Building2, Key, Banknote, Settings, LogOut, Loader2 } from 'lucide-react'

const sidebarLinks = [
  { label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { label: 'Payments', href: '/payments', icon: CreditCard },
  { label: 'Customers', href: '/customers', icon: Users },
  { label: 'Merchants', href: '/merchants', icon: Building2 },
  { label: 'API Keys', href: '/api-keys', icon: Key },
  { label: 'Webhooks', href: '/webhooks', icon: Key },
  { label: 'Settlements', href: '/settlements', icon: Banknote },
  { label: 'Settings', href: '/settings', icon: Settings },
]

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { user, logout, isLoading } = useAuth()
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login')
    }
  }, [isLoading, user, router])

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-900">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!user) return null

  return (
    <div className="flex min-h-screen bg-slate-900">
      <aside className="w-64 border-r border-slate-700 bg-slate-800/50">
        <div className="p-6">
          <Link href="/dashboard" className="text-lg font-bold text-white">
            OpenPG
          </Link>
        </div>
        <nav className="space-y-1 px-3">
          {sidebarLinks.map((link) => {
            const isActive = pathname === link.href || (link.href !== '/dashboard' && pathname.startsWith(link.href))
            return (
              <Link
                key={link.href}
                href={link.href}
                className={`flex items-center gap-3 rounded-lg px-4 py-2.5 text-sm transition-colors ${
                  isActive
                    ? 'bg-primary/10 text-primary font-medium'
                    : 'text-slate-400 hover:bg-slate-700 hover:text-white'
                }`}
              >
                <link.icon className="h-4 w-4" />
                {link.label}
              </Link>
            )
          })}
        </nav>
      </aside>

      <div className="flex flex-1 flex-col">
        <header className="flex items-center justify-between border-b border-slate-700 px-8 py-4">
          <div>
            <h1 className="text-xl font-semibold text-white capitalize">
              {pathname === '/dashboard' ? 'Dashboard' : pathname.split('/')[1]?.replace(/-/g, ' ') || 'Dashboard'}
            </h1>
            <p className="text-sm text-slate-400">
              Welcome back, {user?.merchant_id?.slice(0, 8) || 'User'}
            </p>
          </div>
          <button
            onClick={() => { logout(); router.push('/login') }}
            className="flex items-center gap-2 rounded-lg px-4 py-2 text-sm text-slate-400 hover:bg-slate-700 hover:text-white transition-colors"
          >
            <LogOut className="h-4 w-4" />
            Logout
          </button>
        </header>

        <main className="flex-1 p-8">
          {children}
        </main>
      </div>
    </div>
  )
}
