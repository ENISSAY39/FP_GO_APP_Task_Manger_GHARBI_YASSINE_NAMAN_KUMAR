# Open questions

Ambiguities surfaced while specifying and building the due-date Reminder feature
and the task status vocabulary fix, left open on purpose. Each one is a decision
nobody has made yet, not a bug.

Written 2026-09-01, last reviewed 2026-09-04. Review before the next block of
work.

Settled since: unassigned overdue work now reaches the Project Owner on the
Project card — see ADR 0003, and the **Orphan** entry in `CONTEXT.md`.

Resolved decisions do not live here: they go to `docs/adr/`, and settled
vocabulary goes to `CONTEXT.md`.

The design interview that preceded the spec is recorded below, under **From the
design interview**. That session is not recoverable from the repo — only its
output survives — so those entries were reconstructed from the transcript on
2026-09-04.

## From the design interview

Nine questions were put and answered before any code was written. Six of the
answers closed a door that nobody has since reopened; each is a deliberate
omission, not an oversight.

**Email delivery was deferred, not rejected.**  
The interview settled on in-app only because the repo has no mail
infrastructure: no SMTP library in `go.mod`, no mail variables in `.env.exempl`.
The reasoning was that email is a later layer over the same derived data rather
than a different feature. Whether it is actually wanted was never decided.

**The Reminder concept was kept deliberately narrow.**  
A generic `Notification` model — an event catalogue plus per-type delivery
preferences — was weighed and declined as speculative. The open question is what
happens when a second notification type arrives: whether Reminder becomes a
special case of something larger, or keeps its own path.

**A Reminder cannot be dismissed.**  
This falls straight out of computing on read (ADR 0002): there is no record to
mark as seen, so a Reminder stays until the Task is Done or its due date moves.
It was named during the interview as a real limitation, accepted on purpose.

**There is no dedicated "My Reminders" page.**  
Three tiers of UI placement were weighed and the third — an aggregating page
with its own Vite entry point — was cut, on the grounds that the navbar dropdown
answers the same need for less. Whether the dropdown still suffices as a user's
Task count grows has not been tested.

**Nothing refreshes while a page stays open.**  
Polling was declined because the README treats real-time updates as a separate
future item. A Reminder whose state changes while a tab sits open is corrected
only on the next load or navigation.

**The Owner overview was scoped out, then partly reopened.**  
The interview scoped Reminders to assignees only and named a project-wide Owner
view as a v2. ADR 0003 has since granted the Owner one slice of it — unassigned
overdue work, surfaced as Orphan on the project card — but not the general case
of seeing every assignee's late work. The rest is still open.

## Product

**The 24-hour window is fixed for everyone.**  
No setting, no per-user preference. Fine today; the question is whether a
Project with a different rhythm ever needs its own threshold.

**Clock skew is the viewer's problem.**  
Reminder state is computed in the browser, so a badly-set clock produces wrong
Reminders. Accepted when the feature was designed. The fix is a server-computed
field, which ADR 0002 explicitly declines to build.

**The scope asymmetry is undocumented in the UI.**  
The Task card flags work regardless of assignee; the Project count and navbar
count only the viewer's own. Deliberate, tested from both sides, and invisible
to a user who never reads the tests. Whether it needs explaining in the
interface is unanswered.

## Testing

**Components do not take an injected clock.**  
The spec asked for the rule to accept an injected clock so boundaries could be
tested deterministically, and the rule module does. The three components do not:
they call `reminderStateOf(task)` through its default argument, so their tests
freeze the system clock with `vi.setSystemTime` instead. Same determinism,
different mechanism. Giving each component a `now` prop would make the spec
literally true — worth doing only if something other than tests needs it.

**The API's 400 responses are unproven.**  
`UpdateTask` and `CreateTask` reject an out-of-vocabulary status or priority,
but no test exercises the HTTP layer: that needs a running MySQL and a real
token. The tests cover the refusal message and the predicate boundary only.
Deciding whether this repo wants an HTTP-level test harness (httptest plus a
throwaway database) is the open question.

**Nothing runs the suite automatically.**  
`test.sh` runs both halves in one command, but no CI invokes it, so a red suite
reaches `main` unnoticed. The repo has no CI at all.

## Repo hygiene

**The agent instructions reach nobody.**  
`CLAUDE.md` and `.claude/` are both in `.gitignore` (lines 225–226), so the
project conventions, the skills and the wrap script exist on one machine only. A
second contributor — human or agent — inherits none of it. Tracking them is a
publishing decision, not a technical one.

**The repository name is misspelled.**  
The remote is `FP_GO_APP_GHARBI_YASSINE_NAMAN_KUMAR`, which GitHub redirects to
`FP_GO_APP_Task_Manger_...` — "Manger", not "Manager". The same misspelling is
baked into the module path in `go.mod` and into `REPO_URL` in
`frontend/src/lib/constants.js`. Everything resolves today. Renaming the module
is invasive; leaving it means living with a typo in every import path.

**A build artifact is committed.**  
`FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR.exe` is tracked in git. Almost
certainly unintentional.

**`docs/agents/*.md` fail the wrap check.**  
They wrap earlier than 80 columns rather than wrongly, so
`wrap-markdown.mjs --check` reports them while the only change would be cosmetic
re-filling of skill-installed files. Either normalise them or accept that
`--check` is not clean across the repo.
