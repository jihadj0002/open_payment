import { HTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

interface SkeletonProps extends HTMLAttributes<HTMLDivElement> {
  variant?: 'text' | 'card' | 'table-row'
}

const variantStyles = {
  text: 'h-4 w-full',
  card: 'h-24 w-full rounded-xl',
  'table-row': 'h-12 w-full',
}

const Skeleton = forwardRef<HTMLDivElement, SkeletonProps>(
  ({ className, variant = 'text', ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={clsx(
          'animate-pulse rounded bg-slate-700',
          variantStyles[variant],
          className
        )}
        {...props}
      />
    )
  }
)
Skeleton.displayName = 'Skeleton'

export { Skeleton }
