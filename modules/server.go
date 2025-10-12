package modules

import (
	"fmt"
	"os"
	"os/exec"
) 

// StartServer starts the XyPriss development server in the current directory
func (c *CLITool) StartServer() {
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