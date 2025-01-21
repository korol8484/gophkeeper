package commands

import tea "github.com/charmbracelet/bubbletea"

func WrapCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
