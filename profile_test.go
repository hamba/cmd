package cmd_test

import (
	"io"
	"testing"
	"time"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/logger/v2"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestNewProfiler(t *testing.T) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	tests := []struct {
		name       string
		dsn        string
		uploadRate time.Duration
		tags       []string
		types      []string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:    "ignores no profiling",
			wantErr: assert.NoError,
		},
		{
			name:    "starts profiler",
			dsn:     "https://example.com?token=test&tenantid=me",
			wantErr: assert.NoError,
		},
		{
			name:    "supports tags",
			dsn:     "https://example.com?token=test&tenantid=me",
			tags:    []string{"cluster=test", "namespace=num"},
			wantErr: assert.NoError,
		},
		{
			name:    "supports types",
			dsn:     "https://example.com?token=test&tenantid=me",
			types:   []string{"cpu"},
			wantErr: assert.NoError,
		},
		{
			name:    "handles bad dsn",
			dsn:     ":/",
			wantErr: assert.Error,
		},
		{
			name:    "handles basic and token auth",
			dsn:     "https://test:test@example.com?token=test&tenantid=me",
			wantErr: assert.Error,
		},
		{
			name:    "handles bad tags",
			dsn:     "https://example.com?token=test&tenantid=me",
			tags:    []string{"cluster", "namespace=num"},
			wantErr: assert.Error,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			// t.Parallel()

			c, fs := newTestContext()
			fs.String(cmd.FlagProfilingDSN, test.dsn, "doc")
			fs.Duration(cmd.FlagProfileUploadRate, test.uploadRate, "doc")
			fs.Var(cli.NewStringSlice(test.tags...), cmd.FlagProfilingTags, "doc")
			fs.Var(cli.NewStringSlice(test.types...), cmd.FlagProfilingTypes, "doc")

			profiler, err := cmd.NewProfiler(c, "my-service", log)
			if profiler != nil {
				_ = profiler.Stop()
			}

			test.wantErr(t, err)
		})
	}
}
