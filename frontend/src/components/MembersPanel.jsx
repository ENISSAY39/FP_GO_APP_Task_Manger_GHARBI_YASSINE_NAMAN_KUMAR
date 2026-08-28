import { useState } from 'react'
import Button from './ui/Button.jsx'
import Card, { CardHeader } from './ui/Card.jsx'
import Field, { SelectField } from './ui/Field.jsx'
import { RoleBadge } from './ui/Badge.jsx'
import { MEMBER_ROLES } from '../lib/constants.js'
import { displayUser, initialsOf } from '../lib/format.js'

export default function MembersPanel({ members, isOwner, currentUserId, onAdd, onRemove }) {
  const [userId, setUserId] = useState('')
  const [role, setRole] = useState('MEMBER')
  const [adding, setAdding] = useState(false)
  const [error, setError] = useState('')

  const handleAdd = async (event) => {
    event.preventDefault()
    const parsed = Number(userId)
    if (!Number.isInteger(parsed) || parsed <= 0) {
      setError('Enter a numeric user id.')
      return
    }

    setAdding(true)
    setError('')
    try {
      await onAdd(parsed, role)
      setUserId('')
      setRole('MEMBER')
    } catch (addError) {
      setError(addError.message)
    } finally {
      setAdding(false)
    }
  }

  return (
    <Card>
      <CardHeader title={`Members (${members.length})`} />

      <ul className="divide-y divide-white/5">
        {members.map((member) => {
          const label = displayUser(member.user, member.user_id)
          return (
            <li key={member.id ?? member.user_id} className="flex items-center gap-2.5 px-4 py-3">
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-surface-700 text-xs font-semibold text-slate-200">
                {initialsOf(label)}
              </span>

              <div className="min-w-0 flex-1">
                <p className="truncate text-sm text-slate-200">
                  {label}
                  {member.user_id === currentUserId && (
                    <span className="ml-1.5 text-xs text-slate-500">(you)</span>
                  )}
                </p>
                {member.user?.email && (
                  <p className="truncate text-xs text-slate-500">{member.user.email}</p>
                )}
              </div>

              <RoleBadge role={member.role} />

              {isOwner && member.role !== 'OWNER' && (
                <button
                  type="button"
                  onClick={() => onRemove(member)}
                  title={`Remove ${label}`}
                  aria-label={`Remove ${label}`}
                  className="shrink-0 cursor-pointer rounded-lg p-1.5 text-slate-500 transition-colors hover:bg-white/5 hover:text-rose-300"
                >
                  <svg viewBox="0 0 20 20" fill="none" className="h-4 w-4" aria-hidden="true">
                    <path
                      d="M5 5l10 10M15 5L5 15"
                      stroke="currentColor"
                      strokeWidth="1.8"
                      strokeLinecap="round"
                    />
                  </svg>
                </button>
              )}
            </li>
          )
        })}
      </ul>

      {isOwner && (
        <form onSubmit={handleAdd} className="border-t border-white/8 px-5 py-4">
          <p className="mb-3 text-xs font-medium tracking-wide text-slate-400 uppercase">
            Add a member
          </p>

          <div className="flex flex-wrap items-start gap-2">
            <div className="min-w-[9rem] flex-1">
              <Field
                placeholder="User id"
                inputMode="numeric"
                value={userId}
                onChange={(event) => setUserId(event.target.value)}
                error={error}
                hint={error ? undefined : 'The API identifies users by their numeric id.'}
              />
            </div>

            <div className="w-32">
              <SelectField
                value={role}
                onChange={(event) => setRole(event.target.value)}
                options={MEMBER_ROLES.map((value) => ({ value, label: value }))}
              />
            </div>

            <Button type="submit" loading={adding}>
              Add
            </Button>
          </div>
        </form>
      )}
    </Card>
  )
}
