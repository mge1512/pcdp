# TRANSLATION_REPORT.md — pcd-slice

**Spec-SHA256:** `bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04`
**Spec-SHA256 (host):** `bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04`

**Included-Specs:**

| Path | SHA256 |
|------|--------|
| *(none — the host spec declares no `Includes:`)* | — |

The host spec's META declares `Spec-Schema: 0.4.0` and no `Includes:`
directive, so the merged spec text is byte-identical to the host file and the
merged hash equals the host hash. That single hash is embedded in every
generated artefact.

**LLM-Name:** `claude-opus-5`
**Mode:** `translator`
**Run mode:** single-LLM (no `independent_tests/<other-llm>/` and no
`TEST_REPORT.md` were present in the input directory)

---

## Translation Inputs (provenance)

| Label | File | SHA256 |
|---|---|---|
| `Spec-SHA256 (merged):` | `pcd-slice.spec.md` (merged text) | `bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04` |
| `Spec-SHA256 (host):` | `pcd-slice.spec.md` | `bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04` |
| `Decisions-Hints-SHA256:` | `pcd-slice.go.hints.md` | `77fdb887623323a6c0c630661aac955fe74b32fbbd20097716024e9ca80a97a4` |
| `Milestones-Hints-SHA256:` | `cli-tool.go.milestones.hints.md` | `fe73f2ab89be21a0fd0eb37a818e586b68a338a8392a8f14725c611a2d33b249` |
| `Template-SHA256:` | `cli-tool.template.md` | `13cbe57a345e6be7f6b91f0ee3a6d5d5bfc82a957143e1668325a76bfcac13fa` |
| `Prompt-SHA256:` | `prompt.md` | `2317d61f6e5631af52943243aee40cc5dc0ed29c55984d997dcf72397930d912` |
| `Role-File-SHA256:` | `ROLE.md` | `c1a484b969f2516a1f34ae2baa73e34b8c8e38afa6cbfaba782711406c2b7c66` |
| `Style-Hints-SHA256:` | none | — |
| `Library-Hints-SHA256:` | none | — |
| `Upgrade-Brief-SHA256:` | none | — |
| `Directive-SHA256:` | none | — |

No preset hierarchy was present in this environment (`/etc/pcd/hints/`,
`.pcd/hints/`, `/usr/share/pcd/presets/`, `~/.config/pcd/presets/`,
`<project>/.pcd/` all absent), so no project or company style hints applied
and no preset overrode a template default.

---

## Tests-First-Compliance

**`yes`.**

Order of writing, verifiable from the filesystem timestamps of this run:

1. `independent_tests/claude-opus-5/helpers_test.go`
2. `independent_tests/claude-opus-5/fixtures_test.go`
3. `independent_tests/claude-opus-5/pcd_slice_test.go`
4. *structural guard checked: the directory exists and holds three test files*
5. `internal/pcd-slice/*.go`, `cmd/pcd-slice/main.go`, `go.mod`, `VERSION`
6. packaging, documentation, compile gate, this report

No implementation source file was written before the guard passed. The full
expected output of every verb — the exact list line format, the finding line
format, the summary lines, the provenance line grammar, the
`Requires-Types:` line, the manifest format — was fixed by the tests first
and the implementation was written to satisfy them.

## Continuity-Check

*Not applicable — no test-author input.* No `independent_tests/<other-llm>/`
directory and no `TEST_REPORT.md` were present in `/tmp/pcd-input`, so this
is a single-LLM run, which the prompt declares a fully supported invocation.

---

## Target language and template resolution

**Resolved language: Go** — the `cli-tool` template's `LANGUAGE` row declares
Go as the `default`, and no preset (system, user or project) was present to
override it. No deviation from the default; no language decision was taken
from the evaluation environment.

**Active MILESTONE:** none. The specification contains no `## MILESTONE:`
section, so the full spec was translated as one pass (no scaffold pass, no
deferred BEHAVIORs, no stubs). The generic scaffold-first hints file was read
in full; its non-milestone guidance (static binary with `CGO_ENABLED=0`,
signal handling in `main()`, `filepath.Join` for path construction,
format-string discipline under `go vet`, "no os.Exit below main") was applied.
The hints sections that are specific to scaffold milestones, JSON scope
wrappers, renderers and `OSCommandRunner` do not apply: pcd-slice emits
markdown and TSV, not JSON, and invokes no external command.

**Module identity:** `github.com/mge1512/pcd/tools/pcd-slice`, resolved from
**authoritative source 1**, the spec META `Module:` field. Sources 2–4 were
not consulted for a value beyond confirmation: the Go hints file names no
module, there was no pre-existing manifest in the output directory, and the
spec-title fallback was not needed. No conflict, so no halt. The identity is
propagated to `go.mod`, the import path in `cmd/pcd-slice/main.go`, the RPM
`URL:`, the DEB `Homepage:` and `Source:` fields, `debian/copyright`, the man
page HOMEPAGE section and the README.

**TYPE-BINDINGS:** the `cli-tool` template contains no `## TYPE-BINDINGS`
section, so no mechanical type mapping applied. The spec's logical types were
bound to natural Go types (see *Type bindings chosen* below).

**GENERATED-FILE-BINDINGS:** the template contains no
`## GENERATED-FILE-BINDINGS` section; no generated infrastructure filenames
were invented.

**Spec DELIVERABLES / COMPONENT entries:** the specification has no
`DELIVERABLES` section with `COMPONENT:` entries, so no component-to-filename
mapping was required. The produced file set comes from the template's
`## DELIVERABLES` table alone.

**INTERFACES:** the specification has no `## INTERFACES` section, so no
production/test-double pairs were declared or produced. The tests need none:
they are black-box tests of the binary.

**Delivery mode:** mode 1, *filesystem*. All files were written to
`/tmp/pcd-output` with the filesystem tool. No repository was available to
commit to.

---

## Deliverables produced

Derived from the template's `## DELIVERABLES` table, in the mandated delivery
order.

| OUTPUT-FORMAT | Constraint | Files produced | Status |
|---|---|---|---|
| source | required | `cmd/pcd-slice/main.go` (entry point), `internal/pcd-slice/{types,source,parse,closure,hints,check,render,verbs,write}.go`, `go.mod` | done |
| public-api | required | `## Public API Surface` in this report | done |
| build | required | `Makefile` (`build`, `test`, `install`, `clean`, `man`, `dist`, plus `vendor`, `check-fmt`) | done |
| docs | required | `README.md` | done |
| man | required | `pcd-slice.1.md`, `pcd-slice.1` (pandoc 2.18) | done |
| license | required | `LICENSE` (SPDX `GPL-2.0-only` + authoritative URL, text not reproduced) | done |
| RPM | required | `pcd-slice.spec` | done |
| DEB | required | `debian/control`, `debian/changelog`, `debian/rules`, `debian/copyright` | done |
| OCI | supported | — | **not produced**: no preset activates OCI in this run |
| PKG | supported | — | **not produced**: `PLATFORM` resolves to Linux only; macOS not declared |
| binary | supported | — | no descriptor required by the table |
| report | required | `TRANSLATION_REPORT.md` (this file) | done |
| spec-hash | required | embedded in every source header, `Makefile` `SPEC_SHA256`, RPM `# pcd-spec-sha256:`, DEB `X-PCD-Spec-SHA256:`, `pcd-slice version` output, this report | done |
| *(EXECUTION phase 4)* | required | `translation_report/translation-workflow.pikchr` | done |

Nothing else was written. In particular no `.gitignore`, no `go.sum` (the
module has no dependencies, and the template does not name a lock file), no
`CHANGELOG.md`, no CI configuration, no IDE files and no `build/` directory.
The binary `pcd-slice` at the project root is the compile gate's output at
the location `BINARY-LOCATION: project-root` mandates; `make clean` removes
it and `make dist` excludes it from the tarball.

---

## Template constraints compliance

| Key | Constraint | How it is satisfied |
|---|---|---|
| VERSION | required | `VERSION` file `0.1.0` (spec META `Version:`), single source for the binary, RPM, DEB and tarball name |
| SPEC-SCHEMA | required | `0.4.0`; include resolution implemented (no includes present) |
| AUTHOR / LICENSE | required | Matthias G. Eckermann; SPDX `GPL-2.0-only` in `LICENSE`, RPM, DEP-5 copyright |
| LANGUAGE | default Go | Go; no preset override present |
| BINARY-TYPE | static (default) | `CGO_ENABLED=0` in `Makefile`, RPM `%build` and `debian/rules`; `file ./pcd-slice` reports *statically linked* |
| PERSISTENCE | none (default) | no persistent state of any kind |
| SOURCE-PARTITIONING | modular, one-entry-one-implementation | entry point `cmd/pcd-slice/main.go` does argument parsing, error reporting and dispatch only; nine files in `internal/pcd-slice/` hold the behaviours, partitioned by domain (source loading, parsing, closure, hints, checking, rendering, verbs, disk writing, types) |
| MODULE-IDENTITY | host-specified, propagated, conflict-halts | source 1 (spec META `Module:`); propagated to every artefact; no conflict |
| PUBLIC-API-SURFACE | recorded-in-report | `## Public API Surface` below |
| BINARY-COUNT | 1 | one binary, `pcd-slice` |
| BINARY-LOCATION | project-root | `make build` emits `./pcd-slice` next to `go.mod`; tests address it as `../../pcd-slice` |
| RUNTIME-DEPS | none | standard library only; RPM/DEB declare no runtime dependency |
| CLI-ARG-STYLE | key=value (+ bare-words, subcommand) | verbs `list`/`check`/`slice` plus bare words `version`/`help`; options `spec=`, `hints=`, `out=`; no POSIX flag is accepted (`--version` is an unknown command, exit 2) |
| EXIT-CODE-OK / ERROR / INVOCATION | required | 0 / 1 / 2 exactly as the spec's `ExitCode` type defines them |
| STREAM-DIAGNOSTICS / STREAM-OUTPUT | required | findings and errors on stderr, summaries and listings on stdout |
| SIGNAL-HANDLING | SIGTERM, SIGINT | handler in `main()`; removes any file the current `slice` run wrote, then exits cleanly — no partial output |
| OUTPUT-FORMAT | RPM, DEB required | both produced; OCI/PKG not active in the resolved preset |
| INSTALL-METHOD | OBS required, curl forbidden | README documents zypper/dnf/apt from OBS; no curl anywhere |
| PLATFORM | Linux required | Linux only, as the spec's DEPLOYMENT section states |
| CONFIG-ENV-VARS | forbidden | the tool reads no environment variable at all; no test-injection variable was needed because the tests use per-test working directories and relative paths |
| NETWORK-CALLS | forbidden | no `net` import anywhere; no network at build time either (no dependencies) |
| FILE-MODIFICATION (input files) | forbidden | inputs are opened read-only; `TestInvariantSourcesAreReadOnly` hashes spec and hints before and after every verb |
| IDEMPOTENT | required | byte-identical outputs across runs; `TestSliceIsIdempotent`, `TestExampleTwoRunsAreByteIdentical` |
| PRESET-SYSTEM | systemd-style | no preset file was present; resolution order honoured, template defaults used |

---

## How STEPS ordering was applied

**BEHAVIOR: list** — steps 1–5 in order, in `List()` (`internal/pcd-slice/verbs.go`):
`LoadSources` (step 1; extension, existence and readability predicates first,
`error: cannot read {path}` and exit 2 on failure) → `Analyze` computes the
transitive closure per behavior (step 2, `closure.go`) and attributes
invariants, examples and hints blocks (step 3, `parse.go`/`hints.go`) → one
tab-separated line per behavior in specification order plus
`total: {n} behaviors` (step 4) → return 0 unconditionally (step 5): `list`
never fails on findings, which `TestListNeverFailsOnFindings` pins.

**BEHAVIOR: check** — steps 1–10 in order, in `Check()` and
`Analysis.Findings()`: parse as in list (1); `undefined-type` on the
definition side only (2); `unknown-binding` (3); `unassignable-example` (4);
`unassignable-hints`, hints preamble exempt (5); `duplicate-behavior` (6);
`no-behaviors` (7); the completeness scan over the assignment array, which
would report rule `unassigned` (8); `stale-output`, only when `out=` is given
*and* the directory exists (9); findings sorted by file then line to stderr,
summary to stdout, exit 1 if any finding else 0 (10).

**BEHAVIOR: slice** — steps 1–7 in order, in `Slice()` (`write.go`):
run check without `out=` and refuse on any finding, printing them and writing
nothing (1); compute every output in memory, bundle file names with `/`
replaced by `-` and a collision reported as `duplicate-behavior` (2); one
provenance line per markdown file with the exact field order the spec gives
(3); determinism by construction — sections in source order, closure
alphabetical, manifest sorted, no map iteration in any emitting path (4);
create `out` if absent, refuse with exit 2 on a foreign file *before* the
first write, remove only files carrying a pcd-slice provenance line (5);
write every file then the manifest, removing this run's files on any write
error (6); one summary line on stdout, exit 0 (7).

## MECHANISM notes and hints applied

The Go hints file was followed closely, and its guidance is visible in the
code:

- *One pass, one classifier* — `os.ReadFile`, split on `"\n"` without
  normalising, and a small state machine (`fenceMap` + section walk). No
  markdown library: the spec's SECTION GRAMMAR is the whole grammar.
- *Assignment before output* — `Spec.Assign` is one entry per source line,
  initialised to the `unassigned` sentinel. Outputs are pure projections of
  that array, which makes the completeness invariant a slice scan instead of
  an argument.
- *Type closure* — definitions collected inside the TYPES fence, one
  `\b(A|B|C)\b` alternation compiled from the names sorted longest-first with
  `regexp.QuoteMeta`, worklist over the behavior text and then over each
  newly added definition; alphabetised only at output time.
- *Attribution* — the WHEN line's invoked identifier
  (`^\s*(?:result\s*=\s*)?([a-z][a-z0-9_-]*)\(`), then the example's own
  name; hints headings tested after stripping `**` and back quotes.
- *Determinism mechanics* — never range over a map when emitting;
  `sort.Strings` for closures, `sort.Slice` for the file list;
  `crypto/sha256` over the raw bytes as read; `hex.EncodeToString`; the tool
  version is `main.version`, set with `-ldflags "-X main.version=..."`,
  defaulting to `dev`.
- *Line endings* — a trailing `\r` stays attached to its line for assignment
  and for verbatim projection; only the lines pcd-slice adds itself always
  end in `\n`.
- *Atomic out-directory handling* — everything computed first, foreign-file
  refusal before the first write, files written this run tracked and removed
  on error, the manifest written last as the commit marker.
- *Errors and exits* — a single `os.Exit` at the top of `main`; the verbs
  return an int. Exit-2 messages start with `error: ` and name the path.
- *What not to build* — no configuration file, no environment variables, no
  colour, no concurrency, no YAML front matter.

## Type bindings chosen (no TYPE-BINDINGS table in the template)

| Spec type | Go binding |
|---|---|
| `SpecFile`, `HintsFile` | `string` path + `*Source` after loading; the refinement predicate (exists, readable, `.md`) is checked in `LoadSource` |
| `OutDir` | `string` in `Options.Out`; created with `os.MkdirAll` |
| `BehaviorName` | `Behavior.Name string` |
| `TypeName` | `TypeDef.Name string`, indexed by `map[string]*TypeDef` |
| `InvariantBinding` | `Invariant.Bindings []string` + resolved `Targets []int` |
| `Assignment` | `AssignKind` enum (`unassigned`, `preamble`, `behavior`, `excluded`) + behavior index |
| `Finding` | `struct{ Rule, File string; Line int; Message string }` |
| `ExitCode` | `ExitCode int` with `ExitClean`/`ExitFindings`/`ExitInvocation` |

## BEHAVIOR constraints

All three BEHAVIORs (`list`, `check`, `slice`) carry `Constraint: required`
and are implemented unconditionally. The spec declares no `supported` and no
`forbidden` BEHAVIOR, so no behaviour was gated on a preset and none was
suppressed. No BEHAVIOR is missing from the implementation, and none had to
be flagged "not yet scheduled".

## Parsing approach

A single linear pass over the source lines with three pieces of state: fence
depth, current `##`-level section, and current behavior. Structure is
recognised by prefix and suffix tests plus a handful of anchored regular
expressions; fenced content is never interpreted as structure. Sections are
materialised as half-open line ranges, then sub-parsed: the TYPES fences into
definitions (a definition owning its comment and continuation lines), the
INVARIANTS section into list items (a binding suffix is recognised only at
end of line), the top-level EXAMPLES section into `### EXAMPLE` blocks.
Attribution then stamps the per-line assignment array, and every output file
is a projection of that array plus the attributed material appended under a
heading that names its source. Hints files get the same treatment with the
heading-attribution rule of the grammar.

The consequence of building the assignment array first is that the
completeness invariant is checked, not asserted: any line still carrying the
sentinel after the walk becomes a finding.

## Signal handling approach

`main()` installs one handler for `SIGTERM` and `SIGINT` before dispatch. On
either signal it calls `pcdslice.RemovePartialOutput()`, which removes every
file the current run has already written (tracked under a mutex as each write
succeeds), and then exits cleanly. `list` and `check` never write, so for
them the handler is a plain clean exit; for `slice` it guarantees the
spec's "no partial output" property on interruption, the same cleanup path
the write-error branch uses. There is no `os.Exit` below `main` in the
implementation, so no deferred cleanup is ever skipped.

---

## Phase 6 — Compile gate

| Step | Command | Result |
|---|---|---|
| 1 Dependency resolution | `go mod tidy` | **pass** — no dependencies; `go.mod` declares module and Go version only, no `go.sum` is generated because the standard library needs none |
| 2 Compilation | `go build ./...` | **pass** |
| 2 Vet | `go vet ./...` | **pass** (no diagnostics) |
| 2b Formatting | `gofmt -l cmd internal independent_tests` | **pass** (empty output) |
| 3 Translator test run | `make test` → `go test ./independent_tests/claude-opus-5/...` | **pass** — 50 tests/subtests, 0 failures, 0 skips |
| 4 Test-author test run | — | not applicable (single-LLM run) |
| 5 Record result | this section | done |

Additional gate checks from the Go milestones hints:

```
file ./pcd-slice            → ELF 64-bit, statically linked
./pcd-slice version         → pcd-slice 0.1.0 / spec:bb31cb27…  (exit 0)
./pcd-slice help            → usage                              (exit 0)
./pcd-slice format=bad_value→ error: unknown command             (exit 2)
make dist                   → pcd-slice-0.1.0.tar.gz + pcd-slice-0.1.0-vendor.tar.gz
make clean                  → source tree contains only the intended files
```

`make dist` was executed and its outputs inspected (single top-level
`pcd-slice-0.1.0/` directory; the vendor tarball carries `vendor/modules.txt`
so that an OBS build with `GOFLAGS=-mod=vendor` succeeds for a
dependency-free module), then removed again by `make clean`. The run leaves
no scratch files.

### Dogfooding (not a substitute for the tests, but worth recording)

The binary was run against its own specification and Go hints file:

```
$ pcd-slice list spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md
list	3	1	1	0
check	5	2	2	0
slice	4	5	7	0
total: 3 behaviors

$ pcd-slice check spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md
unassignable-example: pcd-slice.spec.md:401: two_runs_are_byte_identical: attributable to no behavior
pcd-slice check: 1 findings
```

That single finding is correct behaviour, not a defect: EXAMPLE
`two_runs_are_byte_identical` has a prose WHEN line ("slice runs twice into
two fresh directories") with no invoked identifier, and its name contains no
behavior name as a whole word, so neither attribution rule of the grammar
applies. It is recorded under *Specification ambiguities* below. With that
one WHEN line changed locally to `result = slice(spec, out)`, `check` is
clean, `slice` emits `preamble.md`, `list.md`, `check.md`, `slice.md` and
`MANIFEST.tsv`, two runs into two directories are byte-identical, changing one
word in the spec produces one `stale-output` line per emitted file, and a
line-by-line count over the emitted files confirms the completeness invariant
(the only lines appearing more than once are the three invariants whose
`(binds: …)` names more than one behavior — see ambiguity 1).

---

## Test results — translator suite

`independent_tests/claude-opus-5/` — 50 tests and subtests, all **pass**, no
skips, no live external service, no network.

| Test | Covers | Result |
|---|---|---|
| `TestExampleClosureIsTransitive` | EXAMPLE closure_is_transitive | pass |
| `TestExampleUntaggedInvariantIsGlobal` | EXAMPLE untagged_invariant_is_global | pass |
| `TestExampleBoundInvariantLandsInNamedBundles` | EXAMPLE bound_invariant_lands_in_named_bundles | pass |
| `TestExampleBindingToUnknownBehaviorIsAFinding` | EXAMPLE binding_to_unknown_behavior_is_a_finding | pass |
| `TestExampleNestedExampleTravelsWithItsBehavior` | EXAMPLE nested_example_travels_with_its_behavior | pass |
| `TestExampleTopLevelExampleAttributedByWhenLine` | EXAMPLE top_level_example_attributed_by_when_line | pass |
| `TestExampleChangelogIsExcluded` | EXAMPLE changelog_is_excluded | pass |
| `TestExampleSliceRefusesOnFindingsAndWritesNothing` | EXAMPLE slice_refuses_on_findings_and_writes_nothing | pass |
| `TestExampleTwoRunsAreByteIdentical` | EXAMPLE two_runs_are_byte_identical | pass |
| `TestExampleStaleOutputIsOneFindingPerFile` | EXAMPLE stale_output_is_one_finding_per_file | pass |
| `TestExampleForeignFileInOutRefuses` | EXAMPLE foreign_file_in_out_refuses | pass |
| `TestInvariantSourcesAreReadOnly` | INVARIANT sources are read-only (list, check, slice) | pass |
| `TestListAndCheckNeverWriteFiles` | POSTCONDITION list/check never write files | pass |
| `TestInvariantProvenanceHeader` | INVARIANT provenance; hints order in the header | pass |
| `TestProvenanceWithoutHints` | provenance boundary case: zero hints files | pass |
| `TestManifestDescribesEveryEmittedFile` | POSTCONDITIONs on out and the manifest; slice summary line | pass |
| `TestInvariantCompletenessCoversEveryLineExactlyOnce` | INVARIANT completeness | pass |
| `TestInvariantBehaviorAppearsInExactlyOneBundle` | INVARIANT one behavior, one bundle | pass |
| `TestListCountsAndOrder` | BEHAVIOR list, steps 2–4, source order | pass |
| `TestListNeverFailsOnFindings` | BEHAVIOR list, step 5 | pass |
| `TestCheckCleanSpec` | BEHAVIOR check, clean summary and exit 0 | pass |
| `TestCheckDuplicateBehavior` | check rule duplicate-behavior | pass |
| `TestCheckNoBehaviors` | check rule no-behaviors | pass |
| `TestCheckUndefinedType` | check rule undefined-type | pass |
| `TestCheckUnassignableExample` | check rule unassignable-example; slice refusal | pass |
| `TestCheckUnassignableHints` | check rule unassignable-hints | pass |
| `TestHintsPreambleIsNotAFinding` | check step 5 exemption | pass |
| `TestFindingsAreSortedByFileThenLine` | check step 10 ordering | pass |
| `TestHintsBlocksAreAttributed` | hints grammar: attributed block vs hints preamble | pass |
| `TestSliceIsIdempotent` | IDEMPOTENT; re-slice into the same directory | pass |
| `TestSliceRemovesItsOwnStaleFiles` | slice step 5 removal rule | pass |
| `TestSliceCreatesOutDirectory` | slice step 5 creation; summary path | pass |
| `TestBundleIsSelfLocating` | INVARIANT bundle + preamble suffice; bundle layout | pass |
| `TestInvocationErrors` (+10 subtests) | ExitCode 2 paths: no args, unknown verb, missing spec, absent spec, non-`.md`, absent hints, slice without out, unknown option, positional argument, empty value | pass |
| `TestUnreadableSpecExitsTwo` | TYPES predicate "readable" | pass |
| `TestUnwritableOutDirExitsTwo` | PRECONDITION "out is creatable or writable" | pass |
| `TestVersionEmbedsSpecHash` | spec hash in binary version output | pass |
| `TestHelpExitsZero` | bare-word help | pass |
| `TestRefinementPredicatesAreNotTypeReferences` | undefined-type must not fire on `AND`/`OR` in a predicate | pass |
| `TestUndefinedTypeIsReportedOncePerLine` | one finding per undefined name per line | pass |

Test discipline: every test drives the real binary through `exec.Command`;
no test imports the implementation package. Each test runs in its own
`t.TempDir()` with `HOME` and `TMPDIR` redirected into that sandbox, a
minimal environment, and the working directory set to the sandbox, so path
arguments read exactly as in the spec's EXAMPLEs. The fixture writer refuses
to write outside the sandbox. No test asserts only on the exit code. The
suite is re-runnable with no residue and passes from any working directory
and as any unprivileged user; the two permission-based tests skip themselves
if run as root (they were not skipped in this run).

## Test results — test-author suite

None present. Single-LLM run.

## Test Refinements

| Test | Result before | Action | Rationale |
|---|---|---|---|
| `TestSliceRemovesItsOwnStaleFiles` | failed | test edited | The fixture removed behavior `purge` but kept the invariant `(binds: purge)`. The binary then reported `unknown-binding` and `slice` refused — correct per BEHAVIOR: check step 3 and BEHAVIOR: slice step 1. The tool was right, the fixture was structurally incomplete for the test's intent (which is the removal of a stale bundle file); the fixture edit now drops the bound invariant along with the behavior. No assertion changed. |
| `TestRefinementPredicatesAreNotTypeReferences` | *(added after the first run)* | code fixed | Dogfooding against `pcd-slice.spec.md` showed `undefined-type` firing on `AND` inside `SpecFile := path where file_exists AND readable AND extension = ".md"`. The type-shape pattern accepted all-capital words. It now requires at least one lower-case character per camel hump, and this test pins the behaviour. The rule is explicitly definition-side only (check step 2), so a narrow shape is the conservative reading. |
| `TestUndefinedTypeIsReportedOncePerLine` | *(added after the first run)* | code fixed | The same dogfooding run reported the identical finding twice for two occurrences of one undefined name on one line. `Findings()` now drops findings identical in rule, file, line and message; the test pins one finding per undefined name per line. |
| all other 47 tests/subtests | passed | none | — |

No other test was touched after a run. No test was weakened, and no
assertion was changed to match observed behaviour.

---

## Per-EXAMPLE confidence

| EXAMPLE | Confidence | Verification method | Unverified claims |
|---|---|---|---|
| closure_is_transitive | Medium | `TestExampleClosureIsTransitive`: asserts the list line `record 3 0 0 0`, the total line, `Requires-Types: Item, Path, Plan` verbatim, `Unused` present in `preamble.md` and absent from `record.md`; exit 0 | none |
| untagged_invariant_is_global | Medium | `TestExampleUntaggedInvariantIsGlobal`: invariant text present in `preamble.md`, absent from both bundles; exit 0 | none |
| bound_invariant_lands_in_named_bundles | Medium | `TestExampleBoundInvariantLandsInNamedBundles`: present in `reset.md`, absent from `purge.md` and `preamble.md` | none |
| binding_to_unknown_behavior_is_a_finding | Medium | `TestExampleBindingToUnknownBehaviorIsAFinding`: stderr contains `unknown-binding: spec.md:{line}: renmae` with the line number computed from the fixture, stdout equals `pcd-slice check: 1 findings`, exit 1 | none |
| nested_example_travels_with_its_behavior | Medium | `TestExampleNestedExampleTravelsWithItsBehavior`: the nested example is inside `verify.md` and positioned before any appended source-naming heading (i.e. inside the verbatim block), absent from the other bundle and the preamble, and `check` is clean — no attribution rule was consulted | none |
| top_level_example_attributed_by_when_line | Medium | `TestExampleTopLevelExampleAttributedByWhenLine`: the example is in `lint.md` under a heading naming `spec.md`, and in no other file | none |
| changelog_is_excluded | Medium | `TestExampleChangelogIsExcluded`: no changelog row and no `## Changelog` heading in any emitted file; `check` reports no `unassigned` finding, so the lines counted as excluded | none |
| slice_refuses_on_findings_and_writes_nothing | Medium | `TestExampleSliceRefusesOnFindingsAndWritesNothing`: stderr contains `duplicate-behavior`, `spec.d` does not exist afterwards, exit 1 | none |
| two_runs_are_byte_identical | Medium | `TestExampleTwoRunsAreByteIdentical`: SHA-256 of every same-named file compared across two fresh directories, including `MANIFEST.tsv`; no user name or timestamp marker in any output | the *cross-build* case is not tested: two binaries built with different `-ldflags` versions differ in the provenance `version=` field, as the hints file states |
| stale_output_is_one_finding_per_file | Medium | `TestExampleStaleOutputIsOneFindingPerFile`: clean directory reports `clean`; after one word changes, exactly one `stale-output` line per emitted file (bundles, preamble and manifest), exit 1 | none |
| foreign_file_in_out_refuses | Medium | `TestExampleForeignFileInOutRefuses`: stderr contains `error: foreign file in out: spec.d/notes.txt`, the file's bytes are unchanged, nothing else was written, exit 2 | none |

Confidence is **Medium** throughout for one structural reason only: this is a
single-LLM run, and the prompt's confidence scale demotes an EXAMPLE to
Medium when a test-author suite is absent. Every EXAMPLE above is covered by
a named test function that passes without any live external service, and
Tests-First-Compliance is `yes`; a dual-LLM re-run with an independent
test-author suite would lift these rows to High without any code change.

---

## Specification ambiguities

1. **A multi-bound invariant versus "exactly once".** `InvariantBinding` says
   a bound invariant "belongs to exactly the named behaviors' bundles"
   (plural), while the completeness invariant says every line is "assigned to
   exactly one output". The spec's own INVARIANTS section contains three
   invariants binding two or three behaviors. Conservative resolution: the
   `Assignment` array — the accounting device the completeness rule scans —
   stamps such a line onto the *first* named behavior, so every line still
   has exactly one assignment and no line is lost; the bundle projection then
   appends the invariant to *each* named bundle, because
   `InvariantBinding` is explicit that it belongs to all of them and because
   a bundle must be self-locating. The observable consequence is that a
   multi-bound invariant's text appears in more than one bundle. Recorded
   rather than silently resolved.
2. **`two_runs_are_byte_identical` is unattributable under its own grammar.**
   Its WHEN line is prose, and its name contains no behavior name as a whole
   word, so `check` on `pcd-slice.spec.md` reports one
   `unassignable-example`. Implemented exactly as the grammar states (WHEN
   identifier first, example name second, otherwise a finding) rather than
   inventing a third rule such as "a behavior name anywhere in the WHEN
   text". The spec author may wish to reword that WHEN line.
3. **`undefined-type` detection needs a shape heuristic.** Check step 2 says
   the rule fires when "a type definition references a TypeName that is
   defined nowhere", but an undefined name is by definition not in the name
   table, so some shape test is unavoidable. Conservative choice: an
   identifier with two or more camel humps, each with at least one lower-case
   character (`MissingItem`, `OutDir`), outside `//` comments, excluding the
   definition's own name. This deliberately under-reports (an all-capital
   name such as `URLPath` is not flagged) rather than firing on prose or on
   the all-capital connectives of a refinement predicate.
4. **The `unknown-binding` message.** EXAMPLE 4 requires stderr to contain
   `unknown-binding: spec.md:{line}: renmae`, which fixes the message to
   begin with the offending name. The message is exactly the name, nothing
   appended, so the emitted line matches the EXAMPLE literally.
5. **Which examples `list` counts.** The list line counts *attributed
   top-level* examples. Examples nested inside a behavior block are part of
   the block ("no attribution rule is consulted") and are not counted
   separately; they still travel in the bundle.
6. **`slice` output on refusal.** Step 1 says "print its findings, write
   nothing, exit 1" and does not mention the check summary line. `slice`
   therefore prints only the finding lines, on stderr, and nothing on stdout.
7. **`check out=` when the directory does not exist.** Step 9 is scoped to
   "out= is given and exists", so a missing directory produces no
   `stale-output` finding at all (rather than one "everything is missing"
   finding). A directory that exists but lacks `MANIFEST.tsv` is stale in
   full, per the hints file.
8. **`MANIFEST.tsv` has no provenance line** — it is not markdown, and the
   provenance invariant is stated for markdown files. Staleness of the
   manifest is therefore decided by comparing its full content with the
   recomputed manifest, and the manifest doubles as the commit marker of a
   run.
9. **Where a behavior block ends.** The grammar says "to the next `## `
   heading". A `# ` (depth 1) heading also ends the block, since a
   depth-1 heading cannot plausibly be inside a behavior; `### ` and deeper
   never end it.
10. **An unattributable example still needs an assignment** so that the
    completeness scan has no sentinel left. Such lines are assigned to the
    preamble; the `unassignable-example` finding is what makes the situation
    visible, and `slice` refuses before any file is written anyway.
11. **Hints depth-1 headings do not attribute.** The grammar says "a `## ` or
    deeper heading", so a `# Title` line never opens an attributed block,
    though it does close one.
12. **A subdirectory inside `out`** is treated as a foreign entry (exit 2,
    naming it). The spec only anticipates files; refusing is the conservative
    choice because pcd-slice would otherwise have to decide whether to remove
    a directory it did not create.
13. **Version output format.** The spec says nothing; the template requires
    `spec:<hash>` in the version output. `pcd-slice version` prints
    `pcd-slice <version>` and `spec:<hash>` on two lines. `--version` is
    *not* accepted, because POSIX flag style is forbidden by the template.

## Rules that could not be implemented exactly as written

None. Every STEPS entry, every finding rule and every postcondition is
implemented as written, with the interpretations above where the text left a
choice. Two deliberate deviations in *packaging and layout*, both documented:

- The implementation package directory is `internal/pcd-slice/` (the
  template's `internal/<n>/`), while the Go package identifier is
  `pcdslice`: a Go package name cannot contain a hyphen. The import is
  aliased in the entry point.
- `OCI` and `PKG` deliverables were not produced: both are `supported`, not
  `required`, and no preset activates them in this run.

---

## Public API Surface

Recorded per `PUBLIC-API-SURFACE: recorded-in-report`. The next translation
of this spec at Version 0.1.0 must preserve every entry below; it may add.

### Module `github.com/mge1512/pcd/tools/pcd-slice` — package `main` (`cmd/pcd-slice/main.go`)

```
func main()
var version string                                   // set via -ldflags -X main.version
```

No other symbol is exported from the entry point; `dispatch` and
`parseOptions` are unexported by design (the entry point is CLI dispatch
only).

### Package `internal/pcd-slice` (`pcdslice`)

Constants and variables:

```
const SpecSHA256 = "bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04"
const ToolName   = "pcd-slice"
var   Version    = "dev"
const ExitClean, ExitFindings, ExitInvocation ExitCode          // 0, 1, 2
const AssignUnassigned, AssignPreamble, AssignBehavior, AssignExcluded AssignKind
```

Types:

```
type ExitCode int
type AssignKind uint8
    func (k AssignKind) String() string
type Assignment struct { Kind AssignKind; Behavior int }
type Finding struct { Rule, File string; Line int; Message string }
    func (f Finding) String() string
type Options struct { Spec string; Hints []string; Out string; Stdout, Stderr io.Writer }
type Source struct { Path string; Data []byte; SHA256 string; Lines []string }
type Sources struct { Spec *Source; Hints []*Source }
type TypeDef struct { Name string; Line, Start, End int; Text string }
type Invariant struct { Start, End, BindLine int; Bindings []string; Targets []int }
    func (i *Invariant) Global() bool
type Example struct { Name string; Line, Start, End, Behavior int }
type Behavior struct {
        Name, FileName string
        HeadingLine, Start, End int
        Closure []string
        Invariants []*Invariant
        Examples []*Example
        HintsBlocks []*HintsBlock
    }
type Spec struct {
        Src *Source
        Assign []Assignment
        Behaviors []*Behavior
        TypeDefs []*TypeDef
        Invariants []*Invariant
        Examples []*Example
    }
type HintsBlock struct { Heading string; Line, Start, End int; Behavior string; Doc *HintsDoc }
type BadHeading struct { Line int; Text string }
type HintsDoc struct { Src *Source; Blocks []*HintsBlock; PreambleLines []int; BadHeadings []BadHeading }
    func (d *HintsDoc) BlockText(b *HintsBlock) []string
    func (d *HintsDoc) PreambleText() []string
type OutFile struct { Name string; Data []byte }
type Outputs struct { Files []OutFile; Manifest OutFile; Provenance string; Bundles int }
    func (o *Outputs) All() []OutFile
type Analysis struct { Sources *Sources; Spec *Spec; Hints []*HintsDoc }
    func (a *Analysis) BuildOutputs() *Outputs
    func (a *Analysis) Findings() []Finding
    func (a *Analysis) ProvenanceLine() string
    func (a *Analysis) StaleFindings(out string) []Finding
```

Functions:

```
func Analyze(s *Sources) *Analysis
func Check(o Options) int
func List(o Options) int
func Slice(o Options) int
func LoadSource(path string) (*Source, error)
func LoadSources(o Options) (*Sources, error)
func ParseHints(src *Source, behaviors []string) *HintsDoc
func ParseSpec(src *Source) *Spec
func RemovePartialOutput()
func SortFindings(fs []Finding)
func SplitLines(data []byte) []string
```

---

## File inventory

```
VERSION
go.mod
Makefile
LICENSE
README.md
pcd-slice.1.md
pcd-slice.1
pcd-slice.spec
cmd/pcd-slice/main.go
internal/pcd-slice/types.go
internal/pcd-slice/source.go
internal/pcd-slice/parse.go
internal/pcd-slice/closure.go
internal/pcd-slice/hints.go
internal/pcd-slice/check.go
internal/pcd-slice/render.go
internal/pcd-slice/verbs.go
internal/pcd-slice/write.go
debian/control
debian/changelog
debian/rules
debian/copyright
independent_tests/claude-opus-5/helpers_test.go
independent_tests/claude-opus-5/fixtures_test.go
independent_tests/claude-opus-5/pcd_slice_test.go
translation_report/translation-workflow.pikchr
TRANSLATION_REPORT.md
pcd-slice                     (compile-gate build output at BINARY-LOCATION: project-root)
```

Every source, test and packaging file embeds
`sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04`
in the comment syntax of its format. No artefact contains a placeholder
value.
