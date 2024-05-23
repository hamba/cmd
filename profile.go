package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/hamba/logger/v2"
	"github.com/urfave/cli/v2"
)

var allProfilingTypes = []pyroscope.ProfileType{
	pyroscope.ProfileCPU,
	pyroscope.ProfileInuseObjects,
	pyroscope.ProfileAllocObjects,
	pyroscope.ProfileInuseSpace,
	pyroscope.ProfileAllocSpace,
	pyroscope.ProfileGoroutines,
	pyroscope.ProfileMutexCount,
	pyroscope.ProfileMutexDuration,
	pyroscope.ProfileBlockCount,
	pyroscope.ProfileBlockDuration,
}

// Tracing flag constants declared for CLI use.
const (
	FlagProfilingDSN      = "profiling.dsn"
	FlagProfileUploadRate = "profiling.upload-rate"
	FlagProfilingTags     = "profiling.tags"
	FlagProfilingTypes    = "profiling.types"
)

// CategoryProfiling is the profiling category.
const CategoryProfiling = "Profiling"

// ProfilingFlags are flags that configure profiling.
var ProfilingFlags = Flags{
	&cli.StringFlag{
		Name:     FlagProfilingDSN,
		Category: CategoryProfiling,
		Usage: "The address to the Pyroscope server, in the format: " +
			"'http://basic:auth@server:port?token=auth-token&tenantid=tenant-id'.",
		EnvVars: []string{"PROFILING_DSN"},
	},
	&cli.DurationFlag{
		Name:     FlagProfileUploadRate,
		Category: CategoryProfiling,
		Usage:    "The rate at which profiles are uploaded.",
		Value:    15 * time.Second,
		EnvVars:  []string{"PROFILING_UPLOAD_RATE"},
	},
	&cli.StringSliceFlag{
		Name:     FlagProfilingTags,
		Category: CategoryProfiling,
		Usage:    "A list of tags appended to every profile. Format: key=value.",
		EnvVars:  []string{"PROFILING_TAGS"},
	},
	&cli.StringSliceFlag{
		Name:     FlagProfilingTypes,
		Category: CategoryProfiling,
		Usage:    "The type of profiles to include. Defaults to all.",
		EnvVars:  []string{"PROFILING_TYPES"},
	},
}

// NewProfiler returns a profiler configured from the cli.
// If no profiler is configured, nil is returned.
func NewProfiler(c *cli.Context, svc string, log *logger.Logger) (*pyroscope.Profiler, error) {
	dsn := c.String(FlagProfilingDSN)
	if dsn == "" {
		//nolint:nilnil // There is no sentinel in this case.
		return nil, nil
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing profiling DSN: %w", err)
	}

	tenantID := u.Query().Get("tenantid")

	authToken := u.Query().Get("token")
	var username, password string
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	if (username != "" || password != "") && authToken != "" {
		return nil, errors.New("cannot set auth token and basic auth")
	}

	srvURL := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}

	var tags map[string]string
	if pairs := c.StringSlice(FlagProfilingTags); len(pairs) > 0 {
		tags, err = sliceToMap(pairs)
		if err != nil {
			return nil, err
		}
	}

	types := allProfilingTypes
	if newTypes := c.StringSlice(FlagProfilingTypes); len(newTypes) > 0 {
		types = make([]pyroscope.ProfileType, len(newTypes))
		for i, typ := range newTypes {
			types[i] = pyroscope.ProfileType(typ)
		}
	}

	cfg := pyroscope.Config{
		ApplicationName:   svc,
		Tags:              tags,
		ServerAddress:     srvURL.String(),
		AuthToken:         authToken,
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		TenantID:          tenantID,
		UploadRate:        c.Duration(FlagProfileUploadRate),
		Logger:            pyroLogAdapter{log: log},
		ProfileTypes:      types,
	}

	return pyroscope.Start(cfg)
}

type pyroLogAdapter struct {
	log *logger.Logger
}

func (a pyroLogAdapter) Infof(format string, args ...any) {
	a.log.Info(fmt.Sprintf(format, args...))
}

func (a pyroLogAdapter) Debugf(format string, args ...any) {
	a.log.Trace(fmt.Sprintf(format, args...))
}

func (a pyroLogAdapter) Errorf(format string, args ...any) {
	a.log.Error(fmt.Sprintf(format, args...))
}
