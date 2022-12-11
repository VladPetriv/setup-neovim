package colors

import "fmt"

var reset = "\033[0m" //nolint:gochecknoglobals // reset should be global variable to avoid code duplication.

func Red(msg string) {
	red := "\033[31m"
	fmt.Printf("%s%s%s\n", red, msg, reset)
}

func Green(msg string) {
	green := "\033[32m"
	fmt.Printf("%s%s%s\n", green, msg, reset)
}
