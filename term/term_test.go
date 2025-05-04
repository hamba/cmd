package term_test

import (
	"bytes"
	"testing"

	"github.com/hamba/cmd/v3/term"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasic(t *testing.T) {
	trm := term.Basic{}

	assert.Implements(t, (*term.Term)(nil), trm)
}

func TestBasic_Output(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		verbose bool
		want    string
	}{
		{
			name:    "non-verbose",
			msg:     "test string",
			verbose: false,
			want:    "",
		},
		{
			name:    "verbose",
			msg:     "test string",
			verbose: true,
			want:    "test string\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			trm := term.Basic{
				Writer:  &buf,
				Verbose: test.verbose,
			}

			trm.Output(test.msg)

			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestBasic_Info(t *testing.T) {
	var buf bytes.Buffer
	trm := term.Basic{
		Writer: &buf,
	}

	trm.Info("input string")

	assert.Equal(t, "input string\n", buf.String())
}

func TestBasic_Warning(t *testing.T) {
	var buf bytes.Buffer
	trm := term.Basic{
		Writer: &buf,
	}

	trm.Warning("input string")

	assert.Equal(t, "input string\n", buf.String())
}

func TestBasic_Error(t *testing.T) {
	var (
		buf    bytes.Buffer
		errBuf bytes.Buffer
	)
	trm := term.Basic{
		Writer:      &buf,
		ErrorWriter: &errBuf,
	}

	trm.Error("input string")

	assert.Empty(t, buf.String())
	assert.Equal(t, "input string\n", errBuf.String())
}

func TestBasic_ErrorFallsBackToWriter(t *testing.T) {
	var buf bytes.Buffer
	trm := term.Basic{
		Writer: &buf,
	}

	trm.Error("input string")

	assert.Equal(t, "input string\n", buf.String())
}

func TestPrefixed_Output(t *testing.T) {
	m := &mockTerm{}
	m.On("Output", "test: input string")

	trm := term.Prefixed{
		OutputPrefix: "test: ",
		Term:         m,
	}

	trm.Output("input string")

	m.AssertExpectations(t)
}

func TestPrefixed_Info(t *testing.T) {
	m := &mockTerm{}
	m.On("Info", "test: input string")

	trm := term.Prefixed{
		InfoPrefix: "test: ",
		Term:       m,
	}

	trm.Info("input string")

	m.AssertExpectations(t)
}

func TestPrefixed_Warning(t *testing.T) {
	m := &mockTerm{}
	m.On("Warning", "test: input string")

	trm := term.Prefixed{
		WarningPrefix: "test: ",
		Term:          m,
	}

	trm.Warning("input string")

	m.AssertExpectations(t)
}

func TestPrefixed_Error(t *testing.T) {
	m := &mockTerm{}
	m.On("Error", "test: input string")

	trm := term.Prefixed{
		ErrorPrefix: "test: ",
		Term:        m,
	}

	trm.Error("input string")

	m.AssertExpectations(t)
}

type mockTerm struct {
	mock.Mock
}

func (m *mockTerm) Output(msg string) {
	_ = m.Called(msg)
}

func (m *mockTerm) Info(msg string) {
	_ = m.Called(msg)
}

func (m *mockTerm) Warning(msg string) {
	_ = m.Called(msg)
}

func (m *mockTerm) Error(msg string) {
	_ = m.Called(msg)
}
