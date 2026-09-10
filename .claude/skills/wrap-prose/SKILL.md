---
name: wrap-prose
description: Wrap or unwrap Markdown prose for its destination before sending it. Use before writing a commit body, a GitHub issue/PR body or comment, a repo .md file, or a message for Teams/Slack/Discord.
---

# Wrap prose for its destination

Prose has to be folded to the width of wherever it is going. Get it wrong and it
renders badly in a way you cannot see from the terminal, because the terminal
soft-wraps everything and hides the difference.

Two opposite failure modes:

- **Pre-wrapped text in a hard-break renderer** (GitHub issue/PR bodies and
  comments, Teams, Slack): every newline becomes a visible `<br>`, so the text
  renders as a narrow staircase down the left of a wide container.
- **Unwrapped text in a file that gets diffed** (any `.md` in the repo): one
  paragraph is one enormous logical line, so every edit rewrites the whole line
  and the diff is unreadable.

## Widths by destination

| Destination | Format | Why |
| --- | --- | --- |
| `.md` file in the repo | 80 columns | Rendered as CommonMark (single newlines are soft), and the source is read in diffs |
| GitHub issue / PR **body or comment** | **unwrapped** | GitHub enables hard line breaks in comment fields: pre-wrapped text renders as a staircase |
| Commit message body | 72 columns | `git log` indents by 4 (72 + 4 = 76); subject ≤ 50 |
| Teams, Slack, Discord | unwrapped | These clients reflow themselves; pre-wrapped text arrives as a staircase |

The distinction that matters most: a **`.md` file in the repo** and a **GitHub
comment** are both "Markdown on GitHub" but use different renderers. Files get
soft breaks, comment fields get hard breaks. Wrap the first, never the second.

## How

Write the prose to a file, transform it, then hand the **file** to the tool that
sends it. Never dictate a block of prose for the user to copy out of the
terminal: what is displayed is a single logical line that the terminal merely
displays folded, so copying it back produces text folded at the window width, or
not folded at all.

```bash
node scripts/wrap-markdown.mjs body.md              # 80 (default)
node scripts/wrap-markdown.mjs --width 72 body.md   # commit body
node scripts/wrap-markdown.mjs --unwrap body.md     # GitHub / chat
node scripts/wrap-markdown.mjs --check *.md         # exit 1 if not wrapped
node scripts/wrap-markdown.mjs --stdout body.md     # preview, don't write
```

Then:

```bash
git commit -F body.md                        # after --width 72
gh pr create --body-file body.md             # after --unwrap
gh issue create --body-file body.md          # after --unwrap
gh issue comment 78 --body-file body.md      # after --unwrap
```

Wrapping is reversible, so the same source can go to two destinations: unwrap it
for the issue body, wrap it at 80 for the file in the repo.

## What is enforced, and what is not

CI runs `wrap-markdown.mjs --check` over the repo's tracked `.md` files, so a
file that drifts fails the build. That is the only automated net, and it reaches
only files git tracks.

Everything sent rather than committed — a commit body, an issue or PR body, a
chat message — is outside it, because none of those is ever a tracked file. A
`PreToolUse` hook was tried for exactly that case and removed: it never fired,
and it lived in an ignored directory, so it protected nobody else anyway. Write
the prose to a file and transform it yourself.

## What the tool never touches

Fenced and indented code, tables, headings, front matter, link definitions,
horizontal rules, HTML blocks. It never breaks a URL, a `[text](url)` link, or a
code span across lines — an oversized one gets its own line and is allowed to
overflow.

## Checking your work

`--check` exits 1 and names the files that are not already in the target form,
so it works in a pre-commit hook or CI over the repo's `.md` files. To see what
would change without writing, use `--stdout`.
