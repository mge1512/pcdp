// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Analysis is a specification and its hints files, parsed and cross-attributed.
type Analysis struct {
	Sources *Sources
	Spec    *Spec
	Hints   []*HintsDoc
}

// Analyze parses the specification, then every hints file, and attaches each
// attributed hints block to its behavior.
func Analyze(s *Sources) *Analysis {
	a := &Analysis{Sources: s, Spec: ParseSpec(s.Spec)}
	names := make([]string, 0, len(a.Spec.Behaviors))
	for _, b := range a.Spec.Behaviors {
		names = append(names, b.Name)
	}
	for _, h := range s.Hints {
		doc := ParseHints(h, names)
		a.Hints = append(a.Hints, doc)
		for _, blk := range doc.Blocks {
			for _, b := range a.Spec.Behaviors {
				if b.Name == blk.Behavior {
					b.HintsBlocks = append(b.HintsBlocks, blk)
					break
				}
			}
		}
	}
	return a
}

// Findings applies the rules of BEHAVIOR: check in the order the
// specification states them, and returns them sorted by file then line.
func (a *Analysis) Findings() []Finding {
	sp := a.Spec
	var out []Finding

	// 2. undefined-type
	out = append(out, sp.undefinedTypeFindings()...)

	// 3. unknown-binding
	for _, inv := range sp.Invariants {
		for _, name := range inv.Bindings {
			known := false
			for _, b := range sp.Behaviors {
				if b.Name == name {
					known = true
					break
				}
			}
			if !known {
				out = append(out, Finding{
					Rule:    "unknown-binding",
					File:    sp.Src.Path,
					Line:    inv.BindLine,
					Message: name,
				})
			}
		}
	}

	// 4. unassignable-example
	for _, ex := range sp.Examples {
		if ex.Behavior < 0 {
			out = append(out, Finding{
				Rule:    "unassignable-example",
				File:    sp.Src.Path,
				Line:    ex.Line,
				Message: ex.Name + ": attributable to no behavior",
			})
		}
	}

	// 5. unassignable-hints
	for _, doc := range a.Hints {
		for _, bad := range doc.BadHeadings {
			out = append(out, Finding{
				Rule:    "unassignable-hints",
				File:    doc.Src.Path,
				Line:    bad.Line,
				Message: bad.Text + ": names no known behavior",
			})
		}
	}

	// 6. duplicate-behavior, including a collision of bundle file names
	seenName := map[string]int{}
	seenFile := map[string]string{}
	for _, b := range sp.Behaviors {
		if first, dup := seenName[b.Name]; dup {
			out = append(out, Finding{
				Rule:    "duplicate-behavior",
				File:    sp.Src.Path,
				Line:    b.HeadingLine,
				Message: fmt.Sprintf("%s: duplicate behavior heading (first at line %d)", b.Name, first),
			})
			continue
		}
		seenName[b.Name] = b.HeadingLine
		if other, clash := seenFile[b.FileName]; clash {
			out = append(out, Finding{
				Rule:    "duplicate-behavior",
				File:    sp.Src.Path,
				Line:    b.HeadingLine,
				Message: fmt.Sprintf("%s: bundle file name %s collides with behavior %s", b.Name, b.FileName, other),
			})
			continue
		}
		seenFile[b.FileName] = b.Name
	}

	// 7. no-behaviors
	if len(sp.Behaviors) == 0 {
		out = append(out, Finding{
			Rule:    "no-behaviors",
			File:    sp.Src.Path,
			Line:    1,
			Message: "the specification contains no behavior heading",
		})
	}

	// 8. completeness
	for i, as := range sp.Assign {
		if as.Kind == AssignUnassigned {
			out = append(out, Finding{
				Rule:    "unassigned",
				File:    sp.Src.Path,
				Line:    i + 1,
				Message: "line is assigned to no output",
			})
		}
	}

	SortFindings(out)
	return dedupeFindings(out)
}

// dedupeFindings drops findings that are identical in rule, file, line and
// message; two references to the same undefined name on one line are one
// finding, not two.
func dedupeFindings(in []Finding) []Finding {
	seen := map[Finding]bool{}
	out := in[:0]
	for _, f := range in {
		if seen[f] {
			continue
		}
		seen[f] = true
		out = append(out, f)
	}
	return out
}

// StaleFindings compares the output directory with what the current sources
// would produce. It is used only when out= is given and the directory exists.
func (a *Analysis) StaleFindings(out string) []Finding {
	var findings []Finding
	outputs := a.BuildOutputs()

	expected := map[string][]byte{}
	for _, f := range outputs.All() {
		expected[f.Name] = f.Data
	}

	entries, err := os.ReadDir(out)
	if err != nil {
		return nil
	}
	present := map[string]bool{}
	for _, de := range entries {
		present[de.Name()] = true
		if _, want := expected[de.Name()]; want {
			continue
		}
		findings = append(findings, Finding{
			Rule:    "stale-output",
			File:    filepath.Join(out, de.Name()),
			Line:    0,
			Message: "unexpected file in out",
		})
	}

	for name, data := range expected {
		path := filepath.Join(out, name)
		if !present[name] {
			findings = append(findings, Finding{
				Rule:    "stale-output",
				File:    path,
				Line:    0,
				Message: "expected file is missing",
			})
			continue
		}
		got, err := os.ReadFile(path)
		if err != nil {
			findings = append(findings, Finding{
				Rule:    "stale-output",
				File:    path,
				Line:    0,
				Message: "cannot read emitted file",
			})
			continue
		}
		if name == manifestName {
			if string(got) != string(data) {
				findings = append(findings, Finding{
					Rule:    "stale-output",
					File:    path,
					Line:    1,
					Message: "manifest does not match the current sources",
				})
			}
			continue
		}
		if firstLine(got) != firstLine(data) {
			findings = append(findings, Finding{
				Rule:    "stale-output",
				File:    path,
				Line:    1,
				Message: "provenance header does not match the current sources",
			})
			continue
		}
		if string(got) != string(data) {
			findings = append(findings, Finding{
				Rule:    "stale-output",
				File:    path,
				Line:    1,
				Message: "content does not match the current sources",
			})
		}
	}

	SortFindings(findings)
	return findings
}

func firstLine(b []byte) string {
	s := string(b)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
