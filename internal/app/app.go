package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/colors"
	"github.com/VladPetriv/setup-neovim/pkg/errs"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/charmbracelet/huh"
)

func Run(srv service.Services) {
	validator := validation.New()

	utilErrors := srv.File.CheckUtilStatus()
	if len(utilErrors) >= 1 {
		for _, msg := range utilErrors {
			colors.Red(msg)
		}

		os.Exit(0)
	}

	colors.Green("All utilities are installed....")

	handleExistingConfig(srv)

	url, sshKeyPath := promptForRepo(validator)
	cloneAndInstallConfig(srv, url, sshKeyPath)
	handlePackageManagerInstall(srv)
}

func handleExistingConfig(srv service.Services) {
	exists, err := srv.File.CheckConfigExists()
	if err != nil {
		errs.WrapError("Failed to check for existing nvim config!", err)
	}

	if !exists {
		return
	}

	var confirm bool
	err = huh.NewConfirm().
		Title("Existing nvim config detected. Remove it to continue?").
		Value(&confirm).
		Run()
	handlePromptErr(err)

	if !confirm {
		colors.Green("Thank you for using setup-nvim!")
		os.Exit(0)
	}

	if err = srv.File.DeleteConfig(); err != nil {
		errs.WrapError("Failed to delete existing nvim config!", err)
	}

	colors.Green("Successfully removed old nvim config...")
}

func promptForRepo(validator validation.Validator) (string, string) {
	var url string
	err := huh.NewInput().
		Title("Enter URL to your nvim config").
		Validate(validator.ValidateURL).
		Value(&url).
		Run()
	handlePromptErr(err)

	colors.Green("URL is valid...")

	var sshKeyPath string
	if strings.Contains(url, "git@") {
		err = huh.NewInput().
			Title("Enter path to your SSH key (e.g. .ssh/id_ed25519)").
			Value(&sshKeyPath).
			Run()
		handlePromptErr(err)
	}

	return url, sshKeyPath
}

func cloneAndInstallConfig(srv service.Services, url, sshKeyPath string) {
	if err := srv.Repository.CloneRepository(url, sshKeyPath); err != nil {
		errs.WrapError("Failed to clone repository!", err)
	}

	repositoryPath := fmt.Sprintf("./%s", service.DirectoryNameForClonnedRepository)
	if err := srv.Repository.ValidateRepository(repositoryPath); err != nil {
		if errors.Is(err, validation.ErrNoBaseFilesInRepository) {
			errs.WrapError("Repository doesn't have base files for nvim configuration!", err)
		}

		errs.WrapError("Failed to validate repository!", err)
	}

	colors.Green("Repository successfully cloned and validated...")

	if err := srv.File.ExtractAndMoveConfigDirectory(repositoryPath); err != nil {
		if deleteErr := srv.File.DeleteRepositoryDirectory(repositoryPath); deleteErr != nil {
			colors.Yellow("Could not delete cloned repository")
		}

		errs.WrapError("Failed to extract and move config directory!", err)
	}

	if deleteErr := srv.File.DeleteRepositoryDirectory(repositoryPath); deleteErr != nil {
		colors.Yellow("Could not delete cloned repository")
	}

	colors.Green("Config successfully extracted and moved to '.config' directory!")
}

func handlePackageManagerInstall(srv service.Services) {
	detected, err := srv.Manager.DetectInstalledPackageManagers()
	if err != nil {
		errs.WrapError("Failed to detect installed package managers!", err)
	}

	if len(detected) > 0 {
		colors.Yellow(fmt.Sprintf("Detected installed package managers: %s", strings.Join(detected, ", ")))

		var removeExisting bool
		err = huh.NewConfirm().
			Title("Remove existing package managers and install new?").
			Value(&removeExisting).
			Run()
		handlePromptErr(err)

		if !removeExisting {
			colors.Green("Nvim successfully configured!")
			os.Exit(0)
		}

		if err = srv.Manager.DeletePackageManagers(); err != nil {
			errs.WrapError("Failed to delete existing package managers!", err)
		}
	}

	var pm string
	err = huh.NewSelect[string]().
		Title("Choose a package manager to install").
		Options(
			huh.NewOption("Packer", service.PackerPluginManager),
			huh.NewOption("Vim Plug", service.VimPlugPluginManager),
			huh.NewOption("none", service.NonePluginManager),
		).
		Value(&pm).
		Run()
	handlePromptErr(err)

	if err = srv.Manager.InstallPackageManager(pm); err != nil {
		errs.WrapError(fmt.Sprintf("Failed to install %s!", pm), err)
	}

	if pm != service.NonePluginManager {
		colors.Green(fmt.Sprintf("%s successfully installed", pm))
	}

	colors.Green("Nvim successfully configured!")
}

func handlePromptErr(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, huh.ErrUserAborted) {
		colors.Yellow("Installation aborted.")
		os.Exit(0)
	}

	errs.WrapError("Prompt failed!", err)
}
