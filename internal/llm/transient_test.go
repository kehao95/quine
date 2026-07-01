package llm

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

type testNetError struct {
	timeout   bool
	temporary bool
}

func (e testNetError) Error() string   { return "net test error" }
func (e testNetError) Timeout() bool   { return e.timeout }
func (e testNetError) Temporary() bool { return e.temporary }

func TestIsTransientInferenceDecodeError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "missing responses completed event",
			err:  errors.New("parsing SSE stream: no response.completed event found in SSE stream"),
			want: true,
		},
		{
			name: "scanner error",
			err:  errors.New("parsing SSE stream: scanning SSE stream: unexpected EOF"),
			want: true,
		},
		{
			name: "terminal provider failure is not transient",
			err:  errors.New("parsing SSE stream: response.failed: insufficient_quota: quota exhausted"),
			want: false,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTransientInferenceDecodeError(tt.err); got != tt.want {
				t.Fatalf("isTransientInferenceDecodeError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTransientInferenceTransportError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "wrapped unexpected eof",
			err:  fmt.Errorf("reading response body: %w", io.ErrUnexpectedEOF),
			want: true,
		},
		{
			name: "timeout net error",
			err:  testNetError{timeout: true},
			want: true,
		},
		{
			name: "temporary net error",
			err:  testNetError{temporary: true},
			want: true,
		},
		{
			name: "substring broken pipe",
			err:  errors.New("write tcp: broken pipe"),
			want: true,
		},
		{
			name: "permanent provider error",
			err:  errors.New("insufficient quota"),
			want: false,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTransientInferenceTransportError(tt.err); got != tt.want {
				t.Fatalf("isTransientInferenceTransportError() = %v, want %v", got, tt.want)
			}
		})
	}
}
