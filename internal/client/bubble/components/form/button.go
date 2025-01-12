package form

import tea "github.com/charmbracelet/bubbletea"

type button struct {
	text     string
	callback func() tea.Cmd
}
