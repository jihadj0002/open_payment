import { HTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  header?: React.ReactNode
  footer?: React.ReactNode
  padding?: 'none' | 'sm' | 'md' | 'lg'
}

const paddingStyles = {
  none: '',
  sm: 'p-4',
  md: 'p-6',
  lg: 'p-8',
}

const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className, header, footer, padding = 'md', children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={clsx(
          'rounded-xl border border-slate-700 bg-slate-800',
          className
        )}
        {...props}
      >
        {header && (
          <div className="border-b border-slate-700 px-6 py-4">
            {header}
          </div>
        )}
        <div className={paddingStyles[padding]}>{children}</div>
        {footer && (
          <div className="border-t border-slate-700 px-6 py-4">
            {footer}
          </div>
        )}
      </div>
    )
  }
)
Card.displayName = 'Card'

export { Card }
