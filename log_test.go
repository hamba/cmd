package cmd_test

import (
	"testing"

	"github.com/hamba/cmd"
	"github.com/hamba/pkg/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		lvl     string
		format  string
		tags    *cli.StringSlice
		wantErr bool
	}{
		{
			name:    "Json Format",
			lvl:     "info",
			format:  "json",
			tags:    &cli.StringSlice{},
			wantErr: false,
		},
		{
			name:    "Logfmt Format",
			lvl:     "info",
			format:  "logfmt",
			tags:    &cli.StringSlice{},
			wantErr: false,
		},
		{
			name:    "No Format",
			lvl:     "",
			format:  "json",
			tags:    &cli.StringSlice{},
			wantErr: false,
		},
		{
			name:    "Invalid Format",
			lvl:     "info",
			format:  "invalid",
			tags:    &cli.StringSlice{},
			wantErr: false,
		},
		{
			name:    "Valid Level",
			lvl:     "info",
			format:  "",
			tags:    &cli.StringSlice{},
			wantErr: false,
		},
		{
			name:    "Invalid Level",
			lvl:     "invalid",
			format:  "json",
			tags:    &cli.StringSlice{},
			wantErr: true,
		},
		{
			name:    "Tags",
			lvl:     "info",
			format:  "json",
			tags:    &cli.StringSlice{"a=b"},
			wantErr: false,
		},
		{
			name:    "Invalid Tags",
			lvl:     "info",
			format:  "json",
			tags:    &cli.StringSlice{"single"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, fs := newTestContext()
			fs.String(cmd.FlagLogLevel, tt.lvl, "doc")
			fs.String(cmd.FlagLogFormat, tt.format, "doc")
			fs.Var(tt.tags, cmd.FlagLogTags, "doc")

			l, err := cmd.NewLogger(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Implements(t, (*log.Logger)(nil), l)
		})

	}
}
