package observe_test

import (
	"context"

	"github.com/hamba/cmd/v3/observe"
	"github.com/urfave/cli/v3"
)

func ExampleNew() {
	var (
		ctx    context.Context
		cliCmd *cli.Command // Get this from your action
	)

	obsrv, err := observe.New(ctx, cliCmd, "my-service", &observe.Options{})
	if err != nil {
		// Handle error.
		return
	}

	_ = obsrv
}
