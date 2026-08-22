// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"regexp"
	"sort"
	"strings"
)

// undefinedTypeRE matches an identifier of type shape: two or more camel
// humps, each carrying at least one lower-case character, e.g. SpecFile or
// OutDir. The narrow shape keeps ordinary English words and the all-capital
// connectives of the refinement predicates (AND, OR) out of the
// undefined-type rule, which fires only on the definition side of the type
// graph.
var undefinedTypeRE = regexp.MustCompile(`\b([A-Z][a-z0-9_]+(?:[A-Z][a-z0-9_]+)+)\b`)

// typeMatcher compiles one whole-word alternation over the defined type
// names, longest first so that Plan does not shadow PlanPath.
func (sp *Spec) typeMatcher() *regexp.Regexp {
	if len(sp.TypeDefs) == 0 {
		return nil
	}
	names := make([]string, 0, len(sp.typeIndex))
	for n := range sp.typeIndex {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		if len(names[i]) != len(names[j]) {
			return len(names[i]) > len(names[j])
		}
		return names[i] < names[j]
	})
	quoted := make([]string, len(names))
	for i, n := range names {
		quoted[i] = regexp.QuoteMeta(n)
	}
	return regexp.MustCompile(`\b(` + strings.Join(quoted, "|") + `)\b`)
}

// computeClosures computes, per behavior, the transitive type closure: every
// type name whose whole-word form appears in the behavior block or in the
// definition of a type already in the closure.
func (sp *Spec) computeClosures() {
	matcher := sp.typeMatcher()
	for _, b := range sp.Behaviors {
		b.Closure = nil
		if matcher == nil {
			continue
		}
		seen := map[string]bool{}
		work := []string{strings.Join(sp.Src.Lines[b.Start:b.End], "\n")}
		for len(work) > 0 {
			text := work[0]
			work = work[1:]
			for _, m := range matcher.FindAllStringSubmatch(text, -1) {
				name := m[1]
				if seen[name] {
					continue
				}
				seen[name] = true
				if def, ok := sp.typeIndex[name]; ok {
					work = append(work, def.Text)
				}
			}
		}
		for name := range seen {
			b.Closure = append(b.Closure, name)
		}
		sort.Strings(b.Closure)
	}
}

// undefinedTypeFindings reports every type-shaped identifier that a type
// definition references and that is defined nowhere.
func (sp *Spec) undefinedTypeFindings() []Finding {
	var out []Finding
	for _, def := range sp.TypeDefs {
		for i := def.Start; i < def.End; i++ {
			line := stripComment(sp.Src.Lines[i])
			if i == def.Start {
				// Do not report the definition's own name.
				if idx := strings.Index(line, ":="); idx >= 0 {
					line = strings.Repeat(" ", idx) + line[idx:]
				}
			}
			for _, m := range undefinedTypeRE.FindAllStringSubmatch(line, -1) {
				name := m[1]
				if name == def.Name {
					continue
				}
				if _, ok := sp.typeIndex[name]; ok {
					continue
				}
				out = append(out, Finding{
					Rule:    "undefined-type",
					File:    sp.Src.Path,
					Line:    i + 1,
					Message: name + ": no definition (referenced by " + def.Name + ")",
				})
			}
		}
	}
	return out
}

// stripComment removes a "//" comment tail; comments carry prose, not type
// references.
func stripComment(line string) string {
	if i := strings.Index(line, "//"); i >= 0 {
		return line[:i]
	}
	return line
}
