Name:           pcd-lint
Version:        0.3.13
Release:        1
Summary:        Linter and validator for Post-Coding Development specifications

License:        GPL-2.0-only
URL:            https://github.com/mge1512/pcd-lint

Source0:        pcd-lint-0.3.13.tar.gz

BuildRequires:  golang >= 1.21
Requires:       pcd-templates

%description
pcd-lint is a command-line tool that validates specification files written
in the Post-Coding Development (PCD) format. It enforces structural
rules, semantic validation, and cross-section consistency checks.

%prep
%setup -q

%build
CGO_ENABLED=0 go build -o pcd-lint .

%install
install -D -m 0755 pcd-lint %{buildroot}%{_bindir}/pcd-lint

%files
%{_bindir}/pcd-lint

%changelog
* Wed Mar 25 2026 Matthias G. Eckermann <pcd@mailbox.org> - 0.3.13-1
- Initial release of pcd-lint
