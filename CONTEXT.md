# Task Manager

A multi-project task tracker: users join Projects as Owner or Member and
manage Tasks with a status, a priority, and optional assignees.

## Language

**Reminder**:
The signal that a Task needs attention because of its due date — either
Overdue or Due soon. Computed on read from `Task.DueDate` and `Task.Status`;
not a stored record. Scoped narrowly to due dates, not a general
notification system.
_Avoid_: Notification (reserved for a possible future, broader event system
covering things like assignment or membership changes — not what this covers)

**Overdue**:
A Task whose `DueDate` is in the past and whose `Status` is not `DONE`.

**Due soon**:
A Task whose `DueDate` falls within the next 24 hours and whose `Status` is
not `DONE`. The 24-hour window is fixed, not user-configurable.

**Task status**:
The stage a Task is in. Exactly three values, spelled identically in the
database, in the Go constants and in the frontend: `TODO`, `IN_PROGRESS`,
`DONE`. New Tasks default to `TODO`. This list is the source of truth for
both sides — a value outside it is rejected, never stored.
_Avoid_: `DOING` (an old Go spelling that no row ever held; removed rather
than migrated to)

**Task priority**:
How urgent a Task is, independently of its due date. Exactly three values:
`LOW`, `MEDIUM`, `HIGH`, defaulting to `MEDIUM`. Same rule as status: a value
outside the list is rejected, never stored.
