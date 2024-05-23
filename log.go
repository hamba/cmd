package cmd

import (
	"os"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/urfave/cli/v2"
)

// Log flag constants declared for CLI use.
const (
	FlagLogFormat = "log.format"
	FlagLogLevel  = "log.level"
	FlagLogCtx    = "log.ctx"
)

// CategoryLog is the log flag category.
const CategoryLog = "Logging"

// LogFlags are flags that configure logging.
var LogFlags = Flags{
	&cli.StringFlag{
		Name:     FlagLogFormat,
		Category: CategoryLog,
		Usage:    "Specify the format of logs. Supported formats: 'logfmt', 'json', 'console'.",
		EnvVars:  []string{"LOG_FORMAT"},
	},
	&cli.StringFlag{
		Name:     FlagLogLevel,
		Category: CategoryLog,
		Value:    "info",
		Usage:    "Specify the log level. e.g. 'trace', 'debug', 'info', 'error'.",
		EnvVars:  []string{"LOG_LEVEL"},
	},
	&cli.StringSliceFlag{
		Name:     FlagLogCtx,
		Category: CategoryLog,
		Usage:    "A list of context field appended to every log. Format: key=value.",
		EnvVars:  []string{"LOG_CTX"},
	},
}

// NewLogger returns a logger configured from the cli.
func NewLogger(c *cli.Context) (*logger.Logger, error) {
	str := c.String(FlagLogLevel)
	if str == "" {
		str = "info"
	}

	lvl, err := logger.LevelFromString(str)
	if err != nil {
		return nil, err
	}

	fmtr := newLogFormatter(c)

	tags, err := Split(c.StringSlice(FlagLogCtx), "=")
	if err != nil {
		return nil, err
	}

	fields := make([]logger.Field, len(tags))
	for i, t := range tags {
		fields[i] = ctx.Str(t[0], t[1])
	}

	return logger.New(os.Stdout, fmtr, lvl).With(fields...), nil
}

func newLogFormatter(c *cli.Context) logger.Formatter {
	format := c.String(FlagLogFormat)
	switch format {
	case "json":
		return logger.JSONFormat()
	case "console":
		return logger.ConsoleFormat()
	default:
		return logger.LogfmtFormat()
	}
}
