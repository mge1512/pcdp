# pcdp-wizard

## META
Deployment:  cli-tool
Version:     0.3.7
Spec-Schema: 0.3.7
Author:      Matthias G. Eckermann <pcdp@mailbox.org>
License:     GPL-2.0-only
Verification: none
Safety-Level: QM

---

## TYPES

```
WizardState := {
  session_id:   string where non-empty,   // UUID, generated on first run
  component:    string where non-empty,   // component name (spec title)
  started_at:   timestamp,
  last_updated: timestamp,
  sections_done: List<SectionName>,
  partial_spec:  path                     // path to in-progress .md file
}

SectionName := one_of(
  "META" | "TYPES" | "BEHAVIOR" | "PRECONDITIONS" |
  "POSTCONDITIONS" | "INVARIANTS" | "EXAMPLES" | "DEPLOYMENT"
)

StatePath := path
// ~/.config/pcdp/wizard-state/<session_id>.json
// Created on first run, deleted on successful completion.

SpecPath := path where extension = ".md"
// Output path for the generated specification.
// Default: ./<component-name>.md (lowercase, hyphen-separated)
// Override via output=<path> argument.

DeploymentTemplate := string
// Must be a known template name resolvable by pcdp-lint.
// Wizard reads available templates from TEMPLATE_DIR at runtime.

LintResult := passed | failed | skipped
// passed: pcdp-lint exits 0
// failed: pcdp-lint exits 1 (errors present)
// skipped: pcdp-lint not found in PATH
```

---

## BEHAVIOR: start-or-resume

Entry point. Determines whether to start a new session or resume an
existing one for the given component name.

INPUTS:
```
component:  string   // optional; if not provided, wizard asks interactively
output:     SpecPath // optional; default: ./<component-name>.md
list:       bool     // list-sessions command: print resumable sessions and exit
```

OUTPUTS:
```
state: WizardState
```

PRECONDITIONS:
- TEMPLATE_DIR is accessible and contains at least one *.template.md file
- If component is provided, it is non-empty
- If output is provided, its parent directory exists and is writable

POSTCONDITIONS:
- If a state file exists for the given component name in StatePath:
    resume: load existing state, continue from last incomplete section
- If no state file exists:
    new session: generate session_id, create state file, start from META
- state.partial_spec points to the in-progress .md file
- state file is written to ~/.config/pcdp/wizard-state/{session_id}.json

SIDE-EFFECTS:
- Creates ~/.config/pcdp/wizard-state/ if it does not exist
- Writes or updates state file at StatePath
- Writes or updates partial spec at state.partial_spec

---

## BEHAVIOR: interview

Walks the user through each required section interactively.
Sections already completed (in state.sections_done) are skipped.
The wizard processes sections in canonical order:
META → TYPES → BEHAVIOR → PRECONDITIONS → POSTCONDITIONS →
INVARIANTS → EXAMPLES → DEPLOYMENT

INPUTS:
```
state: WizardState
```

OUTPUTS:
```
spec:        string        // complete spec content
lint_result: LintResult
```

PRECONDITIONS:
- state is a valid WizardState
- state.partial_spec is writable

POSTCONDITIONS:
- All sections in REQUIRED_SECTIONS are present in spec
- spec is written to state.partial_spec (always, regardless of lint_result)
- pcdp-lint is run against the written file
- lint_result reflects the pcdp-lint exit code
- If lint_result = passed: state file is deleted from StatePath
- If lint_result = failed: state file is retained; diagnostics printed to stderr;
  summary written to stdout with path to spec file and path to fix
- If lint_result = skipped: warning written to stderr; spec file retained

SIDE-EFFECTS:
- Writes complete spec to state.partial_spec
- Runs pcdp-lint (if available in PATH)
- Deletes state file on success

---

## BEHAVIOR: list-sessions

Prints all resumable wizard sessions found in ~/.config/pcdp/wizard-state/.

INPUTS:
```
none
```

OUTPUTS:
```
stdout: list of resumable sessions
```

POSTCONDITIONS:
- exit_code = 0 always
- For each state file found:
    prints: session_id, component name, started_at, last_updated,
            sections completed, partial spec path
- If no sessions found: prints "No resumable sessions found."
- Output format: one session per block, human-readable

---

## BEHAVIOR/INTERNAL: interview-meta

Collects META section fields interactively.

Questions asked in order:
1. Component name (if not already provided via argument)
2. Deployment template (present list from TEMPLATE_DIR; user selects)
3. Version (default: 0.1.0)
4. Spec-Schema (default: current pcdp-wizard schema version, read-only shown)
5. Author (default: from ~/.config/pcdp/presets/ if set, else ask)
6. License (SPDX identifier; show common choices: Apache-2.0, MIT, GPL-2.0-only,
   CC-BY-4.0; free text entry also accepted)
7. Verification (present options: none, lean4, fstar, dafny, custom;
   default: none)
8. Safety-Level (present options: QM, ASIL-A/B/C/D, DAL-A/B/C/D/E;
   default: QM)

After collection: write ## META section to partial_spec.
Mark META as done in state.sections_done.

---

## BEHAVIOR/INTERNAL: interview-types

Collects TYPES section interactively.

The wizard asks:
1. "Does your component work with custom data types? (Y/n)"
   If no: write a minimal ## TYPES section with a comment and continue.
2. For each type:
   a. Type name
   b. Definition (free text; wizard shows format hint)
   c. Constraints (optional; free text)
   d. "Add another type? (Y/n)"

After collection: write ## TYPES section to partial_spec.
Mark TYPES as done in state.sections_done.

---

## BEHAVIOR/INTERNAL: interview-behavior

Collects one or more BEHAVIOR sections interactively.

The wizard asks:
1. Behavior name (e.g. "transfer", "validate", "parse")
2. Inputs (name: type pairs, one per line, empty to finish)
3. "Add another BEHAVIOR section? (Y/n)"
   If yes: repeat from step 1.
   Wizard notes: BEHAVIOR/INTERNAL sections can be added later by
   editing the spec directly.

After collection: write all ## BEHAVIOR: <n> sections to partial_spec.
Mark BEHAVIOR as done in state.sections_done.

---

## BEHAVIOR/INTERNAL: interview-conditions

Collects PRECONDITIONS and POSTCONDITIONS sections interactively.

For PRECONDITIONS:
  "List the conditions that must be true before your component runs.
   One condition per line. Empty line to finish.
   Example: from.balance >= amount"

For POSTCONDITIONS:
  "List the conditions that must be true after your component runs.
   One condition per line. Empty line to finish.
   Example: from.balance' = from.balance - amount"

After collection: write ## PRECONDITIONS and ## POSTCONDITIONS sections.
Mark both as done in state.sections_done.

---

## BEHAVIOR/INTERNAL: interview-invariants

Collects INVARIANTS section interactively.

"List conditions that must always hold, regardless of which operation runs.
 Prefix with GLOBAL: for system-wide invariants.
 One invariant per line. Empty line to finish.
 Example: GLOBAL: ∀ a: Account. a.balance >= 0"

After collection: write ## INVARIANTS section to partial_spec.
Mark INVARIANTS as done in state.sections_done.

---

## BEHAVIOR/INTERNAL: interview-examples

Collects EXAMPLES section interactively.

For each example:
1. Example name (identifier, no spaces)
2. GIVEN block: "Describe the starting state. One item per line. Empty to finish."
3. WHEN block: "Describe the operation being performed."
4. THEN block: "Describe the expected outcome. One item per line. Empty to finish."
5. "Add another example? (Y/n)"

Minimum: wizard asks for at least one example.
Wizard warns if fewer than two examples are provided:
  "Consider adding a failure/error case example for completeness."

After collection: write ## EXAMPLES section to partial_spec.
Mark EXAMPLES as done in state.sections_done.

---

## BEHAVIOR/INTERNAL: interview-deployment

Collects DEPLOYMENT section interactively.

Questions depend on the deployment template selected in META.
For cli-tool template, asks:
  1. Runtime description (default: "command-line tool, single static binary")
  2. Installation notes (default: "OBS package")
  3. Platform (default: "Linux (primary)")
  4. Any additional deployment notes (optional free text)

For other templates: generic questions about runtime, dependencies,
platform, and any special deployment requirements.

After collection: write ## DEPLOYMENT section to partial_spec.
Mark DEPLOYMENT as done in state.sections_done.

---

## PRECONDITIONS

- pcdp-wizard requires pcdp-lint to be installed for post-write validation.
  If pcdp-lint is not found in PATH, wizard warns and sets lint_result = skipped.
- TEMPLATE_DIR must contain at least one *.template.md file.
  If no templates found, wizard exits 2 with error message.
- ~/.config/pcdp/wizard-state/ must be creatable if it does not exist.
- Output file parent directory must be writable.
- For resume: session_id must correspond to an existing state file.

---

## POSTCONDITIONS

- A valid .md file is always written, even if pcdp-lint reports errors.
- State file is deleted only on successful pcdp-lint exit 0.
- State file is retained on pcdp-lint failure, enabling resume and fix.
- pcdp-wizard never modifies existing spec files without user confirmation.
  If output path already exists: ask "Overwrite {path}? (y/N)" before writing.
- pcdp-wizard does not make network calls.
- pcdp-wizard does not read environment variables for behaviour control.
- pcdp-wizard is idempotent for resume: resuming a completed session
  (state file deleted) starts a new session.

---

## INVARIANTS

- GLOBAL: the partial spec written during the session is always a valid
  subset of a PCDP specification — incomplete sections are omitted,
  never written as malformed content
- GLOBAL: state files are never left in a corrupt state;
  writes are atomic (write to temp, rename)
- GLOBAL: pcdp-wizard never deletes a user's existing spec file
  without explicit confirmation
- GLOBAL: exit codes follow pcdp conventions:
    0 = spec written and pcdp-lint passed
    1 = spec written but pcdp-lint failed (errors in generated spec)
    2 = invocation error (bad arguments, missing templates, unwritable path)
- GLOBAL: session_id is a UUID v4, unique per component per machine

---

## EXAMPLES

EXAMPLE: new_session_cli_tool
GIVEN:
  no existing wizard state for component "spec-checker"
  TEMPLATE_DIR contains cli-tool.template.md
  pcdp-lint is installed
  invocation: pcdp-wizard
WHEN:
  user answers all questions interactively
  selects Deployment: cli-tool
  provides one behavior and two examples
THEN:
  spec-checker.md written to current directory
  pcdp-lint spec-checker.md → exit 0
  state file deleted from ~/.config/pcdp/wizard-state/
  pcdp-wizard exits 0
  stdout = "✓ spec-checker.md: written and valid"

EXAMPLE: resume_incomplete_session
GIVEN:
  state file exists for component "spec-checker"
  sections_done = ["META", "TYPES"]
  invocation: pcdp-wizard
WHEN:
  wizard detects existing session
  prints: "Resuming session for 'spec-checker' (started <date>)"
  prints: "Completed: META, TYPES"
  prints: "Continuing from: BEHAVIOR"
  user completes remaining sections
THEN:
  spec-checker.md written with all sections
  pcdp-lint run on output
  state file deleted on success
  pcdp-wizard exits 0

EXAMPLE: list_sessions
GIVEN:
  two state files exist in ~/.config/pcdp/wizard-state/
  invocation: pcdp-wizard list-sessions
WHEN:
  list-sessions is invoked
THEN:
  stdout shows two session blocks with component name, dates,
  sections completed, and partial spec path
  exit_code = 0

EXAMPLE: lint_failure_retained
GIVEN:
  user completes all wizard questions
  generated spec has a validation error (e.g. invalid SPDX identifier)
  pcdp-lint exits 1
WHEN:
  wizard finishes interview
THEN:
  spec file written regardless
  pcdp-lint diagnostics printed to stderr
  stdout = "✗ mycomponent.md: written with errors — run pcdp-lint mycomponent.md to review"
  state file retained for resume
  pcdp-wizard exits 1

EXAMPLE: overwrite_confirmation
GIVEN:
  mycomponent.md already exists in current directory
  invocation: pcdp-wizard
  user selects component name "mycomponent"
WHEN:
  wizard is about to write output file
THEN:
  wizard prints: "mycomponent.md already exists. Overwrite? (y/N)"
  if user answers N: wizard exits 0 without writing
  if user answers y: wizard writes file and continues normally

EXAMPLE: no_templates_found
GIVEN:
  TEMPLATE_DIR is empty or does not exist
  invocation: pcdp-wizard
WHEN:
  wizard starts
THEN:
  stderr = "error: no deployment templates found in {TEMPLATE_DIR}"
  stdout = (empty)
  exit_code = 2

EXAMPLE: output_argument
GIVEN:
  invocation: pcdp-wizard output=~/specs/transfer.md
  no existing session for this output path
WHEN:
  wizard runs interactively
THEN:
  spec written to ~/specs/transfer.md (not to current directory)
  pcdp-lint run on ~/specs/transfer.md
  all other behaviour identical to new_session_cli_tool

---

## DEPLOYMENT

Runtime: command-line tool, single static binary, no runtime dependencies
Requires: pcdp-lint installed in PATH for post-write validation (optional but recommended)

Invocation:
  pcdp-wizard
  pcdp-wizard output=<path>
  pcdp-wizard list-sessions

Key=value options:
  output=<path>    Write spec to this path instead of default
                   Default: ./<component-name>.md

Commands (bare words):
  list-sessions    List all resumable wizard sessions and exit 0

State storage:
  ~/.config/pcdp/wizard-state/<session_id>.json

Output streams:
  stderr: warnings, pcdp-lint diagnostics
  stdout: prompts, progress, summary line

Summary line format (stdout, always last):
  ✓ {file}: written and valid                    exit 0
  ✗ {file}: written with errors                  exit 1
  ✗ {file}: written, pcdp-lint not found         exit 0 (lint skipped)

Installation:
  OBS package: pcdp-tools
  Available for: openSUSE Leap, SUSE Linux Enterprise, Fedora, Debian/Ubuntu
  No curl-based installation.

Platform:
  Linux (primary)
  macOS (supported)
  Windows (not supported in v1)
