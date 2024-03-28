package observe_test

import (
	"context"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/cmd/v2/observe"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func ExampleNew() {
	var c *cli.Context // Get this from your action

	log, err := cmd.NewLogger(c)
	if err != nil {
		// Handle error.
		return
	}

	stats, err := cmd.NewStatter(c, log)
	if err != nil {
		// Handle error.
		return
	}

	prof, err := cmd.NewProfiler(c, "my-service", log)
	if err != nil {
		return
	}
	profStop := func() {}
	if prof != nil {
		profStop = func() { _ = prof.Stop() }
	}

	tracer, err := cmd.NewTracer(c, log,
		semconv.ServiceNameKey.String("my-service"),
		semconv.ServiceVersionKey.String("1.0.0"),
	)
	if err != nil {
		// Handle error.
		return
	}
	tracerCancel := func() { _ = tracer.Shutdown(context.Background()) }

	obsrv := observe.New(log, stats, tracer, tracerCancel, profStop)

	_ = obsrv
}

func ExampleNewFromCLI() {
	var c *cli.Context // Get this from your action.

	obsrv, err := observe.NewFromCLI(c, "my-service", &observe.Options{
		LogTimestamps: true,
		StatsRuntime:  true,
		TracingAttrs: []attribute.KeyValue{
			semconv.ServiceVersionKey.String("1.0.0"),
		},
	})
	if err != nil {
		// Handle error.
		return
	}

	_ = obsrv
}
