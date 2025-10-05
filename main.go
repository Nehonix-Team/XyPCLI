package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ANSI color codes for beautiful output
const (
	ColorReset     = "\033[0m"
	ColorRed       = "\033[31m"
	ColorGreen     = "\033[32m"
	ColorYellow    = "\033[33m"
	ColorBlue      = "\033[34m"
	ColorMagenta   = "\033[35m"
	ColorCyan      = "\033[36m"
	ColorWhite     = "\033[37m"
	ColorBold      = "\033[1m"
	ColorDim       = "\033[2m"
)

// XyPriss ASCII art logo
const XyPrissLogo = ColorCyan + `
██╗  ██╗██╗   ██╗██████╗ ██████╗ ██╗███████╗███████╗
╚██╗██╔╝╚██╗ ██╔╝██╔══██╗██╔══██╗██║██╔════╝██╔════╝
 ╚███╔╝  ╚████╔╝ ██████╔╝██████╔╝██║███████╗███████╗
 ██╔██╗   ╚██╔╝  ██╔═══╝ ██╔══██╗██║╚════██║╚════██║
██╔╝ ██╗   ██║   ██║     ██║  ██║██║███████║███████║
╚═╝  ╚═╝   ╚═╝   ╚═╝     ╚═╝  ╚═╝╚═╝╚══════╝╚══════╝
` + ColorReset + ColorBlue + `
            ⚡ High-Performance Node.js Framework ⚡
` + ColorReset

// ProjectConfig holds the configuration for a new XyPriss project
// This struct contains all the necessary information to generate a complete
// XyPriss application with the selected features
type ProjectConfig struct {
	Name        string // Project name (used for directory and package.json)
	Description string // Project description
	Version     string // Initial version (defaults to "1.0.0")
	Port        int    // Server port (defaults to 3000)
	WithAuth    bool   // Include JWT authentication system
	WithUpload  bool   // Include file upload functionality with multer
	WithMulti   bool   // Include multi-server configuration
}

// CLITool represents the XyPriss CLI tool with version information
// This tool provides commands for initializing new projects and managing
// XyPriss applications
type CLITool struct {
	version string // CLI version
}

// Template URLs for downloading project templates
const (
	// GitHub releases URL for downloading platform-specific binaries
	GitHubReleasesURL = "https://github.com/Nehonix-Team/XyPCLI/releases/latest/download/"

	// Local template URL for testing (relative to CLI binary)
	LocalTemplatePath = "initdr.zip"
)

// main is the entry point for the XyPriss CLI tool
// It parses command line arguments and dispatches to the appropriate command handler
// The CLI supports the following commands:
//   - init: Initialize a new XyPriss project
//   - start: Start a XyPriss development server
//   - version: Show CLI version information
//   - help: Show help information
func main() {
	cli := &CLITool{version: "1.0.0"}

	if len(os.Args) < 2 {
		cli.showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		cli.initProject()
	case "start":
		cli.startServer()
	case "version", "-v", "--version":
		fmt.Printf("XyPCLI v%s\n", cli.version)
	case "help", "-h", "--help":
		cli.showHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		cli.showHelp()
	}
}

// showHelp displays the CLI help information with beautiful branding
// This function provides comprehensive usage instructions including:
// - XyPriss ASCII art logo
// - Available commands and their descriptions
// - Usage syntax with colored output
// - Example command invocations
// - Version information
func (c *CLITool) showHelp() {
	fmt.Println(XyPrissLogo)
	fmt.Printf("%sCLI Tool v%s%s\n\n", ColorYellow, c.version, ColorReset)
	fmt.Printf("%sUSAGE:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %sxypcli <command> [options]%s\n", ColorCyan, ColorReset)
	fmt.Println()
	fmt.Printf("%sCOMMANDS:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %sinit%s     Initialize a new XyPriss project with all necessary configuration\n", ColorGreen, ColorReset)
	fmt.Printf("  %sstart%s    Start the XyPriss development server in the current directory\n", ColorGreen, ColorReset)
	fmt.Printf("  %sversion%s  Show CLI version information\n", ColorGreen, ColorReset)
	fmt.Printf("  %shelp%s     Show this help message\n", ColorGreen, ColorReset)
	fmt.Println()
	fmt.Printf("%sEXAMPLES:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %sxypcli init%s                    # Create a new XyPriss project interactively\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli start%s                   # Start development server\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli --version%s               # Show CLI version\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli help%s                    # Show this help\n", ColorMagenta, ColorReset)
	fmt.Println()
	fmt.Printf("%sFor more information, visit: %shttps://github.com/Nehonix-Team/XyPriss%s\n", ColorDim, ColorBlue, ColorReset)
}

// initProject initializes a new XyPriss project with all necessary configuration
// This function performs the following steps:
// 1. Prompts user for project configuration (name, features, etc.)
// 2. Downloads the project template from the remote server
// 3. Extracts the template to create the project structure
// 4. Customizes configuration files based on user selections
// 5. Installs dependencies and provides next steps
//
// The template includes:
// - Complete TypeScript project structure
// - XyPriss server configuration
// - Authentication system (optional)
// - File upload support (optional)
// - Multi-server setup (optional)
// - All necessary dependencies and scripts
func (c *CLITool) initProject() {
	fmt.Println(XyPrissLogo)
	fmt.Printf("%s🚀 Initializing new XyPriss project...%s\n\n", ColorGreen, ColorReset)

	// Get project configuration interactively
	config := c.getProjectConfig()

	// Download template
	fmt.Printf("\n%s📥 Downloading project template...%s\n", ColorBlue, ColorReset)
	templatePath, err := c.downloadTemplate()
	if err != nil {
		fmt.Printf("%s❌ Failed to download template:%s %v\n", ColorRed, ColorReset, err)
		os.Exit(1)
	}
	defer os.Remove(templatePath) // Clean up temp file

	// Extract template
	fmt.Printf("%s📦 Extracting template...%s\n", ColorBlue, ColorReset)
	err = c.extractTemplate(templatePath, config.Name)
	if err != nil {
		fmt.Printf("%s❌ Failed to extract template:%s %v\n", ColorRed, ColorReset, err)
		os.Exit(1)
	}

	// Customize configuration files
	fmt.Printf("%s🔧 Customizing configuration...%s\n", ColorYellow, ColorReset)
	c.customizePackageJson(config)
	c.customizeEnvFile(config)
	c.customizeREADME(config)

	// Install dependencies
	fmt.Printf("%s📦 Installing dependencies...%s\n", ColorBlue, ColorReset)
	c.installDependencies(config.Name)

	fmt.Printf("\n%s✅ Project '%s' initialized successfully!%s\n", ColorGreen, config.Name, ColorReset)
	fmt.Printf("\n%sNext steps:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %scd %s%s\n", ColorCyan, config.Name, ColorReset)
	fmt.Printf("  %snpm run dev%s\n", ColorCyan, ColorReset)
	fmt.Printf("\n%s🎉 Happy coding with XyPriss!%s\n", ColorMagenta, ColorReset)
}

// getProjectConfig interactively collects basic project configuration from the user
// This function prompts the user for:
// - Project name (used for directory and package.json)
// - Project description
//
// Returns a ProjectConfig struct with default features enabled for simplicity
func (c *CLITool) getProjectConfig() ProjectConfig {
	reader := bufio.NewReader(os.Stdin)

	config := ProjectConfig{
		Port:       3000,
		Version:    "1.0.0",
		WithAuth:   true,  // Enable by default for better DX
		WithUpload: true,  // Enable by default for better DX
		WithMulti: false, // Keep simple by default
	}

	// Project name - used for directory name and package.json
	fmt.Printf("%sProject name:%s ", ColorCyan, ColorReset)
	name, _ := reader.ReadString('\n')
	config.Name = strings.TrimSpace(name)
	if config.Name == "" {
		config.Name = "my-xypriss-app"
	}

	// Project description - used in package.json and README
	fmt.Printf("%sDescription:%s ", ColorCyan, ColorReset)
	desc, _ := reader.ReadString('\n')
	config.Description = strings.TrimSpace(desc)
	if config.Description == "" {
		config.Description = "A XyPriss application"
	}

	return config
}

// getPlatformInfo detects the current platform and returns the appropriate binary name
func getPlatformInfo() (os string, arch string, binaryName string) {
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

// downloadTemplate downloads the project template from GitHub releases
// This function detects the platform and downloads the appropriate template
func (c *CLITool) downloadTemplate() (string, error) {
	tempFile, err := ioutil.TempFile("", "xypriss-template-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	// Get platform information
	platformOS, arch, _ := getPlatformInfo()
	fmt.Printf("  %sDetected platform: %s/%s%s\n", ColorDim, platformOS, arch, ColorReset)

	// Try to download from GitHub releases first
	templateURL := GitHubReleasesURL + "initdr.zip"
	fmt.Printf("  %sTrying to download template from GitHub releases...%s\n", ColorDim, ColorReset)

	resp, err := http.Get(templateURL)
	if err != nil {
		// Fallback to local template file for testing
		fmt.Printf("  %s⚠️  GitHub releases not available, using local template...%s\n", ColorYellow, ColorReset)
		localTemplate, err := os.Open(LocalTemplatePath)
		if err != nil {
			return "", fmt.Errorf("failed to open local template: %v", err)
		}
		defer localTemplate.Close()

		_, err = io.Copy(tempFile, localTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to copy local template: %v", err)
		}
		fmt.Printf("  %s✅ Local template loaded successfully%s\n", ColorGreen, ColorReset)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to download template: HTTP %d", resp.StatusCode)
		}

		fmt.Printf("  %s✅ Template downloaded from GitHub releases%s\n", ColorGreen, ColorReset)
		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to save template: %v", err)
		}
	}

	return tempFile.Name(), nil
}

// extractTemplate extracts the downloaded zip file to the project directory
// This function creates the complete project structure by extracting
// all files from the template zip
func (c *CLITool) extractTemplate(zipPath, projectName string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(projectName, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// Create directory if it doesn't exist
		dir := filepath.Dir(filePath)
		os.MkdirAll(dir, os.ModePerm)

		// Extract file
		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", filePath, err)
		}

		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return fmt.Errorf("failed to open file in zip %s: %v", file.Name, err)
		}

		_, err = io.Copy(destFile, srcFile)
		destFile.Close()
		srcFile.Close()
		if err != nil {
			return fmt.Errorf("failed to extract file %s: %v", file.Name, err)
		}
	}

	return nil
}

// customizePackageJson modifies the extracted package.json file
// This function updates the package.json with project-specific information:
// - Project name and description
// - Adds optional dependencies based on selected features
// - Maintains the template structure while customizing for the project
func (c *CLITool) customizePackageJson(config ProjectConfig) {
	packagePath := filepath.Join(config.Name, "package.json")

	// Read existing package.json
	data, err := ioutil.ReadFile(packagePath)
	if err != nil {
		log.Printf("Warning: Could not read package.json: %v", err)
		return
	}

	var packageJson map[string]interface{}
	if err := json.Unmarshal(data, &packageJson); err != nil {
		log.Printf("Warning: Could not parse package.json: %v", err)
		return
	}

	// Update project information
	packageJson["name"] = strings.ToLower(strings.ReplaceAll(config.Name, " ", "-"))
	packageJson["description"] = config.Description

	// Add optional dependencies
	dependencies := packageJson["dependencies"].(map[string]interface{})

	if config.WithUpload {
		dependencies["multer"] = "^2.0.2"
	}
	if config.WithAuth {
		dependencies["jsonwebtoken"] = "^9.0.2"
	}

	// Write back to file
	updatedData, _ := json.MarshalIndent(packageJson, "", "  ")
	ioutil.WriteFile(packagePath, updatedData, 0644)
}

// customizeEnvFile modifies the extracted .env file with project-specific settings
// This function updates environment variables like PORT based on user configuration
func (c *CLITool) customizeEnvFile(config ProjectConfig) {
	envPath := filepath.Join(config.Name, ".env")

	// Read existing .env file
	data, err := ioutil.ReadFile(envPath)
	if err != nil {
		log.Printf("Warning: Could not read .env file: %v", err)
		return
	}

	envContent := string(data)

	// Replace PORT if it exists
	envContent = strings.ReplaceAll(envContent, "PORT=8080", fmt.Sprintf("PORT=%d", config.Port))

	// Write back to file
	ioutil.WriteFile(envPath, []byte(envContent), 0644)
}

// customizeREADME modifies the extracted README.md file with project-specific information
// This function updates the README with the correct project name, description, and features
func (c *CLITool) customizeREADME(config ProjectConfig) {
	readmePath := filepath.Join(config.Name, "README.md")

	// Read existing README
	data, err := ioutil.ReadFile(readmePath)
	if err != nil {
		log.Printf("Warning: Could not read README.md: %v", err)
		return
	}

	readmeContent := string(data)

	// Replace placeholders
	readmeContent = strings.ReplaceAll(readmeContent, "{{PROJECT_NAME}}", config.Name)
	readmeContent = strings.ReplaceAll(readmeContent, "{{PROJECT_DESCRIPTION}}", config.Description)
	readmeContent = strings.ReplaceAll(readmeContent, "{{PORT}}", fmt.Sprintf("%d", config.Port))

	// Add feature indicators
	features := ""
	if config.WithAuth {
		features += "- 🔐 **Authentication** - JWT-based authentication\n"
	}
	if config.WithUpload {
		features += "- 📁 **File Upload** - Support for file uploads\n"
	}
	if config.WithMulti {
		features += "- 🌐 **Multi-Server** - Multiple server instances\n"
	}

	if features != "" {
		readmeContent = strings.ReplaceAll(readmeContent, "{{FEATURES}}", features)
	} else {
		readmeContent = strings.ReplaceAll(readmeContent, "{{FEATURES}}", "")
	}

	// Write back to file
	ioutil.WriteFile(readmePath, []byte(readmeContent), 0644)
}

// installDependencies runs npm install in the project directory
// This function installs all the dependencies defined in package.json
func (c *CLITool) installDependencies(projectName string) {
	cmd := exec.Command("npm", "install")
	cmd.Dir = projectName
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Warning: Failed to install dependencies: %v", err)
		fmt.Println("⚠️  You may need to run 'npm install' manually in the project directory")
	}
}

func (c *CLITool) startServer() {
	fmt.Println(XyPrissLogo)
	fmt.Printf("%s🚀 Starting XyPriss development server...%s\n\n", ColorGreen, ColorReset)

	// Check if package.json exists
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Printf("%s❌ No package.json found.%s Are you in a XyPriss project directory?\n", ColorRed, ColorReset)
		fmt.Printf("   Run %s'xypcli init'%s to create a new project.\n", ColorCyan, ColorReset)
		return
	}

	// Check if src/server.ts exists
	if _, err := os.Stat("src/server.ts"); os.IsNotExist(err) {
		fmt.Printf("%s❌ No src/server.ts found.%s Are you in a XyPriss project directory?\n", ColorRed, ColorReset)
		fmt.Printf("   Run %s'xypcli init'%s to create a new project.\n", ColorCyan, ColorReset)
		return
	}

	// Check if node_modules exists
	if _, err := os.Stat("node_modules"); os.IsNotExist(err) {
		fmt.Printf("%s📦 Installing dependencies...%s\n", ColorBlue, ColorReset)
		installCmd := exec.Command("npm", "install")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Printf("%s❌ Failed to install dependencies:%s %v\n", ColorRed, ColorReset, err)
			return
		}
	}

	// Start the server
	fmt.Printf("%s🔥 Starting development server...%s\n", ColorYellow, ColorReset)
	fmt.Printf("%sPress Ctrl+C to stop the server%s\n\n", ColorDim, ColorReset)

	cmd := exec.Command("npm", "run", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n%s❌ Failed to start server:%s %v\n", ColorRed, ColorReset, err)
	}
}
