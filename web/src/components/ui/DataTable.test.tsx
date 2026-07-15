import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { DataTable, Column } from './DataTable'

interface Item {
  id: number
  name: string
}

const columns: Column<Item>[] = [
  { key: 'id', header: 'ID' },
  { key: 'name', header: 'Name' },
]

const data: Item[] = [
  { id: 1, name: 'Alice' },
  { id: 2, name: 'Bob' },
]

describe('DataTable', () => {
  it('renders column headers', () => {
    render(<DataTable columns={columns} data={data} />)
    expect(screen.getByText('ID')).toBeInTheDocument()
    expect(screen.getByText('Name')).toBeInTheDocument()
  })

  it('renders data rows', () => {
    render(<DataTable columns={columns} data={data} />)
    expect(screen.getByText('Alice')).toBeInTheDocument()
    expect(screen.getByText('Bob')).toBeInTheDocument()
  })

  it('renders empty state', () => {
    render(<DataTable columns={columns} data={[]} emptyMessage="No items" />)
    expect(screen.getByText('No items')).toBeInTheDocument()
  })

  it('renders loading skeleton', () => {
    const { container } = render(<DataTable columns={columns} data={[]} loading />)
    const skeletons = container.querySelectorAll('.animate-pulse')
    expect(skeletons.length).toBeGreaterThan(0)
  })

  it('renders pagination controls', () => {
    render(
      <DataTable
        columns={columns}
        data={data}
        pagination={{ page: 1, pageSize: 10, total: 25, onPageChange: vi.fn() }}
      />
    )
    expect(screen.getByText('Prev')).toBeInTheDocument()
    expect(screen.getByText('Next')).toBeInTheDocument()
  })

  it('calls onPageChange on page click', async () => {
    const handlePageChange = vi.fn()
    const user = userEvent.setup()
    render(
      <DataTable
        columns={columns}
        data={data}
        pagination={{ page: 1, pageSize: 10, total: 25, onPageChange: handlePageChange }}
      />
    )
    const pageButtons = screen.getAllByText('2')
    const pageBtn = pageButtons.find((btn) => btn.tagName === 'BUTTON')
    if (pageBtn) await user.click(pageBtn)
    expect(handlePageChange).toHaveBeenCalledWith(2)
  })

  it('renders sortable column with sort icon', () => {
    const sortableColumns: Column<Item>[] = [
      { key: 'id', header: 'ID', sortable: true },
      { key: 'name', header: 'Name' },
    ]
    render(
      <DataTable
        columns={sortableColumns}
        data={data}
        sort={{ key: 'id', direction: 'asc', onSort: vi.fn() }}
      />
    )
    expect(screen.getByText('ID').querySelector('svg')).toBeInTheDocument()
  })

  it('sort toggles direction on click', async () => {
    const handleSort = vi.fn()
    const user = userEvent.setup()
    const sortableColumns: Column<Item>[] = [
      { key: 'id', header: 'ID', sortable: true },
    ]
    render(
      <DataTable
        columns={sortableColumns}
        data={data}
        sort={{ key: 'id', direction: 'asc', onSort: handleSort }}
      />
    )
    await user.click(screen.getByText('ID'))
    expect(handleSort).toHaveBeenCalledWith('id', 'desc')
  })
})
