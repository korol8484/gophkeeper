package form

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type InputId int
type Values map[InputId]string
type inputModel struct {
	textinput.Model
	id InputId
}

type Component struct {
	inputs         []inputModel
	buttons        []button
	focusIndex     int
	style          *style
	formComponents int
}

func NewComponent() *Component {
	m := &Component{
		inputs: []inputModel{},
		style:  defaultStyle(),
	}

	return m
}

func (c *Component) AddButton(title string, callback func() tea.Cmd) *Component {
	c.buttons = append(c.buttons, button{
		text:     title,
		callback: callback,
	})

	c.formComponents++

	return c
}

func (c *Component) AddInput(id InputId, title string, opts ...func(*InputOptions)) *Component {
	opt := &InputOptions{
		charLimit: 30,
	}
	for _, o := range opts {
		o(opt)
	}

	t := textinput.New()
	t.Placeholder = title
	t.CharLimit = opt.charLimit
	t.Cursor.Style = c.style.focusedStyle

	if len(c.inputs) == 0 {
		t.Focus()
		t.PromptStyle = c.style.focusedStyle
		t.TextStyle = c.style.focusedStyle
	}

	if opt.password {
		t.EchoMode = textinput.EchoPassword
		t.EchoCharacter = '*'
	}

	c.inputs = append(c.inputs, inputModel{
		Model: t,
		id:    id,
	})

	c.formComponents++

	return c
}

func (c *Component) Init() tea.Cmd {
	return nil
}

func (c *Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" && len(c.buttons) > 0 {
				bIdx := c.focusIndex - len(c.inputs)
				if bIdx > -1 {
					return c, c.buttons[bIdx].callback()
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				c.focusIndex--
			} else {
				c.focusIndex++
			}

			if c.focusIndex >= c.formComponents {
				c.focusIndex = 0
			} else if c.focusIndex < 0 {
				c.focusIndex = c.formComponents - 1
			}

			cmds := make([]tea.Cmd, len(c.inputs))
			for j := 0; j <= len(c.inputs)-1; j++ {
				if j == c.focusIndex {
					// Set focused state
					cmds[j] = c.inputs[j].Focus()

					c.inputs[j].PromptStyle = c.style.focusedStyle
					c.inputs[j].TextStyle = c.style.focusedStyle
					continue
				}

				// Remove focused state
				c.inputs[j].Blur()
				c.inputs[j].PromptStyle = c.style.noStyle
				c.inputs[j].TextStyle = c.style.noStyle
			}

			return c, tea.Batch(cmds...)
		}
	}

	cmd := c.updateInputs(msg)

	return c, cmd
}

func (c *Component) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	for j := range c.inputs {
		c.inputs[j].Model, cmds[j] = c.inputs[j].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (c *Component) View() string {
	var b strings.Builder

	for j := range c.inputs {
		b.WriteString(c.inputs[j].View())
		if j < len(c.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	if len(c.buttons) > 0 {
		b.WriteRune('\n')
	}

	for j, bt := range c.buttons {
		if c.focusIndex == len(c.inputs)+j {
			b.WriteString(fmt.Sprintf("[ %s ]", c.style.focusedStyle.Render(bt.text)))
		} else {
			b.WriteString(fmt.Sprintf("[ %s ]", c.style.blurredStyle.Render(bt.text)))
		}
	}

	return b.String()
}

func (c *Component) Values() Values {
	vals := make(Values, len(c.inputs))
	for _, inp := range c.inputs {
		vals[inp.id] = inp.Value()
	}

	return vals
}
