Name:           pcdp-lint
Version:        0.3.7
Release:        1%{?dist}
Summary:        Post-Coding Development Paradigm specification linter
License:        GPL-2.0-only
URL:            https://github.com/pcdp/pcdp-lint
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.21
BuildRequires:  make

%description
pcdp-lint is a command-line tool for validating Post-Coding Development
Paradigm specification files. It validates structural requirements,
META field formats, deployment template compatibility, and example
block completeness.

%prep
%setup -q

%build
CGO_ENABLED=0 make build

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 build/pcdp-lint %{buildroot}%{_bindir}/pcdp-lint

%files
%license LICENSE
%doc README.md
%{_bindir}/pcdp-lint

%changelog
* Thu Jan 01 2024 Matthias G. Eckermann <pcdp@mailbox.org> - 0.3.7-1
- Initial package