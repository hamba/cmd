package cmd_test

import (
	"context"
	"io"
	"testing"

	"github.com/hamba/cmd/v3"
	"github.com/hamba/logger/v2"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestNewStats(t *testing.T) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	tests := []struct {
		name    string
		args    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "no stats",
			args:    []string{},
			wantErr: assert.NoError,
		},
		{
			name:    "statsd",
			args:    []string{"--stats.dsn=statsd://localhost:8125?flushBytes=1423&flushInterval=10s", "--stats.prefix=test"},
			wantErr: assert.NoError,
		},
		{
			name:    "l2met",
			args:    []string{"--stats.dsn=l2met://", "--stats.prefix=test"},
			wantErr: assert.NoError,
		},
		{
			name:    "prometheus",
			args:    []string{"--stats.dsn=prometheus://", "--stats.prefix=test"},
			wantErr: assert.NoError,
		},
		{
			name:    "prometheus with server",
			args:    []string{"--stats.dsn=prom://127.0.0.1:51234", "--stats.prefix=test"},
			wantErr: assert.NoError,
		},
		{
			name:    "victoria metrics",
			args:    []string{"--stats.dsn=victoriametrics://", "--stats.prefix=test"},
			wantErr: assert.NoError,
		},
		{
			name:    "victoria metrics with server",
			args:    []string{"--stats.dsn=vm://127.0.0.1:51234", "--stats.prefix=test"},
			wantErr: assert.NoError,
		},
		{
			name:    "unknown stats scheme",
			args:    []string{"--stats.dsn=unknownscheme://"},
			wantErr: assert.Error,
		},
		{
			name:    "invalid DSN",
			args:    []string{"--stats.dsn=://"},
			wantErr: assert.Error,
		},
		{
			name:    "no prefix",
			args:    []string{"--stats.dsn=l2met://"},
			wantErr: assert.NoError,
		},
		{
			name:    "tags",
			args:    []string{"--stats.dsn=l2met://", "--stats.tags=a=b"},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid tags",
			args:    []string{"--stats.dsn=l2met://", "--stats.tags=a"},
			wantErr: assert.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &cli.Command{
				Flags: cmd.StatsFlags,
				Action: func(_ context.Context, c *cli.Command) error {
					_, err := cmd.NewStatter(c, log)
					return err
				},
			}

			err := c.Run(t.Context(), append([]string{"test"}, test.args...))

			test.wantErr(t, err)
		})
	}
}
