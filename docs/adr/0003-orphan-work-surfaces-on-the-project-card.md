# ADR 0003 — Orphan work surfaces on the Project card, not in the navbar

**Status**: accepted, 2026-09-03

## Context

ADR 0002 made every Reminder surface either Task-scoped or viewer-scoped. The
Task card flags any Task that has slipped; the Project card count and the navbar
list count only Tasks assigned to the viewer, because they answer "what do *I*
owe?".

That left a hole, recorded in `docs/open-questions.md`: a Task with a past due
date and **no assignee at all** shows its badge on the Task card, and appears in
nobody's Project count and nobody's navbar list. It is late, and no surface
tells anyone. It is found only by opening the Project and reading the board.

Three shapes were considered.

**Full escalation.** Put orphan Tasks into the Owner's navbar list alongside
their own. Nothing falls through the cracks, but the navbar stops meaning one
thing: the count mixes "assigned to me" with "assigned to nobody". A Project
carrying fifty orphan Tasks makes the personal counter unreadable, which costs
the feature its whole point — a number is only actionable while it stays small
enough to act on.

**No escalation.** Accept the hole and write down why. Zero code, and the navbar
keeps its single meaning, but the Owner still has to go looking.

**Project card only.** The chosen shape.

## Decision

An **Orphan** — a Task that carries a Reminder and has no assignee — is surfaced
to the Project **Owner**, on the **Project card only**, as a badge distinct from
the personal count.

The navbar stays exactly as it is: strictly personal, one meaning, unchanged by
this decision.

Each surface therefore keeps a single answerable question:

| Surface | Question it answers | Scope |
| --- | --- | --- |
| Task card | Has this Task slipped? | the Task |
| Project card, personal count | What do I owe here? | the viewer |
| Project card, orphan count | What does nobody own here? | Owner only |
| Navbar | What do I owe, everywhere? | the viewer |

The Owner is the right recipient because assigning work is already their
responsibility: the signal lands on the person who can act on it, and the action
it invites — assign somebody — is one click away on the board it points at.

No backend change is required. `GetMyProjects` already returns `owner_id` and
preloads `Tasks.Assignees`, so both the ownership test and the orphan test are
answerable from data the client already holds. This decision costs one predicate
and one badge.

## Consequences

- **The hole closes where triage already happens.** An Owner scanning their
  projects list sees which Projects hold unowned late work, without opening each
  one.
- **The navbar keeps one meaning.** Its count remains "mine", so an empty navbar
  stays a trustworthy all-clear. This is the property full escalation would have
  spent.
- **Non-owners see nothing new.** A member of a busy Project is not shown work
  they cannot assign. The badge is Owner-only.
- **An Owner who assigns nobody still gets a growing number.** The badge reports
  the problem; it does not solve it. If orphan counts routinely run large, that
  is evidence about the team's assignment habits, and the answer is a workflow
  change rather than another surface.
- **A Project with no Owner is out of scope.** Every Project currently has an
  `owner_id`, so the case does not arise; if ownership ever becomes optional,
  orphan work in an ownerless Project reaches nobody again, and this ADR is the
  one to reopen.
- **"Orphan" enters the glossary**, so the three surfaces and any future one
  cannot drift into calling it something else.
