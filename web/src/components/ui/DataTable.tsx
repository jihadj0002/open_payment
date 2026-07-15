'use client'

import { useCallback } from 'react'
import { clsx } from 'clsx'
import { ArrowUp, ArrowDown, ArrowUpDown } from 'lucide-react'
import { Skeleton } from './Skeleton'

interface Column<T> {
  key: string
  header: string
  render?: (item: T) => React.ReactNode
  sortable?: boolean
  className?: string
}

interface DataTableProps<T> {
  columns: Column<T>[]
  data: T[]
  loading?: boolean
  emptyMessage?: string
  pagination?: {
    page: number
    pageSize: number
    total: number
    onPageChange: (page: number) => void
  }
  sort?: {
    key: string
    direction: 'asc' | 'desc'
    onSort: (key: string, direction: 'asc' | 'desc') => void
  }
}

function DataTable<T>({
  columns,
  data,
  loading,
  emptyMessage = 'No data found',
  pagination,
  sort,
}: DataTableProps<T>) {
  const totalPages = pagination ? Math.ceil(pagination.total / pagination.pageSize) : 0

  const handleSort = useCallback(
    (key: string) => {
      if (!sort) return
      if (sort.key === key) {
        sort.onSort(key, sort.direction === 'asc' ? 'desc' : 'asc')
      } else {
        sort.onSort(key, 'asc')
      }
    },
    [sort]
  )

  const renderSortIcon = (column: Column<T>) => {
    if (!column.sortable || !sort) return null
    if (sort.key !== column.key) {
      return <ArrowUpDown className="ml-1 h-3 w-3 text-slate-500" />
    }
    return sort.direction === 'asc' ? (
      <ArrowUp className="ml-1 h-3 w-3" />
    ) : (
      <ArrowDown className="ml-1 h-3 w-3" />
    )
  }

  return (
    <div className="overflow-x-auto rounded-lg border border-slate-700">
      <table className="min-w-full divide-y divide-slate-700">
        <thead className="bg-slate-800">
          <tr>
            {columns.map((column) => (
              <th
                key={column.key}
                className={clsx(
                  'px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-400',
                  column.sortable && sort && 'cursor-pointer select-none hover:text-slate-200',
                  column.className
                )}
                onClick={() => column.sortable && handleSort(column.key)}
              >
                <span className="inline-flex items-center">
                  {column.header}
                  {renderSortIcon(column)}
                </span>
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-700 bg-slate-800/50">
          {loading ? (
            Array.from({ length: 5 }).map((_, rowIdx) => (
              <tr key={rowIdx}>
                {columns.map((column) => (
                  <td key={column.key} className="px-4 py-2">
                    <Skeleton className="h-4 w-full" />
                  </td>
                ))}
              </tr>
            ))
          ) : data.length === 0 ? (
            <tr>
              <td
                colSpan={columns.length}
                className="px-4 py-12 text-center text-sm text-slate-400"
              >
                {emptyMessage}
              </td>
            </tr>
          ) : (
            data.map((item, rowIdx) => (
              <tr key={rowIdx} className="transition-colors hover:bg-slate-800">
                {columns.map((column) => (
                  <td
                    key={column.key}
                    className={clsx('whitespace-nowrap px-4 py-3 text-sm text-slate-300', column.className)}
                  >
                      {column.render
                        ? column.render(item)
                        : ((item as Record<string, unknown>)[column.key] as React.ReactNode) ?? '-'}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>

      {pagination && totalPages > 1 && (
        <div className="flex items-center justify-between border-t border-slate-700 px-4 py-3">
          <p className="text-sm text-slate-400">
            Showing {(pagination.page - 1) * pagination.pageSize + 1}-
            {Math.min(pagination.page * pagination.pageSize, pagination.total)} of{' '}
            {pagination.total}
          </p>
          <div className="flex items-center gap-1">
            <button
              onClick={() => pagination.onPageChange(pagination.page - 1)}
              disabled={pagination.page <= 1}
              className="rounded-md px-2 py-1 text-sm text-slate-400 transition-colors hover:bg-slate-700 hover:text-white disabled:pointer-events-none disabled:opacity-40"
            >
              Prev
            </button>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
              <button
                key={page}
                onClick={() => pagination.onPageChange(page)}
                className={clsx(
                  'min-w-[2rem] rounded-md px-2 py-1 text-sm transition-colors',
                  page === pagination.page
                    ? 'bg-primary text-white'
                    : 'text-slate-400 hover:bg-slate-700 hover:text-white'
                )}
              >
                {page}
              </button>
            ))}
            <button
              onClick={() => pagination.onPageChange(pagination.page + 1)}
              disabled={pagination.page >= totalPages}
              className="rounded-md px-2 py-1 text-sm text-slate-400 transition-colors hover:bg-slate-700 hover:text-white disabled:pointer-events-none disabled:opacity-40"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

export { DataTable }
export type { Column, DataTableProps }
