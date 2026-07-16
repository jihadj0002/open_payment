'use client'

import { useState, useEffect, useCallback } from 'react'
import { useParams, useRouter } from 'next/navigation'
import {
  Loader2, XCircle, CheckCircle, Smartphone, Wallet, Building2,
  CreditCard, MoreHorizontal, Clock, ChevronLeft, QrCode,
  ShieldCheck, Lock
} from 'lucide-react'
import { api } from '@/lib/api'
import CardForm, { CardFormData } from '@/components/checkout/CardForm'
import MobileBankingGrid from '@/components/checkout/MobileBankingGrid'
import NetBankingSection from '@/components/checkout/NetBankingSection'

interface CheckoutSession {
  id: string
  amount: number
  currency: string
  status: string
  payment_method: string
  description: string
  return_url: string
  cancel_url: string
  redirect_url: string
  client_secret: string
  metadata: Record<string, string>
}

type PaymentTab = 'card' | 'mobile' | 'netbanking' | 'more'

const tabs: { id: PaymentTab; label: string; icon: React.ReactNode }[] = [
  { id: 'card', label: 'Card', icon: <CreditCard className="h-4 w-4" /> },
  { id: 'mobile', label: 'Mobile Banking', icon: <Smartphone className="h-4 w-4" /> },
  { id: 'netbanking', label: 'Net Banking', icon: <Building2 className="h-4 w-4" /> },
  { id: 'more', label: 'More', icon: <MoreHorizontal className="h-4 w-4" /> },
]

export default function CheckoutPage() {
  const params = useParams()
  const router = useRouter()
  const paymentIntentId = params.payment_intent_id as string

  const [session, setSession] = useState<CheckoutSession | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<PaymentTab>('card')
  const [selectedMobileProvider, setSelectedMobileProvider] = useState<string | null>(null)
  const [customerPhone, setCustomerPhone] = useState('')
  const [cardValid, setCardValid] = useState(false)
  const [cardData, setCardData] = useState<CardFormData | null>(null)
  const [processing, setProcessing] = useState(false)
  const [result, setResult] = useState<{ success: boolean; message: string; redirectUrl?: string } | null>(null)

  const [timeLeft, setTimeLeft] = useState(1800)

  useEffect(() => {
    const fetchSession = async () => {
      try {
        const res = await api.get<CheckoutSession>(`/checkout/${paymentIntentId}`)
        setSession(res.data)
      } catch (err: unknown) {
        const e = err as { message?: string }
        setError(e.message || 'Failed to load checkout session')
      } finally {
        setLoading(false)
      }
    }
    if (paymentIntentId) fetchSession()
  }, [paymentIntentId])

  useEffect(() => {
    if (timeLeft <= 0) return
    const timer = setInterval(() => setTimeLeft((t) => t - 1), 1000)
    return () => clearInterval(timer)
  }, [timeLeft])

  const formatCurrency = (amount: number, currency = 'BDT') =>
    new Intl.NumberFormat('en-BD', { style: 'currency', currency }).format(amount / 100)

  const formatTime = (seconds: number) => {
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }

  const handlePayment = useCallback(async () => {
    if (!session) return
    setProcessing(true)
    setResult(null)

    try {
      const body: Record<string, unknown> = {
        payment_method: 'wallet',
        return_url: `${window.location.origin}/checkout/${paymentIntentId}/success`,
        cancel_url: `${window.location.origin}/checkout/${paymentIntentId}/cancel`,
      }

      if (activeTab === 'mobile' && selectedMobileProvider) {
        body.provider = selectedMobileProvider
        body.customer_phone = customerPhone
      } else if (activeTab === 'card' && cardData) {
        body.provider = 'card'
        body.customer_phone = cardData.cardholderName
      } else {
        body.provider = 'bkash'
      }

      const res = await api.post<{ id: string; status: string; redirect_url: string }>(
        `/checkout/${paymentIntentId}/pay`,
        body
      )

      if (res.data.redirect_url) {
        window.location.href = res.data.redirect_url
      } else {
        setResult({ success: true, message: 'Payment initiated!' })
      }
    } catch (err: unknown) {
      const e = err as { message?: string }
      setResult({ success: false, message: e.message || 'Payment failed. Please try again.' })
    } finally {
      setProcessing(false)
    }
  }, [session, activeTab, selectedMobileProvider, customerPhone, cardData, paymentIntentId])

  const canPay = activeTab === 'card' ? cardValid : activeTab === 'mobile' ? !!selectedMobileProvider : false

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
        <div className="flex items-center gap-3">
          <Loader2 className="h-6 w-6 animate-spin text-primary" />
          <span className="text-slate-300">Loading checkout...</span>
        </div>
      </div>
    )
  }

  if (error || !session) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
        <XCircle className="mb-4 h-12 w-12 text-danger" />
        <h1 className="mb-2 text-2xl font-bold text-white">Checkout Unavailable</h1>
        <p className="mb-6 text-slate-400">{error || 'Payment session not found'}</p>
        <button
          onClick={() => router.push(session?.return_url || '/')}
          className="rounded-lg bg-primary px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          Return to Merchant
        </button>
      </div>
    )
  }

  if (session.status === 'succeeded' || session.status === 'captured') {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
        <CheckCircle className="mb-4 h-16 w-16 text-success" />
        <h1 className="mb-2 text-3xl font-bold text-white">Payment Successful</h1>
        <p className="mb-2 text-lg text-slate-300">{formatCurrency(session.amount, session.currency)}</p>
        <p className="mb-6 text-sm text-slate-400">Transaction completed</p>
        <button
          onClick={() => router.push(session.return_url || '/')}
          className="rounded-lg bg-primary px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          Return to Merchant
        </button>
      </div>
    )
  }

  const renderPaymentContent = () => {
    switch (activeTab) {
      case 'card':
        return (
          <CardForm
            onValidChange={(isValid, data) => {
              setCardValid(isValid)
              setCardData(data)
            }}
          />
        )
      case 'mobile':
        return (
          <MobileBankingGrid
            selectedProvider={selectedMobileProvider}
            onSelect={setSelectedMobileProvider}
            customerPhone={customerPhone}
            onPhoneChange={setCustomerPhone}
          />
        )
      case 'netbanking':
        return <NetBankingSection />
      case 'more':
        return (
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-slate-300">More Payment Options</h3>
            {[
              { title: 'EMI', desc: 'Pay in installments' },
              { title: 'Wallets', desc: 'Digital wallets' },
              { title: 'International Cards', desc: 'Visa, Mastercard abroad' },
              { title: 'Digital Wallets', desc: 'Google Pay, Apple Pay' },
            ].map((item) => (
              <div
                key={item.title}
                className="rounded-lg border border-slate-700 bg-slate-800/50 p-4 opacity-50 cursor-not-allowed"
              >
                <span className="block text-sm font-medium text-slate-300">{item.title}</span>
                <span className="text-xs text-slate-500">{item.desc} — Coming Soon</span>
              </div>
            ))}
          </div>
        )
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="mx-auto max-w-6xl px-4 py-4 lg:py-8">
        {/* Header */}
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button onClick={() => router.push(session.return_url || '/')} className="text-slate-400 hover:text-white">
              <ChevronLeft className="h-5 w-5" />
            </button>
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/20 text-primary text-xs font-bold">
                OP
              </div>
              <span className="text-sm font-medium text-white">Merchant</span>
            </div>
            <span className="hidden sm:inline text-xs text-slate-500">TXN: {paymentIntentId.slice(0, 8)}</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-1.5 text-slate-400">
              <Clock className="h-4 w-4" />
              <span className={`text-sm font-mono ${timeLeft < 300 ? 'text-red-400' : ''}`}>
                {formatTime(timeLeft)}
              </span>
            </div>
            <button className="hidden sm:inline text-xs text-slate-400 hover:text-white">Login</button>
          </div>
        </div>

        {/* Main Content: Desktop two-panel, mobile single column */}
        <div className="flex flex-col lg:flex-row gap-6">

          {/* LEFT PANEL - Payment Summary (Desktop) */}
          <div className="lg:w-80 xl:w-96 space-y-4 order-2 lg:order-1">
            {/* Amount Card */}
            <div className="rounded-xl border border-slate-700 bg-slate-800/80 p-5">
              <p className="text-xs text-slate-400 mb-1">You are paying</p>
              <div className="flex items-center justify-between">
                <span className="text-3xl font-bold text-white">
                  {formatCurrency(session.amount, session.currency)}
                </span>
                <button className="flex items-center gap-1 rounded-lg border border-slate-600 px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-700">
                  <QrCode className="h-4 w-4" />
                  QR Pay
                </button>
              </div>
            </div>

            {/* Cost Breakdown */}
            <div className="rounded-xl border border-slate-700 bg-slate-800/80 p-5 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-slate-400">Subtotal</span>
                <span className="text-white">{formatCurrency(session.amount, session.currency)}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-slate-400">Convenience Fee</span>
                <span className="text-white">{formatCurrency(0)}</span>
              </div>
              <hr className="border-slate-700" />
              <div className="flex justify-between text-sm font-semibold">
                <span className="text-white">TOTAL</span>
                <span className="text-white">{formatCurrency(session.amount, session.currency)}</span>
              </div>
            </div>

            {/* Offers */}
            <div className="rounded-xl border border-slate-700 bg-slate-800/80 p-5">
              <h4 className="text-xs font-semibold uppercase tracking-wider text-slate-400 mb-3">
                Special Offers &amp; Savings
              </h4>
              <div className="flex flex-col items-center justify-center py-4 text-slate-500">
                <div className="text-3xl mb-2">🎉</div>
                <p className="text-xs">No offers available</p>
              </div>
            </div>

            {/* Footer Links (Desktop only) */}
            <div className="hidden lg:flex items-center gap-4 text-xs text-slate-500">
              <span>Support</span>
              <span>FAQ</span>
              <span>Language</span>
            </div>
          </div>

          {/* RIGHT PANEL - Payment Area */}
          <div className="flex-1 order-1 lg:order-2">
            <div className="rounded-xl border border-slate-700 bg-slate-800/80 p-5 lg:p-6">

              {/* Payment Method Tabs */}
              <div className="flex overflow-x-auto gap-1 mb-6 border-b border-slate-700 pb-1">
                {tabs.map((tab) => (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    className={`flex items-center gap-1.5 px-3 py-2 text-xs font-medium whitespace-nowrap border-b-2 transition-all ${
                      activeTab === tab.id
                        ? 'border-primary text-primary'
                        : 'border-transparent text-slate-400 hover:text-slate-300'
                    }`}
                  >
                    {tab.icon}
                    <span className="hidden sm:inline">{tab.label}</span>
                  </button>
                ))}
              </div>

              {/* Dynamic Payment Content */}
              <div className="mb-6">
                {renderPaymentContent()}
              </div>

              {/* Pay Button */}
              <button
                onClick={handlePayment}
                disabled={!canPay || processing}
                className="w-full rounded-lg bg-primary px-6 py-3.5 text-sm font-semibold text-white transition-all hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50 flex items-center justify-center gap-2"
              >
                {processing ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" />
                    Processing...
                  </>
                ) : (
                  <>
                    <Lock className="h-4 w-4" />
                    Pay {formatCurrency(session.amount, session.currency)}
                  </>
                )}
              </button>

              {/* Terms */}
              <p className="mt-4 text-center text-[10px] text-slate-500">
                By clicking Pay, you agree to our{' '}
                <span className="underline cursor-pointer hover:text-slate-400">Terms of Service</span>
              </p>
            </div>

            {/* Footer */}
            <div className="mt-4 flex flex-col items-center gap-2 text-xs text-slate-500">
              <div className="flex items-center gap-2">
                <span>Powered by</span>
                <span className="font-semibold text-slate-400">SSLCommerz</span>
              </div>
              <div className="flex items-center gap-3">
                <ShieldCheck className="h-4 w-4" />
                <span>PCI DSS Compliant</span>
                <Lock className="h-4 w-4" />
                <span>256-bit Secure</span>
              </div>
              <div className="flex lg:hidden items-center gap-4 mt-2">
                <span>Support</span>
                <span>FAQ</span>
                <span>Language</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
