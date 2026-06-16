package cmd

import (
	"fmt"
	"strconv"

	"github.com/ndk123-web/who-is-running/internal/utils"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "at [port]",
	Short: "Scan a single port",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid port")
			return
		}

		utils.ScanSinglePort(port)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}