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

func TestParseExitArgs_Progress(t *testing.T) {
	args := map[string]any{
		"status": "progress",
		"stderr": "context window running low",
	}
	req, err := ParseExitArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Status != StatusProgress {
		t.Errorf("Status = %q, want %q", req.Status, StatusProgress)
	}
	if req.Stderr != "context window running low" {
		t.Errorf("Stderr = %q, want %q", req.Stderr, "context window running low")
	}
}

func TestParseExitArgs_StatusOnly(t *testing.T) {
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

func TestParseExitArgs_MissingStatus(t *testing.T) {
	args := map[string]any{
		"stderr": "error",
	}
	_, err := ParseExitArgs(args)
	if err == nil {
		t.Fatal("expected error for missing status, got nil")
	}
}

func TestParseExitArgs_InvalidStatusType(t *testing.T) {
	args := map[string]any{
		"status": true, // bool, not string
	}
	_, err := ParseExitArgs(args)
	if err == nil {
		t.Fatal("expected error for invalid status type, got nil")
	}
}

func TestParseExitArgs_InvalidStatusValue(t *testing.T) {
	args := map[string]any{
		"status": "partial", // not in enum
	}
	_, err := ParseExitArgs(args)
	if err == nil {
		t.Fatal("expected error for invalid status value, got nil")
	}
}

// --- ExitCode tests ---

func TestExitCode_Success(t *testing.T) {
	req := ExitRequest{Status: StatusSuccess}
	if req.ExitCode() != 0 {
		t.Errorf("ExitCode() = %d, want 0", req.ExitCode())
	}
}

func TestExitCode_Failure(t *testing.T) {
	req := ExitRequest{Status: StatusFailure, Stderr: "reason"}
	if req.ExitCode() != 1 {
		t.Errorf("ExitCode() = %d, want 1", req.ExitCode())
	}
}

func TestExitCode_Progress(t *testing.T) {
	req := ExitRequest{Status: StatusProgress, Stderr: "reason"}
	if req.ExitCode() != 2 {
		t.Errorf("ExitCode() = %d, want 2", req.ExitCode())
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

func TestValidate_ProgressWithStderr_OK(t *testing.T) {
	req := ExitRequest{Status: StatusProgress, Stderr: "context window running low"}
	if err := req.Validate(); err != nil {
		t.Errorf("progress with stderr should be valid, got: %v", err)
	}
}

func TestValidate_ProgressWithoutReason_Rejected(t *testing.T) {
	req := ExitRequest{Status: StatusProgress}
	err := req.Validate()
	if err == nil {
		t.Fatal("progress without reason should be rejected, got nil")
	}
	if !errors.Is(err, ErrProgressWithoutReason) {
		t.Errorf("expected ErrProgressWithoutReason, got: %v", err)
	}
}
