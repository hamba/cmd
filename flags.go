package cmd

import (
	"github.com/urfave/cli/v3"
)

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

// MonitoringFlags are flags that configure logging, stats, profiling and tracing.
var MonitoringFlags = Flags{}.Merge(LogFlags, StatsFlags, ProfilingFlags, TracingFlags)
