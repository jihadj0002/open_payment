'use client'

import { useEffect, useState } from 'react'
import { Plus, Copy, Check, AlertTriangle, X } from 'lucide-react'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { api } from '@/lib/api'

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
      className="flex items-center gap-1.5 rounded bg-slate-700 px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-600 transition-colors"
    >
      {copied ? <Check className="h-3.5 w-3.5 text-success" /> : <Copy className="h-3.5 w-3.5" />}
      {copied ? 'Copied!' : 'Copy'}
    </button>
  )
}

export default function ApiKeysPage() {
  const [keys, setKeys] = useState<ApiKey[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showKeyModal, setShowKeyModal] = useState<CreatedKey | null>(null)
  const [newKeyName, setNewKeyName] = useState('')
  const [creating, setCreating] = useState(false)
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null)

  const fetchKeys = () => {
    setLoading(true)
    api.get<ApiKey[]>('/merchants/api_keys')
      .then((res) => setKeys(res.data))
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchKeys() }, [])

  const handleCreate = async () => {
    if (!newKeyName.trim()) return
    setCreating(true)
    try {
      const res = await api.post<CreatedKey>('/merchants/api_keys', { name: newKeyName.trim() })
      setShowCreateModal(false)
      setNewKeyName('')
      setShowKeyModal(res.data)
      fetchKeys()
    } catch {
    } finally {
      setCreating(false)
    }
  }

  const handleRevoke = async (id: string) => {
    try {
      await api.delete(`/merchants/api_keys/${id}`)
      setDeleteConfirm(null)
      fetchKeys()
    } catch {}
  }

  const statusBadge = (status: string) => {
    if (status === 'active') return 'bg-success/10 text-success'
    if (status === 'revoked') return 'bg-danger/10 text-danger'
    return 'bg-slate-700 text-slate-300'
  }

  return (
    <DashboardLayout>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-white">API Keys</h2>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Create API Key
        </button>
      </div>

      <div className="overflow-hidden rounded-xl border border-slate-700">
        <table className="w-full text-left text-sm">
          <thead className="bg-slate-800">
            <tr>
              <th className="px-6 py-3 font-medium text-slate-400">Name</th>
              <th className="px-6 py-3 font-medium text-slate-400">Prefix</th>
              <th className="px-6 py-3 font-medium text-slate-400">Status</th>
              <th className="px-6 py-3 font-medium text-slate-400">Last Used</th>
              <th className="px-6 py-3 font-medium text-slate-400">Created</th>
              <th className="px-6 py-3 font-medium text-slate-400"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {loading ? (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-slate-500">
                  <div className="flex items-center justify-center gap-2">
                    <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                    Loading...
                  </div>
                </td>
              </tr>
            ) : keys.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-slate-500">
                  No API keys found
                </td>
              </tr>
            ) : (
              keys.map((key) => (
                <tr key={key.id} className="hover:bg-slate-800/50">
                  <td className="px-6 py-4 text-white font-medium">{key.name}</td>
                  <td className="px-6 py-4">
                    <code className="rounded bg-slate-700 px-2 py-0.5 font-mono text-xs text-slate-300">
                      {key.prefix}
                    </code>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusBadge(key.status)}`}>
                      {key.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-400">
                    {key.last_used_at ? new Date(key.last_used_at).toLocaleDateString() : 'Never'}
                  </td>
                  <td className="px-6 py-4 text-slate-400">
                    {new Date(key.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4">
                    {key.status === 'active' && (
                      <button
                        onClick={() => setDeleteConfirm(key.id)}
                        className="rounded px-3 py-1 text-xs font-medium text-danger hover:bg-danger/10 transition-colors"
                      >
                        Revoke
                      </button>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-md rounded-xl border border-slate-700 bg-slate-800 p-6">
            <div className="flex items-center justify-between mb-4">
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
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreate}
                disabled={!newKeyName.trim() || creating}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 disabled:opacity-50 transition-colors"
              >
                {creating ? 'Creating...' : 'Create'}
              </button>
            </div>
          </div>
        </div>
      )}

      {showKeyModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-lg rounded-xl border border-slate-700 bg-slate-800 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-white">API Key Created</h3>
              <button onClick={() => setShowKeyModal(null)} className="text-slate-400 hover:text-white">
                <X className="h-5 w-5" />
              </button>
            </div>
            <div className="flex items-start gap-3 rounded-lg bg-warning/10 p-4 mb-4">
              <AlertTriangle className="h-5 w-5 text-warning shrink-0 mt-0.5" />
              <p className="text-sm text-warning">
                Copy this key now. You won&apos;t be able to see it again.
              </p>
            </div>
            <div className="flex items-center gap-2 rounded-lg bg-slate-900 p-3">
              <code className="flex-1 font-mono text-xs text-green-400 break-all">
                {showKeyModal.full_key}
              </code>
              <CopyButton text={showKeyModal.full_key} />
            </div>
            <div className="mt-4 flex justify-end">
              <button
                onClick={() => setShowKeyModal(null)}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
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
            <h3 className="text-lg font-semibold text-white mb-2">Revoke API Key</h3>
            <p className="text-sm text-slate-400 mb-4">
              Are you sure you want to revoke this API key? This action cannot be undone.
            </p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => handleRevoke(deleteConfirm)}
                className="rounded-lg bg-danger px-4 py-2 text-sm font-medium text-white hover:bg-danger-700 transition-colors"
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
