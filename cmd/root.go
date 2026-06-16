package cmd

import (
	"fmt"
	"os"

	"github.com/ndk123-web/who-is-running/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "who-is-running",
	Short: "Find which ports are running",
	Long:  "A CLI tool to inspect running ports on your machine.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.StartTUI(); err != nil {
			fmt.Printf("Error starting TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}