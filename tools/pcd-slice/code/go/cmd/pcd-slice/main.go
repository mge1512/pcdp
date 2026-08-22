// generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
//
// Command pcd-slice derives per-behavior translation bundles from a PCD
// specification. This file is the entry point: it parses the command line,
// reports invocation errors, and calls into the implementation package.
package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pcdslice "github.com/mge1512/pcd/tools/pcd-slice/internal/pcd-slice"
)

// version is the tool version, injected at build time with
// -ldflags "-X main.version=$(cat VERSION)".
var version = "dev"

const usageText = `usage:
  pcd-slice list  spec=<file> [hints=<file>]...
  pcd-slice check spec=<file> [hints=<file>]... [out=<dir>]
  pcd-slice slice spec=<file> [hints=<file>]... out=<dir>
  pcd-slice version
  pcd-slice help

options:
  spec=<file>   the specification to read; required for every verb
  hints=<file>  a hints file; repeatable, order preserved and recorded
  out=<dir>     output directory; required for slice, optional for check

exit codes:
  0  clean: check found nothing, or slice wrote its bundles
  1  findings: check found at least one, or slice refused and wrote nothing
  2  invocation error: bad arguments, missing, unreadable or unwritable path
`

func main() {
	pcdslice.Version = version
	installSignalHandler()
	os.Exit(dispatch(os.Args[1:], os.Stdout, os.Stderr))
}

// installSignalHandler exits cleanly on SIGTERM and SIGINT, removing any
// files the current run has already written so that no partial output
// survives.
func installSignalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		pcdslice.RemovePartialOutput()
		os.Exit(int(pcdslice.ExitClean))
	}()
}

func dispatch(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usageText)
		return int(pcdslice.ExitInvocation)
	}
	verb := args[0]
	switch verb {
	case "version":
		fmt.Fprintf(stdout, "%s %s\nspec:%s\n", pcdslice.ToolName, version, pcdslice.SpecSHA256)
		return int(pcdslice.ExitClean)
	case "help":
		fmt.Fprint(stdout, usageText)
		return int(pcdslice.ExitClean)
	case "list", "check", "slice":
	default:
		fmt.Fprintf(stderr, "error: unknown command: %s\n", verb)
		fmt.Fprint(stderr, usageText)
		return int(pcdslice.ExitInvocation)
	}

	opts, err := parseOptions(verb, args[1:])
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		fmt.Fprint(stderr, usageText)
		return int(pcdslice.ExitInvocation)
	}
	opts.Stdout = stdout
	opts.Stderr = stderr

	switch verb {
	case "list":
		return pcdslice.List(opts)
	case "check":
		return pcdslice.Check(opts)
	default:
		return pcdslice.Slice(opts)
	}
}

// parseOptions reads the key=value arguments of one verb. POSIX flag style is
// not accepted.
func parseOptions(verb string, args []string) (pcdslice.Options, error) {
	var o pcdslice.Options
	for _, arg := range args {
		key, value, found := strings.Cut(arg, "=")
		if !found {
			return o, fmt.Errorf("expected key=value, got: %s", arg)
		}
		if value == "" {
			return o, fmt.Errorf("empty value for %s=", key)
		}
		switch key {
		case "spec":
			if o.Spec != "" {
				return o, fmt.Errorf("spec= given more than once")
			}
			o.Spec = value
		case "hints":
			o.Hints = append(o.Hints, value)
		case "out":
			if verb == "list" {
				return o, fmt.Errorf("out= is not an option of list")
			}
			if o.Out != "" {
				return o, fmt.Errorf("out= given more than once")
			}
			o.Out = value
		default:
			return o, fmt.Errorf("unknown option: %s", arg)
		}
	}
	if o.Spec == "" {
		return o, fmt.Errorf("missing required option spec=<file>")
	}
	if verb == "slice" && o.Out == "" {
		return o, fmt.Errorf("slice requires out=<dir>")
	}
	return o, nil
}
