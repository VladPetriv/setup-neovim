package input

import "io"

type Inputter interface {
	GetInput(stdin io.Reader) (string, error)
	ProcessInputForPackageManagers(stdin io.Reader) (string, error)
}
