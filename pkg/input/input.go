package input

import (
	"fmt"
)

type input struct{}

var _ Inputter = (*input)(nil)

func New() *input {
	return &input{}
}

func (i input) GetInput(msg string) (string, error) {
	var data string
	fmt.Println(msg + ":")
	_, err := fmt.Scanln(&data)
	if err != nil {
		return "", fmt.Errorf("get input error: %w", err)
	}

	return data, nil
}
