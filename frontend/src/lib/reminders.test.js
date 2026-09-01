import { describe, expect, it } from 'vitest'
import { collectReminders, isAssignedTo, reminderStateOf } from './reminders.js'

// Every test freezes the clock: a reminder is a comparison against "now", so a
// suite that read the real clock would pass or fail depending on the day it ran.
const NOW = new Date('2026-03-15T12:00:00.000Z')

const HOUR = 60 * 60 * 1000
const DAY = 24 * HOUR

/** A due date `ms` from the frozen NOW — negative for the past. */
const dueIn = (ms) => new Date(NOW.getTime() + ms).toISOString()

const task = (overrides = {}) => ({
  id: 1,
  title: 'Task',
  status: 'TODO',
  due_date: dueIn(-DAY),
  assignees: [],
  ...overrides,
})

describe('reminderStateOf', () => {
  it('is overdue when the due date is comfortably in the past', () => {
    expect(reminderStateOf(task({ due_date: dueIn(-30 * DAY) }), NOW)).toBe('overdue')
  })

  it('is overdue one millisecond after the due moment', () => {
    expect(reminderStateOf(task({ due_date: dueIn(-1) }), NOW)).toBe('overdue')
  })

  it('is due soon exactly at the due moment, not yet overdue', () => {
    expect(reminderStateOf(task({ due_date: dueIn(0) }), NOW)).toBe('due_soon')
  })

  it('is due soon just inside the 24-hour window', () => {
    expect(reminderStateOf(task({ due_date: dueIn(23 * HOUR) }), NOW)).toBe('due_soon')
  })

  it('is due soon exactly at 24 hours', () => {
    expect(reminderStateOf(task({ due_date: dueIn(DAY) }), NOW)).toBe('due_soon')
  })

  it('carries no reminder one millisecond beyond the window', () => {
    expect(reminderStateOf(task({ due_date: dueIn(DAY + 1) }), NOW)).toBeNull()
  })

  it('carries no reminder for a task due far in the future', () => {
    expect(reminderStateOf(task({ due_date: dueIn(30 * DAY) }), NOW)).toBeNull()
  })

  describe('returns null', () => {
    it('for a task with no due date', () => {
      expect(reminderStateOf(task({ due_date: null }), NOW)).toBeNull()
    })

    it('for a DONE task inside the due-soon window', () => {
      expect(reminderStateOf(task({ status: 'DONE', due_date: dueIn(HOUR) }), NOW)).toBeNull()
    })

    it('for a DONE task whose due date is long past', () => {
      expect(reminderStateOf(task({ status: 'DONE', due_date: dueIn(-90 * DAY) }), NOW)).toBeNull()
    })

    it('for an unparseable due date, rather than throwing', () => {
      expect(reminderStateOf(task({ due_date: 'not a date' }), NOW)).toBeNull()
    })

    it('for a missing task', () => {
      expect(reminderStateOf(undefined, NOW)).toBeNull()
      expect(reminderStateOf(null, NOW)).toBeNull()
    })
  })
})

describe('isAssignedTo', () => {
  const assigned = task({ assignees: [{ user_id: 7 }, { user_id: 9 }] })

  it('is true when the user is among the assignees', () => {
    expect(isAssignedTo(assigned, 9)).toBe(true)
  })

  it('is false for a user who is not an assignee', () => {
    expect(isAssignedTo(assigned, 42)).toBe(false)
  })

  it('is false when the task has no assignees', () => {
    expect(isAssignedTo(task({ assignees: [] }), 7)).toBe(false)
  })

  it('is false when there is no user id', () => {
    expect(isAssignedTo(assigned, undefined)).toBe(false)
    expect(isAssignedTo(assigned, null)).toBe(false)
  })
})

describe('collectReminders', () => {
  const ME = 7
  const SOMEONE_ELSE = 8

  const mine = (overrides) => task({ assignees: [{ user_id: ME }], ...overrides })
  const theirs = (overrides) => task({ assignees: [{ user_id: SOMEONE_ELSE }], ...overrides })

  it('leaves out tasks assigned to someone else', () => {
    const projects = [
      { id: 1, name: 'Alpha', tasks: [theirs({ id: 1, due_date: dueIn(-DAY) })] },
    ]

    expect(collectReminders(projects, ME, NOW)).toEqual([])
  })

  it('leaves out the viewer’s own tasks that carry no reminder', () => {
    const projects = [
      {
        id: 1,
        name: 'Alpha',
        tasks: [
          mine({ id: 1, due_date: dueIn(30 * DAY) }),
          mine({ id: 2, status: 'DONE', due_date: dueIn(-DAY) }),
          mine({ id: 3, due_date: null }),
        ],
      },
    ]

    expect(collectReminders(projects, ME, NOW)).toEqual([])
  })

  it('spans every project the viewer belongs to', () => {
    const projects = [
      { id: 1, name: 'Alpha', tasks: [mine({ id: 1, due_date: dueIn(HOUR) })] },
      { id: 2, name: 'Beta', tasks: [mine({ id: 2, due_date: dueIn(2 * HOUR) })] },
    ]

    expect(collectReminders(projects, ME, NOW).map((r) => r.task.id)).toEqual([1, 2])
  })

  it('carries the owning project through on each entry', () => {
    const projects = [{ id: 4, name: 'Alpha', tasks: [mine({ id: 1, due_date: dueIn(-DAY) })] }]

    const [reminder] = collectReminders(projects, ME, NOW)
    expect(reminder.project.id).toBe(4)
    expect(reminder.project.name).toBe('Alpha')
    expect(reminder.state).toBe('overdue')
  })

  it('orders every overdue task ahead of every due-soon one', () => {
    const projects = [
      {
        id: 1,
        name: 'Alpha',
        tasks: [
          mine({ id: 1, due_date: dueIn(HOUR) }),
          mine({ id: 2, due_date: dueIn(-HOUR) }),
          mine({ id: 3, due_date: dueIn(2 * HOUR) }),
          mine({ id: 4, due_date: dueIn(-2 * HOUR) }),
        ],
      },
    ]

    expect(collectReminders(projects, ME, NOW).map((r) => r.state)).toEqual([
      'overdue',
      'overdue',
      'due_soon',
      'due_soon',
    ])
  })

  it('orders by due date within each group, soonest first', () => {
    const projects = [
      {
        id: 1,
        name: 'Alpha',
        tasks: [
          mine({ id: 1, due_date: dueIn(-HOUR) }),
          mine({ id: 2, due_date: dueIn(2 * HOUR) }),
          mine({ id: 3, due_date: dueIn(-5 * DAY) }),
          mine({ id: 4, due_date: dueIn(HOUR) }),
        ],
      },
    ]

    expect(collectReminders(projects, ME, NOW).map((r) => r.task.id)).toEqual([3, 1, 4, 2])
  })

  it('returns an empty array for a project with no tasks', () => {
    expect(collectReminders([{ id: 1, name: 'Empty' }], ME, NOW)).toEqual([])
    expect(collectReminders([{ id: 1, name: 'Empty', tasks: [] }], ME, NOW)).toEqual([])
  })

  it('returns an empty array for a null or missing project list', () => {
    expect(collectReminders(null, ME, NOW)).toEqual([])
    expect(collectReminders(undefined, ME, NOW)).toEqual([])
    expect(collectReminders([], ME, NOW)).toEqual([])
  })

  it('returns an empty array when there is no viewer', () => {
    const projects = [{ id: 1, name: 'Alpha', tasks: [mine({ id: 1, due_date: dueIn(-DAY) })] }]

    expect(collectReminders(projects, undefined, NOW)).toEqual([])
  })
})
