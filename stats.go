package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/hamba/statter/v2"
	"github.com/hamba/statter/v2/reporter/l2met"
	"github.com/hamba/statter/v2/reporter/prometheus"
	"github.com/hamba/statter/v2/reporter/statsd"
	"github.com/hamba/statter/v2/reporter/victoriametrics"
	"github.com/urfave/cli/v2"
)

// Stats flag constants declared for CLI use.
const (
	FlagStatsDSN      = "stats.dsn"
	FlagStatsInterval = "stats.interval"
	FlagStatsPrefix   = "stats.prefix"
	FlagStatsTags     = "stats.tags"
)

// StatsFlags are flags that configure stats.
var StatsFlags = Flags{
	&cli.StringFlag{
		Name:    FlagStatsDSN,
		Usage:   "The DSN of a stats backend.",
		EnvVars: []string{"STATS_DSN"},
	},
	&cli.DurationFlag{
		Name:    FlagStatsInterval,
		Usage:   "The frequency at which the stats are reported.",
		Value:   time.Second,
		EnvVars: []string{"STATS_INTERVAL"},
	},
	&cli.StringFlag{
		Name:    FlagStatsPrefix,
		Usage:   "The prefix of the measurements names.",
		EnvVars: []string{"STATS_PREFIX"},
	},
	&cli.StringSliceFlag{
		Name:    FlagStatsTags,
		Usage:   "A list of tags appended to every measurement. Format: key=value.",
		EnvVars: []string{"STATS_TAGS"},
	},
}

// NewStatter returns a statter configured from the cli.
func NewStatter(c *cli.Context, log *logger.Logger, opts ...statter.Option) (*statter.Statter, error) {
	r, err := createReporter(c, log)
	if err != nil {
		return nil, err
	}

	intv := c.Duration(FlagStatsInterval)
	prefix, tags, err := statsWith(c)
	if err != nil {
		return nil, err
	}

	opts = append(opts, statter.WithPrefix(prefix), statter.WithTags(tags...))

	return statter.New(r, intv, opts...), nil
}

func statsWith(c *cli.Context) (string, []statter.Tag, error) {
	strTags, err := Split(c.StringSlice(FlagStatsTags), "=")
	if err != nil {
		return "", nil, err
	}
	tags := make([]statter.Tag, len(strTags))
	for i, st := range strTags {
		tags[i] = statter.Tag{st[0], st[1]}
	}

	return c.String(FlagStatsPrefix), tags, nil
}

func createReporter(c *cli.Context, log *logger.Logger) (statter.Reporter, error) {
	dsn := c.String(FlagStatsDSN)
	if dsn == "" {
		return statter.DiscardReporter, nil
	}

	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch uri.Scheme {
	case "statsd":
		return newStatsd(uri.Host, uri.Query())
	case "l2met":
		return l2met.New(log, ""), nil
	case "prometheus", "prom":
		return newPrometheusStats(uri.Host, log), nil
	case "victoriametrics", "vm":
		return newVictoriaMetricsStats(uri.Host, log), nil
	default:
		return nil, fmt.Errorf("unsupported stats reporter: %s", uri.Scheme)
	}
}

func newStatsd(addr string, qry url.Values) (*statsd.Statsd, error) {
	var opts []statsd.Option
	if s := qry.Get("flushBytes"); s != "" {
		n, err := strconv.Atoi(s)
		if err == nil {
			opts = append(opts, statsd.WithFlushBytes(n))
		}
	}
	if s := qry.Get("flushInterval"); s != "" {
		d, err := time.ParseDuration(s)
		if err == nil {
			opts = append(opts, statsd.WithFlushInterval(d))
		}
	}

	return statsd.New(addr, "", opts...)
}

func newPrometheusStats(addr string, log *logger.Logger) *prometheus.Prometheus {
	r := prometheus.New("")

	if addr != "" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", r.Handler())
		go func() {
			srv := http.Server{
				Addr:              addr,
				Handler:           mux,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      10 * time.Second,
				IdleTimeout:       120 * time.Second,
			}
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error(err.Error(), ctx.Str("server", "prometheus"))
			}
		}()
	}

	return r
}

func newVictoriaMetricsStats(addr string, log *logger.Logger) *victoriametrics.VictoriaMetrics {
	r := victoriametrics.New()

	if addr != "" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", r.Handler())
		go func() {
			srv := http.Server{
				Addr:              addr,
				Handler:           mux,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      10 * time.Second,
				IdleTimeout:       120 * time.Second,
			}
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error(err.Error(), ctx.Str("server", "victoria-metrics"))
			}
		}()
	}

	return r
}
