package cmd_test

import (
	"errors"
	"flag"
	"testing"
	"time"

	"github.com/hamba/cmd"
	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/urfave/cli.v1"
)

func TestNewContext(t *testing.T) {
	c, _ := newTestContext()

	ctx, err := cmd.NewContext(c)

	assert.NoError(t, err)
	assert.IsType(t, &cmd.Context{}, ctx)
}

func TestNewContext_LoggerError(t *testing.T) {
	c, flags := newTestContext()
	flags.String(cmd.FlagLogLevel, "test", "")

	_, err := cmd.NewContext(c)

	assert.Error(t, err)
}

func TestNewContext_StatterError(t *testing.T) {
	c, flags := newTestContext()
	flags.String(cmd.FlagStatsDSN, "test://", "")

	_, err := cmd.NewContext(c)

	assert.Error(t, err)
}

func TestContext_Logger(t *testing.T) {
	l := new(MockLogger)
	c, _ := newTestContext()
	ctx, _ := cmd.NewContext(c)
	ctx.AttachLogger(func(log.Logger) log.Logger { return l })

	assert.Equal(t, l, ctx.Logger())
}

func TestContext_Statter(t *testing.T) {
	s := new(MockStats)
	c, _ := newTestContext()
	ctx, _ := cmd.NewContext(c)
	ctx.AttachStatter(func(stats.Statter) stats.Statter { return s })

	assert.Equal(t, s, ctx.Statter())
}

func TestContext_Close(t *testing.T) {
	c, _ := newTestContext()
	ctx, _ := cmd.NewContext(c)
	ctx.AttachLogger(func(log.Logger) log.Logger { return log.Null })
	ctx.AttachStatter(func(stats.Statter) stats.Statter { return stats.Null })

	err := ctx.Close()

	assert.NoError(t, err)
}

func TestContext_CloseErrors(t *testing.T) {
	tests := []struct {
		name     string
		logErr   error
		statsErr error
	}{
		{
			name:     "No Error",
			logErr:   nil,
			statsErr: nil,
		},
		{
			name:     "Logger Error",
			logErr:   errors.New("test"),
			statsErr: nil,
		},
		{
			name:     "Stats Error",
			logErr:   nil,
			statsErr: errors.New("test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := new(MockLogger)
			l.On("Close").Return(tt.logErr)

			s := new(MockStats)
			s.On("Close").Return(tt.statsErr)

			c, _ := newTestContext()
			ctx, _ := cmd.NewContext(c)
			ctx.AttachLogger(func(log.Logger) log.Logger { return l })
			ctx.AttachStatter(func(stats.Statter) stats.Statter { return s })

			err := ctx.Close()

			if tt.logErr != nil || tt.statsErr != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func newTestContext() (*cli.Context, *flag.FlagSet) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := cli.NewContext(cli.NewApp(), fs, nil)

	return c, fs
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, ctx ...interface{}) {}

func (m *MockLogger) Info(msg string, ctx ...interface{}) {}

func (m *MockLogger) Warn(msg string, ctx ...interface{}) {}

func (m *MockLogger) Error(msg string, ctx ...interface{}) {}

func (m *MockLogger) Crit(msg string, ctx ...interface{}) {}

func (m *MockLogger) Close() error {
	args := m.Called()
	return args.Error(0)
}

type MockStats struct {
	mock.Mock
}

func (m *MockStats) Inc(name string, value int64, rate float32, tags ...interface{}) {}

func (m *MockStats) Gauge(name string, value float64, rate float32, tags ...interface{}) {}

func (m *MockStats) Timing(name string, value time.Duration, rate float32, tags ...interface{}) {}

func (m *MockStats) Close() error {
	args := m.Called()
	return args.Error(0)
}
