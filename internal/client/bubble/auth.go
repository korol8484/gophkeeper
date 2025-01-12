package bubble

//
//import (
//	"fmt"
//	"github.com/charmbracelet/bubbles/textinput"
//	tea "github.com/charmbracelet/bubbletea"
//	"strings"
//)
//
//type authModel struct {
//	focusIndex int
//	inputs     []textinput.Model
//	login      string
//	password   string
//	success    bool
//
//	style *style
//}
//
//func initAuthModel() *authModel {
//	m := &authModel{
//		inputs: make([]textinput.Model, 2),
//		style:  defaultStyle(),
//	}
//
//	var t textinput.Model
//	for i := range m.inputs {
//		t = textinput.New()
//		t.Cursor.Style = m.style.focusedStyle
//		t.CharLimit = 32
//
//		switch i {
//		case 0:
//			t.Placeholder = "Login"
//			t.Focus()
//			t.PromptStyle = m.style.focusedStyle
//			t.TextStyle = m.style.focusedStyle
//		case 1:
//			t.Placeholder = "Password"
//			t.EchoMode = textinput.EchoPassword
//			t.EchoCharacter = '*'
//		}
//
//		m.inputs[i] = t
//	}
//
//	return m
//}
//
//func (i *authModel) Init() tea.Cmd {
//	return nil
//}
//
//func (i *authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch key := msg.String(); key {
//		case "tab", "shift+tab", "enter", "up", "down":
//			s := msg.String()
//
//			// Did the user press enter while the submit button was focused?
//			// If so, exit.
//			if s == "enter" && i.focusIndex == len(i.inputs) {
//				i.login = i.inputs[0].Value()
//				i.password = i.inputs[1].Value()
//				i.success = true
//
//				return i, nil
//			}
//
//			// Cycle indexes
//			if s == "up" || s == "shift+tab" {
//				i.focusIndex--
//			} else {
//				i.focusIndex++
//			}
//
//			if i.focusIndex > len(i.inputs) {
//				i.focusIndex = 0
//			} else if i.focusIndex < 0 {
//				i.focusIndex = len(i.inputs)
//			}
//
//			cmds := make([]tea.Cmd, len(i.inputs))
//			for j := 0; j <= len(i.inputs)-1; j++ {
//				if j == i.focusIndex {
//					// Set focused state
//					cmds[j] = i.inputs[j].Focus()
//					i.inputs[j].PromptStyle = i.style.focusedStyle
//					i.inputs[j].TextStyle = i.style.focusedStyle
//					continue
//				}
//				// Remove focused state
//				i.inputs[j].Blur()
//				i.inputs[j].PromptStyle = i.style.noStyle
//				i.inputs[j].TextStyle = i.style.noStyle
//			}
//
//			return i, tea.Batch(cmds...)
//		}
//	}
//
//	cmd := i.updateInputs(msg)
//
//	return i, cmd
//}
//
//func (i *authModel) Login() string {
//	return i.login
//}
//
//func (i *authModel) Password() string {
//	return i.password
//}
//
//func (i *authModel) IsSuccess() bool {
//	return i.success
//}
//
//func (i *authModel) updateInputs(msg tea.Msg) tea.Cmd {
//	cmds := make([]tea.Cmd, len(i.inputs))
//
//	for j := range i.inputs {
//		i.inputs[j], cmds[j] = i.inputs[j].Update(msg)
//	}
//
//	return tea.Batch(cmds...)
//}
//
//func (i *authModel) View() string {
//	var b strings.Builder
//
//	for j := range i.inputs {
//		b.WriteString(i.inputs[j].View())
//		if j < len(i.inputs)-1 {
//			b.WriteRune('\n')
//		}
//	}
//
//	button := fmt.Sprintf("[ %s ]", i.style.blurredStyle.Render("Submit"))
//	if i.focusIndex == len(i.inputs) {
//		button = i.style.focusedStyle.Render("[ Submit ]")
//	}
//	_, _ = fmt.Fprintf(&b, "\n\n%s\n\n", button)
//
//	return b.String()
//}
