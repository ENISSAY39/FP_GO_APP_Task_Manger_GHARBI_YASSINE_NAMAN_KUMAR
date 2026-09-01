import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProjectCard from './ProjectCard.jsx'

// Like TaskCard, ProjectCard reads the clock through reminderStateOf's default
// argument, so time is frozen at the system level.
const NOW = new Date('2026-03-15T12:00:00.000Z')

const HOUR = 60 * 60 * 1000
const DAY = 24 * HOUR

const dueIn = (ms) => new Date(NOW.getTime() + ms).toISOString()

const ME = 7
const SOMEONE_ELSE = 8

const task = (overrides = {}) => ({
  id: 1,
  title: 'Task',
  status: 'TODO',
  due_date: dueIn(-DAY),
  assignees: [{ user_id: ME }],
  ...overrides,
})

const renderCard = (tasks, props = {}) =>
  render(
    <ProjectCard
      project={{
        id: 1,
        name: 'Apollo',
        description: '',
        members: [{ user_id: ME }],
        tasks,
        created_at: dueIn(-90 * DAY),
      }}
      isOwner={false}
      currentUserId={ME}
      onDelete={() => {}}
      {...props}
    />,
  )

/** The "<n> due" badge, or null when the card shows no count at all. */
const countBadge = () => screen.queryByText(/^\d+ due$/)

beforeEach(() => {
  vi.useFakeTimers()
  vi.setSystemTime(NOW)
})

afterEach(() => {
  vi.useRealTimers()
})

describe('ProjectCard reminder count', () => {
  it('counts the viewer’s overdue and due-soon tasks together', () => {
    renderCard([
      task({ id: 1, due_date: dueIn(-2 * DAY) }),
      task({ id: 2, due_date: dueIn(2 * HOUR) }),
      task({ id: 3, due_date: dueIn(30 * DAY) }),
    ])

    expect(countBadge()).toHaveTextContent('2 due')
  })

  it('renders no count badge when nothing of the viewer’s needs attention', () => {
    renderCard([task({ id: 1, due_date: dueIn(30 * DAY) })])

    expect(countBadge()).toBeNull()
  })

  it('renders no count badge for a project with no tasks', () => {
    renderCard([])

    expect(countBadge()).toBeNull()
  })

  it('leaves out reminder-carrying tasks assigned to someone else', () => {
    renderCard([
      task({ id: 1, due_date: dueIn(-2 * DAY), assignees: [{ user_id: SOMEONE_ELSE }] }),
      task({ id: 2, due_date: dueIn(-3 * DAY), assignees: [{ user_id: ME }] }),
    ])

    expect(countBadge()).toHaveTextContent('1 due')
  })

  it('leaves out unassigned tasks, even overdue ones', () => {
    renderCard([task({ id: 1, due_date: dueIn(-2 * DAY), assignees: [] })])

    expect(countBadge()).toBeNull()
  })

  it('leaves out the viewer’s DONE tasks', () => {
    renderCard([task({ id: 1, status: 'DONE', due_date: dueIn(-2 * DAY) })])

    expect(countBadge()).toBeNull()
  })
})

describe('ProjectCard count alongside the Owner badge', () => {
  it('shows both when the viewer owns a project that needs attention', () => {
    renderCard([task({ id: 1, due_date: dueIn(-2 * DAY) })], { isOwner: true })

    expect(countBadge()).toHaveTextContent('1 due')
    expect(screen.getByText('Owner')).toBeInTheDocument()
  })

  it('still shows the Owner badge when there is no count', () => {
    renderCard([task({ id: 1, due_date: dueIn(30 * DAY) })], { isOwner: true })

    expect(countBadge()).toBeNull()
    expect(screen.getByText('Owner')).toBeInTheDocument()
  })
})
