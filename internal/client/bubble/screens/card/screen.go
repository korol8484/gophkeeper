package card

import (
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/form"
	"github.com/korol8484/gophkeeper/internal/client/bubble/components/valitators"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens"
	"github.com/korol8484/gophkeeper/internal/client/model"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"strconv"
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

	m.form.AddInput(title, "Title", form.WithCharLimit(30), form.WithValidate(valitators.Required("Title")))
	m.form.AddInput(number, "Number", form.WithCharLimit(16), form.WithValidate(m.validateNumber()))
	m.form.AddInput(year, "Year", form.WithCharLimit(2), form.WithValidate(valitators.Length("Year", 2, 2)))
	m.form.AddInput(month, "Month", form.WithCharLimit(2), form.WithValidate(valitators.Length("Month", 2, 2)))
	m.form.AddInput(cvv, "cvv", form.WithCharLimit(3), form.WithValidate(valitators.Length("cvv", 3, 3)))

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

		err := m.service.Save(ctx, model.NewCard(vals[title], vals[number], vals[year], vals[month], vals[cvv]))
		if err != nil {
			return commands.NotifyMsg(err.Error(), 5*time.Second)
		}

		return tea.Batch(
			commands.NotifyMsg("New secret add success", 5*time.Second),
			commands.WrapCmd(commands.GoTo(screens.SecretsScreen)),
		)
	}
}

func (m *Model) validateNumber() textinput.ValidateFunc {
	return func(number string) error {
		var sum int
		var alternate bool

		numberLen := len(number)
		if numberLen < 13 || numberLen > 19 {
			return errors.New("card number length invalid")
		}

		for i := numberLen - 1; i > -1; i-- {
			mod, _ := strconv.Atoi(string(number[i]))
			if alternate {
				mod *= 2
				if mod > 9 {
					mod = (mod % 10) + 1
				}
			}

			alternate = !alternate

			sum += mod
		}

		if sum%10 != 0 {
			return errors.New("card number invalid")
		}

		return nil
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
