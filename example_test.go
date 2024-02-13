package cmd_test

import (
	"context"

	"github.com/hamba/cmd/v3"
	"github.com/urfave/cli/v3"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func ExampleNewLogger() {
	var c *cli.Command // Get this from your action

	log, err := cmd.NewLogger(c)
	if err != nil {
		// Handle error.
		return
	}

	_ = log
}

func ExampleNewStatter() {
	var c *cli.Command // Get this from your action

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
	defer func() { _ = stats.Close() }()

	_ = stats
}

func ExampleNewProfiler() {
	var c *cli.Command // Get this from your action

	log, err := cmd.NewLogger(c)
	if err != nil {
		// Handle error.
		return
	}

	prof, err := cmd.NewProfiler(c, "my-service", log)
	if err != nil {
		// Handle error.
		return
	}
	if prof != nil {
		defer func() { _ = prof.Stop() }()
	}

	_ = prof
}

func ExampleNewTracer() {
	var (
		ctx context.Context
		c   *cli.Command // Get this from your action
	)

	log, err := cmd.NewLogger(c)
	if err != nil {
		// Handle error.
		return
	}

	tracer, err := cmd.NewTracer(ctx, c, log,
		semconv.ServiceNameKey.String("my-service"),
		semconv.ServiceVersionKey.String("1.0.0"),
	)
	if err != nil {
		// Handle error.
		return
	}
	defer func() { _ = tracer.Shutdown(context.Background()) }()

	_ = tracer
}
