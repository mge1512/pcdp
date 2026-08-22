// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const (
	manifestName = "MANIFEST.tsv"
	preambleName = "preamble.md"
)

// OutFile is one file the slice verb would write.
type OutFile struct {
	Name string
	Data []byte
}

// Outputs is the complete derived form, computed in memory before anything
// touches the disk.
type Outputs struct {
	// Files are the markdown files, sorted by name.
	Files []OutFile
	// Manifest is MANIFEST.tsv; it lists every file except itself.
	Manifest OutFile
	// Provenance is the provenance line every markdown file begins with.
	Provenance string
	// Bundles is the number of behavior bundles among Files.
	Bundles int
}

// All returns the markdown files followed by the manifest.
func (o *Outputs) All() []OutFile {
	out := make([]OutFile, 0, len(o.Files)+1)
	out = append(out, o.Files...)
	out = append(out, o.Manifest)
	return out
}

// outBuf accumulates output lines. Source lines are appended verbatim,
// including any trailing "\r"; the lines pcd-slice adds itself end in "\n"
// regardless of the source's line-ending flavour.
type outBuf struct {
	lines []string
}

func (b *outBuf) add(l string)         { b.lines = append(b.lines, l) }
func (b *outBuf) addAll(ls ...string)  { b.lines = append(b.lines, ls...) }
func (b *outBuf) addLines(ls []string) { b.lines = append(b.lines, ls...) }

func (b *outBuf) trimTrailingBlank() {
	for len(b.lines) > 0 && strings.TrimSpace(b.lines[len(b.lines)-1]) == "" {
		b.lines = b.lines[:len(b.lines)-1]
	}
}

// section appends attributed material under a heading that names its source.
func (b *outBuf) section(heading string, body []string) {
	body = trimBlankEdges(body)
	if len(body) == 0 {
		return
	}
	b.trimTrailingBlank()
	b.addAll("", heading, "")
	b.addLines(body)
}

func (b *outBuf) bytes() []byte {
	if len(b.lines) == 0 {
		return nil
	}
	return []byte(strings.Join(b.lines, "\n") + "\n")
}

func trimBlankEdges(ls []string) []string {
	start, end := 0, len(ls)
	for start < end && strings.TrimSpace(ls[start]) == "" {
		start++
	}
	for end > start && strings.TrimSpace(ls[end-1]) == "" {
		end--
	}
	return ls[start:end]
}

// ProvenanceLine renders the single provenance comment of every emitted
// markdown file. Paths appear as given, hashes are lowercase hex, hints
// entries are comma separated in input order, and there is no timestamp.
func (a *Analysis) ProvenanceLine() string {
	paths := make([]string, 0, len(a.Sources.Hints))
	hashes := make([]string, 0, len(a.Sources.Hints))
	for _, h := range a.Sources.Hints {
		paths = append(paths, h.Path)
		hashes = append(hashes, h.SHA256)
	}
	return fmt.Sprintf("<!-- %s source=%s source-sha256=%s hints=%s hints-sha256=%s version=%s -->",
		ToolName,
		a.Sources.Spec.Path,
		a.Sources.Spec.SHA256,
		strings.Join(paths, ","),
		strings.Join(hashes, ","),
		Version)
}

// BuildOutputs computes the preamble, one bundle per behavior, and the
// manifest. Byte content is fully determined by the inputs.
func (a *Analysis) BuildOutputs() *Outputs {
	prov := a.ProvenanceLine()
	out := &Outputs{Provenance: prov}

	out.Files = append(out.Files, OutFile{Name: preambleName, Data: a.buildPreamble(prov)})

	seen := map[string]bool{preambleName: true}
	for _, b := range a.Spec.Behaviors {
		if seen[b.FileName] {
			continue
		}
		seen[b.FileName] = true
		out.Files = append(out.Files, OutFile{Name: b.FileName, Data: a.buildBundle(prov, b)})
		out.Bundles++
	}

	sort.Slice(out.Files, func(i, j int) bool { return out.Files[i].Name < out.Files[j].Name })

	var mb outBuf
	for _, f := range out.Files {
		sum := sha256.Sum256(f.Data)
		mb.add(f.Name + "\t" + hex.EncodeToString(sum[:]))
	}
	out.Manifest = OutFile{Name: manifestName, Data: mb.bytes()}
	return out
}

// buildPreamble projects every line stamped "preamble", in source order, and
// appends the hints preamble of each hints file.
func (a *Analysis) buildPreamble(prov string) []byte {
	var b outBuf
	b.add(prov)
	for i, as := range a.Spec.Assign {
		if as.Kind == AssignPreamble {
			b.add(a.Spec.Src.Lines[i])
		}
	}
	b.trimTrailingBlank()
	for _, doc := range a.Hints {
		b.section(fmt.Sprintf("## Hints Preamble (from %s)", doc.Src.Path), doc.PreambleText())
	}
	return b.bytes()
}

// buildBundle projects one behavior's lines and appends the material
// attributed to it, each part under a heading that names its source.
func (a *Analysis) buildBundle(prov string, beh *Behavior) []byte {
	sp := a.Spec
	var b outBuf
	b.add(prov)
	b.add(requiresTypesLine(beh.Closure))

	var block []string
	for i := beh.Start; i < beh.End; i++ {
		if as := sp.Assign[i]; as.Kind == AssignBehavior && as.Behavior >= 0 && sp.Behaviors[as.Behavior] == beh {
			block = append(block, sp.Src.Lines[i])
		}
	}
	b.addLines(block)

	var examples []string
	for _, ex := range beh.Examples {
		if len(examples) > 0 {
			examples = append(examples, "")
		}
		examples = append(examples, sp.Src.Lines[ex.Start:ex.End]...)
	}
	b.section(fmt.Sprintf("## Examples (from %s)", sp.Src.Path), examples)

	var invariants []string
	for _, inv := range beh.Invariants {
		invariants = append(invariants, sp.Src.Lines[inv.Start:inv.End]...)
	}
	b.section(fmt.Sprintf("## Invariants (from %s)", sp.Src.Path), invariants)

	for _, doc := range a.Hints {
		var blocks []string
		for _, blk := range beh.HintsBlocks {
			if blk.Doc != doc {
				continue
			}
			if len(blocks) > 0 {
				blocks = append(blocks, "")
			}
			blocks = append(blocks, doc.BlockText(blk)...)
		}
		b.section(fmt.Sprintf("## Hints (from %s)", doc.Src.Path), blocks)
	}

	return b.bytes()
}

// requiresTypesLine lists the type closure alphabetically.
func requiresTypesLine(closure []string) string {
	if len(closure) == 0 {
		return "Requires-Types:"
	}
	return "Requires-Types: " + strings.Join(closure, ", ")
}
