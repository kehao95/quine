package tools

import "testing"

func TestBoolFromAny(t *testing.T) {
	truthy := []any{true, float64(1), 1, int64(1), "true", "TRUE", " 1 ", "True"}
	for _, v := range truthy {
		got, err := boolFromAny(v)
		if err != nil || !got {
			t.Fatalf("boolFromAny(%#v) = (%v,%v), want (true,nil)", v, got, err)
		}
	}
	falsy := []any{false, float64(0), 0, "false", "0", "FALSE"}
	for _, v := range falsy {
		got, err := boolFromAny(v)
		if err != nil || got {
			t.Fatalf("boolFromAny(%#v) = (%v,%v), want (false,nil)", v, got, err)
		}
	}
	bad := []any{"yes", "maybe", float64(2), []any{}, map[string]any{}}
	for _, v := range bad {
		if _, err := boolFromAny(v); err == nil {
			t.Fatalf("boolFromAny(%#v) should error, got nil", v)
		}
	}
}

func TestBoolArgAbsentVsBad(t *testing.T) {
	// absent → present=false, no error, caller applies its default
	if v, present, err := BoolArg(map[string]any{}, "fold"); v || present || err != nil {
		t.Fatalf("absent fold = (%v,%v,%v), want (false,false,nil)", v, present, err)
	}
	// present-but-bad → present=true, error (reject, do not default)
	if _, present, err := BoolArg(map[string]any{"fold": "maybe"}, "fold"); !present || err == nil {
		t.Fatalf("bad fold must be present-with-error, got present=%v err=%v", present, err)
	}
	// present-and-stringified → coerced
	if v, present, err := BoolArg(map[string]any{"fold": "true"}, "fold"); !v || !present || err != nil {
		t.Fatalf("stringified fold = (%v,%v,%v), want (true,true,nil)", v, present, err)
	}
}

func TestIntArgAbsentVsBad(t *testing.T) {
	if v, present, err := IntArg(map[string]any{}, "timeout"); v != 0 || present || err != nil {
		t.Fatalf("absent timeout = (%v,%v,%v), want (0,false,nil)", v, present, err)
	}
	if v, present, err := IntArg(map[string]any{"timeout": "30"}, "timeout"); v != 30 || !present || err != nil {
		t.Fatalf("stringified timeout = (%v,%v,%v), want (30,true,nil)", v, present, err)
	}
	if _, present, err := IntArg(map[string]any{"timeout": "soon"}, "timeout"); !present || err == nil {
		t.Fatalf("bad timeout must be present-with-error, got present=%v err=%v", present, err)
	}
}

// TestParseMarkArgsFoldStringified pins A3: a stringified/numeric fold flag is
// coerced rather than silently dropped to false, and a non-boolean is rejected.
func TestParseMarkArgsFoldStringified(t *testing.T) {
	req, err := ParseMarkArgs(map[string]any{"resolution": "x", "fold": "true"})
	if err != nil || !req.Fold {
		t.Fatalf("fold=\"true\" should yield Fold=true, got (%+v,%v)", req, err)
	}
	if _, err := ParseMarkArgs(map[string]any{"resolution": "x", "fold": "maybe"}); err == nil {
		t.Fatal("non-boolean fold must be rejected")
	}
}
