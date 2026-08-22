// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"regexp"
	"strings"
)

// behaviorShapeRE matches a heading whose stripped text is one single
// lowercase identifier. That is the cheap "behavior shape" test the
// unassignable-hints finding needs; anything looser drowns the finding in
// ordinary prose headings.
var behaviorShapeRE = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// HintsBlock is a region of a hints file attributed to one behavior.
type HintsBlock struct {
	Heading    string
	Line       int // 1-based heading line
	Start, End int // half-open range of source line indices
	Behavior   string
	Doc        *HintsDoc
}

// BadHeading is a hints heading of behavior shape that names no known
// behavior.
type BadHeading struct {
	Line int
	Text string
}

// HintsDoc is a parsed hints file: the blocks attributed to behaviors, and
// everything else, which is hints preamble.
type HintsDoc struct {
	Src           *Source
	Blocks        []*HintsBlock
	PreambleLines []int
	BadHeadings   []BadHeading
}

// ParseHints attributes the headings of a hints file to behaviors. A "## " or
// deeper heading whose text contains a behavior name as a whole word opens a
// block for that behavior, extending to the next heading of the same or
// shallower depth. Everything else is hints preamble.
func ParseHints(src *Source, behaviors []string) *HintsDoc {
	doc := &HintsDoc{Src: src}
	fenced, _ := fenceMap(src.Lines)

	type heading struct {
		idx, depth int
		text       string
	}
	var hs []heading
	for i, ln := range src.Lines {
		if fenced[i] {
			continue
		}
		if d := headingDepth(ln); d > 0 {
			hs = append(hs, heading{idx: i, depth: d, text: headingText(ln, d)})
		}
	}

	assigned := make([]bool, len(src.Lines))
	for k, h := range hs {
		if h.depth < 2 || assigned[h.idx] {
			continue
		}
		stripped := stripEmphasis(h.text)
		name := matchBehaviorInHeading(stripped, behaviors)
		if name == "" {
			if behaviorShapeRE.MatchString(strings.TrimSpace(stripped)) {
				doc.BadHeadings = append(doc.BadHeadings, BadHeading{Line: h.idx + 1, Text: strings.TrimSpace(stripped)})
			}
			continue
		}
		end := len(src.Lines)
		for j := k + 1; j < len(hs); j++ {
			if hs[j].depth <= h.depth {
				end = hs[j].idx
				break
			}
		}
		for i := h.idx; i < end; i++ {
			assigned[i] = true
		}
		doc.Blocks = append(doc.Blocks, &HintsBlock{
			Heading:  strings.TrimSpace(h.text),
			Line:     h.idx + 1,
			Start:    h.idx,
			End:      end,
			Behavior: name,
			Doc:      doc,
		})
	}
	for i := range src.Lines {
		if !assigned[i] {
			doc.PreambleLines = append(doc.PreambleLines, i)
		}
	}
	return doc
}

// stripEmphasis removes bold markers and back quotes before a heading is
// tested against the behavior names.
func stripEmphasis(s string) string {
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "`", "")
	return s
}

func matchBehaviorInHeading(text string, behaviors []string) string {
	for _, name := range behaviors {
		if containsWholeToken(text, name) {
			return name
		}
	}
	return ""
}

// PreambleText returns the hints preamble lines in source order.
func (d *HintsDoc) PreambleText() []string {
	out := make([]string, 0, len(d.PreambleLines))
	for _, i := range d.PreambleLines {
		out = append(out, d.Src.Lines[i])
	}
	return out
}

// BlockText returns one attributed block's lines verbatim.
func (d *HintsDoc) BlockText(b *HintsBlock) []string {
	return d.Src.Lines[b.Start:b.End]
}
