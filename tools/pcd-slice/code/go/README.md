<!-- generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f -->

# pcd-slice

Derive per-behavior translation bundles from a Post-Coding Development
specification.

Translation harnesses read files in mechanical windows; a window that cuts a
behavior in half or omits the types it references forces another window or a
guess, and the specification is re-read in full across every pass and every
role. `pcd-slice` computes the semantic form once, deterministically: for
each behavior, one bundle holding that behavior's block, the examples that
exercise it, the invariants bound to it, the hints written for it, and the
names of the types it needs; plus one preamble holding what every bundle
shares.

The source specification is never modified and remains the durable artifact.
Every emitted markdown file carries the SHA-256 of the sources it was derived
from, so a bundle is verifiable against the reproducibility tuple and stale
bundles are detectable.

- Version: 0.1.1
- Specification: `pcd-slice.spec.md`,
  sha256 `c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f`
- License: `GPL-2.0-only`
- Platform: Linux, single static binary, no runtime dependencies

## Installation

`pcd-slice` is distributed through the openSUSE Build Service. Install it
with your distribution's package manager.

openSUSE / SLE:

    sudo zypper install pcd-slice

Debian / Ubuntu (after adding the OBS repository for your release):

    sudo apt update
    sudo apt install pcd-slice

Fedora / RHEL (after adding the OBS repository for your release):

    sudo dnf install pcd-slice

There is no installation script and no download-and-pipe-to-shell
instruction: package installation is the only supported path.

### Building from source

Requires a Go toolchain (1.22 or newer) and `pandoc` for the man page. The
implementation uses the Go standard library only, so no dependency is
fetched at build time.

    make build      # static binary ./pcd-slice
    make test       # run the test suite
    make man        # pcd-slice.1 from pcd-slice.1.md
    make install    # honours DESTDIR and PREFIX
    make dist       # release tarball pcd-slice-0.1.1.tar.gz
    make clean      # remove all build output

Every target runs unprivileged; none invokes `sudo` and none accesses the
network.

## Usage

    pcd-slice version
    pcd-slice help
    pcd-slice list  spec=<file> [hints=<file>]...
    pcd-slice check spec=<file> [hints=<file>]... [out=<dir>]
    pcd-slice slice spec=<file> [hints=<file>]... out=<dir>

Options are `key=value` pairs. POSIX-style flags (`--spec=...`) are not
accepted. The tool reads no configuration file, honours no environment
variable, and makes no network calls.

### Options

| Option         | Verbs                | Meaning |
|----------------|----------------------|---------|
| `spec=<file>`  | every verb, required | The specification to read. Must exist, be readable and end in `.md`. |
| `hints=<file>` | list, check, slice   | A hints file. Repeatable; order is preserved and recorded in the provenance comment. |
| `out=<dir>`    | slice (required), check (optional) | Output directory. Created if absent. Nothing is written outside it. |

### Verbs

`version` prints the tool name and version on one line and
`spec:<sha256-of-the-specification>` on a second.

`list` prints one tab-separated line per behavior in specification order —
name, then the counts of its type closure, bound invariants, attributed
examples and attributed hints blocks — followed by `total: <n> behaviors`.
It never fails on findings; it reports structure as-is.

`check` validates that the specification and hints can be sliced completely
and unambiguously. Findings go to stderr, one per line, as
`{rule}: {file}:{line}: {message}`, sorted by file then line; a one-line
summary goes to stdout. Rules: `undefined-type`, `unknown-binding`,
`unassignable-example`, `unassignable-hints`, `duplicate-behavior`,
`no-behaviors`, `stale-output`. With `out=` given, existing output is
compared against the recomputed form and every stale file is one finding.

`slice` writes the derived form. It runs `check` first: on any finding it
prints the findings, writes nothing and exits 1.

### Example

    $ pcd-slice list spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md
    list	3	1	1	0
    check	4	2	3	0
    slice	4	5	12	0
    total: 3 behaviors

    $ pcd-slice slice spec=pcd-slice.spec.md out=spec.d
    pcd-slice: 3 bundles, preamble, manifest -> spec.d

    $ cat spec.d/MANIFEST.tsv
    check.md	cb33e96b…
    list.md	7a5e726f…
    preamble.md	89bb31fa…
    slice.md	28e7e8da…

## Output layout

| File | Contents |
|------|----------|
| `out/preamble.md` | Shared context: version header, intro prose and every preamble section in source order, the full TYPES section, global invariants, hints preamble. |
| `out/<behavior>.md` | Provenance comment, a `Requires-Types:` line listing the closure alphabetically, the behavior block verbatim, attributed top-level examples, bound invariants and attributed hints blocks, each under a heading naming its source. |
| `out/MANIFEST.tsv` | One line per emitted file — name, tab, sha256 — sorted by name. Lists every file except itself. |

Every emitted markdown file begins with one provenance line:

    <!-- pcd-slice source=<spec> source-sha256=<hash> hints=<paths> hints-sha256=<hashes> version=<tool> -->

There is no timestamp, host name or user name anywhere in the output: two
runs over identical inputs produce byte-identical files.

The consumer contract: a translation pass for behavior *X* reads
`out/preamble.md` followed by `out/X.md` and needs no other file. The
preamble is byte-stable across the behaviors of one slicing run, which makes
it a cacheable prefix for API-side prompt caching.

An existing output directory is refused (exit 2) if it holds any entry that
is not pcd-slice's own — a subdirectory, or a file whose first line is not a
pcd-slice provenance comment. Files that are pcd-slice's own but no longer
in the manifest are removed.

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | Clean: `check` found no finding, `slice` wrote its bundles. |
| `1`  | Findings: `check` found at least one; `slice` refused and wrote nothing. |
| `2`  | Invocation error: bad arguments, missing or unreadable file, foreign entry in the output directory, unwritable output directory. |

Diagnostics and findings go to stderr; normal output goes to stdout.
`SIGTERM` and `SIGINT` cause a clean exit with no partial output.

## Documentation

`man 1 pcd-slice` after installation, or `pcd-slice.1.md` in this source
tree.
