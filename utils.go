package cmd

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// SplitTags splits a slice of strings into a slice using
// the given separator.
func SplitTags(slice []string, sep string) ([]interface{}, error) {
	res := make([]interface{}, 2*len(slice))

	for i, str := range slice {
		parts := strings.SplitN(str, sep, 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid tags string")
		}

		res[2*i] = parts[0]
		res[2*i+1] = parts[1]
	}

	return res, nil
}

// WaitForSignals waits for SIGINT or SIGTERM signals.
func WaitForSignals() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return sigs
}
