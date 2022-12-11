package service

import (
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

type service struct {
	validator validation.Validator
	input     input.Inputter
}

var _ Services = (*service)(nil)

type Options struct {
	Validator validation.Validator
	Inputter  input.Inputter
}

func New(options *Options) Services {
	return &service{
		validator: options.Validator,
		input:     options.Inputter,
	}
}

func (s service) CheckUtilStatus() map[string]string {
	return s.validator.ValidateConsoleUtilities()
}
