// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// The check behavior: rule application, finding order, and staleness of an
// existing output directory.

package pcdslice

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CheckFindings applies every rule of the check behavior that does not need
// an output directory, and returns the findings sorted by file then line.
func CheckFindings(spec *Spec, hints []*Hints) []Finding {
	var fs []Finding
	fs = append(fs, spec.Findings...)
	for _, h := range hints {
		fs = append(fs, h.Findings...)
	}
	if len(spec.Behaviors) == 0 {
		fs = append(fs, Finding{
			Rule: "no-behaviors", File: spec.Path, Line: 1,
			Message: "specification contains no behavior heading",
		})
	}
	fs = append(fs, spec.undefinedTypeFindings()...)
	fs = append(fs, spec.unassignedFindings()...)
	SortFindings(fs)
	return fs
}

// SortFindings orders findings by file, then line, then rule, then message.
func SortFindings(fs []Finding) {
	sort.SliceStable(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		return a.Message < b.Message
	})
}

// StaleFindings compares an existing output directory against the freshly
// recomputed output. An absent out directory produces no findings.
func StaleFindings(out string, o *Output) ([]Finding, error) {
	info, err := os.Stat(out)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot read %s", out)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("out is not a directory: %s", out)
	}
	entries, err := os.ReadDir(out)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s", out)
	}

	var fs []Finding
	present := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		present[name] = true
		p := filepath.Join(out, name)
		if e.IsDir() {
			fs = append(fs, Finding{Rule: "stale-output", File: p, Line: 0,
				Message: "unexpected directory present"})
			continue
		}
		if name == ManifestFile {
			continue
		}
		if _, ok := o.Files[name]; !ok {
			fs = append(fs, Finding{Rule: "stale-output", File: p, Line: 0,
				Message: "unexpected file present"})
		}
	}

	manifestPresent := present[ManifestFile]
	manifestPath := filepath.Join(out, ManifestFile)
	if !manifestPresent {
		fs = append(fs, Finding{Rule: "stale-output", File: manifestPath, Line: 0,
			Message: "MANIFEST.tsv is missing; the output directory is stale in full"})
	} else {
		data, err := os.ReadFile(manifestPath)
		if err != nil || !bytes.Equal(data, o.Manifest) {
			fs = append(fs, Finding{Rule: "stale-output", File: manifestPath, Line: 1,
				Message: "manifest differs from the recomputed manifest"})
		}
	}

	for _, name := range o.Names {
		p := filepath.Join(out, name)
		if !present[name] {
			fs = append(fs, Finding{Rule: "stale-output", File: p, Line: 0,
				Message: "expected file missing"})
			continue
		}
		if !manifestPresent {
			fs = append(fs, Finding{Rule: "stale-output", File: p, Line: 1,
				Message: "output directory has no MANIFEST.tsv; stale in full"})
			continue
		}
		first, err := firstLineOf(p)
		if err != nil {
			fs = append(fs, Finding{Rule: "stale-output", File: p, Line: 1,
				Message: "cannot read emitted file"})
			continue
		}
		if first != o.Provenance {
			fs = append(fs, Finding{Rule: "stale-output", File: p, Line: 1,
				Message: "provenance header differs from the recomputed value"})
		}
	}
	SortFindings(fs)
	return fs, nil
}

// Check implements the check behavior end to end.
func Check(spec *Spec, hints []*Hints, out string, version string, stdout, stderr io.Writer) int {
	findings := CheckFindings(spec, hints)
	if out != "" {
		stale, err := StaleFindings(out, BuildOutput(spec, hints, version))
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return ExitInvocation
		}
		findings = append(findings, stale...)
		SortFindings(findings)
	}
	for _, f := range findings {
		fmt.Fprintln(stderr, f.String())
	}
	if len(findings) == 0 {
		fmt.Fprintln(stdout, "pcd-slice check: clean")
		return ExitOK
	}
	fmt.Fprintf(stdout, "pcd-slice check: %d findings\n", len(findings))
	return ExitFindings
}

func firstLineOf(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := string(data)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return s, nil
}
