package modules

import (
	"archive/zip"
	"bytes"
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
	"sync"
	"time"
)

// npmMutex ensures only one npm install process runs at a time in the same directory
var npmMutex sync.Mutex

// Template URLs for downloading project templates
const (
    NehonixSDKURL = "https://dll.nehonix.com/dl/mds/xypriss/templates/" // production
	// NehonixSDKURL     = "http://127.0.0.1:5500/tools/XyPCLI/"
	LocalTemplatePath = "initdr.zip"
)

// InitProject initializes a new XyPriss project with all necessary configuration
func (c *CLITool) InitProject(flags InitFlags) {
	fmt.Println(XyPrissLogo)
	
	// Show loading animation
	stop := c.showInlineSpinner("🚀 Initializing new XyPriss project...")
	time.Sleep(500 * time.Millisecond)
	c.clearInlineSpinner(stop)
	fmt.Printf("🚀 %sInitializing new XyPriss project...%s\n\n", ColorGreen, ColorReset)

	// Get project configuration interactively or from flags
	config := GetProjectConfig(flags)
	
	// Display configuration in tree format
	c.displayProjectConfig(config)

	// Display configuration in tree format
	c.displayProjectConfig(config)

	// Download template with animation
	fmt.Printf("\n%s┌─────────────────────────────────────────┐%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s│  📥 Downloading project template...    │%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s└─────────────────────────────────────────┘%s\n", ColorBlue, ColorReset)
	
	templatePath, err := c.downloadTemplate()
	if err != nil {
		fmt.Printf("\n%s✗ Failed to download template:%s %v\n", ColorRed, ColorReset, err)
		os.Exit(1)
	}
	defer os.Remove(templatePath)

	// Extract template with animation
	fmt.Printf("\n%s┌─────────────────────────────────────────┐%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s│  📦 Extracting template...             │%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s└─────────────────────────────────────────┘%s\n", ColorBlue, ColorReset)
	
	err = c.extractTemplate(templatePath, config.Name, config.Language)
	if err != nil {
		fmt.Printf("\n%s✗ Failed to extract template:%s %v\n", ColorRed, ColorReset, err)
		os.Exit(1)
	}
	fmt.Printf("  %s✓ Template extracted successfully%s\n", ColorGreen, ColorReset)

	// Customize configuration
	fmt.Printf("\n%s┌─────────────────────────────────────────┐%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s│  🔧 Customizing configuration...       │%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s└─────────────────────────────────────────┘%s\n", ColorYellow, ColorReset)
	
	c.customizePackageJson(config)
	fmt.Printf("  %s✓ package.json configured%s\n", ColorGreen, ColorReset)
	
	c.customizeEnvFile(config)
	fmt.Printf("  %s✓ .env file configured%s\n", ColorGreen, ColorReset)
	
	c.createConfigFile(config)
	fmt.Printf("  %s✓ xypriss.config.json created%s\n", ColorGreen, ColorReset)
	
	c.customizeREADME(config)
	fmt.Printf("  %s✓ README.md configured%s\n", ColorGreen, ColorReset)

	// Install dependencies with tree format
	fmt.Printf("\n%s📦 Installing dependencies...%s\n", ColorMagenta, ColorReset)
	c.installDependencies(config.Name, config.Language, flags.Mode, flags.Strict)

	// Success message with beautiful formatting
	fmt.Printf("\n%s╔═════════════════════════════════════════╗%s\n", ColorGreen, ColorReset)
	fmt.Printf("%s║  ✨ Project '%s' initialized!          ║%s\n", ColorGreen, config.Name, ColorReset)
	fmt.Printf("%s╚═════════════════════════════════════════╝%s\n", ColorGreen, ColorReset)
	
	fmt.Printf("\n%s📋 Next steps:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %s1.%s %scd %s%s\n", ColorCyan, ColorReset, ColorDim, config.Name, ColorReset)
	fmt.Printf("  %s2.%s %snpm run dev%s\n", ColorCyan, ColorReset, ColorDim, ColorReset)
	fmt.Printf("\n%s🎉 Happy coding with XyPriss!%s\n\n", ColorMagenta, ColorReset)
}

// InstallPackage installs a single package using the XyPriss installation system
func (c *CLITool) InstallPackage(packageName string) {
	fmt.Printf("%s📦 Installing dependencies...%s\n", ColorMagenta, ColorReset)

	// Check if we're in a XyPriss project directory
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Printf("  %s✗ No package.json found in current directory%s\n", ColorRed, ColorReset)
		fmt.Printf("%sMake sure you're in a XyPriss project directory%s\n", ColorYellow, ColorReset)
		return
	}

	// Check if Bun is available for faster installation
	useBun := false
	if _, err := exec.LookPath("bun"); err == nil {
		useBun = true
		fmt.Printf("\n  %s⚡ Using 'BMode' for faster installation%s\n", ColorCyan, ColorReset)
	} else {
		fmt.Printf("\n  %s→ Bun not found, using npm%s\n", ColorYellow, ColorReset)
		// Check npm availability
		if _, err := exec.LookPath("npm"); err != nil {
			fmt.Printf("  %s✗ npm is not installed%s\n", ColorRed, ColorReset)
			return
		}
	}

	// Use the same style as installDependencies
	fmt.Printf("%s│%s\n", ColorDim, ColorReset)
	fmt.Printf("%s└─ Package (1)%s\n", ColorDim, ColorReset)
	
	// Install the single package using the existing system
	var failedDeps []string
	c.installSingleDependency(".", packageName, false, useBun, 1, 1, &failedDeps, true, false)

	// Final summary
	fmt.Printf("\n")
	if len(failedDeps) > 0 {
		fmt.Printf("%s⚠ Installation completed with warnings%s\n", ColorYellow, ColorReset)
		fmt.Printf("%s├─ Failed: %d/%d packages%s\n", ColorDim, len(failedDeps), 1, ColorReset)
		for _, dep := range failedDeps {
			prefix := "└─"
			fmt.Printf("%s%s ✗ %s%s\n", ColorDim, prefix, dep, ColorReset)
		}
	} else {
		fmt.Printf("%s✨ Package installed successfully!%s\n", ColorGreen, ColorReset)
		fmt.Printf("%s└─ 1/1 packages%s\n", ColorDim, 1, ColorReset)
	}
}

// InstallPackages installs multiple packages using the XyPriss installation system with intelligent parallelization
func (c *CLITool) InstallPackages(packages []string, mode string) {
	fmt.Printf("%s📦 Installing %d package(s)...%s\n", ColorMagenta, len(packages), ColorReset)

	// Check if we're in a XyPriss project directory
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Printf("  %s✗ No package.json found in current directory%s\n", ColorRed, ColorReset)
		fmt.Printf("%sMake sure you're in a XyPriss project directory%s\n", ColorYellow, ColorReset)
		return
	}

	// Determine installation mode
	useBun := false
	if mode == "b" {
		// Force bun mode
		if _, err := exec.LookPath("bun"); err == nil {
			useBun = true
			fmt.Printf("\n  %s⚡ Using 'BMode' (forced)%s\n", ColorCyan, ColorReset)
		} else {
			fmt.Printf("\n  %s✗ Bun not found, falling back to npm%s\n", ColorRed, ColorReset)
			if _, err := exec.LookPath("npm"); err != nil {
				fmt.Printf("  %s✗ npm is not installed%s\n", ColorRed, ColorReset)
				return
			}
		}
	} else if mode == "n" {
		// Force npm mode
		fmt.Printf("\n  %s→ Using npm (forced)%s\n", ColorCyan, ColorReset)
		if _, err := exec.LookPath("npm"); err != nil {
			fmt.Printf("  %s✗ npm is not installed%s\n", ColorRed, ColorReset)
			return
		}
	} else {
		// Auto-detect mode
		if _, err := exec.LookPath("bun"); err == nil {
			useBun = true
			fmt.Printf("\n  %s⚡ Using 'BMode' for faster installation%s\n", ColorCyan, ColorReset)
		} else {
			fmt.Printf("\n  %s→ Bun not found, using npm%s\n", ColorYellow, ColorReset)
			// Check npm availability
			if _, err := exec.LookPath("npm"); err != nil {
				fmt.Printf("  %s✗ npm is not installed%s\n", ColorRed, ColorReset)
				return
			}
		}
	}

	// Use parallelization for faster installation
	fmt.Printf("  %s⚡ Parallel installation enabled%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s│%s\n", ColorDim, ColorReset)
	fmt.Printf("%s└─ Packages (%d)%s\n", ColorDim, len(packages), ColorReset)

	// Install packages in parallel with intelligent batching
	totalPackages := len(packages)
	failedDeps := make([]string, 0)
	
	// Use channels for parallel installation
	type installResult struct {
		packageName string
		success     bool
		index       int
	}
	
	// Limit concurrent installations to avoid overwhelming the system
	maxConcurrent := 4
	if totalPackages < maxConcurrent {
		maxConcurrent = totalPackages
	}
	
	results := make(chan installResult, totalPackages)
	semaphore := make(chan struct{}, maxConcurrent)
	
	// Start all installations
	for i, pkg := range packages {
		go func(index int, packageName string) {
			semaphore <- struct{}{} // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore
			
			success := c.installPackageParallel(".", packageName, useBun, index+1, totalPackages)
			results <- installResult{packageName: packageName, success: success, index: index}
		}(i, pkg)
	}
	
	// Collect results
	completed := 0
	for completed < totalPackages {
		result := <-results
		if !result.success {
			failedDeps = append(failedDeps, result.packageName)
		}
		completed++
	}

	// Final summary
	fmt.Printf("\n")
	if len(failedDeps) > 0 {
		fmt.Printf("%s⚠ Installation completed with warnings%s\n", ColorYellow, ColorReset)
		fmt.Printf("%s├─ Failed: %d/%d packages%s\n", ColorDim, len(failedDeps), totalPackages, ColorReset)
		for i, dep := range failedDeps {
			prefix := "├─"
			if i == len(failedDeps)-1 {
				prefix = "└─"
			}
			fmt.Printf("%s%s ✗ %s%s\n", ColorDim, prefix, dep, ColorReset)
		}
	} else {
		fmt.Printf("%s✨ All packages installed successfully!%s\n", ColorGreen, ColorReset)
		fmt.Printf("%s└─ %d/%d packages%s\n", ColorDim, totalPackages, totalPackages, ColorReset)
	}
}

// installPackageParallel installs a single package in parallel mode
func (c *CLITool) installPackageParallel(projectDir, packageName string, useBun bool, current, total int) bool {
	var cmd *exec.Cmd
	
	// Special case: nquickdev needs npm for postinstall scripts (Bun ignores them)
	if useBun && packageName == "nquickdev" {
		cmd = exec.Command("npm", "install", packageName)
	} else if useBun {
		cmd = exec.Command("bun", "add", packageName)
	} else {
		cmd = exec.Command("npm", "install", packageName)
	}
	cmd.Dir = projectDir

	// Progress indicator
	progress := fmt.Sprintf("[%d/%d]", current, total)

	// Show inline progress
	fmt.Printf("   %s├─ %s%s %s⚙%s Installing %s...%s\n", 
		ColorDim, progress, ColorReset, ColorCyan, ColorReset, packageName, ColorReset)

	// Concurrent npm installs in the same directory cause race conditions (e.g., ENOTEMPTY, ENOENT)
	// We use a mutex for npm to ensure stability while maintaining the goroutine/channel architecture
	if !useBun {
		npmMutex.Lock()
		defer npmMutex.Unlock()
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("   %s├─ %s%s %s✗%s %s (failed)%s\n", 
			ColorDim, progress, ColorReset, ColorRed, ColorReset, packageName, ColorReset)
		
		// Show detailed error messages
		errOutput := stderr.String()
		if errOutput != "" {
			// Extract and display all relevant error lines (prioritize actual errors)
			lines := strings.Split(errOutput, "\n")
			errorLines := []string{}
			
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					// Collect all error-related lines, but EXCLUDE warnings
					// npm warnings are not fatal errors and shouldn't be highlighted as reasons for failure
					if (strings.Contains(line, "ERR!") || 
					   strings.Contains(line, "error") || 
					   strings.Contains(line, "404") ||
					   strings.Contains(line, "ENOENT") ||
					   strings.Contains(line, "ENOTEMPTY") ||
					   strings.Contains(line, "code")) &&
					   !strings.Contains(strings.ToLower(line), "warn") {
						errorLines = append(errorLines, line)
					}
				}
			}
			
			// If we didn't find specific error lines but the command still failed,
			// fallback to showing the first few lines of stderr
			if len(errorLines) == 0 && errOutput != "" {
				for _, line := range strings.Split(errOutput, "\n") {
					line = strings.TrimSpace(line)
					if line != "" && len(errorLines) < 3 {
						errorLines = append(errorLines, line)
					}
				}
			}
			
			// Display relevant error lines
			maxLines := 5
			if len(errorLines) > maxLines {
				errorLines = errorLines[:maxLines]
			}
			
			for _, errLine := range errorLines {
				fmt.Printf("   %s│  %s→ %s%s\n", ColorDim, ColorYellow, errLine, ColorReset)
			}
		}
		return false
	}

	fmt.Printf("   %s├─ %s%s %s✓%s %s%s\n", 
		ColorDim, progress, ColorReset, ColorGreen, ColorReset, packageName, ColorReset)
	return true
}



// downloadTemplate downloads the project template
func (c *CLITool) downloadTemplate() (string, error) {
	tempFile, err := ioutil.TempFile("", "xypriss-template-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	platformOS, arch, _ := GetPlatformInfo()
	fmt.Printf("  %s→ Platform: %s/%s%s\n", ColorDim, platformOS, arch, ColorReset)

	templateURL := NehonixSDKURL + "initdr.zip"
	fmt.Printf("  %s→ Source: dll.nehonix.com%s\n", ColorDim, ColorReset)

	// Show spinner during download
	stop := c.showInlineSpinner("Downloading...")
	resp, err := http.Get(templateURL)
	c.clearInlineSpinner(stop)

	if err != nil {
		fmt.Printf("  %s⚠ Nehonix SDK unavailable, using local template%s\n", ColorYellow, ColorReset)
		localTemplate, err := os.Open(LocalTemplatePath)
		if err != nil {
			return "", fmt.Errorf("failed to open local template: %v", err)
		}
		defer localTemplate.Close()

		_, err = io.Copy(tempFile, localTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to copy local template: %v", err)
		}
		fmt.Printf("  %s✓ Local template loaded%s\n", ColorGreen, ColorReset)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to download template: HTTP %d", resp.StatusCode)
		}

		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to save template: %v", err)
		}
		fmt.Printf("  %s✓ Template downloaded%s\n", ColorGreen, ColorReset)
	}

	return tempFile.Name(), nil
}

// extractTemplate extracts the downloaded zip file
func (c *CLITool) extractTemplate(zipPath, projectName, language string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer reader.Close()

	templateDir := "TS"
	if language == "js" {
		templateDir = "JS"
	}

	for _, file := range reader.File {
		if !strings.HasPrefix(file.Name, templateDir+"/") && file.Name != templateDir {
			continue
		}

		fileName := strings.TrimPrefix(file.Name, templateDir+"/")
		if fileName == "" || strings.HasPrefix(fileName, "_sys/") || fileName == "_sys" {
			continue
		}

		filePath := filepath.Join(projectName, fileName)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		dir := filepath.Dir(filePath)
		os.MkdirAll(dir, os.ModePerm)

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

// customizePackageJson modifies the package.json file
func (c *CLITool) customizePackageJson(config ProjectConfig) {
	packagePath := filepath.Join(config.Name, "package.json")

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

	packageJson["name"] = strings.ToLower(strings.ReplaceAll(config.Name, " ", "-"))
	packageJson["description"] = config.Description
	packageJson["dependencies"] = make(map[string]interface{})
	packageJson["devDependencies"] = make(map[string]interface{})

	updatedData, _ := json.MarshalIndent(packageJson, "", "  ")
	ioutil.WriteFile(packagePath, updatedData, 0644)
}

// customizeEnvFile modifies the .env file
func (c *CLITool) customizeEnvFile(config ProjectConfig) {
	envPath := filepath.Join(config.Name, ".env")

	data, err := ioutil.ReadFile(envPath)
	if err != nil {
		log.Printf("Warning: Could not read .env file: %v", err)
		return
	}

	envContent := string(data)
	envContent = strings.ReplaceAll(envContent, "{{PORT}}", fmt.Sprintf("%d", config.Port))
	envContent = strings.ReplaceAll(envContent, "PORT=8080", fmt.Sprintf("PORT=%d", config.Port))
	ioutil.WriteFile(envPath, []byte(envContent), 0644)
}

// createConfigFile creates or updates the xypriss.config.json file with system variables
// If the file already exists, it merges the __sys__ section without touching other data
func (c *CLITool) createConfigFile(config ProjectConfig) {
	configPath := filepath.Join(config.Name, "xypriss.config.json")

	// System configuration to add/update
	sysConfig := map[string]interface{}{
		"__version__":     config.Version,
		"__author__":      config.Author,
		"__name__":        config.Name,
		"__description__": config.Description,
		"__alias__":       config.AppAlias,
		"__port__":        config.Port,
		"__PORT__":        config.Port,
	}

	// Try to read existing config file
	var existingConfig map[string]interface{}
	if data, err := ioutil.ReadFile(configPath); err == nil {
		// File exists, parse it
		if err := json.Unmarshal(data, &existingConfig); err != nil {
			log.Printf("Warning: Could not parse existing config, will create new: %v", err)
			existingConfig = make(map[string]interface{})
		}
	} else {
		// File doesn't exist, create new config
		existingConfig = make(map[string]interface{})
	}

	// Merge __sys__ section into existing config
	existingConfig["__sys__"] = sysConfig

	// Write merged config back to file
	data, err := json.MarshalIndent(existingConfig, "", "  ")
	if err != nil {
		log.Printf("Warning: Could not marshal config: %v", err)
		return
	}

	ioutil.WriteFile(configPath, data, 0644)
}

// customizeREADME modifies the README.md file
func (c *CLITool) customizeREADME(config ProjectConfig) {
	readmePath := filepath.Join(config.Name, "README.md")

	data, err := ioutil.ReadFile(readmePath)
	if err != nil {
		log.Printf("Warning: Could not read README.md: %v", err)
		return
	}

	readmeContent := string(data)
	readmeContent = strings.ReplaceAll(readmeContent, "{{PROJECT_NAME}}", config.Name)
	readmeContent = strings.ReplaceAll(readmeContent, "{{PROJECT_DESCRIPTION}}", config.Description)
	readmeContent = strings.ReplaceAll(readmeContent, "{{PORT}}", fmt.Sprintf("%d", config.Port))

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

	ioutil.WriteFile(readmePath, []byte(readmeContent), 0644)
}

// installDependencies installs project dependencies using Bun or npm
func (c *CLITool) installDependencies(projectName string, language string, mode string, strict bool) {
	configPath := filepath.Join(projectName, ".config")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("  %s✗ Failed to read .config file%s\n", ColorRed, ColorReset)
		return
	}

	// Parse dependencies from .config
	lines := strings.Split(string(data), "\n")
	var deps []string
	var devDeps []string
	inDeps := false
	inDevDeps := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Deps:") {
			inDeps = true
			inDevDeps = false
			continue
		}
		if strings.HasPrefix(line, "DevDeps:") {
			inDevDeps = true
			inDeps = false
			continue
		}
		if strings.HasPrefix(line, "- ") {
			dep := strings.TrimPrefix(line, "- ")
			if inDeps {
				deps = append(deps, dep)
			} else if inDevDeps {
				devDeps = append(devDeps, dep)
			}
		}
	}

	// Delete .config file immediately after storing data in memory
	os.Remove(configPath)

	// Determine installation mode
	useBun := false
	if mode == "b" {
		// Force bun mode
		if _, err := exec.LookPath("bun"); err == nil {
			useBun = true
			fmt.Printf("\n  %s⚡ Using 'BMode' (forced)%s\n", ColorCyan, ColorReset)
		} else {
			fmt.Printf("\n  %s✗ Bun not found, falling back to npm%s\n", ColorRed, ColorReset)
			if _, err := exec.LookPath("npm"); err != nil {
				fmt.Printf("  %s✗ npm is not installed%s\n", ColorRed, ColorReset)
				return
			}
		}
	} else if mode == "n" {
		// Force npm mode
		fmt.Printf("\n  %s→ Using npm (forced)%s\n", ColorCyan, ColorReset)
		if _, err := exec.LookPath("npm"); err != nil {
			fmt.Printf("  %s✗ npm is not installed%s\n", ColorRed, ColorReset)
			return
		}
	} else {
		// Auto-detect mode
		if _, err := exec.LookPath("bun"); err == nil {
			useBun = true
			fmt.Printf("\n  %s⚡ Using 'BMode' for faster installation%s\n", ColorCyan, ColorReset)
		} else {
			// Try to install Bun first
			fmt.Printf("\n  %s→ Bun not found, attempting to install...%s\n", ColorYellow, ColorReset)
			if c.installBun() {
				useBun = true
				fmt.Printf("  %s✓ Bun installed successfully%s\n", ColorGreen, ColorReset)
			} else {
				fmt.Printf("  %s→ Falling back to npm%s\n", ColorYellow, ColorReset)
				// Check npm availability
				if _, err := exec.LookPath("npm"); err != nil {
					fmt.Printf("  %s✗ npm is not installed%s\n", ColorRed, ColorReset)
					return
				}
			}
		}
	}

	totalDeps := len(deps) + len(devDeps)
	failedDeps := make([]string, 0)

	// Use parallelization for faster installation
	fmt.Printf("  %s⚡ Parallel installation enabled%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s│%s\n", ColorDim, ColorReset)

	// Install packages in parallel with intelligent batching
	type installResult struct {
		packageName string
		success     bool
		isDev       bool
		index       int
	}
	
	// Limit concurrent installations to avoid overwhelming the system
	maxConcurrent := 4
	if totalDeps < maxConcurrent {
		maxConcurrent = totalDeps
	}
	
	results := make(chan installResult, totalDeps)
	semaphore := make(chan struct{}, maxConcurrent)
	
	// Start all installations (regular dependencies)
	if len(deps) > 0 {
		fmt.Printf("%s├─ Dependencies (%d)%s\n", ColorDim, len(deps), ColorReset)
		for i, dep := range deps {
			go func(index int, packageName string) {
				semaphore <- struct{}{} // Acquire semaphore
				defer func() { <-semaphore }() // Release semaphore
				
				success := c.installPackageParallel(projectName, packageName, useBun, index+1, totalDeps)
				results <- installResult{packageName: packageName, success: success, isDev: false, index: index}
			}(i, dep)
		}
	}
	
	// Start all installations (dev dependencies)
	if len(devDeps) > 0 {
		if len(deps) > 0 {
			fmt.Printf("%s│%s\n", ColorDim, ColorReset)
		}
		fmt.Printf("%s└─ Dev Dependencies (%d)%s\n", ColorDim, len(devDeps), ColorReset)
		for i, dep := range devDeps {
			go func(index int, packageName string) {
				semaphore <- struct{}{} // Acquire semaphore
				defer func() { <-semaphore }() // Release semaphore
				
				success := c.installPackageParallel(projectName, packageName, useBun, len(deps)+index+1, totalDeps)
				results <- installResult{packageName: packageName, success: success, isDev: true, index: index}
			}(i, dep)
		}
	}
	
	// Collect results
	completed := 0
	for completed < totalDeps {
		result := <-results
		if !result.success {
			devLabel := ""
			if result.isDev {
				devLabel = " (dev)"
			}
			failedDeps = append(failedDeps, result.packageName+devLabel)
			
			// In strict mode, exit immediately on first error
			if strict {
				fmt.Printf("\n%s✗ Installation failed in strict mode%s\n", ColorRed, ColorReset)
				fmt.Printf("%s└─ Failed package: %s%s%s\n", ColorDim, ColorRed, result.packageName+devLabel, ColorReset)
				os.Exit(1)
			}
		}
		completed++
	}

	// Final summary
	fmt.Printf("\n")
	if len(failedDeps) > 0 {
		fmt.Printf("%s⚠ Installation completed with warnings%s\n", ColorYellow, ColorReset)
		fmt.Printf("%s├─ Failed: %d/%d packages%s\n", ColorDim, len(failedDeps), totalDeps, ColorReset)
		for i, dep := range failedDeps {
			prefix := "├─"
			if i == len(failedDeps)-1 {
				prefix = "└─"
			}
			fmt.Printf("%s%s ✗ %s%s\n", ColorDim, prefix, dep, ColorReset)
		}
	} else {
		fmt.Printf("%s✨ All dependencies installed successfully!%s\n", ColorGreen, ColorReset)
		fmt.Printf("%s└─ %d/%d packages%s\n", ColorDim, totalDeps, totalDeps, ColorReset)
	}
}

// installBun attempts to install Bun
func (c *CLITool) installBun() bool {
	cmd := exec.Command("npm", "install", "-g", "bun")
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	
	err := cmd.Run()
	return err == nil
}

// installSingleDependency installs a single package with inline progress
func (c *CLITool) installSingleDependency(projectName, dep string, isDev, useBun bool, current, total int, failedDeps *[]string, isLast, isDevSection bool) {
	// Prepare command
	var cmd *exec.Cmd
	
	// Special case: nquickdev needs npm for postinstall scripts (Bun ignores them)
	if useBun && dep == "nquickdev" {
		if isDev {
			cmd = exec.Command("npm", "install", "--save-dev", dep)
		} else {
			cmd = exec.Command("npm", "install", dep)
		}
	} else if useBun {
		if isDev {
			cmd = exec.Command("bun", "add", "-d", dep)
		} else {
			cmd = exec.Command("bun", "add", dep)
		}
	} else {
		if isDev {
			cmd = exec.Command("npm", "install", "--save-dev", dep)
		} else {
			cmd = exec.Command("npm", "install", dep)
		}
	}
	cmd.Dir = projectName

	// Tree branch characters
	branch := "├─"
	if isLast {
		branch = "└─"
	}

	// Progress indicator
	progress := fmt.Sprintf("[%d/%d]", current, total)

	// Show inline progress with spinner
	stop := c.showTreeSpinner(branch, dep, progress, isDev)

	var stdout, stderr bytes.Buffer
	if useBun {
		// Capture Bun output to filter it
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	} else {
		// Hide npm output
		cmd.Stdout = ioutil.Discard
		cmd.Stderr = &stderr
	}

	err := cmd.Run()
	
	// Stop spinner
	c.clearInlineSpinner(stop)

	// Display result
	devLabel := ""
	if isDev {
		devLabel = " %s(dev)%s"
		devLabel = fmt.Sprintf(devLabel, ColorDim, ColorReset)
	}

	if err != nil {
		fmt.Printf("\r\033[K%s   %s %s%s%s %s✗%s %s%s\n", 
			ColorDim, branch, progress, ColorReset, ColorRed, ColorRed, ColorReset, dep, devLabel)
		
		if !useBun {
			outputStr := stderr.String()
			if strings.Contains(outputStr, "404") || strings.Contains(outputStr, "Not Found") {
				subBranch := "├─"
				if isLast {
					subBranch = "   "
				}
				fmt.Printf("%s   %s %s→ Package not found%s\n", ColorDim, subBranch, ColorYellow, ColorReset)
			}
		}
		*failedDeps = append(*failedDeps, dep+devLabel)
	} else {
		fmt.Printf("\r\033[K%s   %s %s%s%s %s✓%s %s%s\n", 
			ColorDim, branch, progress, ColorReset, ColorGreen, ColorGreen, ColorReset, dep, devLabel)
		
		// Show Bun details on next line if available
		if useBun {
			c.displayBunDetails(&stdout, &stderr, isLast)
		}
	}
}

// showInlineSpinner shows a loading spinner on the current line
func (c *CLITool) showInlineSpinner(message string) chan bool {
	stop := make(chan bool)
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-stop:
				return
			case <-time.After(80 * time.Millisecond):
				fmt.Printf("\r  %s%s%s %s", ColorCyan, spinner[i%len(spinner)], ColorReset, message)
				i++
			}
		}
	}()
	return stop
}

// clearInlineSpinner clears the inline spinner line
func (c *CLITool) clearInlineSpinner(stop chan bool) {
	stop <- true
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("\r\033[K")
}

// displayProjectConfig displays project configuration in tree format
func (c *CLITool) displayProjectConfig(config ProjectConfig) {
	fmt.Printf("%s┌─ Project Configuration%s\n", ColorBold, ColorReset)
	fmt.Printf("%s│%s\n", ColorDim, ColorReset)
	fmt.Printf("%s├─%s %sName:%s %s\n", ColorDim, ColorReset, ColorCyan, ColorReset, config.Name)
	
	if config.Description != "" {
		fmt.Printf("%s├─%s %sDescription:%s %s\n", ColorDim, ColorReset, ColorCyan, ColorReset, config.Description)
	}
	
	langName := "TypeScript"
	if config.Language == "js" {
		langName = "JavaScript"
	}
	fmt.Printf("%s├─%s %sLanguage:%s %s\n", ColorDim, ColorReset, ColorCyan, ColorReset, langName)
	fmt.Printf("%s├─%s %sPort:%s %d\n", ColorDim, ColorReset, ColorCyan, ColorReset, config.Port)
	fmt.Printf("%s├─%s %sVersion:%s %s\n", ColorDim, ColorReset, ColorCyan, ColorReset, config.Version)
	fmt.Printf("%s├─%s %sApp Alias:%s %s\n", ColorDim, ColorReset, ColorCyan, ColorReset, config.AppAlias)
	fmt.Printf("%s├─%s %sAuthor:%s %s\n", ColorDim, ColorReset, ColorCyan, ColorReset, config.Author)
	
	// Features
	if config.WithAuth || config.WithUpload || config.WithMulti {
		fmt.Printf("%s└─%s %sFeatures:%s\n", ColorDim, ColorReset, ColorCyan, ColorReset)
		features := []string{}
		if config.WithAuth {
			features = append(features, "Authentication")
		}
		if config.WithUpload {
			features = append(features, "File Upload")
		}
		if config.WithMulti {
			features = append(features, "Multi-Server")
		}
		
		for i, feature := range features {
			if i == len(features)-1 {
				fmt.Printf("   %s└─%s %s\n", ColorDim, ColorReset, feature)
			} else {
				fmt.Printf("   %s├─%s %s\n", ColorDim, ColorReset, feature)
			}
		}
	} else {
		fmt.Printf("%s└─%s %sFeatures:%s None\n", ColorDim, ColorReset, ColorCyan, ColorReset)
	}
}

// displayFilteredBunOutput filters and displays Bun output
func (c *CLITool) displayFilteredBunOutput(stdout, stderr *bytes.Buffer, dep, progress, devLabel string, hasError bool) {
	output := stdout.String() + stderr.String()
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Skip "bun add vX.X.X (XXXXX)" lines - keep only the hash
		if strings.HasPrefix(line, "bun add v") {
			// Extract hash from format "bun add v1.3.3 (274e01c7)"
			if idx := strings.Index(line, "("); idx != -1 {
				if endIdx := strings.Index(line[idx:], ")"); endIdx != -1 {
					hash := line[idx+1 : idx+endIdx]
					fmt.Printf("\r\033[K  %s[%s]%s %s%s%s\n", ColorDim, hash, ColorReset, ColorDim, progress, ColorReset)
				}
			}
			continue
		}
		
		// Display other lines with proper formatting
		if strings.Contains(line, "installed") {
			// Format: "installed package@version"
			fmt.Printf("\r\033[K  %s%s%s %s✓ Installed%s %s%s\n", 
				ColorGreen, progress, ColorReset, ColorGreen, ColorReset, dep, devLabel)
		} else if strings.Contains(line, "packages installed") {
			// Show packages count
			fmt.Printf("  %s%s%s\n", ColorDim, line, ColorReset)
		} else if strings.Contains(line, "done") || strings.Contains(line, "ms]") {
			// Show timing info in dim color
			fmt.Printf("  %s%s%s\n", ColorDim, line, ColorReset)
		} else if !strings.HasPrefix(line, "bun add") && !strings.Contains(line, ".env") {
			// Show other relevant lines
			fmt.Printf("  %s%s%s\n", ColorDim, line, ColorReset)
		}
	}
	
	if hasError {
		fmt.Printf("\r  %s%s%s %s✗ Failed to install %s%s%s\n", 
			ColorRed, progress, ColorReset, ColorRed, dep, devLabel, ColorReset)
	}
}

// displayBunDetails displays filtered Bun output details
func (c *CLITool) displayBunDetails(stdout, stderr *bytes.Buffer, isLast bool) {
	output := stdout.String() + stderr.String()
	lines := strings.Split(output, "\n")
	
	indent := "       "
	if isLast {
		indent = "       "
	}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Extract hash from "bun add vX.X.X (XXXXX)"
		if strings.HasPrefix(line, "bun add v") {
			if idx := strings.Index(line, "("); idx != -1 {
				if endIdx := strings.Index(line[idx:], ")"); endIdx != -1 {
					hash := line[idx+1 : idx+endIdx]
					fmt.Printf("%s%s[%s]%s\n", indent, ColorDim, hash, ColorReset)
				}
			}
			continue
		}
		
		// Show package count and timing
		if strings.Contains(line, "packages installed") || strings.Contains(line, "done") {
			fmt.Printf("%s%s%s%s\n", indent, ColorDim, line, ColorReset)
		}
	}
}

// showTreeSpinner shows a spinner for tree-style installation
func (c *CLITool) showTreeSpinner(branch, dep string, progress string, isDev bool) chan bool {
	stop := make(chan bool)
	devLabel := ""
	if isDev {
		devLabel = fmt.Sprintf(" %s(dev)%s", ColorDim, ColorReset)
	}
	
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-stop:
				return
			case <-time.After(80 * time.Millisecond):
				fmt.Printf("\r%s   %s %s%s %s%s%s %s%s", 
					ColorDim, branch, progress, ColorReset,
					ColorCyan, spinner[i%len(spinner)], ColorReset,
					dep, devLabel)
				i++
			}
		}
	}()
	return stop
}