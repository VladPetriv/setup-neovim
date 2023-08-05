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

	services := service.Services{
		File:       service.NewFile(inputter, validator),
		Manager:    service.NewManager(inputter),
		Repository: service.NewRepository(inputter, validator),
	}

	app.Run(services)
}
