// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

package pcdslice

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Source is one input file read into memory, with the checksum of exactly
// the bytes that were read - the same bytes the provenance line names.
type Source struct {
	// Path is the path as given on the command line.
	Path string
	// Data is the file content as read.
	Data []byte
	// SHA256 is the lowercase hex checksum of Data.
	SHA256 string
	// Lines is Data split on "\n"; a trailing "\r" stays attached to its
	// line so that the emitted bundles are byte-faithful to the source.
	Lines []string
}

// LoadSource reads one input file after checking the refinement predicates
// of the SpecFile and HintsFile types: the path exists, is readable, and
// ends in ".md".
func LoadSource(path string) (*Source, error) {
	if strings.ToLower(filepath.Ext(path)) != ".md" {
		return nil, fmt.Errorf("not a markdown file (must end in .md): %s", path)
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return nil, fmt.Errorf("cannot read %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s", path)
	}
	sum := sha256.Sum256(data)
	return &Source{
		Path:   path,
		Data:   data,
		SHA256: hex.EncodeToString(sum[:]),
		Lines:  SplitLines(data),
	}, nil
}

// SplitLines splits file content into lines on "\n". A trailing "\r" is kept
// attached to its line; a final line without a newline is one line.
func SplitLines(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	s := string(data)
	parts := strings.Split(s, "\n")
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// Sources is a specification together with its hints files, in input order.
type Sources struct {
	Spec  *Source
	Hints []*Source
}

// LoadSources reads the specification and every hints file named in opts.
func LoadSources(o Options) (*Sources, error) {
	spec, err := LoadSource(o.Spec)
	if err != nil {
		return nil, err
	}
	s := &Sources{Spec: spec}
	for _, h := range o.Hints {
		hs, err := LoadSource(h)
		if err != nil {
			return nil, err
		}
		s.Hints = append(s.Hints, hs)
	}
	return s, nil
}
