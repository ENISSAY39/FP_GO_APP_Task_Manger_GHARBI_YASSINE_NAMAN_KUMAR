import Navbar from './Navbar.jsx'

/** Common chrome for the signed-in pages. */
export default function PageShell({ current, children }) {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar current={current} />
      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-8 sm:px-6">{children}</main>
      <footer className="border-t border-white/8 py-6">
        <p className="mx-auto max-w-6xl px-4 text-xs text-slate-600 sm:px-6">
          Task Manager — Go, Gin, GORM &amp; React
        </p>
      </footer>
    </div>
  )
}
