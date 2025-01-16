package bubble

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
)

type WrapScreen interface {
	Update(tea.Msg) tea.Cmd
	View() string
	GetHelp() []key.Binding
}

type screenManager struct {
	screens   map[int]WrapScreen
	current   int
	err       error
	help      help.Model
	helpStyle lipgloss.Style

	width  int
	height int
}

func newManager() *screenManager {
	return &screenManager{
		screens:   make(map[int]WrapScreen),
		current:   0,
		help:      help.New(),
		helpStyle: lipgloss.NewStyle().Padding(1, 0, 0, 2),
	}
}

func (s *screenManager) AddScreen(id int, ss WrapScreen) *screenManager {
	s.screens[id] = ss

	return s
}

func (s *screenManager) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height - lipgloss.Height(s.helpView()) - lipgloss.Height(s.errView())
		s.help.Width = msg.Width

		cmd := s.screens[s.current].Update(tea.WindowSizeMsg{
			Width:  s.width,
			Height: s.height,
		})

		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		cmds = append(cmds, s.screens[s.current].Update(msg))
	case commands.NavigateCmd:
		s.err = nil
		s.current = int(msg)
		cmds = append(cmds, s.screens[s.current].Update(msg))
	case commands.ErrorCmd:
		s.err = msg
	case commands.ClearErrorMsg:
		s.err = nil
	default:
		if ss, ok := s.screens[s.current]; ok {
			return ss.Update(msg)
		}
	}

	return tea.Batch(cmds...)
}

func (s *screenManager) View() string {
	ss, ok := s.screens[s.current]
	if !ok {
		return ""
	}

	contents := []string{
		lipgloss.NewStyle().
			Width(s.width).
			Height(s.height).
			BorderStyle(lipgloss.NormalBorder()).Render(ss.View()),
	}

	if s.err != nil {
		contents = append(contents, s.errView())
	}

	helpView := s.helpView()
	if len(helpView) > 0 {
		contents = append(contents, helpView)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		contents...,
	)
}

func (s *screenManager) ShortHelp() []key.Binding {
	base := s.getBaseHelp()

	ss, ok := s.screens[s.current]
	if !ok {
		return base
	}

	return append(base, ss.GetHelp()...)
}

func (s *screenManager) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{s.getBaseHelp()}

	ss, ok := s.screens[s.current]
	if !ok {
		return kb
	}

	return append(kb, ss.GetHelp())
}

func (s *screenManager) getBaseHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
}

func (s *screenManager) helpView() string {
	return s.helpStyle.Render(s.help.View(s))
}

func (s *screenManager) errView() string {
	if s.err != nil {
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Width(s.width)

		return style.Render(s.err.Error())
	}

	return ""
}
