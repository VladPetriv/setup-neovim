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

func Run(srv service.Services) {
	// commandTimeout represents timeout that should be after completing the previous function.
	commandTimeout := 1 * time.Second

	utilErrors := srv.CheckUtilStatus()
	if len(utilErrors) >= 1 {
		for _, value := range utilErrors {
			colors.Red(value)
		}

		os.Exit(1)
	}

	colors.Green("All utilities are installed....")

	time.Sleep(commandTimeout)

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

	colors.Green("Successfully remove old nvim config...")

	time.Sleep(commandTimeout)

	url, err := srv.ProcessUserURL(os.Stdin)
	if err != nil {
		errs.WrapError("Validation for URL failed! Please try again... ", err)
	}

	colors.Green("URL is valid...")

	time.Sleep(commandTimeout)

	err = srv.CloneAndValidateRepository(url, os.Stdin)
	if err != nil {
		errs.WrapError("Failed to clone repository or repository didn't have base files for nvim configuration", err)
	}

	colors.Green("Repository successfully cloned and checked for base files")

	time.Sleep(commandTimeout)

	err = srv.ExtractAndMoveConfigDirectory("./nvim")
	if err != nil {
		errs.WrapError("Failed to extract and move config directory from repository", err)
	}
	colors.Green("Config successfully extracted!")

	time.Sleep(commandTimeout)

	packageManger, err := srv.ProcessPackageManagers(os.Stdin)
	if err != nil {
		errs.WrapError("Failed to install package manager. Please try again...", err)
	}

	if packageManger != "" {
		colors.Green(fmt.Sprintf("%s successfully installed", packageManger))

		time.Sleep(commandTimeout)
	}

	colors.Green("Nvim successfully configured!")
}
