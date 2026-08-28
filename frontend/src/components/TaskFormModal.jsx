import { useState } from 'react'
import Button from './ui/Button.jsx'
import Field, { SelectField, TextareaField } from './ui/Field.jsx'
import Modal from './ui/Modal.jsx'
import Alert from './ui/Alert.jsx'
import { TASK_PRIORITIES, TASK_STATUSES, PRIORITY_LABELS, STATUS_LABELS } from '../lib/constants.js'
import { toDateInputValue } from '../lib/format.js'

const EMPTY = { title: '', description: '', priority: 'MEDIUM', status: 'TODO', dueDate: '' }

/**
 * One modal for both create and edit. The old frontend chained four
 * window.prompt() calls to edit a task.
 *
 * Rendered only while open (the parent unmounts it on close), so the initial
 * state below is enough — no effect is needed to sync it with the task prop.
 */
export default function TaskFormModal({ task, onClose, onSubmit }) {
  const isEdit = Boolean(task)
  const [values, setValues] = useState(() =>
    task
      ? {
          title: task.title ?? '',
          description: task.description ?? '',
          priority: task.priority ?? 'MEDIUM',
          status: task.status ?? 'TODO',
          dueDate: toDateInputValue(task.due_date),
        }
      : EMPTY,
  )
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)

  const update = (key) => (event) => setValues((current) => ({ ...current, [key]: event.target.value }))

  const handleSubmit = async (event) => {
    event.preventDefault()
    if (!values.title.trim()) {
      setError('A title is required.')
      return
    }

    setSaving(true)
    setError('')
    try {
      await onSubmit({ ...values, title: values.title.trim(), description: values.description.trim() })
      onClose()
    } catch (submitError) {
      setError(submitError.message)
      setSaving(false)
    }
  }

  return (
    <Modal
      open
      onClose={saving ? undefined : onClose}
      title={isEdit ? 'Edit task' : 'New task'}
      footer={
        <>
          <Button variant="ghost" onClick={onClose} disabled={saving}>
            Cancel
          </Button>
          <Button form="task-form" type="submit" loading={saving}>
            {isEdit ? 'Save changes' : 'Create task'}
          </Button>
        </>
      }
    >
      <form id="task-form" onSubmit={handleSubmit} className="space-y-4">
        <Alert tone="error">{error}</Alert>

        <Field
          label="Title"
          placeholder="Write the API documentation"
          value={values.title}
          onChange={update('title')}
          autoFocus
          required
        />

        <TextareaField
          label="Description"
          placeholder="Any detail worth keeping"
          value={values.description}
          onChange={update('description')}
        />

        <div className="grid gap-4 sm:grid-cols-2">
          <SelectField
            label="Priority"
            value={values.priority}
            onChange={update('priority')}
            options={TASK_PRIORITIES.map((value) => ({ value, label: PRIORITY_LABELS[value] }))}
          />

          {isEdit && (
            <SelectField
              label="Status"
              value={values.status}
              onChange={update('status')}
              options={TASK_STATUSES.map((value) => ({ value, label: STATUS_LABELS[value] }))}
            />
          )}

          <Field
            label="Due date"
            type="date"
            value={values.dueDate}
            onChange={update('dueDate')}
            hint="Optional"
          />
        </div>
      </form>
    </Modal>
  )
}
