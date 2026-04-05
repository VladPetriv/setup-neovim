package main

import (
	"github.com/VladPetriv/setup-neovim/internal/app"
	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

func main() {
	validator := validation.New()

	services := service.Services{
		File:       service.NewFile(validator),
		Manager:    service.NewManager(),
		Repository: service.NewRepository(validator),
	}

	app.Run(services)
}
