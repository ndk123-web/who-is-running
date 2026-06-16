package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "who-is-running",
	Short: "Find which ports are running",
	Long:  "A CLI tool to inspect running ports on your machine.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}