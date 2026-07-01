package protocol

import (
	"encoding/json"
	"strings"

	"github.com/kehao95/quine/internal/tape"
)

func marshalToolArguments(args map[string]any) string {
	args = sanitizeToolArguments(args)
	if len(args) == 0 {
		return "{}"
	}
	data, err := json.Marshal(args)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// sanitizeToolArguments strips the internal malformed-arguments sentinel before
// arguments are encoded to any provider. decodeToolArguments stores the sentinel
// (with the raw undecodable payload) so the runtime can reject the call; the
// rejected assistant turn still lives in the replayed context, so without this
// strip the sentinel key and its raw payload would be re-sent on the wire every
// turn. Returns the input unchanged when the sentinel is absent (the common case).
func sanitizeToolArguments(args map[string]any) map[string]any {
	if _, bad := args[tape.MalformedArgumentsKey]; !bad {
		return args
	}
	out := make(map[string]any, len(args))
	for k, v := range args {
		if k == tape.MalformedArgumentsKey {
			continue
		}
		out[k] = v
	}
	return out
}

// decodeToolArguments parses a provider's tool-call arguments string into the
// tape representation. Providers send arguments as an opaque JSON string; an
// empty string legitimately means "no arguments". The previous call sites did
// `_ = json.Unmarshal([]byte(raw), &args)`, which discarded the decode error so
// a truncated/corrupt payload (e.g. a mid-JSON cutoff at max_tokens) silently
// became an empty map and ran as a no-op. Here a non-empty-but-undecodable
// payload is preserved under tape.MalformedArgumentsKey so the runtime can
// reject the call instead of laundering the corruption into a fake success.
func decodeToolArguments(raw string) map[string]any {
	if strings.TrimSpace(raw) == "" {
		return map[string]any{}
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		return map[string]any{tape.MalformedArgumentsKey: raw}
	}
	if args == nil {
		return map[string]any{}
	}
	return args
}
