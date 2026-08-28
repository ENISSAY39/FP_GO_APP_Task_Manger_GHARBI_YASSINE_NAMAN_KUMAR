export default function Logo({ className = 'h-8 w-8' }) {
  return (
    <span
      className={`inline-flex ${className} items-center justify-center rounded-lg bg-gradient-to-br from-accent-400 to-accent-600 text-surface-950 shadow-sm shadow-accent-500/30`}
      aria-hidden="true"
    >
      <svg viewBox="0 0 24 24" fill="none" className="h-1/2 w-1/2">
        <path
          d="M4 12.5l5 5L20 6.5"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
    </span>
  )
}
