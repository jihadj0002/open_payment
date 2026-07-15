'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2, CheckCircle, XCircle } from 'lucide-react'
import { api } from '@/lib/api'

const paymentSchema = z.object({
  amount: z.coerce.number().positive('Amount must be positive'),
  currency: z.enum(['BDT', 'USD'], { required_error: 'Select a currency' }),
  payment_method: z.enum(['card', 'wallet', 'bank_transfer'], { required_error: 'Select a payment method' }),
  description: z.string().optional(),
})

type PaymentForm = z.infer<typeof paymentSchema>

export default function NewPaymentPage() {
  const router = useRouter()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [result, setResult] = useState<{ success: boolean; message: string } | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<PaymentForm>({
    resolver: zodResolver(paymentSchema),
    defaultValues: { currency: 'BDT', payment_method: 'card' },
  })

  const onSubmit = async (data: PaymentForm) => {
    setIsSubmitting(true)
    setResult(null)
    try {
      await api.post('/payments', data)
      setResult({ success: true, message: 'Payment created successfully!' })
      setTimeout(() => router.push('/merchant/payments'), 1500)
    } catch (err: unknown) {
      const e = err as { message?: string }
      setResult({ success: false, message: e.message || 'Payment failed. Please try again.' })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <>
      <h2 className="mb-6 text-lg font-semibold text-white">New Payment</h2>

      {result && (
        <div className={`mb-6 flex items-center gap-3 rounded-lg border px-4 py-3 text-sm ${
          result.success
            ? 'border-success/20 bg-success/10 text-success'
            : 'border-danger/20 bg-danger/10 text-danger'
        }`}>
          {result.success ? <CheckCircle className="h-5 w-5" /> : <XCircle className="h-5 w-5" />}
          {result.message}
        </div>
      )}

      <div className="max-w-lg rounded-xl border border-slate-700 bg-slate-800 p-6">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          <div>
            <label htmlFor="amount" className="block text-sm font-medium text-slate-300">
              Amount
            </label>
            <input
              id="amount"
              type="number"
              step="0.01"
              {...register('amount')}
              className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
              placeholder="0.00"
            />
            {errors.amount && (
              <p className="mt-1 text-sm text-danger">{errors.amount.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="currency" className="block text-sm font-medium text-slate-300">
              Currency
            </label>
            <select
              id="currency"
              {...register('currency')}
              className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
            >
              <option value="BDT">BDT</option>
              <option value="USD">USD</option>
            </select>
            {errors.currency && (
              <p className="mt-1 text-sm text-danger">{errors.currency.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="payment_method" className="block text-sm font-medium text-slate-300">
              Payment Method
            </label>
            <select
              id="payment_method"
              {...register('payment_method')}
              className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
            >
              <option value="card">Card</option>
              <option value="wallet">Wallet</option>
              <option value="bank_transfer">Bank Transfer</option>
            </select>
            {errors.payment_method && (
              <p className="mt-1 text-sm text-danger">{errors.payment_method.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="description" className="block text-sm font-medium text-slate-300">
              Description (optional)
            </label>
            <textarea
              id="description"
              rows={3}
              {...register('description')}
              className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary resize-none"
              placeholder="Payment for..."
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="flex w-full items-center justify-center rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-white hover:bg-primary-700 transition-colors disabled:opacity-50"
          >
            {isSubmitting ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              'Create Payment'
            )}
          </button>
        </form>
      </div>
    </>
  )
}
