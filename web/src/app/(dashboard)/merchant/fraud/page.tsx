'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Badge } from '@/components/ui/Badge'
import { Button } from '@/components/ui/Button'
import { Modal } from '@/components/ui/Modal'
import { Skeleton } from '@/components/ui/Skeleton'
import { toast } from '@/components/ui/Toast'
import { Plus, ToggleLeft, ToggleRight, Edit, Trash2 } from 'lucide-react'

interface FraudRule {
  id: string
  name: string
  type: string
  threshold: number
  action: string
  enabled: boolean
  created_at: string
}

const actionBadgeVariant = (action: string) => {
  switch (action) {
    case 'block': return 'danger'
    case 'review': return 'warning'
    default: return 'success'
  }
}

export default function FraudRulesPage() {
  const queryClient = useQueryClient()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingRule, setEditingRule] = useState<FraudRule | null>(null)
  const [notification, setNotification] = useState<{ type: 'success' | 'error', message: string } | null>(null)

  const { data: rules, isLoading } = useQuery({
    queryKey: queryKeys.fraud.rules,
    queryFn: () => api.get<FraudRule[]>('/merchants/fraud/rules').then(r => r.data),
  })

  const toggleMutation = useMutation({
    mutationFn: ({ id, enabled }: { id: string, enabled: boolean }) =>
      api.patch(`/merchants/fraud/rules/${id}`, { enabled }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.fraud.rules })
      toast.success('Rule updated')
    },
    onError: () => {
      toast.error('Failed to update rule')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.delete(`/merchants/fraud/rules/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.fraud.rules })
      toast.success('Rule deleted')
    },
    onError: () => {
      toast.error('Failed to delete rule')
    },
  })

  const columns: Column<FraudRule>[] = [
    {
      key: 'name',
      header: 'Name',
      render: (rule: FraudRule) => <span>{rule.name}</span>,
    },
    {
      key: 'type',
      header: 'Type',
      render: (rule: FraudRule) => <span>{rule.type}</span>,
    },
    {
      key: 'threshold',
      header: 'Threshold',
      render: (rule: FraudRule) => <span>{rule.threshold}</span>,
    },
    {
      key: 'action',
      header: 'Action',
      render: (rule: FraudRule) => (
        <Badge variant={actionBadgeVariant(rule.action)}>{rule.action}</Badge>
      ),
    },
    {
      key: 'enabled',
      header: 'Enabled',
      render: (rule: FraudRule) => (
        <button
          onClick={() => toggleMutation.mutate({ id: rule.id, enabled: !rule.enabled })}
          className="text-gray-500 hover:text-blue-600"
        >
          {rule.enabled ? <ToggleRight className="w-5 h-5 text-green-500" /> : <ToggleLeft className="w-5 h-5" />}
        </button>
      ),
    },
    {
      key: 'actions',
      header: '',
      render: (rule: FraudRule) => (
        <div className="flex gap-2">
          <button onClick={() => setEditingRule(rule)} className="text-gray-400 hover:text-blue-600">
            <Edit className="w-4 h-4" />
          </button>
          <button onClick={() => deleteMutation.mutate(rule.id)} className="text-gray-400 hover:text-red-600">
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-6">
      {notification && (
        <div className={`rounded-lg border px-4 py-3 text-sm ${
          notification.type === 'success'
            ? 'border-success/20 bg-success/10 text-success'
            : 'border-danger/20 bg-danger/10 text-danger'
        }`}>
          {notification.message}
          <button onClick={() => setNotification(null)} className="float-right ml-2 text-sm">
            ×
          </button>
        </div>
      )}

      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Fraud Detection Rules</h1>
          <p className="text-gray-500 mt-1">Configure automated fraud prevention rules</p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus className="w-4 h-4 mr-2" />
          Add Rule
        </Button>
      </div>

      {isLoading ? (
        <Skeleton variant="card" />
      ) : (
        <DataTable
          columns={columns}
          data={rules || []}
          emptyMessage="No fraud rules configured. Create your first rule to get started."
        />
      )}

      <Modal
        open={showCreateModal || !!editingRule}
        onClose={() => { setShowCreateModal(false); setEditingRule(null) }}
        title={editingRule ? 'Edit Rule' : 'Create Fraud Rule'}
      >
        <div className="space-y-4">
          <p className="text-sm text-gray-500">
            {editingRule
              ? `Editing rule: ${editingRule.name}`
              : 'Configure a new fraud detection rule'}
          </p>
          <div className="flex justify-end gap-3 pt-4">
            <Button variant="outline" onClick={() => { setShowCreateModal(false); setEditingRule(null) }}>
              Cancel
            </Button>
            <Button>Save Rule</Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
