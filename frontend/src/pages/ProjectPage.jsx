import { useCallback, useEffect, useState } from 'react'
import PageShell from '../components/PageShell.jsx'
import MembersPanel from '../components/MembersPanel.jsx'
import TaskCard from '../components/TaskCard.jsx'
import TaskFormModal from '../components/TaskFormModal.jsx'
import Alert from '../components/ui/Alert.jsx'
import Button, { LinkButton } from '../components/ui/Button.jsx'
import Card, { CardHeader } from '../components/ui/Card.jsx'
import ConfirmDialog from '../components/ui/ConfirmDialog.jsx'
import EmptyState from '../components/ui/EmptyState.jsx'
import Spinner from '../components/ui/Spinner.jsx'
import { StatusBadge } from '../components/ui/Badge.jsx'
import { useToast } from '../components/ui/toastContext.js'
import { api } from '../lib/api.js'
import { getCurrentUser } from '../lib/session.js'
import { formatDate, fromDateInputValue } from '../lib/format.js'
import { TASK_STATUSES, STATUS_LABELS } from '../lib/constants.js'

/**
 * Each page is its own document here, so the id in the query string is fixed
 * for the lifetime of the module — read it once, at import time.
 */
function readProjectId() {
  const parsed = Number(new URLSearchParams(window.location.search).get('id'))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : null
}

const PROJECT_ID = readProjectId()

const fetchProject = async (id) => {
  const data = await api.get(`/projects/${id}`)
  return data?.project ?? null
}

export default function ProjectPage() {
  const toast = useToast()
  const currentUser = getCurrentUser()

  const [project, setProject] = useState(null)
  // Nothing to load when the URL carries no id, so we never enter a loading state.
  const [loading, setLoading] = useState(PROJECT_ID !== null)
  const [error, setError] = useState('')

  const [taskModal, setTaskModal] = useState({ open: false, task: null })
  const [pendingTaskDelete, setPendingTaskDelete] = useState(null)
  const [pendingMemberRemove, setPendingMemberRemove] = useState(null)
  const [busy, setBusy] = useState(false)
  const [statusFilter, setStatusFilter] = useState('ALL')

  // Initial load, inlined so no state is set before the first await.
  useEffect(() => {
    if (PROJECT_ID === null) return undefined
    let cancelled = false

    const run = async () => {
      try {
        const loaded = await fetchProject(PROJECT_ID)
        if (!cancelled) setProject(loaded)
      } catch (loadError) {
        if (!cancelled) setError(loadError.message)
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    run()
    return () => {
      cancelled = true
    }
  }, [])

  const refresh = useCallback(async () => {
    setProject(await fetchProject(PROJECT_ID))
    setError('')
  }, [])

  const members = project?.members ?? []
  const tasks = project?.tasks ?? []
  const isOwner = project?.owner_id === currentUser?.id

  const canEditTask = (task) => isOwner || task.creator_id === currentUser?.id

  const visibleTasks =
    statusFilter === 'ALL' ? tasks : tasks.filter((task) => task.status === statusFilter)

  /**
   * Every mutation reloads the project afterwards, so the UI can never drift
   * from what the API actually stored (permissions are enforced server-side).
   */
  const run = async (action, successMessage) => {
    setBusy(true)
    try {
      await action()
      if (successMessage) toast.success(successMessage)
      await refresh()
      return true
    } catch (actionError) {
      toast.error(actionError.message)
      return false
    } finally {
      setBusy(false)
    }
  }

  /**
   * Form submissions let the error bubble up instead of showing a toast: the
   * form displays it inline next to the fields and stays open, so nothing is
   * reported twice.
   */
  const submit = async (action, successMessage) => {
    setBusy(true)
    try {
      await action()
    } finally {
      setBusy(false)
    }
    toast.success(successMessage)
    await refresh().catch((refreshError) => toast.error(refreshError.message))
  }

  const handleTaskSubmit = (values) => {
    const payload = {
      title: values.title,
      description: values.description,
      priority: values.priority,
      due_date: fromDateInputValue(values.dueDate),
    }

    return taskModal.task
      ? submit(
          () => api.put(`/tasks/${taskModal.task.id}`, { ...payload, status: values.status }),
          'Task updated.',
        )
      : submit(() => api.post(`/projects/${PROJECT_ID}/tasks`, payload), 'Task created.')
  }

  const handleStatusChange = (task, status) =>
    run(() => api.put(`/tasks/${task.id}`, { status }), `Moved to “${STATUS_LABELS[status]}”.`)

  const handleAssign = (task, userId) =>
    run(
      () => api.post(`/projects/${PROJECT_ID}/tasks/${task.id}/assign`, { user_id: userId }),
      'Member assigned.',
    )

  const handleUnassign = (task, userId) =>
    run(() => api.put(`/tasks/${task.id}/unassign`, { user_id: userId }), 'Member unassigned.')

  const handleAddMember = (email, role) =>
    submit(() => api.post(`/projects/${PROJECT_ID}/members`, { email, role }), 'Member added.')

  const handleRemoveMember = async () => {
    const ok = await run(
      () => api.del(`/projects/${PROJECT_ID}/members/${pendingMemberRemove.user_id}`),
      'Member removed.',
    )
    if (ok) setPendingMemberRemove(null)
  }

  const handleDeleteTask = async () => {
    const ok = await run(() => api.del(`/tasks/${pendingTaskDelete.id}`), 'Task deleted.')
    if (ok) setPendingTaskDelete(null)
  }

  const renderError = (message) => (
    <PageShell current="projects">
      <div className="mx-auto max-w-lg py-16 text-center">
        <Alert tone="error">{message}</Alert>
        <LinkButton href="/projects.html" variant="ghost" className="mt-6">
          Back to projects
        </LinkButton>
      </div>
    </PageShell>
  )

  if (PROJECT_ID === null) return renderError('No project id in the URL.')

  if (loading) {
    return (
      <PageShell current="projects">
        <div className="flex items-center justify-center gap-3 py-24 text-sm text-slate-500">
          <Spinner className="h-5 w-5" />
          Loading project…
        </div>
      </PageShell>
    )
  }

  if (error || !project) return renderError(error || 'Project not found.')

  const doneCount = tasks.filter((task) => task.status === 'DONE').length

  return (
    <PageShell current="projects">
      <nav className="mb-6 text-sm text-slate-500">
        <a href="/projects.html" className="hover:text-slate-300">
          Projects
        </a>
        <span className="mx-2" aria-hidden="true">
          /
        </span>
        <span className="text-slate-300">{project.name}</span>
      </nav>

      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="min-w-0">
          <h1 className="text-2xl font-semibold break-words text-slate-50">{project.name}</h1>
          <p className="mt-1.5 max-w-2xl text-sm text-slate-400">
            {project.description || 'No description.'}
          </p>
          <p className="mt-2 text-xs text-slate-500">
            Created {formatDate(project.created_at)} · {members.length} member
            {members.length === 1 ? '' : 's'} · {doneCount}/{tasks.length} tasks done
          </p>
        </div>

        <Button onClick={() => setTaskModal({ open: true, task: null })}>New task</Button>
      </div>

      <div className="mt-8 grid gap-6 lg:grid-cols-[1fr_24rem]">
        <Card className="order-2 lg:order-1">
          <CardHeader
            title={`Tasks (${tasks.length})`}
            actions={
              <select
                value={statusFilter}
                onChange={(event) => setStatusFilter(event.target.value)}
                aria-label="Filter by status"
                className="h-8 cursor-pointer rounded-lg border border-white/10 bg-surface-800 px-2 text-xs text-slate-200 focus:outline-none"
              >
                <option value="ALL" className="bg-surface-850">
                  All statuses
                </option>
                {TASK_STATUSES.map((status) => (
                  <option key={status} value={status} className="bg-surface-850">
                    {STATUS_LABELS[status]}
                  </option>
                ))}
              </select>
            }
          />

          <div className="space-y-3 p-5">
            {tasks.length === 0 ? (
              <EmptyState
                title="No task yet"
                description="Break the project down into tasks and assign them to members."
                action={
                  <Button onClick={() => setTaskModal({ open: true, task: null })}>New task</Button>
                }
              />
            ) : visibleTasks.length === 0 ? (
              <div className="flex flex-col items-center gap-3 py-8 text-center">
                <StatusBadge status={statusFilter} />
                <p className="text-sm text-slate-500">No task with this status.</p>
              </div>
            ) : (
              visibleTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  members={members}
                  canEdit={canEditTask(task)}
                  onEdit={(target) => setTaskModal({ open: true, task: target })}
                  onDelete={setPendingTaskDelete}
                  onAssign={handleAssign}
                  onUnassign={handleUnassign}
                  onStatusChange={handleStatusChange}
                />
              ))
            )}
          </div>
        </Card>

        <div className="order-1 lg:order-2">
          <MembersPanel
            members={members}
            isOwner={isOwner}
            currentUserId={currentUser?.id}
            onAdd={handleAddMember}
            onRemove={setPendingMemberRemove}
          />
        </div>
      </div>

      {taskModal.open && (
        <TaskFormModal
          task={taskModal.task}
          onClose={() => setTaskModal({ open: false, task: null })}
          onSubmit={handleTaskSubmit}
        />
      )}

      <ConfirmDialog
        open={Boolean(pendingTaskDelete)}
        title={`Delete “${pendingTaskDelete?.title ?? ''}”?`}
        confirmLabel="Delete task"
        destructive
        loading={busy}
        onCancel={() => setPendingTaskDelete(null)}
        onConfirm={handleDeleteTask}
      />

      <ConfirmDialog
        open={Boolean(pendingMemberRemove)}
        title="Remove this member?"
        description="They will lose access to the project and its tasks."
        confirmLabel="Remove member"
        destructive
        loading={busy}
        onCancel={() => setPendingMemberRemove(null)}
        onConfirm={handleRemoveMember}
      />
    </PageShell>
  )
}
