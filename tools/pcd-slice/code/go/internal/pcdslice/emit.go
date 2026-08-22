// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Emission: the derived form is computed fully in memory before anything
// touches disk. Byte content is a function of the inputs alone - sections in
// source order, closures alphabetical, manifest sorted, no timestamp.

package pcdslice

import (
	"fmt"
	"sort"
	"strings"
)

// PreambleFile is the name of the shared-context file.
const PreambleFile = "preamble.md"

// ManifestFile is the name of the manifest. It carries no provenance line -
// it is not markdown - and is compared by full content.
const ManifestFile = "MANIFEST.tsv"

// Output is the complete derived form of one slicing run.
type Output struct {
	Provenance string
	Names      []string // emitted markdown files, sorted; excludes the manifest
	Files      map[string][]byte
	Manifest   []byte
	Bundles    int
}

type lineBuf struct {
	b strings.Builder
}

func (w *lineBuf) line(s string) {
	w.b.WriteString(s)
	w.b.WriteString("\n")
}

func (w *lineBuf) blank() {
	w.b.WriteString("\n")
}

// ProvenanceLine renders the single provenance comment every emitted
// markdown file begins with.
func ProvenanceLine(spec *Spec, hints []*Hints, version string) string {
	paths := make([]string, 0, len(hints))
	sums := make([]string, 0, len(hints))
	for _, h := range hints {
		paths = append(paths, h.Path)
		sums = append(sums, h.SHA256)
	}
	return fmt.Sprintf("<!-- pcd-slice source=%s source-sha256=%s hints=%s hints-sha256=%s version=%s -->",
		spec.Path, spec.SHA256, strings.Join(paths, ","), strings.Join(sums, ","), version)
}

// BuildOutput computes every emitted file in memory.
func BuildOutput(spec *Spec, hints []*Hints, version string) *Output {
	prov := ProvenanceLine(spec, hints, version)
	out := &Output{
		Provenance: prov,
		Files:      map[string][]byte{},
		Bundles:    len(spec.Behaviors),
	}
	out.Files[PreambleFile] = buildPreamble(spec, hints, prov)
	for _, b := range spec.Behaviors {
		out.Files[b.FileName] = buildBundle(spec, hints, prov, b)
	}
	for name := range out.Files {
		out.Names = append(out.Names, name)
	}
	sort.Strings(out.Names)

	var m lineBuf
	for _, name := range out.Names {
		m.line(name + "\t" + HashBytes(out.Files[name]))
	}
	out.Manifest = []byte(m.b.String())
	return out
}

func buildPreamble(spec *Spec, hints []*Hints, prov string) []byte {
	var w lineBuf
	w.line(prov)

	var idx []int
	for i, a := range spec.Assign {
		if a.Kind == KindPreamble {
			idx = append(idx, i)
		}
	}
	for len(idx) > 0 && strings.TrimSpace(spec.Lines[idx[len(idx)-1]]) == "" {
		idx = idx[:len(idx)-1]
	}
	for _, i := range idx {
		w.line(spec.Line(i))
	}

	for _, h := range hints {
		if len(h.Preamble) == 0 {
			continue
		}
		w.blank()
		w.line("## Hints preamble from " + h.Path)
		w.blank()
		for k, r := range h.Preamble {
			if k > 0 {
				w.blank()
			}
			for i := r.Start; i < r.End; i++ {
				w.line(h.Lines[i])
			}
		}
	}
	return []byte(w.b.String())
}

func buildBundle(spec *Spec, hints []*Hints, prov string, b Behavior) []byte {
	var w lineBuf
	w.line(prov)

	req := "Requires-Types:"
	if len(b.Closure) > 0 {
		req += " " + strings.Join(b.Closure, ", ")
	}
	w.line(req)
	w.blank()

	var idx []int
	end := b.Block.End
	if end > len(spec.Lines) {
		end = len(spec.Lines)
	}
	for i := b.Block.Start; i < end; i++ {
		a := spec.Assign[i]
		if a.Kind == KindBehavior && a.Behavior == b.Name && a.Role == RoleBlock {
			idx = append(idx, i)
		}
	}
	for len(idx) > 0 && strings.TrimSpace(spec.Lines[idx[len(idx)-1]]) == "" {
		idx = idx[:len(idx)-1]
	}
	for _, i := range idx {
		w.line(spec.Line(i))
	}

	var exs []Example
	for _, e := range spec.Examples {
		if e.Target == b.Name {
			exs = append(exs, e)
		}
	}
	if len(exs) > 0 {
		w.blank()
		w.line("## Attributed examples from " + spec.Path)
		w.blank()
		for k, e := range exs {
			if k > 0 {
				w.blank()
			}
			for i := e.Block.Start; i < e.Block.End; i++ {
				w.line(spec.Line(i))
			}
		}
	}

	var invs []Invariant
	for _, inv := range spec.Invariants {
		if containsString(inv.Binds, b.Name) {
			invs = append(invs, inv)
		}
	}
	if len(invs) > 0 {
		w.blank()
		w.line("## Bound invariants from " + spec.Path)
		w.blank()
		for _, inv := range invs {
			for i := inv.Item.Start; i < inv.Item.End; i++ {
				w.line(spec.Line(i))
			}
		}
	}

	for _, h := range hints {
		blocks := h.BlocksFor(b.Name)
		if len(blocks) == 0 {
			continue
		}
		w.blank()
		w.line("## Hints from " + h.Path)
		w.blank()
		for k, blk := range blocks {
			if k > 0 {
				w.blank()
			}
			for i := blk.Block.Start; i < blk.Block.End; i++ {
				w.line(h.Lines[i])
			}
		}
	}
	return []byte(w.b.String())
}
