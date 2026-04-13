
# Contributing to PCD

This document is for people who want to **work on PCD itself** — improving
the specification format, adding deployment templates, fixing `pcd-lint`,
writing examples, or extending the tooling.

If you want to **use PCD** to specify and generate your own software, see
[`doc/user-guide.md`](doc/user-guide.md).

For the reasoning behind design decisions, see
[`doc/technical-reference.md`](doc/technical-reference.md).

---

## Getting Started

Prerequisites: Git, Go (for building tools), basic Markdown familiarity.
Read `doc/user-guide.md` before contributing.

---

## What to Work On

### Deployment templates (`templates/`)

Templates define how specs are translated for a given deployment context.
Use `{curly_braces}` for placeholders. Every template must include a
`## DELIVERABLES` table, an `## EXECUTION` section (or `EXECUTION: none`
in META), and pass `pcd-lint` with zero errors.

Current templates and status are listed in `doc/technical-reference.md`
section 4.

### Hints files (`hints/`)

Hints files contain library-specific implementation knowledge. Five-layer
naming convention — see `doc/technical-reference.md` section 5 for the
full naming rules and lifecycle of each layer, including the style hints
file (`<scope>.<language>.style.hints.md`) and the decisions hints file.

### pcd-lint (`tools/pcd-lint/`)

The reference validator. Spec lives in `tools/pcd-lint/spec/pcd-lint.md`.
To add a validation rule: add the rule definition following the `### RULE-N:`
pattern, add EXAMPLES, update the STEPS list in BEHAVIOR: lint-validation-rules,
then regenerate. Do not hand-edit the generated implementation.

### mcp-server-pcd (`tools/mcp-server-pcd/`)

Spec lives in `tools/mcp-server-pcd/spec/mcp-server-pcd.md`. Same rule:
fix the spec, regenerate.

### Examples (`examples/`)

Every example must pass `pcd-lint` with zero errors, use the current
`Spec-Schema:` version in META, and include at least one non-trivial
BEHAVIOR with error exits and a corresponding negative-path EXAMPLE.

### Prompts (`prompts/`)

Changes here affect every user. Test with at least one real translation
run before submitting.

---

## CLI Conventions

All PCD tools follow these conventions without exception:

- Key=value syntax: `pcd-lint strict=true spec.md`
- Bare words for commands: `pcd-lint list-templates`
- **No `--flag` style ever**
- stderr for diagnostics; stdout for summaries
- Exit codes: 0 = valid, 1 = errors, 2 = invocation error

---

## Licensing

| Artifact | License |
|---|---|
| Specs, templates, examples, docs | CC-BY-4.0 |
| Tools (`pcd-lint`, `mcp-server-pcd`) | GPL-2.0-only |

---

## Pull Requests

1. Fork, create a feature branch
2. Run `pcd-lint` on any spec you modify
3. For Go changes: `go build ./...` must pass
4. Commit with `component: brief description` format
5. Submit pull request

**Contact:** pcd@mailbox.org
