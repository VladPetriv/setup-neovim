package app

import (
	"time"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/logger"
)

// _commandTimeout represents timeout that should be after completing the previous function
var _commandTimeout = 1 * time.Second

func Run(service service.Services, logger *logger.Logger) {
	errs := service.CheckUtilStatus()
	if len(errs) >= 1 {
		logger.Fatal().Interface("errors", errs).Msg("Check for console utilities failed! Please try again...")
	} else {
		logger.Info().Msg("All utilities are installed...")
	}

	time.Sleep(_commandTimeout)
	url, err := service.ProcessUserURL()
	if err != nil {
		logger.Fatal().Err(err).Msg("Process URL failed! Please try again...")
	}
	logger.Info().Msg("URL are valid...")

	time.Sleep(_commandTimeout)
	err = service.CloneAndValidateRepository(url)
	if err != nil {
		logger.Fatal().Err(err).Msg("Process repository failed! Please try again...")
	}
	logger.Info().Msg("Repository successfully cloned and validated...")

	time.Sleep(_commandTimeout)
	err = service.MoveConfigDirectory()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to move repository! Please try again...")
	}

	time.Sleep(_commandTimeout)
	packageManger, err := service.ProcessPackageManagers()
	if err != nil {
		logger.Fatal().Err(err).Msg("Process package managers failed! Please try again...")
	}

	time.Sleep(_commandTimeout)
	err = service.CompleteSetup(packageManger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Setup failed! Please try again...")
	}
	logger.Info().Msg("Editor successfully configured!")
}
