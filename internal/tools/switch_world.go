package tools

import (
	"fmt"
	"strings"
)

const worldHandlePrefix = "world://"

type SwitchWorldRequest struct {
	Target string
}

func ParseSwitchWorldArgs(args map[string]any) (SwitchWorldRequest, error) {
	v, ok := args["target"]
	if !ok {
		return SwitchWorldRequest{}, fmt.Errorf("target is required")
	}

	target, ok := v.(string)
	if !ok {
		return SwitchWorldRequest{}, fmt.Errorf("target must be a string, got %T", v)
	}
	target = strings.TrimSpace(target)
	if target == "" {
		return SwitchWorldRequest{}, fmt.Errorf("target must not be empty")
	}
	if strings.HasPrefix(target, "wr") {
		return SwitchWorldRequest{Target: target}, nil
	}
	if _, _, ok := parseWorldHandle(target); ok {
		return SwitchWorldRequest{Target: target}, nil
	}
	return SwitchWorldRequest{}, fmt.Errorf("target must be a world revision like \"wr3\" or a world handle like %q", worldHandlePrefix+"session/revision")
}

func parseWorldHandle(target string) (session string, revision string, ok bool) {
	if !strings.HasPrefix(target, worldHandlePrefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(target, worldHandlePrefix)
	parts := strings.Split(rest, "/")
	if len(parts) != 2 {
		return "", "", false
	}
	session = strings.TrimSpace(parts[0])
	revision = strings.TrimSpace(parts[1])
	if session == "" || revision == "" || !strings.HasPrefix(revision, "wr") {
		return "", "", false
	}
	return session, revision, true
}

func buildWorldHandle(session string, revision string) string {
	if strings.TrimSpace(session) == "" || strings.TrimSpace(revision) == "" {
		return ""
	}
	return worldHandlePrefix + session + "/" + revision
}
