![Logo](http://svg.wiersma.co.za/hamba/project?title=cmd&tag=Go%20cmd%20helper)

[![Go Report Card](https://goreportcard.com/badge/github.com/hamba/cmd)](https://goreportcard.com/report/github.com/hamba/cmd)
[![Build Status](https://github.com/hamba/cmd/actions/workflows/test.yml/badge.svg)](https://github.com/hamba/cmd/actions)
[![Coverage Status](https://coveralls.io/repos/github/hamba/cmd/badge.svg?branch=master)](https://coveralls.io/github/hamba/cmd?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamba/cmd/v2.svg)](https://pkg.go.dev/github.com/hamba/cmd/v2)
[![GitHub release](https://img.shields.io/github/release/hamba/cmd.svg)](https://github.com/hamba/cmd/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hamba/cmd/master/LICENSE)

Go cmd helper. 

This provides helpers on top of `github.com/urfave/cli`.

## Overview

Install with:

```shell
go get github.com/hamba/cmd/v2
```

## Example

```go
func yourAction(c *cli.Context) error {
    log, err := cmd.NewLogger(c)
	if err != nil {
		// Handle error.
	}

	stats, err := cmd.NewStatter(c, log)
	if err != nil {
		// Handle error.
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

    // Run your application here...
	
	return nil
}
```

## Flags

### Logger

The logger flags are used by `cmd.NewLogger` to create a `hamba.Logger`.

#### FlagLogFormat: *--log.format, $LOG_FORMAT*

This flag sets the log formatter to use. The available options are `logfmt` *(default)*, `json`, `console`.

Example: `--log.format=console`

#### FlagLogLevel: *--log.level, $LOG_LEVEL*

This flag sets the log level to filer on. The available options are `debug`, `info` (default), `warn`, `error`, `crit`.

Example: `--log.level=error`

#### FlagLogCtx: *--log.ctx, $LOG_CTX*

This flag sets contextual key value pairs to set on all log messages. This flag can be specified multiple times.

Example: `--log.ctx="app=my-app" --log.ctx="zone=eu-west"`

### Statter

The statter flags are used by `cmd.NewStatter` to create a new `hamba.Statter.

#### FlagStatsDSN: *--stats.dsn, $STATS_DSN*

This flag sets the DSN describing the stats reporter to use. The available options are `statsd`, `prometheus`, `l2met`, `victoriametrics`.

The DSN can in some situations specify the host and configuration values as shown in the below examples:

**Statsd:** 

`--stats.dsn="statsd://host:port?flushBytes=1432&flushInterval=10s"`

The `host` and `port` are required. Optionally `flushBytes` and `flushInterval` can be set, controlling how often the stats will
be sent to the Statsd server.

**Prometheus:**

`--stats.dsn="prometheus://host:port"`

or

`--stats.dsn="prom://host:port"`

The `host` and `port` are optional. If set they will start a prometheus http server on the specified host and port.

**Victoria Metrics:**

`--stats.dsn="victoriametrics://host:port"`

or

`--stats.dsn="vm://host:port"`

The `host` and `port` are optional. If set they will start a victoria metrics http server on the specified host and port.

**l2met:**

`--stats.dsn="l2met://"`

This report has no exposed options.

#### FlagStatsInterval: *--stats.interval, $STATS_INTERVAL*

This flag sets the interval at which the aggregated stats will be reported to the reporter.

Example: `--stats-interval=10s`

#### FlagStatsPrefix: *--stats.prefix, $STATS_PREFIX*

This flag sets the prefix attached to all stats keys.

Example: `--stats.prefix=my-app.server`

#### FlagStatsTags: *--stats.tags, $STATS_TAGS*

This flag sets tag key value pairs to set on all stats. This flag can be specified multiple times.

Example: `--stats.tags="app=my-app" --stats.tags="zone=eu-west"`

### Tracer

The tracing flags are used by `cmd.NewTracer` to create a new open telemetry `trace.TraceProvider`.

#### FlagTracingExporter: *--tracing.exporter, $TRACING_EXPORTER*

This flag sets the exporter to send spans to. The available options are `zipkin`, `otlphttp` and `otlpgrpc`.

Example: `--tracing.exporter=otlphttp`

#### FlagTracingEndpoint: *--tracing.endpoint, $TRACING_ENDPOINT*

This flag sets the endpoint the exporter should send traces to.

Example: `--tracing.endpoint="agent-host:port"` or `--tracing.endpoint="http://host:port/api/v2"`

#### FlagTracingEndpointInsecure: *--tracing.endpoint-insecure, $TRACING_ENDPOINT_INSECURE*

This flag sets the endpoint the exporter should send traces to.

Example: `--tracing.endpoint-insecure`

#### FlagTracingRatio: *--tracing.ratio, $TRACING_RATIO*

This flag sets the sample ratio of spans that will be reported to the exporter. This should be between 0 and 1.

Example: `--tracing.ratio=0.2`

#### FlagTracingTags: *--tracing.tags, $TRACING_TAGS*

This flag sets a list of tags appended to every trace. This flag can be specified multiple times.

Example: `--tracing.tags="app=my-app" --tracing.tags="zone=eu-west"`

### Observer

The observe package exposes an `Observer` type which is essentially a helper that combines a logger, tracer and statter.
It is useful if you use all three for your services and want to avoid carrying around many arguments.

Here is an example of how one might use it:

```go
func yourAction(c *cli.Context) error {
     obsvr, err := newObserver(c)
    if err != nil {
        return err
    }
    defer obsvr.Close()

	// Run your application here...

	return nil
}

func newObserver(c *cli.Context) (*observe.Observer, error) {
    log, err := cmd.NewLogger(c)
    if err != nil {
        return nil, err
    }

    stats, err := cmd.NewStatter(c, log)
    if err != nil {
        return nil, err
    }

    tracer, err := cmd.NewTracer(c, log,
        semconv.ServiceNameKey.String("my-service"),
        semconv.ServiceVersionKey.String("1.0.0"),
    )
    if err != nil {
        return nil, err
    }
    tracerCancel := func() { _ = tracer.Shutdown(context.Background()) }

    return observe.New(log, stats, tracer, tracerCancel), nil
}
```

It also exposes `NewFake` which allows you to pass fake loggers, tracers and statters in your tests easily.
