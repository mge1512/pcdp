Name:           spec-lint
Version:        0.3.2
Release:        1%{?dist}
Summary:        Post-Coding Development Paradigm specification linter

License:        Apache-2.0
URL:            https://github.com/post-coding/spec-lint
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.19
Requires:       /usr/bin/sh

%description
spec-lint is a command-line tool for validating specification files written
in the Post-Coding Development Paradigm format. It checks structural rules,
validates metadata fields, and ensures compliance with deployment templates.

%prep
%setup -q

%build
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -a -ldflags '-extldflags "-static" -X main.TemplateDir=/usr/share/post-coding/templates/' -o %{name} .

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 %{name} %{buildroot}%{_bindir}/%{name}

mkdir -p %{buildroot}%{_datadir}/post-coding/templates
# Template files would be installed here in a complete package

%files
%{_bindir}/%{name}
%dir %{_datadir}/post-coding
%dir %{_datadir}/post-coding/templates

%changelog
* Thu Jan 01 2024 Matthias G. Eckermann <post-coding-development-paradigm@mailbox.org> - 0.3.2-1
- Initial package for spec-lint 0.3.2