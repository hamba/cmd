package cmd

import (
	"io"
	"os"

	"github.com/ettle/strcase"
	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/urfave/cli/v3"
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
		Usage:    "Specify the format of logs. Supported formats: 'logfmt', 'json', 'console'",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagLogFormat)),
	},
	&cli.StringFlag{
		Name:     FlagLogLevel,
		Category: CategoryLog,
		Value:    "info",
		Usage:    "Specify the log level. e.g. 'debug', 'info', 'error'.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagLogLevel)),
	},
	&cli.StringMapFlag{
		Name:     FlagLogCtx,
		Category: CategoryLog,
		Usage:    "A list of context field appended to every log. Format: key=value.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagLogCtx)),
	},
}

// LoggerOptions are options for creating a logger.
type LoggerOptions struct {
	Writer io.Writer
}

// NewLogger returns a logger configured from the cli.
func NewLogger(cmd *cli.Command) (*logger.Logger, error) {
	return NewLoggerWithOptions(cmd, &LoggerOptions{})
}

// NewLoggerWithOptions returns a logger configured from the cli.
func NewLoggerWithOptions(cmd *cli.Command, opts *LoggerOptions) (*logger.Logger, error) {
	str := cmd.String(FlagLogLevel)
	if str == "" {
		str = "info"
	}

	lvl, err := logger.LevelFromString(str)
	if err != nil {
		return nil, err
	}

	fmtr := newLogFormatter(cmd)

	tags := cmd.StringMap(FlagLogCtx)

	fields := make([]logger.Field, 0, len(tags))
	for k, v := range tags {
		fields = append(fields, ctx.Str(k, v))
	}

	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	return logger.New(w, fmtr, lvl).With(fields...), nil
}

func newLogFormatter(cmd *cli.Command) logger.Formatter {
	format := cmd.String(FlagLogFormat)
	switch format {
	case "json":
		return logger.JSONFormat()
	case "console":
		return logger.ConsoleFormat()
	default:
		return logger.LogfmtFormat()
	}
}
