package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
	"github.com/hamba/statter/l2met"
	"github.com/hamba/statter/statsd"
	"gopkg.in/urfave/cli.v2"
)

// NewStats creates a new statter.
func NewStats(c *cli.Context, l log.Logger) (stats.Statter, error) {
	var s stats.Statter
	var err error

	dsn := c.String(FlagStatsDSN)
	if dsn == "" {
		return stats.Null, nil
	}

	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch uri.Scheme {
	case "statsd":
		s, err = newStatsd(c, uri.Host)
		if err != nil {
			return nil, err
		}

	case "l2met":
		s = newL2met(c, l)

	default:
		return nil, fmt.Errorf("unsupported stats type: %s", uri.Scheme)
	}

	tags, err := SplitTags(c.StringSlice(FlagStatsTags), "=")
	if err != nil {
		return nil, err
	}
	if len(tags) > 0 {
		s = stats.NewTaggedStatter(s, tags...)
	}

	return s, nil
}

func newStatsd(c *cli.Context, addr string) (stats.Statter, error) {
	s, err := statsd.NewBuffered(addr, c.String(FlagStatsPrefix), statsd.WithFlushInterval(1*time.Second))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func newL2met(c *cli.Context, l log.Logger) stats.Statter {
	return l2met.New(l, c.String(FlagStatsPrefix))
}
