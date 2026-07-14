'use client'

import { useEffect, useState } from 'react'
import { Plus, Copy, Check, X } from 'lucide-react'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { api } from '@/lib/api'

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
      className="flex items-center gap-1.5 rounded bg-slate-700 px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-600 transition-colors"
    >
      {copied ? <Check className="h-3.5 w-3.5 text-success" /> : <Copy className="h-3.5 w-3.5" />}
      {copied ? 'Copied!' : 'Copy'}
    </button>
  )
}

export default function WebhooksPage() {
  const [webhooks, setWebhooks] = useState<Webhook[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showSecretModal, setShowSecretModal] = useState<CreatedWebhook | null>(null)
  const [newUrl, setNewUrl] = useState('')
  const [newEvent, setNewEvent] = useState('payment.success')
  const [creating, setCreating] = useState(false)
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null)

  const fetchWebhooks = () => {
    setLoading(true)
    api.get<Webhook[]>('/webhook_endpoints')
      .then((res) => setWebhooks(res.data))
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchWebhooks() }, [])

  const handleCreate = async () => {
    if (!newUrl.trim()) return
    setCreating(true)
    try {
      const res = await api.post<CreatedWebhook>('/webhook_endpoints', {
        url: newUrl.trim(),
        event: newEvent,
      })
      setShowCreateModal(false)
      setNewUrl('')
      setNewEvent('payment.success')
      setShowSecretModal(res.data)
      fetchWebhooks()
    } catch {
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await api.delete(`/webhook_endpoints/${id}`)
      setDeleteConfirm(null)
      fetchWebhooks()
    } catch {}
  }

  const statusBadge = (status: string) => {
    if (status === 'active') return 'bg-success/10 text-success'
    if (status === 'disabled') return 'bg-slate-700 text-slate-300'
    return 'bg-slate-700 text-slate-300'
  }

  return (
    <DashboardLayout>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-white">Webhooks</h2>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Add Webhook
        </button>
      </div>

      <div className="overflow-hidden rounded-xl border border-slate-700">
        <table className="w-full text-left text-sm">
          <thead className="bg-slate-800">
            <tr>
              <th className="px-6 py-3 font-medium text-slate-400">URL</th>
              <th className="px-6 py-3 font-medium text-slate-400">Event</th>
              <th className="px-6 py-3 font-medium text-slate-400">Status</th>
              <th className="px-6 py-3 font-medium text-slate-400">Created</th>
              <th className="px-6 py-3 font-medium text-slate-400"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {loading ? (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-slate-500">
                  <div className="flex items-center justify-center gap-2">
                    <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                    Loading...
                  </div>
                </td>
              </tr>
            ) : webhooks.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-slate-500">
                  No webhooks configured
                </td>
              </tr>
            ) : (
              webhooks.map((wh) => (
                <tr key={wh.id} className="hover:bg-slate-800/50">
                  <td className="px-6 py-4 text-white font-mono text-xs max-w-[200px] truncate">
                    {wh.url}
                  </td>
                  <td className="px-6 py-4">
                    <code className="rounded bg-slate-700 px-2 py-0.5 font-mono text-xs text-slate-300">
                      {wh.event}
                    </code>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusBadge(wh.status)}`}>
                      {wh.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-400">
                    {new Date(wh.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4">
                    <button
                      onClick={() => setDeleteConfirm(wh.id)}
                      className="rounded px-3 py-1 text-xs font-medium text-danger hover:bg-danger/10 transition-colors"
                    >
                      Delete
                    </button>
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
              <h3 className="text-lg font-semibold text-white">Add Webhook</h3>
              <button onClick={() => setShowCreateModal(false)} className="text-slate-400 hover:text-white">
                <X className="h-5 w-5" />
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-slate-400 mb-1.5">Endpoint URL</label>
                <input
                  type="url"
                  placeholder="https://example.com/webhook"
                  value={newUrl}
                  onChange={(e) => setNewUrl(e.target.value)}
                  className="w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-sm text-white placeholder-slate-400 focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-sm text-slate-400 mb-1.5">Event Type</label>
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
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreate}
                disabled={!newUrl.trim() || creating}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 disabled:opacity-50 transition-colors"
              >
                {creating ? 'Creating...' : 'Create'}
              </button>
            </div>
          </div>
        </div>
      )}

      {showSecretModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="w-full max-w-lg rounded-xl border border-slate-700 bg-slate-800 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-white">Webhook Created</h3>
              <button onClick={() => setShowSecretModal(null)} className="text-slate-400 hover:text-white">
                <X className="h-5 w-5" />
              </button>
            </div>
            <p className="text-sm text-slate-400 mb-4">
              Use this signing secret to verify webhook payloads.
            </p>
            <div className="flex items-center gap-2 rounded-lg bg-slate-900 p-3">
              <code className="flex-1 font-mono text-xs text-green-400 break-all">
                {showSecretModal.signing_secret}
              </code>
              <CopyButton text={showSecretModal.signing_secret} />
            </div>
            <div className="mt-4 flex justify-end">
              <button
                onClick={() => setShowSecretModal(null)}
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
            <h3 className="text-lg font-semibold text-white mb-2">Delete Webhook</h3>
            <p className="text-sm text-slate-400 mb-4">
              Are you sure you want to delete this webhook endpoint?
            </p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => handleDelete(deleteConfirm)}
                className="rounded-lg bg-danger px-4 py-2 text-sm font-medium text-white hover:bg-danger-700 transition-colors"
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
