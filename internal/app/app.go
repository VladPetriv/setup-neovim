package app

import (
	"fmt"
	"os"
	"time"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/colors"
	"github.com/VladPetriv/setup-neovim/pkg/errs"
)

func Run(service service.Services) {
	// commandTimeout represents timeout that should be after completing the previous function.
	commandTimeout := 1 * time.Second

	errors := service.CheckUtilStatus()
	if len(errors) >= 1 {
		for _, value := range errors {
			colors.Red(value)
		}

		os.Exit(1)
	}

	colors.Green("All utilities are installed....")

	time.Sleep(commandTimeout)

	url, err := service.ProcessUserURL(os.Stdin)
	if err != nil {
		errs.WrapError("Validation for URL failed! Please try again... ", err)
	}

	colors.Green("URL is valid...")

	time.Sleep(commandTimeout)

	err = service.CloneAndValidateRepository(url, os.Stdin)
	if err != nil {
		errs.WrapError("Failed to clone repository or repository didn't have base files for nvim configuration", err)
	}

	colors.Green("Repository successfully cloned and checked for base files")

	time.Sleep(commandTimeout)

	err = service.ExtractAndMoveConfigDirectory("./nvim")
	if err != nil {
		errs.WrapError("Failed to extract and move config directory from repository", err)
	}
	colors.Green("Config successfully extracted!")

	time.Sleep(commandTimeout)

	packageManger, err := service.ProcessPackageManagers(os.Stdin)
	if err != nil {
		errs.WrapError("Failed to install package manager. Please try again...", err)
	}

	if packageManger != "" {
		colors.Green(fmt.Sprintf("%s successfully installed", packageManger))

		time.Sleep(commandTimeout)
	}

	colors.Green("Nvim successfully configured!")
}
