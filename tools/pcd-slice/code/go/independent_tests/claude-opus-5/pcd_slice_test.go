// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
// tests by: claude-opus-5
//
// One test per EXAMPLE in the specification, plus the declared error paths,
// the INVARIANTS and the boundary conditions implied by the TYPES refinement
// predicates. Every test drives the real binary and asserts the EXAMPLE's
// THEN conditions.
package pcdslice_test

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// provenanceRE is the provenance line as STEPS 3 of BEHAVIOR: slice defines
// it. No timestamp, no hostname, lowercase hex hashes.
var provenanceRE = regexp.MustCompile(
	`^<!-- pcd-slice source=(\S+) source-sha256=([0-9a-f]{64}) hints=(\S*) hints-sha256=(\S*) version=(\S+) -->$`)

// ---------------------------------------------------------------------------
// EXAMPLE: closure_is_transitive
// ---------------------------------------------------------------------------

func TestExampleClosureIsTransitive(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specClosure)

	r := run(t, dir, "list", "spec=spec.md")
	mustExitCode(t, r, 0)

	out := lines(r.stdout)
	if len(out) != 2 {
		t.Fatalf("list printed %d lines, want 2 (one behavior + total):\n%s", len(out), r.stdout)
	}
	// name, types, invariants, examples, hints_blocks - Plan, Item and Path
	// are the transitive closure; Unused is in no closure.
	if want := "record\t3\t0\t0\t0"; out[0] != want {
		t.Errorf("list line = %q, want %q", out[0], want)
	}
	if want := "total: 1 behaviors"; out[1] != want {
		t.Errorf("total line = %q, want %q", out[1], want)
	}

	// The closure is also observable in the bundle's Requires-Types line,
	// listed alphabetically, and Unused stays in the preamble's TYPES section.
	rs := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, rs, 0)

	bundle := readFile(t, filepath.Join(dir, "spec.d", "record.md"))
	bl := lines(bundle)
	if len(bl) < 2 {
		t.Fatalf("bundle too short:\n%s", bundle)
	}
	if want := "Requires-Types: Item, Path, Plan"; bl[1] != want {
		t.Errorf("Requires-Types line = %q, want %q", bl[1], want)
	}
	mustNotContain(t, "record.md", bundle, "Unused := int")

	preamble := readFile(t, filepath.Join(dir, "spec.d", "preamble.md"))
	mustContain(t, "preamble.md", preamble, "Unused := int")
	mustContain(t, "preamble.md", preamble, "Plan := { items: list of Item }")
}

// ---------------------------------------------------------------------------
// EXAMPLE: untagged_invariant_is_global
// ---------------------------------------------------------------------------

func TestExampleUntaggedInvariantIsGlobal(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specGlobalInvariant)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 0)

	const inv = "- The store is never copied back."
	mustContain(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")), inv)
	mustNotContain(t, "fetch.md", readFile(t, filepath.Join(dir, "spec.d", "fetch.md")), inv)
	mustNotContain(t, "store.md", readFile(t, filepath.Join(dir, "spec.d", "store.md")), inv)
}

// ---------------------------------------------------------------------------
// EXAMPLE: bound_invariant_lands_in_named_bundles
// ---------------------------------------------------------------------------

func TestExampleBoundInvariantLandsInNamedBundles(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specBoundInvariant)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 0)

	const inv = "- Reset destroys nothing unfetched."
	mustContain(t, "reset.md", readFile(t, filepath.Join(dir, "spec.d", "reset.md")), inv)
	mustNotContain(t, "purge.md", readFile(t, filepath.Join(dir, "spec.d", "purge.md")), inv)
	mustNotContain(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")), inv)
}

// ---------------------------------------------------------------------------
// EXAMPLE: binding_to_unknown_behavior_is_a_finding
// ---------------------------------------------------------------------------

func TestExampleBindingToUnknownBehaviorIsAFinding(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specUnknownBinding)
	line := lineNumberOf(t, specUnknownBinding, "(binds: renmae)")

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 1)

	want := "unknown-binding: spec.md:" + itoa(line) + ": renmae"
	mustContain(t, "stderr", r.stderr, want)
	if got := strings.TrimRight(r.stdout, "\n"); got != "pcd-slice check: 1 findings" {
		t.Errorf("stdout = %q, want %q", got, "pcd-slice check: 1 findings")
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: nested_example_travels_with_its_behavior
// ---------------------------------------------------------------------------

func TestExampleNestedExampleTravelsWithItsBehavior(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specNestedExample)

	// No attribution rule is consulted for a nested example: check is clean.
	rc := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, rc, 0)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 0)

	verify := readFile(t, filepath.Join(dir, "spec.d", "verify.md"))
	mustContain(t, "verify.md", verify, "### EXAMPLE verify_reports_head")
	mustContain(t, "verify.md", verify, "the report names revision 7")

	// It travels inside the verbatim behavior block, not as appended,
	// attributed material.
	if i := strings.Index(verify, "### EXAMPLE verify_reports_head"); i >= 0 {
		if j := strings.Index(verify, "(from spec.md)"); j >= 0 && j < i {
			t.Errorf("nested example was appended as attributed material, not kept in the block:\n%s", verify)
		}
	}
	mustNotContain(t, "publish.md", readFile(t, filepath.Join(dir, "spec.d", "publish.md")), "verify_reports_head")
	mustNotContain(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")), "verify_reports_head")
}

// ---------------------------------------------------------------------------
// EXAMPLE: top_level_example_attributed_by_when_line
// ---------------------------------------------------------------------------

func TestExampleTopLevelExampleAttributedByWhenLine(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specTopLevelExample)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 0)

	lint := readFile(t, filepath.Join(dir, "spec.d", "lint.md"))
	mustContain(t, "lint.md", lint, "### EXAMPLE: lint_rejects_missing_meta")
	mustContain(t, "lint.md", lint, "the finding names the missing section")
	// Appended under a heading that names its source.
	mustContain(t, "lint.md", lint, "(from spec.md)")
	if i := strings.Index(lint, "(from spec.md)"); i < 0 || i > strings.Index(lint, "### EXAMPLE: lint_rejects_missing_meta") {
		t.Errorf("attributed example is not under a source-naming heading:\n%s", lint)
	}
	mustNotContain(t, "report.md", readFile(t, filepath.Join(dir, "spec.d", "report.md")), "lint_rejects_missing_meta")
	mustNotContain(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")), "lint_rejects_missing_meta")
}

// ---------------------------------------------------------------------------
// EXAMPLE: changelog_is_excluded
// ---------------------------------------------------------------------------

func TestExampleChangelogIsExcluded(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 0)

	rows := []string{"| 0.1.0   | first cut   |", "| 0.0.9   | draft only  |", "## Changelog"}
	for _, name := range dirEntries(t, filepath.Join(dir, "spec.d")) {
		body := readFile(t, filepath.Join(dir, "spec.d", name))
		for _, row := range rows {
			mustNotContain(t, name, body, row)
		}
	}

	// The changelog counts as excluded, not unassigned: check stays clean.
	rc := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, rc, 0)
	mustNotContain(t, "stderr", rc.stderr, "unassigned")
}

// ---------------------------------------------------------------------------
// EXAMPLE: slice_refuses_on_findings_and_writes_nothing
// ---------------------------------------------------------------------------

func TestExampleSliceRefusesOnFindingsAndWritesNothing(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specDuplicateBehavior)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 1)
	mustContain(t, "stderr", r.stderr, "duplicate-behavior")

	if _, err := os.Stat(filepath.Join(dir, "spec.d")); !os.IsNotExist(err) {
		t.Errorf("spec.d exists after a refused slice (err=%v)", err)
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: two_runs_are_byte_identical
// ---------------------------------------------------------------------------

func TestExampleTwoRunsAreByteIdentical(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	r1 := run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=one.d")
	mustExitCode(t, r1, 0)
	r2 := run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=two.d")
	mustExitCode(t, r2, 0)

	one := map[string]string{}
	for _, n := range dirEntries(t, filepath.Join(dir, "one.d")) {
		one[n] = fileSHA256(t, filepath.Join(dir, "one.d", n))
	}
	two := map[string]string{}
	for _, n := range dirEntries(t, filepath.Join(dir, "two.d")) {
		two[n] = fileSHA256(t, filepath.Join(dir, "two.d", n))
	}
	if len(one) == 0 {
		t.Fatalf("first run produced no files")
	}
	if _, ok := one["MANIFEST.tsv"]; !ok {
		t.Errorf("MANIFEST.tsv missing from the first run: %v", one)
	}
	for name, h := range one {
		if two[name] != h {
			t.Errorf("file %s differs between runs: %s vs %s", name, h, two[name])
		}
	}
	if len(one) != len(two) {
		t.Errorf("run file sets differ: %v vs %v", one, two)
	}
	// No output carries a timestamp, hostname or user name.
	for _, n := range dirEntries(t, filepath.Join(dir, "one.d")) {
		body := readFile(t, filepath.Join(dir, "one.d", n))
		for _, forbidden := range []string{os.Getenv("USER"), "timestamp="} {
			if forbidden == "" {
				continue
			}
			mustNotContain(t, n, body, forbidden)
		}
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: stale_output_is_one_finding_per_file
// ---------------------------------------------------------------------------

func TestExampleStaleOutputIsOneFindingPerFile(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)

	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)

	// A clean, current output directory is not stale.
	rc := run(t, dir, "check", "spec=spec.md", "out=spec.d")
	mustExitCode(t, rc, 0)
	if got := strings.TrimRight(rc.stdout, "\n"); got != "pcd-slice check: clean" {
		t.Errorf("stdout = %q, want %q", got, "pcd-slice check: clean")
	}

	// One word in the spec changes.
	write(t, dir, "spec.md", strings.Replace(specFull, "Removes stored entries again.", "Removes stored records again.", 1))

	rs := run(t, dir, "check", "spec=spec.md", "out=spec.d")
	mustExitCode(t, rs, 1)

	emitted := dirEntries(t, filepath.Join(dir, "spec.d"))
	if len(emitted) < 3 {
		t.Fatalf("expected at least preamble, one bundle and the manifest, got %v", emitted)
	}
	for _, name := range emitted {
		want := "stale-output: " + filepath.Join("spec.d", name)
		if !strings.Contains(rs.stderr, want) {
			t.Errorf("stderr has no stale-output line for %s:\n%s", name, rs.stderr)
		}
	}
	// Exactly one line per emitted file.
	n := 0
	for _, ln := range lines(rs.stderr) {
		if strings.HasPrefix(ln, "stale-output: ") {
			n++
		}
	}
	if n != len(emitted) {
		t.Errorf("%d stale-output lines for %d emitted files:\n%s", n, len(emitted), rs.stderr)
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: foreign_file_in_out_refuses
// ---------------------------------------------------------------------------

func TestExampleForeignFileInOutRefuses(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	const notes = "hand written notes, no provenance line\n"
	write(t, dir, "spec.d/notes.txt", notes)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 2)
	mustContain(t, "stderr", r.stderr, "error: foreign file in out: spec.d/notes.txt")

	if got := readFile(t, filepath.Join(dir, "spec.d", "notes.txt")); got != notes {
		t.Errorf("notes.txt was modified: %q", got)
	}
	if got := dirEntries(t, filepath.Join(dir, "spec.d")); len(got) != 1 || got[0] != "notes.txt" {
		t.Errorf("slice wrote into an out directory it refused: %v", got)
	}
}

// ---------------------------------------------------------------------------
// INVARIANT: sources are read-only (binds: list, check, slice)
// POSTCONDITION: list and check never write files
// ---------------------------------------------------------------------------

func TestInvariantSourcesAreReadOnly(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	specBefore := fileSHA256(t, filepath.Join(dir, "spec.md"))
	hintsBefore := fileSHA256(t, filepath.Join(dir, "hints.md"))

	for _, args := range [][]string{
		{"list", "spec=spec.md", "hints=hints.md"},
		{"check", "spec=spec.md", "hints=hints.md"},
		{"slice", "spec=spec.md", "hints=hints.md", "out=spec.d"},
	} {
		r := run(t, dir, args...)
		mustExitCode(t, r, 0)
		if got := fileSHA256(t, filepath.Join(dir, "spec.md")); got != specBefore {
			t.Fatalf("%v modified the specification", args)
		}
		if got := fileSHA256(t, filepath.Join(dir, "hints.md")); got != hintsBefore {
			t.Fatalf("%v modified the hints file", args)
		}
	}
}

func TestListAndCheckNeverWriteFiles(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	before := treeSnapshot(t, dir)
	mustExitCode(t, run(t, dir, "list", "spec=spec.md", "hints=hints.md"), 0)
	mustExitCode(t, run(t, dir, "check", "spec=spec.md", "hints=hints.md"), 0)
	if after := treeSnapshot(t, dir); !sameSnapshot(before, after) {
		t.Errorf("list or check changed the working tree:\nbefore %v\nafter  %v", before, after)
	}

	// check with out= inspects, it does not repair.
	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)
	sliced := treeSnapshot(t, dir)
	mustExitCode(t, run(t, dir, "check", "spec=spec.md", "out=spec.d"), 0)
	if after := treeSnapshot(t, dir); !sameSnapshot(sliced, after) {
		t.Errorf("check with out= changed the working tree")
	}
}

// ---------------------------------------------------------------------------
// INVARIANT: provenance (binds: check, slice)
// ---------------------------------------------------------------------------

func TestInvariantProvenanceHeader(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)
	write(t, dir, "second.md", hintsSecond)

	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "hints=second.md", "out=spec.d"), 0)

	specHash := fileSHA256(t, filepath.Join(dir, "spec.md"))
	hintsHash := fileSHA256(t, filepath.Join(dir, "hints.md"))
	secondHash := fileSHA256(t, filepath.Join(dir, "second.md"))

	names := dirEntries(t, filepath.Join(dir, "spec.d"))
	mdSeen := 0
	for _, name := range names {
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		mdSeen++
		body := readFile(t, filepath.Join(dir, "spec.d", name))
		first := lines(body)[0]
		m := provenanceRE.FindStringSubmatch(first)
		if m == nil {
			t.Fatalf("%s: first line is not a provenance comment: %q", name, first)
		}
		if m[1] != "spec.md" {
			t.Errorf("%s: source=%q, want spec.md", name, m[1])
		}
		if m[2] != specHash {
			t.Errorf("%s: source-sha256=%q, want %q", name, m[2], specHash)
		}
		// Hints entries are comma separated, in input order.
		if m[3] != "hints.md,second.md" {
			t.Errorf("%s: hints=%q, want hints.md,second.md", name, m[3])
		}
		if want := hintsHash + "," + secondHash; m[4] != want {
			t.Errorf("%s: hints-sha256=%q, want %q", name, m[4], want)
		}
		if m[5] == "" {
			t.Errorf("%s: empty version field", name)
		}
	}
	if mdSeen == 0 {
		t.Fatalf("no markdown files emitted: %v", names)
	}
}

func TestProvenanceWithoutHints(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)

	first := lines(readFile(t, filepath.Join(dir, "spec.d", "preamble.md")))[0]
	if m := provenanceRE.FindStringSubmatch(first); m == nil {
		t.Fatalf("preamble.md: first line is not a provenance comment: %q", first)
	} else if m[3] != "" || m[4] != "" {
		t.Errorf("hints fields are not empty without hints files: %q", first)
	}
}

// ---------------------------------------------------------------------------
// POSTCONDITION: out contains exactly the manifest's files plus MANIFEST.tsv,
// each file's sha256 equals its manifest entry, manifest sorted by name.
// ---------------------------------------------------------------------------

func TestManifestDescribesEveryEmittedFile(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	r := run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=spec.d")
	mustExitCode(t, r, 0)
	if got := strings.TrimRight(r.stdout, "\n"); got != "pcd-slice: 2 bundles, preamble, manifest -> spec.d" {
		t.Errorf("stdout = %q, want %q", got, "pcd-slice: 2 bundles, preamble, manifest -> spec.d")
	}

	manifest := readFile(t, filepath.Join(dir, "spec.d", "MANIFEST.tsv"))
	var names []string
	for _, ln := range lines(manifest) {
		fields := strings.Split(ln, "\t")
		if len(fields) != 2 {
			t.Fatalf("manifest line %q is not name<TAB>sha256", ln)
		}
		if fields[0] == "MANIFEST.tsv" {
			t.Errorf("the manifest lists itself")
		}
		got := fileSHA256(t, filepath.Join(dir, "spec.d", fields[0]))
		if got != fields[1] {
			t.Errorf("manifest sha256 for %s = %s, file is %s", fields[0], fields[1], got)
		}
		names = append(names, fields[0])
	}
	if !sort.StringsAreSorted(names) {
		t.Errorf("manifest is not sorted by name: %v", names)
	}

	onDisk := dirEntries(t, filepath.Join(dir, "spec.d"))
	want := append(append([]string{}, names...), "MANIFEST.tsv")
	sort.Strings(want)
	if strings.Join(onDisk, ",") != strings.Join(want, ",") {
		t.Errorf("out contains %v, manifest plus MANIFEST.tsv is %v", onDisk, want)
	}
	for _, n := range []string{"preamble.md", "record.md", "purge.md"} {
		if !contains(names, n) {
			t.Errorf("manifest does not list %s: %v", n, names)
		}
	}
}

// ---------------------------------------------------------------------------
// INVARIANT: completeness (binds: slice)
// ---------------------------------------------------------------------------

func TestInvariantCompletenessCoversEveryLineExactlyOnce(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specCompleteness)

	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)

	src := strings.Split(strings.TrimSuffix(fixture(specCompleteness), "\n"), "\n")
	changelogAt := -1
	for i, ln := range src {
		if strings.HasPrefix(ln, "## Changelog") {
			changelogAt = i
			break
		}
	}
	if changelogAt < 0 {
		t.Fatalf("fixture has no changelog")
	}

	wantCount := map[string]int{}
	for _, ln := range src[:changelogAt] {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		wantCount[ln]++
	}

	gotCount := map[string]int{}
	for _, name := range dirEntries(t, filepath.Join(dir, "spec.d")) {
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		for _, ln := range lines(readFile(t, filepath.Join(dir, "spec.d", name))) {
			gotCount[ln]++
		}
	}

	for ln, want := range wantCount {
		if gotCount[ln] != want {
			t.Errorf("source line %q appears %d times across the outputs, want %d", ln, gotCount[ln], want)
		}
	}
	for _, ln := range src[changelogAt:] {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		if gotCount[ln] != 0 {
			t.Errorf("changelog line %q leaked into the outputs", ln)
		}
	}
}

// ---------------------------------------------------------------------------
// INVARIANT: a behavior appears in exactly one bundle (binds: slice)
// ---------------------------------------------------------------------------

func TestInvariantBehaviorAppearsInExactlyOneBundle(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)

	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)

	seen := 0
	for _, name := range dirEntries(t, filepath.Join(dir, "spec.d")) {
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		if strings.Contains(readFile(t, filepath.Join(dir, "spec.d", name)), "## BEHAVIOR: record") {
			seen++
			if name != "record.md" {
				t.Errorf("behavior record appears in %s", name)
			}
		}
	}
	if seen != 1 {
		t.Errorf("behavior record appears in %d files, want 1", seen)
	}
}

// ---------------------------------------------------------------------------
// BEHAVIOR: list
// ---------------------------------------------------------------------------

func TestListCountsAndOrder(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	r := run(t, dir, "list", "spec=spec.md", "hints=hints.md")
	mustExitCode(t, r, 0)

	got := lines(r.stdout)
	want := []string{
		"record\t3\t0\t1\t1",
		"purge\t0\t1\t0\t0",
		"total: 2 behaviors",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("list output:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if r.stderr != "" {
		t.Errorf("list wrote to stderr: %q", r.stderr)
	}
}

func TestListNeverFailsOnFindings(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specDuplicateBehavior)

	r := run(t, dir, "list", "spec=spec.md")
	mustExitCode(t, r, 0)
	if got := lines(r.stdout); len(got) != 3 || got[2] != "total: 2 behaviors" {
		t.Errorf("list did not report the structure as-is:\n%s", r.stdout)
	}
}

// ---------------------------------------------------------------------------
// BEHAVIOR: check - the finding rules
// ---------------------------------------------------------------------------

func TestCheckCleanSpec(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	r := run(t, dir, "check", "spec=spec.md", "hints=hints.md")
	mustExitCode(t, r, 0)
	if got := strings.TrimRight(r.stdout, "\n"); got != "pcd-slice check: clean" {
		t.Errorf("stdout = %q, want %q", got, "pcd-slice check: clean")
	}
	if r.stderr != "" {
		t.Errorf("clean check wrote findings: %q", r.stderr)
	}
}

func TestCheckDuplicateBehavior(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specDuplicateBehavior)

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 1)
	mustContain(t, "stderr", r.stderr, "duplicate-behavior: spec.md:")
	mustContain(t, "stderr", r.stderr, "emit")
	if got := strings.TrimRight(r.stdout, "\n"); !strings.HasPrefix(got, "pcd-slice check: ") || !strings.HasSuffix(got, " findings") {
		t.Errorf("stdout = %q, want a findings summary", got)
	}
}

func TestCheckNoBehaviors(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specNoBehaviors)

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 1)
	mustContain(t, "stderr", r.stderr, "no-behaviors: spec.md:")
	if got := strings.TrimRight(r.stdout, "\n"); got != "pcd-slice check: 1 findings" {
		t.Errorf("stdout = %q, want %q", got, "pcd-slice check: 1 findings")
	}
}

func TestCheckUndefinedType(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specUndefinedType)
	line := lineNumberOf(t, specUndefinedType, "MissingItem")

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 1)
	mustContain(t, "stderr", r.stderr, "undefined-type: spec.md:"+itoa(line)+": MissingItem")
}

func TestCheckUnassignableExample(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specUnassignableExample)
	line := lineNumberOf(t, specUnassignableExample, "### EXAMPLE: two_runs_agree")

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 1)
	mustContain(t, "stderr", r.stderr, "unassignable-example: spec.md:"+itoa(line)+":")
	mustContain(t, "stderr", r.stderr, "two_runs_agree")

	// slice refuses on that finding and writes nothing.
	rs := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, rs, 1)
	if _, err := os.Stat(filepath.Join(dir, "spec.d")); !os.IsNotExist(err) {
		t.Errorf("spec.d exists after a refused slice")
	}
}

func TestCheckUnassignableHints(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "bad-hints.md", hintsUnassignable)
	line := lineNumberOf(t, hintsUnassignable, "## renmae")

	r := run(t, dir, "check", "spec=spec.md", "hints=bad-hints.md")
	mustExitCode(t, r, 1)
	mustContain(t, "stderr", r.stderr, "unassignable-hints: bad-hints.md:"+itoa(line)+":")
	mustContain(t, "stderr", r.stderr, "renmae")
}

func TestHintsPreambleIsNotAFinding(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	r := run(t, dir, "check", "spec=spec.md", "hints=hints.md")
	mustExitCode(t, r, 0)
	mustNotContain(t, "stderr", r.stderr, "unassignable-hints")
}

func TestFindingsAreSortedByFileThenLine(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specUnassignableExample)
	write(t, dir, "zz-hints.md", hintsUnassignable)

	r := run(t, dir, "check", "spec=spec.md", "hints=zz-hints.md")
	mustExitCode(t, r, 1)

	type fl struct {
		file string
		line int
	}
	var got []fl
	for _, ln := range lines(r.stderr) {
		parts := strings.SplitN(ln, ": ", 3)
		if len(parts) < 3 {
			t.Fatalf("finding line %q is not rule: file:line: message", ln)
		}
		loc := strings.Split(parts[1], ":")
		if len(loc) != 2 {
			t.Fatalf("finding location %q is not file:line", parts[1])
		}
		got = append(got, fl{loc[0], atoi(t, loc[1])})
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].file > got[i].file || (got[i-1].file == got[i].file && got[i-1].line > got[i].line) {
			t.Errorf("findings are not sorted by file then line:\n%s", r.stderr)
		}
	}
}

// ---------------------------------------------------------------------------
// BEHAVIOR: slice - hints attribution, idempotence, stale replacement
// ---------------------------------------------------------------------------

func TestHintsBlocksAreAttributed(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "hints.md", hintsFull)

	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=spec.d"), 0)

	record := readFile(t, filepath.Join(dir, "spec.d", "record.md"))
	mustContain(t, "record.md", record, "Records should be written in one pass.")
	mustContain(t, "record.md", record, "(from hints.md)")

	preamble := readFile(t, filepath.Join(dir, "spec.d", "preamble.md"))
	mustContain(t, "preamble.md", preamble, "General guidance that belongs to no behavior at all.")
	mustContain(t, "preamble.md", preamble, "Nothing behavior specific here.")
	mustNotContain(t, "preamble.md", preamble, "Records should be written in one pass.")

	purge := readFile(t, filepath.Join(dir, "spec.d", "purge.md"))
	mustNotContain(t, "purge.md", purge, "Records should be written in one pass.")
}

func TestSliceIsIdempotent(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)

	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)
	first := treeSnapshot(t, filepath.Join(dir, "spec.d"))
	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)
	second := treeSnapshot(t, filepath.Join(dir, "spec.d"))

	if !sameSnapshot(first, second) {
		t.Errorf("second slice into the same directory differs:\n%v\n%v", first, second)
	}
	mustExitCode(t, run(t, dir, "check", "spec=spec.md", "out=spec.d"), 0)
}

func TestSliceRemovesItsOwnStaleFiles(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)

	// A behavior disappears from the spec: its bundle is a file this run's
	// manifest would not list, and it carries a provenance line, so it goes.
	// The invariant bound to it has to go with it: otherwise check reports
	// unknown-binding and slice refuses (BEHAVIOR: check, step 3).
	reduced := strings.Replace(specFull, `## BEHAVIOR: purge
Constraint: required

Removes stored entries again.

STEPS:
1. Remove them.

`, "", 1)
	reduced = strings.Replace(reduced, "- Purge destroys nothing unfetched. (binds: purge)\n", "", 1)
	if reduced == specFull || strings.Contains(reduced, "binds: purge") {
		t.Fatalf("fixture edit did not apply")
	}
	write(t, dir, "spec.md", reduced)

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	mustExitCode(t, r, 0)
	if _, err := os.Stat(filepath.Join(dir, "spec.d", "purge.md")); !os.IsNotExist(err) {
		t.Errorf("stale bundle purge.md survived the re-slice (err=%v)", err)
	}
	mustExitCode(t, run(t, dir, "check", "spec=spec.md", "out=spec.d"), 0)
}

func TestSliceCreatesOutDirectory(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)

	r := run(t, dir, "slice", "spec=spec.md", "out=nested/deeper.d")
	mustExitCode(t, r, 0)
	if _, err := os.Stat(filepath.Join(dir, "nested", "deeper.d", "MANIFEST.tsv")); err != nil {
		t.Errorf("out directory was not created: %v", err)
	}
	mustContain(t, "stdout", r.stdout, "-> nested/deeper.d")
}

func TestBundleIsSelfLocating(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	mustExitCode(t, run(t, dir, "slice", "spec=spec.md", "out=spec.d"), 0)

	record := readFile(t, filepath.Join(dir, "spec.d", "record.md"))
	rl := lines(record)
	if !provenanceRE.MatchString(rl[0]) {
		t.Errorf("bundle does not start with the provenance line: %q", rl[0])
	}
	if !strings.HasPrefix(rl[1], "Requires-Types:") {
		t.Errorf("second line is not Requires-Types: %q", rl[1])
	}
	mustContain(t, "record.md", record, "## BEHAVIOR: record")
	mustContain(t, "record.md", record, "1. Store the plan.")
	mustContain(t, "record.md", record, "### EXAMPLE: record_accepts_plan")
	mustNotContain(t, "record.md", record, "## BEHAVIOR: purge")
}

// ---------------------------------------------------------------------------
// TYPES refinement predicates and invocation errors (ExitCode = 2)
// ---------------------------------------------------------------------------

func TestInvocationErrors(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	write(t, dir, "notes.txt", "not markdown\n")

	cases := []struct {
		name     string
		args     []string
		contains string
	}{
		{"no arguments", nil, "usage"},
		{"unknown verb", []string{"frobnicate", "spec=spec.md"}, "error:"},
		{"missing spec", []string{"check"}, "error:"},
		{"spec does not exist", []string{"check", "spec=absent.md"}, "error: cannot read absent.md"},
		{"spec is not markdown", []string{"check", "spec=notes.txt"}, "error:"},
		{"hints does not exist", []string{"check", "spec=spec.md", "hints=absent.md"}, "error: cannot read absent.md"},
		{"slice without out", []string{"slice", "spec=spec.md"}, "error:"},
		{"unknown option", []string{"check", "spec=spec.md", "depth=3"}, "error:"},
		{"positional argument", []string{"check", "spec.md"}, "error:"},
		{"empty spec value", []string{"check", "spec="}, "error:"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := run(t, dir, tc.args...)
			mustExitCode(t, r, 2)
			if !strings.Contains(strings.ToLower(r.stderr), tc.contains) {
				t.Errorf("stderr = %q, want it to contain %q", r.stderr, tc.contains)
			}
			if r.stdout != "" {
				t.Errorf("invocation error wrote to stdout: %q", r.stdout)
			}
		})
	}
}

func TestUnreadableSpecExitsTwo(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission bits do not make a file unreadable")
	}
	dir := sandbox(t)
	p := write(t, dir, "spec.md", specFull)
	if err := os.Chmod(p, 0o000); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(p, 0o644) })

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 2)
	mustContain(t, "stderr", r.stderr, "error: cannot read spec.md")
}

func TestUnwritableOutDirExitsTwo(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission bits do not make a directory unwritable")
	}
	dir := sandbox(t)
	write(t, dir, "spec.md", specFull)
	locked := filepath.Join(dir, "locked")
	if err := os.Mkdir(locked, 0o500); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755) })

	r := run(t, dir, "slice", "spec=spec.md", "out=locked/out.d")
	mustExitCode(t, r, 2)
	mustContain(t, "stderr", r.stderr, "error:")
}

// ---------------------------------------------------------------------------
// Bare word commands
// ---------------------------------------------------------------------------

func TestVersionEmbedsSpecHash(t *testing.T) {
	dir := sandbox(t)
	r := run(t, dir, "version")
	mustExitCode(t, r, 0)
	mustContain(t, "stdout", r.stdout, "pcd-slice")
	mustContain(t, "stdout", r.stdout, "spec:"+specSHA256)
}

func TestHelpExitsZero(t *testing.T) {
	dir := sandbox(t)
	r := run(t, dir, "help")
	mustExitCode(t, r, 0)
	low := strings.ToLower(r.stdout)
	for _, needle := range []string{"list", "check", "slice", "spec=", "hints=", "out="} {
		if !strings.Contains(low, needle) {
			t.Errorf("help output does not mention %q:\n%s", needle, r.stdout)
		}
	}
}

// ---------------------------------------------------------------------------
// small local helpers
// ---------------------------------------------------------------------------

func contains(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		return "-" + string(b)
	}
	return string(b)
}

func atoi(t *testing.T, s string) int {
	t.Helper()
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			t.Fatalf("not a number: %q", s)
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// ---------------------------------------------------------------------------
// TYPES refinement predicates: only definition-side references are checked
// ---------------------------------------------------------------------------

func TestRefinementPredicatesAreNotTypeReferences(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specRefinementPredicates)

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 0)
	mustNotContain(t, "stderr", r.stderr, "undefined-type")

	rl := run(t, dir, "list", "spec=spec.md")
	mustExitCode(t, rl, 0)
	if want := "slice\t3\t0\t0\t0"; lines(rl.stdout)[0] != want {
		t.Errorf("list line = %q, want %q", lines(rl.stdout)[0], want)
	}
}

func TestUndefinedTypeIsReportedOncePerLine(t *testing.T) {
	dir := sandbox(t)
	write(t, dir, "spec.md", specUndefinedTypeTwice)

	r := run(t, dir, "check", "spec=spec.md")
	mustExitCode(t, r, 1)
	n := 0
	for _, ln := range lines(r.stderr) {
		if strings.HasPrefix(ln, "undefined-type: ") {
			n++
		}
	}
	if n != 1 {
		t.Errorf("%d undefined-type findings for one line, want 1:\n%s", n, r.stderr)
	}
	if got := strings.TrimRight(r.stdout, "\n"); got != "pcd-slice check: 1 findings" {
		t.Errorf("stdout = %q, want %q", got, "pcd-slice check: 1 findings")
	}
}
