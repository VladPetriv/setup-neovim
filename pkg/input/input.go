package input

import (
	"fmt"

	"github.com/VladPetriv/setup-neovim/pkg/logger"
)

type input struct {
	log *logger.Logger
}

var _ Inputter = (*input)(nil)

func New(log *logger.Logger) *input {
	return &input{
		log: log,
	}
}

func (i input) GetInput(msg string) (string, error) {
	log := i.log

	log.Info().Msgf("%s: ", msg)

	var data string
	_, err := fmt.Scanln(&data)
	if err != nil {
		return "", fmt.Errorf("get input error: %w", err)
	}

	return data, nil
}
