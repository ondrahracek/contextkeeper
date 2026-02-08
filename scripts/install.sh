#!/bin/bash
# ContextKeeper Installation Script
# Installs ContextKeeper CLI tool to your system

set -e

# Configuration
BINARY_NAME="contextkeeper"
REPO="ondrahracek/contextkeeper"
VERSION="${1:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${GREEN}INFO:${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}WARNING:${NC} $1"
}

echo_error() {
    echo -e "${RED}ERROR:${NC} $1"
}

# Detect OS and Architecture
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        *)
            echo_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            echo_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    echo "$OS-$ARCH"
}

# Get latest version from GitHub
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        VERSION=$(curl -sL https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name"' | sed 's/.*": "\(.*\)".*/\1/' || echo "")
    fi
    
    if [ -z "$VERSION" ]; then
        echo_error "Could not determine latest version"
        exit 1
    fi
    
    echo "$VERSION"
}

# Download and install binary
install_binary() {
    local os_arch="$1"
    local version="$2"
    local tmp_dir=$(mktemp -d)
    local binary_url="https://github.com/$REPO/releases/download/$version/$BINARY_NAME-$os_arch"
    local binary_path="$tmp_dir/$BINARY_NAME"
    
    echo_info "Downloading $BINARY_NAME $version for $os_arch..."
    
    if command -v curl >/dev/null 2>&1; then
        curl -sL "$binary_url" -o "$binary_path"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$binary_url" -O "$binary_path"
    else
        echo_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    if [ ! -f "$binary_path" ]; then
        echo_error "Download failed. Binary not found at: $binary_url"
        exit 1
    fi
    
    chmod +x "$binary_path"
    
    echo_info "Installing to $INSTALL_DIR/$BINARY_NAME..."
    sudo mkdir -p "$INSTALL_DIR"
    sudo cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    
    # Cleanup
    rm -rf "$tmp_dir"
    
    echo_info "Successfully installed $BINARY_NAME $version to $INSTALL_DIR/$BINARY_NAME"
}

# Verify installation
verify_install() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo_info "Verifying installation..."
        $BINARY_NAME --version || true
        echo ""
        echo_info "Installation successful!"
    else
        echo_warn "Binary installed but not found in PATH. You may need to restart your shell or add $INSTALL_DIR to your PATH."
    fi
}

# Main
main() {
    echo "=============================================="
    echo "  ContextKeeper Installation Script"
    echo "=============================================="
    echo ""
    
    if [ "$(id -u)" -eq 0 ]; then
        echo_warn "Running as root. Install may require additional steps."
    fi
    
    local os_arch=$(detect_os)
    local version=$(get_latest_version)
    
    echo "OS/Arch: $os_arch"
    echo "Version: $version"
    echo "Install Directory: $INSTALL_DIR"
    echo ""
    
    install_binary "$os_arch" "$version"
    verify_install
    
    echo ""
    echo_info "Usage: $BINARY_NAME add 'Your context note'"
    echo_info "       $BINARY_NAME list"
    echo ""
}

main "$@"
