
# Contributing to PCD

Thank you for your interest in contributing to the Post-Coding Development (PCD) project!

## Project Overview

PCD is a paradigm where:
- Domain experts write specifications in structured Markdown
- AI translates specifications into verified implementations
- Engineers never write implementation code directly
- Target language is derived from deployment templates, not declared by spec authors

**Key distinction from "vibe coding":**
- Vibe coding: humans write code + AI suggests
- PCD: humans write specs only; AI generates all implementation code
- If generated code is wrong: fix the spec, not the code

## Getting Started

### Prerequisites

- Git
- Go (for building tools)
- Basic understanding of Markdown
- Familiarity with specification writing

### Repository Structure

```
pcd/
├── doc/                           # Documentation
│   ├── whitepaper.md             # Canonical whitepaper
│   └── executive-brief.md
├── templates/                     # Deployment templates
│   ├── cli-tool.template.md      # Complete, production-ready
│   ├── mcp-server.template.md    # Complete, production-ready
│   └── ...
├── hints/                         # Translator hints files (see Hints Files below)
├── tools/                         # PCD tooling
│   └── pcd-lint/                 # Specification validator
├── examples/                      # Example specifications
└── prompts/                       # AI translator prompts
```

---

## How to Contribute

### 1. Specification Writing

When writing PCD specifications:

#### Required META Section Fields

Every PCD spec must include:

```markdown
## META

Deployment:        <template-name>
Version:           MAJOR.MINOR.PATCH
Spec-Schema:       MAJOR.MINOR.PATCH
Author:            Name <email>
License:           SPDX-identifier
Verification:      none | lean4 | fstar | dafny | custom
Safety-Level:      QM | ASIL-A | ...
```

For multi-component projects, additionally:
```markdown
Interface-Version: MAJOR.MINOR.PATCH
Imports:           name: path#INTERFACE@>=version
```

#### Behavior Sections

Use these section types:
- `## BEHAVIOR: {n}` — user-facing operation
- `## BEHAVIOR/INTERNAL: {n}` — internal implementation logic

Every BEHAVIOR block **must** include `PRECONDITIONS:`, `STEPS:`, and `POSTCONDITIONS:`.
`STEPS:` are numbered, imperative, and include explicit error exits. They describe the
algorithm — not just the contract. A `MECHANISM:` annotation may follow any step where
the implementation pattern matters for correctness beyond what postconditions capture.

```markdown
## BEHAVIOR: {name}

INPUTS:
  ...

PRECONDITIONS:
  - [condition]

STEPS:
  1. [first action]; on failure → [error action].
  2. [next action].
     MECHANISM: [how, when the how matters for correctness]

POSTCONDITIONS:
  - [outcome]

ERRORS:
  - ERR_X if [condition]
```

#### Examples Section

Must include at least one complete test case. EXAMPLES may use multiple WHEN/THEN pairs
to express multi-pass or multi-step behaviour (e.g. Kubernetes reconcilers). Each
WHEN/THEN pair represents one invocation. Single-pass EXAMPLES remain valid.

```markdown
## EXAMPLES

EXAMPLE: {descriptive_name}
GIVEN:
  [Initial conditions]

WHEN:  [Action — optionally labelled, e.g. "reconcile runs (pass 1)"]
THEN:
  [Expected outcome]

WHEN:  [Action — next pass or step, if multi-pass]
THEN:
  [Expected outcome]
```

#### Interfaces Section (optional, recommended for complex components)

Declare module boundary contracts and required test doubles. This section prevents
translator discretion on abstraction layer decisions and makes independent tests
infrastructure-free.

```markdown
## INTERFACES

AdapterName {
  required-methods:
    MethodName(ctx, InputType) → (OutputType, error)
  implementations-required:
    production:  RealAdapter
    test-double: FakeAdapter {
      configurable fields: ...
      state machine: ...
    }
}
```

#### Dependencies Section (optional)

Declare external library requirements. Translators are bound by these rules and
must not fabricate dependency versions.

```markdown
## DEPENDENCIES

module/path:
  minimum-version: vX.Y.Z
  rationale: [why this version]
  do-not-fabricate: true     # translator must not invent commit hashes
  hints-file: template.language.library.hints.md
```

#### Invariants Section

Tag each invariant with `[observable]` (verifiable by external observation or the
independent test suite) or `[implementation]` (verifiable by code review or static
analysis only). Untagged invariants are a `pcd-lint` warning.

```markdown
## INVARIANTS

- [observable]      [externally verifiable property]
- [implementation]  [code-review-only property]
```

#### Language Neutrality

A spec that may be translated into more than one implementation language must
contain no language-specific constructs in its BEHAVIOR blocks, TYPES,
INTERFACES, INVARIANTS, or MILESTONE acceptance criteria. Language-specific
constructs belong exclusively in hints files. Specifically:

- BEHAVIOR STEPS must use abstract operations: "create directory recursively",
  not `os.MkdirAll()` or `std::fs::create_dir_all()`
- MILESTONE acceptance criteria must use CLI invocations and universal
  tools (`jq`, `grep`, `file`, `stat`), not language build commands such as
  `go build` or `cargo build`
- Type syntax must be pseudocode (`ScopeWrapper<T>`), not any language's
  generic syntax
- Interface method signatures must use neutral notation, not any language's
  method declaration syntax

A useful test: could this spec section be read and understood by a developer
who knows the domain but has not decided on a target language yet? If not,
language-specific content has leaked into the spec.

#### Output Format is a Spec Concern

The spec must fully specify what a BEHAVIOR writes to disk — file format,
field names, schema. Hints files describe only how to implement the writing
in a given language, not what to write. A format decision documented only
in a hints file is invisible to spec reviewers and creates silent
inconsistency risk across language ports.

#### Privileged Components

When a component requires elevated privileges to run (root, sudo, capability
bits, hardware access), the spec author should:

1. Ensure M0 (scaffold) acceptance criteria are fully verifiable without
   privilege — compile gate, `--help`, `--version`, invocation errors.
2. Phrase M1+ acceptance criteria as observable output checks (JSON field
   values, file existence, file sizes) that a human can run in the target
   environment after each milestone pass.
3. Require the INTERFACES section to declare test doubles (e.g. FakeFilesystem,
   FakeCommandRunner) that allow unit tests to exercise logic without privilege.
4. Accept that the translator will report Low confidence for runtime
   verification items. This is correct and honest, not a failure.

---

#### MILESTONE Section (optional)

Declares named, versioned, self-consistent subsets of the full spec that can
be independently translated, tested, and released. Use this for large components
where translating the entire spec in one pass exceeds the context window or is
otherwise impractical. Two complete real-world implementations (Go and Rust) of
a 35-BEHAVIOR, 2900-line specification have been produced using this mechanism.

Each MILESTONE is a projection of the spec — it names which BEHAVIORs are
fully implemented at this stage and which remain as stubs. The spec itself
remains the complete source of truth.

**Status field — pipeline state machine:**

| Status | Meaning |
|---|---|
| `pending` | Not yet attempted. Default for newly written milestones. |
| `active` | Currently being translated. Set by the agent pipeline. Exactly one milestone may be active at a time. |
| `failed` | Compile gate or acceptance criteria did not pass. Set by the agent pipeline. Human reviews. |
| `released` | All gates passed. Set by the agent pipeline. Frozen — do not modify. |

Exactly one MILESTONE may have `Status: active` at any time (RULE-15).
The pipeline agent advances the cursor; humans do not need to edit status
fields manually.

**Full MILESTONE syntax:**

```markdown
## MILESTONE: {version}
Status: pending
Scaffold: true | false          # optional; default false — see below
Hints-file: {filename}          # optional; comma-separated list

Included BEHAVIORs:
  {behavior-name-1}, {behavior-name-2}, ...

Deferred BEHAVIORs:
  {behavior-name-3}, {behavior-name-4}, ...

Acceptance criteria:
  {one criterion per line — shell command that exits 0 on pass}
```

**The scaffold milestone (`Scaffold: true`):**

For any component above roughly 500 lines of generated code, or with more than
10 BEHAVIORs, the first milestone should be a scaffold-only pass. The scaffold
translator creates all files, all types, all function signatures, and all stub
bodies for the **entire component** — not just the milestone's own BEHAVIORs.
The sole acceptance criterion is a clean compile.

All subsequent milestones then operate on a known, stable foundation. They
replace stub bodies with real implementations. They never create new files,
never restructure packages, never add new types.

When `Scaffold: true`:
- `Included BEHAVIORs` lists **all** BEHAVIORs in the spec (the complete set)
- `Deferred BEHAVIORs` is empty or omitted
- The acceptance criteria must include a compile gate as the first criterion
- The translator must not implement any real logic beyond what compiles

When `Scaffold: false` (or the field is absent):
- Behaviour is unchanged from a standard milestone
- The translator fills in real implementations for `Included BEHAVIORs` only
- All files already exist from the scaffold pass; only function bodies change

Empirical calibration: a 35-BEHAVIOR spec produced approximately 1600 lines of
Go scaffold across 4 files in one translator session. The scaffold held without
modification through seven subsequent implementation milestones.

**The `Hints-file:` field:**

Lists hints files the translator must read before beginning work on this
milestone. Multiple files are comma-separated. The translator reads all listed
files before writing any code. For scaffold milestones, the language-specific
milestones hints file is especially important — it specifies file layout, stub
conventions, and infrastructure implementations that must not themselves be
stubbed.

**The stub contract:**

A stub must compile and return the correct zero value for its declared output
type. For any output type that serialises to a JSON object, the stub must
return an initialised empty object — never a null reference. A null reference
serialises to JSON `null`; an initialised empty object serialises to `{}` or
`{"_elements":[]}`. Only the latter is schema-compatible with consumers that
expect an object. The language-specific milestones hints file gives concrete
examples of what "initialised empty object" means in each target language.

**Acceptance criteria format:**

Acceptance criteria should be expressed as shell commands that exit 0 on
pass and non-zero on failure. This makes them automatable by a pipeline
agent without any parsing. Examples:

```
./sitar version | grep -q "^sitar "
./sitar all outdir=/tmp/test && test -s /tmp/test/general.json
jq '.cpu._elements | length > 0' /tmp/test/json/cpu.json | grep -q true
```

For privileged components, M0 criteria must be runnable without privilege;
M1+ criteria may require a privileged runtime environment and should be
phrased accordingly so the human verifier knows what to run.

**Rules:**
- Every BEHAVIOR named in `Included BEHAVIORs` or `Deferred BEHAVIORs` must
  exist in the spec (validated by RULE-16).
- Together, `Included` + `Deferred` need not cover every BEHAVIOR in the spec —
  BEHAVIORs not mentioned in any milestone are always translated in full.
- At most one MILESTONE may have `Scaffold: true`, and if present, it must be
  the first milestone in document order (RULE-17).
- MILESTONE sections are non-normative for pcd-lint rule purposes — they do not
  affect RULE-01 through RULE-14. RULE-15, RULE-16, and RULE-17 apply to them.
- The `## DELTA` section (for single-pass work orders) and `## MILESTONE`
  sections serve different purposes and may coexist. DELTA is ephemeral;
  MILESTONEs are persistent across translation passes.

---

### 2. Template Development

Deployment templates define how specifications are translated to code.

#### Template Structure
- Use `{curly_braces}` for placeholders (not `<angle_brackets>`)
- Include a DELIVERABLES section listing all files to generate
- Specify constraints and requirements clearly

#### Available Templates (v0.3.20)

| Template | Status | Default Lang | Notes |
|---|---|---|---|
| `cli-tool` | Complete | Go | Production-ready + EXECUTION + man pages |
| `mcp-server` | Complete | Go | stdio + streamable-HTTP + man pages |
| `cloud-native` | Complete | Go | SLE-BCI; BUILD-GATE + EXECUTION |
| `backend-service` | Complete | Go | Production-ready + EXECUTION + man pages |
| `gui-tool` | Complete | C++/Rust/Dart | Qt6/Tauri/Flutter; EXECUTION: none |
| `python-tool` | Complete | Python | QM only; POSIX flags; man pages |
| `library-c-abi` | Complete | C | Section 3 man pages; EXECUTION: none |
| `verified-library` | Complete | C | Section 3 man pages; EXECUTION: none |
| `project-manifest` | Complete | N/A | Architect artifact; EXECUTION: none |

#### Hints Files

Hints files contain implementation knowledge that belongs neither in the spec
(which must be language-agnostic) nor in the deployment template (which covers
language and deployment conventions, not library internals or milestone-specific
patterns). They are advisory only — they cannot override spec invariants.

**Three-layer naming convention:**

```
hints/
  <template>.<language>.milestones.hints.md   # generic scaffold patterns;
                                               # reusable across all components
                                               # using this template + language
  <component>.implementation.hints.md         # component-specific, language-neutral;
                                               # file grouping, required field names,
                                               # known failure modes from prior runs
  <template>.<language>.<library>.hints.md    # library-specific API shapes,
                                               # version selection, known gotchas
                                               # (existing convention)
```

The `<template>.<language>.milestones.hints.md` file is the scaffold-first
companion: it specifies file layout, stub conventions, infrastructure
implementations (e.g. how to implement `CommandRunner.Run` as a thin wrapper),
serialisation patterns, and compile gate commands for the target language.
It is reusable across all components using that template and language.

The `<component>.implementation.hints.md` file is component-specific and
language-neutral. It contains recommended file groupings mapped to BEHAVIOR
names, required output field names for JSON schema compatibility, function names
that must exist, and known failure modes from previous translation runs.

Currently shipped:
- `hints/cli-tool.go.milestones.hints.md`
- `hints/cli-tool.rs.milestones.hints.md`
- `hints/cloud-native.go.go-libvirt.hints.md`
- `hints/cloud-native.go.golang-crypto-ssh.hints.md`
- `hints/mcp-server.go.mcp-go.hints.md`
- `hints/python-tool.hints.md`

Specs reference hints files via their `## DEPENDENCIES` section or via the
`Hints-file:` field on individual `## MILESTONE:` sections.

---

### 3. Tool Development

#### pcd-lint Validation Rules

The `pcd-lint` tool validates specifications against these rules:

- **RULE-01**: Required sections present (META, TYPES, BEHAVIOR, PRECONDITIONS, POSTCONDITIONS, INVARIANTS, EXAMPLES)
- **RULE-02**: META fields validation
- **RULE-03**: Deployment template resolution
- **RULE-04**: Deprecated META fields detection
- **RULE-05**: Verification field validation
- **RULE-06**: EXAMPLES structure validation (GIVEN/WHEN/THEN; multi-pass permitted)
- **RULE-07**: EXAMPLES content validation
- **RULE-08**: BEHAVIOR blocks must contain STEPS (v0.3.12+)
- **RULE-09**: INVARIANTS entries must carry `[observable]` or `[implementation]` tag (v0.3.12+, warning)
- **RULE-10**: Negative-path EXAMPLE required for BEHAVIOR with error exits (v0.3.13+)
- **RULE-11**: TOOLCHAIN-CONSTRAINTS section structure (v0.3.13+)
- **RULE-12**: Cross-section consistency: identifiers, types, file names (v0.3.13+)
- **RULE-13**: Constraint: field value on BEHAVIOR headers (v0.3.13+)
- **RULE-14**: EXECUTION section required in deployment templates (v0.3.16+)
- **RULE-15**: MILESTONE section structure and single-active constraint (v0.3.21+)
- **RULE-16**: MILESTONE BEHAVIOR names exist in spec (v0.3.21+)
- **RULE-17**: At most one scaffold milestone; scaffold milestone must appear first (v0.3.21+)

#### CLI Conventions

All PCD tools follow these conventions:
- Key=value syntax: `pcd-lint strict=true spec.md`
- Bare words for commands: `pcd-lint list-templates`
- **NO `--flag` style ever** (firm decision)
- stderr for diagnostics; stdout for summaries and lists
- Exit codes: 0 = valid, 1 = errors/strict warnings, 2 = invocation error

### 4. Code Style and Standards

- **Go**: Follow standard Go conventions
- **Markdown**: Use consistent formatting
- **Filenames**: No version numbers (Git handles versioning)
- **Placeholders**: Always use `{curly_braces}`

---

## Licensing

The project uses a dual-license model:

| Artifact Type | License |
|---|---|
| Whitepaper, specs, templates, examples | CC-BY-4.0 |
| Tools (pcd-lint, etc.) | GPL-2.0-only |

**Rationale for GPL-2.0-only on tools:** Mirrors the Linux kernel model.
Encourages everybody to contribute changes back to the validator toolchain.

---

## Submission Guidelines

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes following the guidelines above
4. Test your changes with `pcd-lint` if applicable
5. Commit with clear, descriptive messages
6. Push to your fork and submit a pull request

### Commit Message Format

```
component: Brief description of change

Longer explanation of what changed and why, if necessary.
Reference any relevant issues.
```

### Testing

- Run `pcd-lint` on any specifications you modify
- For Go code: ensure `go build ./...` succeeds
- Test templates with example specifications

---

## Communication

- **Email**: pcd@mailbox.org
- **Issues**: Use GitHub issues for bug reports and feature requests
- **Discussions**: Use GitHub discussions for questions and ideas

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all
contributors. Please be respectful and professional in all interactions.

## Questions?

If you have questions about contributing, please:
1. Check existing issues and discussions
2. Review the whitepaper in `doc/whitepaper.md`
3. Contact us at pcd@mailbox.org

Thank you for contributing to PCD!
