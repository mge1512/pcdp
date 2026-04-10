
# Contributing to PCD

Thank you for your interest in contributing to the Post-Coding Development
(PCD) project!

This document is for people who want to **work on PCD itself** — improving
the specification format, adding deployment templates, fixing `pcd-lint`,
writing examples, or extending the tooling.

If you want to **use PCD to specify and generate your own software**, see
`doc/guide.md` — the user guide for spec authors.

---

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

---

## Getting Started

### Prerequisites

- Git
- Go (for building tools)
- Basic understanding of Markdown
- Familiarity with the PCD specification format (read `doc/guide.md` first)

### Repository Structure

```
pcd/
├── doc/
│   ├── guide.md                   ← user guide for spec authors  ← START HERE
│   ├── whitepaper.md              ← canonical whitepaper
│   └── executive-brief.md
├── templates/                     ← deployment templates
├── hints/                         ← translator hints files
├── tools/
│   ├── pcd-lint/                  ← specification validator
│   └── mcp-server-pcd/            ← MCP server for PCD toolchain
├── examples/                      ← example specifications
└── prompts/                       ← AI translator and interview prompts
```

---

## How to Contribute

### 1. Adding or Improving Deployment Templates

Deployment templates define how specifications are translated to code for a
given target environment. Templates live in `templates/`.

#### Template Structure

- Use `{curly_braces}` for placeholders (not `<angle_brackets>`)
- Include a `## DELIVERABLES` section listing all files to generate
- Include a `## EXECUTION` section with delivery phases, resume logic,
  and compile gate
- Declare `EXECUTION: none` in META for templates that produce no compiled
  output (e.g. `project-manifest`)

#### Current Templates

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

#### Template Validation

Run `pcd-lint` against any template you modify:
```bash
pcd-lint templates/your-template.template.md
```

Templates must pass pcd-lint with zero errors.

---

### 2. Adding or Improving Hints Files

Hints files contain language-specific or library-specific implementation
knowledge for translators. They live in `hints/`.

**Four-layer naming convention:**

```
hints/
  <template>.<language>.milestones.hints.md   # scaffold-first patterns;
                                               # reusable across all components
  <component>.implementation.hints.md         # component-specific, language-neutral
  <template>.<language>.<library>.hints.md    # library-specific API shapes

# Decisions hints — lives next to the spec, not in hints/:
<specname>.<language>.decisions.hints.md      # implementation decisions from
                                               # prior translation runs;
                                               # language-specific and disposable
```

The first three layers live in `hints/` and are shared across components.
The fourth layer — **decisions hints** — lives alongside the spec itself
(in the same directory as `<specname>.md`), is named after the spec and
the target language, and records implementation decisions that the translator
made but that are not captured in the spec.

**When to use decisions hints:**

A decisions hints file is created or updated by the translator as a required
deliverable of every translation run that makes architectural decisions not
specified in the spec — package layout, routing patterns, error conventions,
type-registry choices. It answers: "what did the previous translator decide
that the next translator should know about?"

**Key properties:**

- Named `<specname>.<language>.decisions.hints.md` — the language qualifier
  makes clear it is disposable when switching target languages
- Generated alongside `TRANSLATION_REPORT.md` as a translation deliverable
- Read by the translator at the start of a **guided regeneration** or
  **incremental update** run; *not* read on a clean full regeneration
- Lives next to the spec in the repository; committed to git but not
  considered part of the spec for review or certification purposes
- Removed or regenerated when switching target languages

**Relationship to the change impact assessment:**

When `assess_change_impact` recommends full regeneration and produces a
"what to preserve" list, that list should be written into
`<specname>.<language>.decisions.hints.md` before the translator runs.
The translator reads the decisions hints file and treats its contents as
normative constraints — the same way it would treat a hints file in `hints/`.

Hints files are advisory only — they cannot override spec invariants or
template constraints. Running `pcd-lint` against a hints file produces
expected structural warnings; this is correct behaviour, not a bug.

---

### 3. Improving pcd-lint

`pcd-lint` is the reference validator. Its specification lives in
`tools/pcd-lint/spec/pcd-lint.md`. Its implementation in
`tools/pcd-lint/code/` was generated from that spec.

**To add a new validation rule:**
1. Add the rule definition to `tools/pcd-lint/spec/pcd-lint.md`
   following the `### RULE-N:` pattern
2. Add EXAMPLES for the new rule (both valid and invalid cases)
3. Add the rule to the `STEPS:` list in `BEHAVIOR: lint-validation-rules`
4. Regenerate the implementation using the translator prompt
5. Verify the new implementation passes all existing examples

**Do not hand-edit the generated implementation** — fix the spec and regenerate.

---

### 4. Improving mcp-server-pcd

The MCP server spec lives in `tools/mcp-server-pcd/spec/mcp-server-pcd.md`.
Same rule: fix the spec, regenerate the implementation.

---

### 5. Adding Examples

Example specifications live in `examples/`. Every example must:
- Pass `pcd-lint` with zero errors
- Use `Spec-Schema: {current version}` in META
- Demonstrate at least one non-trivial BEHAVIOR with error exits and
  a corresponding negative-path EXAMPLE

---

### 6. Improving Prompts

Translator and interview prompts live in `prompts/`. Changes here affect
every user of PCD — review carefully and test with at least one real
translation run before submitting.

---

## Spec Schema Rules

The canonical schema rules for PCD specifications are documented in
`tools/pcd-lint/spec/pcd-lint.md` (RULE-01 through RULE-17) and summarised
in `doc/guide.md`. When adding schema rules, update both.

**Current required sections (RULE-01):**
META, TYPES, BEHAVIOR, PRECONDITIONS, POSTCONDITIONS, INVARIANTS, EXAMPLES

**Current optional sections:**
INTERFACES, DEPENDENCIES, TOOLCHAIN-CONSTRAINTS, DELIVERABLES, MILESTONE, DELTA

---

## Validation Rules Summary

| Rule | Description | Since |
|---|---|---|
| RULE-01 | Required sections present | v0.3.0 |
| RULE-02 | META fields present and non-empty | v0.3.0 |
| RULE-03 | Deployment template resolves | v0.3.0 |
| RULE-04 | Deprecated META fields | v0.3.0 |
| RULE-05 | Verification field value | v0.3.0 |
| RULE-06 | EXAMPLES structure (GIVEN/WHEN/THEN; multi-pass) | v0.3.0 |
| RULE-07 | EXAMPLES minimum content | v0.3.0 |
| RULE-08 | BEHAVIOR blocks contain STEPS | v0.3.12 |
| RULE-09 | INVARIANTS entries carry [observable]/[implementation] tag | v0.3.12 |
| RULE-10 | Negative-path EXAMPLE required for BEHAVIOR with error exits | v0.3.13 |
| RULE-11 | TOOLCHAIN-CONSTRAINTS section structure | v0.3.13 |
| RULE-12 | Cross-section consistency | v0.3.13 |
| RULE-13 | Constraint: field value on BEHAVIOR headers | v0.3.13 |
| RULE-14 | EXECUTION section required in deployment templates | v0.3.16 |
| RULE-15 | MILESTONE section structure and single-active constraint | v0.3.21 |
| RULE-16 | MILESTONE BEHAVIOR names exist in spec | v0.3.21 |
| RULE-17 | Scaffold milestone ordering and uniqueness | v0.3.21 |

---

## CLI Conventions

All PCD tools follow these conventions:
- Key=value syntax: `pcd-lint strict=true spec.md`
- Bare words for commands: `pcd-lint list-templates`
- **NO `--flag` style ever** (firm decision)
- stderr for diagnostics; stdout for summaries and lists
- Exit codes: 0 = valid, 1 = errors/strict warnings, 2 = invocation error

---

## Licensing

| Artifact Type | License |
|---|---|
| Whitepaper, specs, templates, examples, guide | CC-BY-4.0 |
| Tools (pcd-lint, mcp-server-pcd, etc.) | GPL-2.0-only |

**Rationale for GPL-2.0-only on tools:** Mirrors the Linux kernel model.
Encourages everyone to contribute changes back to the validator toolchain.

---

## Submission Guidelines

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes following the guidelines above
4. Run `pcd-lint` on any specifications you modify
5. For Go code: ensure `go build ./...` succeeds
6. Commit with clear, descriptive messages
7. Push to your fork and submit a pull request

### Commit Message Format

```
component: Brief description of change

Longer explanation of what changed and why, if necessary.
```

---

## Communication

- **Email**: pcd@mailbox.org
- **Issues**: Use GitHub issues for bug reports and feature requests
- **Discussions**: Use GitHub discussions for questions and ideas

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all
contributors. Please be respectful and professional in all interactions.

Thank you for contributing to PCD!
