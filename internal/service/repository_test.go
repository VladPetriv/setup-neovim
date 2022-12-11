package service_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func Test_CloneAndValidateRepository(t *testing.T) { //nolint:tparallel // t.Parallel() causes conflicts with go-git
	t.Parallel()

	testService := service.New(&service.Options{
		Inputter:  input.New(),
		Validator: validation.New(),
	})

	tests := []struct {
		name          string
		inputURL      string
		inputFilePath string
		wantErr       bool
	}{
		{
			name:          "CloneAndValidateRepository success with https URL [github]",
			inputURL:      "https://github.com/jdhao/nvim-config.git",
			inputFilePath: "",
		},
		{
			name:          "CloneAndValidateRepository success with https URL [gitlab]",
			inputURL:      "https://gitlab.com/VladPetriv/nvim-config.git",
			inputFilePath: "",
		},

		{
			name:          "CloneAndValidateRepository success with SSH URL",
			inputURL:      "git@github.com:VladPetriv/nvim-config.git",
			inputFilePath: ".ssh/id_ed25519",
		},
		{
			name:          "CloneAndValidateRepository fail with file validation",
			inputURL:      "https://gitlab.com/VladPetriv/tg_scanner.git",
			inputFilePath: "",
			wantErr:       true,
		},
		{
			name:          "CloneAndValidateRepository fail with error when create public key from file",
			inputURL:      "git@github.com:VladPetriv/nvim-config.git",
			inputFilePath: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", tt.inputFilePath)))

			err := testService.CloneAndValidateRepository(tt.inputURL, &stdin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// remove only when repository cloned and validated
			err = os.RemoveAll("./nvim")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
