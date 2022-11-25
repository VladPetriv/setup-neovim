package main

import (
	"github.com/VladPetriv/setup-neovim/internal/app"
	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

func main() {
	validator := validation.New()
	inputter := input.New()

	service := service.New(&service.ServiceOptions{
		Validator: validator,
		Inputter:  inputter,
	})

	app.Run(service)
}
