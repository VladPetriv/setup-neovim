package service_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func Test_ProcesUserURL(t *testing.T) {
	t.Parallel()

	testService := service.New(&service.Options{
		Inputter:  input.New(),
		Validator: validation.New(),
	})

	tests := []struct {
		name        string
		input       string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "ProcessUserURL success with valid host[github]",
			input:       "git@github.com:VladPetriv/setup-neovim.git",
			expectedURL: "git@github.com:VladPetriv/setup-neovim.git",
			expectedErr: nil,
		},
		{
			name:        "ProcessUserURL success with valid host[gitlab]",
			input:       "git@gitlab.com:gitlab-org/gitaly.git",
			expectedURL: "git@gitlab.com:gitlab-org/gitaly.git",
			expectedErr: nil,
		},
		{
			name:        "ProcessUserURL fail with invalid host",
			input:       "test",
			expectedURL: "",
			expectedErr: validation.ErrURLContainsInvalidHost,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", tt.input)))

			url, err := testService.ProcessUserURL(&stdin)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedURL, url)
		})
	}
}
