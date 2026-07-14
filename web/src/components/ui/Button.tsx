import { ButtonHTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

const variantStyles = {
  primary: 'bg-primary text-white hover:bg-primary-700',
  secondary: 'bg-slate-700 text-slate-200 hover:bg-slate-600',
  outline: 'border border-slate-600 text-slate-300 hover:bg-slate-800',
  ghost: 'text-slate-400 hover:text-white hover:bg-slate-800',
  danger: 'bg-danger text-white hover:bg-danger-600',
} as const

const sizeStyles = {
  sm: 'h-8 px-3 text-xs',
  md: 'h-10 px-4 text-sm',
  lg: 'h-12 px-6 text-base',
} as const

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: keyof typeof variantStyles
  size?: keyof typeof sizeStyles
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', ...props }, ref) => {
    return (
      <button
        className={clsx(
          'inline-flex items-center justify-center rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 focus:ring-offset-slate-900 disabled:pointer-events-none disabled:opacity-50',
          variantStyles[variant],
          sizeStyles[size],
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = 'Button'

export { Button }
