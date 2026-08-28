import { useState } from 'react'
import Button from './ui/Button.jsx'
import Logo from './Logo.jsx'
import { logout } from '../lib/auth.js'
import { getCurrentUser } from '../lib/session.js'
import { initialsOf } from '../lib/format.js'

export default function Navbar({ current }) {
  const user = getCurrentUser()
  const [loggingOut, setLoggingOut] = useState(false)

  const handleLogout = async () => {
    setLoggingOut(true)
    await logout()
  }

  const linkClasses = (name) =>
    [
      'rounded-lg px-3 py-1.5 text-sm transition-colors',
      current === name
        ? 'bg-white/10 text-slate-100'
        : 'text-slate-400 hover:bg-white/5 hover:text-slate-200',
    ].join(' ')

  return (
    <header className="sticky top-0 z-40 border-b border-white/8 bg-surface-950/80 backdrop-blur">
      <div className="mx-auto flex h-16 w-full max-w-6xl items-center gap-4 px-4 sm:px-6">
        <a href="/projects.html" className="flex items-center gap-2.5">
          <Logo />
          <span className="hidden text-sm font-semibold text-slate-100 sm:block">Task Manager</span>
        </a>

        <nav className="flex items-center gap-1">
          <a href="/projects.html" className={linkClasses('projects')}>
            Projects
          </a>
        </nav>

        <div className="ml-auto flex items-center gap-3">
          {user && (
            <div className="flex items-center gap-2.5">
              <span
                className="flex h-8 w-8 items-center justify-center rounded-full bg-surface-700 text-xs font-semibold text-slate-200"
                title={user.email}
              >
                {initialsOf(user.name)}
              </span>
              <div className="hidden leading-tight sm:block">
                <p className="text-sm text-slate-200">{user.name}</p>
                <p className="text-xs text-slate-500">{user.email}</p>
              </div>
            </div>
          )}

          <Button variant="ghost" size="sm" onClick={handleLogout} loading={loggingOut}>
            Log out
          </Button>
        </div>
      </div>
    </header>
  )
}
