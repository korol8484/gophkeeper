package main

import (
	"github.com/fatih/color"
	"github.com/korol8484/gophkeeper/internal/server/cli"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	cmd := cli.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		_, _ = color.New(color.FgHiRed, color.Bold).Fprintln(os.Stderr, err.Error())

		return 1
	}

	return 0
}
