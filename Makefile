# ContextKeeper Makefile
# A minimalist CLI tool for managing project context

APP_NAME = contextkeeper
VERSION ?= 0.5.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%d")
DIST_DIR = releases
TMP_DIR = tmp

# Go parameters
GO = go
LDFLAGS = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: all build clean test build-all release install verify help

# Default target
all: build

# Build for current platform
build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .

# Clean build artifacts
clean:
	rm -rf $(DIST_DIR) $(TMP_DIR) $(APP_NAME)

# Run all tests
test:
	$(GO) test ./... -v

# Build all platforms with CORRECT binary names
build-all: clean
	@mkdir -p $(TMP_DIR)
	
	@echo "Building Linux amd64..."
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(TMP_DIR)/ck .
	
	@echo "Building Linux arm64..."
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(TMP_DIR)/ck .
	
	@echo "Building macOS amd64..."
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(TMP_DIR)/ck .
	
	@echo "Building macOS arm64..."
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(TMP_DIR)/ck .
	
	@echo "Building Windows amd64..."
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(TMP_DIR)/ck.exe .
	
	@echo "Creating release packages..."
	@mkdir -p $(DIST_DIR)
	tar -czf $(DIST_DIR)/ck-linux-amd64.tar.gz -C $(TMP_DIR) ck
	tar -czf $(DIST_DIR)/ck-linux-arm64.tar.gz -C $(TMP_DIR) ck
	tar -czf $(DIST_DIR)/ck-darwin-amd64.tar.gz -C $(TMP_DIR) ck
	tar -czf $(DIST_DIR)/ck-darwin-arm64.tar.gz -C $(TMP_DIR) ck
	tar -czf $(DIST_DIR)/ck-windows-amd64.tar.gz -C $(TMP_DIR) ck.exe
	cd $(DIST_DIR) && sha256sum *.tar.gz > checksums.txt
	
	@echo "Build complete! Binaries in $(DIST_DIR)/"
	@echo "Contents:"
	@for f in $(DIST_DIR)/*.tar.gz; do echo "  $$f:"; tar -tzf $$f | head -1; done

# Create release packages (requires build-all first)
release: build-all
	@echo "Release packages ready in $(DIST_DIR)/"

# Install to system (requires sudo)
install: build
	sudo install -d /usr/local/bin
	sudo install -m 755 $(APP_NAME) /usr/local/bin/ck

# Verify build works correctly
verify: build
	@echo "Verifying $(APP_NAME)..."
	@./$(APP_NAME) --help

help:
	@echo "ContextKeeper Makefile:"
	@echo ""
	@echo "  make build         - Build binary for current platform"
	@echo "  make build-all     - Build binaries for ALL platforms with correct names"
	@echo "  make test          - Run all tests"
	@echo "  make release       - Create release packages"
	@echo "  make install       - Install to /usr/local/bin (requires sudo)"
	@echo "  make verify        - Verify build works"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make help          - Show this help"
	@echo ""
	@echo "Version: $(VERSION)"
