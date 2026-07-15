export default function AuthLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-900 px-4">
      <div className="w-full max-w-md rounded-xl border border-slate-700 bg-slate-800 p-8">
        {children}
      </div>
    </div>
  )
}
