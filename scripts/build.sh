#!/bin/bash

# XyPCLI Build Script
# This script builds the CLI binaries for multiple platforms and creates template packages

set -e

echo "🚀 Building XyPCLI..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
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

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go first."
    exit 1
fi

print_status "Go version: $(go version)"

# Create bin directory
print_status "Creating bin directory..."
mkdir -p bin

# Build for multiple platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm"
)

print_status "Building binaries for multiple platforms..."

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a parts <<< "$platform"
    GOOS="${parts[0]}"
    GOARCH="${parts[1]}"

    binary_name="xypcli"
    if [ "$GOOS" = "windows" ]; then
        binary_name="xypcli.exe"
    fi

    output_name="bin/${binary_name}"

    if [ "$GOOS" = "windows" ] && [ "$GOARCH" = "amd64" ]; then
        output_name="bin/xypcli-windows-amd64.exe"
    elif [ "$GOOS" = "windows" ] && [ "$GOARCH" = "arm" ]; then
        output_name="bin/xypcli-windows-arm.exe"
    elif [ "$GOOS" = "linux" ] && [ "$GOARCH" = "amd64" ]; then
        output_name="bin/xypcli-linux-amd64"
    elif [ "$GOOS" = "linux" ] && [ "$GOARCH" = "arm64" ]; then
        output_name="bin/xypcli-linux-arm64"
    elif [ "$GOOS" = "darwin" ] && [ "$GOARCH" = "amd64" ]; then
        output_name="bin/xypcli-darwin-amd64"
    elif [ "$GOOS" = "darwin" ] && [ "$GOARCH" = "arm64" ]; then
        output_name="bin/xypcli-darwin-arm64"
    fi

    print_status "Building for $GOOS/$GOARCH..."

    # Set environment variables for cross-compilation
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    export CGO_ENABLED=0

    # Build the binary
    if go build -ldflags="-s -w" -o "$output_name" .; then
        print_success "Built $output_name"
    else
        print_error "Failed to build for $GOOS/$GOARCH"
        exit 1
    fi
done

# Copy the appropriate binary for the current platform to bin/xypcli
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if [[ $(uname -m) == "x86_64" ]]; then
        cp bin/xypcli-linux-amd64 bin/xypcli
    elif [[ $(uname -m) == "aarch64" ]]; then
        cp bin/xypcli-linux-arm64 bin/xypcli
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ $(uname -m) == "x86_64" ]]; then
        cp bin/xypcli-darwin-amd64 bin/xypcli
    elif [[ $(uname -m) == "arm64" ]]; then
        cp bin/xypcli-darwin-arm64 bin/xypcli
    fi
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    cp bin/xypcli-windows-amd64.exe bin/xypcli.exe
fi

print_status "Building templates package..."

# Remove existing zip file if it exists
rm -f initdr.zip

# Create zip file excluding node_modules and other unwanted files
cd templates
if zip -r ../initdr.zip . \
    -x "node_modules/*" \
    -x "*.log" \
    -x ".DS_Store" \
    -x "*/.git/*" \
    -x "*/node_modules/*"; then
    print_success "Templates zip file created"
else
    print_error "Failed to create templates zip"
    exit 1
fi

cd ..

print_success "✅ XyPCLI build completed successfully!"
echo ""
echo "📦 Generated files:"
ls -lh bin/
ls -lh initdr.zip
echo ""
print_status "Ready for publishing to npm!"
print_status "Use: npm publish"