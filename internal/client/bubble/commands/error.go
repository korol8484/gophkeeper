package commands

import tea "github.com/charmbracelet/bubbletea"

type ErrorCmd struct {
	message string
}

func (e ErrorCmd) Error() string {
	return e.message
}

func Error(msg string) tea.Msg {
	return ErrorCmd{message: msg}
}
