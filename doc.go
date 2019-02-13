/*
Package cmd implements cmd helper.

This provides helpers on top of `github.com/urfave/cli`.

Example usage:
	var c *cli.Context // Get this from your action

	ctx, err := cmd.NewContext(c)
	if err != nil {
		// Handle error
	}

	ctx.Logger()  // Get your logger
	ctx.Statter() // Get your statter

	<-cmd.WaitForSignals()
*/
package cmd
