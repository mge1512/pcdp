% PCD-SLICE(1) pcd-slice 0.1.0 | User Commands
% Matthias G. Eckermann <pcd@mailbox.org>
% 2026

# NAME

pcd-slice - derive per-behavior translation bundles from a PCD specification

# SYNOPSIS

**pcd-slice list** **spec=**_file_ [**hints=**_file_]...

**pcd-slice check** **spec=**_file_ [**hints=**_file_]... [**out=**_dir_]

**pcd-slice slice** **spec=**_file_ [**hints=**_file_]... **out=**_dir_

**pcd-slice version**

**pcd-slice help**

# DESCRIPTION

**pcd-slice** computes the semantic form of a Post-Coding Development
specification once, deterministically. For each behavior it derives one
bundle holding that behavior's block, the examples that exercise it, the
invariants bound to it, the hints written for it, and the names of the types
it needs; plus one preamble holding what every bundle shares.

A translation pass for behavior X reads *out/preamble.md* followed by
*out/X.md* and needs no other file. The preamble is byte-stable across the
behaviors of one slicing run, which makes it a cacheable prefix for API-side
prompt caching.

The source specification is never modified and remains the durable artifact.
Every emitted file carries the SHA-256 of the sources it was derived from, so
a bundle is verifiable against the reproducibility tuple and stale bundles
are detectable.

**pcd-slice** makes no network calls, reads no configuration file, and
honours no environment variable.

# VERBS

**list**
: Prints the structure the other verbs act on: one tab-separated line per
behavior with its name and the counts of its type closure, bound invariants,
attributed examples and attributed hints blocks, in specification order,
followed by a line **total:** _n_ **behaviors**. **list** never fails on
findings; it reports the structure as it is.

**check**
: Validates that the specification and hints can be sliced completely and
unambiguously. Findings go to stderr, one per line, in the form
_rule_**:** _file_**:**_line_**:** _message_, sorted by file then line. A
summary line goes to stdout. When **out=** is given and the directory exists,
**check** additionally recomputes the provenance headers for the current
sources and reports one **stale-output** finding per file that differs, is
missing, or is unexpected.

**slice**
: Runs **check** first; on any finding it prints the findings, writes
nothing and exits 1. Otherwise it computes every output in memory and writes
*preamble.md*, one bundle per behavior, and *MANIFEST.tsv* into **out**.

**version**
: Prints the tool version and the SHA-256 of the specification this binary
was generated from.

**help**
: Prints the usage summary.

# OPTIONS

Options are key=value pairs. POSIX flag style is not accepted.

**spec=**_file_
: The specification to read. Required for every verb. Must exist, be
readable, and end in **.md**.

**hints=**_file_
: A hints file. Repeatable; the order is preserved and recorded in the
provenance line of every emitted file. Must exist, be readable, and end in
**.md**.

**out=**_dir_
: The output directory. Required for **slice**, optional for **check**,
rejected for **list**. Created if absent. **pcd-slice** writes nothing
outside it, and refuses with exit code 2 if it holds a file that does not
carry a pcd-slice provenance line.

# FINDING RULES

**undefined-type**
: A type definition references a type-shaped name that is defined nowhere.

**unknown-binding**
: An invariant binding names no known behavior.

**unassignable-example**
: A top-level example is attributable to no behavior, neither by the invoked
identifier on its **WHEN:** line nor by a behavior name occurring as a whole
word in the example's name.

**unassignable-hints**
: A hints heading of behavior shape names no known behavior. Hints preamble
is not a finding.

**duplicate-behavior**
: Two behavior headings carry the same name, or their bundle file names
collide after **/** is replaced by **-**.

**no-behaviors**
: The specification contains no behavior heading.

**stale-output**
: Reported by **check out=**_dir_ only: an emitted file no longer matches
what the current sources would produce, is missing, or is unexpected.

# OUTPUT FILES

*out/preamble.md*
: The shared context: the header and intro prose, every preamble section in
source order, the full TYPES section, the global invariants, and the hints
preamble of each hints file.

*out/<behavior>.md*
: One bundle per behavior: the provenance line, a **Requires-Types:** line
listing the type closure alphabetically, the behavior block verbatim, the
attributed top-level examples, the bound invariants, and the attributed hints
blocks, each under a heading that names its source.

*out/MANIFEST.tsv*
: One tab-separated line per emitted file - name and SHA-256 - sorted by
name. The manifest lists every file except itself.

Every emitted markdown file begins with one provenance line:

    <!-- pcd-slice source=SPEC source-sha256=HASH hints=PATHS hints-sha256=HASHES version=TOOL -->

It carries no timestamp, no hostname and no user name, so that two runs over
identical inputs produce byte-identical files.

# EXIT STATUS

**0**
: Clean: **check** found nothing, or **slice** wrote its bundles.

**1**
: Findings: **check** found at least one, or **slice** refused and wrote
nothing.

**2**
: Invocation error: bad arguments, or a missing, unreadable or unwritable
path.

# EXAMPLES

List the structure of a specification with one hints file:

    pcd-slice list spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md

Validate before slicing:

    pcd-slice check spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md

Write the bundles, then verify later that they are still current:

    pcd-slice slice spec=pcd-slice.spec.md out=spec.d
    pcd-slice check spec=pcd-slice.spec.md out=spec.d

# ENVIRONMENT

None. **pcd-slice** takes no configuration from the environment.

# FILES

None outside the directory named by **out=**.

# PROVENANCE

This program was generated from the Post-Coding Development specification
*pcd-slice.spec.md*, sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04.
**pcd-slice version** prints that checksum, so a binary can be matched to the
specification it was derived from.

# SEE ALSO

pcd-lint(1)

# HOMEPAGE

https://github.com/mge1512/pcd/tools/pcd-slice
