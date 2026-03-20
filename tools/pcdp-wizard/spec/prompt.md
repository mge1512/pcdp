I am providing two files:

1. cli-tool.template.md — a deployment template that defines
   the conventions, constraints, and defaults for this type of component
   under the Post-Coding Development Paradigm.

2. pcdp-lint.md — a specification for a component, written in the
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

When asked to write a LICENSE file, never attempt to write the full 
LICENSE to disk, but only include name, basic information, and 
link/reference to the autorative wording of the LICENSE.

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

Produce a translation report covering:
- Which deployment template default you used for target language,
  and whether any preset overrides it
- Which delivery mode you used and why
- Any specification ambiguities you encountered
- Any rules you could not implement exactly as written, and why
- Your confidence level per EXAMPLE (0-100%) that the implementation
  satisfies each example in the specification

Do not ask clarifying questions. If the specification is ambiguous,
make the most conservative interpretation, implement it, and note
the ambiguity in the translation report.

