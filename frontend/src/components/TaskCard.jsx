import { useState } from 'react'
import Button from './ui/Button.jsx'
import { PriorityBadge, ReminderBadge, StatusBadge } from './ui/Badge.jsx'
import { TASK_STATUSES, STATUS_LABELS } from '../lib/constants.js'
import { displayUser, formatDate, initialsOf } from '../lib/format.js'
import { reminderStateOf } from '../lib/reminders.js'

export default function TaskCard({
  task,
  members,
  canEdit,
  onEdit,
  onDelete,
  onAssign,
  onUnassign,
  onStatusChange,
}) {
  const [assigning, setAssigning] = useState(false)

  const reminder = reminderStateOf(task)
  const assignees = task.assignees ?? []
  const assignedIds = new Set(assignees.map((assignee) => assignee.user_id))
  const assignable = members.filter((member) => !assignedIds.has(member.user_id))

  const handleAssign = async (event) => {
    const userId = Number(event.target.value)
    if (!userId) return
    event.target.value = ''
    setAssigning(true)
    try {
      await onAssign(task, userId)
    } finally {
      setAssigning(false)
    }
  }

  return (
    <article className="rounded-lg border border-white/8 bg-surface-850/70 p-4 transition-colors hover:border-white/15">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 flex-1">
          <h4 className="text-sm font-semibold break-words text-slate-100">{task.title}</h4>
          {task.description && (
            <p className="mt-1 text-sm break-words text-slate-400">{task.description}</p>
          )}
        </div>

        <div className="flex shrink-0 flex-col items-end gap-1.5">
          <StatusBadge status={task.status} />
          <PriorityBadge priority={task.priority} />
          <ReminderBadge state={reminder} />
        </div>
      </div>

      <div className="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-slate-500">
        <span>#{task.id}</span>
        <span aria-hidden="true">•</span>
        <span>Created {formatDate(task.created_at)}</span>
        {task.due_date && (
          <>
            <span aria-hidden="true">•</span>
            <span
              className={
                reminder === 'overdue'
                  ? 'font-medium text-rose-400'
                  : reminder === 'due_soon'
                    ? 'font-medium text-amber-400'
                    : 'text-slate-400'
              }
            >
              Due {formatDate(task.due_date)}
            </span>
          </>
        )}
      </div>

      <div className="mt-3 flex flex-wrap items-center gap-2">
        {assignees.length === 0 && <span className="text-xs text-slate-500">Nobody assigned</span>}

        {assignees.map((assignee) => {
          const label = displayUser(assignee.user, assignee.user_id)
          return (
            <span
              key={assignee.user_id}
              className="inline-flex items-center gap-1.5 rounded-full border border-white/10 bg-white/5 py-0.5 pr-1 pl-1.5 text-xs text-slate-300"
            >
              <span className="flex h-5 w-5 items-center justify-center rounded-full bg-surface-600 text-[10px] font-semibold">
                {initialsOf(label)}
              </span>
              {label}
              <button
                type="button"
                onClick={() => onUnassign(task, assignee.user_id)}
                aria-label={`Unassign ${label}`}
                className="cursor-pointer rounded-full p-0.5 text-slate-500 transition-colors hover:bg-white/10 hover:text-slate-200"
              >
                <svg viewBox="0 0 20 20" fill="none" className="h-3 w-3" aria-hidden="true">
                  <path d="M5 5l10 10M15 5L5 15" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" />
                </svg>
              </button>
            </span>
          )
        })}

        {assignable.length > 0 && (
          <select
            defaultValue=""
            onChange={handleAssign}
            disabled={assigning}
            aria-label="Assign a member to this task"
            className="h-7 cursor-pointer rounded-full border border-dashed border-white/15 bg-transparent px-2 text-xs text-slate-400 transition-colors hover:border-white/30 hover:text-slate-200 focus:outline-none"
          >
            <option value="" className="bg-surface-850">
              + Assign…
            </option>
            {assignable.map((member) => (
              <option key={member.user_id} value={member.user_id} className="bg-surface-850">
                {displayUser(member.user, member.user_id)}
              </option>
            ))}
          </select>
        )}
      </div>

      {canEdit && (
        <div className="mt-4 flex flex-wrap items-center gap-2 border-t border-white/8 pt-3">
          <select
            value={task.status}
            onChange={(event) => onStatusChange(task, event.target.value)}
            aria-label="Change status"
            className="h-8 cursor-pointer rounded-lg border border-white/10 bg-surface-800 px-2 text-xs text-slate-200 focus:outline-none"
          >
            {TASK_STATUSES.map((status) => (
              <option key={status} value={status} className="bg-surface-850">
                {STATUS_LABELS[status]}
              </option>
            ))}
          </select>

          <Button variant="ghost" size="sm" onClick={() => onEdit(task)}>
            Edit
          </Button>
          <Button variant="quiet" size="sm" onClick={() => onDelete(task)}>
            Delete
          </Button>
        </div>
      )}
    </article>
  )
}
