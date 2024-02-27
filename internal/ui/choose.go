package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func Choose(q string, options ...string) (int, error) {
	var val string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(q).
				Value(&val).
				Options(huh.NewOptions[string](options...)...),
		),
	)

	if err := form.Run(); err != nil {
		return -1, err
	}

	for i, o := range options {
		if o == val {
			return i, nil
		}
	}
	return -1, fmt.Errorf("invalid option %q", val)
}
