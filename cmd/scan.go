package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/ndk123-web/who-is-running/internal/utils"
	"github.com/spf13/cobra"
)

// Styles for the CLI card
var (
	cliAccent  = lipgloss.Color("#7D56F4")
	cliSuccess = lipgloss.Color("#10B981")
	cliError   = lipgloss.Color("#EF4444")
	cliBorder  = lipgloss.Color("#313244")
	cliBg      = lipgloss.Color("#1E1E2E")

	cliCard = lipgloss.NewStyle().
		Background(cliBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cliBorder).
		Padding(1, 2).
		Width(50)

	cliSuccessBadge = lipgloss.NewStyle().
			Bold(true).
			Background(cliSuccess).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	cliErrorBadge = lipgloss.NewStyle().
			Bold(true).
			Background(cliError).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)
)

var scanCmd = &cobra.Command{
	Use:   "at [port]",
	Short: "Scan a single port",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid port. Please specify a numeric port (e.g., 8080)")
			return
		}

		ports, err := utils.ScanListeningPorts()
		if err != nil {
			// Fallback if system commands fail
			utils.ScanSinglePort(port)
			return
		}

		var cardContent string
		if info, exists := ports[port]; exists {
			// In Use
			status := cliErrorBadge.Render(" IN USE ")
			procName := lipgloss.NewStyle().Foreground(cliAccent).Bold(true).Render(info.Process)

			cardContent = fmt.Sprintf(
				"%s Port %d is currently blocked!\n\n🔥 Process:  %s\n🆔 PID:      %d\n🌐 Protocol: %s",
				status, port, procName, info.PID, info.Protocol,
			)
		} else {
			// Free
			status := cliSuccessBadge.Render(" FREE ")
			cardContent = fmt.Sprintf(
				"%s Port %d is free and available!\n\n🚀 You can bind any service to this port.",
				status, port,
			)
		}

		fmt.Println(cliCard.Render(cardContent))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}