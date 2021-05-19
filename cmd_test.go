package cmd_test

import (
	"errors"
	"flag"
	"testing"

	"github.com/hamba/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    [][2]string
		wantErr error
	}{
		{
			name:    "valid input",
			input:   []string{"a=b", "c=d"},
			want:    [][2]string{{"a", "b"}, {"c", "d"}},
			wantErr: nil,
		},
		{
			name:    "invalid input",
			input:   []string{"a"},
			wantErr: errors.New("string \"a\" does not contain separator"),
		},
		{
			name:    "mixed invalid input",
			input:   []string{"a=b", "c"},
			wantErr: errors.New("string \"c\" does not contain separator"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res, err := cmd.Split(test.input, "=")

			assert.Equal(t, res, test.want)
			assert.Equal(t, err, test.wantErr)
		})
	}
}

func newTestContext() (*cli.Context, *flag.FlagSet) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := cli.NewContext(&cli.App{}, fs, nil)
	return c, fs
}
