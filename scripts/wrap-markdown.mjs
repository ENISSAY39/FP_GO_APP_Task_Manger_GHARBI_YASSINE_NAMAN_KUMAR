#!/usr/bin/env node
/**
 * Markdown prose wrapper / unwrapper.
 *
 *   node .claude/hooks/wrap-markdown.mjs file.md              # wrap at 80
 *   node .claude/hooks/wrap-markdown.mjs --width 72 file.md   # commit body
 *   node .claude/hooks/wrap-markdown.mjs --unwrap file.md     # one line/paragraph
 *   node .claude/hooks/wrap-markdown.mjs --check *.md         # exit 1 if not wrapped
 *   node .claude/hooks/wrap-markdown.mjs --stdout file.md     # preview, don't write
 *
 * Wrapping is reversible: --unwrap rejoins each paragraph into a single
 * logical line, so text can move between destinations that want different
 * widths without loss.
 *
 * Never touched: fenced and indented code, tables, headings, front matter,
 * link definitions, horizontal rules, HTML blocks. Never broken across a
 * line: URLs, `[text](url)` links, and code spans — those stay whole even
 * when they overflow the target width.
 */

import { readFileSync, writeFileSync } from 'node:fs'

const DEFAULT_WIDTH = 80

// --------------------------------------------------------------- line kinds

const isBlank = (line) => /^\s*$/.test(line)
const isFence = (line) => /^\s{0,3}(```|~~~)/.test(line)
const isHeading = (line) => /^\s{0,3}#{1,6}(\s|$)/.test(line)
const isTableRow = (line) => /^\s*\|/.test(line)
const isLinkDefinition = (line) => /^\s{0,3}\[[^\]]+\]:\s/.test(line)
const isHorizontalRule = (line) => /^\s{0,3}([-*_])(\s*\1){2,}\s*$/.test(line)
const isHtmlBlock = (line) => /^\s{0,3}<\/?[a-zA-Z!]/.test(line)
const isIndentedCode = (line) => /^(\s{4,}|\t)/.test(line)

/** A line that markdown renders as a forced break must keep its break. */
const endsWithHardBreak = (line) => /(\s{2,}|\\)$/.test(line)

/** Structural lines are emitted verbatim and interrupt a paragraph. */
function isStructural(line) {
  return (
    isFence(line) ||
    isHeading(line) ||
    isTableRow(line) ||
    isLinkDefinition(line) ||
    isHorizontalRule(line) ||
    isHtmlBlock(line)
  )
}

/**
 * Split a prose line into the prefix that must be reprinted (blockquote
 * markers, list bullet) and the text to reflow. Continuation lines keep the
 * quote markers but replace the bullet with matching spaces, so a wrapped
 * list item stays hanging-indented.
 */
function splitPrefix(line) {
  const match = /^(\s*)((?:>\s?)*)([-*+]\s+|\d+[.)]\s+)?([\s\S]*)$/.exec(line)
  const [, indent, quote, bullet = '', text] = match
  return {
    firstPrefix: indent + quote + bullet,
    contPrefix: indent + quote + ' '.repeat(bullet.length),
    hasBullet: bullet.length > 0,
    text,
  }
}

// --------------------------------------------------------------- tokenising

/**
 * Split text into tokens on whitespace, except inside an atom that must not
 * be broken: a code span, or the `[text](url)` of a link or image. Those
 * keep their internal spaces and travel as one token.
 */
function tokenise(text) {
  const tokens = []
  let current = ''
  let backticks = 0 // length of the run that opened the current code span
  let inLinkText = false
  let inLinkUrl = false
  let parenDepth = 0

  for (let i = 0; i < text.length; i += 1) {
    const char = text[i]

    if (char === '`') {
      let run = 0
      while (text[i + run] === '`') run += 1
      if (backticks === 0) backticks = run
      else if (backticks === run) backticks = 0
      current += text.slice(i, i + run)
      i += run - 1
      continue
    }

    if (backticks === 0) {
      if (char === '[') inLinkText = true
      else if (char === ']') inLinkText = false
      else if (char === '(' && text[i - 1] === ']') {
        inLinkUrl = true
        parenDepth = 1
      } else if (inLinkUrl && char === '(') parenDepth += 1
      else if (inLinkUrl && char === ')') {
        parenDepth -= 1
        if (parenDepth === 0) inLinkUrl = false
      }
    }

    const protectedHere = backticks > 0 || inLinkText || inLinkUrl
    if (/\s/.test(char) && !protectedHere) {
      if (current) tokens.push(current)
      current = ''
      continue
    }

    current += char
  }

  if (current) tokens.push(current)
  return tokens
}

function fillTokens(tokens, width, firstPrefix, contPrefix) {
  const lines = []
  let line = firstPrefix
  let empty = true

  for (const token of tokens) {
    const candidate = empty ? line + token : `${line} ${token}`
    if (!empty && candidate.length > width) {
      lines.push(line)
      line = contPrefix + token // an oversized token still gets its own line
      empty = false
    } else {
      line = candidate
      empty = false
    }
  }

  if (!empty || lines.length === 0) lines.push(line)
  return lines
}

// ------------------------------------------------------------- transforming

/**
 * Walk the document, collect runs of prose into paragraph units, and hand
 * each unit to `render`. Everything else is passed through untouched.
 */
function transform(source, render) {
  const lines = source.split('\n')
  const out = []
  let index = 0
  let inFence = false
  let fenceMarker = ''
  // CommonMark only opens an indented code block after a blank line; an
  // indented line under a paragraph is a continuation, not code.
  let previousBlank = true

  // Front matter, only when it opens the file.
  if (lines[0] === '---') {
    const close = lines.indexOf('---', 1)
    if (close !== -1) {
      out.push(...lines.slice(0, close + 1))
      index = close + 1
    }
  }

  while (index < lines.length) {
    const line = lines[index]

    if (inFence) {
      out.push(line)
      if (isFence(line) && line.trim().startsWith(fenceMarker)) inFence = false
      index += 1
      continue
    }

    if (isFence(line)) {
      inFence = true
      fenceMarker = line.trim().slice(0, 3)
      out.push(line)
      previousBlank = false
      index += 1
      continue
    }

    if (isBlank(line) || isStructural(line) || (previousBlank && isIndentedCode(line))) {
      out.push(line)
      previousBlank = isBlank(line)
      index += 1
      continue
    }

    // A paragraph unit: this line plus the continuation lines under it.
    const unit = [line]
    index += 1
    while (index < lines.length) {
      const next = lines[index]
      if (
        isBlank(next) ||
        isStructural(next) ||
        isFence(next) ||
        splitPrefix(next).hasBullet ||
        endsWithHardBreak(unit[unit.length - 1])
      ) {
        break
      }
      unit.push(next)
      index += 1
    }

    out.push(...render(unit))
    previousBlank = false
  }

  return out.join('\n')
}

/**
 * The reflowable text of a paragraph unit. Every line is stripped of its own
 * prefix — indent, blockquote markers, bullet — so a wrapped blockquote does
 * not drag its `>` markers into the middle of the joined text.
 */
function unitText(unit) {
  return unit
    .map((line) => splitPrefix(line).text.trim())
    .join(' ')
    .trim()
}

function wrapDocument(source, width) {
  return transform(source, (unit) => {
    const { firstPrefix, contPrefix } = splitPrefix(unit[0])
    const text = unitText(unit)
    if (!text) return unit

    const lines = fillTokens(tokenise(text), width, firstPrefix, contPrefix)
    if (endsWithHardBreak(unit[unit.length - 1])) lines[lines.length - 1] += '  '
    return lines
  })
}

function unwrapDocument(source) {
  return transform(source, (unit) => {
    const { firstPrefix } = splitPrefix(unit[0])
    const text = unitText(unit)
    if (!text) return unit

    const hardBreak = endsWithHardBreak(unit[unit.length - 1]) ? '  ' : ''
    return [firstPrefix + text + hardBreak]
  })
}

// ---------------------------------------------------------------------- CLI

function parseArgs(argv) {
  const options = { width: DEFAULT_WIDTH, unwrap: false, check: false, stdout: false }
  const files = []

  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i]
    if (arg === '--unwrap') options.unwrap = true
    else if (arg === '--check') options.check = true
    else if (arg === '--stdout') options.stdout = true
    else if (arg === '--width') {
      const value = Number(argv[i + 1])
      if (!Number.isInteger(value) || value < 20) {
        console.error(`wrap-markdown: --width needs an integer >= 20, got "${argv[i + 1]}"`)
        process.exit(2)
      }
      options.width = value
      i += 1
    } else if (arg === '--help' || arg === '-h') {
      options.help = true
    } else if (arg.startsWith('-')) {
      console.error(`wrap-markdown: unknown option "${arg}"`)
      process.exit(2)
    } else {
      files.push(arg)
    }
  }

  return { options, files }
}

const USAGE = `Usage: wrap-markdown.mjs [--width N | --unwrap] [--check] [--stdout] <file.md...>

  --width N   wrap prose at N columns (default ${DEFAULT_WIDTH})
  --unwrap    rejoin each paragraph into one logical line
  --check     don't write; exit 1 if any file is not already in that form
  --stdout    print the result instead of rewriting the file`

function main() {
  const { options, files } = parseArgs(process.argv.slice(2))

  if (options.help || files.length === 0) {
    console.log(USAGE)
    process.exit(files.length === 0 && !options.help ? 2 : 0)
  }

  let changed = false

  for (const file of files) {
    let source
    try {
      source = readFileSync(file, 'utf8')
    } catch (error) {
      console.error(`wrap-markdown: cannot read ${file}: ${error.message}`)
      process.exit(2)
    }

    const result = options.unwrap ? unwrapDocument(source) : wrapDocument(source, options.width)

    if (result !== source) changed = true

    if (options.check) {
      if (result !== source) console.error(`${file}: not ${options.unwrap ? 'unwrapped' : 'wrapped'}`)
    } else if (options.stdout) {
      process.stdout.write(result)
    } else if (result !== source) {
      writeFileSync(file, result)
    }
  }

  if (options.check && changed) process.exit(1)
}

main()
