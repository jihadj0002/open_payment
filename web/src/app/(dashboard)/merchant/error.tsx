'use client'

import { useEffect } from 'react'

export default function MerchantError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error('Merchant dashboard error boundary caught:', error)
  }, [error])

  return (
    <div className="flex flex-col items-center justify-center py-20">
      <div className="mb-4 text-5xl font-bold text-primary/30">!</div>
      <h2 className="mb-2 text-xl font-semibold text-white">Dashboard Error</h2>
      <p className="mb-6 text-sm text-slate-400">
        Something went wrong loading this page. Please try again.
      </p>
      <button
        onClick={reset}
        className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
      >
        Try again
      </button>
    </div>
  )
}
