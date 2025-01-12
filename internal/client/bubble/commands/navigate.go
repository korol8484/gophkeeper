package commands

import tea "github.com/charmbracelet/bubbletea"

type NavigateCmd int

func GoTo(screenId int) tea.Msg {
	return NavigateCmd(screenId)
}
