# Contributing to PCDP

Thank you for your interest in contributing to the Post-Coding Development Paradigm (PCDP) project!

## Project Overview

PCDP is a paradigm where:
- Domain experts write specifications in structured Markdown
- AI translates specifications into verified implementations  
- Engineers never write implementation code directly
- Target language is derived from deployment templates, not declared by spec authors

**Key distinction from "vibe coding":**
- Vibe coding: humans write code + AI suggests
- PCDP: humans write specs only; AI generates all implementation code
- If generated code is wrong: fix the spec, not the code

## Getting Started

### Prerequisites

- Git
- Go (for building tools)
- Basic understanding of Markdown
- Familiarity with specification writing

### Repository Structure

```
pcdp/
├── doc/                           # Documentation
│   ├── whitepaper.md             # Canonical whitepaper
│   └── executive-brief.md
├── templates/                     # Deployment templates
│   ├── cli-tool.template.md      # Complete, production-ready
│   ├── mcp-server.template.md    # Complete, production-ready
│   └── ...                       # Various template stubs
├── tools/                         # PCDP tooling
│   └── pcdp-lint/                # Specification validator
├── examples/                      # Example specifications
└── prompts/                       # AI translator prompts
```

## How to Contribute

### 1. Specification Writing

When writing PCDP specifications:

#### Required META Section Fields

Every PCDP spec must include:

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

#### Examples Section

Must include at least one complete test case:
```markdown
## EXAMPLES

### EXAMPLE: {descriptive_name}

#### GIVEN:
[Initial conditions]

#### WHEN:
[Action performed]

#### THEN:
[Expected outcome]
```

### 2. Template Development

Deployment templates define how specifications are translated to code.

#### Template Structure
- Use `{curly_braces}` for placeholders (not `<angle_brackets>`)
- Include a DELIVERABLES section listing all files to generate
- Specify constraints and requirements clearly

#### Available Templates (v0.3.8)

| Template | Status | Default Lang | Notes |
|---|---|---|---|
| `cli-tool` | Complete | Go | Production-ready |
| `mcp-server` | Complete | Go | Production-ready |
| `verified-library` | Stub | C | Safety/security-critical C-ABI |
| `library-c-abi` | Stub | C | General-purpose C-ABI |
| `python-tool` | Stub | Python | QM only, no formal verification |
| `project-manifest` | Stub | N/A | Multi-component projects |

### 3. Tool Development

#### pcdp-lint Validation Rules

The `pcdp-lint` tool validates specifications against these rules:

- **RULE-01**: Required sections present
- **RULE-02**: META fields validation
- **RULE-03**: Deployment template resolution
- **RULE-04**: Deprecated META fields detection
- **RULE-05**: Verification field validation
- **RULE-06**: EXAMPLES structure validation
- **RULE-07**: EXAMPLES content validation

#### CLI Conventions

All PCDP tools follow these conventions:
- Key=value syntax: `pcdp-lint strict=true spec.md`
- Bare words for commands: `pcdp-lint list-templates`
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
| Tools (pcdp-lint, etc.) | GPL-2.0-only |

**Rationale for GPL-2.0-only on tools:** Mirrors the Linux kernel model. Encourages everybody to contribute changes back to the validator toolchain.

## Submission Guidelines

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes following the guidelines above
4. Test your changes with `pcdp-lint` if applicable
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

- Run `pcdp-lint` on any specifications you modify
- For Go code: ensure `go build ./...` succeeds
- Test templates with example specifications

## Development Priorities (v0.3.9)

### Must Fix
1. Fix pcdp-lint self-validation errors
2. Add compilation gate (`go build ./...` must succeed)
3. Implement atomic delivery sentinel

### Template Completion
4. Complete `verified-library.template.md`
5. Complete `library-c-abi.template.md`
6. Complete `python-tool.template.md`
7. Complete `project-manifest.template.md`

### New Specifications
8. Add SCA CI step for REUSE compliance

## Communication

- **Email**: pcdp@mailbox.org
- **Issues**: Use GitHub issues for bug reports and feature requests
- **Discussions**: Use GitHub discussions for questions and ideas

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and professional in all interactions.

## Questions?

If you have questions about contributing, please:
1. Check existing issues and discussions
2. Review the whitepaper in `doc/whitepaper.md`
3. Contact us at pcdp@mailbox.org

Thank you for contributing to PCDP!
