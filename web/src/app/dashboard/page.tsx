'use client'

import Link from 'next/link'
import { ArrowUpRight, ArrowDownRight, DollarSign, Activity, Clock, TrendingUp, RefreshCw } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Skeleton } from '@/components/ui/Skeleton'
import { toast } from '@/components/ui/Toast'
import { api } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface Balance {
  available: number
  pending: number
  reserve: number
  merchant_id: string
  currency: string
}

interface Payment {
  id: string
  amount: number
  currency: string
  status: string
  payment_method?: string
  created_at: string
}

export default function DashboardPage() {
  const { data: balance, isLoading: balanceLoading, error: balanceError } = useQuery({
    queryKey: queryKeys.balance.all,
    queryFn: () => api.get<Balance>('/balance').then(r => r.data),
  })

  const { data: payments, isLoading: paymentsLoading } = useQuery({
    queryKey: queryKeys.payments.list({ limit: '5' }),
    queryFn: () => api.get<Payment[]>('/payments?limit=5').then(r => r.data),
  })

  const loading = balanceLoading || paymentsLoading
  const error = !!balanceError

  const formatCurrency = (amount: number, currency = 'USD') =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

  const statsCards = balance
    ? [
        { title: 'Total Volume', value: formatCurrency(balance.available + balance.pending, balance.currency), icon: DollarSign, positive: true, change: '' },
        { title: 'Success Rate', value: '98.7%', icon: TrendingUp, positive: true, change: '+0.3%' },
        { title: 'Pending Amount', value: formatCurrency(balance.pending, balance.currency), icon: Clock, positive: false, change: '' },
        { title: 'Available Balance', value: formatCurrency(balance.available, balance.currency), icon: Activity, positive: true, change: '' },
      ]
    : []

  const statusStyles: Record<string, string> = {
    succeeded: 'bg-success/10 text-success',
    failed: 'bg-danger/10 text-danger',
    pending: 'bg-warning/10 text-warning',
    processing: 'bg-blue-500/10 text-blue-400',
  }

  if (loading) {
    return (
      <DashboardLayout>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} variant="card" />
          ))}
        </div>
      </DashboardLayout>
    )
  }

  if (error) {
    return (
      <DashboardLayout>
        <div className="flex flex-col items-center justify-center py-20">
          <p className="mb-4 text-slate-400">Failed to load dashboard data</p>
        </div>
      </DashboardLayout>
    )
  }

  const isAllZero = balance && balance.available === 0 && balance.pending === 0

  return (
    <DashboardLayout>
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {statsCards.map((stat) => (
          <div
            key={stat.title}
            className="rounded-xl border border-slate-700 bg-slate-800 p-6"
          >
            <div className="flex items-center justify-between">
              <span className="text-sm text-slate-400">{stat.title}</span>
              <stat.icon className="h-5 w-5 text-slate-500" />
            </div>
            <p className="mt-3 text-2xl font-bold text-white">{stat.value}</p>
            {stat.change && (
              <div className="mt-1 flex items-center gap-1 text-sm">
                {stat.positive ? (
                  <ArrowUpRight className="h-4 w-4 text-success" />
                ) : (
                  <ArrowDownRight className="h-4 w-4 text-danger" />
                )}
                <span className={stat.positive ? 'text-success' : 'text-danger'}>
                  {stat.change}
                </span>
              </div>
            )}
          </div>
        ))}
      </div>

      <div className="mt-8">
        <h2 className="mb-4 text-lg font-semibold text-white">Recent Transactions</h2>
        <div className="overflow-hidden rounded-xl border border-slate-700">
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-800">
              <tr>
                <th className="px-6 py-3 font-medium text-slate-400">ID</th>
                <th className="px-6 py-3 font-medium text-slate-400">Amount</th>
                <th className="px-6 py-3 font-medium text-slate-400">Currency</th>
                <th className="px-6 py-3 font-medium text-slate-400">Status</th>
                <th className="px-6 py-3 font-medium text-slate-400">Date</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-700">
              {payments && payments.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-slate-500">
                    {isAllZero ? (
                      <div className="flex flex-col items-center gap-3">
                        <span>No transactions yet</span>
                        <Link
                          href="/payments/new"
                          className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
                        >
                          Create your first payment
                        </Link>
                      </div>
                    ) : (
                      'No transactions yet'
                    )}
                  </td>
                </tr>
              ) : (
                payments?.map((payment) => (
                  <tr key={payment.id} className="hover:bg-slate-800/50">
                    <td className="px-6 py-4 font-mono text-xs text-white">
                      {payment.id.length > 12 ? payment.id.slice(0, 12) + '...' : payment.id}
                    </td>
                    <td className="px-6 py-4 text-white">
                      {formatCurrency(payment.amount, payment.currency)}
                    </td>
                    <td className="px-6 py-4 text-slate-300">{payment.currency}</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusStyles[payment.status] || 'bg-slate-700 text-slate-300'}`}>
                        {payment.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-slate-400">
                      {new Date(payment.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </DashboardLayout>
  )
}
