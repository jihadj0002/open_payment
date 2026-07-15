'use client'

import { useState } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, RefreshCw, DollarSign, CreditCard } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Card } from '@/components/ui/Card'
import { Skeleton } from '@/components/ui/Skeleton'
import { Modal } from '@/components/ui/Modal'
import { Input } from '@/components/ui/Input'
import { Button } from '@/components/ui/Button'
import { toast } from '@/components/ui/Toast'
import { api, APIError } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface Customer {
  id: string
  merchant_id: string
  name: string
  email: string
  phone: string
  metadata: Record<string, string>
  created_at: string
  updated_at: string
}

interface CustomerResponse {
  data: Customer
}

interface Payment {
  id: string
  amount: number
  currency: string
  status: string
  customer_id: string
  created_at: string
}

interface PaymentsResponse {
  payments: Payment[]
  total: number
  page: number
  limit: number
}

const editCustomerSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.string().email('Invalid email'),
  phone: z.string().min(1, 'Phone is required'),
})

type EditCustomerForm = z.infer<typeof editCustomerSchema>

export default function CustomerDetailPage() {
  const params = useParams()
  const id = params.id as string
  const queryClient = useQueryClient()
  const [showEditModal, setShowEditModal] = useState(false)

  const { data: customerData, isLoading, error, refetch } = useQuery({
    queryKey: queryKeys.customers.detail(id),
    queryFn: () => api.get<CustomerResponse>(`/v1/customers/${id}`).then(r => r.data),
    enabled: !!id,
  })

  const { data: paymentsData } = useQuery({
    queryKey: [...queryKeys.payments.list({ limit: '5' }), id],
    queryFn: () =>
      api.get<PaymentsResponse>('/payments?limit=5').then(r => r.data),
    enabled: !!id,
  })

  const customer = customerData?.data
  const recentPayments = (paymentsData?.payments || []).filter(
    (p) => p.customer_id === id
  )

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<EditCustomerForm>({
    resolver: zodResolver(editCustomerSchema),
  })

  const openEditModal = () => {
    if (customer) {
      reset({
        name: customer.name,
        email: customer.email,
        phone: customer.phone,
      })
      setShowEditModal(true)
    }
  }

  const editMutation = useMutation({
    mutationFn: (body: EditCustomerForm) =>
      api.patch(`/v1/customers/${id}`, body),
    onSuccess: () => {
      toast.success('Customer updated')
      queryClient.invalidateQueries({ queryKey: queryKeys.customers.detail(id) })
      queryClient.invalidateQueries({ queryKey: queryKeys.customers.list() })
      setShowEditModal(false)
    },
    onError: (err: APIError) => {
      toast.error(err.message || 'Failed to update customer')
    },
  })

  const onSubmit = (data: EditCustomerForm) => {
    editMutation.mutate(data)
  }

  const formatCurrency = (amount: number, currency = 'USD') =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)

  const statusBadge = (status: string) => {
    const colors: Record<string, string> = {
      succeeded: 'bg-success/10 text-success border-success/30',
      failed: 'bg-danger/10 text-danger border-danger/30',
      pending: 'bg-warning/10 text-warning border-warning/30',
      processing: 'bg-blue-500/10 text-blue-400 border-blue-500/30',
    }
    return (
      <span
        className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium capitalize ${colors[status] || 'bg-slate-700 text-slate-300 border-slate-600'}`}
      >
        {status}
      </span>
    )
  }

  if (isLoading) {
    return (
      <>
        <Skeleton className="mb-6 h-4 w-40" />
        <div className="mb-6 flex items-center gap-4">
          <Skeleton className="h-7 w-56" />
          <Skeleton className="h-9 w-24" />
        </div>
        <div className="mb-6 grid gap-6 lg:grid-cols-2">
          <Skeleton variant="card" className="h-56" />
          <Skeleton variant="card" className="h-56" />
        </div>
        <Skeleton variant="card" className="h-48" />
        <Skeleton variant="card" className="mt-6 h-48" />
      </>
    )
  }

  if (error || !customer) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <p className="mb-4 text-slate-400">Customer not found</p>
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
        href="/merchant/customers"
        className="mb-6 inline-flex items-center gap-2 text-sm text-slate-400 transition-colors hover:text-white"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Customers
      </Link>

      <div className="mb-6 flex flex-wrap items-center justify-between gap-4">
        <h2 className="text-xl font-semibold text-white">{customer.name}</h2>
        <Button onClick={openEditModal}>Edit Customer</Button>
      </div>

      <div className="mb-6 grid gap-6 lg:grid-cols-2">
        <Card header={<span className="text-sm font-semibold uppercase tracking-wider text-slate-400">Profile</span>}>
          <div className="space-y-3">
            <DetailRow label="Name" value={customer.name} />
            <DetailRow label="Email" value={customer.email} />
            <DetailRow label="Phone" value={customer.phone} />
            <DetailRow
              label="Created"
              value={new Date(customer.created_at).toLocaleDateString()}
            />
          </div>
        </Card>

        <Card header={<span className="text-sm font-semibold uppercase tracking-wider text-slate-400">Stats</span>}>
          <div className="space-y-3">
            <DetailRow
              label="Total Spent"
              value={
                <span className="inline-flex items-center gap-1">
                  <DollarSign className="h-4 w-4 text-slate-400" />
                  —
                </span>
              }
            />
            <DetailRow
              label="Payment Count"
              value={
                <span className="inline-flex items-center gap-1">
                  <CreditCard className="h-4 w-4 text-slate-400" />
                  —
                </span>
              }
            />
          </div>
        </Card>
      </div>

      <Card className="mb-6" header={<span className="text-sm font-semibold uppercase tracking-wider text-slate-400">Saved Payment Methods</span>}>
        <p className="py-8 text-center text-sm text-slate-500">
          No payment methods saved
        </p>
      </Card>

      <Card header={<span className="text-sm font-semibold uppercase tracking-wider text-slate-400">Recent Payments</span>}>
        {recentPayments.length === 0 ? (
          <p className="py-8 text-center text-sm text-slate-500">
            No recent payments
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-slate-700">
              <thead>
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-400">ID</th>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-400">Amount</th>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-400">Status</th>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-400">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-700">
                {recentPayments.map((p) => (
                  <tr key={p.id} className="transition-colors hover:bg-slate-800">
                    <td className="whitespace-nowrap px-4 py-3">
                      <Link
                        href={`/merchant/payments/${p.id}`}
                        className="font-mono text-xs text-primary transition-colors hover:text-primary-300"
                      >
                        {p.id.length > 12 ? p.id.slice(0, 12) + '...' : p.id}
                      </Link>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-sm text-white">
                      {formatCurrency(p.amount, p.currency)}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3">{statusBadge(p.status)}</td>
                    <td className="whitespace-nowrap px-4 py-3 text-sm text-slate-300">
                      {new Date(p.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      <Modal
        open={showEditModal}
        onClose={() => setShowEditModal(false)}
        title="Edit Customer"
        footer={
          <>
            <Button variant="outline" onClick={() => setShowEditModal(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleSubmit(onSubmit)}
              disabled={editMutation.isPending}
            >
              {editMutation.isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </>
        }
      >
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input
            id="name"
            label="Name"
            error={errors.name?.message}
            {...register('name')}
          />
          <Input
            id="email"
            label="Email"
            type="email"
            error={errors.email?.message}
            {...register('email')}
          />
          <Input
            id="phone"
            label="Phone"
            error={errors.phone?.message}
            {...register('phone')}
          />
        </form>
      </Modal>
    </>
  )
}

function DetailRow({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-slate-400">{label}</span>
      <span className="text-sm text-white">{value}</span>
    </div>
  )
}
