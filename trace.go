package cmd

import (
	"context"
	"fmt"

	"github.com/ettle/strcase"
	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Tracing flag constants declared for CLI use.
const (
	FlagTracingExporter         = "tracing.exporter"
	FlagTracingEndpoint         = "tracing.endpoint"
	FlagTracingEndpointInsecure = "tracing.endpoint-insecure"
	FlagTracingTags             = "tracing.tags"
	FlagTracingHeaders          = "tracing.headers"
	FlagTracingRatio            = "tracing.ratio"
)

// CategoryTracing is the tracing flag category.
const CategoryTracing = "Tracing"

// TracingFlags are flags that configure tracing.
var TracingFlags = Flags{
	&cli.StringFlag{
		Name:     FlagTracingExporter,
		Category: CategoryTracing,
		Usage:    "The tracing backend. Supported: 'zipkin', 'otlphttp', 'otlpgrpc'.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagTracingExporter)),
	},
	&cli.StringFlag{
		Name:     FlagTracingEndpoint,
		Category: CategoryTracing,
		Usage:    "The tracing backend endpoint.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagTracingEndpoint)),
	},
	&cli.BoolFlag{
		Name:     FlagTracingEndpointInsecure,
		Category: CategoryTracing,
		Usage:    "Determines if the endpoint is insecure.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagTracingEndpointInsecure)),
	},
	&cli.StringMapFlag{
		Name:     FlagTracingTags,
		Category: CategoryTracing,
		Usage:    "A list of tags appended to every trace.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagTracingTags)),
	},
	&cli.StringMapFlag{
		Name:     FlagTracingHeaders,
		Category: CategoryTracing,
		Usage:    "A list of headers appended to every trace when supported by the exporter.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagTracingHeaders)),
	},
	&cli.FloatFlag{
		Name:     FlagTracingRatio,
		Category: CategoryTracing,
		Usage:    "The ratio between 0 and 1 of sample traces to take.",
		Value:    0.5,
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagTracingRatio)),
	},
}

// NewTracer returns a tracer configured from the cli.
func NewTracer(ctx context.Context, cmd *cli.Command, log *logger.Logger, attrs ...attribute.KeyValue) (*trace.TracerProvider, error) {
	otel.SetErrorHandler(logErrorHandler{log: log})

	exp, err := createExporter(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if exp == nil {
		return trace.NewTracerProvider(), nil
	}

	proc := trace.NewBatchSpanProcessor(exp)

	ratio := cmd.Float(FlagTracingRatio)
	sampler := trace.ParentBased(trace.TraceIDRatioBased(ratio))

	if tags := cmd.StringMap(FlagTracingTags); len(tags) > 0 {
		for k, v := range tags {
			attrs = append(attrs, attribute.String(k, v))
		}
	}

	return trace.NewTracerProvider(
		trace.WithSampler(sampler),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, attrs...)),
		trace.WithSpanProcessor(proc),
	), nil
}

func createExporter(ctx context.Context, cmd *cli.Command) (trace.SpanExporter, error) {
	backend := cmd.String(FlagTracingExporter)
	endpoint := cmd.String(FlagTracingEndpoint)

	switch backend {
	case "":
		return nil, nil //nolint:nilnil
	case "zipkin":
		return zipkin.New(endpoint)
	case "otlphttp":
		opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithHeaders(cmd.StringMap(FlagTracingHeaders))}
		if cmd.Bool(FlagTracingEndpointInsecure) {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		return otlptracehttp.New(ctx, opts...)
	case "otlpgrpc":
		opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithHeaders(cmd.StringMap(FlagTracingHeaders))}
		if cmd.Bool(FlagTracingEndpointInsecure) {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		return otlptracegrpc.New(ctx, opts...)
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
