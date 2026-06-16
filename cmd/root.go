package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/ndk123-web/who-is-running/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "who-is-running",
	Short:         "Find which ports are running",
	Long:          "A CLI tool to inspect running ports on your machine.",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.StartTUI(); err != nil {
			fmt.Printf("Error starting TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

// Styling for CLI errors and help messages
var (
	helpAccent = lipgloss.Color("#7D56F4")
	helpBorder = lipgloss.Color("#313244")
	helpBg     = lipgloss.Color("#1E1E2E")
	helpErr    = lipgloss.Color("#EF4444")

	helpCard = lipgloss.NewStyle().
			Background(helpBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(helpBorder).
			Padding(1, 2).
			Width(60)

	errCard = lipgloss.NewStyle().
			Background(helpBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(helpErr).
			Padding(1, 2).
			Width(60)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(helpAccent)

	errTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(helpErr)
)

func printHelp() {
	title := titleStyle.Render("WHO IS RUNNING — Help & Usage")
	content := fmt.Sprintf(
		"%s\n\n"+
			"%s\n"+
			"    who-is-running\n\n"+
			"%s\n"+
			"    who-is-running at [port]\n\n"+
			"%s\n"+
			"    -h, --help   Show this help screen",
		title,
		lipgloss.NewStyle().Bold(true).Render("Interactive Dashboard (TUI):"),
		lipgloss.NewStyle().Bold(true).Render("Inspect Single Port (CLI Card):"),
		lipgloss.NewStyle().Bold(true).Render("Global Flags:"),
	)
	fmt.Println(helpCard.Render(content))
}

func printError(err error) {
	title := errTitleStyle.Render("Command Error")
	content := fmt.Sprintf(
		"%s\n\n"+
			"Error: %s\n\n"+
			"Need help? Run %s for usage or run without arguments to start the interactive dashboard.",
		title,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#CDD6F4")).Render(err.Error()),
		lipgloss.NewStyle().Bold(true).Foreground(helpAccent).Render("who-is-running --help"),
	)
	fmt.Println(errCard.Render(content))
}

func Execute() {
	// Set custom help function
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printHelp()
	})

	err := rootCmd.Execute()
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}