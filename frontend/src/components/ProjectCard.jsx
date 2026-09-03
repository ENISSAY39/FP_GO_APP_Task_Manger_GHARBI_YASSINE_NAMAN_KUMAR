import Badge from './ui/Badge.jsx'
import Button, { LinkButton } from './ui/Button.jsx'
import { formatDate } from '../lib/format.js'
import { isAssignedTo, isOrphan, reminderStateOf } from '../lib/reminders.js'

export default function ProjectCard({ project, isOwner, currentUserId, onDelete }) {
  const memberCount = project.members?.length ?? 0
  const taskCount = project.tasks?.length ?? 0
  const doneCount = project.tasks?.filter((task) => task.status === 'DONE').length ?? 0
  const progress = taskCount > 0 ? Math.round((doneCount / taskCount) * 100) : 0
  const reminderCount = (project.tasks ?? []).filter(
    (task) => isAssignedTo(task, currentUserId) && reminderStateOf(task),
  ).length
  // Orphan work is the Owner's to hand out, so only they are shown it
  // (ADR 0003). It is deliberately absent from the personal count above and
  // from the navbar: nobody owes an Orphan, which is the whole problem.
  const orphanCount = isOwner ? (project.tasks ?? []).filter((task) => isOrphan(task)).length : 0

  return (
    <article className="group flex flex-col rounded-xl border border-white/8 bg-surface-900/80 p-5 transition-colors hover:border-white/15">
      <div className="flex items-start justify-between gap-3">
        <h3 className="line-clamp-2 min-w-0 flex-1 text-base font-semibold text-balance text-slate-100">
          {project.name}
        </h3>
        <div className="flex shrink-0 items-center gap-1.5">
          {reminderCount > 0 && (
            <Badge tone="rose">
              {reminderCount} due
            </Badge>
          )}
          {orphanCount > 0 && (
            <Badge tone="amber">{orphanCount} orphaned</Badge>
          )}
          {isOwner && <Badge tone="violet">Owner</Badge>}
        </div>
      </div>

      <p className="mt-2 line-clamp-2 min-h-[2.5rem] text-sm text-slate-400">
        {project.description || 'No description.'}
      </p>

      <div className="mt-4">
        <div className="flex items-center justify-between text-xs text-slate-500">
          <span>
            {doneCount}/{taskCount} tasks done
          </span>
          <span>{progress}%</span>
        </div>
        <div className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-surface-700">
          <div
            className="h-full rounded-full bg-accent-500 transition-[width] duration-500"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>

      <div className="mt-4 flex items-center gap-3 text-xs text-slate-500">
        <span>
          {memberCount} member{memberCount === 1 ? '' : 's'}
        </span>
        <span aria-hidden="true">•</span>
        <span>Created {formatDate(project.created_at)}</span>
      </div>

      <div className="mt-5 flex items-center gap-2 border-t border-white/8 pt-4">
        <LinkButton size="sm" href={`/project.html?id=${project.id}`}>
          Open
        </LinkButton>
        {isOwner && (
          <Button variant="quiet" size="sm" onClick={() => onDelete(project)}>
            Delete
          </Button>
        )}
      </div>
    </article>
  )
}
