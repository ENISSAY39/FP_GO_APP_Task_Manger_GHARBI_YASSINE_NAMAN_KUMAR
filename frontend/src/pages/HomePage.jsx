import { LinkButton } from '../components/ui/Button.jsx'
import Logo from '../components/Logo.jsx'

const FEATURES = [
  {
    title: 'Projects',
    description: 'Create projects, invite members and keep ownership rules enforced server-side.',
  },
  {
    title: 'Tasks',
    description: 'Track status, priority and due dates, and assign the people who work on them.',
  },
  {
    title: 'Secure by default',
    description: 'JWT authentication with role-based permissions on every protected endpoint.',
  },
]

export default function HomePage() {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="mx-auto flex h-16 w-full max-w-6xl items-center px-4 sm:px-6">
        <a href="/" className="flex items-center gap-2.5">
          <Logo />
          <span className="text-sm font-semibold text-slate-100">Task Manager</span>
        </a>
        <div className="ml-auto flex items-center gap-2">
          <LinkButton href="/login.html" variant="ghost" size="sm">
            Log in
          </LinkButton>
          <LinkButton href="/signup.html" size="sm">
            Sign up
          </LinkButton>
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl flex-1 px-4 sm:px-6">
        <section className="py-20 text-center sm:py-28">
          <p className="text-xs font-medium tracking-[0.2em] text-accent-400 uppercase">
            Go · Gin · GORM · React
          </p>
          <h1 className="mx-auto mt-4 max-w-3xl text-4xl font-semibold text-balance text-slate-50 sm:text-5xl">
            Plan projects, split the work, ship it.
          </h1>
          <p className="mx-auto mt-5 max-w-xl text-base text-pretty text-slate-400">
            A collaborative task manager: shared projects, member roles, and tasks you can
            prioritise, schedule and assign.
          </p>
          <div className="mt-9 flex flex-wrap items-center justify-center gap-3">
            <LinkButton href="/signup.html" size="lg">
              Create an account
            </LinkButton>
            <LinkButton href="/login.html" variant="ghost" size="lg">
              I already have one
            </LinkButton>
          </div>
        </section>

        <section className="grid gap-4 pb-20 sm:grid-cols-3">
          {FEATURES.map((feature) => (
            <div
              key={feature.title}
              className="rounded-xl border border-white/8 bg-surface-900/60 p-5"
            >
              <h2 className="text-sm font-semibold text-slate-100">{feature.title}</h2>
              <p className="mt-2 text-sm text-slate-400">{feature.description}</p>
            </div>
          ))}
        </section>
      </main>

      <footer className="border-t border-white/8 py-6">
        <p className="mx-auto max-w-6xl px-4 text-xs text-slate-600 sm:px-6">
          Task Manager — Gharbi Yassine &amp; Naman Kumar
        </p>
      </footer>
    </div>
  )
}
