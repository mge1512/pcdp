
# PCDP: Lessons Learned from the remote-kvm-operator Exercise

**Source material:**
Three independent implementations of the same spec (`remote-kvm-operator.md` v0.1.0)
produced by different AI translators (claude, codex, synth-by-codex), plus the
iterative build error log from the claude implementation.

**Purpose:**
Feed concrete, evidence-based improvements back into the PCDP proposal and its
specification schema. Every finding below is traceable to an observed failure or
divergence — nothing is speculative.

**Status of open items:**
Items marked `OPEN` require a PCDP design decision before they can be resolved.
Items marked `RESOLVED` include a concrete proposed change.

---

## Finding 1 — Logical types diverge at the language binding step

**Observed failure:**
The spec defined `Duration := string where matches "^[0-9]+(s|m|h)$"` and
`Condition := { type, status, lastTransitionTime, reason, message }`.
Claude produced `string` fields. Codex produced `metav1.Duration` and `metav1.Condition`.
These are not interoperable at the Kubernetes API level: a resource created by one
implementation cannot be reliably read or validated by the other.

**Root cause:**
PCDP TYPES define the logical data model. They correctly say nothing about concrete
language bindings — the spec author does not know which language will be used.
But there is no mechanism in the framework for the resolved language to supply
canonical type mappings. The translator fills the gap with discretion, and
discretion produces divergence.

**What this is NOT:**
This is not a problem with the spec. A spec that says `Duration` is correct.
It must not say `metav1.Duration` — that would couple the spec to Go.

**RESOLVED — belongs in the template, not the spec:**
The deployment template is the point where `LANGUAGE=Go` is resolved. It is the
right place to supply a `TYPE-BINDINGS` table keyed by language:

```
## TYPE-BINDINGS

| Spec Type  | LANGUAGE=Go                                   | Notes                          |
|------------|-----------------------------------------------|--------------------------------|
| Duration   | metav1.Duration (k8s.io/apimachinery/meta/v1) | Serialises as "60s", "5m" etc. |
| Timestamp  | metav1.Time                                   |                                |
| Condition  | metav1.Condition                              | Use standard k8s condition type|
| List<T>    | []T                                           |                                |
```

The translator reads this table after resolving LANGUAGE and applies it mechanically
to every TYPES occurrence. The spec author is never involved.

**Proposed PCDP schema change:**
Add a `## TYPE-BINDINGS` section to deployment templates. Format: a table mapping
logical spec types to concrete language types, keyed by LANGUAGE value.
This section is read by the translator during the resolution step, alongside the
TEMPLATE-TABLE. A binding row for a logical type *overrides translator discretion*
for that type when the matching LANGUAGE is active.

---

## Finding 2 — BEHAVIOR sections specify contracts but not algorithms

**Observed failure:**
The `reconcile-remote-vm` spec said:
> "Reconciler waits up to 120s for domain state to become Shutoff."

Claude interpreted this as a blocking `time.Sleep` poll loop inside the reconciler
goroutine. Codex interpreted it as: send shutdown, immediately requeue, check state
on re-entry. Both satisfy the POSTCONDITIONS. The blocking loop is architecturally
wrong for controller-runtime (it starves the work queue), but the spec gave no
way to detect the difference.

**Root cause:**
PCDP BEHAVIOR blocks have PRECONDITIONS and POSTCONDITIONS but no ordered STEPS.
Pre/postconditions describe a *contract*. They do not describe an *algorithm*.
Two algorithms can satisfy the same contract while having entirely different
runtime properties.

**RESOLVED — add an ordered STEPS list as a required BEHAVIOR element:**
STEPS are imperative, numbered, and include explicit error exits. They sit
between PRECONDITIONS and POSTCONDITIONS and are normative.

Example:
```
STEPS:
1. If DeletionTimestamp set → goto DELETION-PATH.
2. Validate spec invariants (e.g. domainXML XOR domainXMLRef).
   On failure → set phase=Error, return RequeueAfter(60s).
3. Ensure finalizer present. If added → requeue immediately.
4. Load host; if not Ready → set phase=Pending, return RequeueAfter(60s). No libvirt call.
5. transport.Connect(); on TransportError → classify, update status, requeue.
6. Dispatch on desiredState.
```

A `MECHANISM:` annotation may accompany individual steps where the *how* matters
for correctness, not just the *what*:

```
STEPS:
...
5. Call domain.Shutdown().
   MECHANISM: do NOT block the reconciler goroutine; requeue immediately.
   Store shutdown-start time in annotation for timeout detection on re-entry.
```

**Proposed PCDP schema change:**
Make `STEPS:` a required element of every `BEHAVIOR:` block.
Add optional `MECHANISM:` inline annotation for steps where the implementation
pattern matters and cannot be inferred from the postconditions alone.
Pre/postconditions remain as the verifiable contract. STEPS are the algorithm.

---

## Finding 3 — Single-pass EXAMPLES cannot express multi-step behaviours

**Observed failure:**
The `vm-graceful-stop` EXAMPLE had a single WHEN/THEN pair showing the final outcome.
The intermediate state — what the reconciler does on the first pass, how it communicates
shutdown-start time to the next pass — was invisible. Both implementations passed the
EXAMPLE while using structurally different mechanisms.

**Root cause:**
PCDP EXAMPLES use the pattern `GIVEN → WHEN → THEN`. This is sufficient for
single-pass behaviours. It is insufficient for multi-pass reconciliation, where
the observable state changes across multiple reconciler invocations, and the
intermediate states are normative.

**RESOLVED — allow multi-pass WHEN/THEN sequences in EXAMPLES:**
An EXAMPLE may contain multiple WHEN/THEN pairs, each representing one reconciler
invocation. This makes the intermediate states explicit and testable:

```
EXAMPLE: vm-graceful-stop
GIVEN:
  RemoteVirtualMachine "testvm-01", spec.desiredState = Stopped, shutdownTimeout = "120s"
  Domain is Running on the remote host

WHEN:  reconcile-remote-vm runs (pass 1)
THEN:
  domain.Shutdown() is called
  A shutdown-requested timestamp is recorded (annotation or condition)
  result = RequeueAfter(short interval)

WHEN:  reconcile-remote-vm runs (pass 2); domain is now Shutoff
THEN:
  status.phase = Stopped
  result = RequeueAfter("60s")

WHEN:  reconcile-remote-vm runs (pass 2 alternate); shutdownTimeout has expired
THEN:
  domain.Destroy() is called
  status.phase = Stopped
  status.conditions includes {type: ForcePowerOff, reason: ShutdownTimeout}
  result = RequeueAfter("60s")
```

**Proposed PCDP schema change:**
Allow EXAMPLES to contain repeated `WHEN:/THEN:` pairs.
A single-pass EXAMPLE remains valid as a special case (one WHEN/THEN pair).
Multi-pass examples are identified by two or more WHEN/THEN pairs.

---

## Finding 4 — Interface types are a first-class design concept with no spec home

**Observed failure:**
The spec defined data types (`RemoteKVMHostSpec`, etc.) in TYPES. It did not define
the transport abstraction boundary (`Dialer`, `Session`, `Domain` interfaces).
Codex invented these interfaces independently; Claude bypassed them entirely and
coupled the reconciler directly to the go-libvirt library. The codex design was
strictly better: it enabled unit testing, clean error handling, and isolated
library API instability behind a boundary. But neither approach was
*required* by the spec.

**Root cause:**
PCDP TYPES was designed for data types. Go interfaces — which define *behavioural
contracts between internal components* — are structurally different. They are
module boundary requirements. A spec author who wants to mandate a clean
abstraction layer has no place to declare it.

**RESOLVED — add an INTERFACES section for behavioural contracts:**
INTERFACES declares the module boundary contracts that translators must implement.
Each interface entry specifies:
- The required method signatures (language-agnostic pseudo-signatures)
- Which implementations must be produced (production, test double, stub)
- What the test double's state machine must do

```
## INTERFACES

Dialer {
  required-methods:
    Dial(ctx, ConnectionSpec) → (Session, TransportError?)
  implementations-required:
    production:   RealDialer  (uses actual SSH + library transport)
    stub:         StubDialer  (always errors; for compilation verification only)
  test-double:    FakeDialer  (configurable per-test)
}

Session {
  required-methods:
    Ping(ctx) → error
    ProbeHost(ctx) → (HostProbe, error)
    LookupDomain(ctx, name) → (Domain, DomainNotFoundError?)
    DefineDomainXML(ctx, xml) → (Domain, error)
    Close() → error
  test-double:    FakeSession {
    configurable fields: PingErr, ProbeResult, ProbeErr, domainMap
    LookupDomain: returns FakeDomain if name in domainMap, else DomainNotFoundError
  }
}

Domain {
  required-methods:
    Name() → string
    GetInfo(ctx) → (DomainInfo, error)
    Create(ctx) → error
    Shutdown(ctx) → error
    Destroy(ctx) → error
    Suspend(ctx) → error
    Undefine(ctx) → error
  test-double:    FakeDomain {
    state machine:
      Create()   → state transitions to Running if no CreateErr
      Shutdown() → state transitions to Shutoff if no ShutdownErr
      Suspend()  → state transitions to Paused  if no SuspendErr
      Destroy()  → state transitions to Shutoff if no DestroyErr
  }
}
```

The INTERFACES section is language-agnostic. The method signatures use the spec's
own type vocabulary. The translator maps them to the resolved language.

**Consequence for INDEPENDENT_TESTS:**
Once INTERFACES are normative, the rule becomes: independent tests must use
*only* the test doubles declared in INTERFACES. They must not import or depend
on the production implementation. This makes `go test ./independent_tests/`
runnable without a live cluster, a live libvirt daemon, or any external service.

**Proposed PCDP schema change:**
Add `## INTERFACES` as a new optional top-level section in the spec schema,
between TYPES and BEHAVIOR. An INTERFACES block declares a named interface,
its required methods, and the implementations that must be produced.

---

## Finding 5 — External library API shapes belong neither in the spec nor in the template

**Observed failure:**
Nine compiler errors in claude's implementation were caused by incorrect
assumptions about go-libvirt function signatures:
`DomainGetInfo` returns 6 individual values, not a struct.
`ConnectGetMaxVcpus` requires `libvirt.OptString`, not a plain `string`.
`NodeGetInfo` returns 9 values.
`ConnectGetLibVersion` returns `uint64`, not `uint32`.
None of this was in the spec or template.

**What this is NOT:**
These are not specification concerns (the spec correctly says "call virDomainGetInfo")
and not template concerns (the template correctly says "use Go"). They are
*library-specific implementation gotchas* that are only observable once you
actually try to compile against that specific library version.

**OPEN question — where does this live in the PCDP world?**
Two candidate locations were considered:

**Candidate A: `EXTERNAL-API:` annotations in BEHAVIOR/INTERNAL blocks.**
Rejected: this would require the spec author to know the library and its API
shapes. Spec authors should not be required to know go-libvirt internals.

**Candidate B: Library hints files in the preset hierarchy.**
A new artefact type: `<template>.<language>.<library>.hints.md`.
Lives at `/usr/share/pcdp/hints/cloud-native.go.go-libvirt.hints.md`.
Read by the translator during code generation, after template resolution.
Contains concrete API gotchas, version-specific notes, and verified
pseudo-version strings for untagged modules.
Purely advisory — cannot override spec invariants.

Example hints file content:
```markdown
# cloud-native · Go · github.com/digitalocean/go-libvirt

## Version selection
This module has no tagged releases. Use a pseudo-version.
Verified good versions (as of 2024):
  v0.0.0-20220804181439-8648fbde413e  (used by containers/podman, lima-vm/lima)
DO NOT fabricate commit hashes or timestamps. Verification:
  git ls-remote https://github.com/digitalocean/go-libvirt.git HEAD

## API shapes that differ from the libvirt C API naming
- DomainGetInfo returns 6 individual values: (state uint8, maxMem uint64,
  memory uint64, nrVirtCPU uint16, cpuTime uint64, err)
  NOT a struct, despite DomainGetInfoRet existing as a type.
- ConnectGetMaxVcpus requires libvirt.OptString{"kvm"}, not a plain string.
- NodeGetInfo returns 9 individual values (model, memory, cpus, mhz,
  nodes, sockets, cores, threads, err).
- NodeGetFreeMemory() returns (uint64, error) — single byte count.
- ConnectGetLibVersion and ConnectGetVersion return uint64, not uint32.
  Format as: major = v/1000000, minor = (v%1000000)/1000, patch = v%1000.
```

Candidate B is the correct home. It is consistent with the preset-hierarchy model
already in PCDP (`/usr/share/pcdp/templates/`, `/etc/pcdp/presets/` etc.).
It separates library-specific knowledge from both the spec and the template.
It is maintainable independently as libraries evolve.

**Proposed PCDP framework change:**
Introduce a `hints/` directory in the preset hierarchy:
  `/usr/share/pcdp/hints/<template>.<language>.<library>.hints.md`

The translator reads all matching hints files after template resolution,
before generating code. Hints files are:
- Advisory only (cannot override spec invariants or template constraints)
- Version-tagged in their META section
- Maintainable without touching specs or templates

---

## Finding 6 — Dependency provenance is a supply-chain concern that needs a spec hook

**Observed failure:**
Claude invented a go-libvirt pseudo-version with a fabricated commit hash
(`v0.0.0-20240220173807-2d6f50e3b5fb`). `go mod tidy` failed with
`invalid version: unknown revision`. This is a supply-chain integrity failure,
not just a build annoyance. A translator that invents dependency versions can also
invent dependency content.

**Root cause:**
PCDP specs say nothing about dependency provenance. The `go.mod` is left entirely
to translator discretion. For a framework whose stated goal includes supply-chain
security (EAL4+/EUCC, OBS packaging), this is a gap.

**RESOLVED — add a DEPENDENCIES section to the spec:**
The spec author declares which external libraries are required and what version
selection rules apply. The translator is bound by these rules.

```
## DEPENDENCIES

github.com/digitalocean/go-libvirt:
  version-strategy: pseudo-version    // no tagged releases exist
  do-not-fabricate: true              // translator must not invent commit hashes
  hints-file: cloud-native.go.go-libvirt.hints.md  // see hints/ hierarchy

sigs.k8s.io/controller-runtime:
  minimum-version: v0.17.0
  rationale: metricsserver.Options API available since v0.17
```

The `do-not-fabricate: true` flag is a constraint on the translator.
A translator that cannot verify a pseudo-version must flag it in the
TRANSLATION_REPORT rather than invent one.

**Proposed PCDP schema change:**
Add `## DEPENDENCIES` as an optional section in the spec, between DEPLOYMENT and
PRECONDITIONS. A DEPENDENCIES block declares module paths, version strategy, and
a reference to the relevant hints file. The translator reads it during dependency
resolution and records any deviations in the TRANSLATION_REPORT.

---

## Finding 7 — The DELIVERABLES file tree must stay language-agnostic

**Observed failure:**
The spec listed concrete filenames (`controllers/libvirt_transport.go`,
`api/v1alpha1/remotekvmhost_types.go`). This implicitly assumed Go and
controller-runtime. The file tree is not portable to a Rust or Java implementation.
Additionally, a useful file (`controllers/common.go`) was invented by all three
implementations because it was the natural factoring point for shared helpers,
but it was absent from the spec's file tree.

**What this is NOT:**
Listing concrete filenames in a spec is wrong not because filenames are unimportant
but because *the spec author does not know the language yet*. The filenames are
a template concern, not a spec concern.

**RESOLVED — two-level DELIVERABLES model:**

**Level 1 — Spec DELIVERABLES (language-agnostic):**
The spec declares *logical component categories*, not filenames.

```
## DELIVERABLES

COMPONENT: api-types
  purpose: CRD type definitions for RemoteKVMHost and RemoteVirtualMachine
  required: true

COMPONENT: transport-layer
  purpose: Implements Dialer/Session/Domain interfaces (production + test double + stub)
  required: true

COMPONENT: host-reconciler
  purpose: Implements reconcile-remote-host behaviour
  required: true

COMPONENT: vm-reconciler
  purpose: Implements reconcile-remote-vm behaviour
  required: true

COMPONENT: shared-helpers
  purpose: Shared constants, requeue durations, condition helpers, spec validators
  required: true
  note: Must be separate from reconcilers to avoid circular imports
```

**Level 2 — Template DELIVERABLES (language-specific):**
The deployment template's existing DELIVERABLES table maps logical components to
concrete filenames for the resolved language. The template already has this table;
it just needs to include the shared-helpers component:

```
| COMPONENT       | Go filename(s)                        | Notes                    |
|-----------------|---------------------------------------|--------------------------|
| api-types       | api/v1alpha1/*_types.go               |                          |
| transport-layer | controllers/libvirt_transport.go      | Dialer+Session+Domain+Fakes |
| host-reconciler | controllers/remotekvmhost_controller.go |                        |
| vm-reconciler   | controllers/remotevirtualmachine_controller.go |               |
| shared-helpers  | controllers/common.go                 | finalizers, requeue consts, helpers |
```

This keeps specs language-agnostic and lets the template own the concrete file
structure for each language.

**Proposed PCDP schema change:**
Replace filename-based DELIVERABLES in specs with component-based DELIVERABLES.
Extend the template DELIVERABLES table with a COMPONENT column that maps
logical components to concrete filenames.

---

## Finding 8 — The TRANSLATION_REPORT encourages overstatement

**Observed failure:**
Claude's report claimed 96% overall confidence. The synth meta-observation
was direct: "several claims need code-level verification before being trusted
as release documentation." In practice, the confidence claims had no backing —
they were the translator's own assessment of its own work, with no verification
mechanism.

**Root cause:**
The TRANSLATION_REPORT template requires confidence levels but not verification
methods. An AI translator will produce confident-sounding numbers. Without a
required verification method, those numbers are not meaningful.

**RESOLVED — require a VERIFICATION-METHOD for each example confidence claim:**

```
## Confidence per EXAMPLE

| EXAMPLE                    | Confidence | Verification method                         | Unverified claims |
|----------------------------|------------|---------------------------------------------|-------------------|
| host-becomes-ready         | 90%        | FakeSession unit test (TestHostBecomesReady) | libvirtVersion format |
| vm-start                   | 85%        | FakeSession unit test (TestVMStart)          | domainUUID population (needs live libvirtd) |
| vm-graceful-stop           | 60%        | No test yet                                  | Entire shutdown sequence |
```

Confidence backed by a named, runnable test carries weight.
Confidence with no test is disclosed as unverified.
The TRANSLATION_REPORT template should make this table format mandatory.

**Proposed PCDP schema change:**
Add a `Verification-method` and `Unverified-claims` column to the
TRANSLATION_REPORT confidence table. A claim is `verified` only if it references
a specific test function in `independent_tests/` that passes without a live
external service. Unverified claims must be listed explicitly, not silently omitted.

---

## Finding 9 — INVARIANTS mix observable and implementation constraints

**Observed failure:**
The INVARIANTS section contained entries like:
- "No QEMU process remains after domain delete" (observable from outside)
- "SSH key bytes never written to the filesystem" (only verifiable by code review)

These have different verification strategies, different enforcement points, and
different relevance for safety assessments. Mixing them makes the section harder
to use as an audit artefact.

**RESOLVED — annotate within a single INVARIANTS section (Option B):**

Add a `[observable]` or `[implementation]` tag to each invariant entry.
This preserves the current single-section structure and keeps spec-author
friction low. Audit consumers and `pcdp-lint` can filter by tag.

Example:

```
## INVARIANTS

- [observable]      No QEMU process remains after domain delete
- [implementation]  SSH key bytes never written to the filesystem
- [observable]      status.phase reflects the last-observed domain state
- [implementation]  All libvirt calls are made through the Session interface
```

`pcdp-lint` must validate that every invariant carries exactly one of
`[observable]` or `[implementation]`. Untagged invariants are a lint error.

**Proposed PCDP schema change:**
Retain the single `## INVARIANTS` section. Require each invariant to carry
an inline `[observable]` or `[implementation]` tag. Add tag validation to
`pcdp-lint`.

---

## Finding 10 — The link between EXAMPLES and INDEPENDENT_TESTS is underspecified

**Observed failure:**
All three implementations produced `INDEPENDENT_TESTS.go` files. Codex and synth
produced near-empty files (package declaration + comment). Claude produced a more
complete file using simulation helpers — but not the `FakeSession`/`FakeDomain`
types that the INTERFACES section now mandates, because those interfaces were not
in the spec.

**Root cause:**
The PCDP template says "second-agent generated tests for specification verification"
but does not specify:
- Which EXAMPLES map to which test functions
- Which interfaces the tests must use
- What "runnable without external services" means in practice

With the INTERFACES section now mandated (Finding 4) and the multi-pass EXAMPLES
format introduced (Finding 3), a clean rule becomes possible.

**RESOLVED — formal link between EXAMPLES and INDEPENDENT_TESTS:**

The following rules apply to the `independent_tests/` deliverable:

1. **One test function per EXAMPLE.** Each EXAMPLE in the spec must have a
   corresponding `Test<CamelCasedExampleName>` function.

2. **Tests use only declared test doubles.** Independent tests must import
   only the interfaces and test doubles declared in the spec's INTERFACES section.
   They must not import the production `RealDialer` or call live external services.

3. **Tests must pass without infrastructure.** `go test ./independent_tests/`
   must succeed with no live Kubernetes cluster, no live libvirt daemon, and
   no network access.

4. **Multi-pass examples generate multi-step tests.** For an EXAMPLE with
   multiple WHEN/THEN pairs, the test function simulates each reconciler pass
   in sequence using the test double's state machine.

**Proposed PCDP schema change:**
Add these four rules to the template's DELIVERABLES section for
`independent-tests`. The template becomes the enforcement point, not the spec.

---

## Summary table

| # | Finding | Proposed change | Location | Status |
|---|---------|-----------------|----------|--------|
| 1 | Duration/Condition type divergence | `TYPE-BINDINGS` table in template | Deployment template | DEFERRED v0.3.13 |
| 2 | BEHAVIOR specifies contract not algorithm | `STEPS:` required in every BEHAVIOR block | PCDP spec schema | RESOLVED |
| 3 | Single-pass EXAMPLES miss intermediate states | Multi-pass WHEN/THEN with `MECHANISM:` annotation | PCDP spec schema | RESOLVED |
| 4 | Interface types have no spec home | `INTERFACES:` section in spec schema | PCDP spec schema | RESOLVED |
| 5 | External library API shapes don't belong in spec or template | `hints/` artefact in preset hierarchy | PCDP framework | RESOLVED |
| 6 | Dependency provenance unspecified | `DEPENDENCIES:` section in spec | PCDP spec schema | RESOLVED |
| 7 | File tree bakes in the language | Component-based DELIVERABLES in spec; filename mapping in template | Both | DEFERRED v0.3.13 |
| 8 | TRANSLATION_REPORT confidence is unverifiable | `Verification-method` column in confidence table | TRANSLATION_REPORT template | DEFERRED v0.3.13 |
| 9 | INVARIANTS mix observable and implementation | Annotate with `[observable]`/`[implementation]` tags | PCDP spec schema | RESOLVED |
| 10 | EXAMPLES–INDEPENDENT_TESTS link is informal | Formal one-test-per-example rule; infrastructure-free requirement | Template DELIVERABLES | DEFERRED v0.3.13 |

---

## The three cleanest wins for the next PCDP iteration

If only one change is made: **Finding 2 — add `STEPS:` to BEHAVIOR blocks.**
This is the highest-leverage change. It eliminates the largest class of
implementation divergence (structural, not just data-type) with no increase in
spec author burden. A spec author already knows the algorithm; writing it as
an ordered list is not additional work.

If only two changes are made: add **Finding 4 — the `INTERFACES:` section.**
Mandating a transport abstraction boundary with test doubles is what makes
independent tests meaningful and makes the reconciler logic testable without
infrastructure.

If three changes are made: add **Finding 5 — the `hints/` artefact.**
This is the cleanest resolution to the library API gotcha problem. It separates
library-specific knowledge from both specs and templates, puts it in the right
place in the preset hierarchy, and gives it a maintenance path independent of
the spec or template lifecycle.

