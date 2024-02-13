package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ettle/strcase"
	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/hamba/statter/v2"
	"github.com/hamba/statter/v2/reporter/l2met"
	"github.com/hamba/statter/v2/reporter/prometheus"
	"github.com/hamba/statter/v2/reporter/statsd"
	"github.com/hamba/statter/v2/reporter/victoriametrics"
	"github.com/urfave/cli/v3"
)

// Stats flag constants declared for CLI use.
const (
	FlagStatsDSN      = "stats.dsn"
	FlagStatsInterval = "stats.interval"
	FlagStatsPrefix   = "stats.prefix"
	FlagStatsTags     = "stats.tags"
)

// CategoryStats is the stats flag category.
const CategoryStats = "Stats"

// StatsFlags are flags that configure stats.
var StatsFlags = Flags{
	&cli.StringFlag{
		Name:     FlagStatsDSN,
		Category: CategoryStats,
		Usage:    "The DSN of a stats backend.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagStatsDSN)),
	},
	&cli.DurationFlag{
		Name:     FlagStatsInterval,
		Category: CategoryStats,
		Usage:    "The frequency at which the stats are reported.",
		Value:    time.Second,
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagStatsInterval)),
	},
	&cli.StringFlag{
		Name:     FlagStatsPrefix,
		Category: CategoryStats,
		Usage:    "The prefix of the measurements names.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagStatsPrefix)),
	},
	&cli.StringMapFlag{
		Name:     FlagStatsTags,
		Category: CategoryStats,
		Usage:    "A list of tags appended to every measurement.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagStatsTags)),
	},
}

// NewStatter returns a statter configured from the cli.
func NewStatter(cmd *cli.Command, log *logger.Logger, opts ...statter.Option) (*statter.Statter, error) {
	r, err := createReporter(cmd, log)
	if err != nil {
		return nil, err
	}

	intv := cmd.Duration(FlagStatsInterval)
	prefix, tags := statsWith(cmd)

	opts = append(opts, statter.WithPrefix(prefix), statter.WithTags(tags...))

	return statter.New(r, intv, opts...), nil
}

func statsWith(cmd *cli.Command) (string, []statter.Tag) {
	strTags := cmd.StringMap(FlagStatsTags)

	tags := make([]statter.Tag, 0, len(strTags))
	for k, v := range strTags {
		tags = append(tags, statter.Tag{k, v})
	}

	return cmd.String(FlagStatsPrefix), tags
}

func createReporter(cmd *cli.Command, log *logger.Logger) (statter.Reporter, error) {
	dsn := cmd.String(FlagStatsDSN)
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
