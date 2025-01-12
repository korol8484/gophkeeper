package auth

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/form"
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

	m.form.AddInput(login, "Login", form.WithCharLimit(30))
	m.form.AddInput(passWord, "Password", form.WithCharLimit(11), form.IsPassword(true))
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
		l, p := m.loadVals()

		if len(l) == 0 || len(p) == 0 {
			return commands.WrapCmd(commands.Error("login or password empty"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.service.Auth(ctx, l, p)
		if err != nil {
			return commands.WrapCmd(commands.Error(err.Error()))
		}

		return commands.WrapCmd(commands.GoTo(screens.SecretsScreen))
	}
}

func (m *Model) register() func() tea.Cmd {
	return func() tea.Cmd {
		l, p := m.loadVals()

		if len(l) == 0 || len(p) == 0 {
			return commands.WrapCmd(commands.Error("login or password empty"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.service.Register(ctx, l, p)
		if err != nil {
			return commands.WrapCmd(commands.Error(err.Error()))
		}

		return commands.WrapCmd(commands.GoTo(screens.SecretsScreen))
	}
}

func (m *Model) loadVals() (string, string) {
	var vl, vp string
	vals := m.form.Values()

	if l, ok := vals[login]; ok {
		vl = l
	}

	if p, ok := vals[passWord]; ok {
		vp = p
	}

	return vl, vp
}

func (m *Model) View() string {
	return m.form.View()
}
