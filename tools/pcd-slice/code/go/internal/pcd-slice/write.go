// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// provenancePrefix marks a file in the output directory as pcd-slice's own.
// A file without it is foreign and stops the run.
const provenancePrefix = "<!-- " + ToolName + " "

var (
	writtenMu sync.Mutex
	written   []string
)

func trackWritten(path string) {
	writtenMu.Lock()
	written = append(written, path)
	writtenMu.Unlock()
}

// RemovePartialOutput removes every file this run has written. It is called
// on a write error and by the signal handler, so that an interrupted run
// leaves no partial output behind.
func RemovePartialOutput() {
	writtenMu.Lock()
	paths := written
	written = nil
	writtenMu.Unlock()
	for _, p := range paths {
		_ = os.Remove(p)
	}
}

// Slice writes the derived form: one preamble, one bundle per behavior, one
// manifest.
func Slice(o Options) int {
	sources, err := LoadSources(o)
	if err != nil {
		return fail(o, err)
	}
	a := Analyze(sources)

	// 1. Run check without out=; refuse on any finding and write nothing.
	if findings := a.Findings(); len(findings) > 0 {
		reportFindings(o, findings)
		return int(ExitFindings)
	}

	// 2. Compute all outputs in memory.
	outputs := a.BuildOutputs()
	expected := map[string]bool{}
	for _, f := range outputs.All() {
		expected[f.Name] = true
	}

	// 5. Refuse on a foreign file; then remove our own files the manifest
	// would not list. Create the directory if it is absent.
	if info, err := os.Stat(o.Out); err == nil {
		if !info.IsDir() {
			return fail(o, fmt.Errorf("out is not a directory: %s", o.Out))
		}
		entries, err := os.ReadDir(o.Out)
		if err != nil {
			return fail(o, fmt.Errorf("cannot read out directory %s", o.Out))
		}
		var removable []string
		for _, de := range entries {
			path := filepath.Join(o.Out, de.Name())
			if de.IsDir() {
				return fail(o, fmt.Errorf("foreign file in out: %s", path))
			}
			if expected[de.Name()] {
				continue
			}
			if !hasProvenance(path) {
				return fail(o, fmt.Errorf("foreign file in out: %s", path))
			}
			removable = append(removable, path)
		}
		for _, path := range removable {
			if err := os.Remove(path); err != nil {
				return fail(o, fmt.Errorf("cannot remove stale file %s", path))
			}
		}
	} else if !os.IsNotExist(err) {
		return fail(o, fmt.Errorf("cannot read %s", o.Out))
	} else if err := os.MkdirAll(o.Out, 0o755); err != nil {
		return fail(o, fmt.Errorf("cannot create out directory %s", o.Out))
	}

	// 6. Write every file, then the manifest. On any write error the files
	// written in this run are removed again.
	for _, f := range outputs.All() {
		path := filepath.Join(o.Out, f.Name)
		if err := os.WriteFile(path, f.Data, 0o644); err != nil {
			RemovePartialOutput()
			return fail(o, fmt.Errorf("cannot write %s", path))
		}
		trackWritten(path)
	}
	forgetWritten()

	// 7. One summary line on stdout.
	fmt.Fprintf(o.Stdout, "%s: %d bundles, preamble, manifest -> %s\n", ToolName, outputs.Bundles, o.Out)
	return int(ExitClean)
}

func forgetWritten() {
	writtenMu.Lock()
	written = nil
	writtenMu.Unlock()
}

// hasProvenance reports whether the first line of a file is a pcd-slice
// provenance comment.
func hasProvenance(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()
	buf := make([]byte, 4096)
	n, _ := f.Read(buf)
	if n <= 0 {
		return false
	}
	first := string(buf[:n])
	if i := strings.IndexByte(first, '\n'); i >= 0 {
		first = first[:i]
	}
	return strings.HasPrefix(first, provenancePrefix)
}
