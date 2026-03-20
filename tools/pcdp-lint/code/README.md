# pcdp-lint

A command-line tool for validating Post-Coding Development Paradigm specification files.

## Installation

### Via Package Manager (Recommended)

Install via the OBS (openSUSE Build Service) package repository:

#### openSUSE/SLES
```bash
sudo zypper install pcdp-tools
```

#### Fedora
```bash
sudo dnf install pcdp-tools
```

#### Debian/Ubuntu
```bash
sudo apt update
sudo apt install pcdp-tools
```

### Manual Installation

Build from source:
```bash
make build
sudo make install
```

## Usage

### Basic Linting

Validate a specification file:
```bash
pcdp-lint myspec.md
```

### Strict Mode

Treat warnings as errors:
```bash
pcdp-lint strict=true myspec.md
```

### List Available Templates

Display all known deployment templates:
```bash
pcdp-lint list-templates
```

### Version Information

Show version and build information:
```bash
pcdp-lint version
```

## Exit Codes

- **0**: Valid (no errors; no warnings when strict=true)
- **1**: Invalid (at least one error; or strict=true and at least one warning)
- **2**: Invocation error (bad arguments, file not found, unreadable file)

## Output Format

### Diagnostic Messages (stderr)

```
ERROR    myspec.md:42   [EXAMPLES]   Example 'foo' missing THEN: marker
WARNING  myspec.md:6    [META]       META field 'Target' is deprecated since v0.3.0
```

### Summary Messages (stdout)

```
✓ myspec.md: valid
✓ myspec.md: valid (1 warning(s))
✗ myspec.md: 1 error(s), 0 warning(s)
✗ myspec.md: 0 error(s), 1 warning(s) [strict mode]
```

## Validation Rules

pcdp-lint validates the following aspects of specification files:

### Required Sections
- `## META`
- `## TYPES` 
- `## BEHAVIOR` (or `## BEHAVIOR: name` or `## BEHAVIOR/INTERNAL: name`)
- `## PRECONDITIONS`
- `## POSTCONDITIONS`
- `## INVARIANTS`
- `## EXAMPLES`

### META Field Validation
- **Required fields**: Deployment, Verification, Safety-Level, Version, Spec-Schema, License
- **Author field**: At least one Author: line required
- **Version format**: Must follow semantic versioning (MAJOR.MINOR.PATCH)
- **License**: Must be a valid SPDX license identifier

### Deployment Template Validation
- Deployment must be a known template name
- Special requirements for specific deployments:
  - `enhance-existing`: requires Language field
  - `manual`: requires Target field
  - `python-tool`: requires Safety-Level: QM and Verification: none
  - `verified-library`: warns if Safety-Level: QM (unusual)

### Examples Section Validation
- Must contain at least one example block
- Each example must have EXAMPLE:, GIVEN:, WHEN:, THEN: markers
- Warns about empty GIVEN, WHEN, or THEN blocks

## Supported Key=Value Options

- `strict=true|false`: Enable/disable strict mode (default: false)

## Platform Support

- **Linux**: Primary platform (fully supported)
- **macOS**: Supported
- **Windows**: Not supported in v1

## License

GPL-2.0-only

## Author

Matthias G. Eckermann <pcdp@mailbox.org>