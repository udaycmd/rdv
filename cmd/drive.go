package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/internal/oauth/providers"
	"github.com/udaycmd/rdv/utils"
)

var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Manage remote drives.",
	Long:  "Add, remove, select, and manage your remote drive connections.",
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
	spinner   spinner.Model
	message   string
	loading   string
	providers []string
	width     int
	height    int
}

func newDriveModel(cfg *internal.RdvConfig) *driveModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0164c"))

	m := &driveModel{
		config:  cfg,
		state:   stateDefault,
		spinner: s,
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

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
		m.message = fmt.Sprintf("Added %s successfully", msg.provider)
	case "revoke":
		m.message = fmt.Sprintf("Disconnected %s successfully", msg.provider)
	case "reconnect":
		m.message = fmt.Sprintf("Reconnected %s successfully", msg.provider)

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
		return m.handleDefault(msg)

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

func (m *driveModel) handleDefault(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		return m, tea.Batch(authorizeCmd(providerName, m.config, "add"), m.spinner.Tick)

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
		return m, tea.Batch(revokeCmd(drive.Name, m.config), m.spinner.Tick)
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
			return m, tea.Batch(authorizeCmd(selected.Name, m.config, "reconnect"), m.spinner.Tick)
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
		msg := authMsg{
			err:      nil,
			provider: providerName,
			action:   action,
		}

		p := providers.Get(providerName)
		if p == nil {
			msg.err = fmt.Errorf("%s is not supported", providerName)
			return msg
		}

		for i, d := range cfg.Drives {
			if d.Name == providerName {
				if d.Status == internal.Revoked && action == "reconnect" {
					if err := oauth.Authorize(p); err != nil {
						msg.err = err
						return msg
					}

					cfg.Drives[i].Status = internal.Default
					if err := cfg.SaveCfg(); err != nil {
						msg.err = err
						return msg
					}

					return msg
				}

				msg.err = fmt.Errorf("%s is already linked", providerName)
				return msg
			}
		}

		if action != "add" {
			msg.err = fmt.Errorf("invalid action %s for new drive", action)
			return msg
		}

		if err := oauth.Authorize(p); err != nil {
			msg.err = err
			return msg
		}

		config := p.GetConfig()
		cfg.Drives = append(cfg.Drives, internal.DriveProviderConfig{
			Name:   p.Name(),
			Id:     config.ClientID,
			Status: internal.Default,
		})

		if err := cfg.SaveCfg(); err != nil {
			msg.err = err
			return msg
		}

		return msg
	}
}

func revokeCmd(providerName string, cfg *internal.RdvConfig) tea.Cmd {
	return func() tea.Msg {
		msg := authMsg{
			err:      nil,
			provider: providerName,
			action:   "revoke",
		}

		found := false
		for i, d := range cfg.Drives {
			if d.Name == providerName {
				found = true

				switch d.Status {
				case internal.Revoked:
					msg.err = fmt.Errorf("%s already disconnected", providerName)
					return msg

				default:
					p := providers.Get(d.Name)
					err := oauth.RevokeToken(p)
					if err != nil {
						msg.err = err
						return msg
					}

					cfg.Drives[i].Status = internal.Revoked
					if err := cfg.SaveCfg(); err != nil {
						msg.err = err
						return msg
					}
				}
			}
		}

		if !found {
			p := providers.Get(providerName)
			if p == nil {
				msg.err = fmt.Errorf("%s is not supported by rdv", providerName)
				return msg
			}

			// If token is present in keyring but configuration is empty
			if err := oauth.RevokeToken(p); err != nil {
				msg.err = err
				return msg
			}
		}

		return msg
	}
}

func (m *driveModel) View() string {
	var s string

	switch m.state {
	case stateDefault:
		s = m.viewDefault()
	case stateAddDrive:
		s = m.viewAddDrive()
	case stateRevokeDrive:
		s = m.viewRevokeDrive()
	case stateSelectDrive:
		s = m.viewSelectDrive()
	case stateLoading:
		s = utils.LoadingView(
			fmt.Sprintf("⏳ %s %s",
				m.loading,
				m.spinner.View(),
			),
			m.width,
			m.height,
		)
	case stateDone:
		s = utils.SuccessView(m.message, m.width, m.height)
	case stateError:
		s = utils.ErrorView(m.err, m.width, m.height)
	}

	if m.state != stateExit {
		s = lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			m.viewHelp(),
		)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		s,
	)
}

var (
	revokeStyle = lipgloss.NewStyle().Foreground(utils.Destructive).Bold(true)
	activeStyle = lipgloss.NewStyle().Foreground(utils.Constructive).Bold(true)
)

func (m *driveModel) viewDefault() string {
	var s string

	s = utils.HeaderStyle.Render("🔌 Connected Drives")
	if len(m.config.Drives) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			lipgloss.NewStyle().Faint(true).Render("(none)"),
			activeStyle.Render("press 'a' to connect a remote drive"),
		)
	}

	for i, d := range m.config.Drives {
		cursor := ""
		if i == m.cursor {
			cursor = utils.CursorStyle.Render("❯ ")
		}

		var status string
		switch d.Status {
		case internal.Selected:
			status = activeStyle.Render("(active)")
		case internal.Revoked:
			status = revokeStyle.Render("(disconnected)")
		default:
			status = lipgloss.NewStyle().Faint(true).Render("(linked)")
		}

		s = lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			fmt.Sprintf("%s%s %s", cursor, d.Name, status),
		)
	}

	return s
}

func (m *driveModel) viewAddDrive() string {
	var s string

	s = utils.HeaderStyle.Render("➕ Add New Drive")

	for i, p := range m.providers {
		cursor := ""
		if i == m.cursor {
			cursor = utils.CursorStyle.Render("❯ ")
		}

		s = lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			fmt.Sprintf("%s%s", cursor, p),
		)
	}

	return s
}

func (m *driveModel) viewRevokeDrive() string {
	var s string

	s = utils.HeaderStyle.Render("🔓 Revoke Drive Access")

	for i, d := range m.config.Drives {
		cursor := ""
		if i == m.cursor {
			cursor = utils.CursorStyle.Render("❯ ")
		}

		var status string
		switch d.Status {
		case internal.Revoked:
			status = revokeStyle.Render("(already disconnected)")
		}

		s = lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			fmt.Sprintf("%s%s %s", cursor, d.Name, status),
		)
	}

	return s
}

func (m *driveModel) viewSelectDrive() string {
	var s string

	s = utils.HeaderStyle.Render("✅ Select Active Drive")

	for i, d := range m.config.Drives {
		cursor := ""
		if i == m.cursor {
			cursor = utils.CursorStyle.Render("❯ ")
		}

		var status string
		switch d.Status {
		case internal.Selected:
			status = activeStyle.Render("(currently active)")
		case internal.Revoked:
			status = revokeStyle.Render("(needs reconnection)")
		default:
			status = lipgloss.NewStyle().Faint(true).Render("(linked)")
		}

		s = lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			fmt.Sprintf("%s%s %s", cursor, d.Name, status),
		)
	}

	return s
}

func (m *driveModel) viewHelp() string {
	switch m.state {
	case stateDefault:
		return utils.HelpStyle.Render(
			"↑↓/j,k: navigate • Enter: select • a: add • r: revoke • u: use • q: quit",
		)
	case stateAddDrive, stateRevokeDrive, stateSelectDrive:
		return utils.HelpStyle.Render(
			"↑↓/j,k: navigate • Enter: confirm • Esc: back",
		)
	default:
		return ""
	}
}

func init() {
	RootCmd.AddCommand(driveCmd)
}
