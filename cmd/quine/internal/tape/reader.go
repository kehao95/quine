package tape

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// TapeSummary holds the parsed header and current state of a tape file.
// It is the read-side counterpart to the write-side Writer.
type TapeSummary struct {
	SessionID       string          `json:"session_id"`
	ParentSessionID string          `json:"parent_session_id"`
	Depth           int             `json:"depth"`
	ModelID         string          `json:"model_id"`
	CreatedAt       int64           `json:"created_at"`
	Entries         []TapeEntry     `json:"entries"`
	Outcome         *SessionOutcome `json:"outcome,omitempty"`
}

// ReadTapeFile reads and parses a complete JSONL tape file from disk.
// It returns a TapeSummary with the header metadata, all entries, and
// the outcome (if the session has terminated).
func ReadTapeFile(path string) (*TapeSummary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening tape file %q: %w", path, err)
	}
	defer f.Close()

	return ReadTape(f)
}

// ReadTape parses a JSONL tape stream. It expects the first line to be a
// "meta" entry. Subsequent lines are collected as entries. If an "outcome"
// entry is found, it is stored separately in TapeSummary.Outcome.
func ReadTape(r io.Reader) (*TapeSummary, error) {
	scanner := bufio.NewScanner(r)
	// Allow large lines (tool results can be big).
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	summary := &TapeSummary{}

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry TapeEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("line %d: unmarshal entry: %w", lineNum, err)
		}

		switch entry.Type {
		case "meta":
			if summary.SessionID != "" {
				// Duplicate meta entry — skip it. This can happen when
				// multiple processes with the same session ID write to the
				// same tape file (a bug that has since been fixed, but we
				// handle it gracefully for backwards compatibility).
				continue
			}
			var meta struct {
				SessionID       string `json:"session_id"`
				ParentSessionID string `json:"parent_session_id"`
				Depth           int    `json:"depth"`
				ModelID         string `json:"model_id"`
				CreatedAt       int64  `json:"created_at"`
			}
			if err := json.Unmarshal(entry.Data, &meta); err != nil {
				return nil, fmt.Errorf("line %d: unmarshal meta: %w", lineNum, err)
			}
			summary.SessionID = meta.SessionID
			summary.ParentSessionID = meta.ParentSessionID
			summary.Depth = meta.Depth
			summary.ModelID = meta.ModelID
			summary.CreatedAt = meta.CreatedAt

		case "outcome":
			var outcome SessionOutcome
			if err := json.Unmarshal(entry.Data, &outcome); err != nil {
				return nil, fmt.Errorf("line %d: unmarshal outcome: %w", lineNum, err)
			}
			summary.Outcome = &outcome
		}

		summary.Entries = append(summary.Entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning tape: %w", err)
	}

	if lineNum == 0 {
		return nil, fmt.Errorf("empty tape file")
	}

	return summary, nil
}

// TailLastEntry reads the last complete line of a tape file and returns
// the parsed TapeEntry. This is used by the watcher to efficiently determine
// the current state of a tape without reading the entire file.
func TailLastEntry(path string) (*TapeEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening tape file %q: %w", path, err)
	}
	defer f.Close()

	// Seek backwards from end to find the last newline-terminated line.
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat tape file: %w", err)
	}

	size := info.Size()
	if size == 0 {
		return nil, fmt.Errorf("empty tape file")
	}

	// Read backwards to find the start of the last line.
	buf := make([]byte, 0, 4096)
	offset := size
	foundNewline := false

	for offset > 0 {
		readSize := int64(4096)
		if readSize > offset {
			readSize = offset
		}
		offset -= readSize

		chunk := make([]byte, readSize)
		if _, err := f.ReadAt(chunk, offset); err != nil {
			return nil, fmt.Errorf("reading chunk: %w", err)
		}

		// Prepend to buf
		buf = append(chunk, buf...)

		// Find last newline (skipping trailing newline)
		for i := len(buf) - 2; i >= 0; i-- {
			if buf[i] == '\n' {
				buf = buf[i+1:]
				foundNewline = true
				break
			}
		}
		if foundNewline {
			break
		}
	}

	// Trim trailing newline
	for len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}

	if len(buf) == 0 {
		return nil, fmt.Errorf("no valid line found")
	}

	var entry TapeEntry
	if err := json.Unmarshal(buf, &entry); err != nil {
		return nil, fmt.Errorf("unmarshal last entry: %w", err)
	}
	return &entry, nil
}
