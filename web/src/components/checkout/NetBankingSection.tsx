'use client'

import { Building2 } from 'lucide-react'

const banks = [
  'Sonali Bank',
  'Janata Bank',
  'Agrani Bank',
  'Rupali Bank',
  'Dutch-Bangla Bank',
  'BRAC Bank',
  'Eastern Bank',
  'City Bank',
]

export default function NetBankingSection() {
  return (
    <div className="space-y-4">
      <h3 className="text-sm font-medium text-slate-300">Pay with Net Banking</h3>

      <div className="grid grid-cols-2 gap-2">
        {banks.map((bank) => (
          <button
            key={bank}
            disabled
            className="flex items-center gap-3 rounded-lg border border-slate-700 bg-slate-800/30 p-3 opacity-40 cursor-not-allowed"
          >
            <Building2 className="h-5 w-5 text-slate-500 shrink-0" />
            <span className="text-xs text-slate-400 text-left">{bank}</span>
          </button>
        ))}
      </div>

      <p className="text-center text-xs text-slate-500">Net Banking coming soon</p>
    </div>
  )
}
