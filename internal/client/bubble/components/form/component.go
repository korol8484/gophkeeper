package form

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
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

type keyMap struct {
	up     key.Binding
	action key.Binding
	down   key.Binding
}

type Component struct {
	inputs         []inputModel
	buttons        []button
	focusIndex     int
	style          *style
	formComponents int
	keyMap         keyMap
}

func NewComponent() *Component {
	m := &Component{
		inputs: []inputModel{},
		style:  defaultStyle(),
		keyMap: defaultKeyMap(),
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

	if opt.validate != nil {
		t.Validate = opt.validate
	}

	if opt.password {
		t.EchoMode = textinput.EchoPassword
		t.EchoCharacter = '*'
	}

	if opt.textarea {
		t.TextStyle.MaxHeight(3)
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
		switch {
		case key.Matches(msg, c.keyMap.up):
			c.up()
			return c, tea.Batch(c.updateStyle()...)
		case key.Matches(msg, c.keyMap.action):
			return c, c.action()
		case key.Matches(msg, c.keyMap.down):
			c.down()
			return c, tea.Batch(c.updateStyle()...)
		}
	}

	return c, c.updateInputs(msg)
}

func (c *Component) Validate() []error {
	var err []error
	for j := range c.inputs {
		if c.inputs[j].Validate == nil {
			continue
		}

		if errStr := c.inputs[j].Validate(c.inputs[j].Value()); errStr != nil {
			err = append(err, errStr)
		}
	}

	return err
}

func (c *Component) action() tea.Cmd {
	if len(c.buttons) > 0 {
		bIdx := c.focusIndex - len(c.inputs)
		if bIdx > -1 {
			return c.buttons[bIdx].callback()
		}
	}

	return nil
}

func (c *Component) up() {
	c.focusIndex--
	if c.focusIndex < 0 {
		c.focusIndex = c.formComponents - 1
	}
}

func (c *Component) down() {
	c.focusIndex++
	if c.focusIndex >= c.formComponents {
		c.focusIndex = 0
	}
}

func (c *Component) updateStyle() []tea.Cmd {
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

	return cmds
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

func defaultKeyMap() keyMap {
	return keyMap{
		up:     key.NewBinding(key.WithKeys("up"), key.WithHelp("tab/shift+tab/↑", "up")),
		action: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "action")),
		down:   key.NewBinding(key.WithKeys("tab", "shift+tab", "down"), key.WithHelp("↓", "down")),
	}
}
