// Package term implements a unified way to present output.
package term

import (
	"io"

	"github.com/fatih/color"
)

// Term represents an object that can output to the terminal.
type Term interface {
	Output(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

// Basic writes messages to a writer. If an error writer
// is provided it is used for error messages.
type Basic struct {
	Writer      io.Writer
	ErrorWriter io.Writer
	Verbose     bool
}

// Output writes general messages that are suppressed when not in
// verbose mode.
func (b Basic) Output(msg string) {
	if !b.Verbose {
		return
	}
	_, _ = b.Writer.Write([]byte(msg + "\n"))
}

// Info writes information messages.
func (b Basic) Info(msg string) {
	_, _ = b.Writer.Write([]byte(msg + "\n"))
}

// Warning writes warning messages.
func (b Basic) Warning(msg string) {
	b.Info(msg)
}

// Error writes error messages.
func (b Basic) Error(msg string) {
	w := b.Writer
	if b.ErrorWriter != nil {
		w = b.ErrorWriter
	}

	_, _ = w.Write([]byte(msg + "\n"))
}

// Prefixed writes prefixed messages to the underlying term.
type Prefixed struct {
	OutputPrefix  string
	InfoPrefix    string
	WarningPrefix string
	ErrorPrefix   string
	Term          Term
}

// Output writes general messages.
func (p Prefixed) Output(msg string) {
	p.Term.Output(p.OutputPrefix + msg)
}

// Info writes information messages.
func (p Prefixed) Info(msg string) {
	p.Term.Info(p.InfoPrefix + msg)
}

// Warning writes warning messages.
func (p Prefixed) Warning(msg string) {
	p.Term.Warning(p.WarningPrefix + msg)
}

// Error writes error messages.
func (p Prefixed) Error(msg string) {
	p.Term.Error(p.ErrorPrefix + msg)
}

// Predefined colors for use with Coloured.
var (
	Blue   = color.New(color.FgHiBlue)
	Red    = color.New(color.FgHiRed)
	Cyan   = color.New(color.FgHiCyan)
	Green  = color.New(color.FgHiGreen)
	Black  = color.New(color.FgHiBlack)
	Yellow = color.New(color.FgHiYellow)
	White  = color.New(color.FgHiWhite)
)

// Colored writer coloured messages to the underlying term.
type Colored struct {
	OutputColor  *color.Color
	InfoColor    *color.Color
	WarningColor *color.Color
	ErrorColor   *color.Color
	Term         Term
}

// Output writes general messages.
func (c Colored) Output(msg string) {
	c.Term.Output(c.OutputColor.Sprint(msg))
}

// Info writes information messages.
func (c Colored) Info(msg string) {
	c.Term.Info(c.InfoColor.Sprint(msg))
}

// Warning writes warning messages.
func (c Colored) Warning(msg string) {
	c.Term.Warning(c.WarningColor.Sprint(msg))
}

// Error writes error messages.
func (c Colored) Error(msg string) {
	c.Term.Error(c.ErrorColor.Sprint(msg))
}
