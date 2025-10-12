package modules

import (
	"runtime"
) 

// GetPlatformInfo detects the current platform and returns the appropriate binary name
func GetPlatformInfo() (os string, arch string, binaryName string) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Normalize architecture names
	switch goarch {
	case "amd64":
		goarch = "amd64"
	case "arm64":
		goarch = "arm64"
	case "arm":
		goarch = "arm"
	default:
		goarch = "amd64" // fallback
	}

	// Normalize OS names
	switch goos {
	case "darwin":
		os = "darwin"
		if goarch == "amd64" {
			binaryName = "xypcli-darwin-amd64"
		} else {
			binaryName = "xypcli-darwin-arm64"
		}
	case "linux":
		os = "linux"
		if goarch == "amd64" {
			binaryName = "xypcli-linux-amd64"
		} else {
			binaryName = "xypcli-linux-arm64"
		}
	case "windows":
		os = "windows"
		if goarch == "amd64" {
			binaryName = "xypcli-windows-amd64.exe"
		} else {
			binaryName = "xypcli-windows-arm.exe"
		}
	default:
		os = "linux"
		binaryName = "xypcli-linux-amd64"
	}

	return os, goarch, binaryName
}