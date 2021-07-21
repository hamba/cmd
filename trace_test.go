package cmd_test

import (
	"io"
	"testing"

	"github.com/hamba/cmd"
	"github.com/hamba/logger/v2"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/semconv"
)

func TestNewTracer(t *testing.T) {
	tests := []struct {
		name     string
		exporter string
		endpoint string
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
			name:    "jaeger",
			exporter: "jaeger",
			endpoint: "localhost:1234",
			ratio:    1.0,
			wantErr: require.NoError,
		},
		{
			name:    "jaeger invalid endpoint",
			exporter: "jaeger",
			endpoint: "localhost",
			ratio:    1.0,
			wantErr: require.Error,
		},
		{
			name:    "zipkin",
			exporter: "zipkin",
			endpoint: "http://localhost:1234/api/v2",
			ratio:    1.0,
			wantErr: require.NoError,
		},
		{
			name:    "unknown exporter",
			exporter: "some-exporter",
			endpoint: "localhost:1234",
			ratio:    1.0,
			wantErr: require.Error,
		},
		{
			name:    "ratio too low",
			exporter: "jaeger",
			endpoint: "localhost:1234",
			ratio:    -1.0,
			wantErr: require.NoError,
		},
		{
			name:    "ratio too high",
			exporter: "jaeger",
			endpoint: "localhost:1234",
			ratio:    2.0,
			wantErr: require.NoError,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c, fs := newTestContext()
			fs.String(cmd.FlagTracingExporter, test.exporter, "doc")
			fs.String(cmd.FlagTracingEndpoint, test.endpoint, "doc")
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
