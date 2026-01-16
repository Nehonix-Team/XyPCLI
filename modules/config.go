package modules

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// ProjectConfig holds the configuration for a new XyPriss project
// This struct contains all the necessary information to generate a complete
// XyPriss application with the selected features
type ProjectConfig struct {
	Name         string // Project name (used for directory and package.json)
	Description  string // Project description
	Version      string // Initial version (defaults to "1.0.0")
	Port         int    // Server port (defaults to 3000)
	Language     string // Programming language: "js" or "ts" (defaults to "ts")
	AppName      string // Application name (defaults to "XyPriss")
	AppAlias     string // Application alias (defaults to "XyP")
	Author       string // Author name (defaults to "Nehonix")
	WithAuth     bool   // Include JWT authentication system
	WithUpload   bool   // Include file upload functionality with multer
	WithMulti    bool   // Include multi-server configuration
}

// handleExistingDirectory checks if a directory exists and handles the case where it's not empty
// Returns true if we can proceed with the directory name, false if user wants to choose another name
func handleExistingDirectory(dirName string, reader *bufio.Reader) bool {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// Directory doesn't exist, we can proceed
		return true
	}

	// Directory exists, check if it's empty
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Printf("%sError reading directory: %v%s\n", ColorRed, err, ColorReset)
		return false
	}

	if len(files) == 0 {
		// Directory exists but is empty, we can proceed
		return true
	}

	// Directory exists and is not empty, ask user what to do
	fmt.Printf("\n%s⚠ Directory '%s' already exists and is not empty.%s\n", ColorYellow, dirName, ColorReset)
	fmt.Printf("%sWhat would you like to do?%s\n", ColorBold, ColorReset)
	fmt.Printf("  %s1.%s Delete the directory and create a new project\n", ColorCyan, ColorReset)
	fmt.Printf("  %s2.%s Choose a different project name\n", ColorCyan, ColorReset)
	fmt.Printf("%sEnter your choice (1 or 2):%s ", ColorBold, ColorReset)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Printf("%s🗑️  Deleting existing directory '%s'...%s\n", ColorRed, dirName, ColorReset)
		err := os.RemoveAll(dirName)
		if err != nil {
			fmt.Printf("%s❌ Failed to delete directory: %v%s\n", ColorRed, err, ColorReset)
			return false
		}
		fmt.Printf("%s✅ Directory deleted successfully%s\n", ColorGreen, ColorReset)
		return true
	case "2":
		return false
	default:
		fmt.Printf("%s❌ Invalid choice. Please choose 1 or 2.%s\n", ColorRed, ColorReset)
		return handleExistingDirectory(dirName, reader) // Recursive call to ask again
	}
}

// GetProjectConfig interactively collects basic project configuration from the user
// This function prompts the user for:
// - Project name (used for directory and package.json)
// - Project description
// - Programming language (JavaScript or TypeScript)
//
// If InitFlags are provided, they will be used instead of prompting the user
// Returns a ProjectConfig struct with default features enabled for simplicity
func GetProjectConfig(flags InitFlags) ProjectConfig {
	reader := bufio.NewReader(os.Stdin)

	config := ProjectConfig{
		Port:       3000,
		Version:    "1.0.0",
		Language:   "ts",    // Default to TypeScript
		AppAlias:   "XyP",
		Author:     "Nehonix-Team",
		WithAuth:   true,    // Enable by default for better DX
		WithUpload: true,    // Enable by default for better DX
		WithMulti:  false,   // Keep simple by default
	}

	// Project name - used for directory name and package.json
	if flags.Name != "" {
		config.Name = flags.Name
	} else {
		fmt.Printf("%sProject name:%s ", ColorCyan, ColorReset)
		name, _ := reader.ReadString('\n')
		config.Name = strings.TrimSpace(name)
		if config.Name == "" {
			config.Name = "my-xypriss-app"
		}
	}

	// Check if directory exists and handle it
	if !handleExistingDirectory(config.Name, reader) {
		// User chose to use a different name
		return GetProjectConfig(InitFlags{}) // Reset flags for new interactive session
	}

	// Project description - used in package.json and README
	if flags.Description != "" {
		config.Description = flags.Description
	} else {
		fmt.Printf("%sDescription:%s ", ColorCyan, ColorReset)
		desc, _ := reader.ReadString('\n')
		config.Description = strings.TrimSpace(desc)
		if config.Description == "" {
			config.Description = "A XyPriss application"
		}
	}

	// Programming language selection
	if flags.Language != "" {
		config.Language = strings.ToLower(flags.Language)
		if config.Language != "js" && config.Language != "ts" {
			config.Language = "ts" // Default to TypeScript
		}
	} else {
		fmt.Printf("%sProgramming language (js/ts):%s ", ColorCyan, ColorReset)
		lang, _ := reader.ReadString('\n')
		config.Language = strings.TrimSpace(strings.ToLower(lang))
		if config.Language != "js" && config.Language != "ts" {
			config.Language = "ts" // Default to TypeScript
		}
	}

	// Server port selection
	if flags.Port != "" {
		if port, err := strconv.Atoi(flags.Port); err == nil && port > 0 && port < 65536 {
			config.Port = port
		} else {
			fmt.Printf("%sInvalid port format, using default 3000%s\n", ColorYellow, ColorReset)
		}
	} else {
		fmt.Printf("%sServer port:%s ", ColorCyan, ColorReset)
		portStr, _ := reader.ReadString('\n')
		portStr = strings.TrimSpace(portStr)
		if portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
				config.Port = port
			} else {
				fmt.Printf("%sInvalid port format, using default 3000%s\n", ColorYellow, ColorReset)
			}
		}
	}

	// Application version
	if flags.Version != "" {
		config.Version = flags.Version
	} else {
		fmt.Printf("%sApplication version:%s ", ColorCyan, ColorReset)
		version, _ := reader.ReadString('\n')
		config.Version = strings.TrimSpace(version)
		if config.Version == "" {
			config.Version = "1.0.0"
		}
	}

	// Application alias
	if flags.Alias != "" {
		config.AppAlias = flags.Alias
	} else {
		fmt.Printf("%sApplication alias:%s ", ColorCyan, ColorReset)
		appAlias, _ := reader.ReadString('\n')
		config.AppAlias = strings.TrimSpace(appAlias)
		if config.AppAlias == "" {
			config.AppAlias = "XyP"
		}
	}

	// Author name
	if flags.Author != "" {
		config.Author = flags.Author
	} else {
		fmt.Printf("%sAuthor name:%s ", ColorCyan, ColorReset)
		author, _ := reader.ReadString('\n')
		config.Author = strings.TrimSpace(author)
		if config.Author == "" {
			config.Author = "Nehonix"
		}
	}

	return config
}