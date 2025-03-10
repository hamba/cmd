package observe

import (
	"context"
	"io"
	"time"

	otelpyroscope "github.com/grafana/otel-profiling-go"
	"github.com/hamba/cmd/v2"
	"github.com/hamba/logger/v2"
	lctx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/statter/v2"
	"github.com/hamba/statter/v2/runtime"
	"github.com/hamba/statter/v2/tags"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Options optionally configures an observer.
type Options struct {
	LogTimeFormat string
	LogTimestamps bool
	LogCtx        []logger.Field
	LogWriter     io.Writer

	StatsRuntime bool
	StatsTags    []statter.Tag

	TracingAttrs []attribute.KeyValue
}

// Observer contains observability primitives.
type Observer struct {
	Log       *logger.Logger
	Stats     *statter.Statter
	TraceProv trace.TracerProvider

	closeFns []func()
}

// New returns an observer with the given observability primitives.
func New(log *logger.Logger, stats *statter.Statter, traceProv trace.TracerProvider, closeFns ...func()) *Observer {
	return &Observer{
		Log:       log,
		Stats:     stats,
		TraceProv: traceProv,
		closeFns:  closeFns,
	}
}

// NewFromCLI returns an observer with the given observability primitives.
//
//nolint:cyclop // Splitting this will not make it simpler.
func NewFromCLI(cliCtx *cli.Context, svc string, opts *Options) (*Observer, error) {
	var closeFns []func()

	if opts == nil {
		opts = &Options{}
	}

	// Logger.
	log, err := cmd.NewLoggerWithOptions(cliCtx, &cmd.LoggerOptions{Writer: opts.LogWriter})
	if err != nil {
		return nil, err
	}
	if opts.LogTimeFormat != "" {
		logger.TimeFormat = opts.LogTimeFormat
	}
	if opts.LogTimestamps {
		closeFns = append(closeFns, log.WithTimestamp())
	}
	opts.LogCtx = append([]logger.Field{lctx.Str("svc", svc)}, opts.LogCtx...)
	log = log.With(opts.LogCtx...)

	// Statter.
	stats, err := cmd.NewStatter(cliCtx, log)
	if err != nil {
		for _, fn := range closeFns {
			fn()
		}
		return nil, err
	}
	closeFns = append(closeFns, func() { _ = stats.Close() })
	if opts.StatsRuntime {
		go runtime.Collect(stats)
	}
	opts.StatsTags = append([]statter.Tag{tags.Str("svc", svc)}, opts.StatsTags...)
	stats = stats.With("", opts.StatsTags...)

	// Profiler.
	prof, err := cmd.NewProfiler(cliCtx, svc, log)
	if err != nil {
		for _, fn := range closeFns {
			fn()
		}
		return nil, err
	}
	if prof != nil {
		closeFns = append(closeFns, func() { _ = prof.Stop() })
	}

	// Tracer.
	opts.TracingAttrs = append(opts.TracingAttrs, semconv.ServiceNameKey.String(svc))
	tracer, err := cmd.NewTracer(cliCtx, log, opts.TracingAttrs...)
	if err != nil {
		for _, fn := range closeFns {
			fn()
		}
		return nil, err
	}
	closeFns = append(closeFns, func() { _ = tracer.Shutdown(context.Background()) })

	var tp trace.TracerProvider = tracer
	if prof != nil && tracer != nil {
		tp = otelpyroscope.NewTracerProvider(tp)
	}

	return &Observer{
		Log:       log,
		Stats:     stats,
		TraceProv: tp,
		closeFns:  closeFns,
	}, nil
}

// Tracer returns a tracer with the given name and options.
// If no trace provider has been set, this function will panic.
func (o *Observer) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	if o.TraceProv == nil {
		panic("calling tracer when no trace provider has been set")
	}
	return o.TraceProv.Tracer(name, opts...)
}

// Close closes the observability primitives.
func (o *Observer) Close() {
	for _, cancel := range o.closeFns {
		cancel()
	}
}

// NewFake returns a fake observer that reports nothing.
// This is useful for tests.
func NewFake() *Observer {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)
	stats := statter.New(statter.DiscardReporter, time.Minute)
	tracer := otel.GetTracerProvider()

	return New(log, stats, tracer)
}
