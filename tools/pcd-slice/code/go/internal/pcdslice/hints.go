// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Hints files: heading attribution and the unassignable-hints rule.

package pcdslice

import (
	"regexp"
	"strings"
)

var singleIdentRe = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// HintsBlock is a heading-delimited block of a hints file attributed to one
// behavior.
type HintsBlock struct {
	Behavior string
	Block    Range
}

// Hints is one parsed hints file.
type Hints struct {
	Path     string
	SHA256   string
	Lines    []string
	Blocks   []HintsBlock
	Preamble []Range
	Findings []Finding

	inFence []bool
	isFence []bool
}

// ParseHints parses a hints file against the behavior names of spec. A "## "
// or deeper heading whose text contains a BehaviorName as a whole word opens
// a block attributed to that behavior, extending to the next heading of the
// same or shallower depth. A depth-1 heading closes a block and never opens
// one. Everything else is hints preamble.
func ParseHints(path string, data []byte, spec *Spec) *Hints {
	h := &Hints{Path: path, SHA256: HashBytes(data), Lines: SplitLines(data)}
	h.inFence, h.isFence = fenceArrays(h.Lines)
	n := len(h.Lines)
	assigned := make([]bool, n)

	openBehavior := ""
	openDepth := 0
	openStart := 0
	closeBlock := func(end int) {
		if openBehavior == "" {
			return
		}
		h.Blocks = append(h.Blocks, HintsBlock{
			Behavior: openBehavior,
			Block:    Range{openStart, trimTrailingBlank(h.Lines, openStart, end)},
		})
		for i := openStart; i < end; i++ {
			assigned[i] = true
		}
		openBehavior = ""
	}

	for _, hd := range headingsOf(h.Lines, h.inFence, h.isFence) {
		if openBehavior != "" && hd.depth <= openDepth {
			closeBlock(hd.line)
		}
		if hd.depth < 2 {
			continue
		}
		if name := behaviorInHeading(spec, hd.text); name != "" {
			openBehavior, openDepth, openStart = name, hd.depth, hd.line
			continue
		}
		if stripped := stripHeadingMarkup(hd.text); singleIdentRe.MatchString(stripped) {
			h.Findings = append(h.Findings, Finding{
				Rule: "unassignable-hints", File: path, Line: hd.line + 1,
				Message: "heading names no known behavior: " + stripped,
			})
		}
	}
	closeBlock(n)

	start := -1
	for i := 0; i <= n; i++ {
		if i < n && !assigned[i] {
			if start < 0 {
				start = i
			}
			continue
		}
		if start >= 0 {
			if e := trimTrailingBlank(h.Lines, start, i); e > start {
				h.Preamble = append(h.Preamble, Range{start, e})
			}
			start = -1
		}
	}
	return h
}

// BlocksFor returns the blocks attributed to a behavior, in file order.
func (h *Hints) BlocksFor(behavior string) []HintsBlock {
	var out []HintsBlock
	for _, b := range h.Blocks {
		if b.Behavior == behavior {
			out = append(out, b)
		}
	}
	return out
}

// behaviorInHeading returns the first behavior name occurring as a whole
// word in a heading, after bold markers and backticks are stripped.
func behaviorInHeading(spec *Spec, text string) string {
	if spec.behaviorMatcher == nil {
		return ""
	}
	t := stripHeadingMarkup(text)
	loc := spec.behaviorMatcher.FindStringIndex(t)
	if loc == nil {
		return ""
	}
	return t[loc[0]:loc[1]]
}

func stripHeadingMarkup(text string) string {
	t := strings.ReplaceAll(text, "**", "")
	t = strings.ReplaceAll(t, "`", "")
	return strings.TrimSpace(t)
}
