'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Plus, Eye } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Modal } from '@/components/ui/Modal'
import { Input } from '@/components/ui/Input'
import { Button } from '@/components/ui/Button'
import { toast } from '@/components/ui/Toast'
import { api, APIError } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface Customer {
  id: string
  name: string
  email: string
  phone: string
  created_at: string
}

interface CustomersResponse {
  data: Customer[]
}

const createCustomerSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.string().email('Invalid email'),
  phone: z.string().min(1, 'Phone is required'),
})

type CreateCustomerForm = z.infer<typeof createCustomerSchema>

export default function CustomersPage() {
  const queryClient = useQueryClient()
  const [showModal, setShowModal] = useState(false)

  const { data, isLoading } = useQuery({
    queryKey: queryKeys.customers.list(),
    queryFn: () => api.get<CustomersResponse>('/v1/customers').then(r => r.data),
  })

  const createMutation = useMutation({
    mutationFn: (body: CreateCustomerForm) =>
      api.post('/v1/customers', body),
    onSuccess: () => {
      toast.success('Customer created')
      queryClient.invalidateQueries({ queryKey: queryKeys.customers.list() })
      setShowModal(false)
      reset()
    },
    onError: (err: APIError) => {
      toast.error(err.message || 'Failed to create customer')
    },
  })

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateCustomerForm>({
    resolver: zodResolver(createCustomerSchema),
  })

  const onSubmit = (data: CreateCustomerForm) => {
    createMutation.mutate(data)
  }

  const columns: Column<Customer>[] = [
    {
      key: 'name',
      header: 'Name',
      render: (c) => <span className="font-medium text-white">{c.name}</span>,
    },
    { key: 'email', header: 'Email' },
    { key: 'phone', header: 'Phone' },
    {
      key: 'created_at',
      header: 'Created',
      render: (c) => new Date(c.created_at).toLocaleDateString(),
    },
    {
      key: 'actions',
      header: '',
      render: (c) => (
        <Link
          href={`/merchant/customers/${c.id}`}
          className="inline-flex items-center gap-1 rounded px-3 py-1 text-xs font-medium text-primary transition-colors hover:bg-primary/10"
        >
          <Eye className="h-3.5 w-3.5" />
          View
        </Link>
      ),
    },
  ]

  const customers = data?.data || []

  return (
    <>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-white">Customers</h2>
        <button
          onClick={() => setShowModal(true)}
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          <Plus className="h-4 w-4" />
          Add Customer
        </button>
      </div>

      <DataTable<Customer>
        columns={columns}
        data={customers}
        loading={isLoading}
        emptyMessage="Add your first customer"
      />

      <Modal
        open={showModal}
        onClose={() => { setShowModal(false); reset() }}
        title="Add Customer"
        footer={
          <>
            <Button variant="outline" onClick={() => { setShowModal(false); reset() }}>
              Cancel
            </Button>
            <Button
              onClick={handleSubmit(onSubmit)}
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? 'Creating...' : 'Create Customer'}
            </Button>
          </>
        }
      >
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input
            id="name"
            label="Name"
            placeholder="John Doe"
            error={errors.name?.message}
            {...register('name')}
          />
          <Input
            id="email"
            label="Email"
            type="email"
            placeholder="john@example.com"
            error={errors.email?.message}
            {...register('email')}
          />
          <Input
            id="phone"
            label="Phone"
            placeholder="+8801234567890"
            error={errors.phone?.message}
            {...register('phone')}
          />
        </form>
      </Modal>
    </>
  )
}
