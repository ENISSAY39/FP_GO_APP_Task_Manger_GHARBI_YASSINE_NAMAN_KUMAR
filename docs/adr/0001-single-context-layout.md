# ADR 0001 — Single-context layout for domain documentation

**Status**: accepted, 2026-09-01

## Context

The repo needed one agreed home for domain vocabulary and for the decisions that
shape it. Two layouts were available:

- **Single context**: one `CONTEXT.md` at the root, with `docs/adr/` beside it,
  covering the whole codebase.
- **Multi context**: a `CONTEXT-MAP.md` at the root pointing at one `CONTEXT.md`
  per bounded context, each with its own `docs/adr/`.

The second exists for codebases where separate areas genuinely speak different
languages, and where the same word means different things in two places.

This codebase is one application: a Go API and a React frontend that share a
single vocabulary end to end. A Task on the server is the Task on the card. The
reminder work made that concrete — the Overdue and Due soon rules are defined
once and consumed by three UI surfaces, and the status vocabulary is one list
that both sides must spell identically.

## Decision

Single-context layout. One `CONTEXT.md` at the repo root holds the glossary for
the whole project; `docs/adr/` holds decisions that affect it.

Both halves of the stack use the same terms. When the Go constant and the
frontend constant disagree, that is a bug to fix (see ADR 0002 and issue #2),
not two contexts to separate.

## Consequences

- One place to look, and one place to change. An agent or a contributor reads
  `CONTEXT.md` before exploring, and nothing points elsewhere.
- The glossary must stay small enough to stay read. Terms earn their place by
  being ambiguous or contested, not by merely existing.
- If the project ever grows an area with its own language — billing, say, or a
  reporting pipeline — this decision is the one to revisit, by introducing
  `CONTEXT-MAP.md` rather than by letting one glossary carry two meanings for a
  word.
- The vocabulary is enforceable because it is single: the server rejects a
  status outside the glossary's three values, so the written definition and the
  runtime behaviour cannot drift apart silently.
