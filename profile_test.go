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

func TestNewProfiler(t *testing.T) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	tests := []struct {
		name    string
		args    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "ignores no profiling",
			wantErr: assert.NoError,
		},
		{
			name: "starts profiler",
			args: []string{
				"--profiling.dsn=https://example.com?token=test&tenantid=me",
			},
			wantErr: assert.NoError,
		},
		{
			name: "supports tags",
			args: []string{
				"--profiling.dsn=https://example.com?token=test&tenantid=me",
				"--profiling.tags=cluster=test",
				"--profiling.tags=namespace=num",
			},
			wantErr: assert.NoError,
		},
		{
			name: "supports types",
			args: []string{
				"--profiling.dsn=https://example.com?token=test&tenantid=me",
				"--profiling.types=cpu",
			},
			wantErr: assert.NoError,
		},
		{
			name: "handles bad dsn",
			args: []string{
				"--profiling.dsn=:/",
			},
			wantErr: assert.Error,
		},
		{
			name: "handles basic and token auth",
			args: []string{
				"--profiling.dsn=https://test:test@example.com?token=test&tenantid=me",
			},
			wantErr: assert.Error,
		},
		{
			name: "handles bad tags",
			args: []string{
				"--profiling.dsn=https://example.com?token=test&tenantid=me",
				"--profiling.tags=cluster",
				"--profiling.tags=namespace=num",
			},
			wantErr: assert.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &cli.Command{
				Flags: cmd.ProfilingFlags,
				Action: func(_ context.Context, c *cli.Command) error {
					_, err := cmd.NewProfiler(c, "my-service", log)
					return err
				},
			}

			err := c.Run(t.Context(), append([]string{"test"}, test.args...))

			test.wantErr(t, err)
		})
	}
}
