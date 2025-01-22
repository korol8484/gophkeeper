package form

import "github.com/charmbracelet/lipgloss"

type style struct {
	focusedStyle        lipgloss.Style
	blurredStyle        lipgloss.Style
	noStyle             lipgloss.Style
	cursorModeHelpStyle lipgloss.Style
}

func defaultStyle() *style {
	return &style{
		focusedStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		blurredStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		noStyle:             lipgloss.NewStyle(),
		cursorModeHelpStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
	}
}
