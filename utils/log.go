package utils

import (
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func InitLogger(debug bool) {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})

	styles := log.DefaultStyles()
	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#ff3030")).
		Foreground(lipgloss.Color("0"))

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#30ff30")).
		Foreground(lipgloss.Color("0"))

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#ffff00")).
		Foreground(lipgloss.Color("0"))

	styles.Keys["message"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#3030ff"))
	styles.Values["message"] = lipgloss.NewStyle().Bold(true)

	if debug {
		Logger.SetLevel(log.DebugLevel)
	}
}
