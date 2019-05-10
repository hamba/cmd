package cmd

import (
	"os"
	"time"

	"github.com/hamba/logger"
	"gopkg.in/urfave/cli.v1"
)

// NewLogger creates a new logger.
func NewLogger(c *cli.Context) (logger.Logger, error) {
	str := c.GlobalString(FlagLogLevel)
	if str == "" {
		str = "info"
	}

	lvl, err := logger.LevelFromString(str)
	if err != nil {
		return nil, err
	}

	fmtr := newLogFormatter(c)
	h := logger.BufferedStreamHandler(os.Stdout, 2000, 1*time.Second, fmtr)
	h = logger.LevelFilterHandler(lvl, h)

	tags, err := SplitTags(c.GlobalStringSlice(FlagLogTags), "=")
	if err != nil {
		return nil, err
	}

	ctx := make([]interface{}, len(tags))
	for i, t := range tags {
		ctx[i] = t
	}

	return logger.New(h, ctx...), nil
}

func newLogFormatter(c *cli.Context) logger.Formatter {
	format := c.GlobalString(FlagLogFormat)
	switch format {

	case "json":
		return logger.JSONFormat()

	default:
		return logger.LogfmtFormat()
	}
}
