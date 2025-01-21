package valitators

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
)

func Length(title string, min, max int) textinput.ValidateFunc {
	return func(s string) error {
		if len(s) < min || len(s) > max {
			return fmt.Errorf("%s length invalid", title)
		}

		return nil
	}
}
