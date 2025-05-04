package cmd_test

import (
	"context"
	"testing"

	"github.com/hamba/cmd/v3"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "json format",
			args:    []string{"--log.level=info", "--log.format=json"},
			wantErr: assert.NoError,
		},
		{
			name:    "logfmt format",
			args:    []string{"--log.level=info", "--log.format=logfmt"},
			wantErr: assert.NoError,
		},
		{
			name:    "console format",
			args:    []string{"--log.level=info", "--log.format=console"},
			wantErr: assert.NoError,
		},
		{
			name:    "no format",
			args:    []string{"--log.level=info"},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid format",
			args:    []string{"--log.level=info", "--log.format=invalid"},
			wantErr: assert.NoError,
		},
		{
			name:    "valid level",
			args:    []string{"--log.level=info"},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid level",
			args:    []string{"--log.level=invalid", "--log.format=json"},
			wantErr: assert.Error,
		},
		{
			name:    "tags",
			args:    []string{"--log.level=info", "--log.format=json", "--log.ctx=a=b"},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid tags",
			args:    []string{"--log.level=info", "--log.format=json", "--log.ctx=a"},
			wantErr: assert.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &cli.Command{
				Flags: cmd.LogFlags,
				Action: func(_ context.Context, c *cli.Command) error {
					_, err := cmd.NewLogger(c)
					return err
				},
			}

			err := c.Run(t.Context(), append([]string{"test"}, test.args...))

			test.wantErr(t, err)
		})
	}
}
