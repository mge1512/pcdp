<!--
  Update chunk, 2026-09. Standalone slides covering what changed since
  the 2026-04 introduction. Build options:
    - cat after pcd-intro-50-technical.md in the Makefile pipeline, or
    - build alone behind pcd-intro-00-header.md.
  Requires pcd-workflow-full.png (from pcd-workflow-full.pikchr).
-->

# Update: April to September 2026

## Whitepaper 0.3.22 to 0.4.6: what changed

| Area | Change |
|---|---|
| Composition | `Includes:` merges shared specs into a host spec; merged spec hash; RULE-19 to RULE-21 |
| TYPES | Tabular TYPE grammar, closed cell vocabularies; RULE-22 to RULE-25; JSON Schema emission |
| Provenance | Reports record a SHA-256 for every translation input, not the spec hash alone |
| Verification | Tests-first enforced structurally; roles renamed to translator / test-author |
| Templates | 14 deployment templates, was 9: abap-report, cockpit-module, kubectl-style-cli, library, spack-package added |
| Tooling | `pcd-slice` shipped; `libpcd` and `pcd-emit-types` specified |

\bigskip

The reproducibility unit is now a tuple: spec + resolved language + hints and template set.\
The spec hash attests intended behaviour; it does not alone determine the binary.

---

## pcd-lint: from 17 to 25 rules - and one engine

**Rules added since April:**

- RULE-18: generated output embeds the correct spec hash (`check-report=true`)
- RULE-19 to RULE-21: composition - paths resolve, no name collisions, inclusion graph acyclic
- RULE-22 to RULE-25: tabular TYPES - structure, cell vocabulary, references, lifecycle

\bigskip

**Consolidation into one engine:**

- Every rule is written once, in `tools/shared/spec/`, and composed via `Includes:` into `libpcd`
- `pcd-lint`, `pcd-emit-types` (TYPE tables to JSON Schema), and `mcp-server-pcd` are thin front ends binding to that engine
- Lint parity between CLI and MCP holds by construction, not by copy discipline

\bigskip

Status: `pcd-lint` shipped - Go, RPM and DEB from the Open Build Service.\
`libpcd` and `pcd-emit-types`: specified, translation pending.

---

## pcd-slice: bundles instead of full-spec reads

A translation reads the spec many times - per milestone, per role.\
One full pass over a 41-BEHAVIOR production spec is on the order of 70K tokens.

\bigskip

`pcd-slice` derives the semantic form once, deterministically:

- One preamble with what every bundle shares: version, intro, full TYPES, global invariants.\
  Byte-identical across bundles, so it becomes a cacheable prompt prefix.
- One bundle per BEHAVIOR: the block verbatim, its examples, its invariants, its hints, its type closure by name

\bigskip

- Every emitted file records the SHA-256 of its sources; `pcd-slice check` detects stale bundles
- The source spec is never modified: bundles are build output, regeneratable
- `skills/pcd-translate` guards a session: translation reads bundles, never the spec

\bigskip

Shipped: Go, single static binary, standard library only, packaged via OBS.\
Read-cost reduction on production specs: typically 50 to 65 percent, before caching.

---

## The workflow, updated

![](pcd-workflow-full.png){width=100%}

1. Domain expert writes the spec, or an AI interview produces it
2. `pcd-lint` validates: 25 rules, before any translation
3. `pcd-slice` derives the preamble and the per-BEHAVIOR bundles
4. The deployment template resolves language, packaging, conventions
5. Translator and test-author produce code, tests, and the audit bundle

If the output is wrong: fix the spec and regenerate. The code is never touched.
