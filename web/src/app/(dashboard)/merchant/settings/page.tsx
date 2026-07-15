'use client'

import { useState, useEffect } from 'react'
import { Key, Shield, Lock, ExternalLink } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import Link from 'next/link'
import { Card } from '@/components/ui/Card'
import { Skeleton } from '@/components/ui/Skeleton'
import { Badge } from '@/components/ui/Badge'
import { Input } from '@/components/ui/Input'
import { Button } from '@/components/ui/Button'
import { toast } from '@/components/ui/Toast'
import { api } from '@/lib/api'
import { queryKeys } from '@/lib/query-keys'

interface MerchantProfile {
  id: string
  name: string
  email: string
  webhook_url: string
  status: string
  created_at: string
  updated_at: string
}

export default function SettingsPage() {
  const queryClient = useQueryClient()
  const [name, setName] = useState('')

  const { data: profile, isLoading, error } = useQuery({
    queryKey: queryKeys.merchant.profile,
    queryFn: () => api.get<MerchantProfile>('/merchants/profile').then(r => r.data),
  })

  useEffect(() => {
    if (profile) {
      setName(profile.name)
    }
  }, [profile])

  const updateMutation = useMutation({
    mutationFn: (newName: string) => api.patch('/merchants/profile', { name: newName }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.merchant.profile })
      toast.success('Profile updated successfully')
    },
    onError: () => {
      toast.error('Failed to update profile')
    },
  })

  const handleSave = () => {
    if (name.trim() && name !== profile?.name) {
      updateMutation.mutate(name.trim())
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton variant="text" className="h-8 w-32" />
        <Skeleton variant="card" className="h-48" />
        <Skeleton variant="card" className="h-32" />
        <Skeleton variant="card" className="h-24" />
      </div>
    )
  }

  if (error) {
    toast.error('Failed to load profile')
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <p className="text-slate-400">Failed to load profile</p>
      </div>
    )
  }

  return (
    <>
      <h1 className="mb-6 text-lg font-semibold text-white">Settings</h1>

      <div className="space-y-6">
        <Card header={<h2 className="text-base font-semibold text-white">Profile</h2>}>
          <div className="space-y-4">
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-300">Merchant ID</label>
              <code className="block rounded-lg bg-slate-900 px-4 py-2.5 font-mono text-xs text-slate-400">
                {profile?.id.slice(0, 8)}...
              </code>
            </div>
            <Input
              label="Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Merchant name"
            />
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-300">Email</label>
              <p className="rounded-lg border border-slate-600 bg-slate-700/50 px-4 py-2.5 text-sm text-slate-400">
                {profile?.email}
              </p>
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-300">Status</label>
              <Badge variant={profile?.status === 'active' ? 'success' : 'neutral'}>
                {profile?.status}
              </Badge>
            </div>
          </div>
          <div className="mt-6 flex justify-end border-t border-slate-700 pt-4">
            <Button
              onClick={handleSave}
              disabled={!name.trim() || name === profile?.name || updateMutation.isPending}
            >
              {updateMutation.isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </Card>

        <Card header={<h2 className="text-base font-semibold text-white">Security</h2>}>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Lock className="h-5 w-5 text-slate-400" />
                <div>
                  <p className="text-sm font-medium text-white">Change Password</p>
                  <p className="text-xs text-slate-400">Update your account password</p>
                </div>
              </div>
              <Button variant="secondary" disabled title="Coming soon">
                Change Password
              </Button>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Shield className="h-5 w-5 text-slate-400" />
                <div>
                  <p className="text-sm font-medium text-white">Two-Factor Authentication</p>
                  <p className="text-xs text-slate-400">Add an extra layer of security</p>
                </div>
              </div>
              <Button variant="secondary" disabled title="Coming soon">
                Enable 2FA
              </Button>
            </div>
          </div>
        </Card>

        <Link href="/merchant/api-keys">
          <Card className="cursor-pointer transition-colors hover:bg-slate-750">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Key className="h-5 w-5 text-slate-400" />
                <div>
                  <p className="text-sm font-medium text-white">API Keys</p>
                  <p className="text-xs text-slate-400">Manage your API keys</p>
                </div>
              </div>
              <ExternalLink className="h-5 w-5 text-slate-400" />
            </div>
          </Card>
        </Link>
      </div>
    </>
  )
}
