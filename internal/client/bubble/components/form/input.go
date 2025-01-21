package form

import "github.com/charmbracelet/bubbles/textinput"

type InputOptions struct {
	charLimit int
	password  bool
	textarea  bool
	validate  textinput.ValidateFunc
}

func WithCharLimit(limit int) func(*InputOptions) {
	return func(i *InputOptions) {
		i.charLimit = limit
	}
}

func IsPassword(password bool) func(*InputOptions) {
	return func(i *InputOptions) {
		i.password = password
	}
}

func IsTextArea(textarea bool) func(*InputOptions) {
	return func(i *InputOptions) {
		i.textarea = textarea
	}
}

func WithValidate(validate textinput.ValidateFunc) func(*InputOptions) {
	return func(i *InputOptions) {
		i.validate = validate
	}
}
