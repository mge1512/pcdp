// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Argument handling and verb dispatch. Options are key=value pairs;
// POSIX-style flags are not accepted.

package pcdslice

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Usage is the help text printed by the help verb and on invocation errors.
const Usage = `pcd-slice - derive per-behavior translation bundles from a PCD specification

Usage:
  pcd-slice version
  pcd-slice help
  pcd-slice list  spec=<file> [hints=<file>]...
  pcd-slice check spec=<file> [hints=<file>]... [out=<dir>]
  pcd-slice slice spec=<file> [hints=<file>]... out=<dir>

Options (key=value only; POSIX-style flags are not accepted):
  spec=<file>   the specification to read; required for every verb
  hints=<file>  a hints file; repeatable, order preserved and recorded
  out=<dir>     output directory; required for slice, optional for check

Exit codes:
  0  clean
  1  findings (check found at least one; slice refused and wrote nothing)
  2  invocation error`

type options struct {
	verb  string
	spec  string
	hints []string
	out   string
}

// Run parses args and executes the selected verb. It returns the process
// exit code and never calls os.Exit itself.
func Run(args []string, version string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "error: no verb given")
		fmt.Fprintln(stderr, Usage)
		return ExitInvocation
	}

	opts := options{verb: args[0]}
	switch opts.verb {
	case "version":
		if len(args) > 1 {
			fmt.Fprintln(stderr, "error: version takes no options")
			return ExitInvocation
		}
		fmt.Fprintf(stdout, "pcd-slice %s\n", version)
		fmt.Fprintf(stdout, "spec:%s\n", SpecSHA256)
		return ExitOK
	case "help":
		fmt.Fprintln(stdout, Usage)
		return ExitOK
	case "list", "check", "slice":
	default:
		fmt.Fprintf(stderr, "error: unknown verb: %s\n", opts.verb)
		fmt.Fprintln(stderr, Usage)
		return ExitInvocation
	}

	for _, a := range args[1:] {
		if strings.HasPrefix(a, "-") {
			fmt.Fprintf(stderr, "error: POSIX-style flags are not accepted: %s\n", a)
			return ExitInvocation
		}
		k, v, ok := strings.Cut(a, "=")
		if !ok {
			fmt.Fprintf(stderr, "error: not a key=value argument: %s\n", a)
			return ExitInvocation
		}
		switch k {
		case "spec":
			if opts.spec != "" {
				fmt.Fprintln(stderr, "error: spec= given more than once")
				return ExitInvocation
			}
			if v == "" {
				fmt.Fprintln(stderr, "error: spec= requires a value")
				return ExitInvocation
			}
			opts.spec = v
		case "hints":
			if v == "" {
				fmt.Fprintln(stderr, "error: hints= requires a value")
				return ExitInvocation
			}
			opts.hints = append(opts.hints, v)
		case "out":
			if opts.out != "" {
				fmt.Fprintln(stderr, "error: out= given more than once")
				return ExitInvocation
			}
			if v == "" {
				fmt.Fprintln(stderr, "error: out= requires a value")
				return ExitInvocation
			}
			opts.out = v
		default:
			fmt.Fprintf(stderr, "error: unknown option: %s\n", a)
			return ExitInvocation
		}
	}

	if opts.spec == "" {
		fmt.Fprintln(stderr, "error: spec= is required")
		return ExitInvocation
	}
	if opts.verb == "slice" && opts.out == "" {
		fmt.Fprintln(stderr, "error: out= is required for slice")
		return ExitInvocation
	}
	if opts.verb == "list" && opts.out != "" {
		fmt.Fprintln(stderr, "error: out= is not accepted by list")
		return ExitInvocation
	}

	specData, err := readMarkdown(opts.spec)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return ExitInvocation
	}
	spec := ParseSpec(opts.spec, specData)

	hints := make([]*Hints, 0, len(opts.hints))
	for _, p := range opts.hints {
		data, err := readMarkdown(p)
		if err != nil {
			fmt.Fprintf(stderr, "%v\n", err)
			return ExitInvocation
		}
		hints = append(hints, ParseHints(p, data, spec))
	}

	switch opts.verb {
	case "list":
		return List(spec, hints, stdout)
	case "check":
		return Check(spec, hints, opts.out, version, stdout, stderr)
	default:
		return Slice(spec, hints, opts.out, version, stdout, stderr)
	}
}

// readMarkdown reads a specification or hints file, enforcing the
// preconditions: it exists, is readable, and ends in ".md".
func readMarkdown(path string) ([]byte, error) {
	if !strings.HasSuffix(path, ".md") {
		return nil, fmt.Errorf("error: not a .md file: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error: cannot read %s", path)
	}
	return data, nil
}
