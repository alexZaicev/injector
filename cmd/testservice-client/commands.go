package main

import (
	"github.com/hashicorp/cli"
	"os"
)

func newCLI() cli.Ui {
	return &cli.BasicUi{
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
		Reader:      os.Stdin,
	}
}

func setupCommands(ui cli.Ui) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"echo": func() (cli.Command, error) {
			return &EchoCommand{
				ui: ui,
			}, nil
		},
	}
}
