import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import RemindersMenu from './RemindersMenu.jsx'
import { api } from '../lib/api.js'

// The network is stubbed at the API client, not at global fetch: these tests
// state what the app asked for, not how the client happens to be built.
vi.mock('../lib/api.js', () => ({
  api: { get: vi.fn() },
}))

const NOW = new Date('2026-03-15T12:00:00.000Z')

const HOUR = 60 * 60 * 1000
const DAY = 24 * HOUR

const dueIn = (ms) => new Date(NOW.getTime() + ms).toISOString()

const ME = 7
const SOMEONE_ELSE = 8

const signIn = (user = { id: ME, name: 'Yassine' }) =>
  window.localStorage.setItem('taskmanager.user', JSON.stringify(user))

const task = (overrides = {}) => ({
  id: 1,
  title: 'Task',
  status: 'TODO',
  due_date: dueIn(-DAY),
  assignees: [{ user_id: ME }],
  ...overrides,
})

/** Resolve the projects call with these projects. */
const respondWith = (projects) => api.get.mockResolvedValue({ projects })

const indicator = () => screen.getByRole('button', { name: /^Reminders/ })

/** Open the dropdown and hand back its list rows. */
const openMenu = async () => {
  await userEvent.click(indicator())
  return screen.queryAllByRole('listitem')
}

beforeEach(() => {
  // shouldAdvanceTime keeps Testing Library's waiting helpers working while the
  // clock stays pinned: the component reads the real clock through
  // reminderStateOf's default argument, so it has to be frozen from outside.
  vi.useFakeTimers({ shouldAdvanceTime: true })
  vi.setSystemTime(NOW)
  signIn()
})

afterEach(() => {
  vi.useRealTimers()
  vi.clearAllMocks()
  window.localStorage.clear()
})

describe('RemindersMenu count', () => {
  it('shows no count when the viewer has nothing due', async () => {
    respondWith([{ id: 1, name: 'Apollo', tasks: [task({ due_date: dueIn(30 * DAY) })] }])

    render(<RemindersMenu />)

    await waitFor(() => expect(api.get).toHaveBeenCalledWith('/projects'))
    expect(indicator()).toHaveAccessibleName('Reminders')
    expect(indicator()).not.toHaveTextContent(/\d/)
  })

  it('shows the total across every project, and names it for a screen reader', async () => {
    respondWith([
      { id: 1, name: 'Apollo', tasks: [task({ id: 1, due_date: dueIn(-DAY) })] },
      {
        id: 2,
        name: 'Borealis',
        tasks: [
          task({ id: 2, due_date: dueIn(2 * HOUR) }),
          task({ id: 3, due_date: dueIn(-3 * DAY) }),
        ],
      },
    ])

    render(<RemindersMenu />)

    await waitFor(() => expect(indicator()).toHaveAccessibleName('Reminders, 3 due'))
    expect(indicator()).toHaveTextContent('3')
  })

  it('counts only the viewer’s own tasks', async () => {
    respondWith([
      {
        id: 1,
        name: 'Apollo',
        tasks: [
          task({ id: 1, due_date: dueIn(-DAY) }),
          task({ id: 2, due_date: dueIn(-DAY), assignees: [{ user_id: SOMEONE_ELSE }] }),
        ],
      },
    ])

    render(<RemindersMenu />)

    await waitFor(() => expect(indicator()).toHaveAccessibleName('Reminders, 1 due'))
  })
})

describe('RemindersMenu dropdown', () => {
  it('says plainly that nothing is due rather than looking like a failure', async () => {
    respondWith([])

    render(<RemindersMenu />)
    await waitFor(() => expect(api.get).toHaveBeenCalled())
    await openMenu()

    expect(screen.getByText('Nothing due soon.')).toBeInTheDocument()
  })

  it('lists overdue tasks first, then due soon, soonest first within each group', async () => {
    respondWith([
      {
        id: 1,
        name: 'Apollo',
        tasks: [
          task({ id: 1, title: 'Due in two hours', due_date: dueIn(2 * HOUR) }),
          task({ id: 2, title: 'Overdue by an hour', due_date: dueIn(-HOUR) }),
          task({ id: 3, title: 'Due in an hour', due_date: dueIn(HOUR) }),
          task({ id: 4, title: 'Overdue by a week', due_date: dueIn(-7 * DAY) }),
        ],
      },
    ])

    render(<RemindersMenu />)
    await waitFor(() => expect(indicator()).toHaveAccessibleName('Reminders, 4 due'))
    const rows = await openMenu()

    expect(rows.map((row) => within(row).getByRole('link').textContent)).toEqual([
      expect.stringContaining('Overdue by a week'),
      expect.stringContaining('Overdue by an hour'),
      expect.stringContaining('Due in an hour'),
      expect.stringContaining('Due in two hours'),
    ])
  })

  it('names the task and its project on each row, and links to that project', async () => {
    respondWith([
      { id: 42, name: 'Borealis', tasks: [task({ id: 1, title: 'Write the report' })] },
    ])

    render(<RemindersMenu />)
    await waitFor(() => expect(indicator()).toHaveAccessibleName('Reminders, 1 due'))
    const [row] = await openMenu()

    expect(within(row).getByText('Write the report')).toBeInTheDocument()
    expect(within(row).getByText('Borealis')).toBeInTheDocument()
    expect(within(row).getByRole('link')).toHaveAttribute('href', '/project.html?id=42')
  })

  it('carries an Overdue or Due soon badge on every row', async () => {
    respondWith([
      {
        id: 1,
        name: 'Apollo',
        tasks: [
          task({ id: 1, title: 'Late', due_date: dueIn(-DAY) }),
          task({ id: 2, title: 'Soon', due_date: dueIn(HOUR) }),
        ],
      },
    ])

    render(<RemindersMenu />)
    await waitFor(() => expect(indicator()).toHaveAccessibleName('Reminders, 2 due'))
    const [late, soon] = await openMenu()

    expect(within(late).getByText('Overdue')).toBeInTheDocument()
    expect(within(soon).getByText('Due soon')).toBeInTheDocument()
  })

  it('closes when the viewer clicks outside it', async () => {
    respondWith([])

    render(
      <div>
        <span data-testid="elsewhere">elsewhere</span>
        <RemindersMenu />
      </div>,
    )
    await waitFor(() => expect(api.get).toHaveBeenCalled())
    await openMenu()
    expect(screen.getByText('Nothing due soon.')).toBeInTheDocument()

    await userEvent.click(screen.getByTestId('elsewhere'))

    expect(screen.queryByText('Nothing due soon.')).not.toBeInTheDocument()
  })
})

describe('RemindersMenu when things go wrong', () => {
  it('stays silently empty when the request fails, rather than throwing', async () => {
    api.get.mockRejectedValue(new Error('network down'))

    render(<RemindersMenu />)

    await waitFor(() => expect(api.get).toHaveBeenCalled())
    expect(indicator()).toHaveAccessibleName('Reminders')
    await openMenu()
    expect(screen.getByText('Nothing due soon.')).toBeInTheDocument()
  })

  it('renders nothing at all for a signed-out visitor, and asks the API for nothing', () => {
    window.localStorage.clear()

    const { container } = render(<RemindersMenu />)

    expect(container).toBeEmptyDOMElement()
    expect(api.get).not.toHaveBeenCalled()
  })
})

describe('RemindersMenu stays personal', () => {
  // ADR 0003 routes Orphan work to the project Owner on the project card and
  // deliberately keeps it out of here: the navbar answers "what do I owe?",
  // and a count mixing that with work nobody owns stops being actionable.
  it('ignores overdue tasks nobody is assigned to, even for the project owner', async () => {
    respondWith([
      {
        id: 1,
        name: 'Apollo',
        owner_id: ME,
        tasks: [
          task({ id: 1, title: 'Mine and late', due_date: dueIn(-DAY) }),
          task({ id: 2, title: 'Nobody owns this', due_date: dueIn(-5 * DAY), assignees: [] }),
        ],
      },
    ])

    render(<RemindersMenu />)

    await waitFor(() => expect(indicator()).toHaveAccessibleName('Reminders, 1 due'))
    const rows = await openMenu()
    expect(rows).toHaveLength(1)
    expect(screen.queryByText('Nobody owns this')).not.toBeInTheDocument()
  })
})
