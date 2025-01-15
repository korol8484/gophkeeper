package bubble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens/add"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens/auth"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens/binary"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens/card"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens/password"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens/secrets"
	"github.com/korol8484/gophkeeper/internal/client/bubble/screens/text"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"os"
	"os/exec"
)

type model struct {
	screen  *screenManager
	style   lipgloss.Style
	service *service.Client
}

func (m *model) Init() tea.Cmd {
	m.screen.AddScreen(screens.AuthScreen, auth.NewAuthScreen(m.service))
	m.screen.AddScreen(screens.SecretsScreen, secrets.NewSecretsScreen(m.service))
	m.screen.AddScreen(screens.PasswordScreen, password.NewPasswordScreen(m.service))
	m.screen.AddScreen(screens.AddScreen, add.NewAddScreen())
	m.screen.AddScreen(screens.TextScreen, text.NewTextScreen(m.service))
	m.screen.AddScreen(screens.CardScreen, card.NewCardScreen(m.service))
	m.screen.AddScreen(screens.BinaryScreen, binary.NewFilePickerScreen(m.service))

	m.style = lipgloss.NewStyle().
		Margin(2)

	return func() tea.Msg {
		return commands.GoTo(screens.AuthScreen)
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = m.style.Width(msg.Width)
		cmd := m.screen.Update(tea.WindowSizeMsg{
			Width:  m.style.GetWidth() - 4,
			Height: msg.Height,
		})

		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		if msg.String() == containsKey(msg.String(), []string{"q", "ctrl+c"}) {
			return m, tea.Quit
		}

		return m, m.screen.Update(msg)
	default:
		cmds = append(cmds, m.screen.Update(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return m.style.Render(m.screen.View())
}

func Run(client *service.Client) error {
	clearScreen()

	m := &model{
		screen:  newManager(),
		service: client,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		return err
	}

	return nil
}

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	_ = c.Run()
}

func containsKey(v string, a []string) string {
	for _, i := range a {
		if i == v {
			return v
		}
	}
	return ""
}
