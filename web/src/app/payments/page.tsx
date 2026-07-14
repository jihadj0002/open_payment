'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { Plus } from 'lucide-react'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { api } from '@/lib/api'

interface Payment {
  id: string
  amount: number
  currency: string
  status: string
  payment_method?: string
  description?: string
  created_at: string
}

interface PaymentsResponse {
  payments: Payment[]
  total: number
  page: number
  limit: number
}

const statusStyles: Record<string, string> = {
  succeeded: 'bg-success/10 text-success',
  failed: 'bg-danger/10 text-danger',
  pending: 'bg-warning/10 text-warning',
  processing: 'bg-blue-500/10 text-blue-400',
}

export default function PaymentsPage() {
  const [data, setData] = useState<PaymentsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)

  useEffect(() => {
    setLoading(true)
    api.get<PaymentsResponse>(`/payments?page=${page}&limit=10`)
      .then((res) => setData(res.data))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [page])

  const formatCurrency = (amount: number, currency = 'USD') =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

  return (
    <DashboardLayout>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-white">Payments</h2>
        <Link
          href="/payments/new"
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          New Payment
        </Link>
      </div>

      <div className="overflow-hidden rounded-xl border border-slate-700">
        <table className="w-full text-left text-sm">
          <thead className="bg-slate-800">
            <tr>
              <th className="px-6 py-3 font-medium text-slate-400">ID</th>
              <th className="px-6 py-3 font-medium text-slate-400">Amount</th>
              <th className="px-6 py-3 font-medium text-slate-400">Currency</th>
              <th className="px-6 py-3 font-medium text-slate-400">Status</th>
              <th className="px-6 py-3 font-medium text-slate-400">Method</th>
              <th className="px-6 py-3 font-medium text-slate-400">Created</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {loading ? (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-slate-500">
                  <div className="flex items-center justify-center gap-2">
                    <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                    Loading...
                  </div>
                </td>
              </tr>
            ) : data?.payments.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-slate-500">
                  No payments found
                </td>
              </tr>
            ) : (
              data?.payments.map((payment) => (
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
                  <td className="px-6 py-4 text-slate-300 capitalize">
                    {payment.payment_method?.replace(/_/g, ' ') || '-'}
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

      {data && data.total > data.limit && (
        <div className="mt-4 flex items-center justify-between text-sm text-slate-400">
          <span>
            Page {data.page} of {Math.ceil(data.total / data.limit)}
          </span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1}
              className="rounded-lg border border-slate-700 px-3 py-1.5 hover:bg-slate-800 disabled:opacity-50"
            >
              Previous
            </button>
            <button
              onClick={() => setPage((p) => p + 1)}
              disabled={page >= Math.ceil(data.total / data.limit)}
              className="rounded-lg border border-slate-700 px-3 py-1.5 hover:bg-slate-800 disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
