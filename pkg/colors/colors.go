package colors

import "fmt"

var reset = "\033[0m" //nolint

func Red(msg string) {
	red := "\033[31m"
	fmt.Printf("%s%s%s\n", string(red), msg, string(reset)) //nolint
}

func Green(msg string) {
	green := "\033[32m"
	fmt.Printf("%s%s%s\n", string(green), msg, string(reset)) //nolint
}
