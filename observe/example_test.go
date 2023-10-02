package observe_test

import (
	"context"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/cmd/v2/observe"
	"github.com/urfave/cli/v2"
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

	tracer, err := cmd.NewTracer(c, log,
		semconv.ServiceNameKey.String("my-service"),
		semconv.ServiceVersionKey.String("1.0.0"),
	)
	if err != nil {
		// Handle error.
		return
	}
	tracerCancel := func() { _ = tracer.Shutdown(context.Background()) }

	obsrv := observe.New(log, stats, tracer, tracerCancel)

	_ = obsrv
}
