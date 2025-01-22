package auth

import (
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/form"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/valitators"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"time"
)

const (
	login form.InputId = iota
	passWord
)

type Model struct {
	form    *form.Component
	service *service.Client
}

func NewAuthScreen(service *service.Client) *Model {
	m := &Model{
		service: service,
		form:    form.NewComponent(),
	}

	m.form.AddInput(login, "Login", form.WithCharLimit(30), form.WithValidate(valitators.Length("Login", 1, 30)))
	m.form.AddInput(
		passWord,
		"Password",
		form.WithCharLimit(11),
		form.IsPassword(true),
		form.WithValidate(valitators.Length("Password", 1, 11)),
	)
	m.form.AddButton("Login", m.login())
	m.form.AddButton("Register", m.register())

	return m
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	ff, cmd := m.form.Update(msg)
	m.form = ff.(*form.Component)

	return cmd
}

func (m *Model) login() func() tea.Cmd {
	return func() tea.Cmd {
		if c := m.validate(); c != nil {
			return c
		}

		vals := m.form.Values()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.service.Auth(ctx, vals[login], vals[passWord])
		if err != nil {
			return commands.WrapCmd(commands.Error(err.Error()))
		}

		return commands.WrapCmd(commands.GoTo(screens.SecretsScreen))
	}
}

func (m *Model) register() func() tea.Cmd {
	return func() tea.Cmd {
		if c := m.validate(); c != nil {
			return c
		}

		vals := m.form.Values()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.service.Register(ctx, vals[login], vals[passWord])
		if err != nil {
			return commands.WrapCmd(commands.Error(err.Error()))
		}

		return commands.WrapCmd(commands.GoTo(screens.SecretsScreen))
	}
}

func (m *Model) validate() tea.Cmd {
	fErr := m.form.Validate()
	if len(fErr) > 0 {
		return commands.WrapCmd(commands.Error(errors.Join(fErr...).Error()))
	}

	return nil
}

func (m *Model) View() string {
	return m.form.View()
}

func (m *Model) GetHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("up"), key.WithHelp("↑/k", "up")),
		key.NewBinding(key.WithKeys("down"), key.WithHelp("↓/j", "down")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
	}
}
