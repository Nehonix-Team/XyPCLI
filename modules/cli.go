package modules

import (
	"fmt"
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

// Run executes the CLI tool with the given command line arguments
func (c *CLITool) Run(args []string) {
	if len(args) < 1 {
		c.ShowHelp()
		return
	}

	command := args[0]

	switch command {
	case "init":
		c.InitProject()
	case "start":
		c.StartServer()
	case "version", "-v", "--version":
		fmt.Printf("XyPCLI v%s\n", c.version)
	case "help", "-h", "--help":
		c.ShowHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		c.ShowHelp()
	}
}