package cmd

import (
	"fmt"
	"strings"
)

// Split splits a slice of strings into an slice of arrays using
// the given separator.
func Split(slice []string, sep string) ([][2]string, error) {
	res := make([][2]string, 0, len(slice))
	for _, str := range slice {
		parts := strings.SplitN(str, sep, 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("string %q does not contain %q", str, sep)
		}

		res = append(res, [2]string{parts[0], parts[1]})
	}

	return res, nil
}
