'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Skeleton } from '@/components/ui/Skeleton'
import { Badge } from '@/components/ui/Badge'
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

interface BalanceTransaction {
  date: string
  type: string
  amount: number
  balance_after: number
}

interface Settlement {
  id: string
  amount: number
  status: string
  date: string
}

interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
}

type Tab = 'transactions' | 'balance' | 'settlements'

const TABS: { key: Tab; label: string }[] = [
  { key: 'transactions', label: 'Transactions' },
  { key: 'balance', label: 'Balance' },
  { key: 'settlements', label: 'Settlements' },
]

export default function ReportsPage() {
  const [activeTab, setActiveTab] = useState<Tab>('transactions')
  const [txPage, setTxPage] = useState(1)
  const [settlementPage, setSettlementPage] = useState(1)

  const { data: balance, isLoading: balanceLoading, error: balanceError } = useQuery({
    queryKey: queryKeys.balance.all,
    queryFn: () => api.get<Balance>('/balance').then(r => r.data),
  })

  const { data: payments, isLoading: paymentsLoading, error: paymentsError } = useQuery({
    queryKey: queryKeys.payments.list({ limit: '100' }),
    queryFn: () => api.get<Payment[]>('/payments?limit=100').then(r => r.data),
  })

  const { data: txResponse, isLoading: txLoading, error: txError } = useQuery({
    queryKey: queryKeys.balance.transactions({ page: String(txPage), per_page: '50' }),
    queryFn: () => api.get<PaginatedResponse<BalanceTransaction>>(`/balance/transactions?page=${txPage}&per_page=50`).then(r => r.data),
  })

  const { data: settlementResponse, isLoading: settlementLoading, error: settlementError } = useQuery({
    queryKey: queryKeys.settlements.list({ page: String(settlementPage), per_page: '50' }),
    queryFn: () => api.get<PaginatedResponse<Settlement>>(`/settlements?page=${settlementPage}&per_page=50`).then(r => r.data),
  })

  if (balanceError) toast.error('Failed to load balance data')
  if (paymentsError) toast.error('Failed to load payments data')
  if (txError) toast.error('Failed to load transactions data')
  if (settlementError) toast.error('Failed to load settlement data')

  const formatCurrency = (amount: number, currency = 'BDT') =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

  const successfulPayments = payments?.filter(p => p.status === 'succeeded') || []
  const totalRevenue = successfulPayments.reduce((sum, p) => sum + p.amount, 0)
  const totalTransactions = payments?.length || 0

  const summaryCards = [
    {
      title: 'Total Revenue',
      value: formatCurrency(totalRevenue, balance?.currency || 'BDT'),
    },
    {
      title: 'Total Transactions',
      value: totalTransactions.toLocaleString(),
    },
    {
      title: 'Available Balance',
      value: balance ? formatCurrency(balance.available, balance.currency) : '-',
    },
    {
      title: 'Pending Balance',
      value: balance ? formatCurrency(balance.pending, balance.currency) : '-',
    },
  ]

  const statusBadge = (status: string) => {
    const variants: Record<string, 'success' | 'danger' | 'warning' | 'info'> = {
      succeeded: 'success',
      failed: 'danger',
      pending: 'warning',
      processing: 'info',
    }
    return <Badge variant={variants[status] || 'neutral'}>{status}</Badge>
  }

  const paymentColumns: Column<Payment>[] = [
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
      header: 'Payment Method',
      render: (p) => <span className="capitalize text-slate-300">{p.payment_method?.replace(/_/g, ' ') || '-'}</span>,
    },
    {
      key: 'created_at',
      header: 'Date',
      render: (p) => new Date(p.created_at).toLocaleDateString(),
    },
  ]

  const txColumns: Column<BalanceTransaction>[] = [
    {
      key: 'date',
      header: 'Date',
      render: (tx) => new Date(tx.date).toLocaleDateString(),
    },
    {
      key: 'type',
      header: 'Type',
      render: (tx) => <span className="capitalize text-slate-300">{tx.type.replace(/_/g, ' ')}</span>,
    },
    {
      key: 'amount',
      header: 'Amount',
      render: (tx) => {
        const isPositive = tx.amount >= 0
        return (
          <span className={isPositive ? 'text-success' : 'text-danger'}>
            {formatCurrency(Math.abs(tx.amount))}
          </span>
        )
      },
    },
    {
      key: 'balance_after',
      header: 'Balance After',
      render: (tx) => <span className="text-white">{formatCurrency(tx.balance_after)}</span>,
    },
  ]

  const settlementColumns: Column<Settlement>[] = [
    {
      key: 'id',
      header: 'ID',
      render: (s) => <span className="font-mono text-xs text-white">{s.id.length > 12 ? s.id.slice(0, 12) + '...' : s.id}</span>,
    },
    {
      key: 'amount',
      header: 'Amount',
      render: (s) => <span className="text-white">{formatCurrency(s.amount)}</span>,
    },
    {
      key: 'status',
      header: 'Status',
      render: (s) => <Badge variant={s.status === 'completed' ? 'success' : s.status === 'pending' ? 'warning' : 'neutral'}>{s.status}</Badge>,
    },
    {
      key: 'date',
      header: 'Date',
      render: (s) => new Date(s.date).toLocaleDateString(),
    },
  ]

  const loading = balanceLoading || paymentsLoading

  if (loading) {
    return (
      <div className="mb-8">
        <Skeleton className="mb-6 h-8 w-48" />
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} variant="card" />
          ))}
        </div>
        <div className="mt-8">
          <div className="mb-4 flex gap-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-10 w-32" />
            ))}
          </div>
          <Skeleton className="h-64 w-full rounded-xl" />
        </div>
      </div>
    )
  }

  return (
    <>
      <h1 className="mb-6 text-2xl font-bold text-white">Reports</h1>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {summaryCards.map((card) => (
          <div
            key={card.title}
            className="rounded-xl border border-slate-700 bg-slate-800 p-6"
          >
            <span className="text-sm text-slate-400">{card.title}</span>
            <p className="mt-2 text-2xl font-bold text-white">{card.value}</p>
          </div>
        ))}
      </div>

      <div className="mt-8">
        <div className="mb-4 flex gap-1 rounded-lg border border-slate-700 bg-slate-800 p-1">
          {TABS.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`rounded-md px-4 py-2 text-sm font-medium transition-colors ${
                activeTab === tab.key
                  ? 'bg-primary text-white'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {activeTab === 'transactions' && (
          <DataTable<Payment>
            columns={paymentColumns}
            data={payments || []}
            loading={paymentsLoading}
            emptyMessage="No transactions found"
          />
        )}

        {activeTab === 'balance' && (
          <DataTable<BalanceTransaction>
            columns={txColumns}
            data={txResponse?.data || []}
            loading={txLoading}
            emptyMessage="No balance transactions found"
            pagination={{
              page: txPage,
              pageSize: 50,
              total: txResponse?.total || 0,
              onPageChange: setTxPage,
            }}
          />
        )}

        {activeTab === 'settlements' && (
          <DataTable<Settlement>
            columns={settlementColumns}
            data={settlementResponse?.data || []}
            loading={settlementLoading}
            emptyMessage="No settlements found"
            pagination={{
              page: settlementPage,
              pageSize: 50,
              total: settlementResponse?.total || 0,
              onPageChange: setSettlementPage,
            }}
          />
        )}
      </div>

      <p className="mt-6 text-center text-sm text-slate-500">
        CSV export coming soon
      </p>
    </>
  )
}
