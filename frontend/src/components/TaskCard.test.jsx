import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import TaskCard from './TaskCard.jsx'

// TaskCard calls reminderStateOf(task) without a clock argument, so the only
// seam for time here is the system clock itself. Freezing it keeps these tests
// from depending on the day they run.
const NOW = new Date('2026-03-15T12:00:00.000Z')

const HOUR = 60 * 60 * 1000
const DAY = 24 * HOUR

const dueIn = (ms) => new Date(NOW.getTime() + ms).toISOString()

const renderCard = (task, props = {}) =>
  render(
    <TaskCard
      task={{
        id: 1,
        title: 'Write the report',
        description: '',
        status: 'TODO',
        priority: 'MEDIUM',
        due_date: dueIn(-DAY),
        assignees: [],
        created_at: dueIn(-30 * DAY),
        ...task,
      }}
      members={[]}
      canEdit={false}
      onEdit={() => {}}
      onDelete={() => {}}
      onAssign={() => {}}
      onUnassign={() => {}}
      onStatusChange={() => {}}
      {...props}
    />,
  )

/**
 * The element holding the "Due <date>" metadata. formatDate always leads with a
 * two-digit day, which is what keeps this from also matching the "Due soon"
 * badge sitting a few nodes away.
 */
const dueDateText = () => screen.getByText(/^Due \d/)

beforeEach(() => {
  vi.useFakeTimers()
  vi.setSystemTime(NOW)
})

afterEach(() => {
  vi.useRealTimers()
})

describe('TaskCard reminder badge', () => {
  it('flags an overdue task', () => {
    renderCard({ due_date: dueIn(-2 * DAY) })

    expect(screen.getByText('Overdue')).toBeInTheDocument()
    expect(screen.queryByText('Due soon')).not.toBeInTheDocument()
  })

  it('flags a task due within the window', () => {
    renderCard({ due_date: dueIn(2 * HOUR) })

    expect(screen.getByText('Due soon')).toBeInTheDocument()
    expect(screen.queryByText('Overdue')).not.toBeInTheDocument()
  })

  it('flags nothing for a DONE task whose due date has passed', () => {
    renderCard({ status: 'DONE', due_date: dueIn(-2 * DAY) })

    expect(screen.queryByText('Overdue')).not.toBeInTheDocument()
    expect(screen.queryByText('Due soon')).not.toBeInTheDocument()
  })

  it('flags nothing for a task due far in the future', () => {
    renderCard({ due_date: dueIn(30 * DAY) })

    expect(screen.queryByText('Overdue')).not.toBeInTheDocument()
    expect(screen.queryByText('Due soon')).not.toBeInTheDocument()
  })

  it('flags nothing for a task with no due date', () => {
    renderCard({ due_date: null })

    expect(screen.queryByText('Overdue')).not.toBeInTheDocument()
    expect(screen.queryByText('Due soon')).not.toBeInTheDocument()
    expect(screen.queryByText(/^Due \d/)).not.toBeInTheDocument()
  })

  it('flags a task assigned to someone else: the badge is task-scoped, not viewer-scoped', () => {
    renderCard({ due_date: dueIn(-2 * DAY), assignees: [{ user_id: 99, user: { name: 'Someone' } }] })

    expect(screen.getByText('Overdue')).toBeInTheDocument()
  })
})

describe('TaskCard due date colour', () => {
  // The colour is the only observable difference here, so these bind to the
  // Tailwind classes. Kept to the three cases that carry meaning.
  it('renders the due date in the overdue colour', () => {
    renderCard({ due_date: dueIn(-2 * DAY) })

    expect(dueDateText()).toHaveClass('text-rose-400')
  })

  it('renders the due date in the due-soon colour', () => {
    renderCard({ due_date: dueIn(2 * HOUR) })

    expect(dueDateText()).toHaveClass('text-amber-400')
  })

  it('renders an unremarkable due date in the neutral colour', () => {
    renderCard({ due_date: dueIn(30 * DAY) })

    expect(dueDateText()).toHaveClass('text-slate-400')
    expect(dueDateText()).not.toHaveClass('text-rose-400')
    expect(dueDateText()).not.toHaveClass('text-amber-400')
  })
})
