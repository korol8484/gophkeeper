package text

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
	"time"
)

const (
	title form.InputId = iota
	text
)

type Model struct {
	keyMap  keyMap
	service *service.Client
	form    *form.Component
}

type keyMap struct {
	lineUp   key.Binding
	lineDown key.Binding
	enter    key.Binding
}

func NewTextScreen(service *service.Client) *Model {
	m := &Model{
		service: service,
		form:    form.NewComponent(),
		keyMap:  defaultKeyMap(),
	}

	m.form.AddInput(title, "Title", form.WithCharLimit(30), form.WithValidate(valitators.Required("Title")))
	m.form.AddInput(text, "Text", form.WithValidate(valitators.Required("Text")))

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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		vals := m.form.Values()
		err := m.service.Save(ctx, model.NewText(vals[title], vals[text]))
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
		m.keyMap.lineUp, m.keyMap.lineDown, m.keyMap.enter,
	}
}

func defaultKeyMap() keyMap {
	return keyMap{
		lineUp:   key.NewBinding(key.WithKeys("up"), key.WithHelp("↑/k", "up")),
		lineDown: key.NewBinding(key.WithKeys("down"), key.WithHelp("↓/j", "down")),
		enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
	}
}
