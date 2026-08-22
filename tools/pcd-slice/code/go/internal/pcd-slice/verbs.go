// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"fmt"
	"os"
)

// List prints the structure the other verbs will act on: each behavior with
// its type closure, bound invariants, attributed examples and attributed
// hints blocks. It never fails on findings.
func List(o Options) int {
	sources, err := LoadSources(o)
	if err != nil {
		return fail(o, err)
	}
	a := Analyze(sources)
	for _, b := range a.Spec.Behaviors {
		fmt.Fprintf(o.Stdout, "%s\t%d\t%d\t%d\t%d\n",
			b.Name, len(b.Closure), len(b.Invariants), len(b.Examples), len(b.HintsBlocks))
	}
	fmt.Fprintf(o.Stdout, "total: %d behaviors\n", len(a.Spec.Behaviors))
	return int(ExitClean)
}

// Check validates that the specification and hints can be sliced completely
// and unambiguously and, when out= is given and exists, that the output is
// current.
func Check(o Options) int {
	sources, err := LoadSources(o)
	if err != nil {
		return fail(o, err)
	}
	a := Analyze(sources)
	findings := a.Findings()
	if o.Out != "" {
		if info, err := os.Stat(o.Out); err == nil && info.IsDir() {
			findings = append(findings, a.StaleFindings(o.Out)...)
			SortFindings(findings)
		}
	}
	reportFindings(o, findings)
	if len(findings) > 0 {
		fmt.Fprintf(o.Stdout, "%s check: %d findings\n", ToolName, len(findings))
		return int(ExitFindings)
	}
	fmt.Fprintf(o.Stdout, "%s check: clean\n", ToolName)
	return int(ExitClean)
}

func reportFindings(o Options, findings []Finding) {
	for _, f := range findings {
		fmt.Fprintln(o.Stderr, f.String())
	}
}

// fail reports an invocation error on stderr and returns the invocation exit
// code. Every such message starts with "error: " and names the path.
func fail(o Options, err error) int {
	fmt.Fprintf(o.Stderr, "error: %v\n", err)
	return int(ExitInvocation)
}
