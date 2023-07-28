package colors

import "fmt"

const reset = "\033[0m"

func Red(msg string) {
	red := "\033[31m"
	fmt.Printf("%s%s%s\n", red, msg, reset)
}

func Green(msg string) {
	green := "\033[32m"
	fmt.Printf("%s%s%s\n", green, msg, reset)
}

func Yellow(msg string) {
	yellow := "\033[33m"
	fmt.Printf("%s%s%s\n", yellow, msg, reset)
}
