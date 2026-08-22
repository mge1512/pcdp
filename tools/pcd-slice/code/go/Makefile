# Makefile for pcd-slice
# generated from spec: pcd-slice.spec.md sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04

NAME        := pcd-slice
VERSION     := $(shell cat VERSION)
SPEC_SHA256 := bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
MODULE      := github.com/mge1512/pcd/tools/pcd-slice
TESTDIR     := independent_tests/claude-opus-5

PREFIX  ?= /usr/local
BINDIR  ?= $(PREFIX)/bin
MANDIR  ?= $(PREFIX)/share/man/man1
DESTDIR ?=

GO       ?= go
LDFLAGS  := -X main.version=$(VERSION)
DISTDIR  := $(NAME)-$(VERSION)

.PHONY: all build test install clean man dist vendor check-fmt

all: build

## build: static binary at the project root, per BINARY-LOCATION: project-root
build:
	CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o $(NAME) ./cmd/$(NAME)

## test: run the translator test suite against the built binary
test: build
	$(GO) test ./$(TESTDIR)/...

check-fmt:
	@test -z "$$(gofmt -l cmd internal $(TESTDIR))" || { gofmt -l cmd internal $(TESTDIR); exit 1; }
	$(GO) vet ./...

## man: troff man page from its markdown source
man: $(NAME).1

$(NAME).1: $(NAME).1.md
	pandoc $(NAME).1.md -s -t man -o $(NAME).1

install: build man
	install -D -m 0755 $(NAME) $(DESTDIR)$(BINDIR)/$(NAME)
	install -D -m 0644 $(NAME).1 $(DESTDIR)$(MANDIR)/$(NAME).1

clean:
	rm -f $(NAME) $(NAME).1
	rm -rf $(DISTDIR) $(DISTDIR).tar.gz $(DISTDIR)-vendor.tar.gz vendor

## vendor: dependency tree for offline package builds. pcd-slice uses the
## standard library only, so the vendor tree is empty but must carry a
## modules.txt for -mod=vendor builds inside OBS.
vendor:
	$(GO) mod vendor
	mkdir -p vendor
	test -f vendor/modules.txt || : > vendor/modules.txt

## dist: release source tarball plus the companion vendor tarball
dist: clean
	mkdir -p $(DISTDIR)
	cp -a VERSION Makefile README.md LICENSE go.mod \
	      $(NAME).1.md $(NAME).spec $(DISTDIR)/
	cp -a cmd internal debian independent_tests translation_report $(DISTDIR)/
	tar czf $(DISTDIR).tar.gz $(DISTDIR)
	$(MAKE) vendor
	tar czf $(DISTDIR)-vendor.tar.gz vendor
	rm -rf $(DISTDIR) vendor
