package input

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type input struct{}

var _ Inputter = (*input)(nil)

func New() Inputter {
	return &input{}
}

func (i input) GetInput(stdin io.Reader) (string, error) {
	reader := bufio.NewReader(stdin)

	data, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("get input error: %w", err)
	}

	return strings.ReplaceAll(data, "\n", ""), nil
}
