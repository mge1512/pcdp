# generated from spec: pcd-slice.spec.md sha256:c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
#
# pcd-slice build orchestration. Every target runs unprivileged; no target
# invokes sudo and no target fetches from the network (the implementation
# uses the Go standard library only).

BINARY_NAME := pcd-slice
VERSION     := $(shell cat VERSION)
SPEC_SHA256 := c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
LLM_NAME    := claude-opus-5

GO          ?= go
PREFIX      ?= /usr/local
BINDIR      ?= $(PREFIX)/bin
MANDIR      ?= $(PREFIX)/share/man/man1
DESTDIR     ?=

LDFLAGS     := -s -w -X main.version=$(VERSION)
DISTDIR     := $(BINARY_NAME)-$(VERSION)
DISTFILES   := cmd internal independent_tests debian translation_report \
               go.mod Makefile VERSION README.md LICENSE \
               $(BINARY_NAME).1.md $(BINARY_NAME).spec

.PHONY: all build test install clean man dist vendor

all: build

## build: static binary at the project root (BINARY-LOCATION: project-root)
build:
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags '$(LDFLAGS)' \
	    -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

## test: run the test suite; exits non-zero on any failure
test: build
	$(GO) test ./independent_tests/$(LLM_NAME)/...

## man: troff man page from its markdown source
man: $(BINARY_NAME).1

$(BINARY_NAME).1: $(BINARY_NAME).1.md
	pandoc $(BINARY_NAME).1.md -s -t man -o $(BINARY_NAME).1

## install: binary and man page under $(DESTDIR)$(PREFIX)
install: build man
	install -D -m 0755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	install -D -m 0644 $(BINARY_NAME).1 $(DESTDIR)$(MANDIR)/$(BINARY_NAME).1

## vendor: materialise the dependency tree (empty: standard library only)
vendor:
	$(GO) mod vendor

## clean: remove every build output; the source tree is left untouched
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).1
	rm -f $(BINARY_NAME)-*.tar.gz
	rm -rf build

## dist: release source tarball with a single top-level directory
dist: clean
	mkdir -p build/$(DISTDIR)
	cp -a $(DISTFILES) build/$(DISTDIR)/
	tar --exclude='.git' --exclude='build' --exclude='vendor' \
	    -czf $(DISTDIR).tar.gz -C build $(DISTDIR)
	@if [ -d vendor ]; then \
	    tar -czf $(BINARY_NAME)-$(VERSION)-vendor.tar.gz vendor; \
	    echo "wrote $(BINARY_NAME)-$(VERSION)-vendor.tar.gz"; \
	else \
	    echo "no vendor/ tree: pcd-slice depends on the standard library only"; \
	fi
	rm -rf build/$(DISTDIR)
	@echo "wrote $(DISTDIR).tar.gz"
