'use client'

import { Smartphone, Wallet } from 'lucide-react'

interface Provider {
  id: string
  name: string
  icon: string
  available: boolean
}

const providers: Provider[] = [
  { id: 'bkash', name: 'bKash', icon: 'B', available: true },
  { id: 'nagad', name: 'Nagad', icon: 'N', available: true },
  { id: 'rocket', name: 'Rocket', icon: 'R', available: false },
  { id: 'upay', name: 'Upay', icon: 'U', available: false },
  { id: 'cellfin', name: 'CellFin', icon: 'C', available: false },
  { id: 'mcash', name: 'MCash', icon: 'M', available: false },
  { id: 'pocket', name: 'Pocket', icon: 'P', available: false },
  { id: 'pathao', name: 'Pathao Pay', icon: 'PP', available: false },
  { id: 'rainbow', name: 'Rainbow', icon: 'R', available: false },
  { id: 'tap', name: 'Tap', icon: 'T', available: false },
  { id: 'ok_wallet', name: 'OK Wallet', icon: 'OK', available: false },
  { id: 'islami', name: 'Islami Wallet', icon: 'IW', available: false },
]

interface MobileBankingGridProps {
  selectedProvider: string | null
  onSelect: (providerId: string) => void
  customerPhone: string
  onPhoneChange: (phone: string) => void
}

export default function MobileBankingGrid({
  selectedProvider,
  onSelect,
  customerPhone,
  onPhoneChange,
}: MobileBankingGridProps) {
  return (
    <div className="space-y-4">
      <h3 className="text-sm font-medium text-slate-300">Pay with Mobile Banking</h3>

      <div className="grid grid-cols-3 gap-3">
        {providers.map((provider) => (
          <button
            key={provider.id}
            onClick={() => provider.available && onSelect(provider.id)}
            disabled={!provider.available}
            className={`flex flex-col items-center gap-2 rounded-lg border p-3 transition-all ${
              selectedProvider === provider.id
                ? 'border-primary bg-primary/10 ring-1 ring-primary'
                : provider.available
                  ? 'border-slate-600 bg-slate-700/50 hover:border-slate-500 hover:bg-slate-700'
                  : 'border-slate-700 bg-slate-800/30 opacity-40 cursor-not-allowed'
            }`}
          >
            <div
              className={`flex h-10 w-10 items-center justify-center rounded-full text-sm font-bold ${
                selectedProvider === provider.id
                  ? 'bg-primary text-white'
                  : 'bg-slate-600 text-slate-300'
              }`}
            >
              {provider.id === 'bkash' ? (
                <Smartphone className="h-5 w-5" />
              ) : provider.id === 'nagad' ? (
                <Wallet className="h-5 w-5" />
              ) : (
                <span>{provider.icon}</span>
              )}
            </div>
            <span className="text-xs text-slate-300 text-center leading-tight">{provider.name}</span>
            {!provider.available && (
              <span className="text-[10px] text-slate-500 -mt-1">Soon</span>
            )}
          </button>
        ))}
      </div>

      {selectedProvider && (
        <div>
          <label htmlFor="mobile-phone" className="mb-1 block text-xs text-slate-400">
            Phone Number
          </label>
          <input
            id="mobile-phone"
            type="tel"
            value={customerPhone}
            onChange={(e) => onPhoneChange(e.target.value.replace(/\D/g, '').slice(0, 11))}
            placeholder="01XXXXXXXXX"
            className="w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
        </div>
      )}
    </div>
  )
}
