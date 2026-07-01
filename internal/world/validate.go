package world

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// ValidationRecord captures the outcome of validating a results file against a
// world spec.
type ValidationRecord struct {
	Accepted bool
	Message  string
}

// EventResult returns the event-log result label for validation.
func (r ValidationRecord) EventResult() string {
	if r.Accepted {
		return "accepted"
	}
	return "rejected"
}

// ExitCode returns the CLI exit code for validation.
func (r ValidationRecord) ExitCode() int {
	if r.Accepted {
		return 0
	}
	return 4
}

// EvaluateResultsFile checks a results file against the current world items.
func EvaluateResultsFile(resultsPath string, items map[string]string) (ValidationRecord, error) {
	f, err := os.Open(resultsPath)
	if err != nil {
		return ValidationRecord{}, fmt.Errorf("read results: %w", err)
	}
	defer f.Close()

	expected := make(map[string]string, len(items))
	for k, v := range items {
		expected[k] = v
	}

	found := make(map[string]string, len(items))
	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			return ValidationRecord{Message: fmt.Sprintf("validate rejected: malformed line %d", lineNo)}, nil
		}

		cell, value, ok := strings.Cut(line, ":")
		if !ok {
			return ValidationRecord{Message: fmt.Sprintf("validate rejected: malformed line %d", lineNo)}, nil
		}
		cell = strings.TrimSpace(cell)
		value = strings.TrimSpace(value)
		if cell == "" || value == "" {
			return ValidationRecord{Message: fmt.Sprintf("validate rejected: malformed line %d", lineNo)}, nil
		}
		if _, exists := expected[cell]; !exists {
			return ValidationRecord{Message: fmt.Sprintf("validate rejected: malformed line %d", lineNo)}, nil
		}
		found[cell] = value
	}
	if err := scanner.Err(); err != nil {
		return ValidationRecord{}, fmt.Errorf("read results: %w", err)
	}

	missing := make([]string, 0)
	incorrect := make([]string, 0)
	for _, cell := range SortedCellIDs(expected) {
		value, ok := found[cell]
		if !ok {
			missing = append(missing, cell)
			continue
		}
		if value != expected[cell] {
			incorrect = append(incorrect, cell)
		}
	}

	if len(missing) == 0 && len(incorrect) == 0 {
		return ValidationRecord{
			Accepted: true,
			Message:  fmt.Sprintf("validate accepted: all %d cells correct", len(expected)),
		}, nil
	}

	parts := make([]string, 0, 2)
	if len(missing) > 0 {
		slices.Sort(missing)
		parts = append(parts, fmt.Sprintf("data incomplete; missing %s", strings.Join(missing, ", ")))
	}
	if len(incorrect) > 0 {
		slices.Sort(incorrect)
		if len(missing) > 0 {
			parts = append(parts, fmt.Sprintf("data mismatch; incorrect %s", strings.Join(incorrect, ", ")))
		} else {
			parts = append(parts, fmt.Sprintf("incorrect %s", strings.Join(incorrect, ", ")))
		}
	}

	return ValidationRecord{
		Message: "validate rejected: " + strings.Join(parts, "; "),
	}, nil
}
