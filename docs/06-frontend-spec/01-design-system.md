# Design System

> **Status:** ✅ Updated 2026-07-16
> **Code Ref:** `web/src/components/ui/`

## Tech Stack
- **Framework:** Next.js 14+ (App Router)
- **Language:** TypeScript (strict mode)
- **Styling:** Tailwind CSS 3+
- **UI Components:** Headless UI (Radix primitives) + custom components
- **Icons:** Lucide React
- **Charts:** Apache ECharts (via echarts-for-react)
- **Forms:** React Hook Form + Zod validation

## Color Palette

| Token | Hex | Usage |
|-------|-----|-------|
| `primary-50` | #EFF6FF | Backgrounds, hover states |
| `primary-100` | #DBEAFE | Light backgrounds |
| `primary-500` | #3B82F6 | Primary buttons, links |
| `primary-600` | #2563EB | Button hover, active states |
| `primary-700` | #1D4ED8 | Focus rings |

| Token | Hex | Usage |
|-------|-----|-------|
| `success-500` | #10B981 | Payment success, green indicators |
| `warning-500` | #F59E0B | Pending states, warnings |
| `danger-500` | #EF4444 | Errors, failed payments, alerts |
| `info-500` | #6366F1 | Info banners, processing states |

| Token | Hex | Usage |
|-------|-----|-------|
| `neutral-50` | #F9FAFB | Page backgrounds |
| `neutral-100` | #F3F4F6 | Card backgrounds, table rows |
| `neutral-200` | #E5E7EB | Borders, dividers |
| `neutral-700` | #374151 | Secondary text |
| `neutral-900` | #111827 | Primary text |

## Typography

| Element | Font Size | Weight | Line Height |
|---------|-----------|--------|-------------|
| h1 (page title) | text-3xl (30px) | font-bold (700) | leading-tight |
| h2 (section title) | text-xl (20px) | font-semibold (600) | leading-snug |
| h3 (card title) | text-base (16px) | font-semibold (600) | leading-normal |
| body | text-sm (14px) | font-normal (400) | leading-normal |
| body-small | text-xs (12px) | font-normal (400) | leading-normal |
| caption | text-xs (12px) | font-medium (500) | leading-normal |
| code | text-sm (14px) | font-mono | leading-normal |

Font family: Inter (sans-serif), JetBrains Mono (monospace for code).

## Spacing Scale
Based on Tailwind's default scale: `0, 0.25, 0.5, 1, 1.5, 2, 2.5, 3, 4, 5, 6, 8, 10, 12, 16, 20, 24` rem.

## Component Library

### Button Variants
```tsx
<Button variant="primary" size="sm" disabled>
  Pay Now
</Button>
```
| Variant | Use | Status |
|---------|-----|--------|
| primary | Primary actions (Pay, Save, Create) | ✅ Implemented |
| secondary | Alternative actions | ✅ Implemented |
| outline | Less emphasis actions | ✅ Implemented |
| danger | Destructive actions (Delete, Suspend) | ✅ Implemented |
| ghost | Minimal, in tables/lists | ✅ Implemented |
| link | Inline navigation | ✅ Implemented |

> **Note:** Current `Button` component uses flat props (`variant`, `size`, `disabled`, `children`). `loading` and `icon` props are **TODO**.

### Status Badges
```tsx
<Badge variant="success">Succeeded</Badge>
<Badge variant="warning">Pending</Badge>
<Badge variant="danger">Failed</Badge>
<Badge variant="info">Processing</Badge>
<Badge variant="neutral">Created</Badge>
```

### Data Table
```tsx
<DataTable
  columns={columns}
  data={transactions}
  pagination={{ pageSize: 25, total: 1042 }}
  onSort={(field, dir) => fetchSorted(field, dir)}
  onFilter={(filters) => fetchFiltered(filters)}
  loading={isLoading}
/>
```

### Card
> **Note:** Current `Card` component uses flat props, not compound component pattern.

```tsx
<Card>
  {/* content */}
</Card>
```

### Stat Card
```tsx
<StatCard
  title="Today's Revenue"
  value="BDT 125,430"
  change="+12.5%"
  changeType="positive"
  icon={<DollarSign />}
/>
```

### Modal/Dialog
> **Note:** Current `Modal` component uses flat props, not compound component pattern.

```tsx
<Modal open={isOpen} onClose={closeModal}>
  <Modal.Header>Confirm Refund</Modal.Header>
  <Modal.Body>
    Are you sure you want to refund BDT 500.00?
  </Modal.Body>
  <Modal.Footer>
    <Button variant="secondary" onClick={closeModal}>Cancel</Button>
    <Button variant="danger" onClick={processRefund}>Confirm Refund</Button>
  </Modal.Footer>
</Modal>
```

### Toast/Notification
```tsx
toast.success('Payment refunded successfully');
toast.error('Refund failed: insufficient balance');
```

## Responsive Breakpoints
| Breakpoint | Width | Layout |
|------------|-------|--------|
| xs | <640px | Single column, stacked nav |
| sm | 640px | Single column, sidebar hidden |
| md | 768px | Two column, sidebar visible |
| lg | 1024px | Full dashboard layout |
| xl | 1280px | Max width container |

## Dark Mode
- Implemented via Tailwind `dark:` variant
- Toggle stored in user preferences (localStorage)
- System preference detection for default
- Dark mode inverts card backgrounds, adjusts text contrast
