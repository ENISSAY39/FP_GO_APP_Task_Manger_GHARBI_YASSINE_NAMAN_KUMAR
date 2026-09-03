# ADR 0002 — A Reminder is derived, never stored

**Status**: accepted, 2026-09-01

## Context

A user belongs to several Projects and can be assigned Tasks in any of them.
Before this work, a due date was plain grey text on a Task card, visible only
after opening the Project that held it: work already late looked exactly like
work due next month.

The feature surfaces a **Reminder** — the signal that a Task needs attention
because of its due date, either **Overdue** or **Due soon**, as defined in
`CONTEXT.md`.

The design question was where that signal lives. Two shapes were available:

- **Stored**: a reminder record per Task, written by a scheduled job or on save,
  read back by the UI.
- **Derived**: no record at all; the signal is computed from the Task's due date
  and status against the current clock, every time something renders.

A stored reminder brings a table, a migration, a job that must run, and a
staleness problem: the record says "Due soon" until something rewrites it,
including after the Task has silently gone Overdue.

## Decision

A Reminder is derived, never stored. No table, no column, no migration, no
reminder record, no read/dismiss state, no persistence of any kind.

One frontend module owns the rules and every surface calls it:
`reminderStateOf(task, now)` returns `'overdue' | 'due_soon' | null`, with
`isAssignedTo` and `collectReminders` beside it. `now` defaults to the current
time and exists as a parameter so tests can freeze it.

The 24-hour Due soon window is a constant, not configuration.

The one backend change this required was a preload, not a contract change: the
"my projects" handler eagerly loads Task assignees so the viewer-scoped surfaces
have assignee data. No handler computes reminder state; the Go side stays
unaware of the concept.

## Consequences

- **Never stale.** The state is recomputed on every render, so a Task that
  crosses from Due soon into Overdue while a tab is open shows the truth on the
  next page load rather than a cached lie.
- **No new API surface.** The navbar indicator reuses the existing projects
  endpoint. Nothing to version, nothing to migrate, nothing to backfill.
- **One rule, three surfaces.** The Task card, the Project card and the navbar
  cannot disagree about whether a Task needs attention, because they ask the
  same function.
- **The clock is the viewer's browser.** A user with a badly skewed clock sees
  skewed Reminders. Accepted: the alternative is a server-computed field, which
  is what this decision declines to build.
- **Reminders cannot be delivered.** No email, push or digest is possible
  without revisiting this, since there is no server-side notion of a Reminder to
  send. That is the decision to reopen if delivery is ever wanted — and the
  point at which "Notification", a word `CONTEXT.md` deliberately reserves,
  would start to mean something.
- **Scope is asymmetric by design.** The Task card flag is task-scoped; the
  Project card count and the navbar list are viewer-scoped. This is a product
  decision, recorded here because it looks like an inconsistency and is the
  thing a later refactor is most likely to flatten by accident.
