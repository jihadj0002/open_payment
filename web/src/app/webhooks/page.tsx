'use client'

import { useState } from 'react'
import { Plus, Copy, Check, X } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Badge } from '@/components/ui/Badge'
import { toast } from '@/components/ui/Toast'
import { api } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface Webhook {
  id: string
  url: string
  event: string
  status: string
  created_at: string
}

interface CreatedWebhook {
  id: string
  url: string
  event: string
  signing_secret: string
}

const EVENT_OPTIONS = [
  { value: 'payment.success', label: 'Payment Succeeded' },
  { value: 'payment.failed', label: 'Payment Failed' },
  { value: 'refund.completed', label: 'Refund Completed' },
  { value: 'chargeback.created', label: 'Chargeback Created' },
]

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <button
      onClick={handleCopy}
      className="flex items-center gap-1.5 rounded bg-slate-700 px-3 py-1.5 text-xs text-slate-300 transition-colors hover:bg-slate-600"
    >
      {copied ? <Check className="h-3.5 w-3.5 text-success" /> : <Copy className="h-3.5 w-3.5" />}
      {copied ? 'Copied!' : 'Copy'}
    </button>
  )
}

export default function WebhooksPage() {
  const queryClient = useQueryClient()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showSecretModal, setShowSecretModal] = useState<CreatedWebhook | null>(null)
  const [newUrl, setNewUrl] = useState('')
  const [newEvent, setNewEvent] = useState('payment.success')
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null)

  const { data: webhooks = [], isLoading } = useQuery({
    queryKey: queryKeys.webhooks.all,
    queryFn: () => api.get<Webhook[]>('/webhook_endpoints').then(r => r.data),
  })

  const createMutation = useMutation({
    mutationFn: (body: { url: string; event: string }) =>
      api.post<CreatedWebhook>('/webhook_endpoints', body),
    onSuccess: (res) => {
      setShowCreateModal(false)
      setNewUrl('')
      setNewEvent('payment.success')
      setShowSecretModal(res.data)
      queryClient.invalidateQueries({ queryKey: queryKeys.webhooks.all })
    },
    onError: () => {
      toast.error('Failed to create webhook')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.delete(`/webhook_endpoints/${id}`),
    onSuccess: () => {
      setDeleteConfirm(null)
      queryClient.invalidateQueries({ queryKey: queryKeys.webhooks.all })
    },
    onError: () => {
      toast.error('Failed to delete webhook')
    },
  })

  const statusBadgeVariant = (status: string) => {
    if (status === 'active') return 'success'
    if (status === 'disabled') return 'neutral'
    return 'neutral'
  }

  const columns: Column<Webhook>[] = [
    {
      key: 'url',
      header: 'URL',
      render: (w) => (
        <span className="max-w-[200px] truncate font-mono text-xs text-white">{w.url}</span>
      ),
    },
    {
      key: 'event',
      header: 'Event',
      render: (w) => (
        <code className="rounded bg-slate-700 px-2 py-0.5 font-mono text-xs text-slate-300">
          {w.event}
        </code>
      ),
    },
    {
      key: 'status',
      header: 'Status',
      render: (w) => <Badge variant={statusBadgeVariant(w.status)}>{w.status}</Badge>,
    },
    {
      key: 'created_at',
      header: 'Created',
      render: (w) => new Date(w.created_at).toLocaleDateString(),
    },
    {
      key: 'actions',
      header: '',
      render: (w) => (
        <button
          onClick={() => setDeleteConfirm(w.id)}
          className="rounded px-3 py-1 text-xs font-medium text-danger transition-colors hover:bg-danger/10"
        >
          Delete
        </button>
      ),
    },
  ]

  return (
    <DashboardLayout>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-white">Webhooks</h2>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          <Plus className="h-4 w-4" />
          Add Webhook
        </button>
      </div>

      <DataTable<Webhook>
        columns={columns}
        data={webhooks}
        loading={isLoading}
        emptyMessage="Configure your first webhook endpoint"
      />

      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-md rounded-xl border border-slate-700 bg-slate-800 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-lg font-semibold text-white">Add Webhook</h3>
              <button onClick={() => setShowCreateModal(false)} className="text-slate-400 hover:text-white">
                <X className="h-5 w-5" />
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="mb-1.5 block text-sm text-slate-400">Endpoint URL</label>
                <input
                  type="url"
                  placeholder="https://example.com/webhook"
                  value={newUrl}
                  onChange={(e) => setNewUrl(e.target.value)}
                  className="w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="mb-1.5 block text-sm text-slate-400">Event Type</label>
                <select
                  value={newEvent}
                  onChange={(e) => setNewEvent(e.target.value)}
                  className="w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-sm text-white focus:border-primary focus:outline-none"
                >
                  {EVENT_OPTIONS.map((opt) => (
                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                  ))}
                </select>
              </div>
            </div>
            <div className="mt-6 flex justify-end gap-3">
              <button
                onClick={() => { setShowCreateModal(false); setNewUrl(''); setNewEvent('payment.success') }}
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 transition-colors hover:bg-slate-700"
              >
                Cancel
              </button>
              <button
                onClick={() => createMutation.mutate({ url: newUrl.trim(), event: newEvent })}
                disabled={!newUrl.trim() || createMutation.isPending}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700 disabled:opacity-50"
              >
                {createMutation.isPending ? 'Creating...' : 'Create'}
              </button>
            </div>
          </div>
        </div>
      )}

      {showSecretModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-lg rounded-xl border border-slate-700 bg-slate-800 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-lg font-semibold text-white">Webhook Created</h3>
              <button onClick={() => setShowSecretModal(null)} className="text-slate-400 hover:text-white">
                <X className="h-5 w-5" />
              </button>
            </div>
            <p className="mb-4 text-sm text-slate-400">
              Use this signing secret to verify webhook payloads.
            </p>
            <div className="flex items-center gap-2 rounded-lg bg-slate-900 p-3">
              <code className="flex-1 break-all font-mono text-xs text-green-400">
                {showSecretModal.signing_secret}
              </code>
              <CopyButton text={showSecretModal.signing_secret} />
            </div>
            <div className="mt-4 flex justify-end">
              <button
                onClick={() => setShowSecretModal(null)}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
              >
                Done
              </button>
            </div>
          </div>
        </div>
      )}

      {deleteConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-sm rounded-xl border border-slate-700 bg-slate-800 p-6">
            <h3 className="mb-2 text-lg font-semibold text-white">Delete Webhook</h3>
            <p className="mb-4 text-sm text-slate-400">
              Are you sure you want to delete this webhook endpoint?
            </p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 transition-colors hover:bg-slate-700"
              >
                Cancel
              </button>
              <button
                onClick={() => deleteMutation.mutate(deleteConfirm)}
                className="rounded-lg bg-danger px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-danger-700"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
