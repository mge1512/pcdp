# mcp-server-pcd

An MCP (Model Context Protocol) server for managing PCD (Post-Coding Development) specifications.

## Overview

`mcp-server-pcd` provides a complete MCP server implementation that enables AI assistants and other MCP hosts to:

- List and retrieve PCD templates
- Access embedded prompts (interview and translator)
- Read hints files for library integration
- Validate PCD specifications using the full rule set
- Lint specification files from disk

The server supports both **stdio** (for CLI-based hosts) and **HTTP** (for web-based hosts) transports in a single binary.

## Installation

### From OBS Package Repository

#### openSUSE Leap / SUSE Linux Enterprise
```bash
zypper install mcp-server-pcd
```

#### Fedora
```bash
dnf install mcp-server-pcd
```

#### Debian / Ubuntu
```bash
apt install mcp-server-pcd
```

### From Source

```bash
git clone https://github.com/mge1512/mcp-server-pcd.git
cd mcp-server-pcd
make build
sudo make install
```

### Docker

```bash
docker build -t mcp-server-pcd .
docker run -p 8080:8080 mcp-server-pcd
```

### Podman (recommended for SUSE/Linux systems)

The Containerfile uses `registry.suse.com/bci/golang:latest` as the builder base image,
ensuring compatibility with SUSE Linux Enterprise and openSUSE systems.

```bash
# Build the container image
make container

# Or manually with podman
podman build -t mcp-server-pcd:latest -f Containerfile .

# Run the container
podman run -p 8080:8080 mcp-server-pcd:latest

# Test the container
make container-test

# Clean up
make container-clean
```

**Container Image Details:**

- **Builder stage:** `registry.suse.com/bci/golang:latest` (SUSE BCI Go image)
- **Final stage:** `FROM scratch` (minimal image, static binary only)
- **Exposed port:** 8080 (HTTP transport)
- **Default entrypoint:** `mcp-server-pcd http` (HTTP mode)
- **Image size:** ~15 MB (static binary only, no runtime dependencies)

## Usage

### Stdio Transport (for CLI-based MCP hosts)

Run the server in stdio mode:

```bash
mcp-server-pcd stdio
```

This is the default transport. The server reads JSON-RPC messages from stdin and writes responses to stdout. Diagnostics and errors go to stderr.

**Configuration in mcphost:**

```yaml
mcpServers:
  pcd:
    command: mcp-server-pcd
    args: [stdio]
```

### HTTP Transport (for web-based hosts)

Run the server in HTTP mode:

```bash
mcp-server-pcd http
```

By default, the server listens on `127.0.0.1:8080`. Override with the `listen=` argument:

```bash
mcp-server-pcd http listen=0.0.0.0:9000
```

The server exposes the MCP protocol on the `/mcp` endpoint. Clients send POST requests with JSON-RPC bodies and receive responses with HTTP status 200 and JSON-RPC result bodies.

**Configuration in mcphost:**

```yaml
mcpServers:
  pcd:
    url: http://127.0.0.1:8080/mcp
```

### Systemd Service (HTTP mode)

If installed via package, a systemd service unit is available:

```bash
systemctl start mcp-server-pcd
systemctl enable mcp-server-pcd
```

The service runs in HTTP mode on `127.0.0.1:8080`.

## MCP Tools

The server exposes the following MCP tools:

### `list_templates`

List all installed PCD templates.

**Arguments:** none

**Returns:** JSON array of template records with `name`, `version`, `language` fields (content omitted).

**Example:**
```json
[
  {"name": "cli-tool", "version": "0.3.17", "language": "go"},
  {"name": "mcp-server", "version": "0.3.17", "language": "go"}
]
```

### `get_template`

Retrieve a PCD template by name and version.

**Arguments:**
- `name` (required): Template name (e.g., `cli-tool`, `mcp-server`)
- `version` (optional): Semantic version or `latest` (default: `latest`)

**Returns:** Full TemplateRecord with `name`, `version`, `language`, and `content` fields.

**Errors:**
- `-32602`: Unknown template name or version not found
- `-32603`: Store read error

### `list_resources`

List all available PCD resources (templates, prompts, hints).

**Arguments:** none

**Returns:** JSON array of resource records with `uri` and `name` fields (content omitted).

**Resource URIs follow the format:** `pcd://<type>/<name>`
- Templates: `pcd://templates/cli-tool`
- Prompts: `pcd://prompts/interview`, `pcd://prompts/translator`
- Hints: `pcd://hints/cloud-native.go.go-libvirt`

### `read_resource`

Read a PCD resource by URI.

**Arguments:**
- `uri` (required): Resource URI in format `pcd://<type>/<name>`

**Returns:** ResourceRecord with `uri`, `name`, and full `content`.

**Errors:**
- `-32602`: Invalid URI, unknown resource type, or resource not found
- `-32603`: Store read error

### `lint_content`

Validate a PCD specification from string content.

**Arguments:**
- `content` (required): Full Markdown text of the PCD specification
- `filename` (optional): Filename for diagnostic references (default: `spec.md`)

**Returns:** LintResult with `valid` boolean, error/warning counts, and diagnostic array.

**Diagnostics include:**
- `severity`: `"error"` or `"warning"`
- `line`: 1-based line number
- `section`: Section name (e.g., `"META"`, `"BEHAVIOR"`)
- `message`: Human-readable diagnostic message
- `rule`: Rule identifier (e.g., `"RULE-01"`)

**Errors:**
- `-32602`: Filename missing `.md` extension

**Linting rules (RULE-01 through RULE-14):**
- RULE-01: Required META section
- RULE-02: Required TYPES section
- RULE-03: Required BEHAVIOR section
- RULE-04: PRECONDITIONS section presence
- RULE-05: POSTCONDITIONS section presence
- RULE-06: Required INVARIANTS section
- RULE-07: EXAMPLES section presence
- RULE-08: DEPLOYMENT section presence
- RULE-09: META section required fields
- RULE-10: BEHAVIOR block required subsections
- RULE-11: INVARIANT annotations
- RULE-12: EXAMPLES structure
- RULE-13: Version field semantic versioning
- RULE-14: Spec-Schema version validation

### `lint_file`

Validate a PCD specification from a file.

**Arguments:**
- `path` (required): Absolute or relative filesystem path to a `.md` file

**Returns:** LintResult (same as `lint_content`).

**Errors:**
- `-32602`: Missing `.md` extension or file not found
- `-32603`: Filesystem read error

### `get_schema_version`

Get the Spec-Schema version this server was built against.

**Arguments:** none

**Returns:** Semantic version string (e.g., `"0.3.17"`).

## MCP Resources

The server also advertises resources for direct access by MCP clients:

- **Templates:** `pcd://templates/{name}` — dynamic resource template
- **Prompts:** `pcd://prompts/{name}` — dynamic resource template
- **Hints:** `pcd://hints/{key}` — dynamic resource template

## Configuration

The server does not read environment variables for behavior control. Configuration is provided via:

1. **Command-line arguments:**
   - `stdio` or `http` — transport selector (bare word)
   - `listen=host:port` — HTTP listen address (default: `127.0.0.1:8080`)

2. **Template and hints locations (production only):**
   - `/usr/share/pcd/templates/` — system templates
   - `/etc/pcd/` — system configuration
   - `~/.config/pcd/` — user configuration
   - `./.pcd/` — project configuration

3. **Embedded prompts:**
   - Interview prompt: embedded at build time from `prompts/interview-prompt.md`
   - Translator prompt: embedded at build time from `prompts/prompt.md`

## Development

### Building from Source

```bash
make build
```

Produces a static binary `mcp-server-pcd` with no runtime dependencies.

### Building Container Images

```bash
# Build with podman (recommended for SUSE/Linux)
make container

# Test the container
make container-test

# Clean up container images
make container-clean
```

The Containerfile uses a multi-stage build:
1. **Builder stage:** `registry.suse.com/bci/golang:latest` — official SUSE BCI Go image
2. **Final stage:** `FROM scratch` — minimal image containing only the static binary
3. **Result:** ~15 MB image with no runtime dependencies

### Running Tests

```bash
make test
```

All tests use in-memory test doubles (FakeTemplateStore, FakePromptStore, FakeFilesystem). No filesystem access or live services are required.

### Code Organization

```
.
├── main.go                          # Transport wiring and tool handlers
├── go.mod                           # Module definition
├── internal/
│   ├── store/
│   │   ├── store.go                 # Interface definitions and implementations
│   │   └── prompts.go               # Embedded prompt constants
│   └── lint/
│       └── lint.go                  # Linting rule engine
├── independent_tests/
│   └── INDEPENDENT_TESTS.go         # Integration tests
├── Makefile                         # Build targets
├── mcp-server-pcd.spec             # RPM spec
├── debian/                          # Debian packaging
├── Containerfile                    # OCI container build
├── mcp-server-pcd.service          # systemd service unit
└── README.md                        # This file
```

## Error Handling

All errors are returned as JSON-RPC 2.0 error responses:

- `-32602` (Invalid params): Malformed requests, missing required arguments, invalid URIs
- `-32603` (Internal error): Store failures, filesystem errors, unhandled exceptions

The server never panics or crashes on invalid input. All errors are caught and returned as proper MCP error responses.

## Security Considerations

- **Static binary:** No runtime dependencies or dynamic linking
- **No environment variables:** Configuration via arguments only
- **Read-only operations:** The server never modifies files or makes outbound network calls
- **Sandboxing:** Systemd service runs with `ProtectSystem=strict` and `ProtectHome=true`

## License

GNU General Public License v2.0 (GPL-2.0-only)

See LICENSE file for details.

## Author

Matthias G. Eckermann <pcd@mailbox.org>

## Links

- PCD Specification: https://github.com/mge1512/pcd
- MCP Protocol: https://modelcontextprotocol.io/
- OBS (Open Build Service): https://build.opensuse.org/
