'use client'

import { usePathname } from 'next/navigation'

const steps = [
  { key: 'business', label: 'Business Info', path: '/onboarding' },
  { key: 'personal', label: 'Personal Info', path: '/onboarding#personal' },
  { key: 'documents', label: 'Documents', path: '/onboarding#documents' },
  { key: 'review', label: 'Review', path: '/onboarding#review' },
]

export default function OnboardingLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()

  const currentStep = 1

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-2xl mx-auto px-4">
        <div className="mb-8">
          <div className="flex items-center justify-between">
            {steps.map((step, index) => (
              <div key={step.key} className="flex items-center">
                <div className={`flex items-center justify-center w-8 h-8 rounded-full text-sm font-medium ${
                  index + 1 <= currentStep
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-200 text-gray-500'
                }`}>
                  {index + 1}
                </div>
                <span className={`ml-2 text-sm ${index + 1 <= currentStep ? 'text-blue-600 font-medium' : 'text-gray-400'}`}>
                  {step.label}
                </span>
                {index < steps.length - 1 && (
                  <div className={`mx-4 w-12 h-0.5 ${index + 1 < currentStep ? 'bg-blue-600' : 'bg-gray-200'}`} />
                )}
              </div>
            ))}
          </div>
        </div>
        {children}
      </div>
    </div>
  )
}
