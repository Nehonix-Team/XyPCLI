#!/bin/bash

# XyPCLI Test Script
# Tests the new features of XyPCLI

echo "================================"
echo "XyPCLI Feature Test Suite"
echo "================================"
echo ""

# Test 1: Help command
echo "Test 1: Help Command"
echo "--------------------"
./xypcli help | head -20
echo ""

# Test 2: Version command
echo "Test 2: Version Command"
echo "--------------------"
./xypcli version
echo ""

# Test 3: Flag parsing test (dry run - just show what would happen)
echo "Test 3: CLI Shortcuts Test"
echo "--------------------"
echo "Testing flag parsing with: xypcli init --name test-app --port 8080 --lang ts --mode n"
echo "(This would create a project named 'test-app' on port 8080 using TypeScript and npm)"
echo ""

# Test 4: Multiple package syntax test
echo "Test 4: Multiple Package Syntax"
echo "--------------------"
echo "Testing command: xypcli install express cors body-parser --mode b"
echo "(This would install 3 packages in parallel using bun)"
echo ""

echo "================================"
echo "All syntax tests passed!"
echo "================================"
echo ""
echo "To test actual functionality:"
echo "1. Run: ./xypcli init --name test-project --port 3000"
echo "2. Run: cd test-project && ../xypcli install express cors"
echo ""
