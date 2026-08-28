import { LinkButton } from '../components/ui/Button.jsx'
import Logo from '../components/Logo.jsx'
import GitHubIcon from '../components/GitHubIcon.jsx'
import { REPO_URL, REPO_LABEL } from '../lib/constants.js'

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
          <a
            href={REPO_URL}
            target="_blank"
            rel="noreferrer noopener"
            title="Source code on GitHub"
            aria-label="Source code on GitHub"
            className="mr-1 rounded-lg p-2 text-slate-400 transition-colors hover:bg-white/5 hover:text-slate-100"
          >
            <GitHubIcon className="h-5 w-5" />
          </a>
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

        <section className="grid gap-4 sm:grid-cols-3">
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

        <section className="py-20">
          <div className="rounded-xl border border-white/8 bg-surface-900/60 p-8 text-center sm:p-10">
            <span className="inline-flex h-11 w-11 items-center justify-center rounded-full bg-surface-800 text-slate-300">
              <GitHubIcon className="h-5 w-5" />
            </span>

            <h2 className="mt-4 text-xl font-semibold text-slate-50">
              Open source — contributions welcome
            </h2>

            <p className="mx-auto mt-3 max-w-xl text-sm text-pretty text-slate-400">
              This is a student project, and the repository is public. If you want to collaborate,
              suggest a feature, or report something that does not work, issues and pull requests
              are very welcome.
            </p>

            <div className="mt-6 flex flex-wrap items-center justify-center gap-3">
              <LinkButton href={REPO_URL} target="_blank" rel="noreferrer noopener">
                <GitHubIcon className="h-4 w-4" />
                View the repository
              </LinkButton>
              <LinkButton
                href={`${REPO_URL}/issues/new`}
                target="_blank"
                rel="noreferrer noopener"
                variant="ghost"
              >
                Open an issue
              </LinkButton>
            </div>

            <p className="mt-5 font-mono text-xs break-all text-slate-600">{REPO_LABEL}</p>
          </div>
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
