package input

import (
	"fmt"
)

type input struct{}

var _ Inputter = (*input)(nil)

func New() *input { //nolint
	return &input{}
}

func (i input) GetInput() (string, error) {
	var data string
	if _, err := fmt.Scanln(&data); err != nil {
		return "", fmt.Errorf("get input error: %w", err)
	}

	return data, nil
}
