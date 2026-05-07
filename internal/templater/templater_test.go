package templater

import (
	"errors"
	"testing"
)

func TestExpand_NoPlaceholders(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Expand(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestExpand_CurlyBraceSyntax(t *testing.T) {
	env := map[string]string{"BASE": "/app", "LOG_DIR": "${BASE}/logs"}
	out, err := Expand(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["LOG_DIR"] != "/app/logs" {
		t.Errorf("got %q, want %q", out["LOG_DIR"], "/app/logs")
	}
}

func TestExpand_BareSyntax(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "DSN": "postgres://$HOST/db"}
	out, err := Expand(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("got %q, want %q", out["DSN"], "postgres://localhost/db")
	}
}

func TestExpand_MissingRef_Lenient(t *testing.T) {
	env := map[string]string{"URL": "http://${MISSING_HOST}/path"}
	out, err := Expand(env, Options{Fallback: "unknown"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://unknown/path" {
		t.Errorf("got %q", out["URL"])
	}
}

func TestExpand_MissingRef_Strict(t *testing.T) {
	env := map[string]string{"URL": "http://${MISSING_HOST}/path"}
	_, err := Expand(env, Options{Strict: true})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var missingErr *ErrMissingVar
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected *ErrMissingVar, got %T", err)
	}
	if missingErr.Ref != "MISSING_HOST" {
		t.Errorf("expected ref MISSING_HOST, got %q", missingErr.Ref)
	}
}

func TestExpand_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "${A} world"}
	out, _ := Expand(env, Options{})
	if env["B"] != "${A} world" {
		t.Error("original map was mutated")
	}
	if out["B"] != "hello world" {
		t.Errorf("got %q", out["B"])
	}
}

func TestMustExpand_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	MustExpand(map[string]string{"X": "${UNDEF}"}, Options{Strict: true})
}

func TestExpand_EmptyMap(t *testing.T) {
	out, err := Expand(map[string]string{}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
