// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Package pcdslice implements the behaviours of the pcd-slice specification:
// list, check and slice. The package is the implementation module; the
// entry point under cmd/pcd-slice contains CLI dispatch only.
package pcdslice

import "fmt"

// SpecSHA256 is the SHA-256 of the specification this implementation was
// generated from. It is printed by the version verb.
const SpecSHA256 = "c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f"

// Exit codes, per the spec's ExitCode type and the cli-tool template.
const (
	ExitOK         = 0 // clean
	ExitFindings   = 1 // findings; slice refused and wrote nothing
	ExitInvocation = 2 // invocation error
)

// Kind is the spec's Assignment type: every source line receives exactly one.
type Kind int

// Assignment kinds.
const (
	KindUnassigned Kind = iota
	KindPreamble
	KindBehavior
	KindExcluded
)

// String renders the assignment kind as the spec names it.
func (k Kind) String() string {
	switch k {
	case KindPreamble:
		return "preamble"
	case KindBehavior:
		return "behavior"
	case KindExcluded:
		return "excluded"
	default:
		return "unassigned"
	}
}

// Role distinguishes why a line belongs to a behavior. It selects the
// emission slot of a line assigned to a behavior; it is not part of the
// spec's Assignment type.
type Role int

// Emission roles.
const (
	RoleNone Role = iota
	RoleBlock
	RoleExample
	RoleInvariant
)

// Assignment is the per-line classification the completeness rule scans.
type Assignment struct {
	Kind     Kind
	Behavior string
	Role     Role
}

// Finding is one diagnostic, rendered as "{rule}: {file}:{line}: {message}".
type Finding struct {
	Rule    string
	File    string
	Line    int
	Message string
}

// String renders the finding in the spec's one-line format.
func (f Finding) String() string {
	return fmt.Sprintf("%s: %s:%d: %s", f.Rule, f.File, f.Line, f.Message)
}

// Range is a half-open line range [Start, End).
type Range struct {
	Start int
	End   int
}

// Behavior is one behavior block of the specification.
type Behavior struct {
	Name           string
	FileName       string
	Heading        int
	Block          Range
	Closure        []string
	NestedExamples int
}

// TypeDef is one definition line plus its attached comment and
// continuation lines inside the TYPES section.
type TypeDef struct {
	Name string
	Line int
	Body Range
	Text string
}

// Invariant is one list item of the INVARIANTS section.
type Invariant struct {
	Item     Range
	BindLine int
	Binds    []string
	Owner    string
}

// Example is one EXAMPLE block, top-level or nested.
type Example struct {
	Name    string
	Heading int
	Block   Range
	Target  string
	Shared  bool
}
