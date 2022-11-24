package service

import (
	"fmt"
	"os/exec"

	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

type service struct {
	validator validation.Validator
	input     input.Inputter
}

var _ Services = (*service)(nil)

type ServiceOptions struct {
	Validator validation.Validator
	Inputter  input.Inputter
}

func New(options *ServiceOptions) *service {
	return &service{
		validator: options.Validator,
		input:     options.Inputter,
	}
}

func (s service) CheckUtilStatus() map[string]string {
	return s.validator.ValidateConsoleUtilities()
}

func (s service) CompleteSetup() error {
	cmd := exec.Command("nvim", "+PackerSync")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to complete setup: %w", err)
	}

	return nil
}