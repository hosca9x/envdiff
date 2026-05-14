package pinner

import (
	"errors"
	"testing"
)

func TestPin_AddsNewKeys(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080"}
	pins := map[string]string{"APP_ENV": "production"}

	out, err := Pin(env, pins, SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", out["APP_ENV"])
	}
	if out["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT=8080, got %q", out["APP_PORT"])
	}
}

func TestPin_SkipExisting_KeepsOriginal(t *testing.T) {
	env := map[string]string{"SECRET": "old"}
	pins := map[string]string{"SECRET": "new"}

	out, err := Pin(env, pins, SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != "old" {
		t.Errorf("expected old value to be kept, got %q", out["SECRET"])
	}
}

func TestPin_OverwriteExisting_ReplacesValue(t *testing.T) {
	env := map[string]string{"SECRET": "old"}
	pins := map[string]string{"SECRET": "new"}

	out, err := Pin(env, pins, OverwriteExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != "new" {
		t.Errorf("expected new value, got %q", out["SECRET"])
	}
}

func TestPin_ErrorOnConflict_ReturnsError(t *testing.T) {
	env := map[string]string{"SECRET": "old"}
	pins := map[string]string{"SECRET": "new"}

	_, err := Pin(env, pins, ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestPin_ErrorOnConflict_NoConflict(t *testing.T) {
	env := map[string]string{"SECRET": "same"}
	pins := map[string]string{"SECRET": "same"}

	_, err := Pin(env, pins, ErrorOnConflict)
	if err != nil {
		t.Fatalf("expected no error for identical value, got %v", err)
	}
}

func TestPin_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"KEY": "original"}
	pins := map[string]string{"KEY": "pinned"}

	_, _ = Pin(env, pins, OverwriteExisting)
	if env["KEY"] != "original" {
		t.Errorf("original map was mutated")
	}
}

func TestPin_NilEnv_ReturnsError(t *testing.T) {
	_, err := Pin(nil, map[string]string{}, SkipExisting)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestPinnedKeys_ReturnsConflictingKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	pins := map[string]string{"A": "99", "B": "2", "C": "77"}

	keys := PinnedKeys(env, pins)
	if len(keys) != 2 {
		t.Fatalf("expected 2 conflicting keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "A" || keys[1] != "C" {
		t.Errorf("expected [A C], got %v", keys)
	}
}
