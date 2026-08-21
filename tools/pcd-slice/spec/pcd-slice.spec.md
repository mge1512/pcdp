# pcd-slice

## META
Deployment:  cli-tool
Version:     0.1.0
Spec-Schema: 0.4.0
Author:      Matthias G. Eckermann <pcd@mailbox.org>
License:     GPL-2.0-only
Verification: none
Safety-Level: QM
Module:       github.com/mge1512/pcd/tools/pcd-slice

---

pcd-slice derives per-behavior translation bundles from a PCD
specification. Translation harnesses read files in mechanical windows;
a window that cuts a behavior in half or omits the types it references
forces another window or a guess, and the specification is re-read in
full across every pass and every role. pcd-slice computes the
semantic form once, deterministically: for each behavior, one bundle
holding that behavior's block, the examples that exercise it, the
invariants bound to it, the hints written for it, and the names of the
types it needs; plus one preamble holding what every bundle shares.
The source specification is never modified and remains the durable
artifact; every emitted file carries the SHA-256 of the sources it was
derived from, so a bundle is verifiable against the reproducibility
tuple and stale bundles are detectable.

The recognized input grammar is deliberately small and stated in this
specification. It covers Spec-Schema specifications and the common
house dialect in which examples nest inside behavior blocks. pcd-slice
does not depend on libpcd in this version; the grammar here binds the
tool, and integration with libpcd's spec model is deferred until both
dialects converge.

---

## TYPES

```
SpecFile := path where file_exists AND readable AND extension = ".md"

HintsFile := path where file_exists AND readable AND extension = ".md"
// Zero or more hints files may be given; order is preserved.

OutDir := path
// Directory the slice verb writes into. Created if absent. pcd-slice
// writes nothing outside OutDir.

BehaviorName := string
// The text after "## BEHAVIOR:" or "## BEHAVIOR/INTERNAL:" on a
// heading line, trimmed. Must be unique within one specification.

TypeName := string
// The identifier left of ":=" on a definition line inside the TYPES
// section. A definition extends until the next definition line or the
// end of the fenced block; comment and continuation lines belong to
// the definition above them.

InvariantBinding := list of BehaviorName
// An invariant line may end with "(binds: name, name)". A bound
// invariant belongs to exactly the named behaviors' bundles. An
// invariant without the suffix is global and belongs to the preamble.
// The suffix is recognized only at end of line, outside code fences.

Assignment := one_of("behavior", "preamble", "excluded", "unassigned")
// Every line of the specification receives exactly one assignment.
// "excluded" covers only the changelog. "unassigned" is a finding.

Finding := { rule: string, file: path, line: int, message: string }
// rule is one_of("undefined-type", "unknown-binding",
// "unassignable-example", "unassignable-hints", "duplicate-behavior",
// "no-behaviors", "stale-output")

ExitCode := 0 | 1 | 2
// 0 = clean (check: no findings; slice: bundles written)
// 1 = findings (check found at least one; slice refused and wrote nothing)
// 2 = invocation error (bad arguments, missing file, unreadable file,
//     unwritable OutDir)
```

---

## SECTION GRAMMAR

The recognized structure of a specification, normative for this tool:

- A heading line is a line beginning with "#" at column 0 outside any
  code fence. Fences are lines beginning with three backticks; fenced
  content is never interpreted as structure.
- `## TYPES` opens the type section; its fenced block holds type
  definition lines `Name := ...` with attached comment and
  continuation lines as defined under TypeName.
- `## BEHAVIOR: <name>` and `## BEHAVIOR/INTERNAL: <name>` open a
  behavior block, which extends to the next `## ` heading. Everything
  inside, including `### ` subsections and nested `### EXAMPLE`
  blocks, belongs to the behavior.
- `## INVARIANTS` holds one invariant per list item; bindings as
  defined under InvariantBinding.
- `## EXAMPLES` (top-level, Spec-Schema form) holds `### EXAMPLE:
  <name>` blocks. An example is attributed to the behavior whose name
  appears as the invoked identifier on its `WHEN:` line, or, failing
  that, as a whole word in the example's name. An example attributable
  to no behavior is a finding.
- A changelog is the section whose `## ` heading contains the word
  "Changelog" in any case, extending to end of file. It is excluded
  from every output and exempt from the completeness rule.
- Every other `## ` section (META, PURPOSE and untitled intro prose,
  INTERFACES, PRECONDITIONS, POSTCONDITIONS, DEPENDENCIES, DEPLOYMENT,
  MILESTONE sections, and any heading this grammar does not name) is
  preamble.

In a hints file, a `## ` or deeper heading whose text contains a
BehaviorName as a whole word opens a block attributed to that
behavior, extending to the next heading of the same or shallower
depth. Hints content before the first attributable heading, and any
block that names no behavior is hints preamble. A heading that looks like
an attribution but names no known behavior is a finding.

---

## BEHAVIOR: list
Constraint: required

Prints the structure the other verbs will act on: each behavior with
its type closure, bound invariants, attributed examples, and
attributed hints blocks.

INPUTS:
```
spec:  SpecFile
hints: list of HintsFile   // zero or more, via repeatable hints=
```

OUTPUTS:
```
lines: one per behavior, tab-separated:
       name, types, invariants, examples, hints_blocks
       (counts for the last four; names on demand are v2)
```

STEPS:
1. Parse spec against the section grammar; on unreadable input ->
   exit 2 with "error: cannot read {path}".
2. Compute, per behavior, the transitive type closure: every TypeName
   whose whole-word form appears in the behavior block or in the
   definition of a type already in the closure.
3. Attribute invariants, examples and hints blocks as the grammar
   defines.
4. Print one line per behavior in specification order; print a final
   line "total: {n} behaviors" to stdout.
5. Exit 0. list never fails on findings; it reports structure as-is.

POSTCONDITIONS:
- no file is created or modified
- output order equals the order of behavior headings in the spec

SIDE-EFFECTS:
- stdout: the listing
- no network calls

---

## BEHAVIOR: check
Constraint: required

Validates that the specification and hints can be sliced completely
and unambiguously, and, when out= is given, that existing output is
current.

INPUTS:
```
spec:  SpecFile
hints: list of HintsFile
out:   OutDir or absent
```

OUTPUTS:
```
findings: list of Finding, written to stderr, one line each:
          "{rule}: {file}:{line}: {message}"
summary:  one line to stdout:
          "pcd-slice check: {n} findings" or "pcd-slice check: clean"
```

STEPS:
1. Parse as in list; exit 2 on unreadable input.
2. undefined-type: a whole-word occurrence of an identifier of type
   shape is not checked; only the reverse holds - every TypeName in a
   computed closure must have a definition. A behavior referencing a
   name that matches no definition is not detectable without a
   language model and is out of scope; this rule fires only when a
   type definition references a TypeName that is defined nowhere.
3. unknown-binding: an invariant binding names no known behavior.
4. unassignable-example: a top-level example attributable to no
   behavior.
5. unassignable-hints: a hints heading that names something of
   behavior shape but no known behavior. Hints preamble is not a
   finding.
6. duplicate-behavior: two behavior headings with the same name.
7. no-behaviors: the spec contains no behavior heading.
8. Completeness: assign every spec line; any line assigned
   "unassigned" is an internal fault of this tool, reported as a
   finding with rule "unassigned" - the grammar above admits none.
9. stale-output, only when out= is given and exists: recompute the
   provenance header values for the current sources and compare with
   the headers of the files present; any difference, and any expected
   file missing or unexpected file present, is one finding that names the
   file.
10. Write findings sorted by file then line; write the summary; exit 1
    if any finding, else 0.

POSTCONDITIONS:
- no file is created or modified
- exit_code = 0 iff findings is empty

SIDE-EFFECTS:
- stderr: finding lines
- stdout: summary line
- no network calls

---

## BEHAVIOR: slice
Constraint: required

Writes the derived form: one preamble, one bundle per behavior, one
manifest.

INPUTS:
```
spec:  SpecFile
hints: list of HintsFile
out:   OutDir
```

OUTPUTS:
```
out/preamble.md          shared context: META or version header line,
                         intro prose and every preamble section in
                         source order, the full TYPES section, global
                         invariants, hints preamble
out/<behavior>.md        one per behavior: provenance header,
                         "Requires-Types:" line listing the closure
                         alphabetically, the behavior block verbatim,
                         attributed top-level examples, bound
                         invariants, attributed hints blocks, each
                         under a heading that names its source
out/MANIFEST.tsv         one line per emitted file: name, sha256;
                         sorted by name; the manifest lists every
                         file except itself
```

STEPS:
1. Run check (without out=); if it finds anything -> print its
   findings, write nothing, exit 1.
2. Compute all outputs in memory. Behavior file names are the
   BehaviorName with "/" replaced by "-"; a collision after
   replacement is a duplicate-behavior finding.
3. Every emitted markdown file begins with one provenance line:
   `<!-- pcd-slice source={spec} source-sha256={h} hints={paths}
   hints-sha256={hashes} version={tool} -->` - paths as given,
   hashes lowercase hex, hints entries comma-separated in input
   order, no timestamp.
4. Byte content is fully determined by the inputs: two runs over
   identical inputs produce identical files. Ordering rules: sections
   in source order, closures alphabetical, manifest sorted.
5. Create OutDir if absent. Remove files the manifest would not list
   only if they carry a pcd-slice provenance line; refuse with exit 2
   if OutDir contains a foreign file, and name the file.
6. Write every file, then the manifest. On any write error -> exit 2;
   files already written in this run are removed.
7. Print to stdout: "pcd-slice: {n} bundles, preamble, manifest ->
   {out}". Exit 0.

POSTCONDITIONS:
- out contains exactly the manifest's files plus MANIFEST.tsv
- each file's sha256 equals its manifest entry
- the concatenation of all assignments covers every spec line outside
  the changelog exactly once, across preamble and bundles
- the source specification and hints files are byte-identical to
  before the run

SIDE-EFFECTS:
- filesystem writes under out only
- stdout: one summary line
- no network calls

---

## PRECONDITIONS
- spec exists, is readable, and ends in ".md"
- every hints path given exists, is readable, and ends in ".md"
- for slice: out is creatable or writable

## POSTCONDITIONS
- list and check never write files
- slice writes only under out, atomically as defined in its steps
- exit codes follow ExitCode for every verb

## INVARIANTS
- Sources are read-only: no verb modifies spec or hints. (binds: list, check, slice)
- Determinism: identical inputs produce byte-identical outputs;
  no output contains a timestamp, hostname, user name, or absolute
  path other than those given on the command line. (binds: slice)
- Completeness: every specification line outside the changelog is
  assigned to exactly one output; the changelog appears in no output.
  (binds: slice)
- Provenance: every emitted markdown file's first line is the
  provenance comment, and stale outputs are detectable from it alone.
  (binds: check, slice)
- A behavior appears in exactly one bundle, and its bundle is
  self-locating: the bundle plus the preamble suffice to translate the
  behavior without opening the source specification. (binds: slice)

## EXAMPLES

### EXAMPLE: closure_is_transitive
GIVEN:
  TYPES defines: `Plan := { items: list of Item }` and
    `Item := { path: Path }` and `Path := string` and `Unused := int`
  one behavior "record" whose block mentions Plan and nothing else
  invocation: pcd-slice list spec=spec.md
WHEN:
  result = list(spec)
THEN:
  the line for "record" counts 3 types (Plan, Item, Path)
  Unused is in no closure and remains in the preamble's TYPES section
  exit_code = 0

### EXAMPLE: untagged_invariant_is_global
GIVEN:
  INVARIANTS holds "- The store is never copied back." with no binding
  two behaviors exist
  invocation: pcd-slice slice spec=spec.md out=spec.d
WHEN:
  result = slice(spec, out)
THEN:
  the invariant appears in preamble.md and in no bundle
  exit_code = 0

### EXAMPLE: bound_invariant_lands_in_named_bundles
GIVEN:
  INVARIANTS holds "- Reset destroys nothing unfetched. (binds: reset)"
  behaviors "reset" and "purge" exist
WHEN:
  result = slice(spec, out)
THEN:
  the invariant appears in reset.md, not in purge.md, not in preamble.md

### EXAMPLE: binding_to_unknown_behavior_is_a_finding
GIVEN:
  INVARIANTS holds "- X holds. (binds: renmae)" and no behavior "renmae"
  invocation: pcd-slice check spec=spec.md
WHEN:
  result = check(spec)
THEN:
  stderr contains "unknown-binding: spec.md:{line}: renmae"
  stdout = "pcd-slice check: 1 findings"
  exit_code = 1

### EXAMPLE: nested_example_travels_with_its_behavior
GIVEN:
  behavior "verify" contains "### EXAMPLE verify_reports_head"
WHEN:
  result = slice(spec, out)
THEN:
  the example is inside verify.md because it is inside the block
  no attribution rule is consulted

### EXAMPLE: top_level_example_attributed_by_when_line
GIVEN:
  "## EXAMPLES" holds "### EXAMPLE: lint_rejects_missing_meta" whose
    WHEN line reads "result = lint(file, strict=false)"
  behavior "lint" exists
WHEN:
  result = slice(spec, out)
THEN:
  the example is appended to lint.md under a heading that names its source

### EXAMPLE: changelog_is_excluded
GIVEN:
  the spec ends with "## Changelog" and forty table rows
WHEN:
  result = slice(spec, out)
THEN:
  no output file contains any changelog row
  the completeness rule counts those lines as excluded, not unassigned

### EXAMPLE: slice_refuses_on_findings_and_writes_nothing
GIVEN:
  the spec has a duplicate behavior heading
  out=spec.d does not exist
WHEN:
  result = slice(spec, out)
THEN:
  stderr contains "duplicate-behavior"
  spec.d does not exist afterwards
  exit_code = 1

### EXAMPLE: two_runs_are_byte_identical
GIVEN:
  a valid spec and one hints file
WHEN:
  slice runs twice into two fresh directories
THEN:
  every pair of same-named files is byte-identical
  MANIFEST.tsv is byte-identical

### EXAMPLE: stale_output_is_one_finding_per_file
GIVEN:
  spec.d was sliced, then one word in the spec changed
  invocation: pcd-slice check spec=spec.md out=spec.d
WHEN:
  result = check(spec, out)
THEN:
  stderr contains one "stale-output" line per emitted file
  exit_code = 1

### EXAMPLE: foreign_file_in_out_refuses
GIVEN:
  spec.d exists and contains notes.txt without a provenance line
WHEN:
  result = slice(spec, out)
THEN:
  stderr contains "error: foreign file in out: spec.d/notes.txt"
  notes.txt is untouched
  exit_code = 2

## DEPENDENCIES

None at build time beyond the target language's standard library.
Deliberately not libpcd in this version: pcd-slice must slice both
Spec-Schema specifications and the nested-example house dialect, and
its grammar is the small one stated above rather than libpcd's full
model. Revisit when the dialects converge; the seam is LoadSpec.

## DEPLOYMENT

Runtime: command-line tool, single static binary.

Invocation:
  pcd-slice list  spec=<file> [hints=<file>]...
  pcd-slice check spec=<file> [hints=<file>]... [out=<dir>]
  pcd-slice slice spec=<file> [hints=<file>]... out=<dir>

Key=value options:
  spec=<file>   required for every verb
  hints=<file>  repeatable; order preserved and recorded in provenance
  out=<dir>     required for slice; optional for check (staleness)

The consumer contract: a translation pass for behavior X reads
out/preamble.md followed by out/X.md and needs no other file. The
preamble is byte-stable across behaviors of one slicing run, which
makes it a cacheable prefix for API-side prompt caching.

Packaging: ships like the other PCD tools; no runtime dependency, no
network access, no configuration file. Linux only in v1.
