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

func (i input) ProcessInputForPackageManagers(stdin io.Reader) (string, error) {
	fmt.Print("Do you have any package managers installed?(y/n): ")

	reader := bufio.NewReader(stdin)
	haveInstalledPackageManagers, err := i.GetInput(reader)
	if err != nil {
		return "", fmt.Errorf("failed to get user input: %w", err)
	}

	switch haveInstalledPackageManagers {
	case "y":
		return "skip", nil
	case "n":
		fmt.Print("Choose package manager(packer/vim-plug): ")

		return i.GetInput(reader)
	default:
		return "", nil
	}
}

func (i input) GetInput(stdin io.Reader) (string, error) {
	reader := bufio.NewReader(stdin)

	data, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("get input error: %w", err)
	}

	return strings.ReplaceAll(data, "\n", ""), nil
}
