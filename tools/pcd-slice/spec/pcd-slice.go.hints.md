# pcd-slice - Go realization hints

Companion to pcd-slice.spec.md for a Go translation. The generic
cli-tool.go.milestones.hints.md applies as usual; this file carries
only what is specific to this tool. Standard library only - no
third-party module earns its place here.

## 1. One pass, one classifier

Read each source once into memory (os.ReadFile; the inputs are
specifications, not gigabytes) and split on "\n" without normalizing
first - see section 6. Walk lines with a small state machine holding
exactly: inFence (bool, toggled by lines whose trimmed prefix is three
backticks), currentSection (the last column-0 "## " heading), and
currentBehavior (set inside behavior blocks, empty elsewhere). Every
structural decision in the spec - heading, fence, definition line,
binding suffix - is a prefix or suffix test on one line plus that
state. Do not reach for a markdown parser: the grammar section of the
spec is the whole grammar, and a parser library would recognize more
than the spec admits.

## 2. Assignment before output

Materialize the completeness invariant directly: build
`assign []Assignment`, one entry per source line, initialized to a
sentinel. The section walk stamps every line. After the walk, any
sentinel left is the spec's own "unassigned" finding - the invariant
becomes a slice scan instead of an argument. Outputs are then pure
projections of assign: preamble.md is the lines stamped preamble, in
order; each bundle is the lines stamped with its behavior, in order,
plus the appended attributed material. Projection from one array keeps
"exactly once" true by construction.

## 3. Type closure

Collect definitions first: inside the TYPES fence, a definition line
matches `^([A-Za-z][A-Za-z0-9_]*) :=`. Attach following comment and
continuation lines to the last definition. Build
`defs map[string]string` (name to full definition text) and
`names []string` sorted longest-first for matching. Reference
detection is whole-word: compile one `\b(A|B|C)\b` alternation from
the sorted names (regexp.QuoteMeta each; longest-first keeps Plan
from shadowing PlanPath). Closure is a worklist over the behavior
block's text and then over each newly added definition's text.
Alphabetize only at output time.

## 4. Attribution

Top-level examples: the WHEN line's invoked identifier is
`^\s*(?:result\s*=\s*)?([a-z][a-z0-9_-]*)\(` - match group against
behavior names; on no match, fall back to whole-word behavior name in
the example's own name; on no match again, finding. Hints headings:
test each behavior name as a whole word against the heading text;
bold markers and backticks are stripped before the test
(strings.ReplaceAll for `**` and the backtick, nothing cleverer).
"Behavior shape but unknown" for the unassignable-hints finding means:
the heading's first word after stripping matches `[a-z][a-z0-9_]*`
and is not a known behavior and not an ordinary English word - keep
this test exactly as cheap as the spec's finding needs: a heading
whose stripped text is one single lowercase identifier. Anything
looser drowns the finding in prose headings.

## 5. Determinism mechanics

- map iteration is randomized in Go: never range over a map when
  emitting. Keep behaviors in a slice in source order; sort closure
  names and manifest lines with sort.Strings.
- SHA-256 via crypto/sha256 over the raw file bytes as read - the
  same bytes the provenance header names. hex.EncodeToString for
  lowercase hex.
- The provenance line is one Sprintf with a fixed field order; the
  tool version is a package-level constant set at build time via
  -ldflags "-X main.version=", defaulting to "dev". A "dev" build
  still emits deterministic bytes for one binary; the acceptance
  examples compare two runs of the same binary, not two builds.

## 6. Line endings and the completeness count

Split on "\n" and keep any trailing "\r" attached to the line for
assignment and for verbatim projection - emitting exactly the bytes
that came in preserves determinism and keeps a CRLF source's bundles
diffable against it. The one normalization allowed: the provenance
line and the headings pcd-slice itself adds end in "\n" regardless of
source flavor. A final line without newline is one line; len(lines)
after split needs the usual empty-tail correction.

## 7. Atomic out-directory handling

slice computes everything, then touches disk. Write into the target
via create-then-rename per file only if simplicity demands it; the
spec's atomicity is per-run, not per-file: on any write error, remove
the files this run created (track them in a slice) and leave earlier
content alone. The foreign-file refusal (exit 2) runs before the
first write: read the first line of every existing file in out;
provenance-marked files are ours to replace or remove, anything else
aborts. os.MkdirAll for out; check writability by creating the
manifest last - it doubles as the commit marker, and check treats a
directory without MANIFEST.tsv as stale in full.

## 8. Errors and exits

Match the house: findings to stderr one per line in the spec's format,
summary to stdout, exit via a single os.Exit at the top of main after
run() returns (int, error) - no os.Exit below main, or deferred
cleanup silently stops running. Exit 2 messages start "error: " and
name the path; exit 1 is reserved for findings.

## 9. Tests from the EXAMPLES

Every EXAMPLE block is one table-driven test with a fixture spec
under testdata/. The two-run determinism example compares
directory trees byte-wise (read both, map relpath to sha256, compare
maps). The closure example asserts the Requires-Types line verbatim.
Golden files for preamble.md and one bundle keep refactors honest;
regenerate goldens only through the tool itself and diff-review the
change.

## 10. What not to build

No configuration file, no environment variables, no color, no
concurrency - the inputs are a few hundred kilobytes and one pass is
microseconds. No YAML front matter parsing: the provenance line is an
HTML comment precisely so bundles stay plain markdown everywhere.
