![Logo](http://svg.wiersma.co.za/hamba/project?title=cmd&tag=Go%20cmd%20helper)

[![Go Report Card](https://goreportcard.com/badge/github.com/hamba/cmd)](https://goreportcard.com/report/github.com/hamba/cmd)
[![Build Status](https://github.com/hamba/cmd/actions/workflows/test.yml/badge.svg)](https://github.com/hamba/cmd/actions)
[![Coverage Status](https://coveralls.io/repos/github/hamba/cmd/badge.svg?branch=master)](https://coveralls.io/github/hamba/cmd?branch=master)
[![GoDoc](https://godoc.org/github.com/hamba/cmd?status.svg)](https://godoc.org/github.com/hamba/cmd)
[![GitHub release](https://img.shields.io/github/release/hamba/cmd.svg)](https://github.com/hamba/cmd/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hamba/cmd/master/LICENSE)

Go cmd helper. 

This provides helpers on top of `github.com/urfave/cli`.

## Overview

Install with:

```shell
go get github.com/hamba/cmd
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

    tracer, err := cmd.NewTracer(c,
        semconv.ServiceNameKey.String("my-service"),
        semconv.ServiceVersionKey.String("1.0.0"),
    )
    if err != nil {
        // Handle error.
        return
    }
    defer tracer.Shutdown(context.Background())

    // Run your application here...
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

This flag sets the DSN describing the stats reporter to use. The available options are `statsd`, `prometheus`, `l2met`.

The DSN can in some situations specify the host and configuration values as shown in the below examples:

**Statsd:** 

`--stats.dsn="statsd://host:port?flushBytes=1432&flushInterval=10s"`

The `host` and `port` are required. Optionally `flushBytes` and `flushInterval` can be set, controlling how often the stats will
be sent to the Statsd server.

**Prometheus:**

`--stats.dsn="prometheus://host:port"`

The `host` and `port` are optional. If set they will start a prometheus http server on the specified host and port.

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

This flag sets the exporter to send spans to. The available options are `jaeger` and `zipkin`.

Example: `--tracing.exporter=jaeger`

#### FlagTracingEndpoint: *--tracing.endpoint, $TRACING_ENDPOINT*

This flag sets the endpoint the exporter should send traces to.

Example: `--tracing.endpoint="host:port"` or `--tracing.endpoint="http://host:port/api/v2"`

#### FlagTracingRatio: *--tracing.ratio, $TRACING_RATIO*

This flag sets the sample ratio of spans that will be reported to the exporter. This should be between 0 and 1.

Example: `--tracing.ratio=0.2`
