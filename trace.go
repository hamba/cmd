package cmd

import (
	"fmt"
	"net"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Tracing flag constants declared for CLI use.
const (
	FlagTracingExporter = "tracing.exporter"
	FlagTracingEndpoint = "tracing.endpoint"
	FlagTracingRatio    = "tracing.ratio"
)

// TracingFlags are flags that configure tracing.
var TracingFlags = Flags{
	&cli.StringFlag{
		Name:    FlagTracingExporter,
		Usage:   "The tracing backend. Supported: 'jaeger', 'zipkin'.",
		EnvVars: []string{"TRACING_EXPORTER"},
	},
	&cli.StringFlag{
		Name:    FlagTracingEndpoint,
		Usage:   "The tracing backend endpoint.",
		EnvVars: []string{"TRACING_ENDPOINT"},
	},
	&cli.Float64Flag{
		Name:    FlagTracingRatio,
		Usage:   "The ratio between 0 and 1 of sample traces to take.",
		Value:   0.5,
		EnvVars: []string{"TRACING_RATIO"},
	},
}

// NewTracer returns a tracer configures from the cli.
func NewTracer(c *cli.Context, log *logger.Logger, resAttributes ...attribute.KeyValue) (*trace.TracerProvider, error) {
	otel.SetErrorHandler(logErrorHandler{log: log})

	exp, err := createExporter(c)
	if err != nil {
		return nil, err
	}
	if exp == nil {
		return trace.NewTracerProvider(), nil
	}

	proc := trace.NewBatchSpanProcessor(exp)

	ratio := c.Float64(FlagTracingRatio)
	sampler := trace.ParentBased(trace.TraceIDRatioBased(ratio))

	return trace.NewTracerProvider(
		trace.WithSampler(sampler),
		trace.WithResource(resource.NewWithAttributes(resAttributes...)),
		trace.WithSpanProcessor(proc),
	), nil
}

func createExporter(c *cli.Context) (trace.SpanExporter, error) {
	backend := c.String(FlagTracingExporter)
	endpoint := c.String(FlagTracingEndpoint)

	switch backend {
	case "":
		return nil, nil
	case "jaeger":
		host, port, err := net.SplitHostPort(endpoint)
		if err != nil {
			return nil, err
		}

		return jaeger.NewRawExporter(
			jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost(host),
				jaeger.WithAgentPort(port),
			),
		)
	case "zipkin":
		return zipkin.NewRawExporter(endpoint)
	default:
		return nil, fmt.Errorf("unsupported tracing backend %q", backend)
	}
}

type logErrorHandler struct {
	log *logger.Logger
}

func (l logErrorHandler) Handle(err error) {
	if err == nil {
		return
	}
	l.log.Error(err.Error(), ctx.Str("component", "otel"))
}
