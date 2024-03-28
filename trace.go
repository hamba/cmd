package cmd

import (
	"fmt"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Tracing flag constants declared for CLI use.
const (
	FlagTracingExporter         = "tracing.exporter"
	FlagTracingEndpoint         = "tracing.endpoint"
	FlagTracingEndpointInsecure = "tracing.endpoint-insecure"
	FlagTracingTags             = "tracing.tags"
	FlagTracingRatio            = "tracing.ratio"
)

// TracingFlags are flags that configure tracing.
var TracingFlags = Flags{
	&cli.StringFlag{
		Name:    FlagTracingExporter,
		Usage:   "The tracing backend. Supported: 'zipkin', 'otlphttp', 'otlpgrpc'.",
		EnvVars: []string{"TRACING_EXPORTER"},
	},
	&cli.StringFlag{
		Name:    FlagTracingEndpoint,
		Usage:   "The tracing backend endpoint.",
		EnvVars: []string{"TRACING_ENDPOINT"},
	},
	&cli.BoolFlag{
		Name:    FlagTracingEndpointInsecure,
		Usage:   "Determines if the endpoint is insecure.",
		EnvVars: []string{"TRACING_ENDPOINT_INSECURE"},
	},
	&cli.StringSliceFlag{
		Name:    FlagTracingTags,
		Usage:   "A list of tags appended to every trace. Format: key=value.",
		EnvVars: []string{"TRACING_TAGS"},
	},
	&cli.Float64Flag{
		Name:    FlagTracingRatio,
		Usage:   "The ratio between 0 and 1 of sample traces to take.",
		Value:   0.5,
		EnvVars: []string{"TRACING_RATIO"},
	},
}

// NewTracer returns a tracer configured from the cli.
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

	attrs := resAttributes
	if tags := c.StringSlice(FlagTracingTags); len(tags) > 0 {
		strTags, err := Split(tags, "=")
		if err != nil {
			return nil, err
		}
		for _, kv := range strTags {
			attrs = append(attrs, attribute.String(kv[0], kv[1]))
		}
	}

	return trace.NewTracerProvider(
		trace.WithSampler(sampler),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, attrs...)),
		trace.WithSpanProcessor(proc),
	), nil
}

func createExporter(c *cli.Context) (trace.SpanExporter, error) {
	backend := c.String(FlagTracingExporter)
	endpoint := c.String(FlagTracingEndpoint)

	switch backend {
	case "":
		return nil, nil
	case "zipkin":
		return zipkin.New(endpoint)
	case "otlphttp":
		opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
		if c.Bool(FlagTracingEndpointInsecure) {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		return otlptracehttp.New(c.Context, opts...)
	case "otlpgrpc":
		opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint)}
		if c.Bool(FlagTracingEndpointInsecure) {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		return otlptracegrpc.New(c.Context, opts...)
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
