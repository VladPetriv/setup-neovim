package colors

import "fmt"

var (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func Red(msg string) {
	fmt.Println(string(red), msg, string(reset))
}

func Green(msg string) {
	fmt.Println(string(green), msg, string(reset))
}
