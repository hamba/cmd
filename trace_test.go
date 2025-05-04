package cmd_test

import (
	"context"
	"io"
	"testing"

	"github.com/hamba/cmd/v3"
	"github.com/hamba/logger/v2"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func TestNewTracer(t *testing.T) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	tests := []struct {
		name    string
		args    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "no tracer",
			args:    []string{},
			wantErr: assert.NoError,
		},
		{
			name: "zipkin",
			args: []string{
				"--tracing.exporter=zipkin",
				"--tracing.endpoint=http://localhost:1234/api/v2",
			},
			wantErr: assert.NoError,
		},
		{
			name: "otelhttp",
			args: []string{
				"--tracing.exporter=otlphttp",
				"--tracing.endpoint=http://localhost:1234/",
			},
			wantErr: assert.NoError,
		},
		{
			name: "otelgrpc",
			args: []string{
				"--tracing.exporter=otlpgrpc",
				"--tracing.endpoint=localhost:1234",
			},
			wantErr: assert.NoError,
		},
		{
			name: "with tags",
			args: []string{
				"--tracing.exporter=zipkin",
				"--tracing.endpoint=http://localhost:1234/api/v2",
				"--tracing.ratio=1",
				"--tracing.tags=cluster=test",
				"--tracing.tags=namespace=num",
			},
			wantErr: assert.NoError,
		},
		{
			name: "with headers",
			args: []string{
				"--tracing.exporter=otlphttp",
				"--tracing.endpoint=http://localhost:1234/",
				"--tracing.ratio=1",
				"--tracing.headers=cluster=test",
				"--tracing.headers=namespace=num",
			},
			wantErr: assert.NoError,
		},
		{
			name: "unknown exporter",
			args: []string{
				"--tracing.exporter=some-exporter",
				"--tracing.endpoint=localhost:1234",
			},
			wantErr: assert.Error,
		},
		{
			name: "ratio too low",
			args: []string{
				"--tracing.exporter=otlpgrpc",
				"--tracing.endpoint=localhost:1234",
				"--tracing.ratio=-1",
			},
			wantErr: assert.NoError,
		},
		{
			name: "ratio too high",
			args: []string{
				"--tracing.exporter=otlpgrpc",
				"--tracing.endpoint=localhost:1234",
				"--tracing.ratio=2",
			},
			wantErr: assert.NoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &cli.Command{
				Flags: cmd.TracingFlags,
				Action: func(ctx context.Context, c *cli.Command) error {
					_, err := cmd.NewTracer(ctx, c, log,
						semconv.ServiceNameKey.String("my-service"),
						semconv.ServiceVersionKey.String("1.0.0"),
					)
					return err
				},
			}

			err := c.Run(t.Context(), append([]string{"test"}, test.args...))

			test.wantErr(t, err)
		})
	}
}
