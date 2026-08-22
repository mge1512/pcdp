// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// The list behavior: print the structure the other verbs act on.

package pcdslice

import (
	"fmt"
	"io"
)

// List prints one line per behavior in specification order - name, then the
// counts of types, invariants, examples and hints blocks - followed by the
// total line. It never fails on findings.
func List(spec *Spec, hints []*Hints, stdout io.Writer) int {
	for _, b := range spec.Behaviors {
		invariants := 0
		for _, inv := range spec.Invariants {
			if containsString(inv.Binds, b.Name) {
				invariants++
			}
		}
		examples := b.NestedExamples
		for _, e := range spec.Examples {
			if e.Target == b.Name {
				examples++
			}
		}
		blocks := 0
		for _, h := range hints {
			blocks += len(h.BlocksFor(b.Name))
		}
		fmt.Fprintf(stdout, "%s\t%d\t%d\t%d\t%d\n",
			b.Name, len(b.Closure), invariants, examples, blocks)
	}
	fmt.Fprintf(stdout, "total: %d behaviors\n", len(spec.Behaviors))
	return ExitOK
}
