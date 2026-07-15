'use client'

import { useState } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, CheckCircle, XCircle, Clock, RefreshCw, DollarSign } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Badge } from '@/components/ui/Badge'
import { Modal } from '@/components/ui/Modal'
import { Skeleton } from '@/components/ui/Skeleton'
import { toast } from '@/components/ui/Toast'
import { api, APIError } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface Transaction {
  id: string
  type: string
  amount: number
  status: string
  processor_response?: string
  created_at: string
}

interface PaymentDetail {
  id: string
  merchant_id: string
  amount: number
  currency: string
  status: string
  payment_method: string
  description: string
  metadata: Record<string, string>
  amount_received: number
  amount_capturable: number
  amount_refunded: number
  customer_id: string
  idempotency_key: string
  failure_reason: string | null
  failure_code: string | null
  created_at: string
  updated_at: string
  transactions: Transaction[]
}

export default function PaymentDetailPage() {
  const params = useParams()
  const id = params.id as string
  const queryClient = useQueryClient()

  const [modalType, setModalType] = useState<'capture' | 'refund' | 'void' | null>(null)
  const [captureAmount, setCaptureAmount] = useState(0)
  const [refundAmount, setRefundAmount] = useState(0)
  const [refundReason, setRefundReason] = useState('')

  const { data: payment, isLoading, error, refetch } = useQuery({
    queryKey: queryKeys.payments.detail(id),
    queryFn: () => api.get<PaymentDetail>(`/payments/${id}`).then(r => r.data),
    enabled: !!id,
  })

  const captureMutation = useMutation({
    mutationFn: (amount: number) =>
      api.post(`/payments/${id}/capture`, { amount_to_capture: amount }),
    onSuccess: () => {
      toast.success('Payment captured')
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.detail(id) })
      setModalType(null)
    },
    onError: (err: APIError) => {
      toast.error(err.message || 'Failed to capture payment')
    },
  })

  const refundMutation = useMutation({
    mutationFn: ({ amount, reason }: { amount: number; reason: string }) =>
      api.post(`/payments/${id}/refund`, { amount, reason }),
    onSuccess: () => {
      toast.success('Refund processed')
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.detail(id) })
      setModalType(null)
    },
    onError: (err: APIError) => {
      toast.error(err.message || 'Failed to process refund')
    },
  })

  const voidMutation = useMutation({
    mutationFn: () => api.post(`/payments/${id}/void`, {}),
    onSuccess: () => {
      toast.success('Payment voided')
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.detail(id) })
      setModalType(null)
    },
    onError: (err: APIError) => {
      toast.error(err.message || 'Failed to void payment')
    },
  })

  const formatCurrency = (amount: number, currency = 'USD') =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

  const statusBadge = (status: string) => {
    const variants: Record<string, 'success' | 'danger' | 'warning' | 'info'> = {
      succeeded: 'success',
      failed: 'danger',
      pending: 'warning',
      processing: 'info',
      authorized: 'info',
      captured: 'success',
      voided: 'danger',
    }
    return <Badge variant={variants[status] || 'neutral'}>{status}</Badge>
  }

  const openModal = (type: 'capture' | 'refund' | 'void') => {
    if (payment) {
      if (type === 'capture') setCaptureAmount(payment.amount_capturable)
      if (type === 'refund')
        setRefundAmount(payment.amount_received - payment.amount_refunded)
      setRefundReason('')
    }
    setModalType(type)
  }

  if (isLoading) {
    return (
      <>
        <Skeleton className="mb-6 h-5 w-32" />
        <div className="mb-6 flex items-center gap-4">
          <Skeleton className="h-6 w-48" />
          <Skeleton variant="card" className="h-8 w-20" />
          <Skeleton className="h-8 w-32" />
        </div>
        <div className="mb-6 flex gap-3">
          <Skeleton variant="card" className="h-10 w-24" />
          <Skeleton variant="card" className="h-10 w-24" />
        </div>
        <div className="mb-6 grid gap-6 lg:grid-cols-2">
          <Skeleton variant="card" className="h-64" />
          <Skeleton variant="card" className="h-64" />
        </div>
        <Skeleton variant="card" className="h-72" />
      </>
    )
  }

  if (error || !payment) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <p className="mb-4 text-slate-400">Payment not found</p>
        <button
          onClick={() => refetch()}
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          <RefreshCw className="h-4 w-4" />
          Retry
        </button>
      </div>
    )
  }

  return (
    <>
      <Link
        href="/merchant/payments"
        className="mb-6 inline-flex items-center gap-2 text-sm text-slate-400 transition-colors hover:text-white"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Payments
      </Link>

      <div className="mb-6 flex flex-wrap items-center gap-4">
        <span className="font-mono text-sm text-white">{payment.id}</span>
        {statusBadge(payment.status)}
        <span className="inline-flex items-center gap-1.5 text-2xl font-bold text-white">
          <DollarSign className="h-6 w-6 text-slate-400" />
          {formatCurrency(payment.amount, payment.currency)}
        </span>
      </div>

      <div className="mb-6 flex flex-wrap gap-3">
        {(payment.status === 'authorized' || payment.status === 'pending') && (
          <button
            onClick={() => openModal('capture')}
            className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
          >
            Capture
          </button>
        )}
        {(payment.status === 'succeeded' ||
          payment.status === 'captured' ||
          payment.status === 'authorized') && (
          <button
            onClick={() => openModal('refund')}
            className="rounded-lg border border-slate-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-slate-700"
          >
            Refund
          </button>
        )}
        {(payment.status === 'authorized' || payment.status === 'pending') && (
          <button
            onClick={() => openModal('void')}
            className="rounded-lg border border-danger/30 px-4 py-2 text-sm font-medium text-danger transition-colors hover:bg-danger/10"
          >
            Void
          </button>
        )}
      </div>

      <div className="mb-6 grid gap-6 lg:grid-cols-2">
        <div className="rounded-xl border border-slate-700 bg-slate-800 p-6">
          <h3 className="mb-4 text-sm font-semibold uppercase tracking-wider text-slate-400">
            Payment Details
          </h3>
          <div className="space-y-3">
            <DetailRow
              label="Original Amount"
              value={formatCurrency(payment.amount, payment.currency)}
            />
            <DetailRow
              label="Amount Received"
              value={formatCurrency(payment.amount_received, payment.currency)}
            />
            <DetailRow
              label="Amount Refunded"
              value={formatCurrency(payment.amount_refunded, payment.currency)}
            />
            <DetailRow
              label="Amount Capturable"
              value={formatCurrency(payment.amount_capturable, payment.currency)}
            />
            <DetailRow label="Currency" value={payment.currency} />
            <DetailRow label="Description" value={payment.description || '-'} />
            <DetailRow
              label="Payment Method"
              value={payment.payment_method?.replace(/_/g, ' ') || '-'}
            />
            <DetailRow
              label="Created"
              value={new Date(payment.created_at).toLocaleString()}
            />
          </div>
        </div>

        <div className="rounded-xl border border-slate-700 bg-slate-800 p-6">
          <h3 className="mb-4 text-sm font-semibold uppercase tracking-wider text-slate-400">
            Customer Info
          </h3>
          <div className="space-y-3">
            <DetailRow
              label="Customer ID"
              value={
                payment.customer_id ? (
                  <Link
                    href={`/merchant/customers/${payment.customer_id}`}
                    className="font-mono text-xs text-primary transition-colors hover:text-primary-300"
                  >
                    {payment.customer_id}
                  </Link>
                ) : (
                  '-'
                )
              }
            />
            {Object.keys(payment.metadata || {}).length > 0
              ? Object.entries(payment.metadata).map(([key, value]) => (
                  <DetailRow key={key} label={key} value={String(value)} />
                ))
              : null}
          </div>
        </div>
      </div>

      <div className="rounded-xl border border-slate-700 bg-slate-800 p-6">
        <h3 className="mb-6 text-sm font-semibold uppercase tracking-wider text-slate-400">
          Transaction Timeline
        </h3>
        <div className="relative">
          {payment.transactions.map((txn, index) => (
            <div key={txn.id} className="relative flex gap-4 pb-8 last:pb-0">
              {index < payment.transactions.length - 1 && (
                <div className="absolute left-[11px] top-6 h-full w-0.5 bg-slate-700" />
              )}
              <div className="relative z-10 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-slate-700">
                {txn.status === 'succeeded' ? (
                  <CheckCircle className="h-5 w-5 text-success" />
                ) : txn.status === 'failed' ? (
                  <XCircle className="h-5 w-5 text-danger" />
                ) : (
                  <Clock className="h-5 w-5 text-warning" />
                )}
              </div>
              <div className="flex-1">
                <div className="flex flex-wrap items-center gap-2">
                  <span className="font-mono text-xs text-white">{txn.id}</span>
                  <span className="inline-flex items-center rounded-full bg-slate-700 px-2 py-0.5 text-xs font-medium capitalize text-slate-300">
                    {txn.type}
                  </span>
                  <Badge
                    variant={
                      txn.status === 'succeeded'
                        ? 'success'
                        : txn.status === 'failed'
                          ? 'danger'
                          : 'warning'
                    }
                  >
                    {txn.status}
                  </Badge>
                  <span className="text-sm font-medium text-white">
                    {formatCurrency(txn.amount, payment.currency)}
                  </span>
                </div>
                <p className="mt-1 text-xs text-slate-400">
                  {new Date(txn.created_at).toLocaleString()}
                </p>
                {txn.processor_response && (
                  <p className="mt-1 text-xs text-slate-500">
                    Processor response: {txn.processor_response}
                  </p>
                )}
              </div>
            </div>
          ))}
          {payment.transactions.length === 0 && (
            <p className="py-8 text-center text-sm text-slate-500">
              No transactions yet
            </p>
          )}
        </div>
      </div>

      <Modal
        open={modalType === 'capture'}
        onClose={() => setModalType(null)}
        title="Capture Payment"
        footer={
          <>
            <button
              onClick={() => setModalType(null)}
              className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 transition-colors hover:bg-slate-700"
            >
              Cancel
            </button>
            <button
              onClick={() => captureMutation.mutate(captureAmount)}
              disabled={captureMutation.isPending || captureAmount <= 0}
              className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {captureMutation.isPending ? 'Capturing...' : 'Capture'}
            </button>
          </>
        }
      >
        <div>
          <label className="mb-1 block text-sm text-slate-400">
            Amount to Capture
          </label>
          <input
            type="number"
            value={captureAmount}
            onChange={(e) => setCaptureAmount(Number(e.target.value))}
            className="w-full rounded-lg border border-slate-600 bg-slate-700 px-3 py-2 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none"
            min={0}
            max={payment.amount_capturable}
          />
          <p className="mt-1 text-xs text-slate-500">
            Maximum capturable:{' '}
            {formatCurrency(payment.amount_capturable, payment.currency)}
          </p>
        </div>
      </Modal>

      <Modal
        open={modalType === 'refund'}
        onClose={() => setModalType(null)}
        title="Process Refund"
        footer={
          <>
            <button
              onClick={() => setModalType(null)}
              className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 transition-colors hover:bg-slate-700"
            >
              Cancel
            </button>
            <button
              onClick={() =>
                refundMutation.mutate({
                  amount: refundAmount,
                  reason: refundReason,
                })
              }
              disabled={refundMutation.isPending || refundAmount <= 0}
              className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {refundMutation.isPending ? 'Refunding...' : 'Refund'}
            </button>
          </>
        }
      >
        <div className="space-y-4">
          <div>
            <label className="mb-1 block text-sm text-slate-400">Amount</label>
            <input
              type="number"
              value={refundAmount}
              onChange={(e) => setRefundAmount(Number(e.target.value))}
              className="w-full rounded-lg border border-slate-600 bg-slate-700 px-3 py-2 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none"
              min={0}
              max={payment.amount_received - payment.amount_refunded}
            />
            <p className="mt-1 text-xs text-slate-500">
              Maximum refundable:{' '}
              {formatCurrency(
                payment.amount_received - payment.amount_refunded,
                payment.currency,
              )}
            </p>
          </div>
          <div>
            <label className="mb-1 block text-sm text-slate-400">Reason</label>
            <textarea
              value={refundReason}
              onChange={(e) => setRefundReason(e.target.value)}
              className="w-full rounded-lg border border-slate-600 bg-slate-700 px-3 py-2 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none"
              rows={3}
              placeholder="Optional reason for refund"
            />
          </div>
        </div>
      </Modal>

      <Modal
        open={modalType === 'void'}
        onClose={() => setModalType(null)}
        title="Void Payment"
        footer={
          <>
            <button
              onClick={() => setModalType(null)}
              className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 transition-colors hover:bg-slate-700"
            >
              Cancel
            </button>
            <button
              onClick={() => voidMutation.mutate()}
              disabled={voidMutation.isPending}
              className="rounded-lg bg-danger px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-danger/80 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {voidMutation.isPending ? 'Voiding...' : 'Void Payment'}
            </button>
          </>
        }
      >
        <p className="text-sm text-slate-300">
          Are you sure you want to void this payment? This action cannot be
          undone.
        </p>
      </Modal>
    </>
  )
}

function DetailRow({
  label,
  value,
}: {
  label: string
  value: React.ReactNode
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-slate-400">{label}</span>
      <span className="text-sm text-white">{value}</span>
    </div>
  )
}
