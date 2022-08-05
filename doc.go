/*
Package cmd implements cmd helpers.

This provides helpers on top of `github.com/urfave/cli`.

Example usage:

	var c *cli.Context // Get this from your action

	log, err := cmd.NewLogger(c)
	if err != nil {
		// Handle error.
	}

	stats, err := cmd.NewStatter(c, log)
	if err != nil {
		// Handle error.
	}
*/
package cmd
