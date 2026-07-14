'use client'

import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

const registerSchema = z
  .object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.string().email('Invalid email address'),
    password: z.string().min(6, 'Password must be at least 6 characters'),
    confirmPassword: z.string().min(6, 'Confirm password must be at least 6 characters'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  })

type RegisterForm = z.infer<typeof registerSchema>

export default function RegisterPage() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
  })

  const onSubmit = (data: RegisterForm) => {
    console.log('Register:', data)
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-900 px-4">
      <div className="w-full max-w-md rounded-xl border border-slate-700 bg-slate-800 p-8">
        <h1 className="text-2xl font-bold text-white">Create your account</h1>
        <p className="mt-2 text-sm text-slate-400">
          Fill in the details below to get started
        </p>

        <form onSubmit={handleSubmit(onSubmit)} className="mt-8 space-y-5">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-slate-300">
              Full Name
            </label>
            <input
              id="name"
              type="text"
              {...register('name')}
              className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
              placeholder="John Doe"
            />
            {errors.name && (
              <p className="mt-1 text-sm text-danger">{errors.name.message}</p>
            )}
          </div>

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

          <div>
            <label
              htmlFor="confirmPassword"
              className="block text-sm font-medium text-slate-300"
            >
              Confirm Password
            </label>
            <input
              id="confirmPassword"
              type="password"
              {...register('confirmPassword')}
              className="mt-1 block w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
              placeholder="••••••••"
            />
            {errors.confirmPassword && (
              <p className="mt-1 text-sm text-danger">{errors.confirmPassword.message}</p>
            )}
          </div>

          <button
            type="submit"
            className="w-full rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
          >
            Create account
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-400">
          Already have an account?{' '}
          <Link href="/login" className="font-medium text-primary hover:text-primary-400">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
