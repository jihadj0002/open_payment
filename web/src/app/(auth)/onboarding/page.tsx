'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Card } from '@/components/ui/Card'
import { Toast } from '@/components/ui/Toast'

const businessSchema = z.object({
  business_name: z.string().min(1, 'Business name is required'),
  business_type: z.string().min(1, 'Business type is required'),
  website: z.string().url('Must be a valid URL').optional().or(z.literal('')),
  tax_id: z.string().optional(),
})

const personalSchema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  phone: z.string().min(10, 'Valid phone number required'),
  address: z.string().min(1, 'Address is required'),
})

type BusinessData = z.infer<typeof businessSchema>
type PersonalData = z.infer<typeof personalSchema>

export default function OnboardingPage() {
  const [step, setStep] = useState<'business' | 'personal' | 'documents' | 'review'>('business')
  const [businessData, setBusinessData] = useState<BusinessData | null>(null)
  const [personalData, setPersonalData] = useState<PersonalData | null>(null)
  const [toast, setToast] = useState<{ type: 'success' | 'error', message: string } | null>(null)

  const businessForm = useForm<BusinessData>({
    resolver: zodResolver(businessSchema),
  })

  const personalForm = useForm<PersonalData>({
    resolver: zodResolver(personalSchema),
  })

  const submitMutation = useMutation({
    mutationFn: (data: BusinessData & PersonalData) =>
      api.post('/merchants/onboarding', data),
    onSuccess: () => {
      setToast({ type: 'success', message: 'Onboarding submitted successfully!' })
    },
    onError: () => {
      setToast({ type: 'error', message: 'Failed to submit onboarding. Please try again.' })
    },
  })

  const handleBusinessSubmit = (data: BusinessData) => {
    setBusinessData(data)
    setStep('personal')
  }

  const handlePersonalSubmit = (data: PersonalData) => {
    setPersonalData(data)
    setStep('documents')
  }

  const handleDocumentsNext = () => {
    setStep('review')
  }

  const handleReviewSubmit = () => {
    if (businessData && personalData) {
      submitMutation.mutate({ ...businessData, ...personalData })
    }
  }

  return (
    <>
      {toast && (
        <Toast
          variant={toast.type}
          message={toast.message}
          onClose={() => setToast(null)}
        />
      )}

      <Card className="p-6">
        {step === 'business' && (
          <form onSubmit={businessForm.handleSubmit(handleBusinessSubmit)}>
            <h2 className="text-xl font-semibold mb-6">Business Information</h2>
            <div className="space-y-4">
              <Input
                label="Business Name"
                {...businessForm.register('business_name')}
                error={businessForm.formState.errors.business_name?.message}
              />
              <Input
                label="Business Type"
                placeholder="e.g., Retail, SaaS, E-commerce"
                {...businessForm.register('business_type')}
                error={businessForm.formState.errors.business_type?.message}
              />
              <Input
                label="Website"
                {...businessForm.register('website')}
                error={businessForm.formState.errors.website?.message}
              />
              <Input
                label="Tax ID (optional)"
                {...businessForm.register('tax_id')}
              />
            </div>
            <div className="mt-6 flex justify-end">
              <Button type="submit">Next: Personal Info</Button>
            </div>
          </form>
        )}

        {step === 'personal' && (
          <form onSubmit={personalForm.handleSubmit(handlePersonalSubmit)}>
            <h2 className="text-xl font-semibold mb-6">Personal Information</h2>
            <div className="space-y-4">
              <Input
                label="First Name"
                {...personalForm.register('first_name')}
                error={personalForm.formState.errors.first_name?.message}
              />
              <Input
                label="Last Name"
                {...personalForm.register('last_name')}
                error={personalForm.formState.errors.last_name?.message}
              />
              <Input
                label="Phone Number"
                {...personalForm.register('phone')}
                error={personalForm.formState.errors.phone?.message}
              />
              <Input
                label="Business Address"
                {...personalForm.register('address')}
                error={personalForm.formState.errors.address?.message}
              />
            </div>
            <div className="mt-6 flex justify-between">
              <Button type="button" variant="outline" onClick={() => setStep('business')}>Back</Button>
              <Button type="submit">Next: Documents</Button>
            </div>
          </form>
        )}

        {step === 'documents' && (
          <div>
            <h2 className="text-xl font-semibold mb-6">Upload Documents</h2>
            <p className="text-gray-600 mb-4">Please upload your identification and business documents for KYC verification.</p>
            <div className="space-y-4">
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
                <p className="text-gray-500">Drag & drop or click to upload</p>
                <p className="text-sm text-gray-400 mt-1">Accepted: PDF, PNG, JPG (max 10MB)</p>
              </div>
            </div>
            <div className="mt-6 flex justify-between">
              <Button type="button" variant="outline" onClick={() => setStep('personal')}>Back</Button>
              <Button onClick={handleDocumentsNext}>Next: Review</Button>
            </div>
          </div>
        )}

        {step === 'review' && (
          <div>
            <h2 className="text-xl font-semibold mb-6">Review Your Information</h2>
            <div className="space-y-4">
              <div>
                <h3 className="font-medium text-gray-700">Business Info</h3>
                <p>{businessData?.business_name} - {businessData?.business_type}</p>
                {businessData?.website && <p className="text-sm text-gray-500">{businessData.website}</p>}
              </div>
              <div>
                <h3 className="font-medium text-gray-700">Personal Info</h3>
                <p>{personalData?.first_name} {personalData?.last_name}</p>
                <p className="text-sm text-gray-500">{personalData?.phone}</p>
              </div>
            </div>
            <div className="mt-6 flex justify-between">
              <Button type="button" variant="outline" onClick={() => setStep('documents')}>Back</Button>
              <Button onClick={handleReviewSubmit} disabled={submitMutation.isPending}>
                {submitMutation.isPending ? 'Submitting...' : 'Submit Application'}
              </Button>
            </div>
          </div>
        )}
      </Card>
    </>
  )
}
