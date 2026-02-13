package tools

import (
	"errors"
	"fmt"
)

// Exit status values.
const (
	StatusSuccess = "success"
	StatusFailure = "failure"
)

// Validation errors for exit tool.
var (
	// ErrSuccessWithStderr is returned when status="success" but stderr is set.
	// Success means the task is done — there's nothing to explain.
	ErrSuccessWithStderr = errors.New(
		"success exit must not include stderr — stderr is only for failure. " +
			"If the task failed, use status=\"failure\"")

	// ErrFailureWithoutReason is returned when status="failure" but stderr is empty.
	ErrFailureWithoutReason = errors.New(
		"failure exit must include a reason in stderr")
)

// ExitRequest represents the parsed arguments from an exit tool call.
// Note: stdout is NOT included here. All stdout output should be written
// via sh commands (e.g., echo "result" > /dev/stdout). This separation
// ensures binary output is not polluted by text from exit.
type ExitRequest struct {
	Status string // "success" or "failure"
	Stderr string
}

// ExitCode maps the status to a POSIX exit code:
//   - success → 0
//   - failure → 1
func (r ExitRequest) ExitCode() int {
	switch r.Status {
	case StatusSuccess:
		return 0
	default:
		return 1
	}
}

// Validate checks semantic constraints on the exit request:
//
//	success:  stderr forbidden
//	failure:  stderr required
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
	case StatusSuccess, StatusFailure:
		// valid
	default:
		return ExitRequest{}, fmt.Errorf("status must be one of %q, %q; got %q",
			StatusSuccess, StatusFailure, status)
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
