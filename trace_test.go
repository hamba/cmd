package cmd_test

import (
	"io"
	"testing"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/logger/v2"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func TestNewTracer(t *testing.T) {
	tests := []struct {
		name     string
		exporter string
		endpoint string
		tags     []string
		ratio    float64
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name:     "no tracer",
			exporter: "",
			endpoint: "",
			ratio:    1.0,
			wantErr:  require.NoError,
		},
		{
			name:     "zipkin",
			exporter: "zipkin",
			endpoint: "http://localhost:1234/api/v2",
			ratio:    1.0,
			wantErr:  require.NoError,
		},
		{
			name:     "with tags",
			exporter: "otlphttp",
			endpoint: "localhost:1234",
			tags:     []string{"cluster=test", "namespace=num"},
			ratio:    1.0,
			wantErr:  require.NoError,
		},
		{
			name:     "unknown exporter",
			exporter: "some-exporter",
			endpoint: "localhost:1234",
			ratio:    1.0,
			wantErr:  require.Error,
		},
		{
			name:     "ratio too low",
			exporter: "otlpgrpc",
			endpoint: "localhost:1234",
			ratio:    -1.0,
			wantErr:  require.NoError,
		},
		{
			name:     "ratio too high",
			exporter: "otlpgrpc",
			endpoint: "localhost:1234",
			ratio:    2.0,
			wantErr:  require.NoError,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c, fs := newTestContext()
			fs.String(cmd.FlagTracingExporter, test.exporter, "doc")
			fs.String(cmd.FlagTracingEndpoint, test.endpoint, "doc")
			fs.Var(cli.NewStringSlice(test.tags...), cmd.FlagTracingTags, "doc")
			fs.Float64(cmd.FlagTracingRatio, test.ratio, "doc")

			log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

			_, err := cmd.NewTracer(c, log,
				semconv.ServiceNameKey.String("my-service"),
				semconv.ServiceVersionKey.String("1.0.0"),
			)

			test.wantErr(t, err)
		})
	}
}
