'use client'

import { useState } from 'react'
import { CreditCard, Lock } from 'lucide-react'

interface CardFormProps {
  onValidChange: (isValid: boolean, data: CardFormData | null) => void
}

export interface CardFormData {
  cardNumber: string
  expiryMonth: string
  expiryYear: string
  cvv: string
  cardholderName: string
  remember: boolean
}

function luhnCheck(cardNumber: string): boolean {
  const digits = cardNumber.replace(/\D/g, '')
  if (digits.length < 13 || digits.length > 19) return false
  let sum = 0
  let alternate = false
  for (let i = digits.length - 1; i >= 0; i--) {
    let n = parseInt(digits[i], 10)
    if (alternate) {
      n *= 2
      if (n > 9) n -= 9
    }
    sum += n
    alternate = !alternate
  }
  return sum % 10 === 0
}

function formatCardNumber(value: string): string {
  const digits = value.replace(/\D/g, '').slice(0, 16)
  const groups = digits.match(/.{1,4}/g)
  return groups ? groups.join(' ') : digits
}

function detectCardBrand(number: string): string {
  const clean = number.replace(/\s/g, '')
  if (/^4/.test(clean)) return 'visa'
  if (/^5[1-5]/.test(clean)) return 'mastercard'
  if (/^3[47]/.test(clean)) return 'amex'
  if (/^6(?:011|5)/.test(clean)) return 'discover'
  return 'unknown'
}

export default function CardForm({ onValidChange }: CardFormProps) {
  const [cardNumber, setCardNumber] = useState('')
  const [expiry, setExpiry] = useState('')
  const [cvv, setCvv] = useState('')
  const [cardholderName, setCardholderName] = useState('')
  const [remember, setRemember] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validate = (): boolean => {
    const errs: Record<string, string> = {}
    const cleanNumber = cardNumber.replace(/\s/g, '')

    if (!cleanNumber) errs.cardNumber = 'Card number is required'
    else if (!luhnCheck(cleanNumber)) errs.cardNumber = 'Invalid card number'

    if (!expiry) errs.expiry = 'Expiry is required'
    else {
      const [month, year] = expiry.split('/').map(s => s.trim())
      if (!month || !year || parseInt(month) < 1 || parseInt(month) > 12) {
        errs.expiry = 'Invalid expiry date'
      } else {
        const now = new Date()
        const expYear = 2000 + parseInt(year)
        const expMonth = parseInt(month)
        if (expYear < now.getFullYear() || (expYear === now.getFullYear() && expMonth < now.getMonth() + 1)) {
          errs.expiry = 'Card has expired'
        }
      }
    }

    if (!cvv) errs.cvv = 'CVV is required'
    else if (cvv.length < 3) errs.cvv = 'Invalid CVV'

    if (!cardholderName.trim()) errs.cardholderName = 'Cardholder name is required'

    setErrors(errs)
    const isValid = Object.keys(errs).length === 0

    if (isValid) {
      const [month, year] = expiry.split('/').map(s => s.trim())
      onValidChange(true, {
        cardNumber: cleanNumber,
        expiryMonth: month,
        expiryYear: year,
        cvv,
        cardholderName: cardholderName.trim(),
        remember,
      })
    } else {
      onValidChange(false, null)
    }

    return isValid
  }

  const brand = detectCardBrand(cardNumber)
  const brandIcon = brand === 'visa' ? 'VISA' : brand === 'mastercard' ? 'MC' : brand === 'amex' ? 'AMEX' : ''

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-slate-300">Pay with New Card</h3>
        <button className="text-xs text-primary hover:underline">Other Cards</button>
      </div>

      <div>
        <div className="relative">
          <input
            type="text"
            inputMode="numeric"
            placeholder="Card Number"
            value={cardNumber}
            onChange={(e) => {
              setCardNumber(formatCardNumber(e.target.value))
              setErrors(prev => ({ ...prev, cardNumber: '' }))
            }}
            onBlur={validate}
            className="w-full rounded-lg border border-slate-600 bg-slate-700 pl-4 pr-16 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
          <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-1">
            {brandIcon && (
              <span className="text-xs font-bold text-slate-400 bg-slate-600 px-1.5 py-0.5 rounded">{brandIcon}</span>
            )}
            <CreditCard className="h-4 w-4 text-slate-400" />
          </div>
        </div>
        {errors.cardNumber && <p className="mt-1 text-xs text-red-400">{errors.cardNumber}</p>}
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <input
            type="text"
            placeholder="MM/YY"
            value={expiry}
            onChange={(e) => {
              let val = e.target.value.replace(/\D/g, '').slice(0, 4)
              if (val.length > 2) val = val.slice(0, 2) + '/' + val.slice(2)
              setExpiry(val)
              setErrors(prev => ({ ...prev, expiry: '' }))
            }}
            onBlur={validate}
            className="w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
          {errors.expiry && <p className="mt-1 text-xs text-red-400">{errors.expiry}</p>}
        </div>
        <div>
          <div className="relative">
            <input
              type="text"
              inputMode="numeric"
              placeholder="CVC/CVV"
              value={cvv}
              onChange={(e) => {
                setCvv(e.target.value.replace(/\D/g, '').slice(0, 4))
                setErrors(prev => ({ ...prev, cvv: '' }))
              }}
              onBlur={validate}
              className="w-full rounded-lg border border-slate-600 bg-slate-700 pl-4 pr-10 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
            />
            <Lock className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          </div>
          {errors.cvv && <p className="mt-1 text-xs text-red-400">{errors.cvv}</p>}
        </div>
      </div>

      <div>
        <input
          type="text"
          placeholder="Card Holder Name"
          value={cardholderName}
          onChange={(e) => {
            setCardholderName(e.target.value)
            setErrors(prev => ({ ...prev, cardholderName: '' }))
          }}
          onBlur={validate}
          className="w-full rounded-lg border border-slate-600 bg-slate-700 px-4 py-2.5 text-white placeholder-slate-400 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        />
        {errors.cardholderName && <p className="mt-1 text-xs text-red-400">{errors.cardholderName}</p>}
      </div>

      <label className="flex items-center gap-2 cursor-pointer">
        <input
          type="checkbox"
          checked={remember}
          onChange={(e) => setRemember(e.target.checked)}
          className="rounded border-slate-600 bg-slate-700 text-primary focus:ring-primary"
        />
        <span className="text-xs text-slate-400">Remember this card</span>
      </label>
    </div>
  )
}
