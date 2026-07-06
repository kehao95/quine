package qcli

import (
	"path/filepath"
	"strings"
)

const humanAuthor = "human"

func FormatClientSignalPayload(endpoint ClientEndpoint, action ControlAction, payload string) string {
	payload = strings.TrimRight(payload, "\r\n")
	lines := []string{
		"[qcli-client]",
		"authority: " + humanAuthor,
		"ctl_action: " + string(action),
		"reply_ctl: " + filepath.Join(endpoint.ControlPath, string(ControlActionPost)),
		"reply_required: false",
		"",
		"message:",
		payload,
	}
	return strings.Join(lines, "\n")
}

func parseClientEnvelope(content string) (authority, action, body string, ok bool) {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "[qcli-client]") && !strings.HasPrefix(trimmed, "[qctl-client]") {
		return "", "", "", false
	}
	authority = humanAuthor
	action = string(ControlActionInject)
	header := content
	if i := strings.Index(content, "message:\n"); i >= 0 {
		header = content[:i]
		body = content[i+len("message:\n"):]
	} else if i := strings.Index(content, "message:"); i >= 0 {
		header = content[:i]
		body = strings.TrimSpace(content[i+len("message:"):])
	}
	for _, line := range strings.Split(header, "\n") {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		value = strings.TrimSpace(value)
		switch strings.TrimSpace(key) {
		case "authority":
			if value != "" {
				authority = value
			}
		case "ctl_action":
			if value != "" {
				action = value
			}
		}
	}
	return authority, action, strings.TrimRight(body, "\r\n"), true
}
