package cmd

import (
	"gopkg.in/urfave/cli.v1"
)

// Flag constants declared for CLI use.
const (
	FlagPort = "port"

	FlagLogFormat = "log.format"
	FlagLogLevel  = "log.level"
	FlagLogTags   = "log.tags"

	FlagStatsDSN    = "stats.dsn"
	FlagStatsPrefix = "stats.prefix"
	FlagStatsTags   = "stats.tags"
)

// Flags represents a set of CLI flags.
type Flags []cli.Flag

// Merge joins one or more Flags together, making a new set.
func (f Flags) Merge(flags ...Flags) Flags {
	var m Flags
	m = append(m, f...)
	for _, flag := range flags {
		m = append(m, flag...)
	}

	return m
}

// ServerFlags are flags that configure a server.
var ServerFlags = Flags{
	cli.StringFlag{
		Name:   FlagPort,
		Value:  "80",
		Usage:  "Port for HTTP server to listen on",
		EnvVar: "PORT",
	},
}

// CommonFlags are flags that configure logging and stats.
var CommonFlags = Flags{
	cli.StringFlag{
		Name:   FlagLogFormat,
		Usage:  "Specify the format of logs. Supported formats: 'logfmt', 'json'",
		EnvVar: "LOG_FORMAT",
	},
	cli.StringFlag{
		Name:   FlagLogLevel,
		Value:  "info",
		Usage:  "Specify the log level. E.g. 'debug', 'warning'.",
		EnvVar: "LOG_LEVEL",
	},
	cli.StringSliceFlag{
		Name:   FlagLogTags,
		Usage:  "A list of tags appended to every log. Format: key=value.",
		EnvVar: "LOG_TAGS",
	},

	cli.StringFlag{
		Name:   FlagStatsDSN,
		Usage:  "The URL of a stats backend.",
		EnvVar: "STATS_DSN",
	},
	cli.StringFlag{
		Name:   FlagStatsPrefix,
		Usage:  "The prefix of the measurements names.",
		EnvVar: "STATS_PREFIX",
	},
	cli.StringSliceFlag{
		Name:   FlagStatsTags,
		Usage:  "A list of tags appended to every measurement. Format: key=value.",
		EnvVar: "STATS_TAGS",
	},
}
