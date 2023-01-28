package errs

import (
	"fmt"
	"os"

	"github.com/VladPetriv/setup-neovim/pkg/colors"
)

// WrapError used to print color error message and exit from program using os.Exit(1).
func WrapError(msg string, err error) {
	colors.Red(msg)
	colors.Red(fmt.Sprintf("Error: %v", err))
	os.Exit(0)
}
