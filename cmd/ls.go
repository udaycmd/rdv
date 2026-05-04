package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal/drive"
	"github.com/udaycmd/rdv/utils"
)

var dirId string

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List contents of a drive.",
	Long:  "List contents of a drive.",
	Run: func(cmd *cobra.Command, args []string) {
		d := config.GetSelectedDrive()
		if d == nil {
			utils.Logger.Fatal("Oops!", "reason", "No selected drive found")
		}

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeoutPeriod)
		defer cancel()

		drive, err := drive.NewDriveFromProvider(ctx, d.Name)
		if err != nil {
			utils.Logger.Fatal("Oops!", "reason", err.Error())
		}

		p := tea.NewProgram(
			newLsModel(drive, dirId),
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
		if _, err := p.Run(); err != nil {
			utils.Logger.Fatal("Oops!", "reason", err.Error())
		}
	},
}

type lsModel struct {
	drive      drive.Drive
	table      table.Model
	spinner    spinner.Model
	dirName    string
	dirId      string
	navigation []navPair
	loading    bool
	err        error
	width      int
	height     int
}

type (
	errMsg  error
	recvMsg []drive.Meta
	navPair struct {
		id   string
		name string
	}
)

func newLsModel(drive drive.Drive, dir string) *lsModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0164c"))

	dirName := ""
	if dir == "root" {
		dirName = "Root"
	}

	return &lsModel{
		drive:      drive,
		dirId:      dir,
		dirName:    dirName,
		spinner:    s,
		loading:    true,
		navigation: []navPair{},
	}
}

func (m *lsModel) fetchFiles() tea.Msg {
	files, err := m.drive.View(m.dirId)
	if err != nil {
		return errMsg(err)
	}

	return recvMsg(files)
}

func (m *lsModel) createTable(files []drive.Meta) table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 40},
		{Title: "Name", Width: 30},
		{Title: "Size", Width: 15},
		{Title: "Type", Width: 20},
		{Title: "Last Modified", Width: 20},
	}

	rows := []table.Row{}

	parentDirName := ""
	if len(m.navigation) > 0 {
		parentDirName = ".. 📁 " + m.navigation[len(m.navigation)-1].name
		rows = append(rows, table.Row{
			m.navigation[len(m.navigation)-1].id,
			parentDirName,
			"",
			"Parent",
			"",
		})
	}

	for _, f := range files {
		size := fmt.Sprintf("%d kb", f.Size/1024)
		modTime := f.LastModified.Format(time.DateTime)

		if f.IsDir {
			f.Name = "📁 " + f.Name
		} else {
			f.Name = "📄 " + f.Name
		}

		rows = append(rows, table.Row{
			f.Id,
			f.Name,
			size,
			f.MimeType,
			modTime,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		Foreground(lipgloss.Color("#665685")).
		Bold(true)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("#c70c28")).
		Bold(true)

	t.SetStyles(tableStyles)
	return t
}

func (m *lsModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchFiles,
	)
}

func (m *lsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	var navigateBack = func(m *lsModel) (tea.Model, tea.Cmd) {
		m.dirId = m.navigation[len(m.navigation)-1].id
		m.dirName = m.navigation[len(m.navigation)-1].name
		m.navigation = m.navigation[:len(m.navigation)-1]
		m.loading = true
		return m, m.fetchFiles
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if !m.loading {
				selected := m.table.SelectedRow()
				if len(selected) >= 4 {
					mimeType := selected[3]
					if mimeType == "Parent" {
						return navigateBack(m)
					}

					id := selected[0]
					if mimeType == "application/vnd.google-apps.folder" {
						m.navigation = append(m.navigation, navPair{id: m.dirId, name: m.dirName})
						m.dirId = id
						m.dirName = selected[1]
						m.loading = true
						return m, m.fetchFiles
					}
				}
			}

		case "backspace":
			if !m.loading && len(m.navigation) > 0 {
				return navigateBack(m)
			}
		}

	case recvMsg:
		m.loading = false
		m.table = m.createTable(msg)
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg
		return m, nil
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *lsModel) View() string {
	if m.err != nil {
		return utils.ErrorView(m.err, m.width, m.height)
	}

	if m.loading {
		return utils.LoadingView(
			fmt.Sprintf("⏳ Fetching files for %s directory %s",
				m.dirName,
				m.spinner.View(),
			),
			m.width,
			m.height,
		)
	}

	return utils.ListView(m.table, m.dirName, m.width, m.height)
}

func init() {
	RootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringVarP(&dirId, "dir", "d", "root", "unique id of the directory")
}
