package tape

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Writer handles JSONL persistence of a Tape.
// It uses open-write-close semantics: the file is only held open during
// each write operation, allowing external tools (tail, panopticon) to
// read the file freely between writes.
type Writer struct {
	dir       string // tape directory path
	sessionID string // session ID for filename
	path      string // full path to the JSONL file
	yamlPath  string // full path to the YAML mirror file
}

// NewWriter creates a Writer. It creates the directory if needed.
// The file is opened only during writes (open-write-close pattern).
func NewWriter(dir string, sessionID string) (*Writer, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating tape dir %q: %w", dir, err)
	}

	return &Writer{
		dir:       dir,
		sessionID: sessionID,
		path:      filepath.Join(dir, sessionID+".jsonl"),
		yamlPath:  filepath.Join(dir, sessionID+".log.yaml"),
	}, nil
}

// WriteEntry appends a TapeEntry as a JSON line to the file.
// Uses open-write-sync-close pattern so the file is not held open.
func (w *Writer) WriteEntry(entry TapeEntry) error {
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshalling tape entry: %w", err)
	}
	line = append(line, '\n')

	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("opening tape file %q: %w", w.path, err)
	}
	defer f.Close()

	if _, err := f.Write(line); err != nil {
		return fmt.Errorf("writing tape entry: %w", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("syncing tape file: %w", err)
	}

	if err := w.writeYAMLEntry(entry); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeYAMLEntry(entry TapeEntry) error {
	var data any
	if len(entry.Data) > 0 {
		if err := json.Unmarshal(entry.Data, &data); err != nil {
			data = map[string]any{
				"_raw":   string(entry.Data),
				"_error": fmt.Sprintf("invalid entry data JSON: %v", err),
			}
		}
	}

	doc := map[string]any{
		"type": entry.Type,
		"data": data,
	}
	yamlBytes, err := yaml.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshalling tape YAML entry: %w", err)
	}

	f, err := os.OpenFile(w.yamlPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("opening tape YAML file %q: %w", w.yamlPath, err)
	}
	defer f.Close()

	if _, err := f.Write([]byte("---\n")); err != nil {
		return fmt.Errorf("writing tape YAML document marker: %w", err)
	}
	if _, err := f.Write(yamlBytes); err != nil {
		return fmt.Errorf("writing tape YAML entry: %w", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("syncing tape YAML file: %w", err)
	}

	return nil
}

// Close is a no-op since we use open-write-close semantics.
// Kept for interface compatibility.
func (w *Writer) Close() error {
	return nil
}
