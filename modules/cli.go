package modules

import (
	"fmt"
	"strings"
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

// CLITool represents the XyPriss CLI tool with version information
// This tool provides commands for initializing new projects and managing
// XyPriss applications
type CLITool struct {
	version string // CLI version
}
 
// NewCLITool creates a new CLI tool instance
func NewCLITool(version string) *CLITool {
	return &CLITool{version: version}
}

// ShowHelp displays the CLI help information with beautiful branding
// This function provides comprehensive usage instructions including:
// - XyPriss ASCII art logo
// - Available commands and their descriptions
// - Usage syntax with colored output
// - Example command invocations
// - Version information
func (c *CLITool) ShowHelp() {
	fmt.Println(XyPrissLogo)
	fmt.Printf("%sCLI Tool v%s%s\n\n", ColorYellow, c.version, ColorReset)
	fmt.Printf("%sUSAGE:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %sxypcli <command> [options]%s\n", ColorCyan, ColorReset)
	fmt.Println()
	fmt.Printf("%sCOMMANDS:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %sinit%s     Initialize a new XyPriss project with all necessary configuration\n", ColorGreen, ColorReset)
	fmt.Printf("  %sstart%s    Start the XyPriss development server in the current directory\n", ColorGreen, ColorReset)
	fmt.Printf("  %sinstall%s  Install one or more packages using the XyPriss installation system\n", ColorGreen, ColorReset)
	fmt.Printf("  %sversion%s  Show CLI version information\n", ColorGreen, ColorReset)
	fmt.Printf("  %shelp%s     Show this help message\n", ColorGreen, ColorReset)
	fmt.Println()
	fmt.Printf("%sINIT OPTIONS:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %s--name <name>%s         Project name (default: interactive prompt)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--desc <description>%s  Project description\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--lang <js|ts>%s        Programming language (default: ts)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--port <port>%s         Server port (default: 3000)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--version <version>%s   Application version (default: 1.0.0)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--alias <alias>%s       Application alias (default: XyP)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--author <author>%s     Author name (default: Nehonix-Team)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--mode <b|n>%s          Installation mode: 'b' for bun, 'n' for npm (default: auto)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--strict%s              Exit immediately if any package installation fails\n", ColorCyan, ColorReset)
	fmt.Println()
	fmt.Printf("%sINSTALL OPTIONS:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %s--mode <b|n>%s          Installation mode: 'b' for bun, 'n' for npm (default: auto)\n", ColorCyan, ColorReset)
	fmt.Println()
	fmt.Printf("%sEXAMPLES:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %sxypcli init%s                                    # Interactive mode\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli init --name my-app --port 8080%s         # Quick init with options\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli init --name my-app --mode n%s            # Force npm installation\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli start%s                                   # Start development server\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli install xypriss cors%s                    # Install multiple packages\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli install xypriss --mode b%s                # Install with bun\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli --version%s                               # Show CLI version\n", ColorMagenta, ColorReset)
	fmt.Printf("  %sxypcli help%s                                    # Show this help\n", ColorMagenta, ColorReset)
	fmt.Println()
	fmt.Printf("%sFor more information, visit: %shttps://github.com/Nehonix-Team/XyPriss%s\n", ColorDim, ColorBlue, ColorReset)
}

// Run executes the CLI tool with the given command line arguments
func (c *CLITool) Run(args []string) {
	if len(args) < 1 {
		c.ShowHelp()
		return
	}

	command := args[0]

	switch command {
	case "init":
		// Parse init flags
		initFlags := parseInitFlags(args[1:])
		c.InitProject(initFlags)
	case "start":
		c.StartServer()
	case "install":
		if len(args) < 2 {
			fmt.Printf("%s❌ Package name required%s\n", ColorRed, ColorReset)
			fmt.Printf("%sUsage:%s xypcli install <package-name> [package-name...] [--mode <b|n>]\n", ColorBold, ColorReset)
			return
		}
		// Parse install flags and packages
		packages, mode := parseInstallArgs(args[1:])
		if len(packages) == 0 {
			fmt.Printf("%s❌ At least one package name required%s\n", ColorRed, ColorReset)
			return
		}
		c.InstallPackages(packages, mode)
	case "version", "-v", "--version":
		fmt.Printf("XyPCLI v%s\n", c.version)
	case "help", "-h", "--help":
		c.ShowHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		c.ShowHelp()
	}
}

// InitFlags holds command-line flags for the init command
type InitFlags struct {
	Name        string
	Description string
	Language    string
	Port        string
	Version     string
	Alias       string
	Author      string
	Mode        string
	Strict      bool   // Exit on first installation error
}

// parseInitFlags parses command-line flags for the init command
func parseInitFlags(args []string) InitFlags {
	flags := InitFlags{}
	
	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "--") {
			continue
		}
		
		flag := args[i]
		var value string
		
		// Get the value (next argument or after =)
		if strings.Contains(flag, "=") {
			parts := strings.SplitN(flag, "=", 2)
			flag = parts[0]
			value = parts[1]
		} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			value = args[i+1]
			i++ // Skip next arg since we used it as value
		}
		
		switch flag {
		case "--name":
			flags.Name = value
		case "--desc", "--description":
			flags.Description = value
		case "--lang", "--language":
			flags.Language = value
		case "--port":
			flags.Port = value
		case "--version":
			flags.Version = value
		case "--alias":
			flags.Alias = value
		case "--author":
			flags.Author = value
		case "--mode":
			flags.Mode = value
		case "--strict":
			flags.Strict = true
		}
	}
	
	return flags
}

// parseInstallArgs parses arguments for the install command
// Returns the list of packages and the installation mode
func parseInstallArgs(args []string) ([]string, string) {
	packages := []string{}
	mode := ""
	
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			if args[i] == "--mode" || strings.HasPrefix(args[i], "--mode=") {
				if strings.Contains(args[i], "=") {
					mode = strings.SplitN(args[i], "=", 2)[1]
				} else if i+1 < len(args) {
					mode = args[i+1]
					i++
				}
			}
		} else {
			packages = append(packages, args[i])
		}
	}
	
	return packages, mode
}