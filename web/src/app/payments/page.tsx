'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Plus } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Badge } from '@/components/ui/Badge'
import { toast } from '@/components/ui/Toast'
import { api } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

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

export default function PaymentsPage() {
  const [page, setPage] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: queryKeys.payments.list({ page: String(page), limit: '10' }),
    queryFn: () => api.get<PaymentsResponse>(`/payments?page=${page}&limit=10`).then(r => r.data),
  })

  const formatCurrency = (amount: number, currency = 'USD') =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

  const statusBadge = (status: string) => {
    const variants: Record<string, 'success' | 'danger' | 'warning' | 'info'> = {
      succeeded: 'success',
      failed: 'danger',
      pending: 'warning',
      processing: 'info',
    }
    return <Badge variant={variants[status] || 'neutral'}>{status}</Badge>
  }

  const columns: Column<Payment>[] = [
    {
      key: 'id',
      header: 'ID',
      render: (p) => <span className="font-mono text-xs text-white">{p.id.length > 12 ? p.id.slice(0, 12) + '...' : p.id}</span>,
    },
    {
      key: 'amount',
      header: 'Amount',
      render: (p) => <span className="text-white">{formatCurrency(p.amount, p.currency)}</span>,
    },
    { key: 'currency', header: 'Currency' },
    {
      key: 'status',
      header: 'Status',
      render: (p) => statusBadge(p.status),
    },
    {
      key: 'payment_method',
      header: 'Method',
      render: (p) => <span className="capitalize text-slate-300">{p.payment_method?.replace(/_/g, ' ') || '-'}</span>,
    },
    {
      key: 'created_at',
      header: 'Created',
      render: (p) => new Date(p.created_at).toLocaleDateString(),
    },
  ]

  return (
    <DashboardLayout>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-white">Payments</h2>
        <Link
          href="/payments/new"
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          <Plus className="h-4 w-4" />
          New Payment
        </Link>
      </div>

      <DataTable<Payment>
        columns={columns}
        data={data?.payments || []}
        loading={isLoading}
        emptyMessage="No payments yet. Try creating one."
        pagination={{
          page,
          pageSize: 10,
          total: data?.total || 0,
          onPageChange: setPage,
        }}
      />
    </DashboardLayout>
  )
}
