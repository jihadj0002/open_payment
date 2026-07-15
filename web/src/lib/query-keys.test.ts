import { describe, it, expect } from 'vitest'
import { queryKeys } from './query-keys'

describe('queryKeys', () => {
  it('payments.all returns correct key', () => {
    expect(queryKeys.payments.all).toEqual(['payments'])
  })

  it('payments.list returns key with params', () => {
    expect(queryKeys.payments.list({ page: '1' })).toEqual(['payments', 'list', { page: '1' }])
  })

  it('payments.detail returns key with id', () => {
    expect(queryKeys.payments.detail('p123')).toEqual(['payments', 'detail', 'p123'])
  })

  it('balance.all returns correct key', () => {
    expect(queryKeys.balance.all).toEqual(['balance'])
  })

  it('balance.transactions returns key with params', () => {
    expect(queryKeys.balance.transactions({ per_page: '20' })).toEqual(['balance', 'transactions', { per_page: '20' }])
  })

  it('settlements.list returns key with params', () => {
    expect(queryKeys.settlements.list({ page: '1' })).toEqual(['settlements', 'list', { page: '1' }])
  })

  it('apiKeys.all returns correct key', () => {
    expect(queryKeys.apiKeys.all).toEqual(['api-keys'])
  })

  it('webhooks.all returns correct key', () => {
    expect(queryKeys.webhooks.all).toEqual(['webhooks'])
  })

  it('customers.list returns correct key', () => {
    expect(queryKeys.customers.list()).toEqual(['customers', 'list'])
  })

  it('customers.detail returns key with id', () => {
    expect(queryKeys.customers.detail('c1')).toEqual(['customers', 'detail', 'c1'])
  })

  it('merchant.profile returns correct key', () => {
    expect(queryKeys.merchant.profile).toEqual(['merchant', 'profile'])
  })

  it('reports.all returns correct key', () => {
    expect(queryKeys.reports.all).toEqual(['reports'])
  })
})
