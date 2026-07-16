'use client'

import { useEffect } from 'react'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error('Root error boundary caught:', error)
  }, [error])

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-slate-900 px-4">
      <div className="w-full max-w-md text-center">
        <div className="mb-6 text-6xl font-bold text-primary/30">500</div>
        <h1 className="mb-2 text-2xl font-bold text-white">Something went wrong</h1>
        <p className="mb-8 text-sm text-slate-400">
          An unexpected error occurred. Please try again or contact support.
        </p>
        <button
          onClick={reset}
          className="rounded-lg bg-primary px-6 py-2.5 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          Try again
        </button>
      </div>
    </div>
  )
}
