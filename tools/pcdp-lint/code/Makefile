.PHONY: build test install clean

# Build configuration
BINARY_NAME = pcdp-lint
BUILD_DIR = build
CGO_ENABLED = 0
LDFLAGS = -s -w

# Default target
all: build

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run tests
test:
	go test -v ./...

# Install the binary to system
install: build
	install -m 755 $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Development targets
fmt:
	go fmt ./...

vet:
	go vet ./...

# Static analysis
lint: fmt vet
	@echo "Static analysis complete"