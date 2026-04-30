package utils

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
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
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff0000")).
			Foreground(lipgloss.Color("#ff3030")).
			Padding(1, 2).
			Bold(true).
			Render("❌ Error\n\n"+content+"\n\nPress q to quit."),
	)
}

func ListView(table table.Model, dir string, width, height int) string {
	infoBox := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#d70070")).
		Padding(1, 2).
		Bold(true).
		Render(lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(width - 10).
			Render(fmt.Sprintf("Viewing directory: '%s'", dir)),
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
			table.HelpView()),
	)
}
