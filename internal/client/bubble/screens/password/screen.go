package password

import (
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/form"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/valitators"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens"
	"github.com/korol8484/gophkeeper/internal/client/model"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"github.com/korol8484/gophkeeper/pkg"
	"time"
)

const (
	title form.InputId = iota
	login
	password
)

type Model struct {
	service *service.Client
	form    *form.Component
}

func NewPasswordScreen(service *service.Client) *Model {
	m := &Model{
		service: service,
		form:    form.NewComponent(),
	}

	m.form.AddInput(title, "Title", form.WithCharLimit(30), form.WithValidate(valitators.Required("Title")))
	m.form.AddInput(login, "Login", form.WithCharLimit(30), form.WithValidate(valitators.Length("Login", 5, 30)))
	m.form.AddInput(
		password,
		"Password",
		form.WithCharLimit(11),
		form.IsPassword(true),
		form.WithValidate(valitators.Length("Password", 5, 11)),
	)

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
		fErr := m.form.Validate()
		if len(fErr) > 0 {
			return commands.WrapCmd(commands.Error(errors.Join(fErr...).Error()))
		}

		vals := m.form.Values()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.service.Save(ctx, model.NewPassword(vals[title], vals[login], vals[password]))
		if err != nil {
			return commands.NotifyMsg(err.Error(), pkg.TimeOut)
		}

		return tea.Batch(
			commands.NotifyMsg("New secret add success", pkg.TimeOut),
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
