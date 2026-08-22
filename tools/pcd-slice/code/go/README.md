<!-- generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04 -->

# pcd-slice

Derive per-behavior translation bundles from a PCD specification.

Translation harnesses read files in mechanical windows. A window that cuts a
behavior in half, or omits the types it references, forces another window or
a guess, and the specification ends up re-read in full on every pass and for
every role. `pcd-slice` computes the semantic form once, deterministically:
for each behavior one bundle holding that behavior's block, the examples that
exercise it, the invariants bound to it, the hints written for it, and the
names of the types it needs; plus one preamble holding what every bundle
shares.

The consumer contract: a translation pass for behavior *X* reads
`out/preamble.md` followed by `out/X.md` and needs no other file. The
preamble is byte-stable across the behaviors of one slicing run, which makes
it a cacheable prefix for API-side prompt caching.

The source specification is never modified and remains the durable artifact.
Every emitted file carries the SHA-256 of the sources it was derived from, so
a bundle is verifiable against the reproducibility tuple and stale bundles
are detectable.

- No network access, at build time or at run time.
- No configuration file, no environment variables.
- Standard library only; single static binary. Linux in v1.

## Installation

`pcd-slice` is distributed through the openSUSE Build Service. Installation
via `curl | sh` is not supported and never will be.

### openSUSE / SLES (zypper)

```
zypper addrepo https://download.opensuse.org/repositories/home:mge1512:pcd/openSUSE_Tumbleweed/home:mge1512:pcd.repo
zypper refresh
zypper install pcd-slice
```

### Fedora / RHEL (dnf)

```
dnf config-manager --add-repo https://download.opensuse.org/repositories/home:mge1512:pcd/Fedora_Rawhide/home:mge1512:pcd.repo
dnf install pcd-slice
```

### Debian / Ubuntu (apt)

```
echo 'deb http://download.opensuse.org/repositories/home:/mge1512:/pcd/Debian_Testing/ /' \
  | sudo tee /etc/apt/sources.list.d/pcd.list
sudo apt update
sudo apt install pcd-slice
```

Substitute the distribution directory that matches your release. The exact
project path is fixed when the package is submitted to OBS.

### From source

```
make build       # CGO_ENABLED=0 go build -o pcd-slice ./cmd/pcd-slice
make test        # runs the test suite against the built binary
make man         # pandoc pcd-slice.1.md -s -t man -o pcd-slice.1
make install     # PREFIX=/usr/local by default; honours DESTDIR
make dist        # release tarball plus companion vendor tarball
```

Requires Go 1.22 or later and `pandoc` for the man page. No network access is
needed: the tool depends on the standard library only.

## Usage

```
pcd-slice list  spec=<file> [hints=<file>]...
pcd-slice check spec=<file> [hints=<file>]... [out=<dir>]
pcd-slice slice spec=<file> [hints=<file>]... out=<dir>
pcd-slice version
pcd-slice help
```

Arguments are `key=value` pairs. POSIX `--flag` style is not accepted.

| Option        | Meaning                                                                 |
|---------------|-------------------------------------------------------------------------|
| `spec=<file>` | The specification to read. Required for every verb; must end in `.md`.  |
| `hints=<file>`| A hints file. Repeatable; order is preserved and recorded in provenance. |
| `out=<dir>`   | Output directory. Required for `slice`, optional for `check`, rejected for `list`. |

### list

Prints one tab-separated line per behavior in specification order - name,
then the counts of type closure, bound invariants, attributed examples and
attributed hints blocks - followed by a total line. `list` never fails on
findings; it reports the structure as it is, and writes no file.

```
$ pcd-slice list spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md
list	3	1	1	0
check	5	2	2	0
slice	4	5	7	0
total: 3 behaviors
```

### check

Validates that the specification and hints can be sliced completely and
unambiguously. Findings go to stderr as `rule: file:line: message`, sorted by
file then line; a summary line goes to stdout. With `out=`, `check` also
recomputes the provenance headers for the current sources and reports one
`stale-output` finding per file that differs, is missing or is unexpected.

```
$ pcd-slice check spec=spec.md
unknown-binding: spec.md:34: renmae
pcd-slice check: 1 findings
```

Finding rules: `undefined-type`, `unknown-binding`, `unassignable-example`,
`unassignable-hints`, `duplicate-behavior`, `no-behaviors`, `stale-output`.

### slice

Runs `check` first. On any finding it prints the findings, writes nothing and
exits 1. Otherwise it computes every output in memory and writes:

| File                | Content                                                                 |
|---------------------|-------------------------------------------------------------------------|
| `preamble.md`       | Header and intro prose, every preamble section in source order, the full TYPES section, the global invariants, and the hints preamble. |
| `<behavior>.md`     | Provenance line, `Requires-Types:` line, the behavior block verbatim, attributed examples, bound invariants, attributed hints blocks. |
| `MANIFEST.tsv`      | One `name<TAB>sha256` line per emitted file, sorted by name; the manifest does not list itself. |

```
$ pcd-slice slice spec=pcd-slice.spec.md hints=pcd-slice.go.hints.md out=spec.d
pcd-slice: 3 bundles, preamble, manifest -> spec.d
```

`slice` writes nothing outside `out`. It removes only files that carry a
pcd-slice provenance line; if `out` holds a foreign file, the run refuses
with exit code 2 and names the file.

Every emitted markdown file begins with one provenance line and no
timestamp, so two runs over identical inputs produce byte-identical files:

```
<!-- pcd-slice source=spec.md source-sha256=8cf3... hints=hints.md hints-sha256=77fd... version=0.1.0 -->
```

## Exit codes

| Code | Meaning                                                                        |
|------|--------------------------------------------------------------------------------|
| `0`  | Clean: `check` found nothing, or `slice` wrote its bundles.                    |
| `1`  | Findings: `check` found at least one, or `slice` refused and wrote nothing.    |
| `2`  | Invocation error: bad arguments, or a missing, unreadable or unwritable path.  |

`SIGTERM` and `SIGINT` cause a clean exit; files written by an interrupted
`slice` run are removed, so no partial output survives.

## Provenance

`pcd-slice version` prints the tool version and the SHA-256 of the
specification this binary was generated from:

```
$ pcd-slice version
pcd-slice 0.1.0
spec:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
```

## Documentation

`man 1 pcd-slice`, or `pcd-slice.1.md` in this repository.

## License

GPL-2.0-only. See `LICENSE`.
