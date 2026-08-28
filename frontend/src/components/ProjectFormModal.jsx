import { useState } from 'react'
import Button from './ui/Button.jsx'
import Field, { TextareaField } from './ui/Field.jsx'
import Modal from './ui/Modal.jsx'
import Alert from './ui/Alert.jsx'

/**
 * Rendered only while open (the parent unmounts it on close), so the form
 * resets simply by being mounted again — no effect needed to clear it.
 */
export default function ProjectFormModal({ onClose, onSubmit }) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (event) => {
    event.preventDefault()
    if (!name.trim()) {
      setError('A project name is required.')
      return
    }

    setSaving(true)
    setError('')
    try {
      await onSubmit({ name: name.trim(), description: description.trim() })
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
      title="New project"
      description="You are added as the owner automatically."
      footer={
        <>
          <Button variant="ghost" onClick={onClose} disabled={saving}>
            Cancel
          </Button>
          <Button form="project-form" type="submit" loading={saving}>
            Create project
          </Button>
        </>
      }
    >
      <form id="project-form" onSubmit={handleSubmit} className="space-y-4">
        <Alert tone="error">{error}</Alert>

        <Field
          label="Name"
          placeholder="Website redesign"
          value={name}
          onChange={(event) => setName(event.target.value)}
          autoFocus
          required
        />

        <TextareaField
          label="Description"
          placeholder="What is this project about?"
          value={description}
          onChange={(event) => setDescription(event.target.value)}
        />
      </form>
    </Modal>
  )
}
