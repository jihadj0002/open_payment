'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useAuth } from '@/context/AuthContext'
import { Loader2 } from 'lucide-react'

const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
})

type LoginForm = z.infer<typeof loginSchema>

export default function LoginPage() {
  const router = useRouter()
  const { login } = useAuth()
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data: LoginForm) => {
    setError('')
    setIsSubmitting(true)
    try {
      await login(data.email, data.password)
      router.push('/merchant/dashboard')
    } catch (err: unknown) {
      const e = err as { message?: string }
      setError(e.message || 'Login failed. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <>
      <h1 className="text-2xl font-bold text-white">Sign in to your account</h1>
      <p className="mt-2 text-sm text-slate-400">
        Enter your credentials to access your dashboard
      </p>

      {error && (
        <div className="mt-4 rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="mt-8 space-y-5">
        <div>
          <label htmlFor="email" className="block text-sm font-medium text-slate-300">
            Email
          </label>
          <input
            id="email"
            type="email"
            {...register('email')}
            className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
            placeholder="you@example.com"
          />
          {errors.email && (
            <p className="mt-1 text-sm text-danger">{errors.email.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-slate-300">
            Password
          </label>
          <input
            id="password"
            type="password"
            {...register('password')}
            className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
            placeholder="••••••••"
          />
          {errors.password && (
            <p className="mt-1 text-sm text-danger">{errors.password.message}</p>
          )}
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="flex w-full items-center justify-center rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-white hover:bg-primary-700 transition-colors disabled:opacity-50"
        >
          {isSubmitting ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            'Sign in'
          )}
        </button>
      </form>

      <p className="mt-6 text-center text-sm text-slate-400">
        Don&apos;t have an account?{' '}
        <Link href="/register" className="font-medium text-primary hover:text-primary-400">
          Register
        </Link>
      </p>
    </>
  )
}
