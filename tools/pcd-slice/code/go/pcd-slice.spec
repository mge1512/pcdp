# pcd-spec-sha256: bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
#
# spec file for package pcd-slice
#
# Copyright (C) 2026 Matthias G. Eckermann <pcd@mailbox.org>
#

Name:           pcd-slice
Version:        0.1.0
Release:        0
Summary:        Derive per-behavior translation bundles from a PCD specification
License:        GPL-2.0-only
Group:          Development/Tools/Other
URL:            https://github.com/mge1512/pcd/tools/pcd-slice
Source0:        pcd-slice-%{version}.tar.gz
Source1:        pcd-slice-%{version}-vendor.tar.gz

BuildRequires:  go >= 1.22
BuildRequires:  make
BuildRequires:  pandoc
BuildRequires:  tar
BuildRequires:  gzip

%description
pcd-slice computes the semantic form of a Post-Coding Development
specification once, deterministically. For each behavior it emits one bundle
holding that behavior's block, the examples that exercise it, the invariants
bound to it, the hints written for it, and the names of the types it needs;
plus one preamble holding what every bundle shares.

The source specification is never modified. Every emitted file carries the
SHA-256 of the sources it was derived from, so a bundle is verifiable against
the reproducibility tuple and stale bundles are detectable. The tool makes no
network calls, reads no configuration file, and honours no environment
variable.

%prep
%autosetup -n pcd-slice-%{version}
tar xf %{SOURCE1} -C .

%build
export GOFLAGS="-mod=vendor"
export CGO_ENABLED=0
go build -ldflags "-X main.version=%{version}" -o pcd-slice ./cmd/pcd-slice
pandoc pcd-slice.1.md -s -t man -o pcd-slice.1

%install
install -D -m 0755 pcd-slice %{buildroot}%{_bindir}/pcd-slice
install -D -m 0644 pcd-slice.1 %{buildroot}%{_mandir}/man1/pcd-slice.1

%check
./pcd-slice version

%files
%license LICENSE
%doc README.md
%{_bindir}/pcd-slice
%{_mandir}/man1/pcd-slice.1*

%changelog
* Sat Aug 22 2026 Matthias G. Eckermann <pcd@mailbox.org> - 0.1.0
- Initial package, generated from pcd-slice.spec.md
  sha256:bb31cb27325030a420f0c1ffd8fd2660c9d9ee768f8b50f475360e451b8e6e04
