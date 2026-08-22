# pcd-spec-sha256: c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f
#
# RPM spec for pcd-slice. The VERSION file is the sole version authority:
# `make dist` derives both the tarball name and the top-level directory from
# it, and the Version: field below must equal its content (0.1.1).

Name:           pcd-slice
Version:        0.1.1
Release:        0
Summary:        Derive per-behavior translation bundles from a PCD specification
License:        GPL-2.0-only
Group:          Development/Tools/Building
URL:            https://github.com/mge1512/pcd
Source0:        pcd-slice-%{version}.tar.gz

BuildRequires:  go >= 1.22
BuildRequires:  make
BuildRequires:  pandoc

%description
pcd-slice derives per-behavior translation bundles from a Post-Coding
Development specification. For each behavior it emits one bundle holding
that behavior's block, the examples that exercise it, the invariants bound
to it, the hints written for it and the names of the types it needs; plus
one preamble holding what every bundle shares.

The source specification is never modified. Every emitted markdown file
carries a provenance comment naming its sources and their SHA-256 checksums,
so bundles are verifiable against the reproducibility tuple and stale
bundles are detectable. The tool reads no configuration, honours no
environment variable and makes no network calls.

%prep
%autosetup -n pcd-slice-%{version}

%build
export CGO_ENABLED=0
export GOFLAGS="-trimpath"
export GOCACHE="%{_builddir}/go-build"
export GOPATH="%{_builddir}/go"
go build -ldflags "-s -w -X main.version=%{version}" -o pcd-slice ./cmd/pcd-slice
pandoc pcd-slice.1.md -s -t man -o pcd-slice.1

%check
export CGO_ENABLED=0
export GOCACHE="%{_builddir}/go-build"
export GOPATH="%{_builddir}/go"
go test ./independent_tests/claude-opus-5/...

%install
install -D -m 0755 pcd-slice %{buildroot}%{_bindir}/pcd-slice
install -D -m 0644 pcd-slice.1 %{buildroot}%{_mandir}/man1/pcd-slice.1

%files
%license LICENSE
%doc README.md
%{_bindir}/pcd-slice
%{_mandir}/man1/pcd-slice.1*

%changelog
* Sat Aug 22 2026 Matthias G. Eckermann <pcd@mailbox.org> - 0.1.1
- Initial packaging, translated from pcd-slice.spec.md version 0.1.1
  (spec sha256 c200af55518fdfe36de34be5a2dfdcc822a16e31410c131bbd8390a6d7897c0f).
