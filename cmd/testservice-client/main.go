package main

import (
	"fmt"
	"github.com/hashicorp/cli"
	"injector/internal/drivers/codes"
	"injector/internal/drivers/logging"
	"os"
)

func run() int {
	ui := newCLI()
	if err := setupCommands(ui); err != nil {
		return codes.Failure
	}

	runner := &cli.CLI{
		Name:     "testservice-client",
		Args:     os.Args[1:],
		Commands: setupCommands(ui),
	}

	exit, err := runner.Run()
	if err != nil {
		ui.Error(fmt.Sprintf("error executing command: %s", err))
		return codes.Failure
	}

	return exit
}

func main() {
	logging.SetupConsoleLogger()
	os.Exit(run())
}
