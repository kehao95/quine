package tools

import (
	"testing"
)

func TestAllToolSchemas_Count(t *testing.T) {
	schemas := AllToolSchemas()
	if len(schemas) != 4 {
		t.Fatalf("AllToolSchemas() returned %d schemas, want 4", len(schemas))
	}
}
