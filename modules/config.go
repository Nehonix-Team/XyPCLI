package modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
) 
 

// ProjectConfig holds the configuration for a new XyPriss project
// This struct contains all the necessary information to generate a complete
// XyPriss application with the selected features
type ProjectConfig struct {
	Name        string // Project name (used for directory and package.json)
	Description string // Project description
	Version     string // Initial version (defaults to "1.0.0")
	Port        int    // Server port (defaults to 3000)
	Language    string // Programming language: "js" or "ts" (defaults to "ts")
	WithAuth    bool   // Include JWT authentication system
	WithUpload  bool   // Include file upload functionality with multer
	WithMulti   bool   // Include multi-server configuration
}

// GetProjectConfig interactively collects basic project configuration from the user
// This function prompts the user for:
// - Project name (used for directory and package.json)
// - Project description
// - Programming language (JavaScript or TypeScript)
//
// Returns a ProjectConfig struct with default features enabled for simplicity
func GetProjectConfig() ProjectConfig {
	reader := bufio.NewReader(os.Stdin)

	config := ProjectConfig{
		Port:       3000,
		Version:    "1.0.0",
		Language:   "ts",  // Default to TypeScript
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

	// Programming language selection
	fmt.Printf("%sProgramming language (js/ts):%s ", ColorCyan, ColorReset)
	lang, _ := reader.ReadString('\n')
	config.Language = strings.TrimSpace(strings.ToLower(lang))
	if config.Language != "js" && config.Language != "ts" {
		config.Language = "ts" // Default to TypeScript
	}

	return config
}