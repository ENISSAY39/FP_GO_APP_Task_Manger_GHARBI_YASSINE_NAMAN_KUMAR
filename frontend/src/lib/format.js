/** Human-readable date, tolerant to the zero dates GORM can return. */
export function formatDate(value) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getFullYear() < 1970) return '—'
  return date.toLocaleDateString(undefined, {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

export function formatDateTime(value) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getFullYear() < 1970) return '—'
  return date.toLocaleString(undefined, {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/** <input type="date"> wants YYYY-MM-DD; the API sends RFC3339. */
export function toDateInputValue(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getFullYear() < 1970) return ''
  return date.toISOString().slice(0, 10)
}

/** ...and the API wants RFC3339 back. */
export function fromDateInputValue(value) {
  if (!value) return null
  const date = new Date(`${value}T12:00:00`)
  return Number.isNaN(date.getTime()) ? null : date.toISOString()
}

/**
 * Members and assignees carry a nested user object, which is only populated
 * when the backend preloads it — fall back to the numeric id rather than
 * rendering an empty name.
 */
export function displayUser(user, userId) {
  if (user && user.name) return user.name
  if (user && user.email) return user.email
  return `User #${userId}`
}

export function initialsOf(label) {
  const parts = String(label || '?')
    .replace(/^User #/, '')
    .trim()
    .split(/\s+/)
    .filter(Boolean)
  if (parts.length === 0) return '?'
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase()
  return (parts[0][0] + parts[1][0]).toUpperCase()
}
