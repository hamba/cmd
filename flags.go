package cmd

import (
	"github.com/urfave/cli/v2"
)

// FlagPort contains the flag name for a server port.
const FlagPort = "port"

// Flags represents a set of CLI flags.
type Flags []cli.Flag

// Merge joins one or more Flags together, making a new set.
func (f Flags) Merge(flags ...Flags) Flags {
	m := make(Flags, 0, len(f))
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

// MonitoringFlags are flags that configure logging, stats and tracing.
var MonitoringFlags = Flags{}.Merge(LogFlags, StatsFlags, TracingFlags)
