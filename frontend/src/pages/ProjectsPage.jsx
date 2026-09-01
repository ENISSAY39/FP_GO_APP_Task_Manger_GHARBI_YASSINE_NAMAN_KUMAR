import { useCallback, useEffect, useState } from 'react'
import PageShell from '../components/PageShell.jsx'
import ProjectCard from '../components/ProjectCard.jsx'
import ProjectFormModal from '../components/ProjectFormModal.jsx'
import Alert from '../components/ui/Alert.jsx'
import Button from '../components/ui/Button.jsx'
import ConfirmDialog from '../components/ui/ConfirmDialog.jsx'
import EmptyState from '../components/ui/EmptyState.jsx'
import Spinner from '../components/ui/Spinner.jsx'
import { useToast } from '../components/ui/toastContext.js'
import { api } from '../lib/api.js'
import { getCurrentUser } from '../lib/session.js'

const fetchProjects = async () => {
  const data = await api.get('/projects')
  return data?.projects ?? []
}

export default function ProjectsPage() {
  const toast = useToast()
  const currentUser = getCurrentUser()

  const [projects, setProjects] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [creating, setCreating] = useState(false)
  const [pendingDelete, setPendingDelete] = useState(null)
  const [deleting, setDeleting] = useState(false)

  // Initial load. The request is inlined here (rather than calling the
  // refresh() below) so no state is set before the first await, and so an
  // in-flight response is ignored if the page unmounts first.
  useEffect(() => {
    let cancelled = false

    const run = async () => {
      try {
        const list = await fetchProjects()
        if (!cancelled) setProjects(list)
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
    try {
      setProjects(await fetchProjects())
      setError('')
    } catch (loadError) {
      setError(loadError.message)
    } finally {
      setLoading(false)
    }
  }, [])

  const handleCreate = async ({ name, description }) => {
    await api.post('/projects', { name, description })
    toast.success(`Project “${name}” created.`)
    await refresh()
  }

  const handleDelete = async () => {
    setDeleting(true)
    try {
      await api.del(`/projects/${pendingDelete.id}`)
      toast.success('Project deleted.')
      setPendingDelete(null)
      await refresh()
    } catch (deleteError) {
      toast.error(deleteError.message)
    } finally {
      setDeleting(false)
    }
  }

  const retry = () => {
    setLoading(true)
    setError('')
    refresh()
  }

  return (
    <PageShell current="projects">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-slate-50">Projects</h1>
          <p className="mt-1 text-sm text-slate-400">
            {currentUser ? `Signed in as ${currentUser.name}. ` : ''}Projects you own or belong to.
          </p>
        </div>
        <Button onClick={() => setCreating(true)}>New project</Button>
      </div>

      <div className="mt-8">
        {loading ? (
          <div className="flex items-center justify-center gap-3 py-20 text-sm text-slate-500">
            <Spinner className="h-5 w-5" />
            Loading projects…
          </div>
        ) : error ? (
          <Alert tone="error">
            {error}{' '}
            <button
              type="button"
              onClick={retry}
              className="cursor-pointer underline underline-offset-2"
            >
              Retry
            </button>
          </Alert>
        ) : projects.length === 0 ? (
          <EmptyState
            title="No project yet"
            description="Create your first project to start adding members and tasks."
            action={<Button onClick={() => setCreating(true)}>New project</Button>}
          />
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {projects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                isOwner={project.owner_id === currentUser?.id}
                currentUserId={currentUser?.id}
                onDelete={setPendingDelete}
              />
            ))}
          </div>
        )}
      </div>

      {creating && (
        <ProjectFormModal onClose={() => setCreating(false)} onSubmit={handleCreate} />
      )}

      <ConfirmDialog
        open={Boolean(pendingDelete)}
        title={`Delete “${pendingDelete?.name ?? ''}”?`}
        description="The project, its members and all its tasks will be removed."
        confirmLabel="Delete project"
        destructive
        loading={deleting}
        onCancel={() => setPendingDelete(null)}
        onConfirm={handleDelete}
      />
    </PageShell>
  )
}
