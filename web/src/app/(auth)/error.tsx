'use client'

import { useEffect } from 'react'

export default function AuthError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error('Auth error boundary caught:', error)
  }, [error])

  return (
    <div className="text-center">
      <div className="mb-4 text-4xl font-bold text-primary/30">!</div>
      <h2 className="mb-2 text-lg font-semibold text-white">Authentication Error</h2>
      <p className="mb-6 text-sm text-slate-400">
        Something went wrong. Please try again.
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
