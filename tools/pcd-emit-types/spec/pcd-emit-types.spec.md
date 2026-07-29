# pcd-emit-types

## META
Deployment:  cli-tool
Version:     0.1.0
Spec-Schema: 0.4.0
Author:      Matthias G. Eckermann <pcd@mailbox.org>
License:     GPL-2.0-only
Verification: none
Safety-Level: QM

---

pcd-emit-types is a thin CLI front end over the shared PCD engine
(libpcd, see DEPENDENCIES). It loads a specification, gates on the
TYPE-table rules, and writes the canonical JSON Schema produced by the
engine's emission behaviour. All PCD semantics - include resolution,
hashing, TYPE-table validation (RULE-22 through RULE-25), the mapping,
and the canonical output form - live in libpcd's merged specification
(lint-rules.md, types-table-rules.md, types-emit-rules.md). This spec
adds argument handling, streams, exit codes, and packaging - nothing
else.

---

## TYPES

```
SpecFile := path where file_exists AND readable AND extension = ".md"

OutPath := path
// Destination for out=<path>. Parent directory must exist.

ExitCode := 0 | 1 | 2
// 0 = emitted (warnings permitted unless strict=true)
// 1 = refused (TYPE-table or load Errors; no TYPE tables;
//     or strict=true and at least one Warning)
// 2 = invocation error (bad arguments, file not found, non-.md
//     extension, unwritable out= destination)

// Contract mirrors - authoritative definitions live in libpcd's merged
// specification (lint-rules.md). Kept here so this spec is
// self-contained for translation; the library definition governs.
Severity := Error | Warning

Diagnostic := {
  severity: Severity,
  section:  string,
  line:     u32 where line > 0,
  message:  string,
  rule:     string
}
```

---

## INTERFACES

```
SpecEngine {
  // Provided by libpcd (see DEPENDENCIES). Subset consumed by this
  // front end; signatures per libpcd.spec.md INTERFACES.
  required-methods:
    LoadSpec(path)          -> (SpecModel, error)
    EmitTypesSchema(model)  -> (text, []Diagnostic, error)
    SchemaVersion()         -> string
  implementations-required:
    production:  libpcd (pinned; see DEPENDENCIES)
    test-double: FakeEngine {
      configurable: Models map[string]SpecModel,
                    EmitText map[string]string,
                    EmitDiags map[string][]Diagnostic,
                    Errs map[string]error
    }
}
```

---

## BEHAVIOR: emit
Constraint: required

The primary operation. Emits the canonical JSON Schema for the TYPE
tables of one specification file.

INPUTS:
```
file:   SpecFile
out:    OutPath   // optional; default is stdout
strict: bool      // strict=true treats warnings as errors;
                  // default false
```

OUTPUTS:
```
exit_code: ExitCode
```

PRECONDITIONS:
- file exists and is readable
- file has `.md` extension

STEPS:
1. Verify file has `.md` extension; on failure → exit 2 with
   "error: file must have .md extension: {path}".
2. Load the spec via LoadSpec(file); on failure → exit 2 with
   "error: cannot open file: {path}".
3. If the model's load diagnostics contain any Error (unresolvable
   Includes, inclusion cycle): write each diagnostic to stderr in the
   defined format → exit 1.
4. Invoke EmitTypesSchema(model).
   a. On error "no TYPE tables found" → write
      "error: no TYPE tables found in {path}" to stderr, exit 1.
   b. On error "TYPE tables have errors; emission refused" → write
      each returned diagnostic to stderr in the defined format,
      exit 1.
5. Write any Warning diagnostics (RULE-25) to stderr. If strict=true
   and at least one Warning was written → exit 1 without producing
   output.
6. If out= was given: write the text to out; on write failure →
   exit 2 with "error: cannot write file: {out}". Then write
   "wrote {out}" to stdout.
   Else: write the text to stdout and nothing else to stdout.
7. Exit 0.

POSTCONDITIONS:
- exit_code = 0 iff schema text was produced and (strict = false OR
  no Warning was emitted)
- in stdout mode, stdout carries exactly the canonical schema text and
  nothing else; all diagnostics go to stderr
- the input file is not modified
- output is byte-identical across runs for identical input
  (engine determinism)

SIDE-EFFECTS:
- stderr: diagnostic lines and error messages, if any
- stdout: the schema text (default) or the confirmation line
  (out= mode)
- file write to out= when given
- no network calls; no environment-variable reads for behaviour
  control

---

## PRECONDITIONS

- file argument must be provided; on missing file argument:
    exit 2, write usage line to stderr (see DEPLOYMENT)
- key=value arguments must use recognised keys (out, strict);
  unrecognised key=value pairs: exit 2,
  "error: unrecognised option: {key}"

---

## POSTCONDITIONS

- pcd-emit-types never modifies its input file
- pcd-emit-types makes no network calls and reads no environment
  variables for behaviour control
- exit code is always 0, 1, or 2; no other values
- on file-not-found or unreadable:
    exit 2, stderr: "error: cannot open file: {path}"
- on file without .md extension:
    exit 2, stderr: "error: file must have .md extension: {path}"

---

## INVARIANTS

- [observable]      idempotent: two runs on the same input produce
  identical output bytes and identical exit codes
- [observable]      exit 0 never occurs when a RULE-22 through RULE-24
  Error is present
- [observable]      Warnings alone never produce exit 1 unless
  strict=true
- [observable]      exit 2 indicates invocation error only, never an
  emission result
- [observable]      in stdout mode the schema is the only bytes on
  stdout; diagnostics never contaminate the artifact stream
- [implementation]  all emission semantics come from the engine; this
  tool contains no mapping logic of its own

---

## EXAMPLES

### EXAMPLE: emit_to_stdout_golden
GIVEN:
  probe.md is a structurally valid spec whose first heading is
  `# example`, whose merged Spec-SHA256 is stipulated as
  0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef,
  and whose ## TYPES section contains exactly:
  ```markdown
  ### TYPE: Probe
  Kind: record

  | Field | Type   | Req | Constraints | Default | Description |
  |-------|--------|-----|-------------|---------|-------------|
  | name  | string | yes | len 1..16   | -       | Probe name  |
  ```
  invocation: pcd-emit-types probe.md
WHEN:
  emit runs
THEN:
  stderr = (empty)
  stdout =
  ```json
  {
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "example",
    "x-pcd-spec-sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    "$defs": {
      "Probe": {
        "type": "object",
        "properties": {
          "name": {
            "description": "Probe name",
            "type": "string",
            "minLength": 1,
            "maxLength": 16
          }
        },
        "required": ["name"],
        "additionalProperties": false
      }
    }
  }
  ```
  followed by exactly one trailing newline
  exit_code = 0

### EXAMPLE: emit_to_file
GIVEN:
  probe.md as in emit_to_stdout_golden
  invocation: pcd-emit-types out=probe.schema.json probe.md
WHEN:
  emit runs
THEN:
  probe.schema.json contains exactly the bytes of
  emit_to_stdout_golden's stdout
  stdout = "wrote probe.schema.json"
  exit_code = 0

### EXAMPLE: emit_no_type_tables
GIVEN:
  plain.md is a structurally valid spec whose ## TYPES section contains
  only prose declarations and no `### TYPE:` heading
  invocation: pcd-emit-types plain.md
WHEN:
  emit runs
THEN:
  stderr = "error: no TYPE tables found in plain.md"
  stdout = (empty)
  exit_code = 1

### EXAMPLE: emit_table_error_refused
GIVEN:
  broken.md contains in ## TYPES:
  ```markdown
  ### TYPE: Broken
  Kind: struct

  | Field | Type   | Req | Constraints | Default | Description |
  |-------|--------|-----|-------------|---------|-------------|
  | name  | string | yes | -           | -       | -           |
  ```
  invocation: pcd-emit-types broken.md
WHEN:
  emit runs
THEN:
  stderr contains one diagnostic:
    severity = Error
    rule = RULE-22
    message contains "invalid Kind: value 'struct'"
  stdout = (empty)
  exit_code = 1

### EXAMPLE: strict_warning_refused
GIVEN:
  aging.md is valid except one record row has Deprecated = yes and
  Replaced-by = -
  invocation: pcd-emit-types strict=true aging.md
WHEN:
  emit runs
THEN:
  stderr contains one diagnostic:
    severity = Warning
    rule = RULE-25
    message contains "Deprecated without Replaced-by"
  stdout = (empty)
  exit_code = 1

### EXAMPLE: file_not_found
GIVEN:
  invocation: pcd-emit-types missing.md
  missing.md does not exist
WHEN:
  pcd-emit-types is invoked
THEN:
  stderr = "error: cannot open file: missing.md"
  stdout = (empty)
  exit_code = 2

### EXAMPLE: non_md_extension
GIVEN:
  invocation: pcd-emit-types types.txt
  types.txt exists and is readable
WHEN:
  pcd-emit-types is invoked
THEN:
  stderr = "error: file must have .md extension: types.txt"
  stdout = (empty)
  exit_code = 2

### EXAMPLE: unrecognised_option
GIVEN:
  invocation: pcd-emit-types format=yaml probe.md
WHEN:
  pcd-emit-types is invoked
THEN:
  stderr = "error: unrecognised option: format"
  stdout = (empty)
  exit_code = 2

---

## DEPENDENCIES

- name: libpcd
  kind: pcd-artifact
  spec: ../../libpcd/spec/libpcd.spec.md
  version: 0.1.0
  do-not-fabricate: true
  provenance: >
    Build-time dependency on a PCD-built library. The
    TRANSLATION_REPORT of this tool must record
    `Dependency-libpcd-Version:` and `Dependency-libpcd-Spec-SHA256:`
    (the merged spec hash of the libpcd spec at the pinned version),
    extending the translation-input provenance chain across the two
    translation units. The build service enforces version alignment
    across all libpcd consumers.
  notes: >
    Public interface: SpecEngine (see INTERFACES here and
    libpcd.spec.md). Per-language import path and linkage are
    hints-layer concerns.

---

## DEPLOYMENT

Runtime: command-line tool, single static binary; libpcd is linked at
build time.

Invocation:
  pcd-emit-types <specfile.md>
  pcd-emit-types out=<path> <specfile.md>
  pcd-emit-types strict=true <specfile.md>
  pcd-emit-types version

Key=value options (all optional, precede the file argument):
  out=<path>       Write the schema to <path> instead of stdout;
                   confirmation line on stdout.
  strict=true      Treat warnings (RULE-25) as errors; exit 1 and
                   produce no output. Default: strict=false.

Commands (bare words, no file argument):
  version    Print tool version and Spec-Schema version, then exit 0.
             Format: pcd-emit-types {version} (schema {spec-schema})
             The schema version is obtained from SchemaVersion().

Output streams:
  stderr: diagnostic lines and error messages
  stdout: the canonical schema text (default) or the confirmation
          line "wrote {out}" (out= mode)

Diagnostic line format (stderr), identical to pcd-lint:
  {SEVERITY}  {file}:{line}  [{section}]  {message}

Installation:
  OBS package: pcd-tools
  Available for: openSUSE Leap, SUSE Linux Enterprise, Fedora,
  Debian/Ubuntu. No curl-based installation.
  Requires: (none at runtime; libpcd is a build-time dependency)

Platform:
  Linux (primary)
  macOS (supported)
  Windows (not supported in v1)

Signal handling note:
  Clean exit on SIGTERM/SIGINT is required; the language runtime's
  default behaviour is acceptable for a short-lived CLI tool.
  Translators document the approach in the translation report.
