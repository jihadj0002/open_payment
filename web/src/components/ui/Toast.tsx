'use client'

import { useSyncExternalStore, useCallback } from 'react'
import { clsx } from 'clsx'
import { X, CheckCircle, AlertCircle, Info } from 'lucide-react'

interface Toast {
  id: string
  type: 'success' | 'error' | 'info'
  message: string
  duration?: number
}

const typeStyles = {
  success: 'border-success/30 bg-success/10 text-success',
  error: 'border-danger/30 bg-danger/10 text-danger',
  info: 'border-blue-500/30 bg-blue-500/10 text-blue-400',
}

const icons = {
  success: CheckCircle,
  error: AlertCircle,
  info: Info,
}

let toasts: Toast[] = []
let listeners: (() => void)[] = []

function subscribe(listener: () => void) {
  listeners.push(listener)
  return () => {
    listeners = listeners.filter((l) => l !== listener)
  }
}

function getSnapshot() {
  return toasts
}

function addToast(toast: Omit<Toast, 'id'>) {
  const id = Math.random().toString(36).slice(2, 10)
  toasts = [...toasts, { ...toast, id }]
  listeners.forEach((l) => l())

  const duration = toast.duration ?? 4000
  if (duration > 0) {
    setTimeout(() => removeToast(id), duration)
  }
}

function removeToast(id: string) {
  toasts = toasts.filter((t) => t.id !== id)
  listeners.forEach((l) => l())
}

const toast = {
  success: (message: string) => addToast({ type: 'success', message }),
  error: (message: string) => addToast({ type: 'error', message }),
  info: (message: string) => addToast({ type: 'info', message }),
}

function Toaster() {
  const currentToasts = useSyncExternalStore(subscribe, getSnapshot, getSnapshot)

  return (
    <div className="fixed bottom-4 right-4 z-[100] flex flex-col-reverse gap-2">
      {currentToasts.map((t) => {
        const Icon = icons[t.type]
        return (
          <div
            key={t.id}
            className={clsx(
              'flex items-center gap-3 rounded-lg border px-4 py-3 shadow-lg transition-all duration-300',
              typeStyles[t.type]
            )}
            style={{
              animation: 'slideInRight 0.3s ease-out',
            }}
          >
            <Icon className="h-5 w-5 shrink-0" />
            <p className="text-sm font-medium">{t.message}</p>
            <button
              onClick={() => removeToast(t.id)}
              className="ml-auto shrink-0 rounded p-0.5 opacity-70 transition-opacity hover:opacity-100"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        )
      })}
      <style jsx global>{`
        @keyframes slideInRight {
          from {
            transform: translateX(100%);
            opacity: 0;
          }
          to {
            transform: translateX(0);
            opacity: 1;
          }
        }
      `}</style>
    </div>
  )
}

export { Toaster, toast }
export type { Toast }
