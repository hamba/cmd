package cmd_test

import (
	"io"
	"testing"
	"time"

	"github.com/hamba/cmd"
	"github.com/hamba/logger/v2"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestNewStats(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		prefix  string
		tags    *cli.StringSlice
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "no stats",
			dsn:     "",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "statsd",
			dsn:     "statsd://localhost:8125?flushBytes=1423&flushInterval=10s",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "l2met",
			dsn:     "l2met://",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "prometheus",
			dsn:     "prometheus://",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "prometheus with server",
			dsn:     "prometheus://:51234",
			prefix:  "test",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "unknown stats scheme",
			dsn:     "unknownscheme://",
			prefix:  "",
			tags:    cli.NewStringSlice(),
			wantErr: require.Error,
		},
		{
			name:    "invalid DSN",
			dsn:     "://",
			prefix:  "",
			tags:    cli.NewStringSlice(),
			wantErr: require.Error,
		},
		{
			name:    "no prefix",
			dsn:     "l2met://",
			prefix:  "",
			tags:    cli.NewStringSlice(),
			wantErr: require.NoError,
		},
		{
			name:    "tags",
			dsn:     "l2met://",
			prefix:  "",
			tags:    cli.NewStringSlice("a=b"),
			wantErr: require.NoError,
		},
		{
			name:    "invalid tags",
			dsn:     "l2met://",
			prefix:  "",
			tags:    cli.NewStringSlice("a"),
			wantErr: require.Error,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c, fs := newTestContext()
			fs.String(cmd.FlagStatsDSN, test.dsn, "doc")
			fs.Duration(cmd.FlagStatsInterval, time.Second, "doc")
			fs.String(cmd.FlagStatsPrefix, test.prefix, "doc")
			fs.Var(test.tags, cmd.FlagStatsTags, "doc")

			log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

			_, err := cmd.NewStatter(c, log)

			test.wantErr(t, err)
		})
	}
}
