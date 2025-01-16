package card

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
	number
	year
	month
	cvv
)

type Model struct {
	service *service.Client
	form    *form.Component
}

func NewCardScreen(service *service.Client) *Model {
	m := &Model{
		service: service,
		form:    form.NewComponent(),
	}

	m.form.AddInput(title, "Title", form.WithCharLimit(30))
	m.form.AddInput(number, "Number", form.WithCharLimit(16))
	m.form.AddInput(year, "Year", form.WithCharLimit(2))
	m.form.AddInput(month, "Month", form.WithCharLimit(2))
	m.form.AddInput(cvv, "cvv", form.WithCharLimit(3))

	m.form.AddButton("Save", m.save())
	m.form.AddButton("Back", m.back())

	return &Model{service: service}
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
		var cTitle, cNumber, cYear, cMonth, cCvv string

		if t, ok := vals[title]; ok {
			cTitle = t
		}

		if t, ok := vals[number]; ok {
			cTitle = t
		}

		if t, ok := vals[year]; ok {
			cTitle = t
		}

		if t, ok := vals[month]; ok {
			cTitle = t
		}

		if t, ok := vals[cvv]; ok {
			cTitle = t
		}

		for _, s := range []string{cTitle, cNumber, cYear, cMonth, cCvv} {
			if len(s) == 0 {
				return commands.WrapCmd(commands.Error("all field required"))
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.service.Save(ctx, model.NewCard(cTitle, cNumber, cYear, cMonth, cCvv))
		if err != nil {
			return commands.WrapCmd(commands.Error(err.Error()))
		}

		return commands.WrapCmd(commands.GoTo(screens.SecretsScreen))
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
