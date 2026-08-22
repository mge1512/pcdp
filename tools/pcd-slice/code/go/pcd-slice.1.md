% PCD-SLICE(1) pcd-slice 0.1.1 | User Commands
% Matthias G. Eckermann <pcd@mailbox.org>
% 2026-08-22

<!-- generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f -->

# NAME

pcd-slice - derive per-behavior translation bundles from a PCD specification

# SYNOPSIS

**pcd-slice** version

**pcd-slice** help

**pcd-slice** list **spec=**_file_ [**hints=**_file_]...

**pcd-slice** check **spec=**_file_ [**hints=**_file_]... [**out=**_dir_]

**pcd-slice** slice **spec=**_file_ [**hints=**_file_]... **out=**_dir_

# DESCRIPTION

**pcd-slice** computes the semantic form of a Post-Coding Development
specification once, deterministically. For each behavior it derives one
bundle holding that behavior's block, the examples that exercise it, the
invariants bound to it, the hints written for it, and the names of the
types it needs; plus one preamble holding what every bundle shares.

The source specification is never modified and remains the durable
artifact. Every emitted markdown file carries a provenance comment naming
the sources it was derived from and their SHA-256 checksums, so a bundle is
verifiable against the reproducibility tuple and stale bundles are
detectable.

The consumer contract: a translation pass for behavior _X_ reads
_out_/**preamble.md** followed by _out_/_X_**.md** and needs no other file.
The preamble is byte-stable across the behaviors of one slicing run, which
makes it a cacheable prefix for API-side prompt caching.

Options are **key=value** pairs. POSIX-style flags are not accepted.
**pcd-slice** reads no configuration file, honours no environment variable,
and makes no network calls.

# VERBS

**version**
: Print the tool name and version on one line, and
  **spec:**_sha256-of-the-specification_ on a second.

**help**
: Print usage on standard output and exit 0.

**list**
: Print the structure the other verbs act on: one tab-separated line per
  behavior in specification order, holding the behavior name and the counts
  of its type closure, bound invariants, attributed examples and attributed
  hints blocks, followed by a final **total:** _n_ **behaviors** line.
  **list** never fails on findings; it reports structure as-is.

**check**
: Validate that the specification and hints can be sliced completely and
  unambiguously. Findings go to standard error, one per line, in the form
  _rule_**:** _file_**:**_line_**:** _message_, sorted by file then line. A
  one-line summary goes to standard output. When **out=** is given and the
  directory exists, existing output is additionally compared against the
  recomputed form and every stale file is one **stale-output** finding.

**slice**
: Write the derived form: one **preamble.md**, one bundle per behavior, and
  **MANIFEST.tsv**. If **check** finds anything, **slice** prints the
  findings, writes nothing and exits 1.

# OPTIONS

**spec=**_file_
: The specification to read. Required for every verb. Must exist, be
  readable and end in **.md**.

**hints=**_file_
: A hints file. Repeatable; the order is preserved and recorded in the
  provenance comment of every emitted file.

**out=**_dir_
: The output directory. Required for **slice**, optional for **check**
  (where it enables the staleness comparison). Created if absent. Nothing
  is ever written outside it.

# OUTPUT FILES

_out_/**preamble.md**
: Shared context: the version header, intro prose and every preamble
  section in source order, the full TYPES section, the global invariants
  and the hints preamble.

_out_/_behavior_**.md**
: One per behavior: the provenance comment, a **Requires-Types:** line
  listing the type closure alphabetically, the behavior block verbatim, the
  attributed top-level examples, the bound invariants and the attributed
  hints blocks, each under a heading naming its source. A **/** in a
  behavior name becomes **-** in the file name.

_out_/**MANIFEST.tsv**
: One line per emitted file - name, tab, SHA-256 - sorted by name. The
  manifest lists every file except itself and carries no provenance line.

An existing output directory is refused with exit code 2 if it holds any
entry that is not pcd-slice's own: a subdirectory, or a file whose first
line is not a pcd-slice provenance comment. Files that are pcd-slice's own
but no longer part of the manifest are removed.

# EXIT STATUS

**0**
: Clean. **check** found no finding; **slice** wrote its bundles.

**1**
: Findings. **check** found at least one; **slice** refused and wrote
  nothing.

**2**
: Invocation error: bad arguments, missing or unreadable file, foreign
  entry in the output directory, or an unwritable output directory.

# ENVIRONMENT

None. Behaviour is controlled by **key=value** arguments only.

# EXAMPLES

Print the structure of a specification together with its hints:

    pcd-slice list spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md

Validate before slicing, including the staleness of an earlier run:

    pcd-slice check spec=pcd-slice.spec.md out=pcd-slice.d

Write the bundles:

    pcd-slice slice spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md \
        out=pcd-slice.d

# SIGNALS

**SIGTERM** and **SIGINT** cause a clean exit; files written by the
interrupted run are removed, so no partial output survives.

# SEE ALSO

The Post-Coding Development specification of this tool, **pcd-slice.spec.md**,
SHA-256 c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f.
