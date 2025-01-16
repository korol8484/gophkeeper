package secrets

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"time"
)

type Model struct {
	service *service.Client
	table   table.Model
	style   lipgloss.Style
}

func NewSecretsScreen(service *service.Client) *Model {
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Title", Width: 30},
		{Title: "Type", Width: 15},
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
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
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
			return tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		case "u":
			return tea.Batch(
				commands.ErrorMsg("blaaat"),
				commands.ClearErrorAfter(10*time.Second),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return cmd
}

func (m *Model) View() string {
	return m.table.View()
}

func (m *Model) GetHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithHelp("↑/k", "up")),
		key.NewBinding(key.WithHelp("↓/j", "down")),
		key.NewBinding(key.WithHelp("←/h/pgup", "prev page")),
		key.NewBinding(key.WithHelp("→/l/pgdn", "next page")),
		key.NewBinding(key.WithHelp("enter", "apply")),
		key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "sync data")),
	}
}
