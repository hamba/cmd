package cmd_test

import (
	"testing"

	"github.com/hamba/cmd"
	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestNewStats(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		prefix  string
		tags    *cli.StringSlice
		wantErr bool
	}{
		{
			name:    "No Stats",
			dsn:     "",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: false,
		},
		{
			name:    "Statsd",
			dsn:     "statsd://localhost:8125",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: false,
		},
		{
			name:    "L2met",
			dsn:     "l2met://",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: false,
		},
		{
			name:    "Prometheus",
			dsn:     "prometheus://",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: false,
		},
		{
			name:    "Prometheus With Server",
			dsn:     "prometheus://:51234",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: false,
		},
		{
			name:    "Unknown Stats",
			dsn:     "unknownscheme://",
			prefix:  "",
			tags:    cli.NewStringSlice(),
			wantErr: true,
		},
		{
			name:    "Invalid DSN",
			dsn:     "://",
			prefix:  "",
			tags:    cli.NewStringSlice(),
			wantErr: true,
		},
		{
			name:    "No Prefix",
			dsn:     "l2met://",
			prefix:  "",
			tags:    cli.NewStringSlice(),
			wantErr: false,
		},
		{
			name:    "Tags",
			dsn:     "l2met://",
			prefix:  "",
			tags:    cli.NewStringSlice("a=b"),
			wantErr: false,
		},
		{
			name:    "Invalid Tags",
			dsn:     "l2met://",
			prefix:  "",
			tags:    cli.NewStringSlice("a"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, fs := newTestContext()
			fs.String(cmd.FlagStatsDSN, tt.dsn, "doc")
			fs.String(cmd.FlagStatsPrefix, tt.prefix, "doc")
			fs.Var(tt.tags, cmd.FlagStatsTags, "doc")

			s, err := cmd.NewStats(c, log.Null)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Implements(t, (*stats.Statter)(nil), s)
		})
	}
}
