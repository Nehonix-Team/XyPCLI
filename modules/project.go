package modules

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io" 
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Template URLs for downloading project templates
const (
	// Nehonix SDK URL for downloading templates
	NehonixSDKURL = "https://sdk.nehonix.space/dl/mds/xypriss/templates/"

	// Local template URL for testing (relative to CLI binary)
	LocalTemplatePath = "initdr.zip"
)

// InitProject initializes a new XyPriss project with all necessary configuration
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
func (c *CLITool) InitProject() {
	fmt.Println(XyPrissLogo)
	fmt.Printf("%s🚀 Initializing new XyPriss project...%s\n\n", ColorGreen, ColorReset)

	// Get project configuration interactively
	config := GetProjectConfig()

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
	err = c.extractTemplate(templatePath, config.Name, config.Language)
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

// downloadTemplate downloads the project template from GitHub releases
// This function detects the platform and downloads the appropriate template
func (c *CLITool) downloadTemplate() (string, error) {
	tempFile, err := ioutil.TempFile("", "xypriss-template-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	// Get platform information
	platformOS, arch, _ := GetPlatformInfo()
	fmt.Printf("  %sDetected platform: %s/%s%s\n", ColorDim, platformOS, arch, ColorReset)

	// Download from Nehonix SDK
	templateURL := NehonixSDKURL + "initdr.zip"
	fmt.Printf("  %sDownloading template from Nehonix SDK...%s\n", ColorDim, ColorReset)

	resp, err := http.Get(templateURL)
	if err != nil {
		// Fallback to local template file for testing
		fmt.Printf("  %s⚠️  Nehonix SDK not available, using local template...%s\n", ColorYellow, ColorReset)
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

		fmt.Printf("  %s✅ Template downloaded from Nehonix SDK%s\n", ColorGreen, ColorReset)
		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to save template: %v", err)
		}
	}

	return tempFile.Name(), nil
}

// extractTemplate extracts the downloaded zip file to the project directory
// This function creates the complete project structure by extracting
// all files from the template zip based on the selected language
func (c *CLITool) extractTemplate(zipPath, projectName, language string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer reader.Close()

	// Determine template path based on language
	templateDir := "TS" // Default to TypeScript
	if language == "js" {
		templateDir = "JS"
	}

	for _, file := range reader.File {
		// Skip files not in the selected language template
		if !strings.HasPrefix(file.Name, templateDir+"/") && file.Name != templateDir {
			continue
		}

		// Remove the language prefix from the file path
		fileName := strings.TrimPrefix(file.Name, templateDir+"/")
		if fileName == "" {
			continue // Skip the directory itself
		}

		filePath := filepath.Join(projectName, fileName)

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