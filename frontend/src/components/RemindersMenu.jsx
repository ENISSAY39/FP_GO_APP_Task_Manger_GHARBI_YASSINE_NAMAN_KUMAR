import { useEffect, useRef, useState } from 'react'
import { ReminderBadge } from './ui/Badge.jsx'
import { api } from '../lib/api.js'
import { getCurrentUser } from '../lib/session.js'
import { collectReminders } from '../lib/reminders.js'

/**
 * Global "what's due" indicator, shared across every protected page via
 * Navbar. Fetches once on mount (no polling — recomputed on the next page
 * load/navigation, same as the rest of the reminder feature) and stays
 * silently empty if the request fails: this is a convenience overlay, not
 * something that should ever break the page around it.
 */
export default function RemindersMenu() {
  const currentUser = getCurrentUser()
  const [reminders, setReminders] = useState([])
  const [open, setOpen] = useState(false)
  const containerRef = useRef(null)

  useEffect(() => {
    if (!currentUser?.id) return undefined
    let cancelled = false

    api
      .get('/projects')
      .then((data) => {
        if (!cancelled) setReminders(collectReminders(data?.projects ?? [], currentUser.id))
      })
      .catch(() => {})

    return () => {
      cancelled = true
    }
  }, [currentUser?.id])

  useEffect(() => {
    if (!open) return undefined
    const handleClickOutside = (event) => {
      if (containerRef.current && !containerRef.current.contains(event.target)) setOpen(false)
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [open])

  if (!currentUser) return null

  return (
    <div className="relative" ref={containerRef}>
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        aria-label={`Reminders${reminders.length ? `, ${reminders.length} due` : ''}`}
        aria-expanded={open}
        className="relative flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg text-slate-400 transition-colors hover:bg-white/5 hover:text-slate-200"
      >
        <BellIcon className="h-4 w-4" />
        {reminders.length > 0 && (
          <span className="absolute -top-0.5 -right-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-rose-500 px-1 text-[10px] font-semibold text-white">
            {reminders.length}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-2 w-80 rounded-xl border border-white/10 bg-surface-900 p-2 shadow-xl shadow-black/40">
          {reminders.length === 0 ? (
            <p className="p-3 text-sm text-slate-500">Nothing due soon.</p>
          ) : (
            <ul className="max-h-80 space-y-1 overflow-y-auto">
              {reminders.map(({ task, project, state }) => (
                <li key={task.id}>
                  <a
                    href={`/project.html?id=${project.id}`}
                    className="block rounded-lg p-2.5 transition-colors hover:bg-white/5"
                  >
                    <div className="flex items-start justify-between gap-2">
                      <span className="min-w-0 flex-1 truncate text-sm text-slate-200">
                        {task.title}
                      </span>
                      <ReminderBadge state={state} />
                    </div>
                    <p className="mt-0.5 truncate text-xs text-slate-500">{project.name}</p>
                  </a>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  )
}

function BellIcon({ className }) {
  return (
    <svg viewBox="0 0 20 20" fill="none" className={className} aria-hidden="true">
      <path
        d="M10 2.5c-2.9 0-4.5 2.1-4.5 4.9v2.1c0 .5-.2 1.2-.5 1.7L4 12.7c-.6 1-.2 2.1.9 2.5 3.6 1.2 7.6 1.2 11.2 0 1-.3 1.5-1.5.9-2.5l-1-1.5c-.3-.5-.5-1.2-.5-1.7V7.4c0-2.7-1.7-4.9-4.5-4.9Z"
        stroke="currentColor"
        strokeWidth="1.4"
        strokeLinejoin="round"
      />
      <path
        d="M7.8 16.2a2.2 2.2 0 0 0 4.4 0"
        stroke="currentColor"
        strokeWidth="1.4"
        strokeLinecap="round"
      />
    </svg>
  )
}
