package input

type Inputter interface {
	GetInput(msg string) (string, error)
}
