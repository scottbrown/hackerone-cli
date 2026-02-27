package main

import (
	"fmt"
	"os"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/scottbrown/hackerone-cli/internal/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd(hackeronecli.Version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
