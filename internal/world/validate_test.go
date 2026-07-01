package world

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluateResultsFileAccepted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.txt")
	if err := os.WriteFile(path, []byte("c01: alpha\nc02: beta\n"), 0o644); err != nil {
		t.Fatalf("write results: %v", err)
	}

	record, err := EvaluateResultsFile(path, map[string]string{"c01": "alpha", "c02": "beta"})
	if err != nil {
		t.Fatalf("EvaluateResultsFile returned error: %v", err)
	}
	if !record.Accepted {
		t.Fatalf("Accepted = false, want true (message=%q)", record.Message)
	}
	if got, want := record.Message, "validate accepted: all 2 cells correct"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
}

func TestEvaluateResultsFileMalformedLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.txt")
	if err := os.WriteFile(path, []byte("c01: alpha\nbad-line\n"), 0o644); err != nil {
		t.Fatalf("write results: %v", err)
	}

	record, err := EvaluateResultsFile(path, map[string]string{"c01": "alpha", "c02": "beta"})
	if err != nil {
		t.Fatalf("EvaluateResultsFile returned error: %v", err)
	}
	if record.Accepted {
		t.Fatalf("Accepted = true, want false")
	}
	if got, want := record.Message, "validate rejected: malformed line 2"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
}

func TestEvaluateResultsFileMissingAndIncorrect(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.txt")
	if err := os.WriteFile(path, []byte("c01: wrong\n"), 0o644); err != nil {
		t.Fatalf("write results: %v", err)
	}

	record, err := EvaluateResultsFile(path, map[string]string{"c01": "alpha", "c02": "beta"})
	if err != nil {
		t.Fatalf("EvaluateResultsFile returned error: %v", err)
	}
	if record.Accepted {
		t.Fatalf("Accepted = true, want false")
	}
	want := "validate rejected: data incomplete; missing c02; data mismatch; incorrect c01"
	if got := record.Message; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
}

func TestEvaluateResultsFileIncorrectOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.txt")
	if err := os.WriteFile(path, []byte("c01: wrong\nc02: beta\n"), 0o644); err != nil {
		t.Fatalf("write results: %v", err)
	}

	record, err := EvaluateResultsFile(path, map[string]string{"c01": "alpha", "c02": "beta"})
	if err != nil {
		t.Fatalf("EvaluateResultsFile returned error: %v", err)
	}
	if record.Accepted {
		t.Fatalf("Accepted = true, want false")
	}
	want := "validate rejected: incorrect c01"
	if got := record.Message; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
}
