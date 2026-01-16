package main

import (
	"os"
	"github.com/Nehonix-Team/XyPCLI/modules"
) 

func main() {
	cli := modules.NewCLITool("1.0.2")

	args := os.Args[1:] // Skip program name
	cli.Run(args)
}
