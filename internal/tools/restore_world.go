package tools

import (
	"fmt"
	"strings"
)

type RestoreWorldRequest struct {
	Revision string
}

func ParseRestoreWorldArgs(args map[string]any) (RestoreWorldRequest, error) {
	v, ok := args["revision"]
	if !ok {
		return RestoreWorldRequest{}, fmt.Errorf("revision is required")
	}

	revision, ok := v.(string)
	if !ok {
		return RestoreWorldRequest{}, fmt.Errorf("revision must be a string, got %T", v)
	}
	revision = strings.TrimSpace(revision)
	if revision == "" {
		return RestoreWorldRequest{}, fmt.Errorf("revision must not be empty")
	}
	if !strings.HasPrefix(revision, "wr") {
		return RestoreWorldRequest{}, fmt.Errorf("revision must start with \"wr\", got %q", revision)
	}

	return RestoreWorldRequest{Revision: revision}, nil
}
