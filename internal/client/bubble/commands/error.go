package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type ErrorCmd struct {
	message string
}

func (e ErrorCmd) Error() string {
	return e.message
}

func Error(msg string) tea.Msg {
	return ErrorCmd{message: msg}
}

func ErrorMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		return Error(msg)
	}
}

type ClearErrorMsg struct{}

func ClearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}

func NotifyMsg(msg string, t time.Duration) tea.Cmd {
	return tea.Batch(
		ErrorMsg(msg),
		ClearErrorAfter(t),
	)
}
