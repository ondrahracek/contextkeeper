# ContextKeeper Makefile
# A minimalist CLI tool for managing project context

BINARY_NAME = ck
OS := $(shell go env GOOS)
ifeq ($(OS),windows)
	BINARY_NAME = ck.exe
endif
VERSION ?= 0.1.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%d")
DIST_DIR = releases

# Go parameters
GO = go
GOBIN = $(shell go env GOPATH)/bin
LDFLAGS = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Directories
SRC = cmd internal
BUILD_DIR = bin

.PHONY: all build clean test test-coverage build-all release install install-homebrew verify version

# Default target
all: build

# Build for current platform
build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

# Build with version info
version: VERSION = $(VERSION)
version:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR) $(BINARY_NAME) $(DIST_DIR)

# Run all tests
test:
	$(GO) test ./... -v

# Run tests with coverage
test-coverage:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

# Build all platforms
build-all: clean
	@echo "Building Linux amd64..."
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/ck-linux-amd64 .
	@echo "Building Linux arm64..."
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/ck-linux-arm64 .
	@echo "Building macOS amd64..."
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/ck-darwin-amd64 .
	@echo "Building macOS arm64..."
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/ck-darwin-arm64 .
	@echo "Building Windows amd64..."
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/ck-windows-amd64.exe .
	@echo "Building source tarball..."
	tar -czf $(DIST_DIR)/ck-$(VERSION)-src.tar.gz --transform "s|^|ck-$(VERSION)/|" cmd/ internal/ Makefile go.mod go.sum README.md LICENSE
	@echo "Build complete! Binaries in $(DIST_DIR)/"

# Create release package (tar.gz for each platform)
release: build-all
	@for f in $(DIST_DIR)/ck-linux-* $(DIST_DIR)/ck-darwin-*; do \
		if [ -f "$$f" ]; then \
			echo "Compressing $$f..."; \
			tar -czf "$$f.tar.gz" "$$f"; \
		fi; \
	done
	@echo "Release packages created in $(DIST_DIR)/"

# Install to system (requires sudo)
install:
	sudo install -d /usr/local/bin
	sudo install -m 755 $(BINARY_NAME) /usr/local/bin/ck

# Create Homebrew bottle and install
install-homebrew: build-all
	@echo "Creating Homebrew bottle..."
	# Create a bottle for each platform
	for f in $(DIST_DIR)/$(BINARY_NAME)-darwin-*; do \
		if [ -f "$$f" ]; then \
			tar -czf "$$f.bottle.tar.gz" "$$f"; \
		fi; \
	done
	@echo "Bottles created. Install with: brew install ./homebrew/contextkeeper.rb"

# Verify build works correctly
verify: build
	@echo "Verifying $(BINARY_NAME)..."
	@./$(BINARY_NAME) --help
	@echo ""
	@echo "Build verification: SUCCESS"

# Show version info
version-info:
	@echo "ContextKeeper $(VERSION)"
	@echo "Git Commit: $(COMMIT)"
	@echo "Build Date: $(DATE)"
	@echo "Go Version: $(shell go version)"
	@echo "Platform: $(shell go env GOOS)/$(shell go env GOARCH)"

# Lint code (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	$(GO) fmt ./...

# Download dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

help:
	@echo "ContextKeeper Makefile Targets:"
	@echo ""
	@echo "  make build              - Build binary for current platform"
	@echo "  make build-all          - Build binaries for all platforms"
	@echo "  make test               - Run all tests"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make release            - Create release packages"
	@echo "  make install            - Install to /usr/local/bin (requires sudo)"
	@echo "  make verify             - Verify build works correctly"
	@echo "  make fmt                - Format code"
	@echo "  make lint               - Lint code (requires golangci-lint)"
	@echo "  make help               - Show this help message"
	@echo ""
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
