package cmd

import (
	"gopkg.in/urfave/cli.v2"
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
	&cli.StringFlag{
		Name:    FlagPort,
		Value:   "80",
		Usage:   "Port for HTTP server to listen on",
		EnvVars: []string{"PORT"},
	},
}

// LogFlags are flags that configure logging.
var LogFlags = Flags{
	&cli.StringFlag{
		Name:    FlagLogFormat,
		Usage:   "Specify the format of logs. Supported formats: 'logfmt', 'json'",
		EnvVars: []string{"LOG_FORMAT"},
	},
	&cli.StringFlag{
		Name:    FlagLogLevel,
		Value:   "info",
		Usage:   "Specify the log level. E.g. 'debug', 'warning'.",
		EnvVars: []string{"LOG_LEVEL"},
	},
	&cli.StringSliceFlag{
		Name:    FlagLogTags,
		Usage:   "A list of tags appended to every log. Format: key=value.",
		EnvVars: []string{"LOG_TAGS"},
	},
}

// StatsFlags are flags that configure stats.
var StatsFlags = Flags{
	&cli.StringFlag{
		Name:    FlagStatsDSN,
		Usage:   "The URL of a stats backend.",
		EnvVars: []string{"STATS_DSN"},
	},
	&cli.StringFlag{
		Name:    FlagStatsPrefix,
		Usage:   "The prefix of the measurements names.",
		EnvVars: []string{"STATS_PREFIX"},
	},
	&cli.StringSliceFlag{
		Name:    FlagStatsTags,
		Usage:   "A list of tags appended to every measurement. Format: key=value.",
		EnvVars: []string{"STATS_TAGS"},
	},
}

// CommonFlags are flags that configure logging and stats.
//
// Common flags include LogFlags and StatsFlags.
var CommonFlags = Flags{}.Merge(LogFlags, StatsFlags)
