package text

import (
	"context"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/form"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens"
	"github.com/korol8484/gophkeeper/internal/client/model"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"time"
)

const (
	title form.InputId = iota
	text
)

type Model struct {
	service *service.Client
	form    *form.Component
}

func NewTextScreen(service *service.Client) *Model {
	m := &Model{
		service: service,
		form:    form.NewComponent(),
	}

	m.form.AddInput(title, "Title", form.WithCharLimit(30))
	m.form.AddInput(text, "Text", form.WithCharLimit(5000), form.IsTextArea(true))

	m.form.AddButton("Save", m.save())
	m.form.AddButton("Back", m.back())

	return m
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	ff, cmd := m.form.Update(msg)
	m.form = ff.(*form.Component)

	return cmd
}

func (m *Model) View() string {
	return m.form.View()
}

func (m *Model) save() func() tea.Cmd {
	return func() tea.Cmd {
		vals := m.form.Values()

		var vp, vt string

		if p, ok := vals[text]; ok {
			vp = p
		}

		if t, ok := vals[title]; ok {
			vt = t
		}

		if len(vp) == 0 || len(vt) == 0 {
			return commands.WrapCmd(commands.Error("all field required"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.service.Save(ctx, model.NewText(vt, vp))
		if err != nil {
			return commands.NotifyMsg(err.Error(), 5*time.Second)
		}

		return tea.Batch(
			commands.NotifyMsg("New secret add success", 5*time.Second),
			commands.WrapCmd(commands.GoTo(screens.SecretsScreen)),
		)
	}
}

func (m *Model) back() func() tea.Cmd {
	return func() tea.Cmd {
		return commands.WrapCmd(commands.GoTo(screens.SecretsScreen))
	}
}

func (m *Model) GetHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("up"), key.WithHelp("↑/k", "up")),
		key.NewBinding(key.WithKeys("down"), key.WithHelp("↓/j", "down")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
	}
}
