package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndk123-web/who-is-running/internal/utils"
)

// Tab indices
const (
	TabInspect = iota
	TabCommon
	TabActive
)

// Model represents the state of the TUI application
type Model struct {
	activeTab     int
	portInput     textinput.Model
	searchInput   textinput.Model
	activePorts   map[int]utils.PortInfo
	commonPorts   []int
	searchPorts   []int
	selectedIndex int
	statusMessage string
	statusType    string // "success", "error", "info"
	width         int
	height        int
	err           error
}

// Styles
var (
	accentColor   = lipgloss.Color("#7D56F4") // Deep Purple
	successColor  = lipgloss.Color("#10B981") // Teal / Green
	errorColor    = lipgloss.Color("#EF4444") // Rose / Red
	bgSecondary   = lipgloss.Color("#1E1E2E") // Dark Charcoal
	borderCol     = lipgloss.Color("#313244") // Border line
	mutedCol      = lipgloss.Color("#7F849C") // Muted gray text
	textCol       = lipgloss.Color("#CDD6F4") // White text

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Background(accentColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2).
			MarginBottom(1)

	tabStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(borderCol).
			Padding(0, 2)

	activeTabStyle = tabStyle.Copy().
			BorderForeground(accentColor).
			Foreground(accentColor).
			Bold(true)

	cardStyle = lipgloss.NewStyle().
			Background(bgSecondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderCol).
			Padding(1, 2).
			MarginBottom(1).
			Width(60)

	successBadge = lipgloss.NewStyle().
			Bold(true).
			Background(successColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	errorBadge = lipgloss.NewStyle().
			Bold(true).
			Background(errorColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	infoBadge = lipgloss.NewStyle().
			Bold(true).
			Background(accentColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(borderCol)

	selectedRowStyle = lipgloss.NewStyle().
				Background(accentColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	rowStyle = lipgloss.NewStyle().
			Foreground(textCol)

	footerStyle = lipgloss.NewStyle().
			Foreground(mutedCol).
			Italic(true).
			MarginTop(1)
)

// StartTUI enters alternate screen mode and runs the Bubble Tea TUI
func StartTUI() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// NewModel initializes the TUI state
func NewModel() Model {
	pi := textinput.New()
	pi.Placeholder = "e.g. 8080"
	pi.Focus()
	pi.CharLimit = 5
	pi.Width = 12

	si := textinput.New()
	si.Placeholder = "Type port or process name..."
	si.Width = 35

	m := Model{
		activeTab:   TabInspect,
		portInput:   pi,
		searchInput: si,
		commonPorts: []int{80, 443, 3000, 3306, 5000, 5432, 8000, 8080, 9000, 27017},
	}
	m.refresh()
	return m
}

// Init sets up input blinking
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// refresh re-scans system listening ports and updates filtered active list
func (m *Model) refresh() {
	ports, err := utils.ScanListeningPorts()
	if err != nil {
		m.err = err
		return
	}
	m.activePorts = ports
	m.updateSearch()
}

// updateSearch filters the active ports based on search input
func (m *Model) updateSearch() {
	query := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	var filtered []int
	for port, info := range m.activePorts {
		if query == "" {
			filtered = append(filtered, port)
			continue
		}
		portStr := strconv.Itoa(port)
		if strings.Contains(portStr, query) || strings.Contains(strings.ToLower(info.Process), query) {
			filtered = append(filtered, port)
		}
	}
	sort.Ints(filtered)
	m.searchPorts = filtered
}

// switchTab moves focus to the selected tab and adjusts textinput focus
func (m *Model) switchTab(tab int) {
	m.activeTab = tab
	m.selectedIndex = 0
	m.statusMessage = ""

	if m.activeTab == TabInspect {
		m.portInput.Focus()
		m.searchInput.Blur()
	} else if m.activeTab == TabActive {
		m.portInput.Blur()
		m.searchInput.Focus()
	} else {
		m.portInput.Blur()
		m.searchInput.Blur()
	}
}

// Update handles terminal keyboard and window events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Reset status message on interaction
		m.statusMessage = ""

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.activeTab == TabActive {
				break
			}
			return m, tea.Quit

		case "tab":
			m.switchTab((m.activeTab + 1) % 3)
			return m, nil

		case "shift+tab":
			m.switchTab((m.activeTab + 2) % 3)
			return m, nil

		case "a", "A":
			if m.activeTab == TabActive {
				break
			}
			m.switchTab(TabInspect)
			return m, nil

		case "b", "B":
			if m.activeTab == TabActive {
				break
			}
			m.switchTab(TabCommon)
			return m, nil

		case "c", "C":
			if m.activeTab == TabActive {
				break
			}
			m.switchTab(TabActive)
			return m, nil

		case "ctrl+r":
			m.refresh()
			m.statusMessage = "🔄 Port list updated successfully!"
			m.statusType = "info"
			return m, nil

		case "r", "R":
			if m.activeTab == TabActive {
				break
			}
			m.refresh()
			m.statusMessage = "🔄 Port list updated successfully!"
			m.statusType = "info"
			return m, nil

		case "up":
			if m.activeTab == TabCommon {
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			} else if m.activeTab == TabActive {
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			}
			return m, nil

		case "down":
			if m.activeTab == TabCommon {
				if m.selectedIndex < len(m.commonPorts)-1 {
					m.selectedIndex++
				}
			} else if m.activeTab == TabActive {
				if m.selectedIndex < len(m.searchPorts)-1 {
					m.selectedIndex++
				}
			}
			return m, nil

		case "ctrl+k", "k", "K":
			keyStr := msg.String()
			if (keyStr == "k" || keyStr == "K") && m.activeTab == TabActive {
				break
			}

			var pidToKill int
			var portKilled int
			var procName string

			if m.activeTab == TabInspect {
				val := m.portInput.Value()
				port, err := strconv.Atoi(val)
				if err == nil {
					if info, exists := m.activePorts[port]; exists {
						pidToKill = info.PID
						portKilled = port
						procName = info.Process
					}
				}
			} else if m.activeTab == TabCommon {
				if m.selectedIndex >= 0 && m.selectedIndex < len(m.commonPorts) {
					port := m.commonPorts[m.selectedIndex]
					if info, exists := m.activePorts[port]; exists {
						pidToKill = info.PID
						portKilled = port
						procName = info.Process
					}
				}
			} else if m.activeTab == TabActive {
				if m.selectedIndex >= 0 && m.selectedIndex < len(m.searchPorts) {
					port := m.searchPorts[m.selectedIndex]
					if info, exists := m.activePorts[port]; exists {
						pidToKill = info.PID
						portKilled = port
						procName = info.Process
					}
				}
			}

			if pidToKill > 0 {
				err := utils.KillProcess(pidToKill)
				if err != nil {
					m.statusMessage = fmt.Sprintf("❌ Failed to kill process %s (PID %d): %v", procName, pidToKill, err)
					m.statusType = "error"
				} else {
					m.statusMessage = fmt.Sprintf("✅ Freed port %d by killing %s (PID %d)!", portKilled, procName, pidToKill)
					m.statusType = "success"
					if m.activeTab == TabInspect {
						m.portInput.SetValue("")
					}
					m.refresh()
				}
			} else {
				m.statusMessage = "ℹ️ No running process selected to terminate."
				m.statusType = "info"
			}
			return m, nil
		}
	}

	// Update text inputs based on active tab
	if m.activeTab == TabInspect {
		m.portInput, cmd = m.portInput.Update(msg)
	} else if m.activeTab == TabActive {
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.updateSearch()
		// Auto-clamp selection index if search updates
		if m.selectedIndex >= len(m.searchPorts) {
			m.selectedIndex = len(m.searchPorts) - 1
		}
		if m.selectedIndex < 0 {
			m.selectedIndex = 0
		}
	}

	return m, cmd
}

// View builds the final TUI render string
func (m Model) View() string {
	var body string

	if m.err != nil {
		return fmt.Sprintf("Error running scanner: %v\nPress Q to quit.", m.err)
	}

	switch m.activeTab {
	case TabInspect:
		body = m.viewInspectTab()
	case TabCommon:
		body = m.viewCommonPortsTab()
	case TabActive:
		body = m.viewActivePortsTab()
	}

	// Status line if any message exists
	statusLine := m.renderStatus()

	// Title
	header := titleStyle.Render("⚡ WHO IS RUNNING? ⚡")

	// Tabs
	tabs := m.renderTabs()

	// Footer instructions
	var footer string
	if m.activeTab == TabInspect {
		footer = footerStyle.Render("a-c / Tab: Switch Tab • R: Refresh • K: Kill Process • Q: Quit")
	} else {
		footer = footerStyle.Render("a-c / Tab: Switch Tab • ↑/↓: Scroll List • R: Refresh • K: Kill Process • Q: Quit")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabs,
		"",
		body,
		statusLine,
		footer,
	)
}

func (m Model) renderTabs() string {
	tabs := []string{"[a] Inspect Port", "[b] Common Ports", "[c] Active Listening"}
	var renderedTabs []string
	for i, t := range tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(t))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m Model) renderStatus() string {
	if m.statusMessage == "" {
		return ""
	}

	var style lipgloss.Style
	switch m.statusType {
	case "success":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(successColor).Padding(0, 1).Bold(true)
	case "error":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(errorColor).Padding(0, 1).Bold(true)
	default:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(accentColor).Padding(0, 1).Bold(true)
	}

	return fmt.Sprintf("\n%s\n", style.Render(m.statusMessage))
}

func (m Model) viewInspectTab() string {
	var body string
	val := m.portInput.Value()

	if val == "" {
		body = cardStyle.Render("Type a port number above to inspect its status.")
	} else {
		port, err := strconv.Atoi(val)
		if err != nil {
			body = cardStyle.Render(fmt.Sprintf("⚠️ Invalid port number: %s", val))
		} else {
			if info, exists := m.activePorts[port]; exists {
				status := errorBadge.Render(" IN USE ")
				processInfo := fmt.Sprintf("🔥 Process:  %s\n🆔 PID:      %d\n🌐 Protocol: %s",
					lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render(info.Process),
					info.PID,
					info.Protocol)
				action := lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render("\n💀 Press [K] to kill this process and free the port.")

				body = cardStyle.Render(fmt.Sprintf("%s Port %d is currently blocked!\n\n%s\n%s", status, port, processInfo, action))
			} else {
				status := successBadge.Render(" FREE ")
				body = cardStyle.Render(fmt.Sprintf("%s Port %d is free and available!\n\n🚀 You can bind any service to this port.", status, port))
			}
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Enter Port to Inspect:",
		m.portInput.View(),
		"",
		body,
	)
}

func (m Model) viewCommonPortsTab() string {
	var rows []string
	rows = append(rows, tableHeaderStyle.Render(fmt.Sprintf("%-6s %-12s %-25s %-8s", "Port", "Status", "Process", "PID")))

	for i, port := range m.commonPorts {
		info, exists := m.activePorts[port]

		var status, procName, pidStr string
		if exists {
			status = errorBadge.Render(" IN USE ")
			procName = info.Process
			pidStr = strconv.Itoa(info.PID)
		} else {
			status = successBadge.Render(" FREE   ")
			procName = "-"
			pidStr = "-"
		}

		rowContent := fmt.Sprintf("%-6d %-12s %-25s %-8s", port, status, procName, pidStr)

		if i == m.selectedIndex {
			rows = append(rows, selectedRowStyle.Render(rowContent))
		} else {
			rows = append(rows, rowStyle.Render(rowContent))
		}
	}

	list := lipgloss.JoinVertical(lipgloss.Left, rows...)

	var help string
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.commonPorts) {
		selectedPort := m.commonPorts[m.selectedIndex]
		if _, exists := m.activePorts[selectedPort]; exists {
			help = lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render("\n💀 Press [K] to kill process on selected port.")
		}
	}

	return fmt.Sprintf("Common Development Ports:\n\n%s%s", list, help)
}

func (m Model) viewActivePortsTab() string {
	searchBar := fmt.Sprintf("Search (Port or Process): %s", m.searchInput.View())

	var rows []string
	rows = append(rows, tableHeaderStyle.Render(fmt.Sprintf("%-6s %-25s %-8s %-10s", "Port", "Process", "PID", "Protocol")))

	if len(m.searchPorts) == 0 {
		rows = append(rows, rowStyle.Render("No active ports found matching search."))
	} else {
		for i, port := range m.searchPorts {
			info := m.activePorts[port]
			rowContent := fmt.Sprintf("%-6d %-25s %-8d %-10s", port, info.Process, info.PID, info.Protocol)

			if i == m.selectedIndex {
				rows = append(rows, selectedRowStyle.Render(rowContent))
			} else {
				rows = append(rows, rowStyle.Render(rowContent))
			}
		}
	}

	list := lipgloss.JoinVertical(lipgloss.Left, rows...)

	var help string
	if len(m.searchPorts) > 0 && m.selectedIndex >= 0 && m.selectedIndex < len(m.searchPorts) {
		selectedPort := m.searchPorts[m.selectedIndex]
		if _, exists := m.activePorts[selectedPort]; exists {
			help = lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render("\n💀 Press [K] to kill process on selected port.")
		}
	}

	return fmt.Sprintf("%s\n\nActive Ports in Use:\n\n%s%s", searchBar, list, help)
}
