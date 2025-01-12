package bubble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/korol8484/gophkeeper/internal/client/bubble/commands"
)

type WrapScreen interface {
	Update(tea.Msg) tea.Cmd
	View() string
}

type screenManager struct {
	screens map[int]WrapScreen
	current int
	err     error

	width  int
	height int
}

func newManager() *screenManager {
	return &screenManager{
		screens: make(map[int]WrapScreen),
		current: 0,
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
		s.height = msg.Height

		cmd := s.screens[s.current].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height,
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
			BorderStyle(lipgloss.NormalBorder()).Render(ss.View()),
	}

	if s.err != nil {
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Width(s.width)

		contents = append(contents, style.Render(s.err.Error()))
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		contents...,
	)
}
