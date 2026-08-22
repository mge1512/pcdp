# TRANSLATION_REPORT.md — pcd-slice

**Spec-SHA256:** `c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f`
**Spec-SHA256 (host):** `c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f`

**Included-Specs:** none — the host spec's META declares no `Includes:`
directives, so the merged spec text is byte-identical to the host spec file
and the merged hash equals the host hash.

| Path | SHA256 |
|------|--------|
| *(none)* | — |

**LLM-Name:** `claude-opus-5`
**Mode:** `translator`
**Spec-Schema:** 0.4.0 (host META) — include resolution implemented; no includes present
**Component version:** 0.1.1 (spec META `Version:`, mirrored in `VERSION`)
**Date:** 2026-08-22

---

## Translation Inputs (provenance)

Every file consumed as a translation input, hashed exactly as read from the
input directory at translation time.

- `Spec-SHA256 (merged):` `c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f`
- `Spec-SHA256 (host):` `c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f`
- `Decisions-Hints-SHA256:` `pcd-slice.go.hints.md` `c1b2685bbc5b395700a2024d0015395325ef0d2e348e6c43671d352a7c818d72`
- `Milestones-Hints-SHA256:` `cli-tool.go.milestones.hints.md` `fe73f2ab89be21a0fd0eb37a818e586b68a338a8392a8f14725c611a2d33b249`
- `Template-SHA256:` `cli-tool.template.md` `13cbe57a345e6be7f6b91f0ee3a6d5d5bfc82a957143e1668325a76bfcac13fa`
- `Prompt-SHA256:` `prompt.md` `2317d61f6e5631af52943243aee40cc5dc0ed29c55984d997dcf72397930d912`
- `Role-SHA256:` `ROLE.md` `c1a484b969f2516a1f34ae2baa73e34b8c8e38afa6cbfaba782711406c2b7c66`
- `Style-Hints-SHA256:` `none` (no `<scope>.<language>.style.hints.md` in the input directory, `/etc/pcd/hints/` or `.pcd/hints/`)
- `Library-Hints-SHA256:` `none`
- `Upgrade-Brief-SHA256:` `none`
- `Directive-SHA256:` `none`

---

## Tests-First-Compliance

**`yes`.** The complete test suite
(`independent_tests/claude-opus-5/pcd_slice_test.go`, 41 190 bytes, 39 test
functions plus 11 sub-tests) was written and committed to disk **before any
implementation source file existed**. The structural guard at step 3 of the
translator flow was satisfied: after writing the suite, the directory
`independent_tests/claude-opus-5/` was verified to exist and contain a test
file; only then was `internal/pcdslice/types.go` — the first implementation
file — written.

Order of writes (observable from the run):

1. `independent_tests/claude-opus-5/pcd_slice_test.go`
2. `internal/pcdslice/{types,parse,closure,attribute,hints,emit,check,list,slice,run}.go`
3. `cmd/pcd-slice/main.go`, `go.mod`
4. packaging, documentation, report

The suite passed on its first execution against the implementation; no test
was edited after any test run (see **Test Refinements**).

## Continuity-Check

**Not applicable — no test-author input.** The input directory
(`/tmp/pcd-input/`) contained only `ROLE.md`, `prompt.md`,
`pcd-slice.spec.md`, `cli-tool.template.md`, `pcd-slice.go.hints.md` and
`cli-tool.go.milestones.hints.md`. There is no
`independent_tests/<other-role-llm-name>/` directory and no `TEST_REPORT.md`,
so this is a **single-LLM run**, which the prompt declares a fully supported
invocation. Steps 6 and 7 of the translator flow's dual-LLM path did not
apply and no test-author suite was executed.

---

## Resolution and constraints

### Target language

**Go**, the `LANGUAGE` default of `cli-tool.template.md` (TEMPLATE-TABLE row
`LANGUAGE | Go | default`). No preset files were present in the input
directory, in `/etc/pcd/presets/`, `~/.config/pcd/presets/` or a project
`.pcd/` directory, so no override applied. Both hints files supplied are
Go-specific, confirming the default. No deviation from the template default.

### Module identity resolved

**`github.com/mge1512/pcd/tools/pcd-slice`**, from **authoritative source 1**:
the spec's META `Module:` field. Sources 2–4 were not consulted for a
competing value (the Go hints file declares no module name; no prior
manifest existed in the output directory; the spec-title fallback was not
needed), so no conflict arose and `MODULE-IDENTITY: conflict-halts` did not
fire. The identity is propagated to `go.mod`, the internal import path in
`cmd/pcd-slice/main.go`, and — as the project home — to the RPM `URL:`,
`debian/control` `Homepage:` and the README.

### Delivery mode

**Mode 1 — filesystem.** All artefacts were written directly to
`/tmp/pcd-output/` with the filesystem write tool. Nothing was emitted to the
terminal as a deliverable. No VCS repository was present, so nothing was
committed.

### Active MILESTONE

**None.** The specification contains no `## MILESTONE:` section, so the full
spec was translated as normal (no scaffold pass, no deferred BEHAVIORs). All
three BEHAVIORs are implemented for real; there are no stubs in the tree.

### BEHAVIOR constraints

| BEHAVIOR | Constraint | Treatment |
|----------|-----------|-----------|
| `list`   | required  | Implemented unconditionally (`internal/pcdslice/list.go`) |
| `check`  | required  | Implemented unconditionally (`internal/pcdslice/check.go`) |
| `slice`  | required  | Implemented unconditionally (`internal/pcdslice/slice.go` + `emit.go`) |

No BEHAVIOR carries `Constraint: supported` or `forbidden`; no BEHAVIOR is
"not yet scheduled".

### INTERFACES / TYPE-BINDINGS / GENERATED-FILE-BINDINGS / spec DELIVERABLES

- The spec has **no `## INTERFACES` section**, so no production/test-double
  pairs were generated. The test suite is black-box over the CLI, which is
  the interface the DEPLOYMENT section declares.
- The template contains **no `## TYPE-BINDINGS` section**; the spec's logical
  types were bound by hand to idiomatic Go equivalents (see *Type mapping*).
- The template contains **no `## GENERATED-FILE-BINDINGS` section**.
- The spec has **no `## DELIVERABLES` section with COMPONENT: entries**; the
  produced file set derives from the template's DELIVERABLES table alone.

### Type mapping (spec TYPES → Go)

| Spec type | Go binding | Notes |
|-----------|------------|-------|
| `SpecFile`, `HintsFile` | `string` path + `readMarkdown()` precondition check | existence, readability and the `.md` extension are validated on load; violation is exit 2 |
| `OutDir` | `string` path | created with `os.MkdirAll`; nothing is written outside it |
| `BehaviorName` | `string` (`Behavior.Name`) | uniqueness enforced by the `duplicate-behavior` rule |
| `TypeName` | `string` (`TypeDef.Name`) | definition body captured as `TypeDef.Body` / `.Text` |
| `InvariantBinding` | `Invariant.Binds []string` | `Owner` is the first named behavior; replication happens in the emitter only |
| `Assignment` | `Kind` (`unassigned`/`preamble`/`behavior`/`excluded`) + `Behavior` + emission `Role` | one entry per source line in `Spec.Assign` |
| `Finding` | `Finding{Rule, File string; Line int; Message string}` | `String()` renders `{rule}: {file}:{line}: {message}` |
| `ExitCode` | `ExitOK`/`ExitFindings`/`ExitInvocation` = 0/1/2 | returned by `Run`, applied by the single `os.Exit` in `main` |

### Template constraints compliance

| Key | Required value | How satisfied |
|-----|----------------|---------------|
| VERSION | MAJOR.MINOR.PATCH | `0.1.1` from spec META, single source in `VERSION`, injected via `-ldflags -X main.version` |
| LANGUAGE | Go (default) | Go 1.22 module, standard library only |
| BINARY-TYPE | static | `CGO_ENABLED=0` in `Makefile`, RPM `%build` and `debian/rules`; verified: `file ./pcd-slice` → "statically linked" |
| PERSISTENCE | none | no state is kept between runs; no config file |
| SOURCE-PARTITIONING | modular, one-entry-one-implementation | entry point `cmd/pcd-slice/main.go` (dispatch + signals only, 33 lines); 10 implementation files in `internal/pcdslice/` |
| SOURCE-PARTITIONING | by-behaviour-domain (supported) | applied: `parse.go`, `closure.go`, `attribute.go`, `hints.go`, `emit.go` (domains) and `list.go`, `check.go`, `slice.go` (one per BEHAVIOR) |
| MODULE-IDENTITY | host-specified, propagated | spec META `Module:`; propagated to `go.mod`, import path, packaging |
| PUBLIC-API-SURFACE | recorded-in-report | see `## Public API Surface` below |
| BINARY-COUNT | 1 | one binary, `pcd-slice` |
| BINARY-LOCATION | project-root | `make build` writes `./pcd-slice` next to `go.mod`; the test suite invokes `../../pcd-slice` and builds it there in `TestMain` from `../../cmd/pcd-slice` |
| RUNTIME-DEPS | none | static binary, standard library only; RPM/DEB declare no runtime dependency beyond `${shlibs:Depends}`/`${misc:Depends}` |
| CLI-ARG-STYLE | key=value | `spec=`, `hints=`, `out=`; any argument starting with `-` is rejected with exit 2 |
| CLI-ARG-STYLE | bare-words / subcommand (supported) | verbs `version`, `help`, `list`, `check`, `slice` as bare first words, per the spec's DEPLOYMENT invocation grammar |
| EXIT-CODE-OK / ERROR / INVOCATION | 0 / 1 / 2 | `ExitOK`, `ExitFindings`, `ExitInvocation`; matches the spec's `ExitCode` type |
| STREAM-DIAGNOSTICS | stderr | findings and `error:` lines only |
| STREAM-OUTPUT | stdout | listing, check summary, slice summary, version, help |
| SIGNAL-HANDLING | SIGTERM, SIGINT | handler in `main()`; calls `pcdslice.Abort()` (removes files written by the run) then exits — no partial output |
| OUTPUT-FORMAT | RPM, DEB (required) | `pcd-slice.spec`; `debian/{control,changelog,rules,copyright}` |
| OUTPUT-FORMAT | OCI, PKG, binary (supported) | **not produced**: no preset activates OCI or PKG, and `PLATFORM: macOS` is not declared (spec DEPLOYMENT says "Linux only in v1"). `binary` requires no descriptor. |
| INSTALL-METHOD | OBS (required), curl (forbidden) | README documents zypper/apt/dnf only; no curl anywhere |
| PLATFORM | Linux | spec DEPLOYMENT declares Linux only in v1 |
| CONFIG-ENV-VARS | forbidden | no environment variable is read anywhere in the implementation (`os.Getenv` appears nowhere) |
| TEST-INJECTION | single-declared-env-var (supported) | **not used** — the tool needs no test-injection variable; tests are hermetic through `key=value` paths and per-test temporary directories |
| NETWORK-CALLS | forbidden | no `net`, `net/http` or exec of network tools; imports are `bytes crypto/sha256 encoding/hex fmt io os os/signal path/filepath regexp sort strconv strings sync syscall` |
| FILE-MODIFICATION | input-files forbidden | sources are opened read-only (`os.ReadFile`); verified by `TestSourcesAreReadOnly` |
| IDEMPOTENT | true | verified by `TestReSliceIntoSameDirectoryIsIdempotent` and `TestExampleTwoRunsAreByteIdentical` |
| PRESET-SYSTEM | systemd-style | no preset consumed; resolution order honoured (none present at any layer) |
| spec-hash | embedded everywhere | see below |

### Spec-hash embedding

`c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f` appears in:
every Go source file header comment (11 files) and the test file header,
`SpecSHA256` in `internal/pcdslice/types.go` (printed by `pcd-slice version`
as `spec:<hash>`), `SPEC_SHA256` in the `Makefile`, `# pcd-spec-sha256:` in
`pcd-slice.spec`, `X-PCD-Spec-SHA256:` in `debian/control`, the header
comments of `debian/rules`, `pcd-slice.1.md`, `README.md` and
`translation_report/translation-workflow.pikchr`, and the `Spec-SHA256:`
field of this report. No `Containerfile` was produced (OCI inactive), so the
`LABEL pcd.spec.sha256=` row does not apply. No placeholder value appears in
any artefact.

---

## STEPS ordering — how each BEHAVIOR was implemented

### BEHAVIOR: list (`internal/pcdslice/list.go`, driven from `run.go`)

1. **Parse spec against the section grammar; exit 2 on unreadable input** —
   `readMarkdown()` enforces the `.md` extension and readability and emits
   `error: cannot read {path}`; `ParseSpec()` runs the one-pass classifier.
2. **Transitive type closure per behavior** — `computeClosures()` →
   `closureOf()`: a worklist over the behavior block text, then over each
   newly added definition's text, matching whole words through one
   longest-first alternation regexp.
3. **Attribute invariants, examples and hints blocks** —
   `collectInvariants()`, `collectExamples()` + `attributeExample()`,
   `ParseHints()`.
4. **One line per behavior in specification order + total line** — behaviors
   are held in a slice in source order (never a map), printed as
   `name\ttypes\tinvariants\texamples\thints_blocks`, then
   `total: {n} behaviors`.
5. **Exit 0** — `List` always returns `ExitOK`; findings are not consulted
   (`TestListExitsZeroEvenWithFindings`).

Postconditions: no file is created or modified (`TestListAndCheckNeverWriteFiles`);
output order equals heading order.

### BEHAVIOR: check (`internal/pcdslice/check.go`)

1. **Parse as in list; exit 2 on unreadable input** — same loader.
2. **undefined-type** — `undefinedTypeFindings()` scans *only* type
   definition bodies, strips `//` comments, excludes the defining name, and
   applies the normative reference shape (`isReferenceShape`: ≥2 camel humps,
   each with at least one lowercase letter). Behavior text is never scanned
   (`TestUndefinedTypeIsAFindingInsideDefinitionsOnly`).
3. **unknown-binding** — raised in `collectInvariants()` at the binding line;
   the message is the unknown name verbatim.
4. **unassignable-example** — raised by the ladder in `attributeExample()`;
   a rung-4 tie lists candidates alphabetically (`sort.Strings`).
5. **unassignable-hints** — raised in `ParseHints()` when a heading's
   stripped text is a single lowercase identifier that names no known
   behavior. Hints preamble is never a finding.
6. **duplicate-behavior** — raised in `collectBehaviors()`, including the
   file-name collision after `/`→`-` replacement (slice step 2, hoisted into
   check so slice never reaches an inconsistent state).
7. **no-behaviors** — raised in `CheckFindings()` when the behavior slice is
   empty.
8. **Completeness** — `unassignedFindings()` scans the assignment array; any
   line still carrying the sentinel is reported with rule `unassigned`.
9. **stale-output** — `StaleFindings()`, only when `out=` is given and the
   directory exists: recomputes the provenance line and manifest for the
   current sources and compares; missing expected file, unexpected file,
   unexpected subdirectory, provenance difference and manifest content
   difference are one finding per file. A directory without `MANIFEST.tsv`
   is stale in full. An absent directory yields no findings.
10. **Write findings sorted by file then line; summary; exit 1 if any** —
    `SortFindings()` (file, line, rule, message), findings to stderr, summary
    `pcd-slice check: {n} findings` / `pcd-slice check: clean` to stdout.

### BEHAVIOR: slice (`internal/pcdslice/slice.go`, emission in `emit.go`)

1. **Run check without out=; on any finding print to stderr, nothing to
   stdout, write nothing, exit 1** — first statement of `Slice()`
   (`TestExampleSliceRefusesOnFindingsAndWritesNothing`).
2. **Compute all outputs in memory** — `BuildOutput()`; behavior file names
   are the name with `/`→`-`; collisions were already caught in step 1.
3. **Provenance line first in every emitted markdown file** —
   `ProvenanceLine()`: one `fmt.Sprintf` with fixed field order, paths as
   given, lowercase hex hashes, hints comma-separated in input order, no
   timestamp.
4. **Byte content fully determined by inputs** — sections in source order,
   closures alphabetical, manifest sorted; no map is ever ranged over during
   emission.
5. **Create OutDir if absent; remove only pcd-slice's own files; refuse
   foreign entries with exit 2** — the directory scan runs *before the first
   write*; a subdirectory, or a file whose first line is not
   `<!-- pcd-slice `, aborts with `error: foreign file in out: {path}`.
6. **Write every file, then the manifest; on write error exit 2 and remove
   files written in this run** — the run tracks written paths; `fail()`
   removes them; the manifest is written last and doubles as the commit
   marker.
7. **Summary to stdout, exit 0** —
   `pcd-slice: {n} bundles, preamble, manifest -> {out}`.

---

## Parsing approach

A single pass with a small explicit state machine, exactly as the Go hints
file prescribes, and deliberately **no markdown library**: a parser library
would recognise more structure than the spec's SECTION GRAMMAR admits.

- The source is read once into memory and split on `"\n"`, with any trailing
  `"\r"` kept attached to the line, so emitted bytes equal the bytes that
  came in and a CRLF source's bundles stay diffable against it. A final line
  without a newline counts as one line.
- Two boolean arrays (`inFence`, `isFence`) mark fenced content once; nothing
  inside a fence is ever interpreted as structure, and the `(binds: …)` and
  `(of: …)` suffixes are recognised only at end of line outside fences.
- Headings are lines beginning with `#` at column 0 outside a fence. A `## `
  or `# ` heading closes the section before it; `### ` and deeper never do,
  so nested `### EXAMPLE` blocks travel inside their behavior block.
- **The completeness invariant is materialised, not argued**: `Spec.Assign`
  holds one entry per source line, initialised to the `unassigned` sentinel.
  The section walk stamps every line; example attribution and invariant
  binding re-stamp their ranges. The completeness rule is then a slice scan.
  Outputs are projections of that array — which is what keeps "exactly once"
  true by construction. The single stated exception, a bound invariant, is
  assigned once to its first named behavior; the *emitter* replicates the
  lines into each further named bundle, so the assignment array stays a
  scan (`TestExampleBoundInvariantReplicatesOncePerBinding`).
- Whole-word matching for type names and behavior names uses one compiled
  alternation per set, sorted longest-first with `regexp.QuoteMeta` on each
  member, so `Plan` never shadows `PlanPath`.
- Determinism mechanics: behaviors are kept in a slice in source order; no
  map is ranged over while emitting; closures and manifest lines are sorted
  with `sort.Strings`; SHA-256 is computed over the raw bytes as read.

## Signal handling approach

`main()` installs one handler for `SIGTERM` and `SIGINT` before dispatch. On
either signal the handler calls `pcdslice.Abort()`, which removes every file
the current run has already written (tracked under a mutex as each write is
attempted) and then exits 0 — a clean exit with no partial output, as
`SIGNAL-HANDLING` requires. There is exactly one `os.Exit` in the normal
path, at the top of `main()` after `Run()` returns its exit code, so no
deferred cleanup is ever skipped. `list` and `check` write nothing, so for
them the handler is a no-op.

---

## Public API Surface

The exported surface of the implementation module. It must remain stable
across translations of spec version 0.1.1; a later translation may add to it
but not remove or rename entries without a spec version increment.

### Module `github.com/mge1512/pcd/tools/pcd-slice/internal/pcdslice`

Constants:

- `const SpecSHA256 = "c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f"`
- `const ExitOK = 0`
- `const ExitFindings = 1`
- `const ExitInvocation = 2`
- `const PreambleFile = "preamble.md"`
- `const ManifestFile = "MANIFEST.tsv"`
- `const Usage string`
- `const KindUnassigned Kind`, `KindPreamble Kind`, `KindBehavior Kind`, `KindExcluded Kind`
- `const RoleNone Role`, `RoleBlock Role`, `RoleExample Role`, `RoleInvariant Role`

Types:

- `type Kind int`
- `type Role int`
- `type Assignment struct { Kind Kind; Behavior string; Role Role }`
- `type Finding struct { Rule string; File string; Line int; Message string }`
- `type Range struct { Start int; End int }`
- `type Behavior struct { Name string; FileName string; Heading int; Block Range; Closure []string; NestedExamples int }`
- `type TypeDef struct { Name string; Line int; Body Range; Text string }`
- `type Invariant struct { Item Range; BindLine int; Binds []string; Owner string }`
- `type Example struct { Name string; Heading int; Block Range; Target string; Shared bool }`
- `type HintsBlock struct { Behavior string; Block Range }`
- `type Hints struct { Path string; SHA256 string; Lines []string; Blocks []HintsBlock; Preamble []Range; Findings []Finding }`
- `type Spec struct { Path string; SHA256 string; Lines []string; Assign []Assignment; Behaviors []Behavior; Types []TypeDef; Invariants []Invariant; Examples []Example; Findings []Finding; Rewrite map[int]string }`
- `type Output struct { Provenance string; Names []string; Files map[string][]byte; Manifest []byte; Bundles int }`

Functions:

- `func Run(args []string, version string, stdout, stderr io.Writer) int`
- `func List(spec *Spec, hints []*Hints, stdout io.Writer) int`
- `func Check(spec *Spec, hints []*Hints, out string, version string, stdout, stderr io.Writer) int`
- `func Slice(spec *Spec, hints []*Hints, out, version string, stdout, stderr io.Writer) int`
- `func CheckFindings(spec *Spec, hints []*Hints) []Finding`
- `func StaleFindings(out string, o *Output) ([]Finding, error)`
- `func SortFindings(fs []Finding)`
- `func BuildOutput(spec *Spec, hints []*Hints, version string) *Output`
- `func ProvenanceLine(spec *Spec, hints []*Hints, version string) string`
- `func ParseSpec(path string, data []byte) *Spec`
- `func ParseHints(path string, data []byte, spec *Spec) *Hints`
- `func SplitLines(data []byte) []string`
- `func HashBytes(b []byte) string`
- `func Abort()`

Methods:

- `func (k Kind) String() string`
- `func (f Finding) String() string`
- `func (s *Spec) BehaviorNames() []string`
- `func (s *Spec) HasBehavior(name string) bool`
- `func (s *Spec) Line(i int) string`
- `func (h *Hints) BlocksFor(behavior string) []HintsBlock`

### Module `main` (`cmd/pcd-slice`)

- `var version = "0.1.1"` (overridable via `-ldflags "-X main.version=…"`)
- `func main()`

---

## Produced artefact set

Every file traces to the template's DELIVERABLES table, the prompt's Reports
section, or the template's EXECUTION phase list. No unsolicited file was
written (no `.gitignore`, no `CHANGELOG.md`, no CI config, no lock file, no
build directory).

| Deliverable row | Files produced |
|---|---|
| source (required) | `cmd/pcd-slice/main.go`, `internal/pcdslice/{types,parse,closure,attribute,hints,emit,check,list,slice,run}.go`, `go.mod` |
| build (required) | `Makefile` (`build`, `test`, `install`, `clean`, `man`, `dist`, plus `vendor`, `all`) |
| docs (required) | `README.md` |
| man (required) | `pcd-slice.1.md`, `pcd-slice.1` (generated with `pandoc 2.18`) |
| license (required) | `LICENSE` (SPDX identifier + authoritative URL; full text not reproduced) |
| RPM (required) | `pcd-slice.spec` |
| DEB (required) | `debian/control`, `debian/changelog`, `debian/rules`, `debian/copyright` |
| public-api (required) | `## Public API Surface` above |
| report (required) | `TRANSLATION_REPORT.md` |
| spec-hash (required) | embedded as listed above |
| Phase 1 (EXECUTION) | `independent_tests/claude-opus-5/pcd_slice_test.go` |
| Phase 4 (EXECUTION) | `translation_report/translation-workflow.pikchr` |
| version single source | `VERSION` |
| OCI / PKG (supported) | **not produced** — not active in the resolved preset |

`go.sum` was not written: the module has no dependencies, so `go mod tidy`
produces none, and the template does not name a lock file as a deliverable.
No `vendor/` tree exists for the same reason — `make dist` therefore emits
the source tarball only and says so; the RPM has no `Source1:`.

---

## Phase 6 — Compile gate

All steps executed in `/tmp/pcd-output/` with Go 1.26.6, unprivileged, with
`GOPATH`/`GOCACHE` under `$HOME`.

| Step | Command | Result |
|---|---|---|
| 1 Dependency resolution | `go mod tidy` | **pass** — no dependencies added; standard library only |
| 2 Compilation | `go build ./...` | **pass** |
| 2 Vet | `go vet ./...` | **pass** |
| 2 Format | `gofmt -l .` | **pass** (empty output) |
| 3 Build target | `make build` | **pass** — `./pcd-slice` at the project root |
| 3 Static check | `file ./pcd-slice` | **pass** — "ELF 64-bit … statically linked, stripped" |
| 3 Translator test run | `make test` → `go test ./independent_tests/claude-opus-5/...` | **pass** — ok, 0 failures, 0 skips |
| 3 Smoke: `./pcd-slice version` | | **pass** — prints `pcd-slice 0.1.1` and `spec:c200af55…` |
| 3 Smoke: `./pcd-slice help` | | **pass** — usage, exit 0 |
| 3 Smoke: `./pcd-slice check spec=missing.md` | | **pass** — `error: cannot read missing.md`, exit 2 |
| 3 Man page | `make man` (pandoc 2.18) | **pass** |
| 3 Release tarball | `make dist` | **pass** — `pcd-slice-0.1.1.tar.gz` with single top-level dir `pcd-slice-0.1.1/` (removed again afterwards; build output is not a deliverable) |
| 4 Test-author test run | — | **not applicable** (single-LLM run) |

**Dogfood check (not a deliverable, run in a scratch directory and removed):**
`pcd-slice check spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md` on its
own specification reports `pcd-slice check: clean`, and
`pcd-slice list` reports `list 3 1 1 0 / check 4 2 3 0 / slice 4 5 12 0 /
total: 3 behaviors` — 16 attributed examples, matching the 16 EXAMPLE blocks
in the spec, and bound-invariant counts matching the five INVARIANTS
bindings. `pcd-slice slice` on its own spec produced a preamble, three
bundles and a manifest.

**Clean-run status:** the delivered tree contains no scratch files, no
`build/` directory, no tarball and no compiled binary; `./pcd-slice` is
rebuilt by `make build` or by the suite's `TestMain`. The test suite creates
every fixture in a per-test `t.TempDir()` and leaves no residue; a second
identical run behaves identically.

---

## Test results — translator suite

`independent_tests/claude-opus-5/pcd_slice_test.go` — 39 test functions and
11 sub-tests. **All pass. 0 fail, 0 skip.** Black-box throughout: every test
runs `../../pcd-slice` through `exec.Command` with the working directory set
to the test's own `t.TempDir()` and a sandboxed `HOME`, and asserts on
stdout, stderr and exit code. No test imports the implementation package.

| Test | Covers | Result |
|---|---|---|
| `TestVersionPrintsToolNameAndSpecHash` | DEPLOYMENT version verb, spec-hash embedding | pass |
| `TestExampleClosureIsTransitive` | EXAMPLE closure_is_transitive | pass |
| `TestExampleUntaggedInvariantIsGlobal` | EXAMPLE untagged_invariant_is_global | pass |
| `TestExampleBoundInvariantLandsInNamedBundles` | EXAMPLE bound_invariant_lands_in_named_bundles | pass |
| `TestExampleBindingToUnknownBehaviorIsAFinding` | EXAMPLE binding_to_unknown_behavior_is_a_finding | pass |
| `TestExampleNestedExampleTravelsWithItsBehavior` | EXAMPLE nested_example_travels_with_its_behavior | pass |
| `TestExampleTopLevelExampleAttributedByWhenLine` | EXAMPLE top_level_example_attributed_by_when_line | pass |
| `TestExampleChangelogIsExcluded` | EXAMPLE changelog_is_excluded | pass |
| `TestExampleSliceRefusesOnFindingsAndWritesNothing` | EXAMPLE slice_refuses_on_findings_and_writes_nothing | pass |
| `TestExampleTwoRunsAreByteIdentical` | EXAMPLE two_runs_are_byte_identical, INVARIANT Determinism | pass |
| `TestExampleStaleOutputIsOneFindingPerFile` | EXAMPLE stale_output_is_one_finding_per_file | pass |
| `TestCheckWithCurrentOutIsClean` | check step 9 negative case | pass |
| `TestCheckAbsentOutDirProducesNoStaleFindings` | check step 9 "absent out produces no findings" | pass |
| `TestExampleForeignFileInOutRefuses` | EXAMPLE foreign_file_in_out_refuses | pass |
| `TestExampleForeignSubdirectoryRefuses` | EXAMPLE foreign_subdirectory_refuses | pass |
| `TestExampleOfSuffixAttributesAndSharedGoesToPreamble` | EXAMPLE of_suffix_attributes_and_shared_goes_to_preamble | pass |
| `TestExampleUniqueTextMatchAttributes` | EXAMPLE unique_text_match_attributes | pass |
| `TestExampleRungFourTieIsAFinding` | EXAMPLE rung_four_tie_is_a_finding | pass |
| `TestUnplaceableExampleIsAFinding` | check step 4, unplaceable example | pass |
| `TestExampleBoundInvariantReplicatesOncePerBinding` | EXAMPLE bound_invariant_replicates_once_per_binding, INVARIANT Completeness | pass |
| `TestNoBehaviorsIsAFinding` | check step 7 | pass |
| `TestUndefinedTypeIsAFindingInsideDefinitionsOnly` | check step 2 incl. the "behavior text is never scanned" clause | pass |
| `TestDuplicateBehaviorIsAFinding` | check step 6 | pass |
| `TestUnassignableHintsHeadingIsAFinding` | check step 5 | pass |
| `TestHintsPreambleIsNotAFinding` | check step 5 negative case | pass |
| `TestFindingsAreSortedByFileThenLine` | check step 10 ordering and line format | pass |
| `TestHintsBlocksTravelWithTheirBehavior` | hints attribution grammar, list hints counts | pass |
| `TestSliceSummaryAndOutputSet` | slice step 7, POSTCONDITION "out contains exactly the manifest's files plus MANIFEST.tsv" | pass |
| `TestProvenanceHeaderIsFirstLineOfEveryMarkdownFile` | slice step 3, INVARIANT Provenance | pass |
| `TestManifestListsEveryFileExceptItselfWithMatchingHashes` | POSTCONDITION "each file's sha256 equals its manifest entry" | pass |
| `TestReSliceIntoSameDirectoryIsIdempotent` | IDEMPOTENT, slice step 5 | pass |
| `TestStaleBundleIsRemovedOnReSlice` | slice step 5 removal clause | pass |
| `TestSourcesAreReadOnly` | INVARIANT "Sources are read-only" | pass |
| `TestListAndCheckNeverWriteFiles` | POSTCONDITION "list and check never write files" | pass |
| `TestListExitsZeroEvenWithFindings` | list step 5 | pass |
| `TestInvocationErrors` (11 sub-tests) | ExitCode 2 paths: no verb, unknown verb, POSIX flags, bare positional, unknown key, missing spec=, slice without out=, missing file, non-.md file, missing hints file | pass |
| `TestHelpExitsZeroAndPrintsUsage` | bare-word help verb | pass |
| `TestBundleCarriesItsBehaviorBlockVerbatim` | INVARIANT "a behavior appears in exactly one bundle" | pass |
| `TestPreambleCarriesSharedContext` | slice OUTPUTS preamble contract | pass |

## Test results — test-author suite

Not applicable: no test-author suite was present at input (single-LLM run).

## Test Refinements

No test was edited after any test run, and no implementation change was made
in response to a failing test: the suite passed on its first execution.

| Test | Result before | Action | Rationale |
|------|---------------|--------|-----------|
| *(all 39 test functions and 11 sub-tests)* | passed | none | — |

---

## Per-EXAMPLE confidence

Confidence is **Medium** for every row: Tests-First-Compliance is `yes` and
each row's named test passes with no live external service, but no
test-author suite exists, which the prompt's definition makes the deciding
factor (High requires an independent test-author suite that also passes).

| EXAMPLE | Confidence | Verification method | Unverified claims |
|---|---|---|---|
| closure_is_transitive | Medium | `TestExampleClosureIsTransitive` — asserts `record 3 1 1 0`, `Requires-Types: Item, Path, Plan` verbatim, `Unused` present in preamble and absent from the bundle, exit 0 | none |
| untagged_invariant_is_global | Medium | `TestExampleUntaggedInvariantIsGlobal` — invariant text present in `preamble.md`, absent from both bundles; slice exit 0 | none |
| bound_invariant_lands_in_named_bundles | Medium | `TestExampleBoundInvariantLandsInNamedBundles` — in `reset.md`, not in `purge.md`, not in `preamble.md` | none |
| binding_to_unknown_behavior_is_a_finding | Medium | `TestExampleBindingToUnknownBehaviorIsAFinding` — stderr `unknown-binding: spec.md:{computed line}: renmae`, stdout `pcd-slice check: 1 findings`, exit 1 | none |
| nested_example_travels_with_its_behavior | Medium | `TestExampleNestedExampleTravelsWithItsBehavior` — example inside `verify.md`, no attribution heading emitted, absent from preamble | none |
| top_level_example_attributed_by_when_line | Medium | `TestExampleTopLevelExampleAttributedByWhenLine` — example under `## Attributed examples from spec.md` in `lint.md` | the spec says "a heading that names its source" without fixing its wording; the wording chosen is asserted, not derived from the spec |
| changelog_is_excluded | Medium | `TestExampleChangelogIsExcluded` — no changelog row and no `## Changelog` in any emitted file; check clean (lines counted excluded, not unassigned) | none |
| slice_refuses_on_findings_and_writes_nothing | Medium | `TestExampleSliceRefusesOnFindingsAndWritesNothing` — stderr contains `duplicate-behavior`, stdout empty, `spec.d` does not exist, exit 1 | none |
| two_runs_are_byte_identical | Medium | `TestExampleTwoRunsAreByteIdentical` — per-file SHA-256 maps of two output trees compared, manifests compared byte-wise | byte-identity across two *builds* is not asserted (the tool version string is a build input; the spec's example compares two runs) |
| stale_output_is_one_finding_per_file | Medium | `TestExampleStaleOutputIsOneFindingPerFile` — exactly one `stale-output` line per emitted file (4/4), each naming its path, exit 1 | a hand-edited output file whose provenance line is intact is not detected: the spec scopes the comparison to provenance headers plus manifest content |
| foreign_file_in_out_refuses | Medium | `TestExampleForeignFileInOutRefuses` — stderr `error: foreign file in out: spec.d/notes.txt`, `notes.txt` byte-unchanged, nothing else written, exit 2 | none |
| of_suffix_attributes_and_shared_goes_to_preamble | Medium | `TestExampleOfSuffixAttributesAndSharedGoesToPreamble` — first example only in `check.md`, second only in `preamble.md` | none |
| unique_text_match_attributes | Medium | `TestExampleUniqueTextMatchAttributes` — example in `verify.md` only, prose WHEN line falls through rung 3 | none |
| rung_four_tie_is_a_finding | Medium | `TestExampleRungFourTieIsAFinding` — stderr contains `unassignable-example` and `purge, reset` (alphabetical), exit 1 | none |
| bound_invariant_replicates_once_per_binding | Medium | `TestExampleBoundInvariantReplicatesOncePerBinding` — exactly one occurrence in `check.md` and in `slice.md`, zero elsewhere, two in total | none |
| foreign_subdirectory_refuses | Medium | `TestExampleForeignSubdirectoryRefuses` — stderr `error: foreign file in out: spec.d/notes`, exit 2, directory untouched | none |

Additional verified claims beyond the EXAMPLES: all five INVARIANTS
(read-only sources, determinism, completeness/replication, provenance,
one-bundle-per-behavior), the PRECONDITIONS (existence, readability, `.md`
extension), the POSTCONDITIONS of all three behaviors, every `Finding` rule,
and every `ExitCode` value.

---

## Specification ambiguities and conservative interpretations

1. **`MANIFEST.tsv` and the foreign-entry rule.** Slice step 5 removes files
   the manifest would not list *only if* they carry a provenance line, and
   calls anything else foreign. Read literally, `MANIFEST.tsv` — which
   carries no provenance line and is not listed in itself — would be foreign,
   making it impossible to re-slice into a directory pcd-slice itself wrote,
   which contradicts `IDEMPOTENT: true` and the staleness design of check
   step 9. Conservative resolution: `MANIFEST.tsv` is recognised as
   pcd-slice's own file by name and overwritten, never treated as foreign.
2. **Heading wording for appended material.** Slice's OUTPUTS require
   attributed material "each under a heading that names its source" but do
   not fix the text. Chosen, and held stable for determinism:
   `## Attributed examples from {spec}`, `## Bound invariants from {spec}`,
   `## Hints from {hints-path}`, `## Hints preamble from {hints-path}`.
3. **`Requires-Types:` with an empty closure.** The format for a behavior
   whose closure is empty is unspecified; the label is emitted alone
   (`Requires-Types:`), keeping the line count of every bundle uniform.
4. **"Every output line traces to exactly one source line."** Taken to govern
   *source-derived* content. The provenance line, the `Requires-Types:` line
   and the four source-naming headings are tool-generated and are required
   by the same STEPS, so they are necessarily exempt.
5. **Suffix stripping.** The Go hints file directs stripping the `(of: …)`
   suffix from an example heading before emitting it. This is the one place
   where an emitted line is not byte-identical to its source line; it is
   implemented through an explicit per-line rewrite map so the assignment
   array (and hence the completeness scan) is untouched.
6. **`(of: …)` naming an unknown behavior.** The ladder says the first rung
   that *fires* wins. A suffix naming neither `shared` nor a known behavior
   is treated as not firing: the ladder continues to rungs 3 and 4, and if
   they cannot place the example it becomes an `unassignable-example`
   finding. The alternative — an immediate finding — would swallow specs
   whose suffix is a typo but whose WHEN line is unambiguous.
7. **Behavior-name occurrences in a hints heading.** The grammar says a
   heading "contains a BehaviorName as a whole word"; it does not say what
   happens when a heading names two. The first occurrence by position wins
   (deterministic, and no finding rule covers the case).
8. **Trailing blank lines.** Each emitted chunk has its trailing blank lines
   trimmed and is separated from the next by exactly one blank line. Within
   a chunk, projection is pure: blank lines left by removed sections stay
   where they were.
9. **Invariant items spanning several lines.** `InvariantBinding` speaks of
   "an invariant line"; real specs (including this one) wrap invariants over
   several lines with the suffix on the last. The unit of assignment and
   replication is therefore the whole list item, which is what "the line
   appears once per named bundle" means for a wrapped item.
10. **Section-heading recognition.** `## TYPES`, `## INVARIANTS` and
    `## EXAMPLES` are recognised by exact heading text (after trimming);
    anything else is preamble, per the grammar's catch-all. A changelog is
    any `## ` heading containing the word "Changelog" as a whole word, in
    any case.
11. **`help` verb.** The spec's DEPLOYMENT section lists `version`, `list`,
    `check` and `slice` but not `help`. It is implemented because the
    milestones hints file makes `<binary> help` part of the compile gate and
    the template's `CLI-ARG-STYLE: bare-words` permits it; it prints usage
    to stdout and exits 0, and affects no other behaviour.

## Rules that could not be implemented exactly as written

None. Every STEPS entry of every BEHAVIOR is implemented in the order
written, with the interpretations recorded above where the text underdetermined
the outcome.

## Deviations from the deployment template

1. **Implementation package directory is `internal/pcdslice`, not
   `internal/pcd-slice`.** The per-language source layout table gives
   `internal/<n>/`, and `<n>` is `pcd-slice`. A Go package identifier cannot
   contain a hyphen; while the *directory* could keep it, the resulting
   split between directory name and package name is a known source of import
   confusion. The hyphen is therefore dropped in the directory as well. The
   module identity itself is unaffected and remains
   `github.com/mge1512/pcd/tools/pcd-slice`.
2. **No vendor tarball from `make dist`, no `Source1:` in the RPM.** The
   template requires a companion vendor tarball "for languages with
   dependency vendoring". The spec's DEPENDENCIES section declares no
   dependency beyond the standard library, so `go mod vendor` produces
   nothing and a vendor tarball would be empty. `make dist` emits the vendor
   tarball if and only if a `vendor/` tree exists, and prints why it did not
   otherwise; a `vendor` target is provided for the day a dependency
   appears.
3. **RPM `Version:` is the literal `0.1.1`, not `%(cat %{_sourcedir}/VERSION)`.**
   The template allows either the `%(cat …)` form or injection at tarball
   build time. Since `make dist` derives the tarball name and top-level
   directory from `VERSION`, and `VERSION` ships inside the tarball rather
   than beside it in `%_sourcedir`, the literal form is the one that
   actually builds in OBS. The spec file carries a comment stating the
   invariant.
4. **No `Containerfile` and no `.pkgbuild`.** `OCI` and `PKG` are `supported`
   OUTPUT-FORMATs, produced only when active in the resolved preset. No
   preset is present at any layer and the spec declares Linux only in v1, so
   neither was produced. This is compliance with, not deviation from, the
   "No unsolicited deliverables" rule; it is listed here so the absence is
   explicit.

## Dependency versions

No third-party dependency is declared, fabricated or required: `go.mod`
contains the module line and `go 1.22` only. Nothing needs manual version
verification before building.
