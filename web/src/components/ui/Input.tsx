import { InputHTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  icon?: React.ReactNode
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, icon, className, id, ...props }, ref) => {
    return (
      <div>
        {label && (
          <label
            htmlFor={id}
            className="mb-1 block text-sm font-medium text-slate-300"
          >
            {label}
          </label>
        )}
        <div className="relative">
          {icon && (
            <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400">
              {icon}
            </div>
          )}
          <input
            id={id}
            ref={ref}
            className={clsx(
              'block w-full rounded-lg border bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 transition-colors focus:outline-none focus:ring-1',
              error
                ? 'border-danger focus:border-danger focus:ring-danger'
                : 'border-slate-600 focus:border-primary focus:ring-primary',
              icon && 'pl-10',
              className
            )}
            {...props}
          />
        </div>
        {error && (
          <p className="mt-1 text-sm text-danger">{error}</p>
        )}
      </div>
    )
  }
)
Input.displayName = 'Input'

export { Input }
