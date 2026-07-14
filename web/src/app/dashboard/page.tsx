import { ArrowUpRight, ArrowDownRight, DollarSign, Activity, Clock, TrendingUp } from 'lucide-react'

const stats = [
  {
    title: 'Total Transactions',
    value: '12,843',
    change: '+12.5%',
    icon: Activity,
    positive: true,
  },
  {
    title: 'Success Rate',
    value: '98.7%',
    change: '+0.3%',
    icon: TrendingUp,
    positive: true,
  },
  {
    title: 'Volume',
    value: '$1.2M',
    change: '+8.2%',
    icon: DollarSign,
    positive: true,
  },
  {
    title: 'Pending',
    value: '23',
    change: '-5.1%',
    icon: Clock,
    positive: false,
  },
]

const recentTransactions = [
  { id: 'TXN-001', customer: 'Acme Corp', amount: '$12,500', status: 'Success', date: '2024-01-15' },
  { id: 'TXN-002', customer: 'Globex Inc', amount: '$3,200', status: 'Pending', date: '2024-01-15' },
  { id: 'TXN-003', customer: 'Initech', amount: '$8,750', status: 'Success', date: '2024-01-14' },
  { id: 'TXN-004', customer: 'Hooli', amount: '$25,000', status: 'Failed', date: '2024-01-14' },
  { id: 'TXN-005', customer: 'Stark Ind', amount: '$6,430', status: 'Success', date: '2024-01-13' },
]

const sidebarLinks = [
  'Dashboard',
  'Payments',
  'Customers',
  'Merchants',
  'Settlements',
  'API Keys',
  'Settings',
]

export default function DashboardPage() {
  return (
    <div className="flex min-h-screen bg-slate-900">
      <aside className="w-64 border-r border-slate-700 bg-slate-800/50 p-6">
        <h2 className="text-lg font-bold text-white">OPG</h2>
        <nav className="mt-8 space-y-1">
          {sidebarLinks.map((link) => (
            <a
              key={link}
              href="#"
              className="block rounded-lg px-4 py-2 text-sm text-slate-400 hover:bg-slate-700 hover:text-white transition-colors"
            >
              {link}
            </a>
          ))}
        </nav>
      </aside>

      <div className="flex-1">
        <header className="border-b border-slate-700 px-8 py-4">
          <h1 className="text-xl font-semibold text-white">Dashboard</h1>
          <p className="text-sm text-slate-400">Welcome back, Admin</p>
        </header>

        <main className="p-8">
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {stats.map((stat) => (
              <div
                key={stat.title}
                className="rounded-xl border border-slate-700 bg-slate-800 p-6"
              >
                <div className="flex items-center justify-between">
                  <span className="text-sm text-slate-400">{stat.title}</span>
                  <stat.icon className="h-5 w-5 text-slate-500" />
                </div>
                <p className="mt-3 text-2xl font-bold text-white">{stat.value}</p>
                <div className="mt-1 flex items-center gap-1 text-sm">
                  {stat.positive ? (
                    <ArrowUpRight className="h-4 w-4 text-success" />
                  ) : (
                    <ArrowDownRight className="h-4 w-4 text-danger" />
                  )}
                  <span className={stat.positive ? 'text-success' : 'text-danger'}>
                    {stat.change}
                  </span>
                </div>
              </div>
            ))}
          </div>

          <div className="mt-8">
            <h2 className="mb-4 text-lg font-semibold text-white">Recent Transactions</h2>
            <div className="overflow-hidden rounded-xl border border-slate-700">
              <table className="w-full text-left text-sm">
                <thead className="bg-slate-800">
                  <tr>
                    <th className="px-6 py-3 font-medium text-slate-400">ID</th>
                    <th className="px-6 py-3 font-medium text-slate-400">Customer</th>
                    <th className="px-6 py-3 font-medium text-slate-400">Amount</th>
                    <th className="px-6 py-3 font-medium text-slate-400">Status</th>
                    <th className="px-6 py-3 font-medium text-slate-400">Date</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-700">
                  {recentTransactions.map((txn) => (
                    <tr key={txn.id} className="hover:bg-slate-800/50">
                      <td className="px-6 py-4 text-white">{txn.id}</td>
                      <td className="px-6 py-4 text-slate-300">{txn.customer}</td>
                      <td className="px-6 py-4 text-white">{txn.amount}</td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${
                            txn.status === 'Success'
                              ? 'bg-success/10 text-success'
                              : txn.status === 'Failed'
                              ? 'bg-danger/10 text-danger'
                              : 'bg-warning/10 text-warning'
                          }`}
                        >
                          {txn.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-slate-400">{txn.date}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}
