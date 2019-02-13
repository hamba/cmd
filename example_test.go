package cmd_test

import (
	"fmt"

	"github.com/hamba/cmd"
	"gopkg.in/urfave/cli.v1"
)

func ExampleNewContext() {
	var c *cli.Context // Get this from your action

	ctx, err := cmd.NewContext(c)
	if err != nil {
		// Handle error
	}

	ctx.Logger()  // Get your logger
	ctx.Statter() // Get your statter

	<-cmd.WaitForSignals()
}

func ExampleSplitTags() {
	input := []string{"a=b", "foo=bar"} // Usually from a cli.StringSlice

	tags, err := cmd.SplitTags(input, "=")
	if err != nil {
		// Handle error
	}

	fmt.Println(tags)
	// Output: [a b foo bar]
}

func ExampleWaitForSignals() {
	<-cmd.WaitForSignals() // Will wait for SIGTERM or SIGINT
}
