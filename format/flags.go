// Package format changes the cli Flag string formatting to be simpler to read.
// The format is as follows:
//
//	 ```
//		Flags [Env Vars] [Default]
//		   Usage
//	 ```
package format

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

func init() {
	cli.FlagStringer = structuredFlagStringer
}

func structuredFlagStringer(f cli.Flag) string {
	df, ok := f.(cli.DocGenerationFlag)
	if !ok {
		return ""
	}

	placeholder, usage := unquoteUsage(df.GetUsage())
	needsPlaceholder := df.TakesValue()
	if needsPlaceholder && placeholder == "" {
		placeholder = "value"
	}

	// Do not set defaults for bool flags.
	defaultValueString := ""
	if bf, ok := f.(*cli.BoolFlag); !ok || !bf.DisableDefaultText {
		if s := df.GetDefaultText(); s != "" {
			defaultValueString = fmt.Sprintf(" (default: %s)", s)
		}
	}

	pn := prefixedNames(df.Names(), placeholder)
	sliceFlag, ok := f.(cli.DocGenerationSliceFlag)
	if ok && sliceFlag.IsSliceFlag() {
		pn = pn + " [" + pn + "]"
	}

	optLine := pn + envHint(df.GetEnvVars()) + defaultValueString
	return optLine + "\n" + nIndent(usage, 7) + "\n"
}

// unquoteUsage finds the first quoted (“) substring, returning it with the
// unquoted usage.
func unquoteUsage(usage string) (string, string) {
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name := usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break
		}
	}
	return "", usage
}

func prefixedNames(names []string, placeholder string) string {
	var prefixed string
	for i, name := range names {
		if name == "" {
			continue
		}

		prefixed += prefixFor(name) + name
		if placeholder != "" {
			prefixed += " " + placeholder
		}
		if i < len(names)-1 {
			prefixed += ", "
		}
	}
	return prefixed
}

func prefixFor(name string) (prefix string) {
	if len(name) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}

	return
}

func envHint(envVars []string) string {
	if runtime.GOOS != "windows" || os.Getenv("PSHOME") != "" {
		return envFormat(envVars, "$", ", $", "")
	}
	return envFormat(envVars, "%", "%, %", "%")
}

func envFormat(envVars []string, prefix, sep, suffix string) string {
	if len(envVars) == 0 {
		return ""
	}
	return fmt.Sprintf(" [%s%s%s]", prefix, strings.Join(envVars, sep), suffix)
}

func nIndent(s string, n int) string {
	indent := strings.Repeat(" ", n)
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return indent + strings.Join(lines, "\n"+indent)
}
