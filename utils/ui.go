package utils

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	Destructive  = lipgloss.Color("#ff3030")
	Constructive = lipgloss.Color("#5aff5a")
)

var (
	HelpStyle  = lipgloss.NewStyle().Faint(true).MarginTop(1)
	ErrorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff0000")).
			Foreground(Destructive).
			Padding(1, 2).
			Bold(true)
	SuccessStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00ff00")).
			Foreground(Constructive).
			Padding(1, 2).
			Bold(true)
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d70070")).
			Padding(1, 2).
			Bold(true)
	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5020")).
			Bold(true)
)

func ErrorView(err error, width, height int) string {
	content := lipgloss.NewStyle().
		Width(width - 10).
		Render(err.Error())

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		ErrorStyle.Render("❌ Error\n\n"+content+"\n\nPress q to quit."),
	)
}

func SuccessView(msg string, width, height int) string {
	content := lipgloss.NewStyle().
		Width(width - 10).
		Render(msg)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		SuccessStyle.Render("✅ "+content+"\n\nPress q to quit."),
	)
}

func LoadingView(message string, width, height int) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		message,
	)
}

func ListView(table table.Model, dir string, width, height int) string {
	infoBox := HeaderStyle.
		Render(lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(width - 10).
			Render(fmt.Sprintf("📁 Directory: '%s'", dir)),
		)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			infoBox,
			table.View(),
			HelpStyle.Render("↑↓/j,k: navigate • Enter: select"),
		),
	)
}
