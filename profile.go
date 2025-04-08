package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ettle/strcase"
	"github.com/grafana/pyroscope-go"
	"github.com/hamba/logger/v2"
	"github.com/urfave/cli/v3"
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
		Sources: cli.EnvVars(strcase.ToSNAKE(FlagProfilingDSN)),
	},
	&cli.DurationFlag{
		Name:     FlagProfileUploadRate,
		Category: CategoryProfiling,
		Usage:    "The rate at which profiles are uploaded.",
		Value:    15 * time.Second,
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagProfileUploadRate)),
	},
	&cli.StringMapFlag{
		Name:     FlagProfilingTags,
		Category: CategoryProfiling,
		Usage:    "A list of tags appended to every profile.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagProfilingTags)),
	},
	&cli.StringSliceFlag{
		Name:     FlagProfilingTypes,
		Category: CategoryProfiling,
		Usage:    "The type of profiles to include. Defaults to all.",
		Sources:  cli.EnvVars(strcase.ToSNAKE(FlagProfilingTypes)),
	},
}

// NewProfiler returns a profiler configured from the cli.
// If no profiler is configured, nil is returned.
func NewProfiler(cmd *cli.Command, svc string, log *logger.Logger) (*pyroscope.Profiler, error) {
	dsn := cmd.String(FlagProfilingDSN)
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

	types := allProfilingTypes
	if newTypes := cmd.StringSlice(FlagProfilingTypes); len(newTypes) > 0 {
		types = make([]pyroscope.ProfileType, len(newTypes))
		for i, typ := range newTypes {
			types[i] = pyroscope.ProfileType(typ)
		}
	}

	cfg := pyroscope.Config{
		ApplicationName:   svc,
		Tags:              cmd.StringMap(FlagProfilingTags),
		ServerAddress:     srvURL.String(),
		AuthToken:         authToken,
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		TenantID:          tenantID,
		UploadRate:        cmd.Duration(FlagProfileUploadRate),
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
