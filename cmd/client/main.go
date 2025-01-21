package main

import (
	"github.com/korol8484/gophkeeper/internal/client/cli"
	"log"
)

func main() {
	rootCmd := cli.Root()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
