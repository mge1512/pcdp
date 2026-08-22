// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// The slice behavior: everything is computed first, then disk is touched.
// The foreign-entry refusal runs before the first write; on any write error
// the files this run created are removed again.

package pcdslice

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// provenanceMark is the prefix that identifies a file as pcd-slice output.
const provenanceMark = "<!-- pcd-slice "

var (
	writtenMu    sync.Mutex
	writtenFiles []string
)

func trackWritten(path string) {
	writtenMu.Lock()
	writtenFiles = append(writtenFiles, path)
	writtenMu.Unlock()
}

func forgetWritten() {
	writtenMu.Lock()
	writtenFiles = nil
	writtenMu.Unlock()
}

// Abort removes every file the current run has written. It is called from
// the signal handler so that a terminated run leaves no partial output.
func Abort() {
	writtenMu.Lock()
	files := writtenFiles
	writtenFiles = nil
	writtenMu.Unlock()
	for _, p := range files {
		_ = os.Remove(p)
	}
}

// Slice implements the slice behavior.
func Slice(spec *Spec, hints []*Hints, out, version string, stdout, stderr io.Writer) int {
	// Step 1: run check without out=; refuse on any finding.
	findings := CheckFindings(spec, hints)
	if len(findings) > 0 {
		for _, f := range findings {
			fmt.Fprintln(stderr, f.String())
		}
		return ExitFindings
	}

	// Step 2-4: compute all outputs in memory.
	o := BuildOutput(spec, hints, version)

	// Step 5: inspect the target directory before writing anything.
	entries, err := os.ReadDir(out)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(stderr, "error: cannot read %s\n", out)
		return ExitInvocation
	}
	for _, e := range entries {
		name := e.Name()
		p := filepath.Join(out, name)
		if e.IsDir() || !isSliceOutput(p, name) {
			fmt.Fprintf(stderr, "error: foreign file in out: %s\n", p)
			return ExitInvocation
		}
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		fmt.Fprintf(stderr, "error: cannot create out directory %s\n", out)
		return ExitInvocation
	}
	for _, e := range entries {
		name := e.Name()
		if name == ManifestFile {
			continue
		}
		if _, keep := o.Files[name]; keep {
			continue
		}
		if err := os.Remove(filepath.Join(out, name)); err != nil {
			fmt.Fprintf(stderr, "error: cannot remove stale file %s\n", filepath.Join(out, name))
			return ExitInvocation
		}
	}

	// Step 6: write every file, then the manifest.
	var written []string
	fail := func(path string) int {
		for _, p := range written {
			_ = os.Remove(p)
		}
		forgetWritten()
		fmt.Fprintf(stderr, "error: cannot write %s\n", path)
		return ExitInvocation
	}
	for _, name := range o.Names {
		p := filepath.Join(out, name)
		trackWritten(p)
		if err := os.WriteFile(p, o.Files[name], 0o644); err != nil {
			return fail(p)
		}
		written = append(written, p)
	}
	manifestPath := filepath.Join(out, ManifestFile)
	trackWritten(manifestPath)
	if err := os.WriteFile(manifestPath, o.Manifest, 0o644); err != nil {
		return fail(manifestPath)
	}
	forgetWritten()

	// Step 7.
	fmt.Fprintf(stdout, "pcd-slice: %d bundles, preamble, manifest -> %s\n", o.Bundles, out)
	return ExitOK
}

// isSliceOutput reports whether an existing entry in the out directory is
// pcd-slice's own: the manifest, or a file whose first line is a pcd-slice
// provenance comment.
func isSliceOutput(path, name string) bool {
	if name == ManifestFile {
		return true
	}
	first, err := firstLineOf(path)
	if err != nil {
		return false
	}
	return strings.HasPrefix(first, provenanceMark)
}
