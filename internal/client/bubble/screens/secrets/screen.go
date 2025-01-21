package secrets

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens"
	cliModel "github.com/korol8484/gophkeeper/internal/client/model"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"time"
)

type viewTableModel interface {
	cliModel.BaseI
	GetTitle() string
	GetId() uuid.UUID
}

type Model struct {
	service  *service.Client
	table    table.Model
	style    lipgloss.Style
	listData []viewTableModel
	init     bool
}

func NewSecretsScreen(service *service.Client) *Model {
	columns := []table.Column{
		{Title: "ID", Width: 36},
		{Title: "Title", Width: 40},
		{Title: "Type", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &Model{
		service: service,
		table:   t,
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
		init: false,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	_, ok := msg.(tea.WindowSizeMsg)

	if !m.init && !ok {
		m.loadTable()
		m.init = true
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "enter":
			m.table.UpdateViewport()
			return tea.Batch(
				commands.ErrorMsg(fmt.Sprintf("Let's go to %s!", m.table.SelectedRow()[1])),
				commands.ClearErrorAfter(5*time.Second),
			)
		case "u":
			return m.loadTable()
		case "a":
			return commands.WrapCmd(commands.GoTo(screens.AddScreen))
		}
	}

	m.table, cmd = m.table.Update(msg)
	return cmd
}

func (m *Model) View() string {
	content := []string{
		m.table.View(),
	}

	selectedR := m.table.SelectedRow()
	if selectedR != nil {
		var selected viewTableModel
		for _, v := range m.listData {
			if v.GetId().UUID.String() == selectedR[0] {
				selected = v
				break
			}
		}

		if selected != nil {
			content = append(content, lipgloss.NewStyle().MarginLeft(2).Render(selected.View()))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, content...)
}

func (m *Model) GetHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "sync data")),
		key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add new data")),
		key.NewBinding(key.WithKeys("up"), key.WithHelp("↑/k", "up")),
		key.NewBinding(key.WithKeys("down"), key.WithHelp("↓/j", "down")),
		//key.NewBinding(key.WithHelp("←/h/pgup", "prev page")),
		//key.NewBinding(key.WithHelp("→/l/pgdn", "next page")),
		//key.NewBinding(key.WithHelp("enter", "apply")),
	}
}

func (m *Model) loadTable() tea.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listData, err := m.service.Load(ctx)
	if err != nil {
		return tea.Batch(
			commands.ErrorMsg(err.Error()),
			commands.ClearErrorAfter(5*time.Second),
		)
	}

	var rows []table.Row
	for _, v := range listData {
		if tm, ok := v.(viewTableModel); ok {
			m.listData = append(m.listData, tm)

			rows = append(rows, table.Row{
				tm.GetId().UUID.String(),
				tm.GetTitle(),
				string(tm.GetType()),
			})
		}
	}

	m.table.SetRows(rows)

	return tea.Batch(
		commands.ErrorMsg("Data updated"),
		commands.ClearErrorAfter(5*time.Second),
	)
}
