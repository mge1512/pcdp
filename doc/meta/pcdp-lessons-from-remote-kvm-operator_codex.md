# PCDP Lessons from the Remote KVM Operator Experiment

This note captures what the three-way implementation comparison
(`claude`, `codex`, `synth-by-codex`) appears to teach us about PCDP
proposal design and specification writing.

It is based on the current `remote-kvm-operator` spec revision `0.2.0`,
which explicitly incorporates lessons from the experiment in its changelog
and normative sections.

## What Changed in the Spec

The updated spec did not merely add more prose. It made previously ambiguous
areas normative:

- It binds core Kubernetes-facing fields to ecosystem-native types such as
  `metav1.Duration` and `metav1.Condition`.
- It makes the API group explicit and normative.
- It defines a required internal abstraction boundary
  (`Dialer` / `Session` / `Domain`) instead of specifying only the external
  behavior.
- It turns vague lifecycle wording into stepwise reconciliation rules,
  especially for deletion and graceful shutdown.
- It specifies build and packaging constraints that turned out to matter in
  practice, such as Containerfile layer order, `COPY --chmod`, and the ban on
  `HEALTHCHECK` in OCI builds.
- It adds negative-path examples, not just success cases.

This is the right direction. The comparison showed that independent
implementations diverge first at exactly these boundaries.

## What PCDP Should Learn

### 1. Canonical platform types must be nameable and normative

If a target platform already has canonical types, a PCDP spec should be able
to bind to them directly rather than describing near-equivalent local types.

Why this matters:

- A custom `Duration` description is not equivalent to `metav1.Duration`.
- A custom `Condition` struct is not equivalent to `metav1.Condition`.
- Two implementations can both be "reasonable" and still become
  non-interoperable.

Proposal consequence:

- PCDP should encourage or require explicit type bindings for ecosystems such
  as Kubernetes, instead of leaving translators to infer them.

### 2. Some internal abstraction boundaries are part of the design

The experiment showed that implementation quality was heavily affected by
whether the transport layer was forced behind an interface seam.

This is not just an implementation preference. It affects:

- testability
- substitutability
- compile-time isolation from external client-library API shapes
- how naturally an independent test suite can be produced

Proposal consequence:

- PCDP should allow a spec to declare selected internal seams as normative
  architecture, not only public behavior.

### 3. Failure paths need the same precision as happy paths

The strongest corrections in the revised spec concern invalid input,
finalizers, shutdown timeouts, requeue behavior, and deletion sequencing.

The original problem was not absence of intent. It was absence of precise
behavior under stress or error.

Proposal consequence:

- PCDP should require explicit failure-path semantics for:
  - mutual-exclusion rules
  - timeouts
  - deletion/finalizer flows
  - retries and requeue behavior
  - unsupported-but-allowed conditions

### 4. Buildability constraints belong in the spec when they are outcome-critical

The comparison found that several seemingly "implementation" details were
actually necessary for a valid deliverable:

- exact Containerfile structure
- OCI-specific constraints
- dependency version verification
- generated versus hand-authored files

Proposal consequence:

- PCDP should have a first-class place for toolchain and packaging
  constraints, including forbidden patterns and verification procedures.

### 5. Examples should include invalid and edge cases

The revised spec improved materially once it added examples for invalid spec
input and forced shutdown behavior.

Proposal consequence:

- PCDP should treat negative examples as a normal requirement for any feature
  with validation, asynchronous convergence, or cleanup semantics.

### 6. Specs need consistency checking across sections, not just syntax checking

The revised spec is better, but it still demonstrates why consistency linting
matters:

- the transport interface names `Dial(...)`, while some behavior text still
  says `transport.Connect(...)`
- `/readyz` is described as `healthz.Ping` and also as "not ready until leader
  elected", which is not the same semantic claim

These are not editorial nits. They can mislead translators and reviewers.

Proposal consequence:

- PCDP should validate cross-section consistency for:
  - identifier names
  - type names
  - required file names
  - endpoint semantics
  - state-machine terminology

### 7. Specs should distinguish required, supported, and optional more sharply

One useful outcome of the revision is the clearer separation between:

- required behavior
- supported-but-not-required behavior
- explicitly forbidden behavior

This helped with items like the uniqueness webhook and deployment patterns.

Proposal consequence:

- PCDP should make these categories first-class and easy to lint.

## Suggested Direction for the PCDP Proposal

If this experiment is representative, the next useful step for PCDP is not
simply "more detailed specs". It is a more disciplined spec model.

The proposal should emphasize:

1. Normative binding to platform-native types and identifiers.
2. Normative internal seams where architecture depends on them.
3. Required error-path and lifecycle semantics.
4. First-class build/toolchain/deployment constraints.
5. Required negative examples for non-trivial state machines.
6. Cross-section consistency linting.
7. Explicit classification of `required`, `supported`, `optional`, and
   `forbidden`.

## Bottom Line

The main lesson is that multi-agent divergence was not caused by weak coding.
It was caused by under-specified decision points in the spec.

PCDP should therefore optimize less for descriptive completeness alone and
more for interoperability-critical precision.
