package tools

import (
	"errors"
	"testing"
)

// --- ParseExitArgs tests ---

func TestParseExitArgs_Success(t *testing.T) {
	args := map[string]any{
		"status": "success",
	}
	req, err := ParseExitArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Status != StatusSuccess {
		t.Errorf("Status = %q, want %q", req.Status, StatusSuccess)
	}
	if req.Stderr != "" {
		t.Errorf("Stderr = %q, want empty", req.Stderr)
	}
}

func TestParseExitArgs_Failure(t *testing.T) {
	args := map[string]any{
		"status": "failure",
		"stderr": "something broke",
	}
	req, err := ParseExitArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Status != StatusFailure {
		t.Errorf("Status = %q, want %q", req.Status, StatusFailure)
	}
	if req.ExitCode() != 1 {
		t.Errorf("ExitCode() = %d, want 1", req.ExitCode())
	}
}

func TestParseExitArgs_ProgressRejected(t *testing.T) {
	args := map[string]any{
		"status": "progress",
		"stderr": "context window running low",
	}
	_, err := ParseExitArgs(args)
	if err == nil {
		t.Fatal("expected error for invalid status 'progress', got nil")
	}
}

func TestParseExitArgs_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
	}{
		{"MissingStatus", map[string]any{"stderr": "error"}},
		{"InvalidStatusType", map[string]any{"status": true}},
		{"InvalidStatusValue", map[string]any{"status": "partial"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseExitArgs(tt.args)
			if err == nil {
				t.Fatalf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

// --- ExitCode tests ---

func TestExitCode(t *testing.T) {
	tests := []struct {
		status string
		want   int
	}{
		{StatusSuccess, 0},
		{StatusFailure, 1},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			req := ExitRequest{Status: tt.status, Stderr: "reason"}
			if got := req.ExitCode(); got != tt.want {
				t.Errorf("ExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

// --- Validate tests ---

func TestValidate_SuccessSilent_OK(t *testing.T) {
	req := ExitRequest{Status: StatusSuccess}
	if err := req.Validate(); err != nil {
		t.Errorf("success should be valid, got: %v", err)
	}
}

func TestValidate_SuccessWithStderr_Rejected(t *testing.T) {
	req := ExitRequest{Status: StatusSuccess, Stderr: "some error"}
	err := req.Validate()
	if err == nil {
		t.Fatal("success with stderr should be rejected, got nil")
	}
	if !errors.Is(err, ErrSuccessWithStderr) {
		t.Errorf("expected ErrSuccessWithStderr, got: %v", err)
	}
}

func TestValidate_FailureWithReason_OK(t *testing.T) {
	req := ExitRequest{Status: StatusFailure, Stderr: "file not found"}
	if err := req.Validate(); err != nil {
		t.Errorf("failure with reason should be valid, got: %v", err)
	}
}

func TestValidate_FailureWithoutReason_Rejected(t *testing.T) {
	req := ExitRequest{Status: StatusFailure}
	err := req.Validate()
	if err == nil {
		t.Fatal("failure without reason should be rejected, got nil")
	}
	if !errors.Is(err, ErrFailureWithoutReason) {
		t.Errorf("expected ErrFailureWithoutReason, got: %v", err)
	}
}
