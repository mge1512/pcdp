// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"regexp"
	"strings"
)

// The whole recognized grammar of a specification, as regular expressions on
// one line plus the fence/section state the walk carries. No markdown parser
// is used: a parser would recognize more than the grammar admits.
var (
	behaviorHeadingRE = regexp.MustCompile(`^BEHAVIOR(?:/[A-Za-z0-9_-]+)?:\s*(\S.*?)\s*$`)
	typeDefRE         = regexp.MustCompile(`^\s*([A-Za-z][A-Za-z0-9_]*)\s*:=`)
	exampleHeadingRE  = regexp.MustCompile(`^(#{3,})\s*EXAMPLE:?\s+(\S.*?)\s*$`)
	bindsRE           = regexp.MustCompile(`\(binds:\s*([^)]*)\)\s*$`)
	invokedRE         = regexp.MustCompile(`^\s*(?:result\s*=\s*)?([a-z][a-z0-9_-]*)\s*\(`)
	listItemRE        = regexp.MustCompile(`^\s*[-*]\s+\S`)
	stanzaLabelRE     = regexp.MustCompile(`^\s*[A-Z][A-Z0-9_ -]*:`)
	changelogRE       = regexp.MustCompile(`(?i)(^|[^A-Za-z0-9])changelog([^A-Za-z0-9]|$)`)
)

// TypeDef is one definition inside the TYPES section. A definition extends
// until the next definition line or the end of the fenced block; comment and
// continuation lines belong to the definition above them.
type TypeDef struct {
	Name       string
	Line       int // 1-based line of the definition line
	Start, End int // half-open range of source line indices
	Text       string
}

// Invariant is one list item of the INVARIANTS section.
type Invariant struct {
	Start, End int      // half-open range of source line indices
	BindLine   int      // 1-based line carrying the "(binds: ...)" suffix, 0 if none
	Bindings   []string // names as written, in order
	Targets    []int    // indices of the behaviors it is bound to
}

// Global reports whether the invariant carries no binding suffix and so
// belongs to the preamble.
func (i *Invariant) Global() bool { return i.BindLine == 0 }

// Example is one "### EXAMPLE" block of the top-level EXAMPLES section.
type Example struct {
	Name       string
	Line       int // 1-based heading line
	Start, End int // half-open range of source line indices
	Behavior   int // index of the attributed behavior, -1 when unattributable
}

// Behavior is one behavior block plus everything attributed to it.
type Behavior struct {
	Name        string
	FileName    string // bundle file name: the name with "/" replaced by "-"
	HeadingLine int    // 1-based
	Start, End  int    // half-open range of source line indices
	Closure     []string
	Invariants  []*Invariant
	Examples    []*Example
	HintsBlocks []*HintsBlock
}

// Spec is a parsed specification: the source, the per-line assignment, and
// the structures the grammar recognizes.
type Spec struct {
	Src        *Source
	Assign     []Assignment
	Behaviors  []*Behavior
	TypeDefs   []*TypeDef
	typeIndex  map[string]*TypeDef
	Invariants []*Invariant
	Examples   []*Example

	fenced   []bool
	sections []section
}

type section struct {
	kind       string // intro, types, behavior, invariants, examples, changelog, other
	name       string
	start, end int // half-open range of source line indices; start is the heading
}

// fenceMap marks, for every line, whether it lies inside a fenced block or is
// a fence delimiter itself. Fenced content is never interpreted as structure.
func fenceMap(lines []string) (fenced []bool, delim []bool) {
	fenced = make([]bool, len(lines))
	delim = make([]bool, len(lines))
	in := false
	for i, ln := range lines {
		if strings.HasPrefix(strings.TrimSpace(ln), "```") {
			delim[i] = true
			fenced[i] = true
			in = !in
			continue
		}
		fenced[i] = in
	}
	return fenced, delim
}

// headingDepth returns the number of leading "#" of a heading line at column
// 0, or 0 when the line is not a heading.
func headingDepth(line string) int {
	n := 0
	for n < len(line) && line[n] == '#' {
		n++
	}
	if n == 0 {
		return 0
	}
	return n
}

func headingText(line string, depth int) string {
	return strings.TrimSpace(line[depth:])
}

// ParseSpec walks the source once and produces the parsed specification with
// every line assigned.
func ParseSpec(src *Source) *Spec {
	sp := &Spec{Src: src, typeIndex: map[string]*TypeDef{}}
	sp.Assign = make([]Assignment, len(src.Lines))
	for i := range sp.Assign {
		sp.Assign[i] = Assignment{Kind: AssignUnassigned, Behavior: -1}
	}
	fenced, delim := fenceMap(src.Lines)
	sp.fenced = fenced

	sp.buildSections()
	sp.collectBehaviors()
	sp.collectTypeDefs(delim)
	sp.collectInvariants()
	sp.collectExamples()
	sp.assignLines()
	sp.computeClosures()
	return sp
}

// buildSections splits the source into "## " level sections. A changelog
// section extends to the end of the file.
func (sp *Spec) buildSections() {
	lines := sp.Src.Lines
	var secs []section
	cur := section{kind: "intro", start: 0}
	for i, ln := range lines {
		if sp.fenced[i] {
			continue
		}
		d := headingDepth(ln)
		if d == 0 || d > 2 {
			continue
		}
		text := headingText(ln, d)
		kind, name := classifyHeading(text)
		if i > 0 || cur.kind != "intro" {
			cur.end = i
			if cur.end > cur.start || cur.kind != "intro" {
				secs = append(secs, cur)
			}
		}
		cur = section{kind: kind, name: name, start: i}
		if kind == "changelog" {
			break
		}
	}
	cur.end = len(lines)
	if cur.end > cur.start {
		secs = append(secs, cur)
	}
	sp.sections = secs
}

func classifyHeading(text string) (kind, name string) {
	switch strings.TrimSpace(text) {
	case "TYPES":
		return "types", ""
	case "INVARIANTS":
		return "invariants", ""
	case "EXAMPLES":
		return "examples", ""
	}
	if m := behaviorHeadingRE.FindStringSubmatch(text); m != nil {
		return "behavior", strings.TrimSpace(m[1])
	}
	if changelogRE.MatchString(text) {
		return "changelog", ""
	}
	return "other", ""
}

func (sp *Spec) collectBehaviors() {
	for _, sec := range sp.sections {
		if sec.kind != "behavior" {
			continue
		}
		sp.Behaviors = append(sp.Behaviors, &Behavior{
			Name:        sec.name,
			FileName:    strings.ReplaceAll(sec.name, "/", "-") + ".md",
			HeadingLine: sec.start + 1,
			Start:       sec.start,
			End:         sec.end,
		})
	}
}

// collectTypeDefs reads the definition lines of every fenced block inside the
// TYPES section.
func (sp *Spec) collectTypeDefs(delim []bool) {
	lines := sp.Src.Lines
	for _, sec := range sp.sections {
		if sec.kind != "types" {
			continue
		}
		var cur *TypeDef
		for i := sec.start; i < sec.end; i++ {
			if delim[i] {
				if cur != nil {
					cur.End = i
					cur = nil
				}
				continue
			}
			if !sp.fenced[i] {
				continue
			}
			if m := typeDefRE.FindStringSubmatch(lines[i]); m != nil {
				if cur != nil {
					cur.End = i
				}
				cur = &TypeDef{Name: m[1], Line: i + 1, Start: i, End: i + 1}
				sp.TypeDefs = append(sp.TypeDefs, cur)
				continue
			}
			if cur != nil {
				cur.End = i + 1
			}
		}
		if cur != nil && cur.End > sec.end {
			cur.End = sec.end
		}
	}
	for _, d := range sp.TypeDefs {
		d.Text = strings.Join(lines[d.Start:d.End], "\n")
		if _, seen := sp.typeIndex[d.Name]; !seen {
			sp.typeIndex[d.Name] = d
		}
	}
}

// collectInvariants groups the INVARIANTS section into list items and reads
// the binding suffix of each.
func (sp *Spec) collectInvariants() {
	lines := sp.Src.Lines
	for _, sec := range sp.sections {
		if sec.kind != "invariants" {
			continue
		}
		var cur *Invariant
		closeItem := func(end int) {
			if cur == nil {
				return
			}
			for end > cur.Start && strings.TrimSpace(lines[end-1]) == "" {
				end--
			}
			cur.End = end
			sp.Invariants = append(sp.Invariants, cur)
			cur = nil
		}
		for i := sec.start + 1; i < sec.end; i++ {
			if sp.fenced[i] {
				continue
			}
			if listItemRE.MatchString(lines[i]) {
				closeItem(i)
				cur = &Invariant{Start: i, End: i + 1}
			}
		}
		closeItem(sec.end)
	}
	for _, inv := range sp.Invariants {
		for i := inv.End - 1; i >= inv.Start; i-- {
			ln := strings.TrimRight(lines[i], "\r")
			if strings.TrimSpace(ln) == "" {
				continue
			}
			if m := bindsRE.FindStringSubmatch(ln); m != nil {
				inv.BindLine = i + 1
				for _, raw := range strings.Split(m[1], ",") {
					if name := strings.TrimSpace(raw); name != "" {
						inv.Bindings = append(inv.Bindings, name)
					}
				}
			}
			break
		}
	}
}

// collectExamples reads the top-level EXAMPLES section and attributes each
// example to a behavior.
func (sp *Spec) collectExamples() {
	lines := sp.Src.Lines
	for _, sec := range sp.sections {
		if sec.kind != "examples" {
			continue
		}
		var cur *Example
		closeItem := func(end int) {
			if cur == nil {
				return
			}
			for end > cur.Start && strings.TrimSpace(lines[end-1]) == "" {
				end--
			}
			cur.End = end
			sp.Examples = append(sp.Examples, cur)
			cur = nil
		}
		for i := sec.start + 1; i < sec.end; i++ {
			if sp.fenced[i] {
				continue
			}
			if m := exampleHeadingRE.FindStringSubmatch(lines[i]); m != nil {
				closeItem(i)
				cur = &Example{Name: m[2], Line: i + 1, Start: i, End: i + 1, Behavior: -1}
				continue
			}
			if headingDepth(lines[i]) > 0 && headingDepth(lines[i]) <= 3 {
				closeItem(i)
			}
		}
		closeItem(sec.end)
	}
	for _, ex := range sp.Examples {
		ex.Behavior = sp.attributeExample(ex)
	}
}

// attributeExample applies the two attribution rules of the grammar: the
// invoked identifier on the WHEN line first, the example's own name second.
func (sp *Spec) attributeExample(ex *Example) int {
	lines := sp.Src.Lines
	for i := ex.Start; i < ex.End; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(trimmed, "WHEN:") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "WHEN:"))
		if idx := sp.behaviorOfInvocation(rest); idx >= 0 {
			return idx
		}
		for j := i + 1; j < ex.End; j++ {
			cand := lines[j]
			if stanzaLabelRE.MatchString(cand) && !strings.HasPrefix(strings.TrimSpace(cand), "WHEN:") {
				break
			}
			if idx := sp.behaviorOfInvocation(cand); idx >= 0 {
				return idx
			}
		}
		break
	}
	for i, b := range sp.Behaviors {
		if containsWholeToken(ex.Name, b.Name) {
			return i
		}
	}
	return -1
}

func (sp *Spec) behaviorOfInvocation(line string) int {
	m := invokedRE.FindStringSubmatch(line)
	if m == nil {
		return -1
	}
	for i, b := range sp.Behaviors {
		if b.Name == m[1] {
			return i
		}
	}
	return -1
}

// containsWholeToken reports whether name occurs in text delimited by
// characters that are neither letters nor digits. Underscores and hyphens
// separate tokens here, because example names are written in snake case.
func containsWholeToken(text, name string) bool {
	if name == "" {
		return false
	}
	for off := 0; ; {
		i := strings.Index(text[off:], name)
		if i < 0 {
			return false
		}
		i += off
		beforeOK := i == 0 || !isAlnum(rune(text[i-1]))
		end := i + len(name)
		afterOK := end == len(text) || !isAlnum(rune(text[end]))
		if beforeOK && afterOK {
			return true
		}
		off = i + 1
		if off >= len(text) {
			return false
		}
	}
}

func isAlnum(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// assignLines stamps every source line with exactly one assignment.
func (sp *Spec) assignLines() {
	set := func(from, to int, a Assignment) {
		for i := from; i < to && i < len(sp.Assign); i++ {
			sp.Assign[i] = a
		}
	}
	for _, sec := range sp.sections {
		switch sec.kind {
		case "changelog":
			set(sec.start, sec.end, Assignment{Kind: AssignExcluded, Behavior: -1})
		case "behavior":
			idx := sp.behaviorIndexAt(sec.start)
			set(sec.start, sec.end, Assignment{Kind: AssignBehavior, Behavior: idx})
		default:
			set(sec.start, sec.end, Assignment{Kind: AssignPreamble, Behavior: -1})
		}
	}
	// A bound invariant travels with the behaviors it names. For the
	// completeness count it is stamped onto the first of them; the bundle
	// projection appends it to every named bundle.
	for _, inv := range sp.Invariants {
		inv.Targets = nil
		for _, name := range inv.Bindings {
			for i, b := range sp.Behaviors {
				if b.Name == name {
					inv.Targets = append(inv.Targets, i)
					break
				}
			}
		}
		for _, t := range inv.Targets {
			sp.Behaviors[t].Invariants = append(sp.Behaviors[t].Invariants, inv)
		}
		if len(inv.Targets) > 0 {
			set(inv.Start, inv.End, Assignment{Kind: AssignBehavior, Behavior: inv.Targets[0]})
		}
	}
	for _, ex := range sp.Examples {
		if ex.Behavior < 0 {
			continue
		}
		sp.Behaviors[ex.Behavior].Examples = append(sp.Behaviors[ex.Behavior].Examples, ex)
		set(ex.Start, ex.End, Assignment{Kind: AssignBehavior, Behavior: ex.Behavior})
	}
}

func (sp *Spec) behaviorIndexAt(start int) int {
	for i, b := range sp.Behaviors {
		if b.Start == start {
			return i
		}
	}
	return -1
}
