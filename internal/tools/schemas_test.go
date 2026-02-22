package tools

import (
	"testing"
)

func TestAllToolSchemas_Count(t *testing.T) {
	schemas := AllToolSchemas()
	if len(schemas) != 5 {
		t.Fatalf("AllToolSchemas() returned %d schemas, want 5", len(schemas))
	}
}
