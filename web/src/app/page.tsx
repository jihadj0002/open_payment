import Link from 'next/link'

const features = [
  {
    title: 'Multi-Processor',
    description:
      'Seamlessly route payments across multiple processors with intelligent failover and auto-retry logic.',
  },
  {
    title: 'Real-Time Ledger',
    description:
      'Double-entry accounting ledger with millisecond settlement tracking and full audit trails.',
  },
  {
    title: 'Fraud Prevention',
    description:
      'ML-powered fraud detection with configurable rules engine and real-time risk scoring.',
  },
]

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <header className="border-b border-white/10">
        <nav className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <span className="text-xl font-bold text-white">Open Payment Gateway</span>
          <div className="flex items-center gap-4">
            <Link
              href="/login"
              className="text-sm text-slate-300 hover:text-white transition-colors"
            >
              Sign in
            </Link>
            <Link
              href="/register"
              className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
            >
              Get started
            </Link>
          </div>
        </nav>
      </header>

      <main>
        <section className="mx-auto max-w-7xl px-6 pt-24 pb-16 text-center">
          <h1 className="text-5xl font-bold tracking-tight text-white sm:text-6xl">
            Open Payment Gateway
          </h1>
          <p className="mt-6 text-lg text-slate-300 max-w-2xl mx-auto">
            Modern, secure payment infrastructure for emerging markets
          </p>
          <div className="mt-10 flex items-center justify-center gap-4">
            <Link
              href="/merchant/dashboard"
              className="rounded-lg bg-primary px-6 py-3 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
            >
              Get Started
            </Link>
            <Link
              href="/docs"
              className="rounded-lg border border-white/20 px-6 py-3 text-sm font-medium text-slate-200 hover:bg-white/5 transition-colors"
            >
              Documentation
            </Link>
          </div>
        </section>

        <section className="mx-auto max-w-7xl px-6 pb-24">
          <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm"
              >
                <h3 className="text-lg font-semibold text-white">{feature.title}</h3>
                <p className="mt-3 text-sm text-slate-400">{feature.description}</p>
              </div>
            ))}
          </div>
        </section>
      </main>
    </div>
  )
}
