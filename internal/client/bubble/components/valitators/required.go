package valitators

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
)

func Required(title string) textinput.ValidateFunc {
	return func(s string) error {
		if len(s) == 0 {
			return fmt.Errorf("%s can't by empty", title)
		}

		return nil
	}
}
