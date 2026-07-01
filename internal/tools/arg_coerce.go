package tools

import (
	"fmt"
	"strings"
)

// Tool arguments arrive as a decoded map[string]any whose values carry whatever
// JSON shape the model emitted. A bare type assertion (`v, _ := args[k].(T)`)
// silently turns a missing key OR a wrong-typed value into the zero value, which
// downstream code then treats as an intentional choice — a fallback that masks
// a malformed tool call. The helpers below make that distinction explicit:
// callers learn whether an argument was absent (apply a default) or present but
// uncoercible (reject the call) instead of laundering both into a zero.

// boolFromAny coerces a decoded JSON value into a bool. It accepts a native
// bool, the JSON-number spellings 1/0, and the canonical string spellings
// ("true"/"false"/"1"/"0"). Anything else is an error so the caller can reject a
// present-but-unparseable argument rather than defaulting to false.
func boolFromAny(v any) (bool, error) {
	switch x := v.(type) {
	case bool:
		return x, nil
	case float64:
		return numericBool(x)
	case float32:
		return numericBool(float64(x))
	case int:
		return numericBool(float64(x))
	case int64:
		return numericBool(float64(x))
	case string:
		switch strings.ToLower(strings.TrimSpace(x)) {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		}
		return false, fmt.Errorf("not a boolean: %q", x)
	default:
		return false, fmt.Errorf("unsupported boolean type %T", v)
	}
}

func numericBool(x float64) (bool, error) {
	switch x {
	case 1:
		return true, nil
	case 0:
		return false, nil
	}
	return false, fmt.Errorf("not a boolean: %v", x)
}

// IntArg extracts an integer-valued argument. present is false when the key is
// absent (apply the caller's default); err is non-nil when the key is present
// but not coercible to an int (reject rather than substitute a zero).
func IntArg(args map[string]any, key string) (val int, present bool, err error) {
	raw, ok := args[key]
	if !ok {
		return 0, false, nil
	}
	v, err := intFromAny(raw)
	if err != nil {
		return 0, true, err
	}
	return v, true, nil
}

// BoolArg extracts a boolean-valued argument with the same absent/present/err
// contract as IntArg.
func BoolArg(args map[string]any, key string) (val bool, present bool, err error) {
	raw, ok := args[key]
	if !ok {
		return false, false, nil
	}
	v, err := boolFromAny(raw)
	if err != nil {
		return false, true, err
	}
	return v, true, nil
}

// StringArg extracts a string-valued argument with the same absent/present/err
// contract as IntArg. A present-but-non-string value is an error rather than an
// empty string.
func StringArg(args map[string]any, key string) (val string, present bool, err error) {
	raw, ok := args[key]
	if !ok {
		return "", false, nil
	}
	s, ok := raw.(string)
	if !ok {
		return "", true, fmt.Errorf("expected string, got %T", raw)
	}
	return s, true, nil
}
