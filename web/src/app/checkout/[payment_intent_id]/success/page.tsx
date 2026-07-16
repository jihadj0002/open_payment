'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { CheckCircle, Loader2 } from 'lucide-react'
import { api } from '@/lib/api'

export default function CheckoutSuccessPage() {
  const params = useParams()
  const router = useRouter()
  const paymentIntentId = params.payment_intent_id as string

  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<{ amount: number; currency: string; returnURL: string } | null>(null)

  useEffect(() => {
    const check = async () => {
      try {
        const res = await api.get<{
          success: boolean
          status: string
          amount: number
          currency: string
          returnURL: string
        }>(`/checkout/${paymentIntentId}/success`)
        setData(res.data)
      } catch {
      } finally {
        setLoading(false)
      }
    }
    if (paymentIntentId) check()
  }, [paymentIntentId])

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="mx-auto max-w-md px-4 text-center">
        <div className="mb-6 inline-flex h-20 w-20 items-center justify-center rounded-full bg-success/20">
          <CheckCircle className="h-10 w-10 text-success" />
        </div>
        <h1 className="mb-2 text-3xl font-bold text-white">Payment Successful!</h1>
        <p className="mb-2 text-lg text-slate-300">Thank you for your payment</p>
        {data && (
          <p className="mb-8 text-sm text-slate-400">
            Amount:{' '}
            {new Intl.NumberFormat('en-BD', {
              style: 'currency',
              currency: data.currency || 'BDT',
            }).format((data.amount || 0) / 100)}
          </p>
        )}
        <button
          onClick={() => router.push(data?.returnURL || '/')}
          className="rounded-lg bg-primary px-8 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          Return to Merchant
        </button>
      </div>
    </div>
  )
}
