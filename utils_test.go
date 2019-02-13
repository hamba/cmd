package cmd_test

import (
	"errors"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/hamba/cmd"
	"github.com/stretchr/testify/assert"
)

func Test_SplitTags(t *testing.T) {
	tests := []struct {
		tags []string

		results []interface{}
		err     error
	}{
		{[]string{"a=b"}, []interface{}{"a", "b"}, nil},
		{[]string{"a=b", "c=d"}, []interface{}{"a", "b", "c", "d"}, nil},
		{[]string{"a"}, nil, errors.New("invalid tags string")},
		{[]string{"a=b", "c"}, nil, errors.New("invalid tags string")},
	}

	for _, tt := range tests {
		res, err := cmd.SplitTags(tt.tags, "=")

		assert.Equal(t, res, tt.results)
		assert.Equal(t, err, tt.err)
	}
}

func TestWaitForSignals(t *testing.T) {
	tests := []struct {
		signal syscall.Signal
	}{
		{signal: syscall.SIGINT},
		{signal: syscall.SIGTERM},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		go func() {
			wg.Add(1)
			defer wg.Done()

			select {
			case s := <-cmd.WaitForSignals():
				assert.Equal(t, tt.signal, s)
			case <-time.After(1 * time.Second):
				assert.Failf(t, "", "Timeout waiting for %v", tt.signal)
			}
		}()

		time.Sleep(time.Millisecond)

		syscall.Kill(syscall.Getpid(), tt.signal)

		wg.Wait()
	}
}
