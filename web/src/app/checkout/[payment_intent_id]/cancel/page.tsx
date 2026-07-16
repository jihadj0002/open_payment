'use client'

import { useParams, useRouter } from 'next/navigation'
import { XCircle } from 'lucide-react'
import CheckoutErrorBoundary from '@/components/checkout/CheckoutErrorBoundary'

export default function CheckoutCancelPage() {
  return (
    <CheckoutErrorBoundary>
      <CheckoutCancelContent />
    </CheckoutErrorBoundary>
  )
}

function CheckoutCancelContent() {
  const params = useParams()
  const router = useRouter()
  const paymentIntentId = params.payment_intent_id as string

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="mx-auto max-w-md px-4 text-center">
        <div className="mb-6 inline-flex h-20 w-20 items-center justify-center rounded-full bg-warning/20">
          <XCircle className="h-10 w-10 text-warning" />
        </div>
        <h1 className="mb-2 text-3xl font-bold text-white">Payment Cancelled</h1>
        <p className="mb-8 text-slate-400">
          You have cancelled the payment. No charges have been made.
        </p>
        <div className="flex items-center justify-center gap-4">
          <button
            onClick={() => router.push(`/checkout/${paymentIntentId}`)}
            className="rounded-lg border border-slate-600 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-slate-700"
          >
            Try Again
          </button>
          <button
            onClick={() => router.push('/')}
            className="rounded-lg bg-primary px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-700"
          >
            Go Home
          </button>
        </div>
      </div>
    </div>
  )
}
