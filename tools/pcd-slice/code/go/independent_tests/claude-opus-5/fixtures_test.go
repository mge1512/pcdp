// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
// tests by: claude-opus-5
//
// Specification fixtures. Each fixture is a structurally complete input for
// the grammar pcd-slice recognises: a title, preamble sections, a TYPES
// block where the test needs one, and at least one BEHAVIOR heading. Only
// fixtures whose test purpose is a structural finding are deliberately
// incomplete, and those are named accordingly.
//
// "@@@" stands for a markdown code fence (see fixture() in helpers_test.go).
package pcdslice_test

// specClosure is the fixture of EXAMPLE: closure_is_transitive. One behavior
// "record" whose block mentions Plan and nothing else; Unused is defined but
// referenced by nobody.
const specClosure = `# demo

## META
Deployment:  cli-tool
Version:     0.1.0

Intro prose for the demo specification.

## TYPES

@@@
Plan := { items: list of Item }
Item := { path: Path }
Path := string
Unused := int
@@@

## BEHAVIOR: record
Constraint: required

Records a Plan into the store.

STEPS:
1. Store it.
`

// specFull is a complete two-behavior specification with a global invariant,
// a bound invariant, a nested example, a top-level example and a changelog.
const specFull = `# demo

## META
Deployment:  cli-tool
Version:     0.1.0

Intro prose for the demo specification.

## TYPES

@@@
Plan := { items: list of Item }
Item := { path: Path }
Path := string
Unused := int
@@@

## BEHAVIOR: record
Constraint: required

Records a Plan into the store.

STEPS:
1. Store the plan.

### EXAMPLE record_writes_plan
GIVEN:
  an empty store
THEN:
  the plan is on disk

## BEHAVIOR: purge
Constraint: required

Removes stored entries again.

STEPS:
1. Remove them.

## INVARIANTS
- The store is never copied back.
- Purge destroys nothing unfetched. (binds: purge)

## EXAMPLES

### EXAMPLE: record_accepts_plan
GIVEN:
  a plan with one item
WHEN:
  result = record(plan)
THEN:
  the plan is stored

## DEPLOYMENT

Runtime: command-line tool.

## Changelog

| version | note        |
|---------|-------------|
| 0.1.0   | first cut   |
| 0.0.9   | draft only  |
`

// hintsFull carries one attributable block ("record") and two blocks that
// name no behavior; the latter are hints preamble, not findings.
const hintsFull = `# demo hints

General guidance that belongs to no behavior at all.

## Notes on record

Records should be written in one pass.

## Other guidance

Nothing behavior specific here.
`

// hintsSecond is a second hints file, used to prove that hints entries are
// recorded in input order in the provenance line.
const hintsSecond = `# second hints

Second file guidance for everyone.

## More about purge

Purging twice is allowed.
`

// specGlobalInvariant is the fixture of EXAMPLE: untagged_invariant_is_global.
const specGlobalInvariant = `# demo

## META
Version: 0.1.0

## TYPES

@@@
Path := string
@@@

## BEHAVIOR: fetch
Constraint: required

Fetches things.

## BEHAVIOR: store
Constraint: required

Stores things.

## INVARIANTS
- The store is never copied back.
`

// specBoundInvariant is the fixture of EXAMPLE:
// bound_invariant_lands_in_named_bundles.
const specBoundInvariant = `# demo

## META
Version: 0.1.0

## BEHAVIOR: reset
Constraint: required

Resets the working area.

## BEHAVIOR: purge
Constraint: required

Purges the cache.

## INVARIANTS
- Reset destroys nothing unfetched. (binds: reset)
`

// specUnknownBinding is the fixture of EXAMPLE:
// binding_to_unknown_behavior_is_a_finding.
const specUnknownBinding = `# demo

## META
Version: 0.1.0

## BEHAVIOR: rename
Constraint: required

Renames a thing.

## INVARIANTS
- X holds. (binds: renmae)
`

// specNestedExample is the fixture of EXAMPLE:
// nested_example_travels_with_its_behavior.
const specNestedExample = `# demo

## META
Version: 0.1.0

## BEHAVIOR: verify
Constraint: required

Verifies the head revision.

### EXAMPLE verify_reports_head
GIVEN:
  a repository at revision 7
WHEN:
  the head is verified
THEN:
  the report names revision 7

## BEHAVIOR: publish
Constraint: required

Publishes the result.
`

// specTopLevelExample is the fixture of EXAMPLE:
// top_level_example_attributed_by_when_line.
const specTopLevelExample = `# demo

## META
Version: 0.1.0

## BEHAVIOR: lint
Constraint: required

Lints a specification file.

## BEHAVIOR: report
Constraint: required

Writes a report.

## EXAMPLES

### EXAMPLE: lint_rejects_missing_meta
GIVEN:
  a specification without a META section
WHEN:
  result = lint(file, strict=false)
THEN:
  the finding names the missing section
`

// specDuplicateBehavior is the fixture of EXAMPLE:
// slice_refuses_on_findings_and_writes_nothing. It is deliberately faulty:
// the same behavior heading appears twice.
const specDuplicateBehavior = `# demo

## META
Version: 0.1.0

## BEHAVIOR: emit
Constraint: required

Emits once.

## BEHAVIOR: emit
Constraint: required

Emits again.
`

// specNoBehaviors is deliberately without any behavior heading.
const specNoBehaviors = `# demo

## META
Version: 0.1.0

## TYPES

@@@
Path := string
@@@

## DEPLOYMENT

Runtime: nothing at all.
`

// specUndefinedType is deliberately faulty: a definition references a type
// that is defined nowhere.
const specUndefinedType = `# demo

## META
Version: 0.1.0

## TYPES

@@@
Plan := { items: list of MissingItem }
Path := string
@@@

## BEHAVIOR: record
Constraint: required

Records a Plan.
`

// specUnassignableExample is deliberately faulty: the top-level example can
// be attributed to no behavior.
const specUnassignableExample = `# demo

## META
Version: 0.1.0

## BEHAVIOR: record
Constraint: required

Records things.

## EXAMPLES

### EXAMPLE: two_runs_agree
GIVEN:
  two runs
WHEN:
  both runs happen
THEN:
  they agree
`

// hintsUnassignable is deliberately faulty: a heading of behavior shape that
// names no known behavior.
const hintsUnassignable = `# demo hints

Preamble guidance.

## renmae

This block names a behavior that does not exist.
`

// specCompleteness exercises the completeness invariant. Every non-blank
// line is unique, no invariant binds more than one behavior, and a changelog
// closes the file.
const specCompleteness = `# completeness demo

## META
Version: 0.2.0

Prose introducing the completeness demo specification.

## TYPES

@@@
Alpha := { beta: Beta }
Beta := string
Gamma := int
@@@

## INTERFACES

An interface paragraph that belongs to the preamble only.

## BEHAVIOR: alpha-verb
Constraint: required

Consumes an Alpha value and returns nothing of interest.

STEPS:
1. Consume the value quietly.

## BEHAVIOR: beta-verb
Constraint: required

Does something entirely different with no declared types.

STEPS:
1. Do the different thing.

## INVARIANTS
- Global truth number one about the whole component.
- Bound truth about the first verb only. (binds: alpha-verb)

## EXAMPLES

### EXAMPLE: alpha_verb_consumes
GIVEN:
  an alpha value
WHEN:
  result = alpha-verb(value)
THEN:
  nothing observable happens

## Changelog

| 0.2.0 | changelog row one that must never be emitted |
| 0.1.0 | changelog row two that must never be emitted |
`

// specRefinementPredicates uses the refinement-predicate style of the
// specification's own TYPES section. The all-capital connectives of a
// predicate are not type references and must not be reported.
const specRefinementPredicates = `# demo

## META
Version: 0.1.0

## TYPES

@@@
SpecFile := path where file_exists AND readable AND extension = ".md"
HintsFile := path where file_exists AND readable AND extension = ".md"
OutDir := path
// Directory the slice verb writes into. Created if absent.
@@@

## BEHAVIOR: slice
Constraint: required

Reads a SpecFile plus every HintsFile and writes into OutDir.
`

// specUndefinedTypeTwice references the same undefined name twice on one
// line; that is one finding, not two.
const specUndefinedTypeTwice = `# demo

## META
Version: 0.1.0

## TYPES

@@@
Pair := { left: MissingItem, right: MissingItem }
@@@

## BEHAVIOR: record
Constraint: required

Records a Pair.
`
