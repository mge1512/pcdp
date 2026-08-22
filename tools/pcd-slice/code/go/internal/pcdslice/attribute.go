// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Attribution: invariant bindings and the four-rung example ladder.

package pcdslice

import (
	"regexp"
	"sort"
	"strings"
)

var (
	// A binding suffix is recognized only at end of line, outside fences.
	bindsRe = regexp.MustCompile(`\(binds:\s*([^)]*)\)\s*$`)
	// Rung 2: the attribution suffix on an EXAMPLE heading.
	ofSuffixRe = regexp.MustCompile(`\(of:\s*([A-Za-z][A-Za-z0-9_/-]*)\)\s*$`)
	// Rung 3: the identifier invoked on the WHEN line.
	whenIdentRe = regexp.MustCompile(`^(?:result\s*=\s*)?([A-Za-z][A-Za-z0-9_/-]*)\s*\(`)
)

// listItems returns the list items of a section: a line whose trimmed form
// starts with "- " or "* " opens an item that extends to the next item or
// the end of the section, trailing blank lines excluded.
func (s *Spec) listItems(r Range) []Range {
	var items []Range
	start := -1
	flush := func(end int) {
		if start < 0 {
			return
		}
		items = append(items, Range{start, trimTrailingBlank(s.Lines, start, end)})
		start = -1
	}
	end := r.End
	if end > len(s.Lines) {
		end = len(s.Lines)
	}
	for i := r.Start; i < end; i++ {
		if s.inFence[i] || s.isFence[i] {
			continue
		}
		t := strings.TrimSpace(s.Lines[i])
		if strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ") {
			flush(i)
			start = i
		}
	}
	flush(end)
	return items
}

func (s *Spec) collectInvariants() {
	for _, sec := range s.sections {
		if sec.kind != secInvariants {
			continue
		}
		for _, item := range s.listItems(sec.r) {
			inv := Invariant{Item: item, BindLine: -1}
			for i := item.Start; i < item.End; i++ {
				if s.inFence[i] {
					continue
				}
				line := strings.TrimRight(s.Lines[i], " \t\r")
				m := bindsRe.FindStringSubmatch(line)
				if m == nil {
					continue
				}
				inv.BindLine = i
				for _, raw := range strings.Split(m[1], ",") {
					name := strings.TrimSpace(raw)
					if name == "" {
						continue
					}
					if !s.HasBehavior(name) {
						s.Findings = append(s.Findings, Finding{
							Rule: "unknown-binding", File: s.Path, Line: i + 1, Message: name,
						})
						continue
					}
					if !containsString(inv.Binds, name) {
						inv.Binds = append(inv.Binds, name)
					}
				}
			}
			if len(inv.Binds) > 0 {
				// Assigned once, to the first named behavior; emission
				// replicates it into each further named bundle.
				inv.Owner = inv.Binds[0]
				for i := item.Start; i < item.End; i++ {
					s.Assign[i] = Assignment{Kind: KindBehavior, Behavior: inv.Owner, Role: RoleInvariant}
				}
			}
			s.Invariants = append(s.Invariants, inv)
		}
	}
}

func (s *Spec) collectExamples() {
	hs := s.headings()
	for _, sec := range s.sections {
		if sec.kind != secExamples {
			continue
		}
		var deep []heading
		for _, h := range hs {
			if h.line >= sec.r.Start && h.line < sec.r.End && h.depth >= 3 {
				deep = append(deep, h)
			}
		}
		for k, h := range deep {
			if !exampleHeadR.MatchString(h.text) {
				continue
			}
			end := sec.r.End
			if k+1 < len(deep) {
				end = deep[k+1].line
			}
			if end > len(s.Lines) {
				end = len(s.Lines)
			}
			ex := Example{
				Heading: h.line,
				Block:   Range{h.line, trimTrailingBlank(s.Lines, h.line+1, end)},
			}
			s.attributeExample(&ex, h)
			s.Examples = append(s.Examples, ex)
		}
	}
}

// attributeExample runs the attribution ladder; the first rung that fires
// wins. Nesting (rung 1) is already handled by the section walk: a nested
// example never reaches this function.
func (s *Spec) attributeExample(ex *Example, h heading) {
	raw := strings.TrimRight(s.Lines[h.line], "\r")
	text := h.text
	ex.Name = exampleName(text)

	// Rung 2: the heading suffix.
	if loc := ofSuffixRe.FindStringSubmatchIndex(raw); loc != nil {
		arg := raw[loc[2]:loc[3]]
		if arg == "shared" || s.HasBehavior(arg) {
			s.Rewrite[h.line] = strings.TrimRight(raw[:loc[0]], " \t")
			ex.Name = exampleName(strings.TrimSpace(strings.TrimLeft(s.Rewrite[h.line], "#")))
			if arg == "shared" {
				ex.Shared = true // belongs to the preamble
				return
			}
			s.assignExample(ex, arg)
			return
		}
	}

	// Rung 3: the identifier invoked on the WHEN line.
	if name := s.whenIdentifier(ex); name != "" {
		s.assignExample(ex, name)
		return
	}

	// Rung 4: a behavior name appearing as a whole word in the example text.
	names := s.behaviorNamesIn(ex.Block.Start+1, ex.Block.End)
	switch len(names) {
	case 1:
		s.assignExample(ex, names[0])
	case 0:
		s.Findings = append(s.Findings, Finding{
			Rule: "unassignable-example", File: s.Path, Line: h.line + 1,
			Message: "example " + ex.Name + " cannot be attributed to any behavior",
		})
	default:
		s.Findings = append(s.Findings, Finding{
			Rule: "unassignable-example", File: s.Path, Line: h.line + 1,
			Message: "example " + ex.Name + " matches several behaviors: " + strings.Join(names, ", "),
		})
	}
}

func (s *Spec) assignExample(ex *Example, name string) {
	ex.Target = name
	for i := ex.Block.Start; i < ex.Block.End; i++ {
		s.Assign[i] = Assignment{Kind: KindBehavior, Behavior: name, Role: RoleExample}
	}
}

func (s *Spec) whenIdentifier(ex *Example) string {
	inWhen := false
	for i := ex.Block.Start + 1; i < ex.Block.End; i++ {
		if s.isFence[i] {
			continue
		}
		t := strings.TrimSpace(s.Lines[i])
		up := strings.ToUpper(t)
		if strings.HasPrefix(up, "WHEN:") {
			inWhen = true
			t = strings.TrimSpace(t[len("WHEN:"):])
		} else if strings.HasPrefix(up, "THEN:") || strings.HasPrefix(up, "GIVEN:") {
			inWhen = false
			continue
		}
		if !inWhen || t == "" {
			continue
		}
		if m := whenIdentRe.FindStringSubmatch(t); m != nil && s.HasBehavior(m[1]) {
			return m[1]
		}
	}
	return ""
}

// behaviorNamesIn returns the distinct behavior names occurring as whole
// words in the given line range, alphabetically sorted.
func (s *Spec) behaviorNamesIn(start, end int) []string {
	if s.behaviorMatcher == nil {
		return nil
	}
	if end > len(s.Lines) {
		end = len(s.Lines)
	}
	seen := map[string]bool{}
	for i := start; i < end; i++ {
		for _, m := range s.behaviorMatcher.FindAllString(s.Lines[i], -1) {
			seen[m] = true
		}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// unassignedFindings implements the completeness rule of check step 8.
func (s *Spec) unassignedFindings() []Finding {
	var out []Finding
	for i, a := range s.Assign {
		if a.Kind == KindUnassigned {
			out = append(out, Finding{
				Rule: "unassigned", File: s.Path, Line: i + 1,
				Message: "line received no assignment",
			})
		}
	}
	return out
}

func exampleName(text string) string {
	t := strings.TrimSpace(strings.TrimPrefix(text, "EXAMPLE"))
	t = strings.TrimSpace(strings.TrimPrefix(t, ":"))
	if loc := ofSuffixRe.FindStringIndex(t); loc != nil {
		t = strings.TrimSpace(t[:loc[0]])
	}
	return t
}

func containsString(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
