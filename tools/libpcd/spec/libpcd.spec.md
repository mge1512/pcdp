# libpcd

## META
Deployment:   library
Version:      0.1.0
Spec-Schema:  0.4.0
Author:       Matthias G. Eckermann <pcd@mailbox.org>
License:      GPL-2.0-only
Verification: none
Safety-Level: QM
Includes:     ../../shared/spec/lint-rules.md
Includes:     ../../shared/spec/types-table-rules.md
Includes:     ../../shared/spec/types-emit-rules.md

---

libpcd is the shared PCD engine: one specified, translated library that
carries every PCD-specific behaviour the toolchain needs - spec loading
and include resolution, merged-hash computation, the structural lint
rules, TYPE-table parsing and validation, JSON Schema emission,
milestone state transitions, and change-impact assessment.

Front-end tools do not include this specification. They bind to the
built artifact through their DEPENDENCIES section and call the public
interface declared here (SpecEngine). Consumers at this version:

- `pcd-lint` - CLI validation front end
- `pcd-emit-types` - CLI schema emission front end
- `mcp-server-pcd` - MCP front end (tools and resources)

The rule content itself is composed in via `Includes:`; this spec is
the single translation unit for it. The binding rule, its rationale,
and the provenance obligations for depending on a PCD-built artifact
are documented in `doc/technical-reference.md`, section 20.

Note: Multiple BEHAVIOR and BEHAVIOR/INTERNAL sections are permitted.
Each describes a distinct operation or internal rule of the component.
All BEHAVIOR sections share the TYPES, INVARIANTS, and EXAMPLES
sections of the merged specification.

---

## TYPES

```
SpecPath := path where extension = ".md"

HexDigest := string where matches("^[0-9a-f]{64}$")
// SHA256, lowercase hexadecimal.

IncludeRecord := {
  path:   string,      // Includes: value as written in the host META
  file:   string,      // resolved path relative to the host spec
  sha256: HexDigest    // hash of the included file contents as read
}

SpecModel := {
  host_path:        string,               // "" for content-only loads
  title:            string,               // first "#" heading of the host
  meta:             Map<string, string>,  // host META; authoritative
  merged_text:      string,               // canonical merged spec text
  host_text:        string,               // host file text as read
  includes:         List<IncludeRecord>,  // empty when no Includes
  load_diagnostics: List<Diagnostic>      // RULE-19/RULE-21 findings
                                          // collected during resolution
}
// Diagnostic, Severity, MilestoneStatus, DeploymentTemplate and the
// other rule-domain types are contributed by the included
// lint-rules.md. TypesTableModel and the emission types are
// contributed by types-table-rules.md and types-emit-rules.md.

LintOptions := {
  check_report: bool      // evaluate RULE-18; default false
}

SpecHashResult := {
  spec_path:      string    // path of the spec file
  spec_hash:      HexDigest // SHA256 of the merged spec text (equals
                            // the host file hash when the spec
                            // declares no Includes)
  report_hash:    string    // Spec-SHA256 value from
                            // TRANSLATION_REPORT.md, or "" if absent
  match:          boolean   // true if spec_hash = report_hash
  status:         string    // "current" | "stale" | "no-report"
                            //           | "no-hash-in-report"
}

SetMilestoneResult := {
  spec_path:        string          // path of the modified spec file
  milestone_name:   string          // e.g. "0.1.0"
  previous_status:  MilestoneStatus
  new_status:       MilestoneStatus
}

ChangeImpactRecommendation := "full-regeneration" | "incremental"

ChangeImpactResult := {
  recommendation:  ChangeImpactRecommendation
  primary_factor:  string     // one-sentence reason
  structural_impact: string   // high | medium | low | none
  blast_radius:    string     // description of affected scope
  scaffold_affected: boolean
  released_milestone_affected: boolean
  consistency_risk: string    // high | medium | low
  if_incremental:  string     // what to change (empty if regeneration)
  if_regeneration: string     // what to preserve (empty if incremental)
  reasoning:       string     // full assessment narrative
}
```

---

## INTERFACES

```
Filesystem {
  // Used by every path-taking operation.
  required-methods:
    ReadFile(path)              -> (content: string, error)
    WriteFile(path, content)    -> error
  implementations-required:
    production:  OSFilesystem
    test-double: FakeFilesystem {
      configurable: Files map[string]string,
                    ReadErr map[string]error,
                    WriteErr map[string]error
    }
}

SpecEngine {
  // The public interface of libpcd. Front ends bind to the built
  // artifact and call exactly these methods; every method's semantics
  // are defined by a BEHAVIOR in this merged specification. The
  // interface is versioned by this spec's Version; the PUBLIC-API
  // rules of the library deployment template apply.
  required-methods:
    LoadSpec(path)                          -> (SpecModel, error)
    LoadSpecContent(content, filename)      -> (SpecModel, error)
    Lint(model, options)                    -> []Diagnostic
    SpecHash(model)                         -> (merged: HexDigest,
                                                host:   HexDigest)
    VerifySpecHash(path)                    -> (SpecHashResult, error)
    SetMilestoneStatus(path, name, status)  -> (SetMilestoneResult, error)
    AssessChangeImpact(change_description,
                       old_spec, new_spec,
                       existing_code)       -> (ChangeImpactResult, error)
    ParseTypes(model)                       -> (TypesTableModel,
                                                []Diagnostic)
    EmitTypesSchema(model)                  -> (text: CanonicalSchemaText,
                                                []Diagnostic, error)
    SchemaVersion()                         -> string

  implementations-required:
    production:  the library itself; no alternative production
                 implementation exists
    test-double: none required at library level. Front ends define
                 their own SpecEngine test doubles for hermetic tests.
}
```

---

## BEHAVIOR: load-spec
Constraint: required

Reads a host spec, resolves its `Includes:` directives recursively, and
produces the SpecModel with the canonical merged text. Merge semantics
follow `doc/spec-composition.md`: host META is authoritative; included
TYPES, BEHAVIORs, INVARIANTS, EXAMPLES, INTERFACES, DEPENDENCIES,
TOOLCHAIN-CONSTRAINTS, PRECONDITIONS, and POSTCONDITIONS are appended
in declaration order before the host's own content; included specs may
not declare MILESTONE or DEPLOYMENT sections.

INPUTS:
```
path: SpecPath
```

OUTPUTS:
```
model: SpecModel
```

PRECONDITIONS:
- path has .md extension

STEPS:
1. Read the host file via Filesystem.ReadFile(path); on failure →
   error "cannot open file: {path}".
2. Parse the host META; record title (first "#" heading).
3. For each `Includes:` value, in declaration order:
   a. Resolve the path relative to the host spec's directory.
   b. On unresolvable or unreadable path: append the RULE-19
      diagnostic to load_diagnostics and continue with the next value.
   c. Read the included file; record its IncludeRecord with sha256.
   d. Recurse into its own `Includes:`. On a cycle: append the RULE-21
      cycle diagnostic to load_diagnostics and do not re-enter the
      cycling file.
   e. If the included spec declares a MILESTONE or DEPLOYMENT section:
      append the RULE-21 diagnostic to load_diagnostics.
4. Construct the canonical merged text: host META first, then all
   merged sections in their defined order, included contributions in
   declaration order before the host's own. The merge is
   deterministic: same host plus same included files always produces
   the same bytes.
5. Return the SpecModel.

POSTCONDITIONS:
- model.merged_text equals model.host_text when the spec declares no
  Includes and no load diagnostics were collected
- model.includes carries one record per included file, transitive
  closure, in resolution order
- no file is modified; no network access; no environment-variable reads

LoadSpecContent(content, filename) behaves identically with the content
supplied directly; host_path = "" and relative Includes resolve against
the current working directory of the caller.

---

## BEHAVIOR: compute-spec-hash
Constraint: required

INPUTS:
```
model: SpecModel
```

OUTPUTS:
```
merged: HexDigest    // SHA256 of model.merged_text
host:   HexDigest    // SHA256 of model.host_text
```

PRECONDITIONS:
- model was produced by load-spec or LoadSpecContent

STEPS:
1. merged = SHA256(model.merged_text), lowercase hexadecimal.
2. host   = SHA256(model.host_text), lowercase hexadecimal.
3. Return both.

POSTCONDITIONS:
- merged = host iff merged_text = host_text
- repeated calls on the same model return identical values

---

## BEHAVIOR: run-lint
Constraint: required

Applies the full structural rule set to a loaded spec. This is the
engine behind `pcd-lint` and the `lint_content` / `lint_file` MCP
tools; those front ends add argument handling, rendering, and exit
codes, nothing else.

INPUTS:
```
model:   SpecModel
options: LintOptions
```

OUTPUTS:
```
diagnostics: List<Diagnostic>
```

PRECONDITIONS:
- model was produced by load-spec or LoadSpecContent

STEPS:
1. Start with model.load_diagnostics.
2. Apply the lint-validation-rules BEHAVIOR (from the included
   lint-rules.md): RULE-01 through RULE-21 in order, RULE-18 only when
   options.check_report = true. All rules run; no short-circuit.
3. If the merged TYPES section contains at least one `### TYPE:`
   heading: apply the types-table-validation-rules BEHAVIOR (from the
   included types-table-rules.md): RULE-22 through RULE-25.
4. Return all collected diagnostics in evaluation order. Presentation
   ordering (for example sorting by line) is the caller's concern.

POSTCONDITIONS:
- every diagnostic carries severity, section, message, rule identifier,
  and a 1-based line number
- the result for identical inputs is identical across calls
- the input model is not modified

---

## BEHAVIOR: verify-spec-hash
Constraint: required

Computes the merged spec hash for a spec file and compares it to the
`Spec-SHA256:` field recorded in the adjacent `TRANSLATION_REPORT.md`.
Reports whether generated artifacts are current with the spec. The
recomputed hash is the merged-text hash - equal to the host file hash
when the spec declares no Includes (RULE-18 semantics).

INPUTS:
```
spec_path: string
```

OUTPUTS:
```
result: SpecHashResult
```

PRECONDITIONS:
- spec_path is non-empty and points to a readable .md file

STEPS:
1. If spec_path is empty or does not end in ".md" → error
   "spec path must reference a .md file: {spec_path}".
2. Load the spec via LoadSpec(spec_path); on read failure → error
   "cannot open file: {spec_path}".
3. spec_hash = merged hash from SpecHash(model).
4. Look for TRANSLATION_REPORT.md in the same directory as spec_path,
   then in a `code/` subdirectory of the spec's parent directory.
   If not found: return result with status = "no-report".
5. If found: search for a line matching `Spec-SHA256: <hex>`.
   If not found: return result with status = "no-hash-in-report".
6. Extract report_hash. If spec_hash = report_hash: return result with
   match = true, status = "current"; else match = false,
   status = "stale".

POSTCONDITIONS:
- result.spec_hash is always the current merged hash of the spec
- result.status is one of: "current" | "stale" | "no-report"
  | "no-hash-in-report"
- result.match is true only when status = "current"

---

## BEHAVIOR: set-milestone-status
Constraint: required

Sets the `Status:` field of a named MILESTONE section in a spec file on
disk. Used by the agent pipeline to advance the milestone cursor.

INPUTS:
```
spec_path:      string
milestone_name: string          // exact MILESTONE label, e.g. "0.1.0"
new_status:     MilestoneStatus
```

OUTPUTS:
```
result: SetMilestoneResult
```

PRECONDITIONS:
- spec_path exists, is readable and writable, and has .md extension
- new_status is a valid MilestoneStatus value

STEPS:
1. Read spec_path via Filesystem.ReadFile; on error → error
   "cannot open file: {spec_path}".
2. Locate the `## MILESTONE: {milestone_name}` section; on not found →
   error "MILESTONE '{milestone_name}' not found in {spec_path}".
3. If new_status = "active": scan all other MILESTONE sections.
   If any other section already has `Status: active` → error
   "Cannot set MILESTONE '{milestone_name}' to active:
    MILESTONE '{other}' is already active.
    Set it to released or failed first."
4. Record previous_status (current Status: value, or "pending" if
   absent).
5. Replace or insert the `Status: {value}` line within the located
   MILESTONE section.
   MECHANISM: the Status: line must be the first non-blank line after
   the ## MILESTONE: header line. If no Status: line is present,
   insert one. Do not modify any other content in the file.
6. Write the modified content back via Filesystem.WriteFile; on error →
   error "cannot write file: {spec_path}".
7. Return SetMilestoneResult.

POSTCONDITIONS:
- spec_path on disk has exactly the Status: value changed for the named
  milestone; all other content is byte-for-byte identical
- if new_status = "active", no other milestone in the file has
  Status: active
- result.previous_status reflects the status before this call
- result.new_status = new_status

---

## BEHAVIOR: assess-change-impact
Constraint: required

Analyses a specification change and recommends full regeneration or
incremental update, applying the PCD regeneration strategy framework
(whitepaper A.19) to estimate structural impact, blast radius, scaffold
involvement, and consistency risk.

INPUTS:
```
change_description: string   // unified diff or plain-language
                             // description - required
old_spec:           string   // optional; improves blast radius analysis
new_spec:           string   // optional
existing_code:      string   // optional; enables file/function-level
                             // affected-scope identification
```

OUTPUTS:
```
result: ChangeImpactResult
```

PRECONDITIONS:
- change_description is non-empty

STEPS:
1. If change_description is empty → error
   "change_description must not be empty".
2. Parse change_description to identify affected spec sections
   (TYPES, INTERFACES, INVARIANTS, BEHAVIOR, EXAMPLES, MILESTONE, META).
3. If old_spec or new_spec provided: extract the full set of changed
   elements with their section types.
4. Evaluate structural impact:
   a. TYPES, INTERFACES, or INVARIANTS affected → "high"
   b. only BEHAVIOR STEPS or EXAMPLES affected → "low" or "medium"
      depending on count
   c. only META affected → "none"
5. scaffold_affected = true iff a MILESTONE with Scaffold: true is in
   the changed scope.
6. released_milestone_affected = true iff a MILESTONE with
   Status: released is in the changed scope.
7. Estimate blast radius: from existing_code when provided (files and
   functions referencing changed elements), else from spec
   cross-references. Classify: "1-2 BEHAVIORs" | "3-5 BEHAVIORs"
   | "5+ BEHAVIORs or shared types".
8. Assess consistency risk from codebase provenance when inferable.
9. Apply decision rules:
   - structural_impact = "high" OR scaffold_affected OR
     released_milestone_affected → "full-regeneration"
   - structural_impact = "low" AND blast radius <= 2 BEHAVIORs AND
     neither scaffold nor released milestone affected → "incremental"
   - otherwise → "full-regeneration" (conservative default)
10. Compose the reasoning narrative; populate if_incremental or
    if_regeneration as appropriate. Return ChangeImpactResult.

POSTCONDITIONS:
- result.recommendation is always set
- result.primary_factor states the single most important factor
- result.reasoning is a complete narrative suitable for the audit
  bundle

---

## BEHAVIOR: emit-types
Constraint: required

High-level emission entry point: gates on TYPE-table validity, then
lowers the tables to the canonical JSON Schema text per the included
types-emit-rules.md.

INPUTS:
```
model: SpecModel
```

OUTPUTS:
```
text:        CanonicalSchemaText
diagnostics: List<Diagnostic>     // RULE-22 .. RULE-25 findings
```

PRECONDITIONS:
- model was produced by load-spec or LoadSpecContent

STEPS:
1. Apply parse-types-tables (from the included types-table-rules.md)
   to obtain the TypesTableModel and its diagnostics.
2. If the model contains no TYPE tables → error
   "no TYPE tables found in {host_path}".
3. If the diagnostics contain any Error from RULE-22 through RULE-24 →
   error "TYPE tables have errors; emission refused", returning the
   diagnostics alongside the error.
4. Compute the merged hash via SpecHash(model); assemble EmitInputs
   from the TypesTableModel, model.title, and the merged hash.
5. Apply emit-types-schema (from the included types-emit-rules.md) to
   produce the canonical text.
6. Return the text and any Warning-level diagnostics (RULE-25).

POSTCONDITIONS:
- on success, text is byte-identical across runs for identical input
- the embedded `x-pcd-spec-sha256` equals the merged hash of the model
- emission never proceeds past RULE-22 through RULE-24 Errors

---

## TOOLCHAIN-CONSTRAINTS

```
RUNTIME-DEPENDENCIES:
  - external-libraries: forbidden
    reason: the engine is certifiable surface for every consumer;
      the implementation uses the target language's standard library
      only. Hash computation, JSON rendering, and Markdown parsing are
      implemented in-tree.

DETERMINISM:
  - emit-output: required
    // Byte-identical output for identical inputs, across runs,
    // platforms, and conforming implementations. Enables the CI
    // byte-diff gate and the cross-language differential self-test.

LINKAGE:
  - consumption: required
    // Consumed by PCD front ends at build time via the target
    // language's native linkage. Version alignment across consumers
    // is enforced by the build service (see technical-reference,
    // section 20).
```

---

## DEPENDENCIES

none. libpcd depends on the target language's standard library only
(TOOLCHAIN-CONSTRAINTS: RUNTIME-DEPENDENCIES).

---

## DELIVERABLES

COMPONENT: implementation
  files: library source modules per the resolved language's layout
  notes: >
    Partitioned by behavioural domain: spec loading and merge, hash
    computation, lint rule engine (RULE-01 through RULE-21), TYPE-table
    parsing and validation (RULE-22 through RULE-25), schema emission,
    milestone editing, change-impact assessment. The public interface
    is exactly SpecEngine plus the types it exposes; everything else is
    internal. Entry-point, per-language file names, and import paths
    are hints-layer concerns.

COMPONENT: build
  files: Makefile
  notes: >
    Targets: build, test, clean, dist. The test target runs the full
    suite including the golden-pair conformance tests and exits
    non-zero on any failure.

COMPONENT: packaging
  files: per the library deployment template's Deliverables Table
  notes: >
    Signed packages from the build service; library and development
    artifacts per the target language's packaging convention.

COMPONENT: license
  files: LICENSE
  notes: GPL-2.0-only - SPDX identifier and URL only, do not reproduce
    full text.

COMPONENT: tests
  files: independent_tests/ per the library template
  notes: >
    All tests hermetic: FakeFilesystem, no network, no live binaries.
    Must include the golden-pair emission tests (byte comparison
    against the EXAMPLES in this spec) and a lint-parity fixture set
    shared conceptually with the front ends.

COMPONENT: documentation
  files: README.md
  notes: >
    Documents the SpecEngine interface, the consumers, and the
    version-alignment rule. API reference generated per the resolved
    language's documentation convention.

COMPONENT: report
  files: TRANSLATION_REPORT.md

---

## PRECONDITIONS

- For LoadSpec: path must have .md extension
- For Lint, ParseTypes, EmitTypesSchema: the model argument was
  produced by LoadSpec or LoadSpecContent
- For VerifySpecHash and SetMilestoneStatus: the path argument
  references a readable (and, for SetMilestoneStatus, writable)
  .md file
- For AssessChangeImpact: change_description is non-empty

---

## POSTCONDITIONS

- No SpecEngine method other than SetMilestoneStatus modifies any file
  on disk; SetMilestoneStatus modifies exactly the Status: line of one
  named milestone
- No method makes network calls or reads environment variables for
  behaviour control
- All methods are deterministic: identical inputs produce identical
  outputs

---

## INVARIANTS

- [observable]      lint parity: for identical spec input, the
  diagnostics returned by run-lint are identical regardless of which
  front end invoked the engine
- [observable]      hash discipline: the merged hash returned by
  compute-spec-hash equals the host hash exactly when the spec declares
  no Includes and resolution collected no diagnostics
- [observable]      emission determinism: emit-types output is
  byte-identical across runs and conforming implementations for
  identical input
- [observable]      emission is gated: no schema text is produced while
  RULE-22 through RULE-24 Errors are present
- [observable]      read-only engine: apart from set-milestone-status,
  no operation writes to disk
- [implementation]  the public surface is exactly the SpecEngine
  interface; internal partitioning is invisible to consumers

---

## EXAMPLES

### EXAMPLE: run_lint_valid_spec
GIVEN:
  a structurally valid spec file valid.md with all required sections,
  no Includes, no TYPE tables
WHEN:
  model = LoadSpec("valid.md"); diagnostics = run-lint with
  check_report = false
THEN:
  diagnostics = (empty)
  merged hash from SpecHash(model) = host hash

### EXAMPLE: load_spec_unreadable_host
GIVEN:
  missing.md does not exist
WHEN:
  load-spec is invoked as LoadSpec("missing.md")
THEN:
  error = "cannot open file: missing.md"
  no model is returned

### EXAMPLE: emit_record_golden
GIVEN:
  a spec whose first heading is `# example`, whose merged Spec-SHA256
  is stipulated as
  0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef,
  and whose ## TYPES section contains exactly:
  ```markdown
  ### TYPE: PackageRecord
  Kind: record

  | Field   | Type    | Req | Constraints | Default | Description |
  |---------|---------|-----|-------------|---------|-------------|
  | name    | string  | yes | len 1..255  | -       | Package name as reported by rpm |
  | epoch   | integer | no  | 0..         | 0       | RPM epoch |
  | enabled | boolean | no  | -           | true    | Included in the migration set |
  ```
WHEN:
  emit-types runs on the loaded model
THEN:
  diagnostics = (empty)
  text =
  ```json
  {
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "example",
    "x-pcd-spec-sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    "$defs": {
      "PackageRecord": {
        "type": "object",
        "properties": {
          "name": {
            "description": "Package name as reported by rpm",
            "type": "string",
            "minLength": 1,
            "maxLength": 255
          },
          "epoch": {
            "description": "RPM epoch",
            "type": "integer",
            "minimum": 0,
            "default": 0
          },
          "enabled": {
            "description": "Included in the migration set",
            "type": "boolean",
            "default": true
          }
        },
        "required": ["name"],
        "additionalProperties": false
      }
    }
  }
  ```
  followed by exactly one trailing newline

### EXAMPLE: emit_variant_golden
GIVEN:
  a spec titled `# example` with the stipulated merged Spec-SHA256 of
  emit_record_golden, whose ## TYPES section contains exactly:
  ```markdown
  ### TYPE: SourceMode
  Kind: variant
  Tag: mode

  | Variant | Payload            | Description |
  |---------|--------------------|-------------|
  | Live    | -                  | Inspect the running system |
  | Chroot  | image_path: string | Inspect an offline image; mount is read-only |
  ```
WHEN:
  emit-types runs on the loaded model
THEN:
  diagnostics = (empty)
  text =
  ```json
  {
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "example",
    "x-pcd-spec-sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    "$defs": {
      "SourceMode": {
        "oneOf": [
          {
            "description": "Inspect the running system",
            "type": "object",
            "properties": {
              "mode": {
                "const": "Live"
              }
            },
            "required": ["mode"],
            "additionalProperties": false
          },
          {
            "description": "Inspect an offline image; mount is read-only",
            "type": "object",
            "properties": {
              "mode": {
                "const": "Chroot"
              },
              "image_path": {
                "type": "string"
              }
            },
            "required": ["mode", "image_path"],
            "additionalProperties": false
          }
        ]
      }
    }
  }
  ```
  followed by exactly one trailing newline

### EXAMPLE: emit_enum_golden
GIVEN:
  a spec titled `# example` with the stipulated merged Spec-SHA256 of
  emit_record_golden, whose ## TYPES section contains exactly:
  ```markdown
  ### TYPE: Channel
  Kind: enum

  | Value   | Description         | Since | Deprecated | Replaced-by |
  |---------|---------------------|-------|------------|-------------|
  | stable  | Production channel  | 0.1.0 | -          | -           |
  | testing | Pre-release channel | 0.1.0 | -          | -           |
  | factory | -                   | 0.1.0 | yes        | testing     |
  ```
WHEN:
  emit-types runs on the loaded model
THEN:
  the RULE-25 set is empty (factory carries Replaced-by)
  text =
  ```json
  {
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "example",
    "x-pcd-spec-sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    "$defs": {
      "Channel": {
        "type": "string",
        "oneOf": [
          {
            "const": "stable",
            "description": "Production channel",
            "x-pcd-since": "0.1.0"
          },
          {
            "const": "testing",
            "description": "Pre-release channel",
            "x-pcd-since": "0.1.0"
          },
          {
            "const": "factory",
            "x-pcd-since": "0.1.0",
            "deprecated": true,
            "x-pcd-replaced-by": "testing"
          }
        ]
      }
    }
  }
  ```
  followed by exactly one trailing newline

### EXAMPLE: emit_types_no_tables
GIVEN:
  a structurally valid spec notables.md whose ## TYPES section contains
  only prose declarations and no `### TYPE:` heading
WHEN:
  emit-types runs on the loaded model
THEN:
  error = "no TYPE tables found in notables.md"
  no text is produced

### EXAMPLE: emit_types_table_errors_refused
GIVEN:
  a spec whose ## TYPES section contains:
  ```markdown
  ### TYPE: Broken
  Kind: struct

  | Field | Type   | Req | Constraints | Default | Description |
  |-------|--------|-----|-------------|---------|-------------|
  | name  | string | yes | -           | -       | -           |
  ```
WHEN:
  emit-types runs on the loaded model
THEN:
  diagnostics contain one RULE-22 Error:
    message contains "invalid Kind: value 'struct'"
  error = "TYPE tables have errors; emission refused"
  no text is produced

### EXAMPLE: verify_spec_hash_not_md
GIVEN:
  spec_path = "notes.txt"
WHEN:
  verify-spec-hash is invoked
THEN:
  error = "spec path must reference a .md file: notes.txt"

### EXAMPLE: verify_spec_hash_stale
GIVEN:
  spec.md with no Includes; an adjacent TRANSLATION_REPORT.md whose
  Spec-SHA256: field records a hash different from the current file
  hash of spec.md
WHEN:
  verify-spec-hash is invoked on spec.md
THEN:
  result.match = false
  result.status = "stale"

### EXAMPLE: set_milestone_unknown_name
GIVEN:
  spec.md contains ## MILESTONE: 0.1.0 but no ## MILESTONE: 0.9.0
WHEN:
  set-milestone-status is invoked with milestone_name = "0.9.0",
  new_status = "active"
THEN:
  error = "MILESTONE '0.9.0' not found in spec.md"
  spec.md is unmodified

### EXAMPLE: assess_change_impact_empty
GIVEN:
  change_description = ""
WHEN:
  assess-change-impact is invoked
THEN:
  error = "change_description must not be empty"

---

## MILESTONE: 0.1.0
Status: active
Scaffold: true
Included BEHAVIORs: load-spec, compute-spec-hash, run-lint,
  verify-spec-hash, set-milestone-status, assess-change-impact,
  emit-types, lint-validation-rules, code-fence-tracking,
  parse-types-tables, types-table-validation-rules, emit-types-schema
Acceptance criteria:
  make test passes, including byte-exact reproduction of the three
  emission goldens (emit_record_golden, emit_variant_golden,
  emit_enum_golden) and the lint fixture set

---

## DEPLOYMENT

Runtime: linkable library. No executable entry point, no CLI, no
daemon. Consumed at build time by PCD front ends via the target
language's native linkage; the resolved language and artifact form are
governed by the library deployment template and the hints layer.

Parsing approach:
  This specification describes rule and merge semantics, not the
  internal parsing implementation. Translators are free to choose any
  parsing strategy - line-by-line state machine, AST, or other. The
  EXAMPLES section plus the consuming front ends' EXAMPLES are the
  acceptance tests. Two conventions are binding for all strategies:

  Code-fence exclusion: all content between fence markers is excluded
  from structural parsing, tracked with a fence-depth counter, not a
  boolean (see BEHAVIOR/INTERNAL: code-fence-tracking).

  Column-0 requirement: structural markers (section headers, EXAMPLE:,
  GIVEN:, WHEN:, THEN:, STEPS:, Constraint:, ### TYPE:, Kind:, Tag:)
  are recognised only at column 0 of the original untrimmed line.
  Exception: fence detection uses the trimmed line so that indented
  fences are recognised.

Consumers and version alignment:
  pcd-lint, pcd-emit-types, and mcp-server-pcd bind to this library
  through their DEPENDENCIES sections, pinning version and spec hash.
  The build service builds the library from source and enforces version
  alignment across all consumers; a front end never links a library
  version other than the one its spec pins.

Installation:
  Built and distributed as signed packages via the build service as a
  build-time dependency of the front ends. No standalone runtime
  package is required by end users.

Platform:
  Linux (primary). Front-end platform statements govern end-user
  visibility.
