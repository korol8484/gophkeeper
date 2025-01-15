package add

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens"
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type Model struct {
	list list.Model
}

func NewAddScreen() *Model {
	items := []list.Item{
		item("Password"),
		item("Text"),
		item("Card"),
		item("Binary"),
		item("Back"),
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 0, 0)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Title = "Select you gophKeeper type to create"

	return &Model{
		list: l,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if !ok {
				return nil
			}

			switch string(i) {
			case "Password":
				return commands.WrapCmd(commands.GoTo(screens.PasswordScreen))
			case "Text":
				return commands.WrapCmd(commands.GoTo(screens.TextScreen))
			case "Card":
				return commands.WrapCmd(commands.GoTo(screens.CardScreen))
			case "Binary":
			default:
				return commands.WrapCmd(commands.GoTo(screens.SecretsScreen))
			}
		}
	}

	return nil
}

func (m *Model) View() string {
	return m.list.View()
}
