// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
// tests by: claude-opus-5
//
// Black-box test suite for the pcd-slice CLI binary. Every test invokes
// the binary as an external process (exec.Command) and asserts on stdout,
// stderr and exit code. No implementation package is imported.
package tests

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// specSHA256 is the SHA-256 of the merged pcd-slice specification the
// implementation under test was generated from. The binary must print it
// under the `version` verb.
const specSHA256 = "c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f"

// binaryPath is the canonical binary location mandated by the cli-tool
// template's BINARY-LOCATION constraint: ../../<binary-name>, i.e. the
// project root relative to this test directory.
var binaryPath string

func TestMain(m *testing.M) {
	abs, err := filepath.Abs(filepath.Join("..", "..", "pcd-slice"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot resolve binary path: %v\n", err)
		os.Exit(1)
	}
	build := exec.Command("go", "build", "-o", abs, "./cmd/pcd-slice")
	build.Dir = filepath.Join("..", "..")
	build.Stdout = os.Stderr
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot build binary under test: %v\n", err)
		os.Exit(1)
	}
	binaryPath = abs
	os.Exit(m.Run())
}

// ---------------------------------------------------------------------------
// harness helpers
// ---------------------------------------------------------------------------

type result struct {
	stdout string
	stderr string
	code   int
}

// run executes the binary under test inside the sandbox directory dir.
func run(t *testing.T, dir string, args ...string) result {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir
	// Hermetic environment: HOME points into the sandbox; the tool honours
	// no environment variable of its own (CONFIG-ENV-VARS is forbidden).
	cmd.Env = []string{"HOME=" + dir, "PATH=" + os.Getenv("PATH")}
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	code := 0
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}
	if err != nil && cmd.ProcessState == nil {
		t.Fatalf("cannot run binary: %v", err)
	}
	return result{stdout: out.String(), stderr: errb.String(), code: code}
}

// sandbox returns a fresh per-test temporary directory. Go's testing package
// removes it (including on the failure path) when the test finishes.
func sandbox(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// mustBeInside is defence in depth: no helper may create a file outside the
// test's own temporary directory.
func mustBeInside(t *testing.T, root, target string) {
	t.Helper()
	a, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("abs(%s): %v", root, err)
	}
	b, err := filepath.Abs(target)
	if err != nil {
		t.Fatalf("abs(%s): %v", target, err)
	}
	rel, err := filepath.Rel(a, b)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		t.Fatalf("refusing to touch a path outside the sandbox: %s", target)
	}
}

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	mustBeInside(t, dir, p)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func mkdir(t *testing.T, dir, name string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	mustBeInside(t, dir, p)
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", p, err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func fileSHA(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func dirEntries(t *testing.T, dir string) []string {
	t.Helper()
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir %s: %v", dir, err)
	}
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names
}

// md converts the fixture sentinel @@@ into a markdown code fence. Go raw
// string literals cannot contain backticks, so fixtures use the sentinel.
func md(s string) string {
	return strings.ReplaceAll(s, "@@@", "```")
}

func lineOf(t *testing.T, content, needle string) int {
	t.Helper()
	for i, l := range strings.Split(content, "\n") {
		if strings.Contains(l, needle) {
			return i + 1
		}
	}
	t.Fatalf("fixture does not contain %q", needle)
	return 0
}

func countOccurrences(hay, needle string) int {
	return strings.Count(hay, needle)
}

func requireContains(t *testing.T, what, hay, needle string) {
	t.Helper()
	if !strings.Contains(hay, needle) {
		t.Fatalf("%s does not contain %q\n---- actual ----\n%s", what, needle, hay)
	}
}

func requireNotContains(t *testing.T, what, hay, needle string) {
	t.Helper()
	if strings.Contains(hay, needle) {
		t.Fatalf("%s unexpectedly contains %q\n---- actual ----\n%s", what, needle, hay)
	}
}

func requireCode(t *testing.T, r result, want int) {
	t.Helper()
	if r.code != want {
		t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s", r.code, want, r.stdout, r.stderr)
	}
}

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// demoSpec is the canonical, structurally complete fixture: title, META,
// intro prose, a TYPES block with a transitive chain plus one unused type,
// two behaviors, one global and one bound invariant, and one top-level
// example attributable by its WHEN line. `check` on this fixture is clean.
const demoSpecRaw = `# demo-tool

## META
Deployment:  cli-tool
Version:     0.1.0

Intro prose for the demo tool.

## TYPES

@@@
Plan := { items: list of Item }
Item := { path: Path }
Path := string
Unused := int
@@@

## BEHAVIOR: record
Constraint: required

Records a Plan and returns.

INPUTS:
@@@
plan: Plan
@@@

STEPS:
1. Store the plan.

## BEHAVIOR: purge
Constraint: required

Removes everything that is expired.

STEPS:
1. Drop expired entries.

## INVARIANTS
- The store is never copied back.
- Nothing is destroyed unfetched. (binds: record)

## EXAMPLES

### EXAMPLE: record_writes_plan
GIVEN:
  a plan with one item
WHEN:
  result = record(plan)
THEN:
  exit_code = 0
`

func demoSpec() string { return md(demoSpecRaw) }

const demoHintsRaw = `# demo-tool hints

General guidance that belongs to no behavior at all.

## record

Write the plan atomically.

## purge

Drop by age, oldest first.
`

func demoHints() string { return md(demoHintsRaw) }

// sliceDemo writes the canonical fixture and slices it into out.
func sliceDemo(t *testing.T) (dir string, out string) {
	t.Helper()
	dir = sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	return dir, filepath.Join(dir, "spec.d")
}

// ---------------------------------------------------------------------------
// DEPLOYMENT: version verb (spec hash embedding)
// ---------------------------------------------------------------------------

func TestVersionPrintsToolNameAndSpecHash(t *testing.T) {
	dir := sandbox(t)
	r := run(t, dir, "version")
	requireCode(t, r, 0)
	lines := strings.Split(strings.TrimRight(r.stdout, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("version must print exactly two lines, got %d: %q", len(lines), r.stdout)
	}
	if !regexp.MustCompile(`^pcd-slice \S+$`).MatchString(lines[0]) {
		t.Fatalf("first version line = %q, want \"pcd-slice <version>\"", lines[0])
	}
	if lines[1] != "spec:"+specSHA256 {
		t.Fatalf("second version line = %q, want %q", lines[1], "spec:"+specSHA256)
	}
	if r.stderr != "" {
		t.Fatalf("version wrote to stderr: %q", r.stderr)
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: closure_is_transitive
// ---------------------------------------------------------------------------

func TestExampleClosureIsTransitive(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())

	r := run(t, dir, "list", "spec=spec.md")
	requireCode(t, r, 0)
	if r.stderr != "" {
		t.Fatalf("list wrote to stderr: %q", r.stderr)
	}
	lines := strings.Split(strings.TrimRight(r.stdout, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("list output must be 2 behavior lines + total, got %d:\n%s", len(lines), r.stdout)
	}
	// name, types, invariants, examples, hints_blocks - tab separated.
	if lines[0] != "record\t3\t1\t1\t0" {
		t.Fatalf("list line 1 = %q, want %q", lines[0], "record\t3\t1\t1\t0")
	}
	if lines[1] != "purge\t0\t0\t0\t0" {
		t.Fatalf("list line 2 = %q, want %q", lines[1], "purge\t0\t0\t0\t0")
	}
	if lines[2] != "total: 2 behaviors" {
		t.Fatalf("final list line = %q, want %q", lines[2], "total: 2 behaviors")
	}

	// The closure is transitive (Plan -> Item -> Path) and Unused is in no
	// closure: it stays in the preamble's TYPES section.
	rs := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, rs, 0)
	record := readFile(t, filepath.Join(dir, "spec.d", "record.md"))
	requireContains(t, "record.md", record, "Requires-Types: Item, Path, Plan")
	requireNotContains(t, "record.md", record, "Unused")
	preamble := readFile(t, filepath.Join(dir, "spec.d", "preamble.md"))
	requireContains(t, "preamble.md", preamble, "Unused := int")
	purge := readFile(t, filepath.Join(dir, "spec.d", "purge.md"))
	requireContains(t, "purge.md", purge, "Requires-Types:")
	requireNotContains(t, "purge.md", purge, "Requires-Types: Plan")
}

// ---------------------------------------------------------------------------
// EXAMPLE: untagged_invariant_is_global
// ---------------------------------------------------------------------------

func TestExampleUntaggedInvariantIsGlobal(t *testing.T) {
	dir, out := sliceDemo(t)
	_ = dir
	const global = "The store is never copied back."
	requireContains(t, "preamble.md", readFile(t, filepath.Join(out, "preamble.md")), global)
	requireNotContains(t, "record.md", readFile(t, filepath.Join(out, "record.md")), global)
	requireNotContains(t, "purge.md", readFile(t, filepath.Join(out, "purge.md")), global)
}

// ---------------------------------------------------------------------------
// EXAMPLE: bound_invariant_lands_in_named_bundles
// ---------------------------------------------------------------------------

func TestExampleBoundInvariantLandsInNamedBundles(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: reset
Constraint: required

Clears the Store.

## BEHAVIOR: purge
Constraint: required

Deletes expired entries.

## INVARIANTS
- Reset destroys nothing unfetched. (binds: reset)
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	const inv = "Reset destroys nothing unfetched."
	requireContains(t, "reset.md", readFile(t, filepath.Join(dir, "spec.d", "reset.md")), inv)
	requireNotContains(t, "purge.md", readFile(t, filepath.Join(dir, "spec.d", "purge.md")), inv)
	requireNotContains(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")), inv)
}

// ---------------------------------------------------------------------------
// EXAMPLE: binding_to_unknown_behavior_is_a_finding
// ---------------------------------------------------------------------------

func TestExampleBindingToUnknownBehaviorIsAFinding(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: rename
Constraint: required

Renames the Store.

## INVARIANTS
- X holds. (binds: renmae)
`)
	writeFile(t, dir, "spec.md", spec)
	want := fmt.Sprintf("unknown-binding: spec.md:%d: renmae", lineOf(t, spec, "(binds: renmae)"))
	r := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r, 1)
	requireContains(t, "stderr", r.stderr, want)
	if strings.TrimRight(r.stdout, "\n") != "pcd-slice check: 1 findings" {
		t.Fatalf("stdout = %q, want %q", r.stdout, "pcd-slice check: 1 findings\n")
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: nested_example_travels_with_its_behavior
// ---------------------------------------------------------------------------

func TestExampleNestedExampleTravelsWithItsBehavior(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Head := string
@@@

## BEHAVIOR: verify
Constraint: required

Verifies the Head.

### EXAMPLE verify_reports_head
GIVEN:
  a repository
WHEN:
  the tool runs
THEN:
  the head is reported
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	verify := readFile(t, filepath.Join(dir, "spec.d", "verify.md"))
	requireContains(t, "verify.md", verify, "### EXAMPLE verify_reports_head")
	requireContains(t, "verify.md", verify, "the head is reported")
	// It travels inside the block: no attribution heading is added for it.
	requireNotContains(t, "verify.md", verify, "## Attributed examples")
	requireNotContains(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")), "verify_reports_head")
}

// ---------------------------------------------------------------------------
// EXAMPLE: top_level_example_attributed_by_when_line
// ---------------------------------------------------------------------------

func TestExampleTopLevelExampleAttributedByWhenLine(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
File := string
@@@

## BEHAVIOR: lint
Constraint: required

Lints a File.

## EXAMPLES

### EXAMPLE: lint_rejects_missing_meta
GIVEN:
  a document without META
WHEN:
  result = lint(file, strict=false)
THEN:
  exit_code = 1
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	lint := readFile(t, filepath.Join(dir, "spec.d", "lint.md"))
	requireContains(t, "lint.md", lint, "## Attributed examples from spec.md")
	requireContains(t, "lint.md", lint, "### EXAMPLE: lint_rejects_missing_meta")
	requireContains(t, "lint.md", lint, "result = lint(file, strict=false)")
	requireNotContains(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")),
		"lint_rejects_missing_meta")
}

// ---------------------------------------------------------------------------
// EXAMPLE: changelog_is_excluded
// ---------------------------------------------------------------------------

func TestExampleChangelogIsExcluded(t *testing.T) {
	dir := sandbox(t)
	spec := demoSpec() + md(`
## Changelog

| Version | Change |
|---------|--------|
| 0.1.0   | initial release of the demo tool |
| 0.1.1   | second changelog row goes nowhere |
`)
	writeFile(t, dir, "spec.md", spec)

	rc := run(t, dir, "check", "spec=spec.md")
	requireCode(t, rc, 0)
	if strings.TrimRight(rc.stdout, "\n") != "pcd-slice check: clean" {
		t.Fatalf("check stdout = %q, want clean", rc.stdout)
	}

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	for _, name := range dirEntries(t, filepath.Join(dir, "spec.d")) {
		body := readFile(t, filepath.Join(dir, "spec.d", name))
		requireNotContains(t, name, body, "initial release of the demo tool")
		requireNotContains(t, name, body, "second changelog row goes nowhere")
		requireNotContains(t, name, body, "## Changelog")
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: slice_refuses_on_findings_and_writes_nothing
// ---------------------------------------------------------------------------

func TestExampleSliceRefusesOnFindingsAndWritesNothing(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: reset
Constraint: required

Clears the Store.

## BEHAVIOR: reset
Constraint: required

Clears the Store again.
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 1)
	requireContains(t, "stderr", r.stderr, "duplicate-behavior")
	if r.stdout != "" {
		t.Fatalf("slice wrote to stdout while refusing: %q", r.stdout)
	}
	if _, err := os.Stat(filepath.Join(dir, "spec.d")); !os.IsNotExist(err) {
		t.Fatalf("spec.d must not exist after a refused slice (stat err = %v)", err)
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: two_runs_are_byte_identical
// ---------------------------------------------------------------------------

func TestExampleTwoRunsAreByteIdentical(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	writeFile(t, dir, "hints.md", demoHints())

	r1 := run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=out1")
	requireCode(t, r1, 0)
	r2 := run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=out2")
	requireCode(t, r2, 0)

	tree := func(root string) map[string]string {
		m := map[string]string{}
		for _, n := range dirEntries(t, root) {
			m[n] = fileSHA(t, filepath.Join(root, n))
		}
		return m
	}
	a := tree(filepath.Join(dir, "out1"))
	b := tree(filepath.Join(dir, "out2"))
	if len(a) == 0 {
		t.Fatal("no files emitted")
	}
	if len(a) != len(b) {
		t.Fatalf("file sets differ: %v vs %v", a, b)
	}
	for name, sum := range a {
		if b[name] != sum {
			t.Fatalf("file %s differs between runs: %s vs %s", name, sum, b[name])
		}
	}
	if readFile(t, filepath.Join(dir, "out1", "MANIFEST.tsv")) != readFile(t, filepath.Join(dir, "out2", "MANIFEST.tsv")) {
		t.Fatal("MANIFEST.tsv is not byte-identical between runs")
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: stale_output_is_one_finding_per_file
// ---------------------------------------------------------------------------

func TestExampleStaleOutputIsOneFindingPerFile(t *testing.T) {
	dir, out := sliceDemo(t)
	emitted := dirEntries(t, out)
	if len(emitted) != 4 {
		t.Fatalf("expected preamble + 2 bundles + manifest, got %v", emitted)
	}

	// change one word in the spec
	spec := strings.Replace(demoSpec(), "Records a Plan and returns.", "Records a Plan and exits.", 1)
	writeFile(t, dir, "spec.md", spec)

	r := run(t, dir, "check", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 1)
	var stale int
	for _, l := range strings.Split(strings.TrimRight(r.stderr, "\n"), "\n") {
		if strings.HasPrefix(l, "stale-output: ") {
			stale++
		}
	}
	if stale != len(emitted) {
		t.Fatalf("want one stale-output finding per emitted file (%d), got %d:\n%s",
			len(emitted), stale, r.stderr)
	}
	for _, name := range emitted {
		requireContains(t, "stderr", r.stderr, "spec.d/"+name)
	}
	if strings.TrimRight(r.stdout, "\n") != fmt.Sprintf("pcd-slice check: %d findings", stale) {
		t.Fatalf("stdout = %q, want %d findings", r.stdout, stale)
	}
}

func TestCheckWithCurrentOutIsClean(t *testing.T) {
	dir, _ := sliceDemo(t)
	r := run(t, dir, "check", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	if strings.TrimRight(r.stdout, "\n") != "pcd-slice check: clean" {
		t.Fatalf("stdout = %q, want clean", r.stdout)
	}
	if r.stderr != "" {
		t.Fatalf("clean check wrote to stderr: %q", r.stderr)
	}
}

func TestCheckAbsentOutDirProducesNoStaleFindings(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	r := run(t, dir, "check", "spec=spec.md", "out=nowhere.d")
	requireCode(t, r, 0)
	if strings.TrimRight(r.stdout, "\n") != "pcd-slice check: clean" {
		t.Fatalf("stdout = %q, want clean", r.stdout)
	}
	if _, err := os.Stat(filepath.Join(dir, "nowhere.d")); !os.IsNotExist(err) {
		t.Fatalf("check must not create the out directory (stat err = %v)", err)
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: foreign_file_in_out_refuses
// ---------------------------------------------------------------------------

func TestExampleForeignFileInOutRefuses(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	mkdir(t, dir, "spec.d")
	writeFile(t, dir, "spec.d/notes.txt", "hand written notes\n")

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 2)
	requireContains(t, "stderr", r.stderr, "error: foreign file in out: spec.d/notes.txt")
	if got := readFile(t, filepath.Join(dir, "spec.d", "notes.txt")); got != "hand written notes\n" {
		t.Fatalf("notes.txt was modified: %q", got)
	}
	if names := dirEntries(t, filepath.Join(dir, "spec.d")); len(names) != 1 {
		t.Fatalf("slice wrote files despite refusing: %v", names)
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: foreign_subdirectory_refuses
// ---------------------------------------------------------------------------

func TestExampleForeignSubdirectoryRefuses(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	mkdir(t, dir, "spec.d/notes")

	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 2)
	requireContains(t, "stderr", r.stderr, "error: foreign file in out: spec.d/notes")
	if names := dirEntries(t, filepath.Join(dir, "spec.d")); len(names) != 1 || names[0] != "notes" {
		t.Fatalf("slice modified the out directory: %v", names)
	}
}

// ---------------------------------------------------------------------------
// EXAMPLE: of_suffix_attributes_and_shared_goes_to_preamble
// ---------------------------------------------------------------------------

func TestExampleOfSuffixAttributesAndSharedGoesToPreamble(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Trail := string
@@@

## BEHAVIOR: check
Constraint: required

Validates the Trail.

## BEHAVIOR: record
Constraint: required

Appends to the Trail.

## EXAMPLES

### EXAMPLE: audit_trail_holds (of: check)
GIVEN:
  an audit trail
WHEN:
  the tool runs
THEN:
  the trail holds

### EXAMPLE: end_to_end_story (of: shared)
GIVEN:
  a whole workflow
WHEN:
  the workflow runs
THEN:
  the story is told
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)

	checkMD := readFile(t, filepath.Join(dir, "spec.d", "check.md"))
	recordMD := readFile(t, filepath.Join(dir, "spec.d", "record.md"))
	preamble := readFile(t, filepath.Join(dir, "spec.d", "preamble.md"))

	requireContains(t, "check.md", checkMD, "audit_trail_holds")
	requireContains(t, "check.md", checkMD, "the trail holds")
	requireNotContains(t, "preamble.md", preamble, "audit_trail_holds")
	requireNotContains(t, "record.md", recordMD, "audit_trail_holds")

	requireContains(t, "preamble.md", preamble, "end_to_end_story")
	requireContains(t, "preamble.md", preamble, "the story is told")
	requireNotContains(t, "check.md", checkMD, "end_to_end_story")
	requireNotContains(t, "record.md", recordMD, "end_to_end_story")
}

// ---------------------------------------------------------------------------
// EXAMPLE: unique_text_match_attributes
// ---------------------------------------------------------------------------

func TestExampleUniqueTextMatchAttributes(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Head := string
@@@

## BEHAVIOR: verify
Constraint: required

Checks the Head.

## BEHAVIOR: record
Constraint: required

Writes an entry.

## EXAMPLES

### EXAMPLE: head_is_reported
GIVEN:
  a repository with one commit
WHEN:
  the tool is invoked on it
THEN:
  verify prints the head and exits zero
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	requireContains(t, "verify.md", readFile(t, filepath.Join(dir, "spec.d", "verify.md")), "head_is_reported")
	requireNotContains(t, "record.md", readFile(t, filepath.Join(dir, "spec.d", "record.md")), "head_is_reported")
	requireNotContains(t, "preamble.md", readFile(t, filepath.Join(dir, "spec.d", "preamble.md")), "head_is_reported")
}

// ---------------------------------------------------------------------------
// EXAMPLE: rung_four_tie_is_a_finding
// ---------------------------------------------------------------------------

func TestExampleRungFourTieIsAFinding(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: reset
Constraint: required

Clears the Store.

## BEHAVIOR: purge
Constraint: required

Deletes expired entries.

## EXAMPLES

### EXAMPLE: both_are_mentioned
GIVEN:
  a populated store
WHEN:
  the tool is invoked twice
THEN:
  reset happens first and purge happens second
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r, 1)
	requireContains(t, "stderr", r.stderr, "unassignable-example")
	requireContains(t, "stderr", r.stderr, "purge, reset")
	if strings.TrimRight(r.stdout, "\n") != "pcd-slice check: 1 findings" {
		t.Fatalf("stdout = %q, want 1 findings", r.stdout)
	}
}

func TestUnplaceableExampleIsAFinding(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: reset
Constraint: required

Clears the Store.

## EXAMPLES

### EXAMPLE: nothing_matches_here
GIVEN:
  an empty situation
WHEN:
  nothing at all happens
THEN:
  nothing is observed
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r, 1)
	requireContains(t, "stderr", r.stderr, "unassignable-example")
	requireContains(t, "stderr", r.stderr, "nothing_matches_here")
}

// ---------------------------------------------------------------------------
// EXAMPLE: bound_invariant_replicates_once_per_binding
// ---------------------------------------------------------------------------

func TestExampleBoundInvariantReplicatesOncePerBinding(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: check
Constraint: required

Validates the Store.

## BEHAVIOR: slice
Constraint: required

Derives bundles.

## BEHAVIOR: list
Constraint: required

Prints structure.

## INVARIANTS
- Nothing is lost. (binds: check, slice)
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)

	const inv = "Nothing is lost."
	total := 0
	for _, name := range dirEntries(t, filepath.Join(dir, "spec.d")) {
		if name == "MANIFEST.tsv" {
			continue
		}
		body := readFile(t, filepath.Join(dir, "spec.d", name))
		n := countOccurrences(body, inv)
		switch name {
		case "check.md", "slice.md":
			if n != 1 {
				t.Fatalf("%s contains the bound invariant %d times, want exactly 1", name, n)
			}
		default:
			if n != 0 {
				t.Fatalf("%s contains the bound invariant %d times, want 0", name, n)
			}
		}
		total += n
	}
	if total != 2 {
		t.Fatalf("bound invariant appears %d times overall, want 2 (once per binding)", total)
	}
}

// ---------------------------------------------------------------------------
// findings: remaining rules
// ---------------------------------------------------------------------------

func TestNoBehaviorsIsAFinding(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## META
Deployment: cli-tool

## TYPES

@@@
Store := string
@@@
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r, 1)
	requireContains(t, "stderr", r.stderr, "no-behaviors: spec.md:")
	requireContains(t, "stdout", r.stdout, "pcd-slice check: ")
}

func TestUndefinedTypeIsAFindingInsideDefinitionsOnly(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Plan := { item: MissingThing }
Store := string
@@@

## BEHAVIOR: record
Constraint: required

Records a Plan; prose mentioning OtherThing is never scanned.
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r, 1)
	wantLine := lineOf(t, spec, "MissingThing")
	requireContains(t, "stderr", r.stderr, fmt.Sprintf("undefined-type: spec.md:%d:", wantLine))
	requireContains(t, "stderr", r.stderr, "MissingThing")
	// Behavior prose is never scanned by this rule.
	requireNotContains(t, "stderr", r.stderr, "OtherThing")
}

func TestDuplicateBehaviorIsAFinding(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: reset
Constraint: required

Clears the Store.

## BEHAVIOR: reset
Constraint: required

Clears it again.
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r, 1)
	requireContains(t, "stderr", r.stderr, "duplicate-behavior: spec.md:")
	requireContains(t, "stderr", r.stderr, "reset")
}

func TestUnassignableHintsHeadingIsAFinding(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	hints := md(`# hints

Some general guidance.

## renmae

This heading names no known behavior.
`)
	writeFile(t, dir, "hints.md", hints)
	r := run(t, dir, "check", "spec=spec.md", "hints=hints.md")
	requireCode(t, r, 1)
	requireContains(t, "stderr", r.stderr,
		fmt.Sprintf("unassignable-hints: hints.md:%d:", lineOf(t, hints, "## renmae")))
	requireContains(t, "stderr", r.stderr, "renmae")
}

func TestHintsPreambleIsNotAFinding(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	writeFile(t, dir, "hints.md", demoHints())
	r := run(t, dir, "check", "spec=spec.md", "hints=hints.md")
	requireCode(t, r, 0)
	if strings.TrimRight(r.stdout, "\n") != "pcd-slice check: clean" {
		t.Fatalf("stdout = %q, want clean", r.stdout)
	}
}

func TestFindingsAreSortedByFileThenLine(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Plan := { item: MissingThing }
@@@

## BEHAVIOR: reset
Constraint: required

Clears the Plan.

## INVARIANTS
- A holds. (binds: nosuch)
- B holds. (binds: alsonosuch)
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r, 1)
	lines := strings.Split(strings.TrimRight(r.stderr, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 findings, got:\n%s", r.stderr)
	}
	re := regexp.MustCompile(`^[a-z-]+: ([^:]+):(\d+): `)
	prevFile, prevLine := "", -1
	for _, l := range lines {
		m := re.FindStringSubmatch(l)
		if m == nil {
			t.Fatalf("finding line does not match \"{rule}: {file}:{line}: {message}\": %q", l)
		}
		file := m[1]
		var line int
		fmt.Sscanf(m[2], "%d", &line)
		if file < prevFile || (file == prevFile && line < prevLine) {
			t.Fatalf("findings not sorted by file then line:\n%s", r.stderr)
		}
		prevFile, prevLine = file, line
	}
}

// ---------------------------------------------------------------------------
// hints attribution in emitted bundles
// ---------------------------------------------------------------------------

func TestHintsBlocksTravelWithTheirBehavior(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	writeFile(t, dir, "hints.md", demoHints())

	rl := run(t, dir, "list", "spec=spec.md", "hints=hints.md")
	requireCode(t, rl, 0)
	requireContains(t, "list stdout", rl.stdout, "record\t3\t1\t1\t1")
	requireContains(t, "list stdout", rl.stdout, "purge\t0\t0\t0\t1")

	r := run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=spec.d")
	requireCode(t, r, 0)
	record := readFile(t, filepath.Join(dir, "spec.d", "record.md"))
	purge := readFile(t, filepath.Join(dir, "spec.d", "purge.md"))
	preamble := readFile(t, filepath.Join(dir, "spec.d", "preamble.md"))

	requireContains(t, "record.md", record, "## Hints from hints.md")
	requireContains(t, "record.md", record, "Write the plan atomically.")
	requireNotContains(t, "record.md", record, "Drop by age, oldest first.")
	requireContains(t, "purge.md", purge, "Drop by age, oldest first.")
	requireContains(t, "preamble.md", preamble, "General guidance that belongs to no behavior at all.")
	requireNotContains(t, "preamble.md", preamble, "Write the plan atomically.")
}

// ---------------------------------------------------------------------------
// slice output structure, provenance, manifest
// ---------------------------------------------------------------------------

func TestSliceSummaryAndOutputSet(t *testing.T) {
	dir, out := sliceDemo(t)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	if strings.TrimRight(r.stdout, "\n") != "pcd-slice: 2 bundles, preamble, manifest -> spec.d" {
		t.Fatalf("summary = %q", r.stdout)
	}
	if r.stderr != "" {
		t.Fatalf("clean slice wrote to stderr: %q", r.stderr)
	}
	got := dirEntries(t, out)
	want := []string{"MANIFEST.tsv", "preamble.md", "purge.md", "record.md"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("out contains %v, want %v", got, want)
	}
}

func TestProvenanceHeaderIsFirstLineOfEveryMarkdownFile(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	writeFile(t, dir, "hints.md", demoHints())
	r := run(t, dir, "slice", "spec=spec.md", "hints=hints.md", "out=spec.d")
	requireCode(t, r, 0)

	specSum := fileSHA(t, filepath.Join(dir, "spec.md"))
	hintsSum := fileSHA(t, filepath.Join(dir, "hints.md"))
	re := regexp.MustCompile(`^<!-- pcd-slice source=spec\.md source-sha256=([0-9a-f]{64}) hints=hints\.md hints-sha256=([0-9a-f]{64}) version=(\S+) -->$`)

	for _, name := range dirEntries(t, filepath.Join(dir, "spec.d")) {
		body := readFile(t, filepath.Join(dir, "spec.d", name))
		first := strings.SplitN(body, "\n", 2)[0]
		if name == "MANIFEST.tsv" {
			if strings.HasPrefix(first, "<!--") {
				t.Fatalf("MANIFEST.tsv must carry no provenance line, got %q", first)
			}
			continue
		}
		m := re.FindStringSubmatch(first)
		if m == nil {
			t.Fatalf("%s first line is not a well-formed provenance line: %q", name, first)
		}
		if m[1] != specSum {
			t.Fatalf("%s: source-sha256 = %s, want %s", name, m[1], specSum)
		}
		if m[2] != hintsSum {
			t.Fatalf("%s: hints-sha256 = %s, want %s", name, m[2], hintsSum)
		}
		// determinism: no timestamp, host or user names in the output
		requireNotContains(t, name, body, os.Getenv("USER")+"@")
	}
}

func TestManifestListsEveryFileExceptItselfWithMatchingHashes(t *testing.T) {
	_, out := sliceDemo(t)
	manifest := readFile(t, filepath.Join(out, "MANIFEST.tsv"))
	lines := strings.Split(strings.TrimRight(manifest, "\n"), "\n")
	var names []string
	for _, l := range lines {
		parts := strings.Split(l, "\t")
		if len(parts) != 2 {
			t.Fatalf("manifest line is not name<TAB>sha256: %q", l)
		}
		if parts[0] == "MANIFEST.tsv" {
			t.Fatal("manifest must not list itself")
		}
		if got := fileSHA(t, filepath.Join(out, parts[0])); got != parts[1] {
			t.Fatalf("manifest hash for %s = %s, file hash = %s", parts[0], parts[1], got)
		}
		names = append(names, parts[0])
	}
	if !sort.StringsAreSorted(names) {
		t.Fatalf("manifest is not sorted by name: %v", names)
	}
	onDisk := dirEntries(t, out)
	if len(onDisk) != len(names)+1 {
		t.Fatalf("out contains %v but manifest lists %v", onDisk, names)
	}
}

func TestReSliceIntoSameDirectoryIsIdempotent(t *testing.T) {
	dir, out := sliceDemo(t)
	before := map[string]string{}
	for _, n := range dirEntries(t, out) {
		before[n] = fileSHA(t, filepath.Join(out, n))
	}
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	after := map[string]string{}
	for _, n := range dirEntries(t, out) {
		after[n] = fileSHA(t, filepath.Join(out, n))
	}
	if len(before) != len(after) {
		t.Fatalf("file set changed: %v -> %v", before, after)
	}
	for n, s := range before {
		if after[n] != s {
			t.Fatalf("file %s changed on identical re-run", n)
		}
	}
}

func TestStaleBundleIsRemovedOnReSlice(t *testing.T) {
	dir, out := sliceDemo(t)
	// remove one behavior from the spec: its bundle must disappear
	spec := strings.Replace(demoSpec(), `## BEHAVIOR: purge
Constraint: required

Removes everything that is expired.

STEPS:
1. Drop expired entries.

`, "", 1)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "slice", "spec=spec.md", "out=spec.d")
	requireCode(t, r, 0)
	for _, n := range dirEntries(t, out) {
		if n == "purge.md" {
			t.Fatal("stale bundle purge.md was not removed")
		}
	}
	requireNotContains(t, "MANIFEST.tsv", readFile(t, filepath.Join(out, "MANIFEST.tsv")), "purge.md")
}

// ---------------------------------------------------------------------------
// INVARIANT: sources are read-only; list and check never write files
// ---------------------------------------------------------------------------

func TestSourcesAreReadOnly(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	writeFile(t, dir, "hints.md", demoHints())
	specBefore := fileSHA(t, filepath.Join(dir, "spec.md"))
	hintsBefore := fileSHA(t, filepath.Join(dir, "hints.md"))

	for _, args := range [][]string{
		{"list", "spec=spec.md", "hints=hints.md"},
		{"check", "spec=spec.md", "hints=hints.md"},
		{"slice", "spec=spec.md", "hints=hints.md", "out=spec.d"},
	} {
		r := run(t, dir, args...)
		requireCode(t, r, 0)
	}
	if fileSHA(t, filepath.Join(dir, "spec.md")) != specBefore {
		t.Fatal("spec.md was modified")
	}
	if fileSHA(t, filepath.Join(dir, "hints.md")) != hintsBefore {
		t.Fatal("hints.md was modified")
	}
}

func TestListAndCheckNeverWriteFiles(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	before := dirEntries(t, dir)
	r1 := run(t, dir, "list", "spec=spec.md")
	requireCode(t, r1, 0)
	r2 := run(t, dir, "check", "spec=spec.md")
	requireCode(t, r2, 0)
	after := dirEntries(t, dir)
	if strings.Join(before, ",") != strings.Join(after, ",") {
		t.Fatalf("list/check created files: %v -> %v", before, after)
	}
}

func TestListExitsZeroEvenWithFindings(t *testing.T) {
	dir := sandbox(t)
	spec := md(`# tool

## TYPES

@@@
Store := string
@@@

## BEHAVIOR: reset
Constraint: required

Clears the Store.

## INVARIANTS
- X holds. (binds: renmae)
`)
	writeFile(t, dir, "spec.md", spec)
	r := run(t, dir, "list", "spec=spec.md")
	requireCode(t, r, 0)
	requireContains(t, "stdout", r.stdout, "total: 1 behaviors")
}

// ---------------------------------------------------------------------------
// invocation errors (exit 2)
// ---------------------------------------------------------------------------

func TestInvocationErrors(t *testing.T) {
	dir := sandbox(t)
	writeFile(t, dir, "spec.md", demoSpec())
	writeFile(t, dir, "notes.txt", "not markdown\n")

	cases := []struct {
		name string
		args []string
		want string
	}{
		{"no arguments", []string{}, "error:"},
		{"unknown verb", []string{"frobnicate", "spec=spec.md"}, "error:"},
		{"posix flag rejected", []string{"check", "--spec=spec.md"}, "error:"},
		{"posix short flag rejected", []string{"check", "-s", "spec.md"}, "error:"},
		{"bare positional rejected", []string{"check", "spec.md"}, "error:"},
		{"unknown key", []string{"check", "spec=spec.md", "depth=2"}, "error:"},
		{"missing spec", []string{"check"}, "error:"},
		{"slice without out", []string{"slice", "spec=spec.md"}, "error:"},
		{"missing file", []string{"check", "spec=missing.md"}, "error: cannot read missing.md"},
		{"not markdown", []string{"check", "spec=notes.txt"}, "error:"},
		{"missing hints file", []string{"check", "spec=spec.md", "hints=missing.md"}, "error: cannot read missing.md"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := run(t, dir, c.args...)
			requireCode(t, r, 2)
			requireContains(t, "stderr", r.stderr, c.want)
			if r.stdout != "" {
				t.Fatalf("invocation error wrote to stdout: %q", r.stdout)
			}
		})
	}
}

func TestHelpExitsZeroAndPrintsUsage(t *testing.T) {
	dir := sandbox(t)
	r := run(t, dir, "help")
	requireCode(t, r, 0)
	requireContains(t, "stdout", r.stdout, "pcd-slice")
	requireContains(t, "stdout", r.stdout, "spec=")
	requireContains(t, "stdout", r.stdout, "out=")
}

// ---------------------------------------------------------------------------
// bundle self-location: bundle + preamble suffice
// ---------------------------------------------------------------------------

func TestBundleCarriesItsBehaviorBlockVerbatim(t *testing.T) {
	_, out := sliceDemo(t)
	record := readFile(t, filepath.Join(out, "record.md"))
	for _, want := range []string{
		"## BEHAVIOR: record",
		"Constraint: required",
		"Records a Plan and returns.",
		"INPUTS:",
		"plan: Plan",
		"STEPS:",
		"1. Store the plan.",
	} {
		requireContains(t, "record.md", record, want)
	}
	requireNotContains(t, "record.md", record, "## BEHAVIOR: purge")
	// The behavior block appears in exactly one bundle.
	requireNotContains(t, "purge.md", readFile(t, filepath.Join(out, "purge.md")), "## BEHAVIOR: record")
	requireNotContains(t, "preamble.md", readFile(t, filepath.Join(out, "preamble.md")), "## BEHAVIOR: record")
}

func TestPreambleCarriesSharedContext(t *testing.T) {
	_, out := sliceDemo(t)
	preamble := readFile(t, filepath.Join(out, "preamble.md"))
	for _, want := range []string{
		"# demo-tool",
		"## META",
		"Intro prose for the demo tool.",
		"## TYPES",
		"Plan := { items: list of Item }",
		"Unused := int",
	} {
		requireContains(t, "preamble.md", preamble, want)
	}
}
