package cmd

import (
	"io"

	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
	"github.com/urfave/cli/v2"
)

// Context represents an application context.
//
// Context implements both log.Loggable and stats.Statable.
type Context struct {
	*cli.Context

	logger  log.Logger
	statter stats.Statter
}

// NewContext creates a new Context from the CLI Context.
func NewContext(c *cli.Context) (*Context, error) {
	l, err := NewLogger(c)
	if err != nil {
		return nil, err
	}

	s, err := NewStats(c, l)
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		Context: c,
		logger:  l,
		statter: s,
	}

	return ctx, nil
}

// Logger returns the Logger attached to the Context.
func (c *Context) Logger() log.Logger {
	return c.logger
}

// AttachLogger attaches a Logger to the Context.
func (c *Context) AttachLogger(fn func(l log.Logger) log.Logger) {
	c.logger = fn(c.logger)
}

// Statter returns the Statter attached to the Context.
func (c *Context) Statter() stats.Statter {
	return c.statter
}

// AttachStatter attaches a Statter to the Context.
func (c *Context) AttachStatter(fn func(s stats.Statter) stats.Statter) {
	c.statter = fn(c.statter)
}

// Close closes the context.
func (c *Context) Close() error {
	if err := c.statter.Close(); err != nil {
		return err
	}

	if l, ok := c.logger.(io.Closer); ok {
		return l.Close()
	}

	return nil
}
