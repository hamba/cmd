/*
Package observe implements a type that combines statter, logger and tracer.

Example usage:

	func newObserver(c *cli.Context) (*observe.Observer, error) {
	    log, err := cmd.NewLogger(c)
	    if err != nil {
	    	return nil, err
	    }

	    stats, err := cmd.NewStatter(c, log)
	    if err != nil {
	    	return nil, err
	    }

	    tracer, err := cmd.NewTracer(c, log,
	    	semconv.ServiceNameKey.String("my-service"),
	    	semconv.ServiceVersionKey.String("1.0.0"),
	    )
	    if err != nil {
	    	return nil, err
	    }
	    tracerCancel := func() { _ = tracer.Shutdown(context.Background()) }

	    return observe.New(log, stats, tracer, tracerCancel), nil
    }
*/
package observe
