package input_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/stretchr/testify/assert"
)

func Test_GetInput(t *testing.T) {
	t.Parallel()

	testInput := input.New()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "GetInput success",
			input: "test",
			want:  "test",
		},
	}

	for _, tt := range tests {
		tt := tt //nolint
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", tt.input)))

			got, err := testInput.GetInput(&stdin)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, "", got)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
