package input

type Inputter interface {
	GetInput() (string, error)
}
