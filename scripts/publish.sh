#!/bin/bash

# XyPCLI Publishing Script
# This script handles the publishing process for XyPCLI to npm

set -e

echo "🚀 Publishing XyPCLI to npm..."

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

# Check if npm is logged in
if ! npm whoami &> /dev/null; then
    print_error "You are not logged in to npm. Please run 'npm login' first."
    exit 1
fi

print_status "Logged in as: $(npm whoami)"

# Check if we're in the right directory
if [ ! -f "package.json" ] || [ ! -f "main.go" ]; then
    print_error "This script must be run from the XyPCLI directory."
    exit 1
fi

# Run build script
print_status "Running build script..."
if [ -f "build.sh" ]; then
    chmod +x build.sh
    ./build.sh
else
    print_error "build.sh not found!"
    exit 1
fi

# Check if install.js exists
if [ ! -f "install.js" ]; then
    print_error "install.js not found."
    exit 1
fi

# Check if templates zip exists
if [ ! -f "initdr.zip" ]; then
    print_error "Templates zip not found. Build may have failed."
    exit 1
fi

# Get version from package.json
VERSION=$(node -p "require('./package.json').version")
print_status "Publishing version: $VERSION"

# Ask for confirmation
echo ""
print_warning "About to publish xypriss-cli@$VERSION to npm"
echo ""
echo "📦 Package contents:"
echo "  - Auto-installer script (install.js)"
echo "  - Templates package (initdr.zip)"
echo "  - Package metadata"
echo ""
echo "📦 Binaries will be downloaded on-demand from GitHub releases"
echo ""

read -p "Do you want to continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_status "Publishing cancelled."
    exit 0
fi

# Determine if this is a beta or stable release
if [[ $VERSION == *"beta"* ]] || [[ $VERSION == *"alpha"* ]] || [[ $VERSION == *"rc"* ]]; then
    TAG="beta"
    print_status "Publishing as beta release (tag: beta)"
else
    TAG="latest"
    print_status "Publishing as stable release (tag: latest)"
fi

# Publish to npm
print_status "Publishing to npm with tag: $TAG"
if npm publish --tag "$TAG"; then
    print_success "✅ Successfully published xypriss-cli@$VERSION to npm!"
    echo ""
    print_status "Users can now install with:"
    echo "  npm install -g xypriss-cli"
    echo "  # or"
    echo "  npx xypriss-cli"
    echo ""
    print_status "Don't forget to create a GitHub release with the binaries:"
    echo "  https://github.com/Nehonix-Team/XyPCLI/releases"
else
    print_error "❌ Failed to publish to npm"
    exit 1
fi

print_success "🎉 Publishing process completed!"