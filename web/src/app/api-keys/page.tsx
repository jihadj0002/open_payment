'use client'

import { useState } from 'react'
import { Plus, Copy, Check, AlertTriangle, X } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { DataTable, Column } from '@/components/ui/DataTable'
import { Badge } from '@/components/ui/Badge'
import { toast } from '@/components/ui/Toast'
import { api } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface ApiKey {
  id: string
  name: string
  prefix: string
  status: 'active' | 'revoked'
  last_used_at: string | null
  created_at: string
}

interface CreatedKey {
  id: string
  name: string
  full_key: string
}

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

export default function ApiKeysPage() {
  const queryClient = useQueryClient()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showKeyModal, setShowKeyModal] = useState<CreatedKey | null>(null)
  const [newKeyName, setNewKeyName] = useState('')
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null)

  const { data: keys = [], isLoading } = useQuery({
    queryKey: queryKeys.apiKeys.all,
    queryFn: () => api.get<ApiKey[]>('/merchants/api_keys').then(r => r.data),
  })

  const createMutation = useMutation({
    mutationFn: (name: string) => api.post<CreatedKey>('/merchants/api_keys', { name }),
    onSuccess: (res) => {
      setShowCreateModal(false)
      setNewKeyName('')
      setShowKeyModal(res.data)
      queryClient.invalidateQueries({ queryKey: queryKeys.apiKeys.all })
    },
    onError: () => {
      toast.error('Failed to create API key')
    },
  })

  const revokeMutation = useMutation({
    mutationFn: (id: string) => api.delete(`/merchants/api_keys/${id}`),
    onSuccess: () => {
      setDeleteConfirm(null)
      queryClient.invalidateQueries({ queryKey: queryKeys.apiKeys.all })
    },
    onError: () => {
      toast.error('Failed to revoke API key')
    },
  })

  const columns: Column<ApiKey>[] = [
    {
      key: 'name',
      header: 'Name',
      render: (k) => <span className="font-medium text-white">{k.name}</span>,
    },
    {
      key: 'prefix',
      header: 'Prefix',
      render: (k) => (
        <code className="rounded bg-slate-700 px-2 py-0.5 font-mono text-xs text-slate-300">
          {k.prefix}
        </code>
      ),
    },
    {
      key: 'status',
      header: 'Status',
      render: (k) => (
        <Badge variant={k.status === 'active' ? 'success' : 'danger'}>{k.status}</Badge>
      ),
    },
    {
      key: 'last_used_at',
      header: 'Last Used',
      render: (k) => (k.last_used_at ? new Date(k.last_used_at).toLocaleDateString() : 'Never'),
    },
    {
      key: 'created_at',
      header: 'Created',
      render: (k) => new Date(k.created_at).toLocaleDateString(),
    },
    {
      key: 'actions',
      header: '',
      render: (k) =>
        k.status === 'active' ? (
          <button
            onClick={() => setDeleteConfirm(k.id)}
            className="rounded px-3 py-1 text-xs font-medium text-danger transition-colors hover:bg-danger/10"
          >
            Revoke
          </button>
        ) : null,
    },
  ]

  return (
    <DashboardLayout>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-white">API Keys</h2>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          <Plus className="h-4 w-4" />
          Create API Key
        </button>
      </div>

      <DataTable<ApiKey>
        columns={columns}
        data={keys}
        loading={isLoading}
        emptyMessage="Generate your first API key"
      />

      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-md rounded-xl border border-slate-700 bg-slate-800 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-lg font-semibold text-white">Create API Key</h3>
              <button onClick={() => setShowCreateModal(false)} className="text-slate-400 hover:text-white">
                <X className="h-5 w-5" />
              </button>
            </div>
            <input
              type="text"
              placeholder="Key name (e.g. Production, Staging)"
              value={newKeyName}
              onChange={(e) => setNewKeyName(e.target.value)}
              className="w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none"
            />
            <div className="mt-4 flex justify-end gap-3">
              <button
                onClick={() => { setShowCreateModal(false); setNewKeyName('') }}
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 transition-colors hover:bg-slate-700"
              >
                Cancel
              </button>
              <button
                onClick={() => createMutation.mutate(newKeyName.trim())}
                disabled={!newKeyName.trim() || createMutation.isPending}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700 disabled:opacity-50"
              >
                {createMutation.isPending ? 'Creating...' : 'Create'}
              </button>
            </div>
          </div>
        </div>
      )}

      {showKeyModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-lg rounded-xl border border-slate-700 bg-slate-800 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-lg font-semibold text-white">API Key Created</h3>
              <button onClick={() => setShowKeyModal(null)} className="text-slate-400 hover:text-white">
                <X className="h-5 w-5" />
              </button>
            </div>
            <div className="mb-4 flex items-start gap-3 rounded-lg bg-warning/10 p-4">
              <AlertTriangle className="mt-0.5 h-5 w-5 shrink-0 text-warning" />
              <p className="text-sm text-warning">
                Copy this key now. You won&apos;t be able to see it again.
              </p>
            </div>
            <div className="flex items-center gap-2 rounded-lg bg-slate-900 p-3">
              <code className="flex-1 break-all font-mono text-xs text-green-400">
                {showKeyModal.full_key}
              </code>
              <CopyButton text={showKeyModal.full_key} />
            </div>
            <div className="mt-4 flex justify-end">
              <button
                onClick={() => setShowKeyModal(null)}
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
            <h3 className="mb-2 text-lg font-semibold text-white">Revoke API Key</h3>
            <p className="mb-4 text-sm text-slate-400">
              Are you sure you want to revoke this API key? This action cannot be undone.
            </p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 transition-colors hover:bg-slate-700"
              >
                Cancel
              </button>
              <button
                onClick={() => revokeMutation.mutate(deleteConfirm)}
                className="rounded-lg bg-danger px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-danger-700"
              >
                Revoke
              </button>
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
