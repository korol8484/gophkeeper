package form

type InputOptions struct {
	charLimit int
	password  bool
	textarea  bool
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
