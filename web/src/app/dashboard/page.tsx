'use client'

import { useEffect, useState } from 'react'
import { ArrowUpRight, ArrowDownRight, DollarSign, Activity, Clock, TrendingUp } from 'lucide-react'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { api } from '@/lib/api'

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
  const [balance, setBalance] = useState<Balance | null>(null)
  const [payments, setPayments] = useState<Payment[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get<Balance>('/balance'),
      api.get<Payment[]>('/payments?limit=5'),
    ])
      .then(([balanceRes, paymentsRes]) => {
        setBalance(balanceRes.data)
        setPayments(paymentsRes.data)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

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

  return (
    <DashboardLayout>
      {loading ? (
        <div className="flex items-center justify-center py-20">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        </div>
      ) : (
        <>
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
                  {payments.length === 0 ? (
                    <tr>
                      <td colSpan={5} className="px-6 py-8 text-center text-slate-500">
                        No transactions yet
                      </td>
                    </tr>
                  ) : (
                    payments.map((payment) => (
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
        </>
      )}
    </DashboardLayout>
  )
}
