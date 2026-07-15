'use client'

import { useEffect, useCallback, useRef, useState } from 'react'
import { clsx } from 'clsx'
import { X } from 'lucide-react'

interface ModalProps {
  open: boolean
  onClose: () => void
  title?: string
  children: React.ReactNode
  footer?: React.ReactNode
  size?: 'sm' | 'md' | 'lg'
}

const sizeStyles = {
  sm: 'max-w-md',
  md: 'max-w-lg',
  lg: 'max-w-xl',
}

function Modal({ open, onClose, title, children, footer, size = 'md' }: ModalProps) {
  const overlayRef = useRef<HTMLDivElement>(null)
  const [visible, setVisible] = useState(false)
  const [animating, setAnimating] = useState(false)

  useEffect(() => {
    if (open) {
      setVisible(true)
      requestAnimationFrame(() => setAnimating(true))
    } else {
      setAnimating(false)
      const timer = setTimeout(() => setVisible(false), 200)
      return () => clearTimeout(timer)
    }
  }, [open])

  const handleEscape = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    },
    [onClose]
  )

  useEffect(() => {
    if (open) {
      document.addEventListener('keydown', handleEscape)
      return () => document.removeEventListener('keydown', handleEscape)
    }
  }, [open, handleEscape])

  if (!visible) return null

  return (
    <div
      ref={overlayRef}
      className={clsx(
        'fixed inset-0 z-50 flex items-center justify-center p-4 transition-colors duration-200',
        animating ? 'bg-black/60' : 'bg-transparent'
      )}
      onClick={(e) => {
        if (e.target === overlayRef.current) onClose()
      }}
    >
      <div
        className={clsx(
          'w-full rounded-xl border border-slate-700 bg-slate-800 shadow-xl transition-all duration-200',
          sizeStyles[size],
          animating ? 'scale-100 opacity-100' : 'scale-95 opacity-0'
        )}
      >
        {title && (
          <div className="flex items-center justify-between border-b border-slate-700 px-6 py-4">
            <h2 className="text-lg font-semibold text-white">{title}</h2>
              <button
                onClick={onClose}
                aria-label="Close"
                className="rounded-lg p-1 text-slate-400 transition-colors hover:bg-slate-700 hover:text-white"
              >
              <X className="h-5 w-5" />
            </button>
          </div>
        )}
        <div className="px-6 py-4">{children}</div>
        {footer && (
          <div className="flex items-center justify-end gap-3 border-t border-slate-700 px-6 py-4">
            {footer}
          </div>
        )}
      </div>
    </div>
  )
}

export { Modal }
export type { ModalProps }
