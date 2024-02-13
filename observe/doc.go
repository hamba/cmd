/*
Package observe implements a type that combines statter, logger and tracer.

Example usage:

	obsrv, err := observe.New(ctx, cliCmd, &observe.Options{})
	if err != nil {
		// Handle error.
		return
	}
*/
package observe
