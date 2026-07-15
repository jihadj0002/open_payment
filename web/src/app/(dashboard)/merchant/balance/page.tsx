'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Badge } from '@/components/ui/Badge'
import { Skeleton } from '@/components/ui/Skeleton'
import { Card } from '@/components/ui/Card'
import { toast } from '@/components/ui/Toast'
import { api } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface Balance {
  merchant_id: string
  currency: string
  available: number
  pending: number
  reserve: number
}

interface Transaction {
  id: string
  type: string
  amount: number
  currency: string
  balance_after: number
  description?: string
  created_at: string
}

interface Settlement {
  id: string
  amount: number
  currency: string
  status: string
  created_at: string
}

interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
  has_more: boolean
}

export default function BalancePage() {
  const [txPage, setTxPage] = useState(1)
  const [settlementPage, setSettlementPage] = useState(1)

  const { data: balance, isLoading: balanceLoading, error: balanceError } = useQuery({
    queryKey: queryKeys.balance.all,
    queryFn: () => api.get<Balance>('/balance').then(r => r.data),
  })

  const { data: txResponse, isLoading: txLoading, error: txError } = useQuery({
    queryKey: queryKeys.balance.transactions({ page: String(txPage), per_page: '20' }),
    queryFn: () => api.get<PaginatedResponse<Transaction>>(`/balance/transactions?page=${txPage}&per_page=20`).then(r => r.data),
  })

  const { data: settlementResponse, isLoading: settlementLoading, error: settlementError } = useQuery({
    queryKey: queryKeys.settlements.list({ page: String(settlementPage), per_page: '20' }),
    queryFn: () => api.get<PaginatedResponse<Settlement>>(`/settlements?page=${settlementPage}&per_page=20`).then(r => r.data),
  })

  useEffect(() => {
    if (balanceError) toast.error('Failed to load balance')
  }, [balanceError])

  useEffect(() => {
    if (txError) toast.error('Failed to load transactions')
  }, [txError])

  useEffect(() => {
    if (settlementError) toast.error('Failed to load settlements')
  }, [settlementError])

  const formatCurrency = (amount: number, currency = 'USD') =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

  const txColumns: Column<Transaction>[] = [
    {
      key: 'created_at',
      header: 'Date',
      render: (t) => new Date(t.created_at).toLocaleDateString(),
    },
    {
      key: 'type',
      header: 'Type',
      render: (t) => {
        const variants: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
          capture: 'success',
          refund: 'warning',
          settlement: 'info',
          fee: 'danger',
        }
        return <Badge variant={variants[t.type] || 'neutral'}>{t.type}</Badge>
      },
    },
    {
      key: 'amount',
      header: 'Amount',
      render: (t) => <span className="text-white">{formatCurrency(t.amount, t.currency)}</span>,
    },
    {
      key: 'balance_after',
      header: 'Balance After',
      render: (t) => <span className="text-slate-300">{formatCurrency(t.balance_after, t.currency)}</span>,
    },
  ]

  const settlementColumns: Column<Settlement>[] = [
    {
      key: 'id',
      header: 'ID',
      render: (s) => (
        <span className="font-mono text-xs text-white">
          {s.id.length > 12 ? `${s.id.slice(0, 12)}...` : s.id}
        </span>
      ),
    },
    {
      key: 'amount',
      header: 'Amount',
      render: (s) => <span className="text-white">{formatCurrency(s.amount, s.currency)}</span>,
    },
    {
      key: 'status',
      header: 'Status',
      render: (s) => {
        const variants: Record<string, 'success' | 'warning' | 'info'> = {
          completed: 'success',
          pending: 'warning',
          processing: 'info',
        }
        return <Badge variant={variants[s.status] || 'neutral'}>{s.status}</Badge>
      },
    },
    {
      key: 'created_at',
      header: 'Date',
      render: (s) => new Date(s.created_at).toLocaleDateString(),
    },
    {
      key: 'actions',
      header: 'Actions',
      render: () => (
        <button className="text-sm text-primary transition-colors hover:text-primary-700">
          View
        </button>
      ),
    },
  ]

  return (
    <>
      <div className="mb-6">
        <h2 className="text-lg font-semibold text-white">Balance</h2>
      </div>

      {balanceLoading ? (
        <div className="grid gap-4 sm:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} variant="card" />
          ))}
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-3">
          <Card>
            <p className="text-sm text-slate-400">Available</p>
            <p className="mt-2 text-2xl font-bold text-success">
              {formatCurrency(balance?.available ?? 0, balance?.currency || 'BDT')}
            </p>
          </Card>
          <Card>
            <p className="text-sm text-slate-400">Pending</p>
            <p className="mt-2 text-2xl font-bold text-warning">
              {formatCurrency(balance?.pending ?? 0, balance?.currency || 'BDT')}
            </p>
          </Card>
          <Card>
            <p className="text-sm text-slate-400">Reserve</p>
            <p className="mt-2 text-2xl font-bold text-blue-400">
              {formatCurrency(balance?.reserve ?? 0, balance?.currency || 'BDT')}
            </p>
          </Card>
        </div>
      )}

      <div className="mt-8">
        <h3 className="mb-4 text-lg font-semibold text-white">Transaction History</h3>
        <DataTable<Transaction>
          columns={txColumns}
          data={txResponse?.data || []}
          loading={txLoading}
          emptyMessage="No transactions yet"
          pagination={{
            page: txPage,
            pageSize: 20,
            total: txResponse?.total || 0,
            onPageChange: setTxPage,
          }}
        />
      </div>

      <div className="mt-8">
        <h3 className="mb-4 text-lg font-semibold text-white">Settlement History</h3>
        <DataTable<Settlement>
          columns={settlementColumns}
          data={settlementResponse?.data || []}
          loading={settlementLoading}
          emptyMessage="No settlements yet"
          pagination={{
            page: settlementPage,
            pageSize: 20,
            total: settlementResponse?.total || 0,
            onPageChange: setSettlementPage,
          }}
        />
      </div>
    </>
  )
}
