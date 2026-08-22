// generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
//
// Entry point: CLI dispatch only. Every behaviour lives in
// internal/pcdslice. There is exactly one os.Exit, at the top level, so no
// cleanup path is ever skipped.

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mge1512/pcd/tools/pcd-slice/internal/pcdslice"
)

// version is the tool version; the Makefile injects the content of the
// VERSION file via -ldflags "-X main.version=...".
var version = "0.1.1"

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		// Clean exit: remove anything this run has written so far, so no
		// partial output survives.
		pcdslice.Abort()
		os.Exit(0)
	}()

	os.Exit(pcdslice.Run(os.Args[1:], version, os.Stdout, os.Stderr))
}
