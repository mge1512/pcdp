// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
// tests by: claude-opus-5
//
// Black-box test harness for the pcd-slice CLI binary.
//
// The suite invokes the binary through the interface declared in the spec's
// DEPLOYMENT section (a command line tool taking a verb plus key=value
// options) and asserts on stdout, stderr and exit code only. It never
// imports the implementation packages.
//
// The binary is discovered at the canonical BINARY-LOCATION of the cli-tool
// deployment template: the project root, which is two directories up from
// independent_tests/<llm-name>/.
package pcdslice_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// specSHA256 is the SHA-256 of pcd-slice.spec.md, the specification these
// tests were derived from. The binary must report it in its version output.
const specSHA256 = "bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04"

// binaryRelPath is the canonical binary location for a cli-tool deployment,
// expressed relative to this test directory.
const binaryRelPath = "../../pcd-slice"

// binaryPath is binaryRelPath resolved to an absolute path. Tests run the
// binary with a per-test working directory, so the executable itself has to
// be addressed absolutely; the discovery path is still binaryRelPath.
var binaryPath string

func TestMain(m *testing.M) {
	abs, err := filepath.Abs(binaryRelPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot resolve %s: %v\n", binaryRelPath, err)
		os.Exit(1)
	}
	root := filepath.Dir(abs)

	args := []string{"build"}
	if v, err := os.ReadFile(filepath.Join(root, "VERSION")); err == nil {
		args = append(args, "-ldflags", "-X main.version="+strings.TrimSpace(string(v)))
	}
	args = append(args, "-o", abs, "./cmd/pcd-slice")

	build := exec.Command("go", args...)
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot build binary under test: %v\n%s\n", err, out)
		os.Exit(1)
	}
	binaryPath = abs
	os.Exit(m.Run())
}

// result holds one observed invocation of the binary.
type result struct {
	stdout string
	stderr string
	code   int
}

// sandbox returns a fresh per-test temporary directory. Every file the test
// or the binary touches lives underneath it; the directory is removed when
// the test finishes, including on the failure path.
func sandbox(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// run invokes the binary with dir as its working directory, so that path
// arguments in the test read exactly as they do in the spec's EXAMPLEs.
func run(t *testing.T, dir string, args ...string) result {
	t.Helper()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("cannot create sandbox HOME: %v", err)
	}
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir
	cmd.Env = []string{
		"HOME=" + home,
		"TMPDIR=" + dir,
		"PATH=/usr/bin:/bin",
	}
	var so, se bytes.Buffer
	cmd.Stdout = &so
	cmd.Stderr = &se
	err := cmd.Run()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatalf("running %v: %v", args, err)
	}
	return result{stdout: so.String(), stderr: se.String(), code: code}
}

// fixture expands the "@@@" marker to a markdown code fence. Go raw string
// literals cannot contain back quotes, so fixtures spell fences this way.
func fixture(s string) string {
	return strings.ReplaceAll(s, "@@@", "```")
}

// write materialises a fixture inside the sandbox. It refuses, loudly, to
// write anywhere outside it.
func write(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	rel, err := filepath.Rel(dir, p)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		t.Fatalf("refusing to write outside the sandbox: %s", p)
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatalf("cannot create fixture directory: %v", err)
	}
	if err := os.WriteFile(p, []byte(fixture(content)), 0o644); err != nil {
		t.Fatalf("cannot write fixture %s: %v", p, err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	return string(b)
}

// lines splits captured output into lines, dropping the trailing empty
// element produced by a final newline.
func lines(s string) []string {
	s = strings.TrimSuffix(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	return sha256hex(b)
}

// dirEntries returns the sorted names of the regular files in dir.
func dirEntries(t *testing.T, dir string) []string {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read directory %s: %v", dir, err)
	}
	var names []string
	for _, de := range des {
		names = append(names, de.Name())
	}
	sort.Strings(names)
	return names
}

// treeSnapshot maps every file under root (except the sandbox HOME) to its
// content hash, so a test can prove that a verb wrote nothing.
func treeSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	snap := map[string]string{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, rerr := filepath.Rel(root, p)
		if rerr != nil {
			return rerr
		}
		if rel == "home" || strings.HasPrefix(rel, "home"+string(filepath.Separator)) {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		snap[rel] = sha256hex(b)
		return nil
	})
	if err != nil {
		t.Fatalf("cannot snapshot %s: %v", root, err)
	}
	return snap
}

func sameSnapshot(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func mustContain(t *testing.T, what, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("%s does not contain %q\n--- %s ---\n%s", what, needle, what, haystack)
	}
}

func mustNotContain(t *testing.T, what, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Errorf("%s unexpectedly contains %q\n--- %s ---\n%s", what, needle, what, haystack)
	}
}

func mustExitCode(t *testing.T, r result, want int) {
	t.Helper()
	if r.code != want {
		t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s", r.code, want, r.stdout, r.stderr)
	}
}

// lineNumberOf returns the 1-based line number of the first fixture line
// containing needle. Findings name source line numbers; tests compute the
// expected number from the fixture rather than hard-coding it.
func lineNumberOf(t *testing.T, content, needle string) int {
	t.Helper()
	for i, ln := range strings.Split(fixture(content), "\n") {
		if strings.Contains(ln, needle) {
			return i + 1
		}
	}
	t.Fatalf("fixture does not contain %q", needle)
	return 0
}
