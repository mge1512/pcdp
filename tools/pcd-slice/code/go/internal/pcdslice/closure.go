// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Type definitions, the transitive type closure, and the undefined-type
// rule. The closure is a worklist over the behavior block text and then
// over each newly added definition's text; alphabetisation happens only at
// output time.

package pcdslice

import (
	"regexp"
	"sort"
	"strings"
)

// wordAlternation compiles one whole-word alternation over names, longest
// first so that a longer name is never shadowed by a prefix of it. It
// returns nil for an empty name set.
func wordAlternation(names []string) *regexp.Regexp {
	if len(names) == 0 {
		return nil
	}
	sorted := append([]string(nil), names...)
	sort.Slice(sorted, func(i, j int) bool {
		if len(sorted[i]) != len(sorted[j]) {
			return len(sorted[i]) > len(sorted[j])
		}
		return sorted[i] < sorted[j]
	})
	quoted := make([]string, 0, len(sorted))
	for _, n := range sorted {
		quoted = append(quoted, regexp.QuoteMeta(n))
	}
	return regexp.MustCompile(`\b(?:` + strings.Join(quoted, "|") + `)\b`)
}

func (s *Spec) computeClosures() {
	for i := range s.Behaviors {
		b := &s.Behaviors[i]
		b.Closure = s.closureOf(strings.Join(s.Lines[b.Block.Start:b.Block.End], "\n"))
	}
}

// closureOf returns the transitive type closure of a text, alphabetically
// sorted.
func (s *Spec) closureOf(text string) []string {
	if s.typeMatcher == nil {
		return nil
	}
	found := map[string]bool{}
	queue := []string{text}
	for len(queue) > 0 {
		t := queue[0]
		queue = queue[1:]
		for _, m := range s.typeMatcher.FindAllString(t, -1) {
			if found[m] {
				continue
			}
			found[m] = true
			if idx, ok := s.typeByName[m]; ok {
				queue = append(queue, s.Types[idx].Text)
			}
		}
	}
	names := make([]string, 0, len(found))
	for n := range found {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

var tokenSplitRe = regexp.MustCompile(`[^A-Za-z0-9_]+`)

// isReferenceShape implements the normative reference shape of the
// undefined-type rule: an identifier with two or more camel humps, each
// containing at least one lowercase letter.
func isReferenceShape(tok string) bool {
	if tok == "" || strings.ContainsRune(tok, '_') {
		return false
	}
	if tok[0] < 'A' || tok[0] > 'Z' {
		return false
	}
	humps := 0
	lower := false
	for i := 0; i < len(tok); i++ {
		c := tok[i]
		switch {
		case c >= 'A' && c <= 'Z':
			if humps > 0 && !lower {
				return false // previous hump carried no lowercase letter
			}
			humps++
			lower = false
		case c >= 'a' && c <= 'z':
			lower = true
		case c >= '0' && c <= '9':
			// digits belong to the current hump
		default:
			return false
		}
	}
	return humps >= 2 && lower
}

// stripLineComment removes a "//" comment from a line.
func stripLineComment(line string) string {
	if i := strings.Index(line, "//"); i >= 0 {
		return line[:i]
	}
	return line
}

// undefinedTypeFindings implements check step 2. It fires only inside type
// definitions; behavior text is never scanned.
func (s *Spec) undefinedTypeFindings() []Finding {
	var out []Finding
	for _, def := range s.Types {
		for i := def.Body.Start; i < def.Body.End; i++ {
			line := stripLineComment(strings.TrimRight(s.Lines[i], "\r"))
			if i == def.Line {
				// exclude the defining name itself
				if j := strings.Index(line, ":="); j >= 0 {
					line = strings.Repeat(" ", j) + line[j:]
				}
			}
			for _, tok := range tokenSplitRe.Split(line, -1) {
				if tok == def.Name || !isReferenceShape(tok) {
					continue
				}
				if _, ok := s.typeByName[tok]; ok {
					continue
				}
				out = append(out, Finding{
					Rule:    "undefined-type",
					File:    s.Path,
					Line:    i + 1,
					Message: "undefined type " + tok + " referenced in definition of " + def.Name,
				})
			}
		}
	}
	return out
}
