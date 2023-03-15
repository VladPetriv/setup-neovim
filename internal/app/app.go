package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/colors"
	"github.com/VladPetriv/setup-neovim/pkg/errs"
)

// commandTimeout represents timeout that should be after completing the previous function.
const commandTimeout = 1 * time.Second

func Run(srv service.Services) {
	utilErrors := srv.CheckUtilStatus()
	if len(utilErrors) >= 1 {
		for _, value := range utilErrors {
			colors.Red(value)
		}

		os.Exit(1)
	}

	successfulCommand("All utilities are installed....")

	err := srv.DeleteConfigOrStopInstallation(os.Stdin)
	if err != nil {
		if errors.Is(err, service.ErrStopInstallation) {
			errs.WrapError("Thank you for using setup-nvim!", err)
		}

		if errors.Is(err, service.ErrEnterValidAnswer) {
			errs.WrapError("Please choose correct answer for question!", err)
		}

		errs.WrapError("Failed to delete config or stop installation!", err)
	}

	successfulCommand("Successfully remove old nvim config...")

	url, err := srv.ProcessUserURL(os.Stdin)
	if err != nil {
		errs.WrapError("Validation for URL failed! Please try again... ", err)
	}

	successfulCommand("URL is valid...")

	err = srv.CloneAndValidateRepository(url, os.Stdin)
	if err != nil {
		errs.WrapError("Failed to clone repository or repository didn't have base files for nvim configuration", err)
	}

	successfulCommand("Repository successfully cloned and checked for base files")

	err = srv.ExtractAndMoveConfigDirectory("./nvim")
	if err != nil {
		errs.WrapError("Failed to extract and move config directory from repository", err)
	}

	successfulCommand("Config successfully extracted!")

	packageManger, err := srv.ProcessPackageManagers(os.Stdin)
	if err != nil {
		errs.WrapError("Failed to install package manager. Please try again...", err)
	}

	if packageManger != "" {
		successfulCommand(fmt.Sprintf("%s successfully installed", packageManger))
	}

	colors.Green("Nvim successfully configured!")
}

// successfulCommand print message with green color and wait command timeout [3s].
func successfulCommand(text string) {
	colors.Green(text)

	time.Sleep(commandTimeout)
}
