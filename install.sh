#!/bin/bash

# XyPriss CLI Installer
# This script downloads and installs the XyPriss CLI tool globally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
SDK_BASE_URL="https://dll.nehonix.com/dl/mds/xypriss/bin"
CLI_NAME="xypcli"
TEMP_DIR="/tmp/xypriss-cli-install"
INSTALL_DIR="${HOME}/.local/bin"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}  XyPriss CLI Installer${NC}"
    echo -e "${CYAN}================================${NC}"
    echo ""
}

# Function to detect platform
detect_platform() {
    local os=""
    local arch=""

    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          print_error "Unsupported OS: $(uname -s)"; exit 1 ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        arm*)        arch="arm" ;;
        *)           print_error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac

    # Special handling for macOS ARM64
    if [ "$os" = "darwin" ] && [ "$arch" = "arm64" ]; then
        binary_name="${CLI_NAME}-darwin-arm64"
    elif [ "$os" = "darwin" ] && [ "$arch" = "amd64" ]; then
        binary_name="${CLI_NAME}-darwin-amd64"
    elif [ "$os" = "linux" ] && [ "$arch" = "amd64" ]; then
        binary_name="${CLI_NAME}-linux-amd64"
    elif [ "$os" = "linux" ] && [ "$arch" = "arm64" ]; then
        binary_name="${CLI_NAME}-linux-arm64"
    elif [ "$os" = "windows" ] && [ "$arch" = "amd64" ]; then
        binary_name="${CLI_NAME}-windows-amd64.exe"
        CLI_NAME="${CLI_NAME}.exe"
    elif [ "$os" = "windows" ] && [ "$arch" = "arm" ]; then
        binary_name="${CLI_NAME}-windows-arm.exe"
        CLI_NAME="${CLI_NAME}.exe"
    else
        print_error "Unsupported platform combination: $os/$arch"
        exit 1
    fi

    echo "$os:$arch:$binary_name"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if curl or wget is available
check_download_tool() {
    if command_exists curl; then
        echo "curl"
    elif command_exists wget; then
        echo "wget"
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Function to download file
download_file() {
    local url="$1"
    local output="$2"
    local tool=$(check_download_tool)

    print_info "Downloading $url..."

    if [ "$tool" = "curl" ]; then
        if ! curl -L -o "$output" "$url" 2>/dev/null; then
            print_error "Failed to download using curl"
            return 1
        fi
    else
        if ! wget -O "$output" "$url" 2>/dev/null; then
            print_error "Failed to download using wget"
            return 1
        fi
    fi

    return 0
}

# Function to create installation directory
create_install_dir() {
    if [ ! -d "$INSTALL_DIR" ]; then
        print_info "Creating installation directory: $INSTALL_DIR"
        if ! mkdir -p "$INSTALL_DIR" 2>/dev/null; then
            print_warning "Cannot create $INSTALL_DIR, trying system directories..."

            # Try system directories
            for dir in "/usr/local/bin" "/usr/bin" "/opt/bin"; do
                if [ -w "$(dirname "$dir")" ] 2>/dev/null || [ -w "$dir" ] 2>/dev/null; then
                    INSTALL_DIR="$dir"
                    print_info "Using system directory: $INSTALL_DIR"
                    break
                fi
            done

            if [ ! -d "$INSTALL_DIR" ]; then
                mkdir -p "$INSTALL_DIR" 2>/dev/null || {
                    print_error "Cannot create installation directory. Please run with sudo or specify a different location."
                    exit 1
                }
            fi
        fi
    fi
}

# Function to check if PATH includes install directory
check_path() {
    case ":$PATH:" in
        *:"$INSTALL_DIR":*)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Function to add to PATH
add_to_path() {
    local shell_rc=""

    # Detect shell
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    else
        shell_rc="$HOME/.profile"
    fi

    if [ -f "$shell_rc" ]; then
        print_info "Adding $INSTALL_DIR to PATH in $shell_rc"

        # Check if already in PATH
        if ! grep -q "export PATH.*$INSTALL_DIR" "$shell_rc" 2>/dev/null; then
            echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$shell_rc"
            print_success "Added to PATH. Please restart your terminal or run: source $shell_rc"
        else
            print_info "PATH already configured in $shell_rc"
        fi
    else
        print_warning "Could not detect shell configuration file. Please manually add $INSTALL_DIR to your PATH."
        print_info "Add this line to your shell configuration:"
        echo "export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

# Function to verify installation
verify_installation() {
    if command_exists "$CLI_NAME"; then
        print_success "XyPriss CLI installed successfully!"
        echo ""
        print_info "Run '$CLI_NAME --version' to verify the installation"
        print_info "Run '$CLI_NAME --help' to see available commands"
        return 0
    else
        print_warning "CLI not found in PATH. You may need to restart your terminal."
        print_info "Try running: $INSTALL_DIR/$CLI_NAME --version"
        return 1
    fi
}

# Main installation function
main() {
    print_header

    print_info "Detecting platform..."
    local platform_info=$(detect_platform)
    IFS=':' read -r os arch binary_name <<< "$platform_info"

    print_info "Platform detected: $os/$arch"
    print_info "Binary to download: $binary_name"

    # Create temp directory
    print_info "Creating temporary directory..."
    rm -rf "$TEMP_DIR"
    mkdir -p "$TEMP_DIR"

    # Download CLI
    local download_url="$SDK_BASE_URL/$binary_name"
    local temp_file="$TEMP_DIR/$CLI_NAME"

    if ! download_file "$download_url" "$temp_file"; then
        print_error "Failed to download XyPriss CLI from $download_url"
        print_info "Please check your internet connection and try again."
        rm -rf "$TEMP_DIR"
        exit 1
    fi

    print_success "Downloaded XyPriss CLI successfully"

    # Make executable
    chmod +x "$temp_file"
    print_info "Made binary executable"

    # Create installation directory
    create_install_dir

    # Install binary
    print_info "Installing to $INSTALL_DIR/$CLI_NAME..."
    if ! mv "$temp_file" "$INSTALL_DIR/$CLI_NAME" 2>/dev/null; then
        print_error "Failed to install binary. Try running with sudo:"
        echo "sudo mv \"$temp_file\" \"$INSTALL_DIR/$CLI_NAME\""
        rm -rf "$TEMP_DIR"
        exit 1
    fi

    # Clean up
    rm -rf "$TEMP_DIR"

    # Check PATH
    if ! check_path; then
        print_warning "$INSTALL_DIR is not in your PATH"
        add_to_path
    fi

    # Verify installation
    echo ""
    if verify_installation; then
        echo ""
        print_info "🎉 Installation completed successfully!"
        print_info "You can now use 'xypcli' to create new XyPriss projects."
        echo ""
        print_info "Example usage:"
        echo "  xypcli init          # Create a new project"
        echo "  xypcli start         # Start development server"
        echo "  xypcli --help        # Show all commands"
    fi
}

# Function to show usage
show_usage() {
    cat << EOF
XyPriss CLI Installer

USAGE:
    ./install.sh [OPTIONS]

OPTIONS:
    -h, --help          Show this help message
    -v, --version       Show version information
    --install-dir DIR   Specify custom installation directory
    --no-path           Don't modify PATH environment variable

DESCRIPTION:
    This script downloads and installs the XyPriss CLI tool globally on your system.

    The CLI will be installed to: $INSTALL_DIR
    Make sure this directory is in your PATH.

SUPPORTED PLATFORMS:
    - Linux (amd64, arm64)
    - macOS (amd64, arm64)
    - Windows (amd64, arm)

EXAMPLES:
    ./install.sh                           # Install with default settings
    ./install.sh --install-dir /usr/local/bin  # Install to custom directory
    sudo ./install.sh                      # Install system-wide (Linux/macOS)

For more information, visit: https://github.com/Nehonix-Team/XyPriss
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -v|--version)
            echo "XyPriss CLI Installer v1.0.0"
            exit 0
            ;;
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --no-path)
            NO_PATH=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Run main installation
main