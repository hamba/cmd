package cmd_test

import (
	"testing"

	"github.com/hamba/cmd"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		lvl     string
		format  string
		tags    *cli.StringSlice
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "json format",
			lvl:     "info",
			format:  "json",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "logfmt format",
			lvl:     "info",
			format:  "logfmt",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "console format",
			lvl:     "info",
			format:  "console",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "no format",
			lvl:     "",
			format:  "json",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "invalid format",
			lvl:     "info",
			format:  "invalid",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "valid Level",
			lvl:     "info",
			format:  "",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "invalid Level",
			lvl:     "invalid",
			format:  "json",
			tags:    cli.NewStringSlice(),
			wantErr: require.Error,
		},
		{
			name:    "tags",
			lvl:     "info",
			format:  "json",
			tags:    cli.NewStringSlice("a=b"),
			wantErr: require.NoError,
		},
		{
			name:    "invalid tags",
			lvl:     "info",
			format:  "json",
			tags:    cli.NewStringSlice("single"),
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, fs := newTestContext()
			fs.String(cmd.FlagLogLevel, tt.lvl, "doc")
			fs.String(cmd.FlagLogFormat, tt.format, "doc")
			fs.Var(tt.tags, cmd.FlagLogCtx, "doc")

			_, err := cmd.NewLogger(c)

			tt.wantErr(t, err)
		})

	}
}
