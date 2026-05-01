package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/internal/oauth/providers"
	"github.com/udaycmd/rdv/utils"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#66B3FF"))
	subtitleStyle = lipgloss.NewStyle().Faint(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#66B3FF")).Bold(true)
	statusStyle   = lipgloss.NewStyle().Padding(0, 1)
	revokedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#69DB7C")).Bold(true)
	defaultStyle  = lipgloss.NewStyle().Faint(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#69DB7C"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	helpStyle     = lipgloss.NewStyle().Faint(true).MarginTop(1)
)

var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Interactive TUI to manage remote drives",
	Long:  "Launch an interactive terminal UI to add, remove, select, and manage your remote drive connections.",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(
			newDriveModel(config),
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
		if _, err := p.Run(); err != nil {
			utils.Logger.Fatal("Oops!", "reason", err.Error())
		}
	},
}

type state int

const (
	stateDefault state = iota
	stateAddDrive
	stateRevokeDrive
	stateSelectDrive
	stateLoading
	stateDone
	stateError
	stateExit
)

type (
	authMsg struct {
		err      error
		provider string
		action   string
	}
)

type driveModel struct {
	config    *internal.RdvConfig
	state     state
	cursor    int
	err       error
	message   string
	loading   string
	providers []string
}

func newDriveModel(cfg *internal.RdvConfig) *driveModel {
	m := &driveModel{
		config: cfg,
		state:  stateDefault,
	}

	for _, p := range providers.GetAll() {
		m.providers = append(m.providers, p.Name())
	}

	return m
}

func (m *driveModel) Init() tea.Cmd {
	if len(m.config.Drives) == 0 {
		m.message = "No drives connected. Press 'a' to add one."
	}

	return nil
}

func (m *driveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case authMsg:
		return m.handleAuthMsg(msg)

	case tea.WindowSizeMsg:
		return m, nil
	}

	return m, nil
}

func (m *driveModel) handleAuthMsg(msg authMsg) (tea.Model, tea.Cmd) {
	m.state = stateDefault

	if msg.err != nil {
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	switch msg.action {
	case "add":
		m.message = fmt.Sprintf("✅ Added %s successfully", msg.provider)
	case "revoke":
		m.message = fmt.Sprintf("✅ Disconnected %s successfully", msg.provider)
	case "reconnect":
		m.message = fmt.Sprintf("✅ Reconnected %s successfully", msg.provider)

		for i := range m.config.Drives {
			if m.config.Drives[i].Name == msg.provider {
				m.config.Drives[i].Status = internal.Selected
			} else if m.config.Drives[i].Status == internal.Selected {
				m.config.Drives[i].Status = internal.Default
			}
		}
	}

	if err := m.config.SaveCfg(); err != nil {
		m.err = err
		m.state = stateError
		return m, nil
	}

	return m, nil
}

func (m *driveModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateDefault:
		return m.handleMain(msg)

	case stateAddDrive:
		return m.handleAddDrive(msg)

	case stateRevokeDrive:
		return m.handleRevokeDrive(msg)

	case stateSelectDrive:
		return m.handleSelectDrive(msg)

	case stateDone, stateError:
		return m.handleDone(msg)

	case stateExit:
		return m, tea.Quit
	}

	return m, nil
}

func (m *driveModel) handleMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.state = stateExit
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 && m.cursor < len(m.config.Drives) {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.config.Drives)-1 {
			m.cursor++
		}

	case "a":
		m.state = stateAddDrive
		m.cursor = 0
		return m, nil

	case "r":
		if len(m.config.Drives) == 0 {
			m.err = fmt.Errorf("no drives to revoke")
			m.state = stateError
			return m, nil
		}
		m.state = stateRevokeDrive
		m.cursor = 0
		return m, nil

	case "u", "enter":
		if len(m.config.Drives) == 0 {
			m.err = fmt.Errorf("no drives to select")
			m.state = stateError
			return m, nil
		}
		m.state = stateSelectDrive
		return m, nil
	}

	return m, nil
}

func (m *driveModel) handleAddDrive(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateDefault
		m.message = ""
		return m, nil

	case "enter":
		if len(m.providers) == 0 {
			return m, nil
		}
		providerName := m.providers[m.cursor]
		m.state = stateLoading
		m.loading = "Authorizing"
		return m, authorizeCmd(providerName, m.config, "add")

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.providers)-1 {
			m.cursor++
		}
	}
	return m, nil
}

func (m *driveModel) handleRevokeDrive(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateDefault
		return m, nil

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.config.Drives)-1 {
			m.cursor++
		}

	case "enter":
		drive := &m.config.Drives[m.cursor]
		if drive.Status == internal.Revoked {
			m.err = fmt.Errorf("%s is already disconnected", drive.Name)
			m.state = stateError
			return m, nil
		}

		m.state = stateLoading
		m.loading = fmt.Sprintf("Revoking %s", drive.Name)
		return m, revokeCmd(drive.Name, m.config)
	}

	return m, nil
}

func (m *driveModel) handleSelectDrive(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateDefault
		return m, nil

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.config.Drives)-1 {
			m.cursor++
		}

	case "enter":
		selected := &m.config.Drives[m.cursor]
		if selected.Status == internal.Revoked {
			m.state = stateLoading
			m.loading = fmt.Sprintf("Reconnecting %s", selected.Name)
			return m, authorizeCmd(selected.Name, m.config, "reconnect")
		}

		selected.Status = internal.Selected
		for i := range m.config.Drives {
			if i != m.cursor && m.config.Drives[i].Status == internal.Selected {
				m.config.Drives[i].Status = internal.Default
			}
		}

		if err := m.config.SaveCfg(); err != nil {
			m.err = err
			m.state = stateError
			return m, nil
		}

		m.message = fmt.Sprintf("Now using %s", selected.Name)
		m.state = stateDone
		return m, nil
	}

	return m, nil
}

func (m *driveModel) handleDone(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ", "esc", "q", "ctrl+c":
		m.state = stateDefault
		m.cursor = 0
		if m.err != nil {
			m.err = nil
		}
		return m, nil
	}

	return m, nil
}

func authorizeCmd(providerName string, cfg *internal.RdvConfig, action string) tea.Cmd {
	return func() tea.Msg {
		p := providers.Get(providerName)
		if p == nil {
			return authMsg{
				err:      fmt.Errorf("%s is not supported", providerName),
				action:   action,
				provider: providerName,
			}
		}

		for i, d := range cfg.Drives {
			if d.Name == providerName {
				if d.Status == internal.Revoked && action == "reconnect" {
					if err := oauth.Authorize(p); err != nil {
						return authMsg{
							err:      err,
							provider: providerName,
							action:   action,
						}
					}

					cfg.Drives[i].Status = internal.Default
					if err := cfg.SaveCfg(); err != nil {
						return authMsg{
							err:      err,
							provider: providerName,
							action:   action,
						}
					}

					return authMsg{
						err:      nil,
						provider: providerName,
						action:   action,
					}
				}

				return errMsg(fmt.Errorf("%s is already linked", providerName))
			}
		}

		if action != "add" {
			return authMsg{
				err:      fmt.Errorf("invalid action %s for new drive", action),
				provider: providerName,
				action:   action,
			}
		}

		if err := oauth.Authorize(p); err != nil {
			return authMsg{
				err:      err,
				provider: providerName,
				action:   action,
			}
		}

		config := p.GetConfig()
		cfg.Drives = append(cfg.Drives, internal.DriveProviderConfig{
			Name:   p.Name(),
			Id:     config.ClientID,
			Status: internal.Default,
		})

		if err := cfg.SaveCfg(); err != nil {
			return authMsg{
				err:      err,
				provider: providerName,
				action:   action,
			}
		}

		return authMsg{
			err:      nil,
			provider: providerName,
			action:   action,
		}
	}
}

func revokeCmd(providerName string, cfg *internal.RdvConfig) tea.Cmd {
	return func() tea.Msg {
		found := false
		for i, d := range cfg.Drives {
			if d.Name == providerName {
				found = true

				switch d.Status {
				case internal.Revoked:
					return authMsg{
						err:      fmt.Errorf("%s already disconnected", providerName),
						provider: providerName,
						action:   "revoke",
					}
				default:
					p := providers.Get(d.Name)
					err := oauth.RevokeToken(p)
					if err != nil {
						return authMsg{
							err:      err,
							provider: providerName,
							action:   "revoke",
						}
					}

					cfg.Drives[i].Status = internal.Revoked
					if err := cfg.SaveCfg(); err != nil {
						return authMsg{
							err:      err,
							provider: providerName,
							action:   "revoke",
						}
					}
				}
			}
		}

		if !found {
			p := providers.Get(providerName)
			if p == nil {
				return authMsg{
					err:      fmt.Errorf("%s is not supported by rdv", providerName),
					provider: providerName,
					action:   "revoke",
				}
			}

			// If token is present in keyring but configuration is empty
			if err := oauth.RevokeToken(p); err != nil {
				return authMsg{
					err:      err,
					provider: providerName,
					action:   "revoke",
				}
			}
		}

		return authMsg{
			err:      nil,
			provider: providerName,
			action:   "revoke",
		}
	}
}

func (m *driveModel) View() string {
	var s strings.Builder

	switch m.state {
	case stateDefault:
		s.WriteString(m.viewDefault())
	case stateAddDrive:
		s.WriteString(m.viewAddDrive())
	case stateRevokeDrive:
		s.WriteString(m.viewRevokeDrive())
	case stateSelectDrive:
		s.WriteString(m.viewSelectDrive())
	case stateLoading:
		s.WriteString(m.viewLoading())
	case stateDone:
		s.WriteString(m.viewDone())
	case stateError:
		s.WriteString(m.viewError())
	}

	if m.state != stateExit {
		s.WriteString(m.helpView())
	}

	return s.String()
}

// func (m *driveModel) viewDefault() string {
// 	var b strings.Builder

// 	// Connected drives section
// 	b.WriteString("📦 Connected Drives:\n")
// 	if len(m.config.Drives) == 0 {
// 		b.WriteString(defaultStyle.Render("   (none connected)") + "\n")
// 	} else {
// 		for i, d := range m.config.Drives {
// 			cursor := "  "
// 			if i == m.cursor {
// 				cursor = cursorStyle.Render("❯ ")
// 			}

// 			var status string
// 			switch d.Status {
// 			case internal.Selected:
// 				status = activeStyle.Render("[active]")
// 			case internal.Revoked:
// 				status = revokedStyle.Render("[disconnected]")
// 			default:
// 				status = defaultStyle.Render("[linked]")
// 			}

// 			b.WriteString(fmt.Sprintf("   %s%-20s %s\n", cursor, d.Name, status))
// 		}
// 	}
// 	b.WriteString("\n")

// 	// Supported providers preview
// 	b.WriteString("🔌 Supported Providers: ")
// 	b.WriteString(defaultStyle.Render(fmt.Sprintf("%d available", len(m.providers))))
// 	b.WriteString("\n\n")

// 	// Status message
// 	if m.message != "" {
// 		b.WriteString(successStyle.Render("✓ "+m.message) + "\n\n")
// 	}

// 	return b.String()
// }

// func (m *driveModel) viewAddDrive() string {
// 	var b strings.Builder

// 	b.WriteString("➕ Add New Drive\n")
// 	b.WriteString(subtitleStyle.Render("Select a provider to authorize") + "\n\n")

// 	if m.search != "" {
// 		b.WriteString(fmt.Sprintf("Search: %s\n\n", m.search))
// 	}

// 	for i, p := range m.filteredProviders {
// 		cursor := "  "
// 		if i == m.cursor {
// 			cursor = cursorStyle.Render("❯ ")
// 		}
// 		b.WriteString(fmt.Sprintf("   %s%s\n", cursor, p))
// 	}

// 	if len(m.filteredProviders) == 0 {
// 		b.WriteString(defaultStyle.Render("   No providers match your search") + "\n")
// 	}

// 	return b.String()
// }

// func (m *driveModel) viewRevokeDrive() string {
// 	var b strings.Builder

// 	b.WriteString("🔓 Revoke Drive Access\n")
// 	b.WriteString(subtitleStyle.Render("Select a drive to disconnect") + "\n\n")

// 	for i, d := range m.config.Drives {
// 		cursor := "  "
// 		if i == m.cursor {
// 			cursor = cursorStyle.Render("❯ ")
// 		}

// 		var status string
// 		switch d.Status {
// 		case internal.Revoked:
// 			status = revokedStyle.Render("[already disconnected]")
// 		default:
// 			status = "✓ linked"
// 		}

// 		b.WriteString(fmt.Sprintf("   %s%-20s %s\n", cursor, d.Name, status))
// 	}

// 	return b.String()
// }

// func (m *driveModel) viewSelectDrive() string {
// 	var b strings.Builder

// 	b.WriteString("🎯 Select Active Drive\n")
// 	b.WriteString(subtitleStyle.Render("Choose which drive to use") + "\n\n")

// 	for i, d := range m.config.Drives {
// 		cursor := "  "
// 		if i == m.cursor {
// 			cursor = cursorStyle.Render("❯ ")
// 		}

// 		var status string
// 		switch d.Status {
// 		case internal.Selected:
// 			status = activeStyle.Render("[currently active]")
// 		case internal.Revoked:
// 			status = revokedStyle.Render("[needs reconnection]")
// 		default:
// 			status = defaultStyle.Render("[available]")
// 		}

// 		b.WriteString(fmt.Sprintf("   %s%-20s %s\n", cursor, d.Name, status))
// 	}

// 	return b.String()
// }

// func (m *driveModel) viewLoading() string {
// 	var b strings.Builder
// 	b.WriteString("⏳ " + m.loading + "\n\n")
// 	b.WriteString(defaultStyle.Render("Please complete authorization in your browser...") + "\n")
// 	return b.String()
// }

// func (m *driveModel) viewDone() string {
// 	var b strings.Builder
// 	b.WriteString(successStyle.Render("✓ "+m.message) + "\n\n")
// 	b.WriteString(helpStyle.Render("Press any key to continue...") + "\n")
// 	return b.String()
// }

// func (m *driveModel) viewError() string {
// 	var b strings.Builder
// 	b.WriteString(errorStyle.Render("✗ Error: "+m.err.Error()) + "\n\n")
// 	b.WriteString(helpStyle.Render("Press any key to continue...") + "\n")
// 	return b.String()
// }

func (m *driveModel) helpView() string {
	switch m.state {
	case stateDefault:
		return helpStyle.Render(
			"↑↓/j,k: navigate  |  Enter: select  |  a: add  |  r: revoke  |  u: use  |  ?: help  |  q: quit",
		)
	case stateAddDrive, stateRevokeDrive, stateSelectDrive:
		return helpStyle.Render(
			"↑↓/j,k: navigate  |  Enter: confirm  |  Esc: back  |  q: quit",
		)
	default:
		return ""
	}
}

func init() {
	rootCmd.AddCommand(driveCmd)
}
