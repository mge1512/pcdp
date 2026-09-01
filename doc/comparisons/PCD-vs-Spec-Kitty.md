# PCD vs Spec Kitty

Version: 2026.09.01.01
Date:    2026-09-01
Author:  Matthias G. Eckermann <pcd@mailbox.org>
Status:  Draft, intended for public use
Scope:   doc/comparisons/

---

## 1. Purpose

This is the second entry in `doc/comparisons/`. It positions PCD against
Spec Kitty, an open-source workflow and governance layer for agentic
coding.

The first entry compared PCD with a project that agrees with it on the
paradigm and optimizes a different end of it. This entry does the opposite.
Spec Kitty's author read the spec-driven literature, understood the
spec-as-source-of-truth premise, and rejected it in print. His conclusion is
that the code is the software and the specification is a change request plus
a decision ledger. That is the premise PCD inverts.

A comparison set that lists only projects converging on PCD's thesis is a
selection artifact. Spec Kitty states the opposite position in print and gives
reasons for it, from inside the same movement, so answering it does more for
the paradigm than agreeing with a third project would. It also has more users
than the other entries in this directory, so a reader may already know
its argument.

Sources are listed in section 8. Claims that rest on a single article or a
single repository read are flagged in section 7 as "to verify" before any
public use.

---

## 2. What Spec Kitty is

Spec Kitty is a command-line tool that wraps a governed workflow around AI
coding agents. It turns product intent into a repeatable loop, keeps the
context of that loop in the repository, cuts implementation into work
packages that agents execute, and isolates each package in a git worktree so
parallel work does not collide.

The loop is: charter, specify, plan, tasks, next, implement, review, accept,
merge, retrospective. Each step has a slash command generated for the user's
agent, and a corresponding command in the CLI.

| Item | Value |
|---|---|
| Origin | Fork of GitHub Spec Kit, since heavily diverged |
| License | MIT |
| Language and runtime | Python, requires 3.11 or newer |
| Distribution | Python package index, installed with `pipx`, `uv` or `pip` |
| Repository artifacts | `kitty-specs/`, `.kittify/config.yaml`, `.worktrees/` |
| Agent coverage | Roughly fifteen agent front ends, including Claude Code, Codex, Cursor, Gemini CLI, Copilot, Windsurf, OpenCode, Qwen, Kiro |
| Governance | Per-project charter, doctrine layer (directives, tactics, profiles), architecture decision records, `advise` / `ask` / `do`, trail model |
| Additional surfaces | Local kanban dashboard, mission retrospectives, optional hosted workspace |
| First release | 2025-11-02 (`spec-kitty-cli` 0.2.20) |
| Current release | 3.2.5, 2026-07-08; pre-release 3.2.6rc2, 2026-08-21 |
| Release count | 165 published versions in ten months |

Three constructs account for most of the design.

The **mission** is the unit of work: a purpose-specific workflow with defined
inputs, process, outcomes and states. Coding is one mission kind;
research and documentation are others.

The **charter** is a per-project governance document that records what an
agent cannot infer from the code. The published example includes which of two
container workflows is policy rather than convenience, which continuous
integration failure is expected after a merge and must not be treated as a
defect, which domain noun is current and which is a legacy alias, and which
user-authored files an upgrade must preserve. Every project gets its own
charter rather than inheriting a generic rule set.

**Cross-review with a fresh context** is the verification step. When an
implementation lands, a second agent, or the same agent with its context
stripped, reviews it. The stated reason is that an agent which has just
finished writing code has spent the preceding hour convincing itself the code
is correct, and is therefore the worst available reviewer of it.

The project publishes a security position dated 2026-02-18. It states that
the core is a deterministic workflow tool which does not set approval,
sandbox or network flags on the agent front ends it drives; that autonomous
orchestration is a separate, optional component that policy can disallow;
that workflow state changes must go through the host commands; and that there
is no direct write path from an agent provider to the hosted service. The
document reads as a response to earlier criticism about bundled autonomous
execution.

---

## 3. The disagreement

Spec Kitty exists because of a disagreement with Spec Kit, and its author
states it plainly. On his reading, Spec Kit places the specification at the
center as the source of truth for what the software is, with code downstream
as a generated artifact faithful to that source. He rejects that. The code is
what compiles, what runs in production, and what customers experience;
changing the code changes the software whether or not any specification
agrees.

In his model the specification is a different object: a change request, a
history of how the team arrived here, a ledger of decisions with the
reasoning attached, and a guardrail for future work. It is the source of
truth for where the software is going, not for what it is. He declines to
place his tool between the agent and the codebase, on the argument that
agents read code well; the failure he set out to fix is intent translation
and lost project memory, not code comprehension.

PCD asserts the inverse. The specification is the sole durable human
artifact, generated code is a build output, and a defect is fixed in the
specification and the code regenerated.

The two positions sit at opposite ends of one axis, and the axis is which
artifact answers the question *what is this software*:

```
Spec Kitty            Spec Kit / Kiro      SPDD / Plain          PCD
code is truth,        spec-anchored,       spec-first with       spec-only,
spec is a change      code maintained      sync back to spec     code is a
request and ledger                                               build output
```

Stated as one sentence for talks: Spec Kitty makes the specification the
source of truth for where the software is going, PCD makes it the source of
truth for what the software is, and everything else in this table follows
from that one choice.

---

## 4. Where PCD differs

The differences all descend from section 3. They are listed in the order a
skeptical reader will raise them.

### 4.1 The audited object

Under Spec Kitty the merged code remains the object a reviewer, an auditor or
a certifier reads. The charter, the review lane and the merge record are
evidence about how that code was produced, and an evaluator can work with
them, but the artifact under examination is still a codebase that changes
with every mission.

Under PCD the specification is the object under examination and the binary is
its derivative. The specification changes when behavior changes and not
otherwise, which makes the reviewed object stable in a way a codebase is not.
PCD's assurance argument rests on this, and 4.2 to 4.6 are its
consequences.

### 4.2 Reproducibility and attestation

PCD pins the input tuple - specification, resolved language, hints and
template set - and records a labeled hash of every file consumed as a
translation input. The specification hash is embedded in source headers, in
the binary's version output, in package metadata and in container labels, so
the question "was this binary produced from the certified specification?" is
answered by recomputing a hash rather than by asking a person.

No comparable mechanism is documented in Spec Kitty. Nothing in the published
material records which model, which prompt version and which specification
revision produced a given work package in a form that survives the merge.
Provenance in that system is the git history plus the decision ledger, which
is a record of intent, not an attestation of derivation.

This is a difference of purpose, not a defect. A team shipping product
features does not need a derivation proof.

### 4.3 Verification: the same instinct, placed differently

Both systems refuse to let the author of an implementation be its reviewer,
and both use a second model with a clean context to enforce that. The
argument for it is the same in both projects, and Spec Kitty states it well.

The placement differs. Spec Kitty's reviewer reads a diff and decides whether
to accept it; the reviewer may propose changes to that diff, and the accepted
diff is merged. PCD's second model does not review the code at all. It
authors an independent test suite from the EXAMPLES section of the
specification, before the translator writes anything, and any finding it
raises must be resolved in the specification, the hints, the template or the
translator prompt. A reviewer that is permitted to patch code has repaired
the symptom and left the specification wrong.

### 4.4 Language neutrality

Spec Kitty operates on the repository as it already is. The target language,
frameworks and libraries are properties of the existing codebase, which is
appropriate for its target case and leaves no retargeting question to answer.

PCD keeps language-specific realization out of the behavioral specification
and places it in hints and deployment templates, so that one specification
produces implementations in several languages without specification edits.
This is what makes a specification outlive the language ecosystem it was
first translated into, and it only has meaning if the specification, rather
than the code, is the durable artifact.

### 4.5 Supply chain and sovereignty

Spec Kitty is distributed through a language package index and installed with
a language-ecosystem installer. Its runtime is an interpreter with a version
floor, its agent front ends are largely hosted services, and a hosted
workspace is available as an option. Nothing here is unusual, and the
security position document gives reasonable deployment guidance for
organizations that need to control the agent front ends.

PCD's target is a pipeline an organization can run itself: signed packages
from a controlled build service, offline and reproducible builds, an SPDX
software bill of materials, and models that can be hosted in a chosen
jurisdiction or air-gapped entirely. For a customer whose constraint is
sovereignty rather than velocity, this is where the two diverge furthest in
practice, and it is a distribution and deployment question rather than a
paradigm question.

### 4.6 Certification: the same goal, a different distribution of evidence

Neither approach makes certification impossible, and the comparison is easy
to overstate in PCD's favor. Common Criteria evaluations of hand-written code
are routine and have been for decades. At EAL4 the developer supplies the
design documentation, the development-environment controls and the
implementation representation of the security functionality, and the
evaluator examines a sample of it. Introducing an AI coding agent into that
process does not close the path. It changes what the developer has to show.

Under Spec Kitty's model the code stays the evaluated object and the agent
becomes part of the development environment. The team must show that its
review step, its work-package boundaries and its acceptance gates are
adequate controls over what the agent produced. That arrangement is familiar to evaluators,
and the charter, the review lane and the merge record are usable evidence. The
cost is that the evaluated object keeps moving: every agent-authored change
re-enters code-level review scope, and the evidence that a given change was
understood is a human sign-off rather than a machine-checkable link.

Under PCD the primary human-reviewed object moves to the specification and
the specification-to-artifact link is made machine-checkable. The claim is
narrower than it is sometimes read to be. PCD does not assert that a
certifiable system must be built this way; it redistributes the evidence,
trading code-level review effort for provenance and pipeline evidence, and
betting that the trade is accepted.

The residue is stated honestly in the whitepaper and repeated here: PCD's
route depends on a certifier accepting specification-level comprehension plus
an auditable pipeline in place of code-level comprehension, and that has not
yet been demonstrated at certification scale. Spec Kitty's route is the
proven one and is expensive per change. PCD's route is cheaper per change if
it is accepted, and acceptance is exactly what is still open.

---

## 5. Where Spec Kitty leads

Stated plainly, so the comparison stays honest.

- Brownfield work. Spec Kitty is built for a codebase with twenty years of
  decisions in it, most of them undocumented. PCD assumes a component whose
  behavior can be written down in full, which is a much narrower starting
  condition. Where a specification cannot be written, Spec Kitty has an
  answer and PCD does not.
- The charter as a named artifact. Project policy that the agent reads, kept
  separate from the specification and from the framework's own rules, is a
  construct PCD lacks under any name. PCD has hints, templates, prompts and
  contribution rules, which cover parts of the same ground without being
  recognized as one thing.
- Parallel execution. Work packages with lifecycle lanes, worktree isolation
  per package and a dashboard across all of them are an orchestration layer
  PCD has never built. The per-behavior slicing PCD does for token cost is
  the same decomposition applied to a different problem, and the two could
  meet.
- Multi-agent reach. Command surfaces generated for roughly fifteen agent
  front ends from one source is the pattern any framework needs if it wants
  adoption beyond its author's own toolchain.
- Interview-first authoring. Spec Kitty interviews the human until the
  boundaries and contracts of the intended change are settled, and presses
  when the answers are vague or contradictory. This is the second independent
  arrival at conversational specification authoring recorded in
  `doc/comparisons/`, which moves it from an idea to a pattern.
- Retrospectives as artifacts. Every completed mission produces one by
  default, and there is machinery to synthesize across them. PCD produces
  findings and change briefs per translation run but has no loop that feeds
  outcomes back into the framework.
- Adoption. The project is roughly ten months old with 165 published
  releases, an active issue tracker, third-party training offerings and
  visible community use. A single-author framework cannot manufacture that,
  and the gap belongs in the comparison rather than outside it.

---

## 6. Net positioning

PCD and Spec Kitty disagree about which artifact is authoritative, and they
are right about different situations.

- Spec Kitty optimizes for a large existing codebase whose behavior nobody
  can write down, where the achievable win is that agents stop misreading
  intent and the team stops losing its decisions. It keeps human
  accountability at the code, which is the conservative and currently
  provable position.
- PCD optimizes for components whose behavior can be specified completely and
  must stay verifiable for a decade, where the achievable win is that the
  binary remains a checkable function of a certified specification produced
  by a pipeline the operator can run.

One-line framing for talks: both projects agree that intent has to be written
down; they disagree about whether writing it down replaces the code or
governs the people who change it, and the answer depends on whether the
component can be specified in full.

The related question is where the boundary lies. A system usually contains
both kinds of component: a long-lived, specifiable core and a large body of
integration code that nobody will ever specify completely. Nothing prevents
the two approaches from meeting at that boundary, and the open design
question is what the interface between them looks like.

---

## 7. Open questions and claims to verify

Before section 4 is used publicly, confirm the following against the
documentation and a working trial rather than against the article and the
repository page alone.

- Does any part of the pipeline record model identity, prompt version and
  specification revision per work package in a form that survives the merge?
  Absence in the documentation is not proof of absence, and this single
  answer decides 4.2.
- Is the charter machine-checkable, or is it text placed in an agent's
  context and honored at the model's discretion? The distinction decides how
  much of 4.6 stands.
- For a greenfield component, does the specification ever grow into a
  description of the component in full, or does it stay a change request by
  construction? This is where the two paradigms would touch.
- Are there users under regulated audit, and what do they present in place of
  code-level review evidence?
- Is the hosted workspace available self-hosted, with a choice of model and
  jurisdiction? This affects 4.5 and nothing else.
- Confirm the fork lineage from GitHub Spec Kit from a primary source. Three
  secondary sources agree, which is suggestive rather than conclusive.
- Confirm the release and version figures at the time of use. They were read
  on 2026-09-01 and this project changes weekly.

---

## 8. Sources

- Robert Douglass, "Why I Built Spec Kitty", Medium, 2026-04-27. The
  statement of the disagreement in section 3 comes from here.
- github.com/Priivacy-ai/spec-kitty - README, `SECURITY-POSITION.md`
  (dated 2026-02-18), repository structure. Read 2026-09-01.
- docs.spec-kitty.ai - documentation index and glossary, version 3.2.
  Read 2026-09-01.
- Python package index metadata for `spec-kitty-cli`, retrieved 2026-09-01:
  license, Python floor, release list and upload dates.
- Public profile material for the project's founders, used only for section 2
  attribution.

---

## 9. Follow-up actions (non-normative)

- Decide whether PCD wants a named per-project governance artifact in the
  sense of section 5, distinct from hints and templates, or whether the
  existing layers already cover it under different names.
- Complete section 7 verification before citing 4.x publicly.
- Add the axis sentence from section 3 to the conference talk. It positions
  PCD without attacking a project whose author argued his case in public and
  argued it well.
- Consider whether the boundary question at the end of section 6 deserves its
  own note: which components in a mixed system are specifiable, and what the
  interface between a specified core and an unspecified periphery looks like.

---

## Changelog

- 2026.09.01.01 - Initial draft. Based on the founder's article of 2026-04-27,
  the repository README and security position document, the documentation
  index, and package index metadata, all read on 2026-09-01. Written for
  public use: no third-party or internal project names appear. Section 4.6
  rewritten from an earlier draft that claimed certification was closed to
  hand-maintained AI-generated code; the accurate statement is that both
  routes are open and distribute the evidence differently.
