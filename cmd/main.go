package main

import (
	"github.com/VladPetriv/setup-neovim/internal/app"
	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/logger"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

func main() {
	logger := logger.Get()

	validator := validation.New(logger)
	inputter := input.New(logger)

	service := service.New(&service.ServiceOptions{
		Validator: validator,
		Inputter:  inputter,
	})

	app.Run(service, logger)
}
