package cmd_test

import (
	"context"
	"fmt"

	"github.com/hamba/cmd/v2"
	"github.com/urfave/cli/v2"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func ExampleNewLogger() {
	var c *cli.Context // Get this from your action

	log, err := cmd.NewLogger(c)
	if err != nil {
		// Handle error.
		return
	}

	_ = log
}

func ExampleNewStatter() {
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
	defer stats.Close()

	_ = stats
}

func ExampleNewProfiler() {
	var c *cli.Context // Get this from your action

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
	var c *cli.Context // Get this from your action

	log, err := cmd.NewLogger(c)
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
	defer tracer.Shutdown(context.Background())

	_ = tracer
}

func ExampleSplit() {
	input := []string{"a=b", "foo=bar"} // Usually from a cli.StringSlice

	tags, err := cmd.Split(input, "=")
	if err != nil {
		// Handle error
	}

	fmt.Println(tags)
	// Output: [[a b] [foo bar]]
}
