/*
Package observe implements a type that combines statter, logger and tracer.

Example usage:

	func New(c *cli.Context, svc, version string) (*observe.Observer, error) {
		log, logCancel, err := NewLogger(c, svc)
		if err != nil {
			return nil, err
		}

		stats, statsCancel, err := NewStatter(c, log, svc)
		if err != nil {
			logCancel()
			return nil, err
		}

		tracer, traceCancel, err := NewTracer(c, log, svc, version)
		if err != nil {
			logCancel()
			statsCancel()
			return nil, err
		}

		return observe.New(log, stats, tracer, traceCancel, statsCancel, logCancel), nil
	}
*/
package observe
