
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
│   └── ...                       # Various template stubs
├── tools/                         # PCD tooling
│   └── pcd-lint/                # Specification validator
├── examples/                      # Example specifications
└── prompts/                       # AI translator prompts
```

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

#### MILESTONE Section (optional)

Declares named, versioned, self-consistent subsets of the full spec that can
be independently translated, tested, and released. Use this for large components
where translating the entire spec in one pass exceeds the context window or is
otherwise impractical.

Each MILESTONE is a projection of the spec — it names which BEHAVIORs are
fully implemented at this stage and which are deferred stubs. The spec itself
remains the complete source of truth.

**Status field — pipeline state machine:**

| Status | Meaning |
|---|---|
| `pending` | Not yet attempted. Default for newly written milestones. |
| `active` | Currently being translated. Set by the agent pipeline. Exactly one milestone may be active at a time. |
| `failed` | Compile gate or acceptance criteria did not pass. Set by the agent pipeline. Human reviews. |
| `released` | All gates passed. Set by the agent pipeline. Frozen — do not modify. |

Exactly one MILESTONE may have `Status: active` at any time. This is validated
by RULE-15. The pipeline agent advances the cursor; humans do not need to edit
status fields manually.

```markdown
## MILESTONE: {version}
Status: pending

Included BEHAVIORs:
  {behavior-name-1}, {behavior-name-2}, ...

Deferred BEHAVIORs:
  {behavior-name-3}, {behavior-name-4}, ...

Acceptance criteria:
  {one criterion per line — concrete, testable, CLI-invocable where possible}
```

**Rules:**
- Every BEHAVIOR named in `Included BEHAVIORs` or `Deferred BEHAVIORs` must
  exist in the spec (validated by RULE-16).
- Together, `Included` + `Deferred` need not cover every BEHAVIOR in the spec —
  BEHAVIORs not mentioned in any milestone are always included in full
  (they have no phasing constraint).
- `Acceptance criteria` are free-form but should be concrete enough for an agent
  to evaluate: prefer CLI invocations, file existence checks, or observable outputs.
- MILESTONE sections are non-normative for pcd-lint rule purposes — they do not
  affect RULE-01 through RULE-14. Only RULE-15 and RULE-16 apply to them.
- The `## DELTA` section (for single-pass work orders) and `## MILESTONE` sections
  serve different purposes and may coexist. DELTA is ephemeral; MILESTONEs are
  persistent across translation passes.

### 2. Template Development


Deployment templates define how specifications are translated to code.

#### Template Structure
- Use `{curly_braces}` for placeholders (not `<angle_brackets>`)
- Include a DELIVERABLES section listing all files to generate
- Specify constraints and requirements clearly

#### Available Templates (v0.3.14)

| Template | Status | Default Lang | Notes |
|---|---|---|---|
| `cli-tool` | Complete | Go | Production-ready |
| `mcp-server` | Complete | Go | Production-ready |
| `cloud-native` | Complete | Go | TYPE-BINDINGS; kit findings fixed (v0.3.14) |
| `verified-library` | Stub | C | Safety/security-critical C-ABI |
| `library-c-abi` | Stub | C | General-purpose C-ABI |
| `python-tool` | Stub | Python | QM only, no formal verification |
| `project-manifest` | Stub | N/A | Multi-component projects |

#### Hints Files

Library-specific API shapes, version selection rules, and known gotchas
live in `hints/` files, separate from templates and specs.
Naming convention: `<template>.<language>.<library>.hints.md`

Current shipped hints:
- `hints/cloud-native.go.go-libvirt.hints.md`
- `hints/cloud-native.go.golang-crypto-ssh.hints.md`

Specs reference hints files via their `## DEPENDENCIES` section.
Hints are advisory only — they cannot override spec invariants.

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

## Licensing

The project uses a dual-license model:

| Artifact Type | License |
|---|---|
| Whitepaper, specs, templates, examples | CC-BY-4.0 |
| Tools (pcd-lint, etc.) | GPL-2.0-only |

**Rationale for GPL-2.0-only on tools:** Mirrors the Linux kernel model. Encourages everybody to contribute changes back to the validator toolchain.

## Submission Guidelines

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes following the guidelines above
4. Test your changes with `pcd-lint` if applicable
5. Commit with clear, descriptive messages
6. Push to your fork and submit a pull request

### Commit Message Format

Use clear, descriptive commit messages:
```
component: Brief description of change

Longer explanation of what changed and why, if necessary.
Reference any relevant issues.
```

### Testing

- Run `pcd-lint` on any specifications you modify
- For Go code: ensure `go build ./...` succeeds
- Test templates with example specifications

## Development Priorities (v0.3.14+)

### v0.3.14 completed
- cloud-native template: INDEPENDENT_TESTS Go naming note, operator.yaml dedup,
  HEALTHCHECK contradiction fixed, CRD scope note, go.sum as generated file
- Compiler gate (Phase 7) added to translator prompt
- Two hints files shipped: cloud-native.go.go-libvirt, cloud-native.go.golang-crypto-ssh
- remote-kvm-operator.md revised to language-neutral v0.3.0

### Carry-forward (templates)
1. Complete `verified-library.template.md`
2. Complete `library-c-abi.template.md`
3. Complete `python-tool.template.md`
4. Complete `project-manifest.template.md` (full BEHAVIOR STEPS beyond stub)
5. Add `independent_tests/` deliverable to cli-tool and mcp-server templates

### Tooling
6. Regenerate `pcd-lint` implementation from v0.3.13 spec (RULE-08–13 new)
7. Update generic `prompts/prompt.md` (A.13) to include compile gate, TYPE-BINDINGS
   guidance, and v0.3.13 confidence table format

## Communication

- **Email**: pcd@mailbox.org
- **Issues**: Use GitHub issues for bug reports and feature requests
- **Discussions**: Use GitHub discussions for questions and ideas

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and professional in all interactions.

## Questions?

If you have questions about contributing, please:
1. Check existing issues and discussions
2. Review the whitepaper in `doc/whitepaper.md`
3. Contact us at pcd@mailbox.org

Thank you for contributing to PCD!

