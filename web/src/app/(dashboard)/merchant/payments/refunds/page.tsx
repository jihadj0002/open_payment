'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, RotateCcw, Loader2 } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Badge } from '@/components/ui/Badge'
import { Modal } from '@/components/ui/Modal'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { toast } from '@/components/ui/Toast'
import { api, APIError } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface Payment {
  id: string
  amount: number
  currency: string
  status: string
  amount_received: number
  amount_refunded: number
  metadata: Record<string, string>
  created_at: string
}

interface PaymentDetail {
  id: string
  amount: number
  currency: string
  status: string
  amount_received: number
  amount_refunded: number
  metadata: Record<string, string>
}

interface PaymentsListResponse {
  data: Payment[]
  total: number
  limit: number
  offset: number
}

const formatCurrency = (amount: number, currency = 'USD') =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

const statusBadge = (status: string) => {
  const variants: Record<string, 'success' | 'danger' | 'warning' | 'info' | 'neutral'> = {
    refunded: 'success',
    partially_refunded: 'warning',
  }
  return <Badge variant={variants[status] || 'neutral'}>{status.replace(/_/g, ' ')}</Badge>
}

export default function RefundsPage() {
  const queryClient = useQueryClient()

  const [modalOpen, setModalOpen] = useState(false)
  const [paymentId, setPaymentId] = useState('')
  const [amount, setAmount] = useState(0)
  const [reason, setReason] = useState('')
  const [fetchedPayment, setFetchedPayment] = useState<PaymentDetail | null>(null)
  const [fetching, setFetching] = useState(false)

  const { data: paymentsResponse, isLoading, error } = useQuery({
    queryKey: queryKeys.payments.list({ limit: '100' }),
    queryFn: () => api.get<PaymentsListResponse>(`/payments?limit=100`).then(r => r.data),
  })

  const refundedPayments = (paymentsResponse?.data || []).filter(
    (p) => p.status === 'refunded' || p.status === 'partially_refunded'
  )

  const refundMutation = useMutation({
    mutationFn: (body: { amount: number; reason: string }) =>
      api.post(`/payments/${paymentId}/refund`, body),
    onSuccess: () => {
      toast.success('Refund processed')
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.all })
      handleCloseModal()
    },
    onError: (err: APIError) => {
      toast.error(err.message || 'Failed to process refund')
    },
  })

  const handleCloseModal = () => {
    setModalOpen(false)
    setPaymentId('')
    setAmount(0)
    setReason('')
    setFetchedPayment(null)
  }

  const handleFetchPayment = async () => {
    if (!paymentId.trim()) return
    setFetching(true)
    setFetchedPayment(null)
    try {
      const payment = await api.get<PaymentDetail>(`/payments/${paymentId}`).then(r => r.data)
      setFetchedPayment(payment)
      setAmount(payment.amount_received - payment.amount_refunded)
    } catch {
      toast.error('Payment not found')
    } finally {
      setFetching(false)
    }
  }

  const columns: Column<Payment>[] = [
    {
      key: 'id',
      header: 'Payment ID',
      render: (p) => (
        <span className="font-mono text-xs text-white">
          {p.id.length > 12 ? p.id.slice(0, 12) + '...' : p.id}
        </span>
      ),
    },
    {
      key: 'amount',
      header: 'Amount',
      render: (p) => (
        <span className="text-white">{formatCurrency(p.amount, p.currency)}</span>
      ),
    },
    { key: 'currency', header: 'Currency' },
    {
      key: 'status',
      header: 'Status',
      render: (p) => statusBadge(p.status),
    },
    {
      key: 'reason',
      header: 'Reason',
      render: (p) => (
        <span className="text-slate-300">
          {p.metadata?.refund_reason || '-'}
        </span>
      ),
    },
    {
      key: 'created_at',
      header: 'Date',
      render: (p) => new Date(p.created_at).toLocaleDateString(),
    },
    {
      key: 'actions',
      header: 'Actions',
      render: (p) => (
        <Link
          href={`/merchant/payments/${p.id}`}
          className="text-sm text-primary transition-colors hover:text-primary-300"
        >
          View Payment
        </Link>
      ),
    },
  ]

  return (
    <>
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link
            href="/merchant/payments"
            className="rounded-lg p-1.5 text-slate-400 transition-colors hover:bg-slate-800 hover:text-white"
          >
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <h2 className="text-lg font-semibold text-white">Refunds</h2>
        </div>
        <Button onClick={() => setModalOpen(true)}>
          <RotateCcw className="mr-2 h-4 w-4" />
          Initiate Refund
        </Button>
      </div>

      {error && (
        <div className="mb-6 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
          Failed to load refunds
        </div>
      )}

      <DataTable<Payment>
        columns={columns}
        data={refundedPayments}
        loading={isLoading}
        emptyMessage="No refunds yet"
      />

      <Modal
        open={modalOpen}
        onClose={handleCloseModal}
        title="Initiate Refund"
        footer={
          <>
            <Button variant="secondary" onClick={handleCloseModal}>
              Cancel
            </Button>
            <Button
              onClick={() => refundMutation.mutate({ amount, reason })}
              disabled={refundMutation.isPending || amount <= 0 || !fetchedPayment}
            >
              {refundMutation.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Refunding...
                </>
              ) : (
                'Confirm Refund'
              )}
            </Button>
          </>
        }
      >
        <div className="space-y-4">
          <div className="flex items-end gap-2">
            <div className="flex-1">
              <Input
                label="Payment ID"
                value={paymentId}
                onChange={(e) => setPaymentId(e.target.value)}
                placeholder="Enter payment ID"
              />
            </div>
            <Button
              variant="secondary"
              onClick={handleFetchPayment}
              disabled={fetching || !paymentId.trim()}
            >
              {fetching ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                'Fetch'
              )}
            </Button>
          </div>

          {fetchedPayment && (
            <>
              <div>
                <label className="mb-1 block text-sm text-slate-400">Amount</label>
                <input
                  type="number"
                  value={amount}
                  onChange={(e) => setAmount(Number(e.target.value))}
                  className="w-full rounded-lg border border-slate-600 bg-slate-700 px-3 py-2 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none"
                  min={0}
                  max={fetchedPayment.amount_received - fetchedPayment.amount_refunded}
                />
                <p className="mt-1 text-xs text-slate-500">
                  Maximum refundable:{' '}
                  {formatCurrency(
                    fetchedPayment.amount_received - fetchedPayment.amount_refunded,
                    fetchedPayment.currency,
                  )}
                </p>
              </div>
              <div>
                <label className="mb-1 block text-sm text-slate-400">
                  Reason <span className="text-slate-500">(optional)</span>
                </label>
                <textarea
                  value={reason}
                  onChange={(e) => setReason(e.target.value)}
                  className="w-full rounded-lg border border-slate-600 bg-slate-700 px-3 py-2 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none resize-none"
                  rows={3}
                  placeholder="Reason for refund"
                />
              </div>
            </>
          )}
        </div>
      </Modal>
    </>
  )
}
