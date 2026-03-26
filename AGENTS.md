# AI Agent Guidelines

This project accepts contributions from AI agents and AI-assisted workflows.
This document defines expectations, conventions, and architecture context.

## Principles

1. **Don't break what works.** Minimise changes. If a test suite exists,
   run it before submitting. If the compile gate fails, fix only the
   identified errors — do not rewrite unaffected files.

2. **Be honest about limitations.** "Applies and builds" is not
   "semantically correct." Do not overstate what a tool verifies.
   In `TRANSLATION_REPORT.md`, use the defined confidence levels honestly:
   High requires a named passing test; Low means reasoning only.
   Unverified claims must be listed explicitly — never silently omitted.

3. **Supply chain security is not optional.** This project targets
   regulated and safety-critical deployment contexts. Concretely:
   - Never use `curl` to download software or dependencies
   - Never use unqualified container image names (e.g. `golang:1.24`)
     Always use fully-qualified names with registry
     (e.g. `registry.suse.com/bci/golang:latest`)
   - Never fabricate dependency version strings or commit hashes
   - Use OBS (build.opensuse.org) for packaging; signed packages only
   - Formal certification frameworks (Common Criteria, ISO 26262, etc.)
     are reference points for the level of rigour expected — not
     necessarily required for every component

4. **No internal references.** No hostnames, IPs, internal paths,
   personal names, or employer-specific content in committed material.
   Use project-neutral framing. Product names are permitted where
   technically precise (e.g. base container image registries).

5. **Credit your work.** Use `Co-Authored-By: <Model> <contact>` in
   commits and documents where AI made a substantive contribution.

6. **Fix the spec, not the code.** This is the core PCDP invariant.
   If generated code is wrong, update the specification and regenerate.
   Never hand-edit generated implementation files.

---

## What This Project Is

The **Post-Coding Development Paradigm (PCDP)** is an open specification
for a software development paradigm where:

- Domain experts write specifications in structured Markdown
- AI translates specifications into verified implementations
- Engineers never write implementation code directly
- The target language is derived from deployment templates — never
  declared by the spec author
- AI translation is probabilistic; correctness comes from multiple
  complementary mechanisms (human-reviewable specs, formal verification
  when used, independent tests, audit trails) — not from spec structure alone

**This is not vibe coding.** In vibe coding, if the output is wrong, you
edit the code. In PCDP, you fix the specification and regenerate. The spec
is always the source of truth.

**Key artifacts:**
- `pcdp-lint` — validates PCDP specification files (RULE-01 through RULE-14)
- `mcp-server-pcdp` — MCP server serving templates, prompts, and hints;
  exposes `lint_content` and `lint_file` tools for in-session validation
- Deployment templates — define target language, packaging, conventions,
  and the full translation execution recipe per deployment type

**Licenses:**
- Whitepaper, specs, templates, examples: CC-BY-4.0
- Tools (`pcdp-lint`, `mcp-server-pcdp`): GPL-2.0-only

---

## What You Can Do

- Fix bugs in `pcdp-lint` rule implementations
- Add or improve deployment templates
- Add or improve library hints files
- Improve `mcp-server-pcdp` tool and resource implementations
- Update documentation — whitepaper, README, slide content
- Add EXAMPLES to existing specs
- Translate a spec to code using the standard translator prompt

## What Requires Human Review

- Changes to pcdp-lint RULE definitions — rules affect all downstream
  translation runs and must be reviewed for correctness and consistency
- New deployment templates — the EXECUTION section governs how AI
  translators behave; errors here affect every translation for that type
- Changes to the spec schema (required sections, field names, constraints)
- Any claim about model accuracy, security properties, or certification
  readiness
- Publication decisions — the repository is currently private

---

## Architecture Quick Reference

### Spec format

Required sections: `META`, `TYPES`, `BEHAVIOR`, `PRECONDITIONS`,
`POSTCONDITIONS`, `INVARIANTS`, `EXAMPLES`

Optional sections: `INTERFACES`, `DEPENDENCIES`, `TOOLCHAIN-CONSTRAINTS`,
`DELIVERABLES`

Every `BEHAVIOR` block requires:
- `STEPS:` — ordered algorithm with explicit error exits on each step
- `Constraint:` — `required` (default) | `supported` | `forbidden`
- Optional `MECHANISM:` annotation where the *how* is normative

Every `INVARIANTS` entry should carry `[observable]` or `[implementation]`.

### Two-layer prompt architecture

```
prompts/prompt.md              — universal, language-agnostic principles
                                 delegates execution recipe to the template

templates/<n>.template.md      — deployment-specific ## EXECUTION section:
  ## EXECUTION                   input files, delivery phases (ordered),
    ### Input files              resume logic, compile gate
    ### Delivery phases
    ### Resume logic
    ### Compile gate
```

`pcdp-lint` RULE-14 validates that every deployment template has
a `## EXECUTION` section with the required subsections.

### Deployment templates (current)

| Template | Status | Default language |
|---|---|---|
| `cli-tool` | Complete | Go |
| `mcp-server` | Complete | Go |
| `cloud-native` | Complete | Go |
| `verified-library` | Stub | C |
| `library-c-abi` | Stub | C |
| `python-tool` | Stub | Python |
| `project-manifest` | Stub | — |

### mcp-server-pcdp

7 tools: `list_templates`, `get_template`, `list_resources`,
`read_resource`, `lint_content`, `lint_file`, `get_schema_version`

Native MCP resources (browseable without tool calls):
- `pcdp://templates/{name}` — full template Markdown
- `pcdp://prompts/interview` — interview prompt (embedded at build time)
- `pcdp://prompts/translator` — translator prompt (embedded at build time)
- `pcdp://hints/{template}.{lang}.{lib}` — library hints

Transports — same binary, bare-word selection:
```bash
mcp-server-pcdp stdio   # for mcphost, Claude Desktop, VS Code
mcp-server-pcdp http    # default: 127.0.0.1:8080
```

### Repository layout

```
prompts/          — translator prompt, interview prompt, usage guides
templates/        — deployment templates (*.template.md)
hints/            — library hints files (*.hints.md — not PCDP specs)
tools/
  pcdp-lint/
    spec/         — canonical pcdp-lint specification
    code/         — generated Go implementation
  mcp-server-pcdp/
    spec/         — canonical mcp-server-pcdp specification
    code/         — generated Go implementation
doc/              — whitepaper, executive brief, presentation slides
examples/         — example PCDP specs
```

---

## Conventions

### CLI style (all pcdp tools)
- `key=value` for options, bare words for commands
- No `--flag` style. Ever.
- `stderr` for diagnostics, `stdout` for summaries
- Exit codes: `0` = valid/success, `1` = errors, `2` = invocation error

### Containerfiles
- Builder stage: `FROM registry.suse.com/bci/golang:latest`
- Final stage: `FROM scratch` (static binary, no runtime deps)
- Never use unqualified names (`golang:1.24`, `docker.io/golang`)
- Layer order: `COPY go.mod go.sum` → `RUN go mod download` → `COPY . .`

### Go modules
- Declare direct dependencies only in `go.mod`
- Never hand-write indirect dependencies — use `go mod tidy`
- Never fabricate pseudo-versions or commit hashes for untagged modules
- Use verified versions from hints files when available

### Diagrams
- README.md: Mermaid (GitHub native rendering)
- Whitepaper and audit bundles: Pikchr
- Slides: Pikchr (converted via `pikchr --svg-only | sed ... | magick`)

### Hints files
- Named: `<template>.<language>.<library>.hints.md`
- Live in `hints/` — they are **not** PCDP specs
- Running `pcdp-lint` against a hints file produces expected errors;
  this is correct behaviour, not a bug
- Advisory only — cannot override spec invariants or template constraints

---

## Tooling Notes

- **pcdp-lint:** `pcdp-lint myspec.md` / `pcdp-lint strict=true myspec.md` /
  `pcdp-lint list-templates`
- **mcp-go:** v0.46.0. Use `NewStreamableHTTPServer` (not `NewSSEServer`).
  Use `mcp.NewToolResultError` for domain errors, not Go `error` returns.
- **pikchr:** system dependency; install via OBS. Font fix required:
  `sed 's/<text /<text font-family="DejaVu Sans, Liberation Sans, sans-serif" /g'`
- **Slides:** pandoc → pdflatex. Use `\textrightarrow{}` not `$\rightarrow$`
  in list contexts. UTF-8 em-dashes require `---`. Consider XeLaTeX for
  native UTF-8 support.
- **max_tokens:** ≥ 32000 for complete translation runs
- **Filesystem MCP:** must allow subdirectory creation for packaging artifacts
