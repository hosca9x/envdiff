package trimmer

import (
	"testing"
)

func TestTrim_TrimAll_RemovesWhitespace(t *testing.T) {
	env := map[string]string{
		"KEY1": "  hello  ",
		"KEY2": "world",
		"KEY3": "\t spaced\t ",
	}
	out, results, err := Trim(env, TrimAll)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY1"] != "hello" {
		t.Errorf("KEY1: got %q, want %q", out["KEY1"], "hello")
	}
	if out["KEY2"] != "world" {
		t.Errorf("KEY2 should be unchanged")
	}
	if out["KEY3"] != "spaced" {
		t.Errorf("KEY3: got %q, want %q", out["KEY3"], "spaced")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 changed results, got %d", len(results))
	}
}

func TestTrim_TrimLeading_OnlyLeft(t *testing.T) {
	env := map[string]string{"K": "  value  "}
	out, _, err := Trim(env, TrimLeading)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "value  " {
		t.Errorf("got %q, want %q", out["K"], "value  ")
	}
}

func TestTrim_TrimTrailing_OnlyRight(t *testing.T) {
	env := map[string]string{"K": "  value  "}
	out, _, err := Trim(env, TrimTrailing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "  value" {
		t.Errorf("got %q, want %q", out["K"], "  value")
	}
}

func TestTrim_TrimInner_CollapsesSpaces(t *testing.T) {
	env := map[string]string{"K": "hello   world  foo"}
	out, _, err := Trim(env, TrimInner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "hello world foo" {
		t.Errorf("got %q, want %q", out["K"], "hello world foo")
	}
}

func TestTrim_NilMap_ReturnsError(t *testing.T) {
	_, _, err := Trim(nil, TrimAll)
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"K": "  v  "}
	_, _, _ = Trim(env, TrimAll)
	if env["K"] != "  v  " {
		t.Error("original map was mutated")
	}
}

func TestMustTrim_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil map")
		}
	}()
	MustTrim(nil, TrimAll)
}
