// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
//
// Package pcdslice derives per-behavior translation bundles from a PCD
// specification. It implements the three verbs the specification declares -
// list, check and slice - over the small section grammar the specification
// states. It depends on the Go standard library only.
package pcdslice

import (
	"fmt"
	"io"
	"sort"
)

// SpecSHA256 is the SHA-256 of the merged specification this implementation
// was generated from.
const SpecSHA256 = "bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04"

// ToolName is the component name; it appears in output summaries.
const ToolName = "pcd-slice"

// Version is the tool version recorded in the provenance line of every
// emitted file. The entry point overrides it from the VERSION file via
// -ldflags "-X main.version=...".
var Version = "dev"

// ExitCode enumerates the three exit codes the specification allows.
type ExitCode int

// The exit codes of every verb.
const (
	// ExitClean is 0: check found nothing, slice wrote its bundles.
	ExitClean ExitCode = 0
	// ExitFindings is 1: check found at least one finding, or slice
	// refused and wrote nothing.
	ExitFindings ExitCode = 1
	// ExitInvocation is 2: bad arguments, missing, unreadable or
	// unwritable path.
	ExitInvocation ExitCode = 2
)

// AssignKind is the assignment a specification line receives. Every line of
// the specification receives exactly one.
type AssignKind uint8

// The four assignment kinds of the specification's Assignment type.
const (
	// AssignUnassigned is the sentinel; a line left with it is a finding.
	AssignUnassigned AssignKind = iota
	// AssignPreamble sends the line to preamble.md.
	AssignPreamble
	// AssignBehavior sends the line to one behavior's bundle.
	AssignBehavior
	// AssignExcluded covers the changelog only.
	AssignExcluded
)

// String renders the assignment kind with the name the specification uses.
func (k AssignKind) String() string {
	switch k {
	case AssignPreamble:
		return "preamble"
	case AssignBehavior:
		return "behavior"
	case AssignExcluded:
		return "excluded"
	default:
		return "unassigned"
	}
}

// Assignment is the assignment of one specification line.
type Assignment struct {
	Kind AssignKind
	// Behavior is the index into Spec.Behaviors when Kind is
	// AssignBehavior, and -1 otherwise.
	Behavior int
}

// Finding is one diagnostic of the check verb.
type Finding struct {
	Rule    string
	File    string
	Line    int
	Message string
}

// String renders the finding in the one-line form the specification
// declares: "{rule}: {file}:{line}: {message}".
func (f Finding) String() string {
	return fmt.Sprintf("%s: %s:%d: %s", f.Rule, f.File, f.Line, f.Message)
}

// SortFindings orders findings by file, then line, then rule and message, so
// that two runs over identical inputs report in identical order.
func SortFindings(fs []Finding) {
	sort.SliceStable(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		return a.Message < b.Message
	})
}

// Options carries one invocation of a verb.
type Options struct {
	// Spec is the specification path exactly as given on the command line.
	Spec string
	// Hints are the hints paths in input order, exactly as given.
	Hints []string
	// Out is the output directory; empty when absent.
	Out string
	// Stdout receives normal output, Stderr findings and errors.
	Stdout io.Writer
	Stderr io.Writer
}
