package observe

import (
	"io"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/statter/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

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
