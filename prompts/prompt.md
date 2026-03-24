
I am providing two files:

1. <deployment-template>.template.md — a deployment template that defines
   the conventions, constraints, and defaults for this type of component
   under the Post-Coding Development Paradigm.

2. <spec-name>.md — a specification for a component, written in the
   Post-Coding Development Paradigm format.

Your task:
Implement the component in full, exactly as specified. Do not add
features not described in the specification. Do not omit any
specified behaviour.

Derive the target language from the deployment template:
the template declares the default language and valid alternatives.
Use the default unless a project preset overrides it — if you deviate
from the default, state why explicitly.

Produce all deliverables required by the deployment template's
DELIVERABLES section. For each OUTPUT-FORMAT marked required or
supported in the template, produce the files listed in the DELIVERABLES
table. Do not enumerate these files yourself — read them from the
template's DELIVERABLES section.

If the spec contains a DELIVERABLES section with COMPONENT: entries,
map each COMPONENT to the concrete filenames defined in the template's
DELIVERABLES table. Do not invent filenames not listed there.

If the template contains a TYPE-BINDINGS table, apply it mechanically:
every logical type named in the spec (e.g. Duration, Condition) maps
to the concrete language type given in the table for the resolved
LANGUAGE. Do not substitute your own type judgement.

For each BEHAVIOR block:
- The Constraint: field (required | supported | forbidden) governs
  whether to implement it. Implement all required behaviors.
  Implement supported behaviors only if the resolved preset activates
  them. Never implement forbidden behaviors.
- Follow the STEPS: list in order. Do not reorder or skip steps.
  Implement MECHANISM: annotations exactly where specified.

If the spec contains an INTERFACES section, produce all declared
implementations: production and test doubles. Independent tests must
use only the declared test doubles — never the production implementation.

When asked to write a LICENSE file, never attempt to write the full
LICENSE to disk, but only include name, basic information, and
link/reference to the authoritative wording of the LICENSE.

Deliver the implementation as follows, depending on your environment:

1. If you have access to a filesystem or MCP server (git, GitHub,
   or similar): write the source files directly. Commit or push
   if possible, and report the location.

2. If you have code execution capability but no persistent storage:
   write the files within your execution environment and present
   them as downloadable artifacts.

3. If you are running in a browser sandbox or have no filesystem
   access: deliver the complete source code inline in your response,
   as clearly separated files with explicit filenames.

Do not attempt to compile, execute, or install anything unless
explicitly asked. Do not invent a delivery mechanism not listed above.

The deployment template describes the target runtime environment of
the generated artifact, not the environment where this prompt is
being evaluated. Do not make language or toolchain decisions based
on what is available in your current execution environment.

Produce a TRANSLATION_REPORT.md covering:
- Which deployment template default you used for target language,
  and whether any preset overrides it
- Which delivery mode you used and why
- How you applied the STEPS: lists in each BEHAVIOR block
- Which INTERFACES test doubles you produced (if INTERFACES present)
- How you applied the TYPE-BINDINGS table (if present in template)
- Which BEHAVIOR blocks had Constraint: supported or forbidden, and
  how that affected code generation
- Which COMPONENT entries from DELIVERABLES mapped to which filenames
- Any specification ambiguities you encountered
- Any rules you could not implement exactly as written, and why
- Your confidence per EXAMPLE as a table with these columns:

  | EXAMPLE | Confidence | Verification method | Unverified claims |

  A claim is verified only if it references a specific named test
  function in independent_tests/ that passes without a live external
  service. Unverified claims must be listed explicitly — never
  silently omitted.

Do not ask clarifying questions. If the specification is ambiguous,
make the most conservative interpretation, implement it, and note
the ambiguity in the translation report.

