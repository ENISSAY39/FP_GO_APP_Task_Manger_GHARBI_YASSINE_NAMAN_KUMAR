/**
 * Due-date reminders.
 *
 * A reminder is never stored: it's derived from a task's `due_date` and
 * `status` at read time, against the current clock. See CONTEXT.md for the
 * Overdue / Due soon definitions this mirrors.
 */

const DUE_SOON_WINDOW_MS = 24 * 60 * 60 * 1000

export const REMINDER_LABELS = {
  overdue: 'Overdue',
  due_soon: 'Due soon',
}

/** 'overdue' | 'due_soon' | null for a single task, as of `now`. */
export function reminderStateOf(task, now = new Date()) {
  if (!task?.due_date || task.status === 'DONE') return null

  const due = new Date(task.due_date)
  if (Number.isNaN(due.getTime())) return null

  const msUntilDue = due.getTime() - now.getTime()
  if (msUntilDue < 0) return 'overdue'
  if (msUntilDue <= DUE_SOON_WINDOW_MS) return 'due_soon'
  return null
}

export function isAssignedTo(task, userId) {
  if (!userId) return false
  return (task?.assignees ?? []).some((assignee) => assignee.user_id === userId)
}

/**
 * An Orphan: a task that needs attention and that nobody is on the hook for.
 *
 * Both halves are required — an unassigned task that is not yet due is not an
 * Orphan, and neither is a late task somebody owns. See the CONTEXT.md entry,
 * which reserves "unassigned" for the weaker, assignee-only sense.
 *
 * Surfaced to the project Owner on the project card only; the navbar stays
 * personal (ADR 0003).
 */
export function isOrphan(task, now = new Date()) {
  if ((task?.assignees ?? []).length > 0) return false
  return reminderStateOf(task, now) !== null
}

/**
 * Every task across `projects` that is assigned to `userId` and currently
 * has a reminder, soonest-due first with overdue tasks ahead of due-soon.
 */
export function collectReminders(projects, userId, now = new Date()) {
  const reminders = []

  for (const project of projects ?? []) {
    for (const task of project.tasks ?? []) {
      if (!isAssignedTo(task, userId)) continue
      const state = reminderStateOf(task, now)
      if (!state) continue
      reminders.push({ task, project, state })
    }
  }

  return reminders.sort((a, b) => {
    if (a.state !== b.state) return a.state === 'overdue' ? -1 : 1
    return new Date(a.task.due_date).getTime() - new Date(b.task.due_date).getTime()
  })
}
