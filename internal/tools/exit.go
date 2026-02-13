package tools

import (
	"errors"
	"fmt"
)

// Exit status values.
const (
	StatusSuccess  = "success"
	StatusFailure  = "failure"
	StatusProgress = "progress"
)

// Validation errors for exit tool.
var (
	// ErrSuccessWithStderr is returned when status="success" but stderr is set.
	// Success means the task is done — there's nothing to explain.
	ErrSuccessWithStderr = errors.New(
		"success exit must not include stderr — stderr is only for failure/progress. " +
			"If the task isn't fully done, use status=\"progress\"")

	// ErrFailureWithoutReason is returned when status="failure" but stderr is empty.
	ErrFailureWithoutReason = errors.New(
		"failure exit must include a reason in stderr")

	// ErrProgressWithoutReason is returned when status="progress" but stderr is empty.
	// Progress means "not done yet" — you must explain why.
	ErrProgressWithoutReason = errors.New(
		"progress exit must include stderr explaining why the task is incomplete")
)

// ExitRequest represents the parsed arguments from an exit tool call.
// Note: stdout is NOT included here. All stdout output should be written
// via sh commands (e.g., echo "result" > /dev/stdout). This separation
// ensures binary output is not polluted by text from exit.
type ExitRequest struct {
	Status string // "success", "failure", or "progress"
	Stderr string
}

// ExitCode maps the three-state status to a POSIX exit code:
//   - success  → 0
//   - failure  → 1
//   - progress → 2
func (r ExitRequest) ExitCode() int {
	switch r.Status {
	case StatusSuccess:
		return 0
	case StatusProgress:
		return 2
	default:
		return 1
	}
}

// Validate checks semantic constraints on the exit request:
//
//	success:  stderr forbidden
//	failure:  stderr required
//	progress: stderr required
func (r ExitRequest) Validate() error {
	switch r.Status {
	case StatusSuccess:
		if r.Stderr != "" {
			return ErrSuccessWithStderr
		}
	case StatusFailure:
		if r.Stderr == "" {
			return ErrFailureWithoutReason
		}
	case StatusProgress:
		if r.Stderr == "" {
			return ErrProgressWithoutReason
		}
	}
	return nil
}

// ParseExitArgs extracts ExitRequest from a ToolCall's Arguments map.
// Arguments: status (string enum, required), stderr (string, optional).
func ParseExitArgs(args map[string]any) (ExitRequest, error) {
	raw, ok := args["status"]
	if !ok {
		return ExitRequest{}, fmt.Errorf("missing required argument: status")
	}

	status, ok := raw.(string)
	if !ok {
		return ExitRequest{}, fmt.Errorf("status must be a string, got %T", raw)
	}

	switch status {
	case StatusSuccess, StatusFailure, StatusProgress:
		// valid
	default:
		return ExitRequest{}, fmt.Errorf("status must be one of %q, %q, %q; got %q",
			StatusSuccess, StatusFailure, StatusProgress, status)
	}

	req := ExitRequest{Status: status}

	if v, ok := args["stderr"]; ok {
		s, ok := v.(string)
		if !ok {
			return ExitRequest{}, fmt.Errorf("stderr must be a string, got %T", v)
		}
		req.Stderr = s
	}

	return req, nil
}
