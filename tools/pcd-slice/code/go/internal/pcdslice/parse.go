// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Section grammar and one-pass line classifier. Every structural decision is
// a prefix or suffix test on one line plus the walk state; no markdown
// library is used, because the spec's SECTION GRAMMAR is the whole grammar.

package pcdslice

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
)

type sectionKind int

const (
	secIntro sectionKind = iota
	secOther
	secTypes
	secInvariants
	secExamples
	secBehavior
	secChangelog
)

type section struct {
	kind sectionKind
	name string
	head int
	r    Range
}

type heading struct {
	line  int
	depth int
	text  string
}

// Spec is the parsed specification: its lines, the per-line assignment
// array, and the structures the verbs project from.
type Spec struct {
	Path       string
	SHA256     string
	Lines      []string
	Assign     []Assignment
	Behaviors  []Behavior
	Types      []TypeDef
	Invariants []Invariant
	Examples   []Example
	Findings   []Finding

	// Rewrite holds emission-time replacements for individual source lines
	// (an EXAMPLE heading whose "(of: ...)" suffix is stripped).
	Rewrite map[int]string

	byName          map[string]int
	typeByName      map[string]int
	inFence         []bool
	isFence         []bool
	sections        []section
	typeMatcher     *regexp.Regexp
	behaviorMatcher *regexp.Regexp
}

var (
	typeDefRe    = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9_]*)\s*:=`)
	changelogRe  = regexp.MustCompile(`(?i)\bchangelog\b`)
	exampleHeadR = regexp.MustCompile(`^EXAMPLE\b`)
)

// HashBytes returns the lowercase hex SHA-256 of b.
func HashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// SplitLines splits on "\n" and keeps any trailing "\r" attached to the
// line, so that emitted bytes match the bytes that came in. A final line
// without a newline counts as one line.
func SplitLines(data []byte) []string {
	s := string(data)
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	if n := len(lines); n > 0 && lines[n-1] == "" {
		lines = lines[:n-1]
	}
	return lines
}

// ParseSpec parses data as a specification located at path.
func ParseSpec(path string, data []byte) *Spec {
	s := &Spec{
		Path:       path,
		SHA256:     HashBytes(data),
		Lines:      SplitLines(data),
		Rewrite:    map[int]string{},
		byName:     map[string]int{},
		typeByName: map[string]int{},
	}
	s.Assign = make([]Assignment, len(s.Lines))
	s.markFences()
	s.buildSections()
	s.collectBehaviors()
	s.stampSections()
	s.collectTypes()
	s.buildMatchers()
	s.computeClosures()
	s.collectInvariants()
	s.collectExamples()
	return s
}

// BehaviorNames returns the behavior names in specification order.
func (s *Spec) BehaviorNames() []string {
	names := make([]string, 0, len(s.Behaviors))
	for _, b := range s.Behaviors {
		names = append(names, b.Name)
	}
	return names
}

// HasBehavior reports whether name is a known behavior.
func (s *Spec) HasBehavior(name string) bool {
	_, ok := s.byName[name]
	return ok
}

// Line returns the emitted text of source line i (the source line itself,
// unless an emission-time rewrite applies).
func (s *Spec) Line(i int) string {
	if r, ok := s.Rewrite[i]; ok {
		return r
	}
	return s.Lines[i]
}

// fenceArrays marks, for every line, whether it lies inside a code fence
// and whether it is a fence delimiter itself.
func fenceArrays(lines []string) (inFence, isFence []bool) {
	n := len(lines)
	inFence = make([]bool, n)
	isFence = make([]bool, n)
	open := false
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "```") {
			isFence[i] = true
			open = !open
			continue
		}
		inFence[i] = open
	}
	return inFence, isFence
}

// headingsOf returns every heading line: a line beginning with "#" at
// column 0 outside any code fence.
func headingsOf(lines []string, inFence, isFence []bool) []heading {
	var hs []heading
	for i, l := range lines {
		if inFence[i] || isFence[i] {
			continue
		}
		if !strings.HasPrefix(l, "#") {
			continue
		}
		d := 0
		for d < len(l) && l[d] == '#' {
			d++
		}
		hs = append(hs, heading{line: i, depth: d, text: strings.TrimSpace(l[d:])})
	}
	return hs
}

// trimTrailingBlank shrinks the half-open range [start,end) so that it does
// not end in blank lines.
func trimTrailingBlank(lines []string, start, end int) int {
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return end
}

func (s *Spec) markFences() {
	s.inFence, s.isFence = fenceArrays(s.Lines)
}

func (s *Spec) headings() []heading {
	return headingsOf(s.Lines, s.inFence, s.isFence)
}

func classifySection(h heading) section {
	sec := section{kind: secOther, head: h.line}
	if h.depth != 2 {
		return sec
	}
	t := h.text
	switch {
	case strings.HasPrefix(t, "BEHAVIOR/INTERNAL:"):
		name := strings.TrimSpace(strings.TrimPrefix(t, "BEHAVIOR/INTERNAL:"))
		if name != "" {
			sec.kind, sec.name = secBehavior, name
		}
	case strings.HasPrefix(t, "BEHAVIOR:"):
		name := strings.TrimSpace(strings.TrimPrefix(t, "BEHAVIOR:"))
		if name != "" {
			sec.kind, sec.name = secBehavior, name
		}
	case t == "TYPES":
		sec.kind = secTypes
	case t == "INVARIANTS":
		sec.kind = secInvariants
	case t == "EXAMPLES":
		sec.kind = secExamples
	case changelogRe.MatchString(t):
		sec.kind = secChangelog
	}
	return sec
}

// buildSections splits the source into sections. A "## " or "# " heading
// closes the section before it; "### " and deeper never do. A changelog
// section extends to end of file.
func (s *Spec) buildSections() {
	n := len(s.Lines)
	cur := section{kind: secIntro, head: -1, r: Range{0, n}}
	for _, h := range s.headings() {
		if cur.kind == secChangelog {
			break
		}
		if h.depth > 2 {
			continue
		}
		cur.r.End = h.line
		s.sections = append(s.sections, cur)
		cur = classifySection(h)
		cur.r = Range{h.line, n}
	}
	s.sections = append(s.sections, cur)
}

func (s *Spec) collectBehaviors() {
	for _, sec := range s.sections {
		if sec.kind != secBehavior {
			continue
		}
		if first, ok := s.byName[sec.name]; ok {
			s.Findings = append(s.Findings, Finding{
				Rule:    "duplicate-behavior",
				File:    s.Path,
				Line:    sec.head + 1,
				Message: sec.name + " is already defined at line " + strconv.Itoa(s.Behaviors[first].Heading+1),
			})
			continue
		}
		s.byName[sec.name] = len(s.Behaviors)
		s.Behaviors = append(s.Behaviors, Behavior{
			Name:     sec.name,
			FileName: strings.ReplaceAll(sec.name, "/", "-") + ".md",
			Heading:  sec.head,
			Block:    sec.r,
		})
	}
	// A file-name collision after "/" replacement, or a collision with the
	// preamble file, is a duplicate-behavior finding as well.
	seen := map[string]int{"preamble.md": -1, "MANIFEST.tsv": -1}
	for i, b := range s.Behaviors {
		if j, ok := seen[b.FileName]; ok {
			msg := b.Name + " maps to the reserved output file " + b.FileName
			if j >= 0 {
				msg = b.Name + " maps to the same output file as " + s.Behaviors[j].Name + " (" + b.FileName + ")"
			}
			s.Findings = append(s.Findings, Finding{
				Rule: "duplicate-behavior", File: s.Path, Line: b.Heading + 1, Message: msg,
			})
			continue
		}
		seen[b.FileName] = i
	}
}

func (s *Spec) stampSections() {
	for _, sec := range s.sections {
		for i := sec.r.Start; i < sec.r.End && i < len(s.Lines); i++ {
			switch sec.kind {
			case secChangelog:
				s.Assign[i] = Assignment{Kind: KindExcluded}
			case secBehavior:
				s.Assign[i] = Assignment{Kind: KindBehavior, Behavior: sec.name, Role: RoleBlock}
			default:
				s.Assign[i] = Assignment{Kind: KindPreamble}
			}
		}
	}
	// Count nested examples per behavior.
	for _, h := range s.headings() {
		if h.depth < 3 || !exampleHeadR.MatchString(h.text) {
			continue
		}
		a := s.Assign[h.line]
		if a.Kind != KindBehavior {
			continue
		}
		if idx, ok := s.byName[a.Behavior]; ok {
			s.Behaviors[idx].NestedExamples++
		}
	}
}

func (s *Spec) collectTypes() {
	for _, sec := range s.sections {
		if sec.kind != secTypes {
			continue
		}
		cur := -1
		for i := sec.r.Start; i < sec.r.End && i < len(s.Lines); i++ {
			if !s.inFence[i] {
				cur = -1
				continue
			}
			line := strings.TrimRight(s.Lines[i], "\r")
			if m := typeDefRe.FindStringSubmatch(line); m != nil {
				s.Types = append(s.Types, TypeDef{Name: m[1], Line: i, Body: Range{i, i + 1}})
				cur = len(s.Types) - 1
				continue
			}
			if cur >= 0 {
				s.Types[cur].Body.End = i + 1
			}
		}
	}
	for i := range s.Types {
		t := &s.Types[i]
		t.Text = strings.Join(s.Lines[t.Body.Start:t.Body.End], "\n")
		if _, ok := s.typeByName[t.Name]; !ok {
			s.typeByName[t.Name] = i
		}
	}
}

func (s *Spec) buildMatchers() {
	names := make([]string, 0, len(s.Types))
	for _, t := range s.Types {
		names = append(names, t.Name)
	}
	s.typeMatcher = wordAlternation(names)
	s.behaviorMatcher = wordAlternation(s.BehaviorNames())
}
