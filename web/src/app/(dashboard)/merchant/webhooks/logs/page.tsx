'use client'

import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { Card } from '@/components/ui/Card'
import { Badge } from '@/components/ui/Badge'
import { Skeleton } from '@/components/ui/Skeleton'
import { DataTable } from '@/components/ui/DataTable'
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

function statusBadgeVariant(status: string) {
  if (status === 'active') return 'success' as const
  if (status === 'disabled') return 'neutral' as const
  return 'neutral' as const
}

export default function WebhookLogsPage() {
  const { data: endpoints = [], isLoading, isError, error } = useQuery({
    queryKey: queryKeys.webhooks.all,
    queryFn: () => api.get<Webhook[]>('/webhook_endpoints').then(r => r.data),
  })

  if (isError) {
    toast.error(error instanceof Error ? error.message : 'Failed to load webhook endpoints')
  }

  return (
    <>
      <div className="mb-6">
        <Link
          href="/merchant/webhooks"
          className="mb-4 inline-flex items-center gap-1.5 text-sm text-slate-400 transition-colors hover:text-white"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to Webhooks
        </Link>
        <h2 className="mt-2 text-lg font-semibold text-white">Webhook Logs</h2>
      </div>

      <Card header={<span className="text-sm font-medium text-white">Endpoints Overview</span>}>
        {isLoading ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} variant="card" />
            ))}
          </div>
        ) : endpoints.length === 0 ? (
          <p className="py-6 text-center text-sm text-slate-400">No webhook endpoints configured</p>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {endpoints.map((ep) => (
              <div
                key={ep.id}
                className="rounded-lg border border-slate-700 bg-slate-800/50 p-4"
              >
                <div className="mb-2 flex items-center justify-between">
                  <Badge variant={statusBadgeVariant(ep.status)}>{ep.status}</Badge>
                </div>
                <p className="mb-1 truncate font-mono text-xs text-slate-300">{ep.url}</p>
                <code className="inline-block rounded bg-slate-700 px-2 py-0.5 font-mono text-xs text-slate-400">
                  {ep.event}
                </code>
              </div>
            ))}
          </div>
        )}
      </Card>

      <div className="mt-6">
        <Card header={<span className="text-sm font-medium text-white">Recent Deliveries</span>}>
          <div className="py-6 text-center">
            <p className="mb-4 text-sm text-slate-400">
              Webhook delivery logging is being implemented. Check back soon.
            </p>
            <div className="mx-auto max-w-2xl">
              <DataTable
                columns={[
                  { key: 'endpoint', header: 'Endpoint' },
                  { key: 'event', header: 'Event' },
                  { key: 'status', header: 'Status' },
                  { key: 'attempted_at', header: 'Attempted At' },
                ]}
                data={[]}
                emptyMessage="Delivery logs are not yet available"
              />
            </div>
          </div>
        </Card>
      </div>
    </>
  )
}
